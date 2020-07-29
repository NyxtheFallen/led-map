/*
 * Package lights provides a compatibility layer between the API and the main process,
 * allowing the API to define how its data should be represented on the light strip.
 */
package lights

type Light struct {
	Color int //A color, eg. 0xffffff (white) or 0x000000 (black)
	//Define how to fade between two colors by providing a function of form func(startColor int, endColor int, numSteps int) []int,
	//where the return is a slice of len(numSteps) of colors between startColor and endColor.
	FadeBetween func(int, int, int) []int
	Blink       bool //Should the light blink?
	BlinkSpeed  int  //Time in milliseconds between blinks
}

type Lighter interface {
	ToLight() Light
}
