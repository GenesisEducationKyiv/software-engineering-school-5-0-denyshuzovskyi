package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQResources struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func InitRabbitMQ(url string) (*RabbitMQResources, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	res := &RabbitMQResources{
		Conn:    conn,
		Channel: ch,
	}

	return res, nil
}

func (r *RabbitMQResources) Close() error {
	// also closes channel
	err := r.Conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close RabbitMQ connection %w", err)
	}

	return nil
}
