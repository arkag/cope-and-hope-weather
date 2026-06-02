package handler

import (
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
				mu.Lock()
				results = append(results, weather)
				mu.Unlock()
			}
		}(city)
	}
	wg.Wait()

	return results
}
