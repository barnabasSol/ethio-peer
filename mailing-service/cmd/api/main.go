package main

import (
	"bytes"
	"encoding/json"
	broker "ep-mailing-service/internal/broker/rabbitmq"
	"ep-mailing-service/internal/mail"
	"log"
	"os"
)

func main() {
	mailtrap_host := os.Getenv("HOST")
	mailtrap_token := os.Getenv("TOKEN")
	sender := os.Getenv("SENDER")

	mailing_service := mail.NewService(
		mailtrap_host,
		mailtrap_token,
		sender,
	)
	rmq, err := broker.InitRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := rmq.Subscribe("mailer_que", "email.*")

	if err != nil {
		log.Fatal("failed to subscribe:", err)
	}

	log.Println("Mailer service is waiting for messages...")

	for msg := range msgs {
		log.Printf("Received [%s]: %s", msg.RoutingKey, msg.Body)

		switch msg.RoutingKey {
		case "email.otp":
			var payload broker.OtpPayload
			decoder := json.NewDecoder(bytes.NewReader(msg.Body))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&payload); err != nil {
				msg.Nack(false, false)
				log.Println("Failed to decode OTP payload:", err)
				continue
			}
			if err := mailing_service.SendOTP(payload); err != nil {
				log.Println(err)
				msg.Nack(false, false)
				continue
			}
			log.Println(payload)
			msg.Ack(false)

		case "email.welcome":
			var payload broker.WelcomePayload
			decoder := json.NewDecoder(bytes.NewReader(msg.Body))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&payload); err != nil {
				msg.Nack(false, false)
				log.Println("Failed to decode Welcome payload:", err)
				continue
			}
			if err := mailing_service.SendWelcome(payload); err != nil {
				log.Println(err)
				msg.Nack(false, false)
				continue
			}
			msg.Ack(false)

		default:
			log.Printf("Unknown email event: %s", msg.RoutingKey)
		}
	}

}
