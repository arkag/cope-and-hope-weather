package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/arkag/cope-and-hope-weather/client"
	"github.com/arkag/cope-and-hope-weather/models"
)

func TestFetchAlternativesConcurrently(t *testing.T) {
	// Create a fake server that purposefully sleeps for 100ms before responding.
	// If 3 requests are made sequentially, it will take 300ms.
	// If made concurrently, it will take ~100ms.
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)

		city := r.URL.Query().Get("q")
		if city == "FailCity" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(models.OWMWeatherResponse{
			Name: city,
			Main: models.OWMMain{Temp: 20.0},
		})
	}))
	defer slowServer.Close()

	s := &Server{
		WeatherClient: client.Client{
			APIKey:  "test-key",
			BaseURL: slowServer.URL,
		},
	}

	cities := []string{"London", "Tokyo", "FailCity", "Sydney"}

	start := time.Now()
	results := s.FetchAlternativesConcurrently(cities)
	duration := time.Since(start)

	// Verify duration proves concurrency.
	// Allow up to 250ms for overhead, but definitely fail if > 300ms.
	if duration > 250*time.Millisecond {
		t.Fatalf("Fetch took too long (%v), requests are not concurrent!", duration)
	}

	// FailCity should be ignored, so we expect exactly 3 results.
	if len(results) != 3 {
		t.Fatalf("expected 3 successful results, got %d", len(results))
	}

	// Check that all 3 successful cities are in the result slice
	foundCities := make(map[string]bool)
	for _, res := range results {
		foundCities[res.City] = true
	}

	for _, expected := range []string{"London", "Tokyo", "Sydney"} {
		if !foundCities[expected] {
			t.Errorf("expected %q in results, but it was missing", expected)
		}
	}
}
