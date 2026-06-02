package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/arkag/cope-and-hope-weather/cache"
	"github.com/arkag/cope-and-hope-weather/logic"
	"github.com/arkag/cope-and-hope-weather/models"
)

type Server struct {
	WeatherClient cache.WeatherFetcher
}

func (s *Server) HandleWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	mode := r.URL.Query().Get("mode")

	if city == "" {
		http.Error(w, `{"error":"city parameter is required"}`, http.StatusBadRequest)
		return
	}

	if mode != "cope" && mode != "hope" {
		http.Error(w, `{"error":"mode must be 'cope' or 'hope'"}`, http.StatusBadRequest)
		return
	}

	weather, err := s.WeatherClient.FetchWeather(city)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch weather"}`, http.StatusBadGateway)
		return
	}

	// 1. Pick 3 potential alternative cities
	var altCities []string
	for i := 0; i < 3; i++ {
		altCities = append(altCities, logic.PickAlternativeCity(mode, weather))
	}

	// 2. Fetch all 3 simultaneously!
	slog.Debug("starting concurrent fetch for alternative cities", "cities", altCities)
	altResults := s.FetchAlternativesConcurrently(altCities)

	if len(altResults) == 0 {
		http.Error(w, `{"error":"failed to fetch alternative weather"}`, http.StatusBadGateway)
		return
	}

	// 3. Find the best match
	var alternative models.WeatherData
	foundMatch := false

	// Determine what extreme condition we are trying to match for 'cope' mode
	targetCondition := logic.Classify(weather)

	for _, alt := range altResults {
		if mode == "hope" {
			// Hope cities are always good, take the first one that succeeded
			alternative = alt
			foundMatch = true
			break
		}

		// For 'cope' mode, verify it ACTUALLY matches the miserable condition right now
		if logic.Classify(alt) == targetCondition {
			alternative = alt
			foundMatch = true
			break
		}
	}

	// 4. Fallback if none perfectly matched the extreme (e.g. it was unusually nice in all 3)
	if !foundMatch {
		alternative = altResults[0]
	}

	// 5. Generate the conversational message
	msg := logic.GenerateMessage(mode, weather, alternative)

	resp := models.APIResponse{
		Requested:   weather,
		Alternative: alternative,
		Mode:        mode,
		Message:     msg,
	}

	slog.Info("weather processed",
		"mode", mode,
		"requested_city", weather.City,
		"alternative_city", alternative.City,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
