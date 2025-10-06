package broker

import amqp "github.com/rabbitmq/amqp091-go"

func NewNotificationExchange(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		"notification_exchange",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	return err
}
