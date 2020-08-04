package main

import (
	"fmt"
	owmapi "led-map/src/apis/owmapi"
	templed "led-map/src/compatibility/templed"
)

type test struct{}

func (t test) ListTemps() [][]float64 {
	temps := [][]float64{
		{10.0, 10.0, 10.0, 10.0},
		{23.4, 23.6, 67.5, 40.2},
	}
	return temps
}

func main() {
	forecast, err := owmapi.Get("", []string{"2172797"})
	if err != nil {
		panic(err)
	}
	colors, _ := templed.GetColors(forecast, 3)
	fmt.Println(colors)
}
