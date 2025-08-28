package broker

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func InitRabbitMQ() (*RabbitMQ, error) {
	uri := os.Getenv("RABBITMQ_URI")
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = NewNotificationExchange(ch)

	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn: conn,
		ch:   ch,
	}, nil
}
