package weather

import (
	"TGbot/config"
	"TGbot/errors"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type WeatherData struct {
	Data     Data     `json:"data"`
	Location Location `json:"location"`
}

type Data struct {
	Values Values `json:"values"`
}

type Location struct {
	Name string `json:"name"`
}

type Values struct {
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"windSpeed"`
	Humidity    float64 `json:"humidity"`
}

func Weather() (*WeatherData, error) {
	_, err := config.GetKey("")
	errors.CheckError(err)

	var apiRealtimeWeather = "https://api.tomorrow.io/v4/weather/realtime?location=kaunas&apikey=9JSN8g563Duyv4IBlexCrozX7iDzVM4N"

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
