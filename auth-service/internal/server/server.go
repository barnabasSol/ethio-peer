package server

import (
	"context"
	broker "ep-auth-service/internal/broker/rabbitmq"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Server struct {
	addr   string
	echo   *echo.Echo
	db     *mongo.Client
	broker *broker.RabbitMQ
}

func New(
	addr string,
	db *mongo.Client,
	broker *broker.RabbitMQ,
) *Server {
	return &Server{
		addr:   addr,
		echo:   echo.New(),
		db:     db,
		broker: broker,
	}
}

func (s *Server) Run() error {

	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.Recover())

	s.echo.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	s.echo.Static("/static", "public")

	srv := &http.Server{
		Addr:         s.addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		if err := s.echo.StartServer(srv); err != nil {
			log.Fatalf("failed to start the auth server %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGINT)

	s.bootstrap()

	for _, route := range s.echo.Routes() {
		fmt.Printf("%s \t %s\n", route.Method, route.Path)
	}

	<-quit

	log.Println("auth service has gracefully shutdown")
	ctx, shutdown := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer shutdown()

	return s.echo.Shutdown(ctx)
}
