package cron

import (
	"context"
	"log/slog"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
	"github.com/robfig/cron/v3"
)

type notificationService interface {
	SendNotifications(context.Context, model.Frequency, config.EmailData)
}

func SetUpCronJobs(ctx context.Context, notificationService notificationService, weatherEmailData config.EmailData, log *slog.Logger) (*cron.Cron, error) {
	c := cron.New()

	// daily 09:00
	if _, err := c.AddFunc("0 9 * * *", func() {
		notificationService.SendNotifications(ctx, model.Frequency_Daily, weatherEmailData)
	}); err != nil {
		log.Error("failed to schedule daily notification service", "error", err)
		return nil, err
	}

	// hourly
	if _, err := c.AddFunc("0 * * * *", func() {
		notificationService.SendNotifications(ctx, model.Frequency_Hourly, weatherEmailData)
	}); err != nil {
		log.Error("failed to schedule hourly notification service", "error", err)
		return nil, err
	}

	return c, nil
}
