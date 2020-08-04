package templed

import (
	"fmt"
	"math"
)

//TempLister is anything that provides a list of lists of temperatures
type TempLister interface {
	ListTemps() [][]float64
}

//GetColors receives something that can provide a table of temperatures, and returns a 3-dimensional array of colors.
//For example, the array input [[12.5, 23.5, 32.0, 40.5], [23.4, 23.6, 67.5, 40.2]] will return
func GetColors(t TempLister, fadeSteps int) ([][][]int, error) {
	temps := t.ListTemps()
	if len(temps) == 0 {
		return [][][]int{}, fmt.Errorf("0 temps returned from TempLister")
	}
	forecastColors := make([][][]int, len(temps))
	for i := range forecastColors {
		forecastColors[i] = make([][]int, len(temps[0]))
	}
	for i, forecast := range temps {
		for j, temp := range forecast {
			if j == 0 {
				continue
			}
			colorFade, err := fade(forecast[j-1], temp, fadeSteps)
			if err != nil {
				return [][][]int{}, err
			}
			forecastColors[i][j-1] = colorFade
		}
		lastTemp := forecast[len(forecast)-1]
		colorFade, err := fade(lastTemp, lastTemp, fadeSteps)
		if err != nil {
			return [][][]int{}, err
		}
		forecastColors[i][len(forecast)-1] = colorFade
	}

	return reorderForecastColorHierarchy(forecastColors), nil
}

func reorderForecastColorHierarchy(forecastColors [][][]int) [][][]int {
	//At this point, we have a hierarchy: location > forecast > colors to fade between
	//To be most useful to our map, we need to set all locations, one color at a time, meaning that
	//we have to reorder this hierarchy to forecast > colors > locations. This is not easy. Buckle in.

	forecastLength := len(forecastColors[0]) //len(forecastColors[0]) is the number of weather instances at a single location, i.e. the length of the forecast
	fadeLength := len(forecastColors[0][0])  //len(forecastColors[0][0]) is the number of colors per weather instance
	numLocations := len(forecastColors)      //len(forecastColors) is the number of locations, which will coincide with the number of LEDs

	//The base array must be long enough to hold the forecasts
	unpivotedColors := make([][][]int, forecastLength)
	for i := range unpivotedColors {
		//The second-level array must be long enough to hold each fade-group
		unpivotedColors[i] = make([][]int, fadeLength)
		for j := range unpivotedColors[i] {
			//The third-level array must be long enough to hold one int for each location
			unpivotedColors[i][j] = make([]int, numLocations)
		}
	}

	for i := 0; i < forecastLength; i++ {
		for j := 0; j < fadeLength; j++ {
			for k := 0; k < numLocations; k++ {
				unpivotedColors[i][j][k] = forecastColors[k][i][j]
			}
		}
	}

	return unpivotedColors
}

//fade returns a list of color values representing a fade between two temperatures.
//numSteps indicates how many values to return in the list.
func fade(fromTemp, toTemp float64, numSteps int) ([]int, error) {
	colors := make([]int, numSteps)
	tempTransitions, err := linspace(fromTemp, toTemp, numSteps)
	if err != nil {
		return []int{}, err
	}
	for i, temp := range tempTransitions {
		colors[i], err = getTempColor(temp)
		if err != nil {
			return []int{}, err
		}
	}
	return colors, nil
}

//Find what color a a temperature should be. Zero degrees and colder are 100% blue,
//0-32 fade from blue to white, 32-70 fade from white to yellow, and from 70-110, fade from
//yellow to red.
//Math is done using HSV, because keeping brightness at 100% while adjusting hue is much easier this way.
func getTempColor(temp float64) (int, error) {
	if temp > 110.0 {
		temp = 110.0
	}
	if temp < 0.0 {
		temp = 0.0
	}

	//If below freezing, we want blue of varying saturations
	var tempHue float64
	var tempSat float64
	if temp < 32 {
		tempHue = (2.0 / 3.0) //this is blue
		//saturation needs to decrease inverse-exponentially from 1 to 0 as x approaches 32
		x := temp * 10.0 / 32.0 //force x to be between 0 and 10
		y := math.Pow(1.4, x-7.0) / math.Pow(1.4, 3.0)
		tempSat = 1.0 - y
	} else if temp < 70 { //If between freezing and 70, fade to a pleasant yellow
		tempHue = (1.0 / 6.0) //this is yellow
		//Saturation needs to increase exponentially from 0 to 1 as x approaches 70 from 32
		tempRange := 70.0 - 32.0
		tempValue := temp - 32.0
		x := tempValue * 10.0 / tempRange
		y := 1 - (math.Pow(1.5, -x-3.0) / math.Pow(1.5, -3.0))
		tempSat = y
	} else { //Between 70 and 110, fade from yellow to red
		tempHue = (1.0 / 6.0) - ((temp-70)/40.0)*(1.0/6.0)
		tempSat = 1.0
	}

	grb, err := hsvToGrb(tempHue, tempSat, 1.0)
	if err != nil {
		return 0, err
	}
	return grb, nil
}

//converts hsv to rgb. RGB is returned as a single int, in the format expected by our map.
//All input values should be between 0 and 1.
func hsvToGrb(h, s, v float64) (int, error) {
	if h < 0 || h > 1 || s < 0 || s > 1 || v < 0 || v > 1 {
		return 0, fmt.Errorf("all arguments to hsvToRgb should be in the range 0, 1, got: %v %v %v", h, s, v)
	}
	if s == 0 {
		//it's gray
		v *= 255
		intV := int(v)
		grb := intV<<16 + intV<<8 + intV
		return grb, nil
	}
	//Weird color wheel magic
	h *= 6
	if h == 6 {
		h = 0
	}
	i := math.Floor(h)
	//Set three values that will be assigned to g, r, and b depending on
	//where in the color wheel the hue falls
	//h-i represents the relative shade of a given section of the color wheel
	v1 := v * (1 - s)
	v2 := v * (1 - s*(h-i))
	v3 := v * (1 - s*(1-(h-i)))

	var g, r, b float64
	switch i {
	case 0:
		g = v3
		r = v
		b = v1
	case 1:
		g = v
		r = v2
		b = v1
	case 2:
		g = v
		r = v1
		b = v3
	case 3:
		g = v2
		r = v1
		b = v
	case 4:
		g = v1
		r = v3
		b = v
	default:
		g = v1
		r = v
		b = v2
	}
	g, r, b = g*255, r*255, b*255

	grb := int(g)<<16 + int(r)<<8 + int(b)
	return grb, nil
}

//return an evenly-distributed array of values between start and end, inclusive of start and end
//not guaranteed to be perfect, but guaranteed not to overflow start or end.
func linspace(start, end float64, steps int) ([]float64, error) {
	if steps < 0 {
		return []float64{}, fmt.Errorf("negative steps not allowed, got: %v", steps)
	}
	linspaceRange := end - start
	increment := linspaceRange / float64(steps-1)
	returnArray := make([]float64, steps)
	returnArray[0] = start
	for i := 1; i < steps-1; i++ {
		returnArray[i] = start + increment*float64(i)
	}
	returnArray[steps-1] = end
	return returnArray, nil
}
