package weatherupd

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/model"
)

type SubscriptionService interface {
	PrepareSubscriptionDataForFrequency(context.Context, model.Frequency) ([]dto.SubscriptionData, error)
}

type WeatherService interface {
	GetCurrentWeatherForLocation(context.Context, string) (*dto.WeatherDTO, error)
}

type NotificationService interface {
	SendWeatherUpdateNotification(context.Context, dto.SubscriptionData, dto.WeatherDTO) error
}

type WeatherUpdJobMetrics interface {
	RecordJobRun(string)
	RecordNotificationSent(string)
	RecordNotificationFailed(string)
	RecordSubscriptionsHandled(string, int)
	ObserveJobDuration(string, float64)
}

type WeatherUpdateSendingService struct {
	subscriptionService  SubscriptionService
	weatherService       WeatherService
	notificationService  NotificationService
	weatherUpdJobMetrics WeatherUpdJobMetrics
	log                  *slog.Logger
}

func NewWeatherUpdateSendingService(subscriptionService SubscriptionService, weatherService WeatherService, notificationService NotificationService, weatherUpdJobMetrics WeatherUpdJobMetrics, log *slog.Logger) *WeatherUpdateSendingService {
	return &WeatherUpdateSendingService{
		subscriptionService:  subscriptionService,
		weatherService:       weatherService,
		notificationService:  notificationService,
		weatherUpdJobMetrics: weatherUpdJobMetrics,
		log:                  log,
	}
}

func (s *WeatherUpdateSendingService) SendWeatherUpdates(ctx context.Context, frequency model.Frequency) {
	s.log.Info("started SendWeatherUpdates", "frequency", frequency)
	start := time.Now()
	s.weatherUpdJobMetrics.RecordJobRun(string(frequency))

	subscriptionData, err := s.subscriptionService.PrepareSubscriptionDataForFrequency(ctx, frequency)
	if err != nil {
		s.log.Error("failed to get subscription data", "error", err)
		return
	}
	s.weatherUpdJobMetrics.RecordSubscriptionsHandled(string(frequency), len(subscriptionData))

	failures := 0
	for _, subsData := range subscriptionData {
		if err := s.fetchWeatherAndSendNotification(ctx, subsData); err != nil {
			s.log.Error("failed to send weather update",
				"subscriber", subsData.Email,
				"location", subsData.Location,
				"error", err,
			)
			failures++
			s.weatherUpdJobMetrics.RecordNotificationFailed(string(frequency))
			continue
		}
		s.weatherUpdJobMetrics.RecordNotificationSent(string(frequency))
		s.log.Debug("weather update sent",
			"subscriber", subsData.Email,
			"location", subsData.Location,
		)
	}

	duration := time.Since(start).Seconds()
	s.weatherUpdJobMetrics.ObserveJobDuration(string(frequency), duration)
	s.log.Info("finished SendWeatherUpdates", "frequency", frequency, "total", len(subscriptionData), "failures", failures)
}

func (s *WeatherUpdateSendingService) fetchWeatherAndSendNotification(ctx context.Context, subscriptionData dto.SubscriptionData) error {
	weatherDTO, err := s.weatherService.GetCurrentWeatherForLocation(ctx, subscriptionData.Location)
	if err != nil {
		return fmt.Errorf("failed to get current weather for location %s: %w", subscriptionData.Location, err)
	}
	err = s.notificationService.SendWeatherUpdateNotification(ctx, subscriptionData, *weatherDTO)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}
