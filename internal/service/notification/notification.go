package notification

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
)

type EmailSender interface {
	Send(context.Context, dto.SimpleEmail) error
}

type NotificationService struct {
	emailSender EmailSender
}

func NewNotificationService(emailSender EmailSender) *NotificationService {
	return &NotificationService{
		emailSender: emailSender,
	}
}

func (s *NotificationService) SendWeatherUpdateNotification(ctx context.Context, emailData config.EmailData, subscriptionData dto.SubscriptionData, weather dto.WeatherDTO) error {
	email := dto.SimpleEmail{
		From:    emailData.From,
		To:      subscriptionData.Email,
		Subject: emailData.Subject,
		Text: fmt.Sprintf(
			emailData.Text,
			subscriptionData.Location,
			weather.Temperature,
			weather.Humidity,
			weather.Description,
			subscriptionData.Token,
		),
	}

	if err := s.emailSender.Send(ctx, email); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}
