package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arkag/cope-and-hope-weather/client"
	"github.com/arkag/cope-and-hope-weather/models"
)

// newFakeOWMServer creates a fake OWM API that responds based on the city.
// Known cities return 200 with weather data; unknown cities return 404.
func newFakeOWMServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		city := r.URL.Query().Get("q")

		responses := map[string]models.OWMWeatherResponse{
			"London": {
				Name:    "London",
				Main:    models.OWMMain{Temp: 15.0, Humidity: 80},
				Weather: []models.OWMWeatherEntry{{Description: "overcast clouds"}},
			},
			"Tokyo": {
				Name:    "Tokyo",
				Main:    models.OWMMain{Temp: 28.0, Humidity: 65},
				Weather: []models.OWMWeatherEntry{{Description: "clear sky"}},
			},
		}

		resp, ok := responses[city]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

func TestHandleWeather(t *testing.T) {
	fakeOWM := newFakeOWMServer()
	defer fakeOWM.Close()

	s := &Server{
		WeatherClient: client.Client{
			APIKey:  "test-key",
			BaseURL: fakeOWM.URL,
		},
	}

	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantMode   string
		wantCity   string
	}{
		{
			name:       "missing city returns 400",
			url:        "/weather?mode=cope",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing mode returns 400",
			url:        "/weather?city=London",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid mode returns 400",
			url:        "/weather?city=London&mode=yolo",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "cope mode success",
			url:        "/weather?city=London&mode=cope",
			wantStatus: http.StatusOK,
			wantMode:   "cope",
			wantCity:   "London",
		},
		{
			name:       "hope mode success",
			url:        "/weather?city=Tokyo&mode=hope",
			wantStatus: http.StatusOK,
			wantMode:   "hope",
			wantCity:   "Tokyo",
		},
		{
			name:       "unknown city returns 502",
			url:        "/weather?city=FakeCity&mode=cope",
			wantStatus: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			rec := httptest.NewRecorder()

			s.HandleWeather(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status code: got %d, want %d", rec.Code, tt.wantStatus)
			}

			// For successful responses, verify the JSON body
			if tt.wantStatus == http.StatusOK {
				var resp models.APIResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if resp.Mode != tt.wantMode {
					t.Errorf("mode: got %q, want %q", resp.Mode, tt.wantMode)
				}
				if resp.Requested.City != tt.wantCity {
					t.Errorf("city: got %q, want %q", resp.Requested.City, tt.wantCity)
				}
			}
		})
	}
}
