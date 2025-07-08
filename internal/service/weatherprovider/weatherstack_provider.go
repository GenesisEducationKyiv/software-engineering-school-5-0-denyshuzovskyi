package weatherprovider

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/client/weatherstack"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	weatherstackmapper "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/mapper/weatherprovider/weatherstack"
)

type WeatherstackClient interface {
	GetCurrentWeather(context.Context, string) (*weatherstack.CurrentWeather, error)
}

type WeatherstackProvider struct {
	client WeatherstackClient
}

func NewWeatherstackProvider(client WeatherstackClient) *WeatherstackProvider {
	return &WeatherstackProvider{
		client: client,
	}
}

func (p *WeatherstackProvider) GetCurrentWeather(ctx context.Context, location string) (*dto.WeatherWithLocationDTO, error) {
	currentWeather, err := p.client.GetCurrentWeather(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("weatherstack provider: %w", err)
	}

	var weatherWithLocationDTO dto.WeatherWithLocationDTO
	weatherWithLocationDTO, err = weatherstackmapper.CurrentWeatherToWeatherWithLocationDTO(*currentWeather)
	if err != nil {
		return nil, fmt.Errorf("weatherstack provider mapping error: %w", err)
	}
	return &weatherWithLocationDTO, nil
}
