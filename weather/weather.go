package weather

import (
	"TGbot/config"
	"TGbot/errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WeatherData struct {
	Data     Data     `json:"data"`
	Location Location `json:"location"`
}

type Data struct {
	Values Values `json:"values"`
	Time   string `json:"time"`
}

type Location struct {
	Name string `json:"name"`
}

type Values struct {
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"windSpeed"`
	Humidity    float64 `json:"humidity"`
	CloudCover  float64 `json:"cloudCover"`
	Visibility  float64 `jsin:"visibility"`
}

func Weather(cityWeather string) (*WeatherData, error) {
	var apiRealtimeWeather = fmt.Sprintf("https://api.tomorrow.io/v4/weather/realtime?location=%v&apikey=%v", cityWeather, config.WEATHERKEY)

	weather, err := RealtimeWeather(apiRealtimeWeather)
	errors.CheckError(err)

	return *&weather, nil
}

func RealtimeWeather(apiWeather string) (*WeatherData, error) {
	resp, err := http.Get(apiWeather)
	errors.CheckError(err)

	body, err := ioutil.ReadAll(resp.Body)
	errors.CheckError(err)

	defer resp.Body.Close()

	var weatherData WeatherData
	er := json.Unmarshal(body, &weatherData)
	errors.CheckError(er)

	return &weatherData, nil
}
