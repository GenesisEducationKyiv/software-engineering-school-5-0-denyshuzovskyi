package weatherstack

type Location struct {
	Name      string `json:"name"`
	Country   string `json:"country"`
	Region    string `json:"region"`
	Lat       string `json:"lat"`
	Lon       string `json:"lon"`
	Timezone  string `json:"timezone_id"`
	Localtime string `json:"localtime"`
}

type Current struct {
	ObservationTime     string   `json:"observation_time"`
	Temperature         int      `json:"temperature"`
	WeatherCode         int      `json:"weather_code"`
	WeatherDescriptions []string `json:"weather_descriptions"`
	WindSpeed           int      `json:"wind_speed"`
	WindDegree          int      `json:"wind_degree"`
	WindDir             string   `json:"wind_dir"`
	Pressure            int      `json:"pressure"`
	Precip              int      `json:"precip"`
	Humidity            int      `json:"humidity"`
	Cloudcover          int      `json:"cloudcover"`
	Feelslike           int      `json:"feelslike"`
	UVIndex             int      `json:"uv_index"`
	Visibility          int      `json:"visibility"`
}

type CurrentWeather struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}
