package broker

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) Subscribe(que_name, binding_key string) (<-chan amqp.Delivery, error) {
	q, err := r.ch.QueueDeclare(
		que_name,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}
	err = r.ch.QueueBind(
		q.Name,
		binding_key,
		"new_peer_exchange",
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return r.ch.Consume(
		q.Name,
		"peer-consumer",
		false, // auto-ack
		false, // exclusive
		false,
		false,
		nil,
	)
}

func (r *RabbitMQ) Listen(msgs <-chan amqp.Delivery) {
	log.Println("Mailer service is waiting for messages...")
	for msg := range msgs {
		log.Printf("Received [%s]: %s", msg.RoutingKey, msg.Body)
		switch msg.RoutingKey {
		case "peer.new":
			var payload PeerPayload
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				msg.Nack(false, true)
				log.Println("failed to unmarshall new peer payload")
				continue
			}
			log.Println(payload)
			msg.Ack(false)
		default:
			log.Printf("⚠️ Unknown peer event: %s", msg.RoutingKey)
		}
	}
}
