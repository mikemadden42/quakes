package main

// https://tutorialedge.net/golang/consuming-restful-api-with-go/
// https://mholt.github.io/json-to-go/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type quakes struct {
	Type     string `json:"type"`
	Metadata struct {
		Generated int64  `json:"generated"`
		URL       string `json:"url"`
		Title     string `json:"title"`
		Status    int    `json:"status"`
		API       string `json:"api"`
		Count     int    `json:"count"`
	} `json:"metadata"`
	Features []struct {
		Type       string `json:"type"`
		Properties struct {
			Mag     float64     `json:"mag"`
			Place   string      `json:"place"`
			Time    int64       `json:"time"`
			Updated int64       `json:"updated"`
			Tz      int         `json:"tz"`
			URL     string      `json:"url"`
			Detail  string      `json:"detail"`
			Felt    interface{} `json:"felt"`
			Cdi     interface{} `json:"cdi"`
			Mmi     interface{} `json:"mmi"`
			Alert   interface{} `json:"alert"`
			Status  string      `json:"status"`
			Tsunami int         `json:"tsunami"`
			Sig     int         `json:"sig"`
			Net     string      `json:"net"`
			Code    string      `json:"code"`
			Ids     string      `json:"ids"`
			Sources string      `json:"sources"`
			Types   string      `json:"types"`
			Nst     interface{} `json:"nst"`
			Dmin    interface{} `json:"dmin"`
			Rms     float64     `json:"rms"`
			Gap     interface{} `json:"gap"`
			MagType string      `json:"magType"`
			Type    string      `json:"type"`
			Title   string      `json:"title"`
		} `json:"properties"`
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		ID string `json:"id"`
	} `json:"features"`
	Bbox []float64 `json:"bbox"`
}

func main() {
	response, err := http.Get("https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_day.geojson")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject quakes
	err = json.Unmarshal(responseData, &responseObject)
	checkErr(err)
	fmt.Println(responseObject.Metadata.Title)
	for _, quake := range responseObject.Features {
		fmt.Println(
			time.Unix(quake.Properties.Time/1000, 0),
			quake.Properties.Type,
			quake.Properties.Mag,
			quake.Properties.Place,
			quake.Geometry.Coordinates[2])
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
