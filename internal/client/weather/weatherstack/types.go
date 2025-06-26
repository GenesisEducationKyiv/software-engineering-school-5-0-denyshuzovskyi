package weatherstack

type Request struct {
	Type     string `json:"type"`
	Query    string `json:"query"`
	Language string `json:"language"`
	Unit     string `json:"unit"`
}

type Location struct {
	Name           string `json:"name"`
	Country        string `json:"country"`
	Region         string `json:"region"`
	Lat            string `json:"lat"`
	Lon            string `json:"lon"`
	Timezone       string `json:"timezone_id"`
	Localtime      string `json:"localtime"`
	LocaltimeEpoch int64  `json:"localtime_epoch"`
}

type Astro struct {
	Sunrise          string `json:"sunrise"`
	Sunset           string `json:"sunset"`
	Moonrise         string `json:"moonrise"`
	Moonset          string `json:"moonset"`
	MoonPhase        string `json:"moon_phase"`
	MoonIllumination int    `json:"moon_illumination"`
}

type AirQuality struct {
	CO           string `json:"co"`
	NO2          string `json:"no2"`
	O3           string `json:"o3"`
	SO2          string `json:"so2"`
	PM25         string `json:"pm2_5"`
	PM10         string `json:"pm10"`
	USEPAIndex   string `json:"us-epa-index"`
	GBDEFRAIndex string `json:"gb-defra-index"`
}

type Current struct {
	ObservationTime     string     `json:"observation_time"`
	Temperature         int        `json:"temperature"`
	WeatherCode         int        `json:"weather_code"`
	WeatherIcons        []string   `json:"weather_icons"`
	WeatherDescriptions []string   `json:"weather_descriptions"`
	Astro               Astro      `json:"astro"`
	AirQuality          AirQuality `json:"air_quality"`
	WindSpeed           int        `json:"wind_speed"`
	WindDegree          int        `json:"wind_degree"`
	WindDir             string     `json:"wind_dir"`
	Pressure            int        `json:"pressure"`
	Precip              int        `json:"precip"`
	Humidity            int        `json:"humidity"`
	Cloudcover          int        `json:"cloudcover"`
	Feelslike           int        `json:"feelslike"`
	UVIndex             int        `json:"uv_index"`
	Visibility          int        `json:"visibility"`
}

type CurrentWeather struct {
	Request  Request  `json:"request"`
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}
