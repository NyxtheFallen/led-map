package ledmap

import "time"

//ColorFiller is something that mimics the behavior of an LED strip
type ColorFiller interface {
	FillSingle(int) error
	Fill([]int) error
	Set(int, int) error
	Render() error
}

//LedMap represents
type LedMap struct {
	leds            ColorFiller
	initialPause    int
	subsequentPause int
}

//New creates a new LedMap and returns it.
func New(filler ColorFiller, initialPause int, subsequentPause int) *LedMap {
	l := &LedMap{
		leds:            filler,
		initialPause:    initialPause,
		subsequentPause: subsequentPause,
	}
	return l
}

//StartMap initializes the map and starts the LED cycle.
func (l *LedMap) StartMap(colors [][][]int) {
	for i, colorStrandSet := range colors {
		//Set initial color in group, then wait the prescribed amount of time
		l.leds.Fill(colors[i][0])
		l.leds.Render()
		if i == 0 {
			time.Sleep(time.Duration(l.initialPause) * time.Second)
		} else {
			time.Sleep(time.Duration(l.subsequentPause) * time.Second)
		}
		for j := 1; j < len(colorStrandSet); j++ {
			l.leds.Fill(colorStrandSet[j])
			l.leds.Render()
			time.Sleep(time.Duration(25) * time.Millisecond)
			//May have to add a very short sleep here, depending on how fast the LEDS end up fading
		}
	}
}
