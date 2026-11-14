// Package main prints recent earthquakes from the USGS GeoJSON summary feeds.
package main

// https://tutorialedge.net/golang/consuming-restful-api-with-go/
// https://mholt.github.io/json-to-go/

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const feedBase = "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/"

// feedURL returns the summary feed URL for a name like "all_day" or "2.5_week".
func feedURL(name string) (string, error) {
	level, period, ok := strings.Cut(name, "_")
	validLevel := map[string]bool{"all": true, "significant": true, "4.5": true, "2.5": true, "1.0": true}
	validPeriod := map[string]bool{"hour": true, "day": true, "week": true, "month": true}
	if !ok || !validLevel[level] || !validPeriod[period] {
		return "", fmt.Errorf("unknown feed %q: want {all,significant,4.5,2.5,1.0}_{hour,day,week,month}", name)
	}
	return feedBase + name + ".geojson", nil
}

type quakes struct {
	Metadata struct {
		Title string `json:"title"`
	} `json:"metadata"`
	Features []struct {
		Properties struct {
			Mag   *float64 `json:"mag"`
			Place string   `json:"place"`
			Time  int64    `json:"time"`
			Type  string   `json:"type"`
		} `json:"properties"`
		Geometry struct {
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

func main() {
	feed := flag.String("feed", "all_day", "USGS summary feed, e.g. all_hour, 2.5_day, significant_week")
	minMag := flag.Float64("min-mag", 0, "only show quakes at or above this magnitude")
	flag.Parse()

	url, err := feedURL(*feed)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{Timeout: 15 * time.Second}
	q, err := fetch(client, url)
	if err != nil {
		log.Fatal(err)
	}
	filter(q, *minMag)
	if err := report(os.Stdout, q); err != nil {
		log.Fatal(err)
	}
}

func fetch(client *http.Client, url string) (*quakes, error) {
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close() //nolint:errcheck // read-only body

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching %s: %s", url, response.Status)
	}
	return decode(response.Body)
}

func decode(r io.Reader) (*quakes, error) {
	var q quakes
	if err := json.NewDecoder(r).Decode(&q); err != nil {
		return nil, fmt.Errorf("decoding feed: %w", err)
	}
	return &q, nil
}

// filter drops quakes below minMag. A quake with no magnitude is kept only
// when minMag is 0.
func filter(q *quakes, minMag float64) {
	if minMag <= 0 {
		return
	}
	kept := q.Features[:0]
	for _, quake := range q.Features {
		if quake.Properties.Mag != nil && *quake.Properties.Mag >= minMag {
			kept = append(kept, quake)
		}
	}
	q.Features = kept
}

func report(w io.Writer, q *quakes) error {
	if _, err := fmt.Fprintln(w, q.Metadata.Title); err != nil {
		return err
	}
	for _, quake := range q.Features {
		mag := "?"
		if quake.Properties.Mag != nil {
			mag = fmt.Sprintf("%.1f", *quake.Properties.Mag)
		}
		depth := "?"
		if c := quake.Geometry.Coordinates; len(c) >= 3 {
			depth = fmt.Sprintf("%.1f km", c[2])
		}
		_, err := fmt.Fprintf(w, "%s  %s  %4s  %s  %s\n",
			time.UnixMilli(quake.Properties.Time).Format("2006-01-02 15:04:05 MST"),
			quake.Properties.Type,
			mag,
			quake.Properties.Place,
			depth)
		if err != nil {
			return err
		}
	}
	return nil
}
