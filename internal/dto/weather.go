package dto

import (
	"log/slog"
	"time"
)

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

// LogValue implements slog.LogValuer
func (w WeatherWithLocationDTO) LogValue() slog.Value {
	t := time.Unix(w.LastUpdated, 0).UTC()

	return slog.GroupValue(
		slog.Float64("temperature", float64(w.Weather.Temperature)),
		slog.Float64("humidity", float64(w.Weather.Humidity)),
		slog.String("description", w.Weather.Description),
		slog.String("location", w.Location.Name),
		slog.String("last_updated", t.Format(time.RFC3339)),
	)
}
