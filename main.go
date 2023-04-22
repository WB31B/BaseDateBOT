package main

import (
	"TGbot/bot"
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
}

type Location struct {
	Name string `json:"name"`
}

type Values struct {
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"windSpeed"`
	Humidity    float64 `json:"humidity"`
}

func main() {
	weatherKey, err := config.GetKey("")
	errors.CheckError(err)

	var apiRealtimeWeather = fmt.Sprintf("https://api.tomorrow.io/v4/weather/realtime?location=kaunas&apikey=%v", weatherKey)

	RealtimeWeather(apiRealtimeWeather)

	bot.StartBot()
}

func RealtimeWeather(apiWeather string) {
	resp, err := http.Get(apiWeather)
	errors.CheckError(err)

	body, err := ioutil.ReadAll(resp.Body)
	errors.CheckError(err)

	defer resp.Body.Close()

	var weatherData WeatherData
	er := json.Unmarshal(body, &weatherData)
	errors.CheckError(er)

	fmt.Printf("%+v\n", weatherData)
}
