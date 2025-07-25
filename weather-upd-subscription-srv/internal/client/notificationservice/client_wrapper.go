package notificationservice

import (
	"context"

	v1 "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-proto/gen/go/notification/v1"
)

type NotificationServiceClientWrapper struct {
	client v1.NotificationServiceClient
}

func NewNotificationServiceClient(client v1.NotificationServiceClient) *NotificationServiceClientWrapper {
	return &NotificationServiceClientWrapper{
		client: client,
	}
}

func (w *NotificationServiceClientWrapper) SendConfirmation(ctx context.Context, req *v1.SendConfirmationRequest) error {
	_, err := w.client.SendConfirmation(ctx, req)

	return err
}

func (w *NotificationServiceClientWrapper) SendConfirmationSuccess(ctx context.Context, req *v1.SendConfirmationSuccessRequest) error {
	_, err := w.client.SendConfirmationSuccess(ctx, req)

	return err
}

func (w *NotificationServiceClientWrapper) SendUnsubscribeSuccess(ctx context.Context, req *v1.SendUnsubscribeSuccessRequest) error {
	_, err := w.client.SendUnsubscribeSuccess(ctx, req)

	return err
}

func (w *NotificationServiceClientWrapper) SendWeatherUpdate(ctx context.Context, req *v1.SendWeatherUpdateRequest) error {
	_, err := w.client.SendWeatherUpdate(ctx, req)

	return err
}
