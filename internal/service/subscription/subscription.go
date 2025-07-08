package subscription

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	commonerrors "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/error"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/sqlutil"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
	"github.com/google/uuid"
)

type WeatherProvider interface {
	GetCurrentWeather(context.Context, string) (*dto.WeatherWithLocationDTO, error)
}

type SubscriberRepository interface {
	Save(context.Context, sqlutil.SQLExecutor, *model.Subscriber) (int32, error)
	FindByEmail(context.Context, sqlutil.SQLExecutor, string) (*model.Subscriber, error)
	FindById(context.Context, sqlutil.SQLExecutor, int32) (*model.Subscriber, error)
}

type SubscriptionRepository interface {
	Save(context.Context, sqlutil.SQLExecutor, *model.Subscription) (int32, error)
	FindBySubscriberIdAndLocationName(context.Context, sqlutil.SQLExecutor, int32, string) (*model.Subscription, error)
	FindById(context.Context, sqlutil.SQLExecutor, int32) (*model.Subscription, error)
	DeleteById(context.Context, sqlutil.SQLExecutor, int32) error
	Update(context.Context, sqlutil.SQLExecutor, *model.Subscription) (*model.Subscription, error)
	FindAllByFrequencyAndConfirmedStatus(context.Context, sqlutil.SQLExecutor, model.Frequency) ([]*model.Subscription, error)
}

type TokenRepository interface {
	Save(context.Context, sqlutil.SQLExecutor, *model.Token) error
	FindByToken(context.Context, sqlutil.SQLExecutor, string) (*model.Token, error)
	FindBySubscriptionIdAndType(context.Context, sqlutil.SQLExecutor, int32, model.TokenType) (*model.Token, error)
}

type EmailSender interface {
	Send(context.Context, dto.SimpleEmail) error
}

type SubscriptionService struct {
	db                      *sql.DB
	weatherProvider         WeatherProvider
	subscriberRepository    SubscriberRepository
	subscriptionRepository  SubscriptionRepository
	tokenRepository         TokenRepository
	emailSender             EmailSender
	confirmEmailData        config.EmailData
	confirmSuccessEmailData config.EmailData
	unsubEmailData          config.EmailData
	log                     *slog.Logger
}

func NewSubscriptionService(
	db *sql.DB,
	weatherProvider WeatherProvider,
	subscriberRepository SubscriberRepository,
	subscriptionRepository SubscriptionRepository,
	tokenRepository TokenRepository,
	emailSender EmailSender,
	confirmEmailData config.EmailData,
	confirmSuccessEmailData config.EmailData,
	unsubEmailData config.EmailData,
	log *slog.Logger,
) *SubscriptionService {
	return &SubscriptionService{
		db:                      db,
		weatherProvider:         weatherProvider,
		subscriberRepository:    subscriberRepository,
		subscriptionRepository:  subscriptionRepository,
		tokenRepository:         tokenRepository,
		emailSender:             emailSender,
		confirmEmailData:        confirmEmailData,
		confirmSuccessEmailData: confirmSuccessEmailData,
		unsubEmailData:          unsubEmailData,
		log:                     log,
	}
}

func (s *SubscriptionService) Subscribe(ctx context.Context, subReq dto.SubscriptionRequest) error {
	weatherWithLocationDTO, err := s.weatherProvider.GetCurrentWeather(ctx, subReq.City)
	if err != nil {
		if errors.Is(err, commonerrors.ErrLocationNotFound) {
			return err
		} else {
			return fmt.Errorf("unable to validate location err:%w", err)
		}
	}

	err = sqlutil.WithTx(ctx, s.db, nil, func(tx *sql.Tx) error {

		subscriber, errIn := s.subscriberRepository.FindByEmail(ctx, tx, subReq.Email)
		if errIn != nil {
			return errIn
		}
		var subscriberId int32
		if subscriber != nil {
			subscriberId = subscriber.Id
		} else {
			subscriberToSave := model.Subscriber{
				Email:     subReq.Email,
				CreatedAt: time.Now().UTC(),
			}
			subscriberId, errIn = s.subscriberRepository.Save(ctx, tx, &subscriberToSave)
			if errIn != nil {
				return errIn
			}
		}

		subscription, errIn := s.subscriptionRepository.FindBySubscriberIdAndLocationName(ctx, tx, subscriberId, weatherWithLocationDTO.Location.Name)
		if errIn != nil {
			return errIn
		}
		if subscription != nil {
			return commonerrors.ErrSubscriptionAlreadyExists
		}

		subscription = &model.Subscription{
			Id:           0,
			SubscriberId: subscriberId,
			LocationName: weatherWithLocationDTO.Location.Name,
			Frequency:    model.Frequency(subReq.Frequency),
			Status:       model.SubscriptionStatus_Pending,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		}
		subscriptionId, errIn := s.subscriptionRepository.Save(ctx, tx, subscription)
		if errIn != nil {
			return errIn
		}

		token := model.Token{
			Token:          uuid.NewString(),
			SubscriptionId: subscriptionId,
			Type:           model.TokenType_Confirmation,
			CreatedAt:      time.Now().UTC(),
			ExpiresAt:      time.Now().UTC().Add(15 * time.Minute),
			UsedAt:         time.Unix(0, 0),
		}
		if errIn = s.tokenRepository.Save(ctx, tx, &token); errIn != nil {
			return errIn
		}

		email := dto.SimpleEmail{
			From:    s.confirmEmailData.From,
			To:      subReq.Email,
			Subject: s.confirmEmailData.Subject,
			Text: fmt.Sprintf(
				s.confirmEmailData.Text,
				token.Token,
			),
		}

		errIn = s.emailSender.Send(ctx, email)
		if errIn != nil {
			return errIn
		}
		s.log.Info("confirmation email is send")

		return nil
	})
	if err != nil {
		s.log.Info("rollback transaction")
		return err
	}
	s.log.Info("transaction commited successfully")

	return nil
}

