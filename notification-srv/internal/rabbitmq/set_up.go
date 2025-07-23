package rabbitmq

import (
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-lib/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SetUpQueue(ch *amqp.Channel, exchangeName string, queueName string) error {
	queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	keys := []string{
		rabbitmq.SendConfirmationKey,
		rabbitmq.SendConfirmationSuccessKey,
		rabbitmq.SendWeatherUpdateKey,
		rabbitmq.SendUnsubscribeSuccessKey,
	}

	for _, key := range keys {
		err = ch.QueueBind(queue.Name, key, exchangeName, false, nil)
		if err != nil {
			return fmt.Errorf("failed to bind queue to exchange: %w", err)
		}
	}

	return nil
}
