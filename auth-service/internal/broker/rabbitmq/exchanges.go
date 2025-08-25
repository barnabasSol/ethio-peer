package broker

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewNotificationExchange(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		"notification_exchange", // name
		"topic",                 // type
		true,                    // durable
		false,                   // auto-deleted
		false,                   // internal
		false,                   // no-wait
		nil,                     // arguments
	)

	return err
}
