package model

import "time"

type Weather struct {
	LocationId  int32
	LastUpdated time.Time
	FetchedAt   time.Time
	Temperature float32
	Humidity    float32
	Description string
}

type WeatherWithLocation struct {
	Weather
	Location
}
