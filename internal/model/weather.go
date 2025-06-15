package model

import "time"

type Weather struct {
	LocationName string
	LastUpdated  time.Time
	FetchedAt    time.Time
	Temperature  float32
	Humidity     float32
	Description  string
}
