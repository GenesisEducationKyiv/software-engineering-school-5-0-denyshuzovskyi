package weatherapi

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/client/weatherapi"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
)

func CurrentWeatherToWeatherWithLocationDTO(currentWeather weatherapi.CurrentWeather) dto.WeatherWithLocationDTO {
	return dto.WeatherWithLocationDTO{
		Weather: dto.WeatherDTO{
			Temperature: currentWeather.Current.TempC,
			Humidity:    float32(currentWeather.Current.Humidity),
			Description: currentWeather.Current.Condition.Text,
		},
		Location: dto.Location{
			Name: currentWeather.Location.Name,
		},
		LastUpdated: currentWeather.Current.LastUpdated,
	}
}
