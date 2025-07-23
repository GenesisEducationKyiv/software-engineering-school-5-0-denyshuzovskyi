package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-lib/pkg/message"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-lib/pkg/notification/command"
	commonerrors "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/error"
	amqp "github.com/rabbitmq/amqp091-go"
)

type NotificationSendingService interface {
	SendConfirmation(context.Context, command.SendConfirmation) error
	SendConfirmationSuccess(context.Context, command.SendConfirmationSuccess) error
	SendWeatherUpdate(context.Context, command.SendWeatherUpdate) error
	SendUnsubscribeSuccess(context.Context, command.SendUnsubscribeSuccess) error
}

type NotificationCommandConsumer struct {
	ch                         *amqp.Channel
	queue                      string
	notificationSendingService NotificationSendingService
	log                        *slog.Logger
}

func NewNotificationCommandConsumer(ch *amqp.Channel, queue string, notificationSendingService NotificationSendingService, log *slog.Logger) *NotificationCommandConsumer {
	return &NotificationCommandConsumer{ch: ch, queue: queue, notificationSendingService: notificationSendingService, log: log}
}

func (c *NotificationCommandConsumer) StartConsuming(ctx context.Context) error {
	msgs, err := c.ch.ConsumeWithContext(ctx, c.queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to start consuming queue %s: %w", c.queue, err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case msg, ok := <-msgs:
			if !ok {
				return commonerrors.ErrClosedMessageChannel
			}

			if err := c.processCommand(ctx, msg); err != nil {
				c.log.Error("failed to process message", "error", err)
			}
		}
	}
}

func (c *NotificationCommandConsumer) processCommand(ctx context.Context, msg amqp.Delivery) error {
	var env message.Envelope
	if err := json.Unmarshal(msg.Body, &env); err != nil {
		errIn := msg.Nack(false, false)
		return errors.Join(fmt.Errorf("failed to unmarshal envelope %w", err), errIn)
	}

	if err := c.dispatch(ctx, env); err != nil {
		errIn := msg.Nack(false, true)
		return errors.Join(fmt.Errorf("dispatch type %s err: %w", env.Type, err), errIn)
	}

	return msg.Ack(false)
}

func (c *NotificationCommandConsumer) dispatch(ctx context.Context, envelope message.Envelope) error {
	switch envelope.Type {
	case command.Confirmation:
		var cmd command.SendConfirmation
		if err := json.Unmarshal(envelope.Payload, &cmd); err != nil {
			return fmt.Errorf("failed to unmarshal %s: %w", envelope.Type, err)
		}
		return c.notificationSendingService.SendConfirmation(ctx, cmd)

	case command.ConfirmationSuccess:
		var cmd command.SendConfirmationSuccess
		if err := json.Unmarshal(envelope.Payload, &cmd); err != nil {
			return fmt.Errorf("failed to unmarshal %s: %w", envelope.Type, err)
		}
		return c.notificationSendingService.SendConfirmationSuccess(ctx, cmd)

	case command.WeatherUpdate:
		var cmd command.SendWeatherUpdate
		if err := json.Unmarshal(envelope.Payload, &cmd); err != nil {
			return fmt.Errorf("failed to unmarshal %s: %w", envelope.Type, err)
		}
		return c.notificationSendingService.SendWeatherUpdate(ctx, cmd)

	case command.UnsubscribeSuccess:
		var cmd command.SendUnsubscribeSuccess
		if err := json.Unmarshal(envelope.Payload, &cmd); err != nil {
			return fmt.Errorf("failed to unmarshal %s: %w", envelope.Type, err)
		}
		return c.notificationSendingService.SendUnsubscribeSuccess(ctx, cmd)

	default:
		return fmt.Errorf("%w type: %s", commonerrors.ErrUnsupportedCommand, envelope.Type)
	}
}
