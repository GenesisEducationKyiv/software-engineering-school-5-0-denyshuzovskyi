package mapper

import (
	"time"

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

func WeatherWithLocationDTOToWeather(weatherWithLocationDTO dto.WeatherWithLocationDTO) model.Weather {
	return model.Weather{
		LocationName: weatherWithLocationDTO.Location.Name,
		LastUpdated:  time.Unix(weatherWithLocationDTO.LastUpdated, 0).UTC(),
		FetchedAt:    time.Now().UTC(),
		Temperature:  weatherWithLocationDTO.Weather.Temperature,
		Humidity:     weatherWithLocationDTO.Weather.Humidity,
		Description:  weatherWithLocationDTO.Weather.Description,
	}
}
