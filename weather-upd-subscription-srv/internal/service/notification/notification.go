package notification

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-proto/gen/go/notification/v1"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/dto"
)

type NotificationSender interface {
	SendWeatherUpdate(context.Context, *v1.SendWeatherUpdateRequest) error
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
	weatherUpdReq := &v1.SendWeatherUpdateRequest{
		WeatherUpdateNotification: &v1.WeatherUpdateNotificationRequest{
			NotificationWithToken: &v1.NotificationWithTokenRequest{
				Notification: &v1.NotificationRequest{
					To: subscriptionData.Email,
				},
				Token: subscriptionData.Token,
			},
			Weather: &v1.WeatherRequest{
				Location:    subscriptionData.Location,
				Temperature: weather.Temperature,
				Humidity:    weather.Humidity,
				Description: weather.Description,
			},
		},
	}

	if err := s.notificationSender.SendWeatherUpdate(ctx, weatherUpdReq); err != nil {
		return fmt.Errorf("send weather update notification: %w", err)
	}

	return nil
}
