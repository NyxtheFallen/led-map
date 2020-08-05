//Package ledstrip wraps rpi-ws281x-go to provide simple functions for working with LED strips.
//It provides convenient methods for setting the colors of all
//LEDs in the map at the same time.
//It assumes you're using a string of LEDs, not a matrix, so it is one-dimensional.
package ledstrip

import (
	"fmt"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

//LedStrip represents a string of LEDs
type LedStrip struct {
	leds     *ws2811.WS2811
	autoFill bool
}

//Init returns an initialized led map with the supplied number of LEDs and brightness.
//Brightness is an integer between 0 and 255, and cannot be changed after initialization.
//If you want fine control over brightness, initalize the map with a brightness of 255
//and control the brightness through the color variables you pass.
//AutoFill gives you control over the rendering of your lights. If you want to have control over
//when color changes you apply actually get pushed to the LED strip, set this to false.
//Example:
//	myMap := ledmap.Init(100, 255, false) //100 LEDs with full brightness, has to be manually rendered
func Init(ledCount int, brightness int, autoFill bool) (*LedStrip, error) {
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = ledCount
	dev, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		return &LedStrip{}, err
	}
	err = dev.Init()
	if err != nil {
		return &LedStrip{}, err
	}

	return &LedStrip{dev, autoFill}, nil
}

//FillSingle applies a color to all LEDs on the map.
func (l *LedStrip) FillSingle(color int) error {
	leds := l.leds.Leds(0)
	for i := 0; i < len(leds); i++ {
		leds[i] = uint32(color)
	}
	if l.autoFill {
		err := l.leds.Render()
		if err != nil {
			return err
		}
	}
	return nil
}

//Fill applies an array of colors to all LEDs on the map. The array of colors must be the same length as the LED strip.
func (l *LedStrip) Fill(colors []int) error {
	leds := l.leds.Leds(0)
	if len(colors) != len(leds) {
		return fmt.Errorf("mismatch between number of colors and number of LEDs. colors = %v, LEDs = %v", len(colors), len(leds))
	}

	for i := 0; i < len(leds); i++ {
		leds[i] = uint32(colors[i])
	}
	if l.autoFill {
		err := l.leds.Render()
		if err != nil {
			return err
		}
	}
	return nil
}

//Set sets a single LED's color.
func (l *LedStrip) Set(index int, color int) error {
	leds := l.leds.Leds(0)
	if index > len(leds) || index < 0 {
		return fmt.Errorf("index is out of bounds")
	}
	leds[index] = uint32(color)
	if l.autoFill {
		err := l.leds.Render()
		if err != nil {
			return err
		}
	}
	return nil
}

//Render pushes all pending color changes to the LED strip.
//If autoFill is true, there's no reason to use this.
func (l *LedStrip) Render() error {
	err := l.leds.Render()
	if err != nil {
		return err
	}
	return nil
}

//Deinit shuts down the LEDs and releases their memory.
func (l *LedStrip) Deinit() {
	l.leds.Fini()
}
