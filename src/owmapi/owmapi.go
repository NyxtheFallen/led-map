//Package owmapi provides Get methods for the Open Weather Map API, along with an implementation of
//lights.Lighter for each weather entry, allowing for easy compatibility with the main process.
package owmapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

//Weather represents simple information about the weather at a given time.
type Weather struct {
	Datetime      time.Time //Unix DT in the host's local tz
	Temp          float32
	Precipitation bool
}

//Forecast is a wrapper around a list of Weathers
type Forecast []Weather

//Type needed to unmarshal a weather JSON response
type weatherResponse struct {
	Weather []struct{ Id int }
	Main    struct{ Feels_like float32 }
	Dt      int
}

//Type needed to unmarshal a forecast JSON response
type forecastResponse struct {
	List []weatherResponse
}

const (
	weather  = "weather"
	forecast = "forecast"
)

var client = &http.Client{Timeout: 120 * time.Second}

//Get the current weather as a weather object
func getCurrentWeather(apiKey string, locationID string) (Weather, error) {
	resp, err := getOpenWeatherMapPayload(apiKey, locationID, weather)
	if err != nil {
		return Weather{}, err
	}
	defer resp.Body.Close()

	//Decode and unmarshal the response
	weather := &weatherResponse{}
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Weather{}, nil
	}
	err = json.Unmarshal(response, weather)
	if err != nil {
		return Weather{}, err
	}

	//Weather IDs less than 700 are all different kinds of precipitation
	precipitation := false
	if weather.Weather[0].Id < 700 {
		precipitation = true
	}
	currentWeather := Weather{
		Datetime:      time.Unix(int64(weather.Dt), 0),
		Temp:          weather.Main.Feels_like,
		Precipitation: precipitation,
	}

	return currentWeather, nil
}

//Get the current forecast as a Forecast object
func getForecast(apiKey string, locationID string) (Forecast, error) {
	resp, err := getOpenWeatherMapPayload(apiKey, locationID, forecast)
	if err != nil {
		return Forecast{}, err
	}
	defer resp.Body.Close()

	//Decode and unmarshal the response
	forecast := &forecastResponse{}
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Forecast{}, err
	}
	err = json.Unmarshal(response, forecast)

	//Take the unmarshaled forecastResponse and load it into a Forecast
	weatherList := make(Forecast, len(forecast.List))
	for i, prediction := range forecast.List {
		precipitation := false
		if prediction.Weather[0].Id < 700 {
			precipitation = true
		}
		weatherList[i] = Weather{
			Datetime:      time.Unix(int64(prediction.Dt), 0),
			Temp:          prediction.Main.Feels_like,
			Precipitation: precipitation,
		}
	}

	return weatherList, nil
}

func getOpenWeatherMapPayload(apiKey string, locationID string, requestType string) (*http.Response, error) {
	// Build the base URL
	base, err := url.Parse("http://api.openweathermap.org/data/2.5/" + requestType)
	if err != nil {
		return &http.Response{}, err
	}

	//Add the required parameters to a Values object, then add them to the base URL
	params := url.Values{}
	params.Add("appid", apiKey)
	params.Add("id", locationID)
	params.Add("units", "imperial")

	base.RawQuery = params.Encode()

	//Get our data as a http Response
	resp, err := client.Get(base.String())
	if err != nil {
		return &http.Response{}, err
	}

	return resp, nil
}
