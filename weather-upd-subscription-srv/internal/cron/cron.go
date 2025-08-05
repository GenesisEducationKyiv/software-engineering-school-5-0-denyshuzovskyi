package cron

import (
	"context"
	"log/slog"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/model"
	"github.com/robfig/cron/v3"
)

type weatherUpdateSendingService interface {
	SendWeatherUpdates(ctx context.Context, frequency model.Frequency)
}

func SetUpCronJobs(ctx context.Context, weatherUpdateSendingService weatherUpdateSendingService, log *slog.Logger) (*cron.Cron, error) {
	c := cron.New()

	// daily 09:00
	if _, err := c.AddFunc("0 9 * * *", func() {
		weatherUpdateSendingService.SendWeatherUpdates(ctx, model.Frequency_Daily)
	}); err != nil {
		log.Error("failed to schedule daily notification service", "error", err)
		return nil, err
	}

	// hourly
	if _, err := c.AddFunc("0 * * * *", func() {
		weatherUpdateSendingService.SendWeatherUpdates(ctx, model.Frequency_Hourly)
	}); err != nil {
		log.Error("failed to schedule hourly notification service", "error", err)
		return nil, err
	}

	return c, nil
}
