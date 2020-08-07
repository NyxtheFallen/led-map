package ledmap

import (
	"led-map/ledstrip"
)

//MapController is any function that takes a list of color values and uses them to set the map's LED colors.
type MapController func([][][]int, ColorFiller) error

//ColorFiller is something that mimics the behavior of an LED strip
type ColorFiller interface {
	FillSingle(int) error
	Fill([]int) error
	Set(int, int) error
	Render() error
}

//Function used to set options
type option func(*LedMap)

//LedMap represents
type LedMap struct {
	leds       ColorFiller
	colors     [][][]int
	controller MapController
}

//New creates a new LedMap and returns it.
func New() (*LedMap, error) {
	stripLength := 100
	strip, err := ledstrip.Init(stripLength, 255, false)
	if err != nil {
		return &LedMap{}, err
	}
	colors := make([][][]int, 0)
	l := &LedMap{
		leds:       strip,
		colors:     colors,
		controller: nil,
	}
	return l, nil
}

//Option applies an option function to the LedMap
func (l *LedMap) Option(opts ...option) {
	for _, opt := range opts {
		opt(l)
	}
}

//LEDs provides an option for setting the ColorFiller on the LedMap.
func LEDs(leds ColorFiller) option {
	return func(l *LedMap) {
		l.leds = leds
	}
}

//Colors provides an option for setting the array of colors on the LedMap.
func Colors(colors [][][]int) option {
	return func(l *LedMap) {
		l.colors = colors
	}
}

//Controller provides an option for setting the MapController function on the LedMap.
func Controller(controller MapController) option {
	return func(l *LedMap) {
		l.controller = controller
	}
}

//RunMapController runs the LedMap's controller on the LedMap's colors.
func (l *LedMap) RunMapController() error {
	err := l.controller(l.colors, l.leds)
	if err != nil {
		return err
	}
	return nil
}

// for i, colorStrandSet := range colors {
// 	//Set initial color in group, then wait the prescribed amount of time
// 	l.leds.Fill(colors[i][0])
// 	l.leds.Render()
// 	if i == 0 {
// 		time.Sleep(time.Duration(l.initialPause) * time.Second)
// 	} else {
// 		time.Sleep(time.Duration(l.subsequentPause) * time.Second)
// 	}
// 	for j := 1; j < len(colorStrandSet); j++ {
// 		l.leds.Fill(colorStrandSet[j])
// 		l.leds.Render()
// 		time.Sleep(time.Duration(25) * time.Millisecond)
// 		//May have to add a very short sleep here, depending on how fast the LEDS end up fading
// 	}
// }
