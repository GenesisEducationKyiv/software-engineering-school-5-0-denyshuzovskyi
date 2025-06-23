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

func (s *NotificationService) SendWeatherUpdateNotification(ctx context.Context, emailData config.EmailData, subscriptionDate dto.SubscriptionData, weather dto.WeatherDTO) error {
	email := dto.SimpleEmail{
		From:    emailData.From,
		To:      subscriptionDate.Email,
		Subject: emailData.Subject,
		Text: fmt.Sprintf(
			emailData.Text,
			subscriptionDate.Location,
			weather.Temperature,
			weather.Humidity,
			weather.Description,
			subscriptionDate.Token,
		),
	}

	if err := s.emailSender.Send(ctx, email); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}
