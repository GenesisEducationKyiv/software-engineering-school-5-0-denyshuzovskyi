package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-lib/pkg/message"
	commonerrors "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/error"
	amqp "github.com/rabbitmq/amqp091-go"
)

type CommandDispatcher interface {
	Dispatch(context.Context, message.Envelope) error
}

type NotificationCommandConsumer struct {
	ch                *amqp.Channel
	queue             string
	commandDispatcher CommandDispatcher
	log               *slog.Logger
}

func NewNotificationCommandConsumer(ch *amqp.Channel, queue string, commandDispatcher CommandDispatcher, log *slog.Logger) *NotificationCommandConsumer {
	return &NotificationCommandConsumer{
		ch:                ch,
		queue:             queue,
		commandDispatcher: commandDispatcher,
		log:               log,
	}
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

	if err := c.commandDispatcher.Dispatch(ctx, env); err != nil {
		errIn := msg.Nack(false, true)
		return errors.Join(fmt.Errorf("dispatch type %s err: %w", env.Type, err), errIn)
	}

	return msg.Ack(false)
}
