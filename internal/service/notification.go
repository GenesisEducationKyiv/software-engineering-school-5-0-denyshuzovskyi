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

func (s *NotificationService) SendDailyNotifications(emailData config.EmailData) {
	s.log.Info("triggered SendDailyNotifications")
	ctx := context.Background()

	err := sqlutil.WithTx(ctx, s.db, &sql.TxOptions{ReadOnly: true}, func(tx *sql.Tx) error {
		subscriptions, errIn := s.subscriptionRepository.FindAllByFrequencyAndConfirmedStatus(ctx, tx, model.Frequency_Daily)
		if errIn != nil {
			return errIn
		}

		if errIn = s.sendNotificationsToAllSubscribers(ctx, tx, subscriptions, emailData); errIn != nil {
			return errIn
		}

		return nil
	})
	if err != nil {
		s.log.Error("rolled back transaction because of ", "error", err)
		return
	}
	s.log.Info("transaction commited successfully")
}

func (s *NotificationService) SendHourlyNotifications(emailData config.EmailData) {
	s.log.Info("triggered SendHourlyNotifications")
	ctx := context.Background()

	err := sqlutil.WithTx(ctx, s.db, nil, func(tx *sql.Tx) error {
		subscriptions, errIn := s.subscriptionRepository.FindAllByFrequencyAndConfirmedStatus(ctx, tx, model.Frequency_Hourly)
		if errIn != nil {
			return errIn
		}

		if errIn = s.sendNotificationsToAllSubscribers(ctx, tx, subscriptions, emailData); errIn != nil {
			return errIn
		}

		return nil
	})
	if err != nil {
		s.log.Error("rolled back transaction because of ", "error", err)
		return
	}
	s.log.Info("transaction commited successfully")
}

func (s *NotificationService) sendNotificationsToAllSubscribers(ctx context.Context, tx *sql.Tx, subscriptions []*model.Subscription, emailData config.EmailData) error {
	for i := 0; i < len(subscriptions); i++ {
		subscriber, err := s.subscriberRepository.FindById(ctx, tx, subscriptions[i].SubscriberId)
		if err != nil {
			return err
		}
		token, err := s.tokenRepository.FindBySubscriptionIdAndType(ctx, tx, subscriptions[i].Id, model.TokenType_Unsubscribe)
		if err != nil {
			return err
		}
		lastWeather, err := s.weatherRepository.FindLastUpdatedByLocation(ctx, tx, subscriptions[i].LocationName)
		if err != nil {
			return err
		}

		if lastWeather == nil || lastWeather.LastUpdated.Add(15*time.Minute).Before(time.Now()) {
			weather, err := s.weatherProvider.GetCurrentWeather(subscriptions[i].LocationName)
			if err != nil {
				return err
			}

			weather.FetchedAt = time.Now().UTC()

			err = s.weatherRepository.Save(ctx, tx, weather)
			if err != nil {
				return err
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

		err = s.emailSender.Send(ctx, email)
		if err != nil {
			return err
		}

		s.log.Info("weather email is send")
	}

	return nil
}
