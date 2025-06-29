package dto

type WeatherDTO struct {
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
	Description string  `json:"description"`
}

type Location struct {
	Name string
}

type WeatherWithLocationDTO struct {
	Weather     WeatherDTO
	Location    Location
	LastUpdated int64
}
