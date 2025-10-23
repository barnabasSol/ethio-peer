package main

import (
	"context"
	broker "ep-auth-service/internal/broker/rabbitmq"
	"ep-auth-service/internal/db"
	"ep-auth-service/internal/server"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	ctx, close := context.WithTimeout(context.Background(), 10*time.Second)
	defer close()

	rmq, err := broker.InitRabbitMQ()

	if err != nil {
		log.Fatalln(err.Error())
	}

	mongo, err := db.NewMongoDbClient(ctx)

	if err != nil {
		log.Fatalln(err.Error())
	}

	port := os.Getenv("PORT")
	grpc_port := os.Getenv("GRPC_PORT")

	app := server.New(port, grpc_port, mongo, rmq)

	if err := app.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}

}
