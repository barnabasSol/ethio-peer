package broker

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RabbitMQ struct {
	db   *mongo.Client
	conn *amqp.Connection
	ch   *amqp.Channel
}

func InitRabbitMQ(db *mongo.Client) (*RabbitMQ, error) {
	uri := os.Getenv("RABBITMQ_URI")
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = NewPeerExchange(ch)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn: conn,
		db:   db,
		ch:   ch,
	}, nil
}

func (r *RabbitMQ) Publish(msg Message) error {
	return r.ch.Publish(
		msg.Exchange,
		msg.Topic,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg.Data,
		},
	)
}