func (s *SubscriptionService) Confirm(ctx context.Context, tokenStr string) error {
	err := sqlutil.WithTx(ctx, s.db, nil, func(tx *sql.Tx) error {
		token, errIn := s.tokenRepository.FindByToken(ctx, tx, tokenStr)
		if errIn != nil {
			return errIn
		}
		if token == nil {
			return commonerrors.ErrTokenNotFound
		}
		if time.Now().UTC().After(token.ExpiresAt) || token.Type != model.TokenType_Confirmation {
			return commonerrors.ErrInvalidToken
		}

		subscription, errIn := s.subscriptionRepository.FindById(ctx, tx, token.SubscriptionId)
		if errIn != nil {
			return errIn
		}
		if subscription == nil {
			return commonerrors.ErrUnexpectedState
		}

		subscription.Status = model.SubscriptionStatus_Confirmed
		subscription.UpdatedAt = time.Now().UTC()

		_, errIn = s.subscriptionRepository.Update(ctx, tx, subscription)
		if errIn != nil {
			return errIn
		}

		subscriber, errIn := s.subscriberRepository.FindById(ctx, tx, subscription.SubscriberId)
		if errIn != nil {
			return errIn
		}

		unsubToken := model.Token{
			Token:          uuid.NewString(),
			SubscriptionId: token.SubscriptionId,
			Type:           model.TokenType_Unsubscribe,
			CreatedAt:      time.Now().UTC(),
			ExpiresAt:      time.Now().UTC().AddDate(0, 0, 1),
			UsedAt:         time.Unix(0, 0),
		}
		if errIn = s.tokenRepository.Save(ctx, tx, &unsubToken); errIn != nil {
			return errIn
		}

		email := dto.SimpleEmail{
			From:    s.confirmSuccessEmailData.From,
			To:      subscriber.Email,
			Subject: s.confirmSuccessEmailData.Subject,
			Text: fmt.Sprintf(
				s.confirmSuccessEmailData.Text,
				unsubToken.Token,
			),
		}

		errIn = s.emailSender.Send(ctx, email)
		if errIn != nil {
			return errIn
		}
		s.log.Info("confirmation success email is send")

		return nil
	})
	if err != nil {
		s.log.Info("rollback transaction")
		return err
	}
	s.log.Info("transaction commited successfully")

	return nil
}

func (s *SubscriptionService) Unsubscribe(ctx context.Context, tokenStr string) error {
	err := sqlutil.WithTx(ctx, s.db, nil, func(tx *sql.Tx) error {
		token, errIn := s.tokenRepository.FindByToken(ctx, tx, tokenStr)
		if errIn != nil {
			return errIn
		}
		if token == nil {
			return commonerrors.ErrTokenNotFound
		} else if time.Now().UTC().After(token.ExpiresAt) || token.Type != model.TokenType_Unsubscribe {
			return commonerrors.ErrInvalidToken
		}

		subscription, errIn := s.subscriptionRepository.FindById(ctx, tx, token.SubscriptionId)
		if errIn != nil {
			return errIn
		}

		subscriber, errIn := s.subscriberRepository.FindById(ctx, tx, subscription.SubscriberId)
		if errIn != nil {
			return errIn
		}

		errIn = s.subscriptionRepository.DeleteById(ctx, tx, token.SubscriptionId)
		if errIn != nil {
			return errIn
		}

		email := dto.SimpleEmail{
			From:    s.unsubEmailData.From,
			To:      subscriber.Email,
			Subject: s.unsubEmailData.Subject,
			Text:    s.unsubEmailData.Text,
		}

		errIn = s.emailSender.Send(ctx, email)
		if errIn != nil {
			return errIn
		}
		s.log.Info("unsubscribe success email is send")

		return nil
	})
	if err != nil {
		s.log.Info("rollback transaction")
		return err
	}
	s.log.Info("transaction commited successfully")

	return nil
}

func (s *SubscriptionService) PrepareSubscriptionDataForFrequency(ctx context.Context, frequency model.Frequency) ([]dto.SubscriptionData, error) {
	var subscriptionData []dto.SubscriptionData

	subscriptions, err := s.subscriptionRepository.FindAllByFrequencyAndConfirmedStatus(ctx, s.db, frequency)
	if err != nil {
		return subscriptionData, err
	}

	for _, subscription := range subscriptions {
		subscriber, errIn := s.subscriberRepository.FindById(ctx, s.db, subscription.SubscriberId)
		if errIn != nil {
			s.log.Error("failed to find subscriber", "error", errIn, "id", subscription.SubscriberId)
			continue
		}
		token, errIn := s.tokenRepository.FindBySubscriptionIdAndType(ctx, s.db, subscription.Id, model.TokenType_Unsubscribe)
		if errIn != nil {
			s.log.Error("failed to find unsub token for subscription", "error", errIn, "id", subscription.Id)
			continue
		}

		data := dto.SubscriptionData{
			Email:    subscriber.Email,
			Location: subscription.LocationName,
			Token:    token.Token,
		}
		subscriptionData = append(subscriptionData, data)
	}

	return subscriptionData, nil
}
