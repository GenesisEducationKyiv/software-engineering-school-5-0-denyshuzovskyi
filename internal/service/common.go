package service

import (
	"context"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
)

type WeatherProvider interface {
	GetCurrentWeather(context.Context, string) (*model.Weather, error)
}

type EmailSender interface {
	Send(context.Context, dto.SimpleEmail) error
}
