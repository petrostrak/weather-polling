package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	endpoint = "https://api.open-meteo.com/v1/forecast"
)

var (
	pollInterval = time.Second * 2
)

type Data struct {
	Elevation float64        `json:"elevation"`
	Hourly    map[string]any `json:"hourly"`
}

type Poller struct {
	closeChan chan struct{}
}

func NewPoller() *Poller {
	return &Poller{
		closeChan: make(chan struct{}),
	}
}

func (p *Poller) start() {
	fmt.Println("starting poller")

	ticker := time.NewTicker(pollInterval)
outer:
	for {
		select {
		case <-ticker.C:
			data, err := getResults(52.52, 13.41)
			if err != nil {
				log.Fatal(err)
			}
			err = p.handleData(data)
			if err != nil {
				log.Fatal(err)
			}
		case <-p.closeChan:
			break outer
		}
	}
	fmt.Println("poller stopped")
}

func (p *Poller) close() {
	p.closeChan <- struct{}{}
}

func (p *Poller) handleData(data *Data) error {
	fmt.Println(data)
	return nil
}

func main() {
	poller := NewPoller()
	go poller.start()

	time.Sleep(time.Second * 7)
	poller.close()
	select {}
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
