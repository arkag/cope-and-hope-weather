package models

import (
	"encoding/json"
	"testing"
)

func TestWeatherDataMarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    WeatherData
		wantJSON string
	}{
		{
			name: "sunny city",
			input: WeatherData{
				City:        "Phoenix",
				Temperature: 45.2,
				Description: "clear sky",
				Humidity:    10,
			},
			wantJSON: `{"city":"Phoenix","temperature":45.2,"description":"clear sky","humidity":10}`,
		},
		{
			name: "zero values",
			input: WeatherData{
				City:        "",
				Temperature: 0,
				Description: "",
				Humidity:    0,
			},
			wantJSON: `{"city":"","temperature":0,"description":"","humidity":0}`,
		},
		{
			name: "negative temperature",
			input: WeatherData{
				City:        "Yakutsk",
				Temperature: -42.5,
				Description: "freezing fog",
				Humidity:    85,
			},
			wantJSON: `{"city":"Yakutsk","temperature":-42.5,"description":"freezing fog","humidity":85}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("unexpected error marshaling: %v", err)
			}
			if string(got) != tt.wantJSON {
				t.Errorf("got  %s\nwant %s", string(got), tt.wantJSON)
			}
		})
	}
}

func TestWeatherDataUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    WeatherData
	}{
		{
			name:    "full payload",
			jsonStr: `{"city":"Denver","temperature":22.3,"description":"scattered clouds","humidity":40}`,
			want: WeatherData{
				City:        "Denver",
				Temperature: 22.3,
				Description: "scattered clouds",
				Humidity:    40,
			},
		},
		{
			name:    "ignores unknown fields",
			jsonStr: `{"city":"Oslo","temperature":-5.0,"description":"snow","humidity":90,"wind_speed":15.2}`,
			want: WeatherData{
				City:        "Oslo",
				Temperature: -5.0,
				Description: "snow",
				Humidity:    90,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got WeatherData
			if err := json.Unmarshal([]byte(tt.jsonStr), &got); err != nil {
				t.Fatalf("unexpected error unmarshaling: %v", err)
			}
			if got != tt.want {
				t.Errorf("got  %+v\nwant %+v", got, tt.want)
			}
		})
	}
}

func TestAPIResponseMarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    APIResponse
		wantJSON string
	}{
		{
			name: "cope mode",
			input: APIResponse{
				Requested: WeatherData{
					City:        "Seattle",
					Temperature: 8.5,
					Description: "light rain",
					Humidity:    88,
				},
				Alternative: WeatherData{
					City:        "Dhaka",
					Temperature: 43.0,
					Description: "haze",
					Humidity:    95,
				},
				Mode:    "cope",
				Message: "It's rainy in Seattle, but at least you're not melting in Dhaka.",
			},
			wantJSON: `{"requested":{"city":"Seattle","temperature":8.5,"description":"light rain","humidity":88},"alternative":{"city":"Dhaka","temperature":43,"description":"haze","humidity":95},"mode":"cope","message":"It's rainy in Seattle, but at least you're not melting in Dhaka."}`,
		},
		{
			name: "hope mode",
			input: APIResponse{
				Requested: WeatherData{
					City:        "Chicago",
					Temperature: -15.0,
					Description: "blizzard",
					Humidity:    70,
				},
				Alternative: WeatherData{
					City:        "Honolulu",
					Temperature: 27.0,
					Description: "clear sky",
					Humidity:    55,
				},
				Mode:    "hope",
				Message: "Chicago is brutal right now. Dream of Honolulu!",
			},
			wantJSON: `{"requested":{"city":"Chicago","temperature":-15,"description":"blizzard","humidity":70},"alternative":{"city":"Honolulu","temperature":27,"description":"clear sky","humidity":55},"mode":"hope","message":"Chicago is brutal right now. Dream of Honolulu!"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("unexpected error marshaling: %v", err)
			}
			if string(got) != tt.wantJSON {
				t.Errorf("got  %s\nwant %s", string(got), tt.wantJSON)
			}
		})
	}
}

func TestOWMWeatherResponseUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		jsonStr  string
		wantCity string
		wantTemp float64
		wantHum  int
		wantDesc string
	}{
		{
			name: "standard response",
			jsonStr: `{
				"name": "London",
				"main": { "temp": 18.5, "humidity": 72 },
				"weather": [{ "description": "light rain" }]
			}`,
			wantCity: "London",
			wantTemp: 18.5,
			wantHum:  72,
			wantDesc: "light rain",
		},
		{
			name: "multiple weather entries uses first",
			jsonStr: `{
				"name": "Tokyo",
				"main": { "temp": 30.0, "humidity": 80 },
				"weather": [
					{ "description": "thunderstorm" },
					{ "description": "heavy rain" }
				]
			}`,
			wantCity: "Tokyo",
			wantTemp: 30.0,
			wantHum:  80,
			wantDesc: "thunderstorm",
		},
		{
			name: "negative temperature",
			jsonStr: `{
				"name": "Murmansk",
				"main": { "temp": -28.3, "humidity": 95 },
				"weather": [{ "description": "heavy snow" }]
			}`,
			wantCity: "Murmansk",
			wantTemp: -28.3,
			wantHum:  95,
			wantDesc: "heavy snow",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got OWMWeatherResponse
			if err := json.Unmarshal([]byte(tt.jsonStr), &got); err != nil {
				t.Fatalf("unexpected error unmarshaling: %v", err)
			}

			wd := got.ToWeatherData()

			if wd.City != tt.wantCity {
				t.Errorf("City: got %q, want %q", wd.City, tt.wantCity)
			}
			if wd.Temperature != tt.wantTemp {
				t.Errorf("Temperature: got %v, want %v", wd.Temperature, tt.wantTemp)
			}
			if wd.Humidity != tt.wantHum {
				t.Errorf("Humidity: got %d, want %d", wd.Humidity, tt.wantHum)
			}
			if wd.Description != tt.wantDesc {
				t.Errorf("Description: got %q, want %q", wd.Description, tt.wantDesc)
			}
		})
	}
}

func TestOWMWeatherResponseEmptyWeatherSlice(t *testing.T) {
	jsonStr := `{
		"name": "Atlantis",
		"main": { "temp": 20.0, "humidity": 50 },
		"weather": []
	}`

	var got OWMWeatherResponse
	if err := json.Unmarshal([]byte(jsonStr), &got); err != nil {
		t.Fatalf("unexpected error unmarshaling: %v", err)
	}

	wd := got.ToWeatherData()

	if wd.Description != "" {
		t.Errorf("Description: got %q, want empty string for empty weather slice", wd.Description)
	}
}
