package handler

import (
	"encoding/json"
	"net/http"

	"github.com/arkag/cope-and-hope-weather/client"
	"github.com/arkag/cope-and-hope-weather/models"
)

type Server struct {
	WeatherClient client.Client
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

	resp := models.APIResponse{
		Requested: weather,
		Mode:      mode,
		Message:   "Weather fetched successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
