package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	channel  *amqp.Channel
	exchange string
}

func NewPublisher(ch *amqp.Channel, exchange string) *Publisher {
	return &Publisher{
		channel:  ch,
		exchange: exchange,
	}
}

func (p *Publisher) Publish(ctx context.Context, routingKey RoutingKey, body []byte) error {
	return p.channel.PublishWithContext(
		ctx,
		p.exchange,
		string(routingKey),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
