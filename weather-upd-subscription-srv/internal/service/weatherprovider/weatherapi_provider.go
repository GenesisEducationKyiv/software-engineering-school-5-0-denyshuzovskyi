package weatherprovider

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/client/weatherapi"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/dto"
	weatherapimapper "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/mapper/weatherprovider/weatherapi"
)

type WeatherapiClient interface {
	GetCurrentWeather(context.Context, string) (*weatherapi.CurrentWeather, error)
}

type WeatherapiProvider struct {
	client WeatherapiClient
}

func NewWeatherapiProvider(client WeatherapiClient) *WeatherapiProvider {
	return &WeatherapiProvider{
		client: client,
	}
}

func (p *WeatherapiProvider) GetCurrentWeather(ctx context.Context, location string) (*dto.WeatherWithLocationDTO, error) {
	currentWeather, err := p.client.GetCurrentWeather(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("weatherapi provider: %w", err)
	}
	weatherWithLocationDTO := weatherapimapper.CurrentWeatherToWeatherWithLocationDTO(*currentWeather)
	return &weatherWithLocationDTO, nil
}
