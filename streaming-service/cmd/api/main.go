package main

import (
	"context"
	"ep-streaming-service/internal/db"
	server "ep-streaming-service/internal/server/http"
	"log"
	"os"
)

func main() {
	mongo, err := db.NewMongoDbClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	srv := server.New(os.Getenv("APP_PORT"), mongo)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}
}
