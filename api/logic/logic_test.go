package logic

import (
	"strings"
	"testing"

	"github.com/arkag/cope-and-hope-weather/models"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		name    string
		weather models.WeatherData
		want    Condition
	}{
		{
			name: "hot and humid",
			weather: models.WeatherData{
				Temperature: 35.0,
				Humidity:    85,
				Description: "clear sky",
			},
			want: HotHumid,
		},
		{
			name: "hot and dry",
			weather: models.WeatherData{
				Temperature: 42.0,
				Humidity:    15,
				Description: "sunny",
			},
			want: HotDry,
		},
		{
			name: "cold",
			weather: models.WeatherData{
				Temperature: -5.0,
				Humidity:    80,
				Description: "snow",
			},
			want: Cold,
		},
		{
			name: "rainy dominates hot",
			weather: models.WeatherData{
				Temperature: 32.0,
				Humidity:    90,
				Description: "heavy thunderstorm",
			},
			want: Rainy,
		},
		{
			name: "mild",
			weather: models.WeatherData{
				Temperature: 20.0,
				Humidity:    50,
				Description: "partly cloudy",
			},
			want: Mild,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Classify(tt.weather); got != tt.want {
				t.Errorf("Classify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPickAlternativeCity(t *testing.T) {
	// Test Hope mode
	hopeCity := PickAlternativeCity("hope", models.WeatherData{})
	found := false
	for _, c := range hopeCities {
		if c == hopeCity {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("PickAlternativeCity('hope') returned %q, not in hopeCities", hopeCity)
	}

	// Test Cope mode (Cold)
	coldWeather := models.WeatherData{Temperature: -10}
	copeCity := PickAlternativeCity("cope", coldWeather)
	found = false
	for _, c := range copeCities[Cold] {
		if c == copeCity {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("PickAlternativeCity('cope', cold) returned %q, not in Cold cities", copeCity)
	}
}

func TestGenerateMessage(t *testing.T) {
	req := models.WeatherData{
		City:        "Seattle",
		Temperature: 8.5,
		Description: "light rain",
	}
	alt := models.WeatherData{
		City:        "Honolulu",
		Temperature: 27.2,
		Description: "clear sky",
	}

	copeMsg := GenerateMessage("cope", req, alt)
	if !strings.Contains(copeMsg, "dealing with clear sky at 27.2°C") {
		t.Errorf("GenerateMessage(cope) missing expected text. Got: %s", copeMsg)
	}

	hopeMsg := GenerateMessage("hope", req, alt)
	if !strings.Contains(hopeMsg, "perfect 27.2°C") {
		t.Errorf("GenerateMessage(hope) missing expected text. Got: %s", hopeMsg)
	}
}
