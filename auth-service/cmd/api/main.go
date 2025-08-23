package main

import (
	"context"
	"ep-auth-service/internal/db"
	"ep-auth-service/internal/server"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	ctx, close := context.WithTimeout(context.Background(), 10*time.Second)
	defer close()
	mongo, err := db.NewMongoDbClient(ctx)

	if err != nil {
		log.Fatalln(err.Error())
	}

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))

	app := server.New(port, mongo)

	if err := app.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}

}
