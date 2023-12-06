package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type apiConfigData struct {
	OpenWeatherApiKey string `json:"OpenWeatherApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kalvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := os.ReadFile(filename)

	// CHECKPOINT
	fmt.Println("This is Bytes Data: ", bytes)

	if err != nil {
		// CHECKPOINT
		fmt.Println(err)
		return apiConfigData{}, err
	}

	var configData apiConfigData

	err1 := json.Unmarshal(bytes, &configData)
	if err1 != nil {
		// CHECKPOINT
		fmt.Println(err1)
		return apiConfigData{}, err1
	}
	// CHECKPOINT
	fmt.Println("This is Config Data: ", configData)
	return configData, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w,"Hello User!")
	w.Write([]byte("Hello User!"))
}

func query(city string) (weatherData, error) {

	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		// CHECKPOINT
		fmt.Println(err)
		return weatherData{}, err
	}
	fmt.Println(apiConfig.OpenWeatherApiKey)
	//CALLING THE EXTERNAL API
	// "https://api.openweathermap.org/data/2.5/weather?" + "q=" + city + "&appid=" + apiConfig.OpenWeatherApiKey
	// res, err1 := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID" + apiConfig.OpenWeatherApiKey + "&q=" + city)
	res, err1 := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city + "&APPID=" + apiConfig.OpenWeatherApiKey)

	if err1 != nil {
		// CHECKPOINT
		fmt.Println(err1)
		return weatherData{}, err1
	}

	defer res.Body.Close()

	var data weatherData
	err2 := json.NewDecoder(res.Body).Decode(&data)
	if err2 != nil {
		// CHECKPOINT
		fmt.Println(err2)
		return weatherData{}, err2
	}

	return data, err
}

func getWeatherReport(w http.ResponseWriter, r *http.Request) {
	//Refer to google doc to understand whats going on in "city"
	city := strings.SplitN(r.URL.Path, "/", 3)[2]
	data, err := query(city)
	if err != nil {
		// CHECKPOINT
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", getWeatherReport)

	http.ListenAndServe(":8080", nil)
}
