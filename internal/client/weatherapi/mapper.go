package weatherapi

import (
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
)

func CurrentWeatherToWeather(currentWeather CurrentWeather) model.Weather {
	return model.Weather{
		LocationName: currentWeather.Location.Name,
		LastUpdated:  time.Unix(currentWeather.Current.LastUpdated, 0).UTC(),
		FetchedAt:    time.Unix(0, 0),
		Temperature:  currentWeather.Current.TempC,
		Humidity:     float32(currentWeather.Current.Humidity),
		Description:  currentWeather.Current.Condition.Text,
	}
}
