package logic

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/arkag/cope-and-hope-weather/models"
)

// Condition represents a classified weather state.
type Condition string

const (
	HotHumid Condition = "hot_humid"
	HotDry   Condition = "hot_dry"
	Cold     Condition = "cold"
	Rainy    Condition = "rainy"
	Mild     Condition = "mild"
)

// copeCities maps each condition to cities with even WORSE versions of that condition.
var copeCities = map[Condition][]string{
	HotHumid: {
		"Dhaka", "Bangkok", "Jakarta", "Manila", "Singapore", "Kuala Lumpur", "Colombo", "Chennai", "Kolkata", "Yangon",
		"Ho Chi Minh City", "Phnom Penh", "Lagos", "Accra", "Douala", "Guayaquil", "Panama City", "Cartagena", "Havana", "San Juan",
		"Kingston", "Port-au-Prince", "Mumbai", "Cochin", "Malé", "Mogadishu", "Mombasa", "Dar es Salaam", "Beira", "Darwin",
	},
	HotDry: {
		"Kuwait City", "Death Valley", "Dallol", "Riyadh", "Baghdad", "Phoenix", "Las Vegas", "Khartoum", "Djibouti", "Mecca",
		"Doha", "Abu Dhabi", "Dubai", "Muscat", "Ahvaz", "Basra", "Timbuktu", "Agadez", "Alice Springs", "Coober Pedy",
		"Yuma", "Tucson", "Mexicali", "Hermosillo", "Aswan", "Luxor", "Niamey", "N'Djamena", "Nouakchott", "Kidal",
	},
	Cold: {
		"Yakutsk", "Oymyakon", "Norilsk", "Vostok Station", "Yellowknife", "Iqaluit", "Barrow", "Fairbanks", "Harbin", "Ulaanbaatar",
		"Astana", "Dudinka", "Dikson", "Khatanga", "Tiksi", "Verkhoyansk", "Olenyok", "Saskatoon", "Winnipeg", "Regina",
		"Tromso", "Svalbard", "Kiruna", "Rovaniemi", "Murmansk", "Arkhangelsk", "Nuuk", "Ilulissat", "McMurdo Station", "Amundsen-Scott",
	},
	Rainy: {
		"Cherrapunji", "Mawsynram", "Buenaventura", "Quibdo", "Monrovia", "Hilo", "Ketchikan", "Bergen", "Milford Sound", "Yakutat",
		"Prince Rupert", "Taipei", "Baguio", "Kuching", "Cali", "Valdivia", "Puerto Montt", "Manaus", "Belem", "Iquitos",
		"Cayenne", "Georgetown", "Paramaribo", "Douala", "Malabo", "Bata", "Libreville", "Conakry", "Freetown", "Kikori",
	},
	Mild: {
		// Cities known for persistent gloom, overcast, or boring moderate weather
		"London", "Seattle", "Portland", "Vancouver", "Dublin", "Edinburgh", "Amsterdam", "Brussels", "Copenhagen", "Stockholm",
		"Oslo", "Helsinki", "Reykjavik", "Torshavn", "Nuuk", "Wellington", "Hobart", "Falkland Islands", "Punta Arenas", "Ushuaia",
		"Keflavik", "Trondheim", "Bodo", "Gothenburg", "Malmo", "Aarhus", "Odense", "Aalborg", "Glasgow", "Belfast",
	},
}

// hopeCities are places with generally ideal weather, regardless of current conditions.
var hopeCities = []string{
	"Honolulu", "San Diego", "Lisbon", "Medellín", "Málaga", "Nice", "Barcelona", "Valencia", "Alicante", "Palma",
	"Ibiza", "Canary Islands", "Tenerife", "Madeira", "Faro", "Funchal", "Ponta Delgada", "Bermuda", "Nassau", "Grand Cayman",
	"Turks and Caicos", "Aruba", "Curacao", "Bonaire", "St. Barts", "St. Lucia", "Barbados", "Antigua", "Maui", "Sydney",
}

// Classify determines the weather condition from the full WeatherData.
func Classify(w models.WeatherData) Condition {
	desc := strings.ToLower(w.Description)
	isRainy := strings.Contains(desc, "rain") ||
		strings.Contains(desc, "drizzle") ||
		strings.Contains(desc, "thunderstorm") ||
		strings.Contains(desc, "shower")

	switch {
	case isRainy:
		return Rainy
	case w.Temperature > 30 && w.Humidity > 70:
		return HotHumid
	case w.Temperature > 30:
		return HotDry
	case w.Temperature < 10:
		return Cold
	default:
		return Mild
	}
}

// PickAlternativeCity selects a city based on the mode and the user's current weather.
func PickAlternativeCity(mode string, weather models.WeatherData) string {
	if mode == "hope" {
		return hopeCities[rand.IntN(len(hopeCities))]
	}

	condition := Classify(weather)
	cities := copeCities[condition]
	return cities[rand.IntN(len(cities))]
}

// GenerateMessage creates a conversational response based on mode and both cities' weather.
func GenerateMessage(mode string, requested, alternative models.WeatherData) string {
	if mode == "cope" {
		return fmt.Sprintf(
			"Sure, it's %s in %s at %.1f°C — but at least you're not in %s dealing with %s at %.1f°C!",
			requested.Description, requested.City, requested.Temperature,
			alternative.City, alternative.Description, alternative.Temperature,
		)
	}

	return fmt.Sprintf(
		"Hang in there, %s! Right now %s is enjoying %s at a perfect %.1f°C.",
		requested.City, alternative.City, alternative.Description, alternative.Temperature,
	)
}
