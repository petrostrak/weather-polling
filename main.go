package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const endpoint = "https://api.open-meteo.com/v1/forecast"

type Data struct {
	Elevation float64        `json:"elevation"`
	Hourly    map[string]any `json:"hourly"`
}

func main() {
	data, err := getResults(52.52, 13.41)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data)
}

func getResults(lat, long float64) (*Data, error) {
	uri := fmt.Sprintf("%s?latitude=%.2f&longitude=%.2f&hourly=temperature_2m", endpoint, lat, long)

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data Data
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
