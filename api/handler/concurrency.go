package handler

import (
	"log/slog"
	"sync"

	"github.com/arkag/cope-and-hope-weather/models"
)

func (s *Server) FetchAlternativesConcurrently(cities []string) []models.WeatherData {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []models.WeatherData
	for _, city := range cities {
		wg.Add(1)

		go func(c string) {
			defer wg.Done()

			weather, err := s.WeatherClient.FetchWeather(c)
			if err == nil {
				slog.Debug("concurrent fetch succeeded", "city", c)
				mu.Lock()
				results = append(results, weather)
				mu.Unlock()
			} else {
				slog.Debug("concurrent fetch failed", "city", c, "error", err)
			}
		}(city)
	}
	wg.Wait()

	return results
}
