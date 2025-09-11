package broker

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewPeerExchange(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		"new_peer_exchange", // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)

	return err
}
