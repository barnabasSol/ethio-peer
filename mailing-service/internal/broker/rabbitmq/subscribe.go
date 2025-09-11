package broker

import amqp "github.com/rabbitmq/amqp091-go"

func (r *RabbitMQ) Subscribe(que_name, binding_key string) (<-chan amqp.Delivery, error) {
	q, err := r.ch.QueueDeclare(
		que_name,
		true,  // durable
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
		"notification_exchange",
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return r.ch.Consume(
		q.Name,
		"",
		false, // auto-ack
		false, // exclusive
		false,
		false,
		nil,
	)
}
