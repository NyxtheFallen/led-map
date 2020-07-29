package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type Weather struct {
	Datetime      time.Time //Unix DT in the host's local tz
	Temp          float32
	Precipitation bool
}

//Type needed to unmarshal JSON response
type weatherResponse struct {
	weather  []struct{ id int }
	main     struct{ feels_like float32 }
	dt       int
	timezone int
}

type Forecast []Weather

var client = &http.Client{Timeout: 120 * time.Second}

// func Get(apiKey string) Forecast {

// }

func GetCurrentWeather(apiKey string, locationId string) (Weather, error) {
	// Build the base URL
	base, err := url.Parse("http://api.openweathermap.org/data/2.5/weather")
	if err != nil {
		return Weather{}, err
	}

	//Add the required parameters to a Values object, then add them to the base URL
	params := url.Values{}
	params.Add("appid", apiKey)
	params.Add("id", locationId)
	params.Add("units", "imperial")

	base.RawQuery = params.Encode()

	//Get our data as a http Response
	resp, err := client.Get(base.String())
	if err != nil {
		return Weather{}, err
	}
	defer resp.Body.Close()

	//Decode the few bits we actually need
	weather := &weatherResponse{}
	err = json.NewDecoder(resp.Body).Decode(weather)
	if err != nil {
		return Weather{}, err
	}

	//Weather IDs less than 700 are all different kinds of precipitation
	precipitation := false
	if weather.weather[0].id < 700 {
		precipitation = true
	}

	currentWeather := Weather{
		Datetime:      time.Unix(int64(weather.dt+weather.timezone), 0),
		Temp:          weather.main.feels_like,
		Precipitation: precipitation,
	}

	return currentWeather, nil
}
