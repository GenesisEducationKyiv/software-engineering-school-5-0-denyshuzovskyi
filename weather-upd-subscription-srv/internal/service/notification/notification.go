package notification

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-lib/pkg/command/notification"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/dto"
)

type NotificationSender interface {
	SendWeatherUpdate(context.Context, notification.SendWeatherUpdate) error
}

type NotificationService struct {
	notificationSender NotificationSender
}

func NewNotificationService(notificationSender NotificationSender) *NotificationService {
	return &NotificationService{
		notificationSender: notificationSender,
	}
}

func (s *NotificationService) SendWeatherUpdateNotification(ctx context.Context, subscriptionData dto.SubscriptionData, weather dto.WeatherDTO) error {
	sendWeatherUpd := notification.SendWeatherUpdate{
		NotificationWithToken: notification.NotificationWithToken{
			Notification: notification.Notification{
				To: subscriptionData.Email,
			},
			Token: subscriptionData.Token,
		},
		Weather: notification.Weather{
			Location:    subscriptionData.Location,
			Temperature: weather.Temperature,
			Humidity:    weather.Humidity,
			Description: weather.Description,
		},
	}

	if err := s.notificationSender.SendWeatherUpdate(ctx, sendWeatherUpd); err != nil {
		return fmt.Errorf("send weather update notification: %w", err)
	}

	return nil
}
