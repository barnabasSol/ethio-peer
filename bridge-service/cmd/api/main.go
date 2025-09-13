package main

import (
	"ep-bridge-service/internal/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	app_port := os.Getenv("APP_PORT")
	s := server.New(app_port)
	if err := s.Run(); err != nil {
		log.Println(err)
	}
}
