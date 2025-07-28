package publisher

import (
	"context"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-lib/pkg/command/notification"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-lib/pkg/rabbitmq"
)

type CommandPublisher interface {
	Publish(context.Context, string, []byte) error
}

type NotificationCommandSender struct {
	cmdPublisher CommandPublisher
}

func NewNotificationCommandSender(cmdPublisher CommandPublisher) *NotificationCommandSender {
	return &NotificationCommandSender{cmdPublisher}
}

func (n *NotificationCommandSender) SendConfirmation(ctx context.Context, sendConf notification.SendConfirmation) error {
	body, err := notification.MarshalEnvelopeFromCommand(&sendConf)
	if err != nil {
		return err
	}
	return n.cmdPublisher.Publish(ctx, rabbitmq.SendConfirmationKey, body)
}

func (n *NotificationCommandSender) SendConfirmationSuccess(ctx context.Context, sendConfSuccess notification.SendConfirmationSuccess) error {
	body, err := notification.MarshalEnvelopeFromCommand(&sendConfSuccess)
	if err != nil {
		return err
	}

	return n.cmdPublisher.Publish(ctx, rabbitmq.SendUnsubscribeSuccessKey, body)
}

func (n *NotificationCommandSender) SendUnsubscribeSuccess(ctx context.Context, sendUnsubSuccess notification.SendUnsubscribeSuccess) error {
	body, err := notification.MarshalEnvelopeFromCommand(&sendUnsubSuccess)
	if err != nil {
		return err
	}

	return n.cmdPublisher.Publish(ctx, rabbitmq.SendUnsubscribeSuccessKey, body)
}

func (n *NotificationCommandSender) SendWeatherUpdate(ctx context.Context, sendWeatherUpd notification.SendWeatherUpdate) error {
	body, err := notification.MarshalEnvelopeFromCommand(&sendWeatherUpd)
	if err != nil {
		return err
	}

	return n.cmdPublisher.Publish(ctx, rabbitmq.SendWeatherUpdateKey, body)
}
