package models

type WeatherData struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
	Humidity    int     `json:"humidity"`
}

type APIResponse struct {
	Requested   WeatherData `json:"requested"`
	Alternative WeatherData `json:"alternative"`
	Mode        string      `json:"mode"`
	Message     string      `json:"message"`
}

type owmMain struct {
	Temp     float64 `json:"temp"`
	Humidity int     `json:"humidity"`
}

type owmWeatherEntry struct {
	Description string `json:"description"`
}

type OWMWeatherResponse struct {
	Name    string            `json:"name"`
	Main    owmMain           `json:"main"`
	Weather []owmWeatherEntry `json:"weather"`
}

func (o OWMWeatherResponse) ToWeatherData() WeatherData {
	desc := ""
	if len(o.Weather) > 0 {
		desc = o.Weather[0].Description
	}

	return WeatherData{
		City:        o.Name,
		Temperature: o.Main.Temp,
		Description: desc,
		Humidity:    o.Main.Humidity,
	}
}
