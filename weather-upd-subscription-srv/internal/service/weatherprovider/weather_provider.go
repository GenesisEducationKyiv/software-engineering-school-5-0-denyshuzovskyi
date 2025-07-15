package weatherprovider

import (
	"context"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/dto"
)

type WeatherProvider interface {
	GetCurrentWeather(context.Context, string) (*dto.WeatherWithLocationDTO, error)
}
