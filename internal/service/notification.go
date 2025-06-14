package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/sqlutil"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
)

type NotificationService struct {
	db                     *sql.DB
	weatherProvider        WeatherProvider
	weatherRepository      WeatherRepository
	subscriberRepository   SubscriberRepository
	subscriptionRepository SubscriptionRepository
	tokenRepository        TokenRepository
	emailSender            EmailSender
	log                    *slog.Logger
}

func NewNotificationService(
	db *sql.DB,
	weatherProvider WeatherProvider,
	weatherRepository WeatherRepository,
	subscriberRepository SubscriberRepository,
	subscriptionRepository SubscriptionRepository,
	tokenRepository TokenRepository,
	emailSender EmailSender,
	log *slog.Logger) *NotificationService {
	return &NotificationService{
		db:                     db,
		weatherProvider:        weatherProvider,
		weatherRepository:      weatherRepository,
		subscriberRepository:   subscriberRepository,
		subscriptionRepository: subscriptionRepository,
		tokenRepository:        tokenRepository,
		emailSender:            emailSender,
		log:                    log,
	}
}

func (s *NotificationService) SendNotifications(frequency model.Frequency, emailData config.EmailData) {
	s.log.Info("triggered SendNotifications", "frequency", frequency)
	ctx := context.Background()

	var subscriptions []*model.Subscription
	err := sqlutil.WithTx(ctx, s.db, &sql.TxOptions{ReadOnly: true}, func(tx *sql.Tx) error {
		var err error
		subscriptions, err = s.subscriptionRepository.FindAllByFrequencyAndConfirmedStatus(ctx, tx, frequency)
		return err
	})
	if err != nil {
		s.log.Error("failed to fetch subscriptions", "error", err)
		return
	}

	failures := 0
	for _, sub := range subscriptions {
		if err := s.processAndSendNotification(ctx, sub, emailData); err != nil {
			s.log.Error("failed to send notification",
				"subscriber_id", sub.SubscriberId,
				"location", sub.LocationName,
				"error", err,
			)
			failures++
			continue
		}
		s.log.Debug("notification sent",
			"subscriber_id", sub.SubscriberId,
			"location", sub.LocationName,
		)
	}

	s.log.Info("finished sending notifications", "total", len(subscriptions), "failures", failures)
}

func (s *NotificationService) processAndSendNotification(ctx context.Context, sub *model.Subscription, emailData config.EmailData) error {
	subscriber, err := s.subscriberRepository.FindById(ctx, s.db, sub.SubscriberId)
	if err != nil {
		return fmt.Errorf("fetch subscriber: %w", err)
	}

	token, err := s.tokenRepository.FindBySubscriptionIdAndType(ctx, s.db, sub.Id, model.TokenType_Unsubscribe)
	if err != nil {
		return fmt.Errorf("fetch token: %w", err)
	}

	lastWeather, err := s.weatherRepository.FindLastUpdatedByLocation(ctx, s.db, sub.LocationName)
	if err != nil {
		return fmt.Errorf("fetch weather: %w", err)
	}

	if lastWeather == nil || lastWeather.LastUpdated.Add(15*time.Minute).Before(time.Now()) {
		weather, err := s.weatherProvider.GetCurrentWeather(sub.LocationName)
		if err != nil {
			return fmt.Errorf("update weather: %w", err)
		}
		weather.FetchedAt = time.Now().UTC()
		if err := s.weatherRepository.Save(ctx, s.db, weather); err != nil {
			return fmt.Errorf("save weather: %w", err)
		}
		lastWeather = weather
	}

	email := dto.SimpleEmail{
		From:    emailData.From,
		To:      subscriber.Email,
		Subject: emailData.Subject,
		Text: fmt.Sprintf(
			emailData.Text,
			lastWeather.LocationName,
			lastWeather.Temperature,
			lastWeather.Humidity,
			lastWeather.Description,
			token.Token,
		),
	}

	if err := s.emailSender.Send(ctx, email); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}
