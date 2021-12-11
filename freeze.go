package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Main struct {
	Temp      float64 `json:"temp"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	FeelsLike float64 `json:"feels_like"`
	Pressure  float64 `json:"pressure"`
	SeaLevel  float64 `json:"sea_level"`
	GrndLevel float64 `json:"grnd_level"`
	Humidity  int     `json:"humidity"`
}

type CurrentWeatherData struct {
	Main `json:"main"`
}

func main() {
	// we dont care if it fails, only used for dev
	godotenv.Load()

	apiKey := os.Getenv("OPEN_WEATHER_API_KEY")
	baseUrl := os.Getenv("OPEN_WEATHER_BASE_URL")
	location := os.Getenv("OPEN_WEATHER_LOCATION")

	temp := getTemp(apiKey, baseUrl, location)

	log.Printf("Min temperature is %g", temp)

	if temp > 1 {
		log.Print("Min temperature is above 1c, not sending message")
		return
	}

	webhook := os.Getenv("WEBHOOK_URL")

	sendMessage(webhook, temp)

	log.Println("done")
}

func getTemp(apiKey string, baseUrl string, location string) float64 {
	v := url.Values{}
	v.Set("q", location)
	v.Set("appid", apiKey)
	v.Set("units", "metric")
	url := baseUrl + "?" + v.Encode()

	log.Printf("Sending request to %s", url)

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalln("OpenWeather request failed with status code", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var data = new(CurrentWeatherData)
	err = json.Unmarshal([]byte(body), &data)

	if err != nil {
		log.Fatalln(err)
	}

	return data.Main.TempMin
}

func sendMessage(webhook string, temp float64) {
	requestBody, err := json.Marshal(map[string]string{
		"content": fmt.Sprintf("It's freezing tonight! With a low of %g", temp),
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Going to send message to %s", webhook)

	http.Post(webhook, "application/json", bytes.NewBuffer(requestBody))
}
