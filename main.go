package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/arkag/cope-and-hope-weather/client"
	"github.com/arkag/cope-and-hope-weather/handler"
)

func main() {
	apiKey := os.Getenv("OWM_API_KEY")
	if apiKey == "" {
		fmt.Println("OWM_API_KEY environment variable is required")
		os.Exit(1)
	}

	s := &handler.Server{
		WeatherClient: client.Client{
			APIKey:  apiKey,
			BaseURL: "https://api.openweathermap.org",
		},
	}

	http.HandleFunc("/weather", s.HandleWeather)

	fmt.Println("Cope and Hope Weather API listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
