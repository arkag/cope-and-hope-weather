package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arkag/cope-and-hope-weather/models"
)

// newTestServer creates a fake OWM server that returns the given response.
func newTestServer(statusCode int, body any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(body)
	}))
}

func TestFetchWeather(t *testing.T) {
	tests := []struct {
		name       string
		city       string
		statusCode int
		response   models.OWMWeatherResponse
		want       models.WeatherData
		wantErr    bool
	}{
		{
			name:       "successful fetch",
			city:       "London",
			statusCode: http.StatusOK,
			response: models.OWMWeatherResponse{
				Name: "London",
				Main: models.OWMMain{Temp: 18.5, Humidity: 72},
				Weather: []models.OWMWeatherEntry{
					{Description: "light rain"},
				},
			},
			want: models.WeatherData{
				City:        "London",
				Temperature: 18.5,
				Description: "light rain",
				Humidity:    72,
			},
			wantErr: false,
		},
		{
			name:       "city not found returns error",
			city:       "FakeCity123",
			statusCode: http.StatusNotFound,
			response:   models.OWMWeatherResponse{},
			want:       models.WeatherData{},
			wantErr:    true,
		},
		{
			name:       "server error returns error",
			city:       "Tokyo",
			statusCode: http.StatusInternalServerError,
			response:   models.OWMWeatherResponse{},
			want:       models.WeatherData{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newTestServer(tt.statusCode, tt.response)
			defer server.Close()

			c := Client{
				APIKey:  "test-key",
				BaseURL: server.URL,
			}

			got, err := c.FetchWeather(tt.city)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FetchWeather() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("FetchWeather()\ngot  %+v\nwant %+v", got, tt.want)
			}
		})
	}
}

func TestFetchWeatherEmptyCity(t *testing.T) {
	c := Client{
		APIKey:  "test-key",
		BaseURL: "http://localhost",
	}

	_, err := c.FetchWeather("")
	if err == nil {
		t.Error("FetchWeather(\"\") should return an error for empty city")
	}
}

func TestFetchWeatherBadURL(t *testing.T) {
	c := Client{
		APIKey:  "test-key",
		BaseURL: "http://localhost:99999",
	}

	_, err := c.FetchWeather("London")
	if err == nil {
		t.Error("FetchWeather() should return an error when the server is unreachable")
	}
}

func TestFetchWeatherQueryParams(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(models.OWMWeatherResponse{
			Name: "Denver",
		})
	}))
	defer server.Close()

	c := Client{
		APIKey:  "my-secret-key",
		BaseURL: server.URL,
	}

	_, _ = c.FetchWeather("Denver")

	wantParams := "q=Denver&appid=my-secret-key&units=metric"
	if capturedQuery != wantParams {
		t.Errorf("query params\ngot  %q\nwant %q", capturedQuery, wantParams)
	}
}
