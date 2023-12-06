package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type apiaConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (apiaConfigData, error) {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		return apiaConfigData{}, err
	}

	var configData apiaConfigData

	err = json.Unmarshal(bytes, &configData)
	if err != nil {
		return apiaConfigData{}, err
	}

	return configData, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w,"Hello User!")
	w.Write([]byte("Hello User!"))
}

func query(city string) (weatherData, error) {

	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}
	//CALLING THE EXTERNAL API
	// "https://api.openweathermap.org/data/2.5/weather?" + "q=" + city + "&appid=" + apiConfig.OpenWeatherMapApiKey
	// res, err1 := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	res, err1 := http.Get("https://api.openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + apiConfig.OpenWeatherMapApiKey)

	if err1 != nil {
		fmt.Println("Error Detected")
		fmt.Println(err1)
		return weatherData{}, err1
	}

	defer res.Body.Close()

	var data weatherData
	err2 := json.NewDecoder(res.Body).Decode(&data)
	if err2 != nil {
		return weatherData{}, err2
	}

	return data, err
}

func getWeatherReport(w http.ResponseWriter, r *http.Request) {
	//Refer to google doc to understand whats going on in "city"
	city := strings.SplitN(r.URL.Path, "/", 3)[2]
	data, err := query(city)
	if err != nil {
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
