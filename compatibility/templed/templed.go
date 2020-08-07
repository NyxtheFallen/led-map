package templed

import (
	"fmt"
	"led-map/utilities"
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

	//This hurts my brain. It works, but it hurts.
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
	tempTransitions, err := utilities.Linspace(fromTemp, toTemp, numSteps)
	if err != nil {
		return []int{}, err
	}
	for i, temp := range tempTransitions {
		colors[i], err = utilities.GetFahrenheitTempColor(temp)
		if err != nil {
			return []int{}, err
		}
	}
	return colors, nil
}
