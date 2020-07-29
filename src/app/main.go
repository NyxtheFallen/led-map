package main

import (
	"fmt"
	"led-map/src/api"

)

func main() {
	fmt.Println(api.GetCurrentWeather("ac295b00b689d2da17af1dbb659d3ffa", "4076784"))
}
