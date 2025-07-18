package weatherupd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
)

type SubscriptionService interface {
	PrepareSubscriptionDataForFrequency(context.Context, model.Frequency) ([]dto.SubscriptionData, error)
}

type WeatherService interface {
	GetCurrentWeatherForLocation(context.Context, string) (*dto.WeatherDTO, error)
}

type NotificationService interface {
	SendWeatherUpdateNotification(context.Context, config.EmailData, dto.SubscriptionData, dto.WeatherDTO) error
}

type WeatherUpdateSendingService struct {
	subscriptionService SubscriptionService
	weatherService      WeatherService
	notificationService NotificationService
	log                 *slog.Logger
}

func NewWeatherUpdateSendingService(subscriptionService SubscriptionService, weatherService WeatherService, notificationService NotificationService, log *slog.Logger) *WeatherUpdateSendingService {
	return &WeatherUpdateSendingService{
		subscriptionService: subscriptionService,
		weatherService:      weatherService,
		notificationService: notificationService,
		log:                 log,
	}
}

func (s *WeatherUpdateSendingService) SendWeatherUpdates(ctx context.Context, frequency model.Frequency, emailData config.EmailData) {
	s.log.Info("started SendWeatherUpdates", "frequency", frequency)

	subscriptionData, err := s.subscriptionService.PrepareSubscriptionDataForFrequency(ctx, frequency)
	if err != nil {
		s.log.Error("failed to get subscription data", "error", err)
		return
	}

	failures := 0
	for _, subsData := range subscriptionData {
		if err := s.fetchWeatherAndSendNotification(ctx, subsData, emailData); err != nil {
			s.log.Error("failed to send weather update",
				"subscriber", subsData.Email,
				"location", subsData.Location,
				"error", err,
			)
			failures++
			continue
		}
		s.log.Debug("weather update sent",
			"subscriber", subsData.Email,
			"location", subsData.Location,
		)
	}

	s.log.Info("finished SendWeatherUpdates", "frequency", frequency, "total", len(subscriptionData), "failures", failures)
}

func (s *WeatherUpdateSendingService) fetchWeatherAndSendNotification(ctx context.Context, subscriptionData dto.SubscriptionData, emailData config.EmailData) error {
	weatherDTO, err := s.weatherService.GetCurrentWeatherForLocation(ctx, subscriptionData.Location)
	if err != nil {
		return fmt.Errorf("failed to get current weather for location %s: %w", subscriptionData.Location, err)
	}
	err = s.notificationService.SendWeatherUpdateNotification(ctx, emailData, subscriptionData, *weatherDTO)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}
