// Package main prints the past day's earthquakes from the USGS GeoJSON feed.
package main

// https://tutorialedge.net/golang/consuming-restful-api-with-go/
// https://mholt.github.io/json-to-go/

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const feedURL = "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_day.geojson"

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
	client := &http.Client{Timeout: 15 * time.Second}
	q, err := fetch(client, feedURL)
	if err != nil {
		log.Fatal(err)
	}
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
