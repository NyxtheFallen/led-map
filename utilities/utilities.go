package utilities

import (
	"fmt"
	"math"
)

//GetFahrenheitTempColor finds what color a a temperature should be. Zero degrees and colder are 100% blue,
//0-32 fade from blue to white, 32-70 fade from white to yellow, and from 70-110, fade from
//yellow to red.
//Math is done using HSV, because keeping brightness at 100% while adjusting hue is much easier this way.
func GetFahrenheitTempColor(temp float64) (int, error) {
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

	grb, err := HsvToGrb(tempHue, tempSat, 1.0)
	if err != nil {
		return 0, err
	}
	return grb, nil
}

//HsvToGrb converts hsv to rgb. RGB is returned as a single int, in the format expected by our map.
//All input values should be between 0 and 1.
func HsvToGrb(h, s, v float64) (int, error) {
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

//Linspace return an evenly-distributed array of values between start and end, inclusive of start and end
//not guaranteed to be perfect, but guaranteed not to overflow start or end.
func Linspace(start, end float64, steps int) ([]float64, error) {
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
