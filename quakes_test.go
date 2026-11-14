package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func loadFixture(t *testing.T) *quakes {
	t.Helper()
	f, err := os.Open("testdata/feed.geojson")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close() //nolint:errcheck // read-only file
	q, err := decode(f)
	if err != nil {
		t.Fatal(err)
	}
	return q
}

func TestDecode(t *testing.T) {
	q := loadFixture(t)
	if q.Metadata.Title != "USGS All Earthquakes, Past Day" {
		t.Errorf("title = %q", q.Metadata.Title)
	}
	if len(q.Features) != 3 {
		t.Fatalf("got %d features, want 3", len(q.Features))
	}
	if q.Features[1].Properties.Mag != nil {
		t.Errorf("null mag decoded as %v, want nil", *q.Features[1].Properties.Mag)
	}
}

func TestDecodeInvalid(t *testing.T) {
	if _, err := decode(strings.NewReader("<html>oops</html>")); err == nil {
		t.Error("expected error for non-JSON input")
	}
}

func TestReport(t *testing.T) {
	time.Local = time.UTC
	var buf bytes.Buffer
	if err := report(&buf, loadFixture(t)); err != nil {
		t.Fatal(err)
	}
	want := []string{
		"USGS All Earthquakes, Past Day",
		"2023-11-14 22:13:20 UTC  earthquake   2.4  8 km W of Holtville, CA  11.3 km",
		"2023-11-14 22:14:20 UTC  quarry blast     ?  Somewhere  0.5 km",
		"2023-11-14 22:15:20 UTC  earthquake   1.0  No depth  ?",
	}
	got := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d:\n%s", len(got), len(want), buf.String())
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d:\n got %q\nwant %q", i, got[i], want[i])
		}
	}
}

func TestFetch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/fail" {
			http.Error(w, "unavailable", http.StatusServiceUnavailable)
			return
		}
		http.ServeFile(w, r, "testdata/feed.geojson")
	}))
	defer srv.Close()

	q, err := fetch(srv.Client(), srv.URL+"/feed")
	if err != nil {
		t.Fatal(err)
	}
	if len(q.Features) != 3 {
		t.Errorf("got %d features, want 3", len(q.Features))
	}

	if _, err := fetch(srv.Client(), srv.URL+"/fail"); err == nil || !strings.Contains(err.Error(), "503") {
		t.Errorf("err = %v, want 503 error", err)
	}
}

func TestFeedURL(t *testing.T) {
	for _, name := range []string{"all_day", "all_hour", "significant_week", "4.5_month", "2.5_day", "1.0_hour"} {
		got, err := feedURL(name)
		if err != nil {
			t.Errorf("feedURL(%q): %v", name, err)
			continue
		}
		if want := feedBase + name + ".geojson"; got != want {
			t.Errorf("feedURL(%q) = %q, want %q", name, got, want)
		}
	}
	for _, name := range []string{"", "all", "day", "all_year", "3.0_day", "all_day_extra", "../all_day"} {
		if _, err := feedURL(name); err == nil {
			t.Errorf("feedURL(%q): expected error", name)
		}
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		minMag float64
		want   []string
	}{
		{0, []string{"8 km W of Holtville, CA", "Somewhere", "No depth"}},
		{1.0, []string{"8 km W of Holtville, CA", "No depth"}},
		{2.0, []string{"8 km W of Holtville, CA"}},
		{5.0, nil},
	}
	for _, tt := range tests {
		q := loadFixture(t)
		filter(q, tt.minMag)
		var got []string
		for _, f := range q.Features {
			got = append(got, f.Properties.Place)
		}
		if strings.Join(got, "|") != strings.Join(tt.want, "|") {
			t.Errorf("filter(%v) = %q, want %q", tt.minMag, got, tt.want)
		}
	}
}
