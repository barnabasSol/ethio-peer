package broker

import amqp "github.com/rabbitmq/amqp091-go"

func NewUserExchange(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		"new_user_exchange",
		"topic",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)

	return err
}

func NewResourceExchange(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		"new_resource_exchange",
		"topic",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)

	return err
}
