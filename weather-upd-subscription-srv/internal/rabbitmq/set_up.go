package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SetUpExchange(ch *amqp.Channel, exchangeName string) error {
	err := ch.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("exchange declare: %w", err)
	}

	return nil
}
