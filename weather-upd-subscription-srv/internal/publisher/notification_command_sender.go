package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-lib/pkg/message"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-lib/pkg/notification/command"
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

func (n *NotificationCommandSender) SendConfirmation(ctx context.Context, sendConf command.SendConfirmation) error {
	body, err := preparePayloadForCommand(&sendConf)
	if err != nil {
		return err
	}
	return n.cmdPublisher.Publish(ctx, rabbitmq.SendConfirmationKey, body)
}

func (n *NotificationCommandSender) SendConfirmationSuccess(ctx context.Context, sendConfSuccess command.SendConfirmationSuccess) error {
	body, err := preparePayloadForCommand(&sendConfSuccess)
	if err != nil {
		return err
	}

	return n.cmdPublisher.Publish(ctx, rabbitmq.SendUnsubscribeSuccessKey, body)
}

func (n *NotificationCommandSender) SendUnsubscribeSuccess(ctx context.Context, sendUnsubSuccess command.SendUnsubscribeSuccess) error {
	body, err := preparePayloadForCommand(&sendUnsubSuccess)
	if err != nil {
		return err
	}

	return n.cmdPublisher.Publish(ctx, rabbitmq.SendUnsubscribeSuccessKey, body)
}

func (n *NotificationCommandSender) SendWeatherUpdate(ctx context.Context, sendWeatherUpd command.SendWeatherUpdate) error {
	body, err := preparePayloadForCommand(&sendWeatherUpd)
	if err != nil {
		return err
	}

	return n.cmdPublisher.Publish(ctx, rabbitmq.SendWeatherUpdateKey, body)
}

func preparePayloadForCommand(cmd command.NotificationCommand) ([]byte, error) {
	var empty []byte

	payload, err := json.Marshal(cmd)
	if err != nil {
		return empty, fmt.Errorf("marshal command: %w", err)
	}

	env := message.Envelope{
		Type:    cmd.Type(),
		Payload: payload,
	}

	body, err := json.Marshal(env)
	if err != nil {
		return empty, fmt.Errorf("marshal envelope: %w", err)
	}

	return body, nil
}
