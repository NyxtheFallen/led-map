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
func (l *LedMap) StartMap(apiKey string, locationIds []string, fadeSteps int) error {
	resp, err := l.api.Get(apiKey, locationIds)
	if err != nil {
		return err
	}
	colors, err := l.compatibility.GetColors(resp, fadeSteps)
	if err != nil {
		return err
	}
	for {
		for k := 0; k < len(colors[0]); k++ {
			for j := 0; j < fadeSteps; j++ {
				for i, locationColors := range colors {
					l.leds.Set(i, locationColors[k][j])
					if i == 0 {
						time.Sleep(15 * time.Second)
					} else {

						time.Sleep(5 * time.Second)
					}
				}
				l.leds.Render()
			}
		}
	}
}
