package notification

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-proto/gen/go/notification/v1"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/dto"
	"google.golang.org/grpc"
)

type NotificationSender interface {
	SendWeatherUpdate(context.Context, *v1.SendWeatherUpdateRequest, ...grpc.CallOption) (*v1.SendWeatherUpdateResponse, error)
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
		WeatherUpdateNotification: &v1.WeatherUpdateNotification{
			NotificationWithToken: &v1.NotificationWithToken{
				Notification: &v1.Notification{
					To: subscriptionData.Email,
				},
				Token: subscriptionData.Token,
			},
			Weather: &v1.Weather{
				Location:    subscriptionData.Location,
				Temperature: weather.Temperature,
				Humidity:    weather.Humidity,
				Description: weather.Description,
			},
		},
	}

	_, err := s.notificationSender.SendWeatherUpdate(ctx, weatherUpdReq)
	if err != nil {
		return fmt.Errorf("send weather update notification: %w", err)
	}

	return nil
}
