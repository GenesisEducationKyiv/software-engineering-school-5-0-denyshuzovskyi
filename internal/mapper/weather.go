package mapper

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
)

func WeatherToWeatherDTO(weather model.Weather) dto.WeatherDTO {
	return dto.WeatherDTO{
		Temperature: weather.Temperature,
		Humidity:    weather.Humidity,
		Description: weather.Description,
	}
}
