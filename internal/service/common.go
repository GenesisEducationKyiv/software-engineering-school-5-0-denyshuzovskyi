package service

import (
	"context"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/sqlutil"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
)

type WeatherProvider interface {
	GetCurrentWeather(string) (*model.WeatherWithLocation, error)
}

type LocationRepository interface {
	Save(context.Context, sqlutil.SQLExecutor, *model.Location) (int32, error)
	FindByName(context.Context, sqlutil.SQLExecutor, string) (*model.Location, error)
	FindById(context.Context, sqlutil.SQLExecutor, int32) (*model.Location, error)
}

type EmailSender interface {
	Send(context.Context, dto.SimpleEmail) error
}
