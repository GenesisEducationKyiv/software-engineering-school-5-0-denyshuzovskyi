package weatherstack

type Location struct {
	Name      string `json:"name"`
	Timezone  string `json:"timezone_id"`
	Localtime string `json:"localtime"`
}

type Current struct {
	ObservationTime     string   `json:"observation_time"`
	Temperature         int      `json:"temperature"`
	WeatherDescriptions []string `json:"weather_descriptions"`
	Humidity            int      `json:"humidity"`
}

type CurrentWeather struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}
