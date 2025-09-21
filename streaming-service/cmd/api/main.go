package main

import (
	server "ep-streaming-service/internal/server/http"
	"log"
	"os"
)

func main() {
	srv := server.New(os.Getenv("APP_PORT"))
	if err := srv.Run(); err != nil {
		log.Println(err)
	}
}
