package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/arkag/cope-and-hope-weather/models"
)

type Client struct {
	APIKey  string
	BaseURL string
}

func (c Client) FetchWeather(city string) (models.WeatherData, error) {
	if city == "" {
		return models.WeatherData{}, fmt.Errorf("city cannot be empty")
	}

	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("%s/data/2.5/weather?q=%s&appid=%s&units=metric", c.BaseURL, encodedCity, c.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		return models.WeatherData{}, fmt.Errorf("fetching weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return models.WeatherData{}, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var owmResp models.OWMWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&owmResp); err != nil {
		return models.WeatherData{}, fmt.Errorf("decoding response: %w", err)
	}

	return owmResp.ToWeatherData(), nil
}
