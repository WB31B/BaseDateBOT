package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const myKeyWeather = ""

// var apiWeather = fmt.Sprintf("https://api.tomorrow.io/v4/weather/forecast?location=kaunas&apikey=%v", myKeyWeather)
var apiRealtimeWeather = fmt.Sprintf("https://api.tomorrow.io/v4/weather/realtime?location=kaunas&apikey=%v", myKeyWeather)

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

	RealtimeWeather(apiRealtimeWeather)

}

func RealtimeWeather(apiWeather string) {
	resp, err := http.Get(apiWeather)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	defer resp.Body.Close()

	var weatherData WeatherData
	er := json.Unmarshal(body, &weatherData)
	if er != nil {
		panic(er.Error)
	}

	fmt.Printf("%+v\n", weatherData)
}
