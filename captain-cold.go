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

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

type Hourly struct {
	// we only want the first 32 results (8am the next day)
	Temperature [32]float64 `json:"temperature_2m"`
}

type WeatherData struct {
	Hourly `json:"hourly"`
}

const MIN_TEMP = 1
const NUM_HOURS = 11

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest() {
	defer log.Println("Finished running Captain Cold")

	// we dont care if it fails, only used for dev
	godotenv.Load()

	weatherUrl := os.Getenv("WEATHER_BASE_URL")
	if weatherUrl == "" {
		log.Fatalln("WEATHER_BASE_URL not set")
	}

	lat := os.Getenv("WEATHER_LAT")
	if lat == "" {
		log.Fatalln("WEATHER_LAT not set")
	}

	lng := os.Getenv("WEATHER_LNG")
	if lng == "" {
		log.Fatalln("WEATHER_LNG not set")
	}

	temps := getTemps(weatherUrl, lat, lng)

	lowerThanMin := isLowerThanMin(temps)

	if !lowerThanMin {
		log.Print("Min temperature is above 1c, not sending message")
		return
	}

	log.Println("Attempting to send discord message")

	webhook := os.Getenv("WEBHOOK_URL")

	sendMessage(webhook)
}

func getTemps(baseUrl string, lat string, lng string) *[NUM_HOURS]float64 {
	v := url.Values{}
	v.Set("latitude", lat)
	v.Set("longitude", lng)
	v.Set("hourly", "temperature_2m")
	url := baseUrl + "?" + v.Encode()

	log.Printf("Sending request to %s", url)

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalln("Weather request failed with status code", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var data = new(WeatherData)
	err = json.Unmarshal([]byte(body), &data)

	if err != nil {
		log.Fatalln(err)
	}

	// get the 15 times we care about 9pm - 8am
	tempSlice := data.Hourly.Temperature[21:32]

	// convert the slice to fixed array pointer so we aren't returning a slice and is more memory efficent as pointer
	temps := (*[11]float64)(tempSlice)

	log.Print("Temperatures returned from the API are", temps)

	return temps
}

func isLowerThanMin(temps *[NUM_HOURS]float64) bool {
	for i := 0; i < NUM_HOURS; i++ {
		if temps[i] <= MIN_TEMP {
			return true
		}
	}

	return false
}

func sendMessage(webhook string) {
	requestBody, err := json.Marshal(map[string]string{
		"content": fmt.Sprintf("It's freezing tonight! It will drop lower than %d before 8am tomorrow.", MIN_TEMP),
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Going to send message to %s", webhook)

	http.Post(webhook, "application/json", bytes.NewBuffer(requestBody))
}
