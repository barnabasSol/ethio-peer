package main

import (
	"context"
	broker "ep-peer-service/internal/broker/rabbitmq"
	"ep-peer-service/internal/db"
	"ep-peer-service/internal/server"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongo, err := db.NewMongoDbClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	rmq, err := broker.InitRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	grpc_port := os.Getenv("GRPC_PORT")
	gRPC := server.NewGrpcServer(grpc_port, mongo, rmq)
	go gRPC.Run()

	msgs, err := rmq.Subscribe("new_peer_que", "peer.*")
	if err != nil {
		log.Fatal(err)
	}
	go rmq.Listen(msgs)

	http_port := os.Getenv("PORT")
	srv := server.NewHttpServer(http_port, mongo, rmq)
	if err := srv.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
