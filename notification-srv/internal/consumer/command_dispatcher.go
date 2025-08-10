package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-lib/pkg/command/notification"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-lib/pkg/message"
	commonerrors "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/error"
)

type NotificationSendingService interface {
	SendConfirmation(context.Context, notification.SendConfirmation) error
	SendConfirmationSuccess(context.Context, notification.SendConfirmationSuccess) error
	SendWeatherUpdate(context.Context, notification.SendWeatherUpdate) error
	SendUnsubscribeSuccess(context.Context, notification.SendUnsubscribeSuccess) error
}

type NotificationCommandDispatcher struct {
	notificationSendingService NotificationSendingService
}

func NewNotificationCommandDispatcher(notificationSendingService NotificationSendingService) *NotificationCommandDispatcher {
	return &NotificationCommandDispatcher{notificationSendingService: notificationSendingService}
}

func (d *NotificationCommandDispatcher) Dispatch(ctx context.Context, envelope message.Envelope) error {
	switch notification.CommandType(envelope.Type) {
	case notification.Confirmation:
		var cmd notification.SendConfirmation
		if err := json.Unmarshal(envelope.Payload, &cmd); err != nil {
			return fmt.Errorf("failed to unmarshal %s: %w", envelope.Type, err)
		}
		return d.notificationSendingService.SendConfirmation(ctx, cmd)

	case notification.ConfirmationSuccess:
		var cmd notification.SendConfirmationSuccess
		if err := json.Unmarshal(envelope.Payload, &cmd); err != nil {
			return fmt.Errorf("failed to unmarshal %s: %w", envelope.Type, err)
		}
		return d.notificationSendingService.SendConfirmationSuccess(ctx, cmd)

	case notification.WeatherUpdate:
		var cmd notification.SendWeatherUpdate
		if err := json.Unmarshal(envelope.Payload, &cmd); err != nil {
			return fmt.Errorf("failed to unmarshal %s: %w", envelope.Type, err)
		}
		return d.notificationSendingService.SendWeatherUpdate(ctx, cmd)

	case notification.UnsubscribeSuccess:
		var cmd notification.SendUnsubscribeSuccess
		if err := json.Unmarshal(envelope.Payload, &cmd); err != nil {
			return fmt.Errorf("failed to unmarshal %s: %w", envelope.Type, err)
		}
		return d.notificationSendingService.SendUnsubscribeSuccess(ctx, cmd)

	default:
		return fmt.Errorf("%w type: %s", commonerrors.ErrUnsupportedCommand, envelope.Type)
	}
}
