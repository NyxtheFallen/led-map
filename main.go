package main

import (
	"led-map/compatibility/templed"
	"led-map/datastore/owmapi"
	"led-map/ledmap"
	"led-map/ledstrip"
	"os"
)

const apiBasePath string = "/api"

func main() {
	theSamePlace := make([]string, 100)
	for i := 0; i < 100; i++ {
		theSamePlace[i] = "2172797"
	}
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("API key required for proper execution!")
	}
	forecast, err := owmapi.Get(apiKey, theSamePlace)
	if err != nil {
		panic(err)
	}
	colors, _ := templed.GetColors(forecast, 40)
	leds, err := ledstrip.Init(100, 255, false)
	if err != nil {
		panic(err)
	}
	defer leds.Deinit()
	weathermap := ledmap.New(leds, 10, 3)
	for {
		weathermap.StartMap(colors)
	}
}
