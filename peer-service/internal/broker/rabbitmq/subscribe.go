package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"ep-peer-service/internal/db"
	"ep-peer-service/internal/models"
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *RabbitMQ) SubscribeNewScore(que_name, binding_key string) (<-chan amqp.Delivery, error) {
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
		"score_exchange",
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return r.ch.Consume(
		q.Name,
		"score-consumer",
		false, // auto-ack
		false, // exclusive
		false,
		false,
		nil,
	)
}

func (r *RabbitMQ) SubscribeNewPeer(que_name, binding_key string) (<-chan amqp.Delivery, error) {
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

func (r *RabbitMQ) ListenNewScore(msgs <-chan amqp.Delivery) {
	log.Println("peer service is waiting for new peer score messages...")
	for msg := range msgs {
		switch msg.RoutingKey {
		case "score.new":
			var score ScorePayload
			decoder := json.NewDecoder(bytes.NewReader(msg.Body))
			if err := decoder.Decode(&score); err != nil {
				msg.Nack(false, false)
				log.Println("failed to unmarshall new peer score payload:", err)
				continue
			}

			collection := r.db.Database(db.Name).Collection(models.PeerCollection)
			oid, err := bson.ObjectIDFromHex(score.UserId)
			if err != nil {
				log.Println("invalid object id")
				return
			}
			str_score := strconv.FormatFloat(float64(score.Score), 'f', 2, 32)

			_, err = collection.UpdateOne(
				context.Background(),
				bson.D{{Key: "_id", Value: oid}},
				bson.D{
					{Key: "$set", Value: bson.D{
						{Key: "overall_score", Value: str_score},
						{Key: "updated_at", Value: time.Now().UTC()},
					}},
				},
			)
			log.Println(score)
			if err != nil {
				log.Println("failed to update peer score:", err)
			}
		}
	}
}

func (r *RabbitMQ) ListenNewPeer(msgs <-chan amqp.Delivery) {
	log.Println("peer service is waiting for new peer messages...")
	for msg := range msgs {
		log.Printf(
			"Received [%s]: %s",
			msg.RoutingKey,
			msg.Body,
		)
		switch msg.RoutingKey {
		case "peer.new":
			var payload PeerPayload
			decoder := json.NewDecoder(bytes.NewReader(msg.Body))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&payload); err != nil {
				msg.Nack(false, false)
				log.Println(
					"failed to unmarshall new peer payload with strict checking:",
					err,
				)
				continue
			}
			obj_id, err := bson.ObjectIDFromHex(payload.UserId)
			if err != nil {
				msg.Nack(false, false)
				continue
			}
			collection := r.db.Database(db.Name).Collection(models.PeerCollection)
			zero := "0"
			result, err := collection.InsertOne(
				context.Background(),
				models.Peer{
					UserId:       obj_id,
					OverallScore: zero,
					OnlineStatus: false,
					Bio:          payload.Bio,
					Interests:    payload.Interests,
					UpdatedAt:    time.Now().UTC(),
				},
			)
			if err != nil || !result.Acknowledged {
				msg.Nack(false, false)
				continue
			}
			msg.Ack(false)
		default:
			log.Printf("Unknown peer event: %s", msg.RoutingKey)
		}
	}
}
