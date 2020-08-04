package main

import (
	owmapi "led-map/src/apis/owmapi"
	templed "led-map/src/compatibility/templed"
	ledmap "led-map/src/ledmap"
	ledstrip "led-map/src/ledstrip"
)

func main() {
	forecast, err := owmapi.Get("", []string{"2172797"})
	if err != nil {
		panic(err)
	}
	colors, _ := templed.GetColors(forecast, 3)
	leds, err := ledstrip.Init(100, 255, false)
	if err != nil {
		panic(err)
	}
	weathermap := ledmap.New(leds, 10, 5)
	for {
		weathermap.StartMap(colors)
	}
}
