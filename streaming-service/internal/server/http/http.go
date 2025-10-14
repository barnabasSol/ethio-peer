package server

import (
	"context"
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
	addr string
	db   *mongo.Client
	echo *echo.Echo
}

func New(
	addr string,
	db *mongo.Client,
) *Server {
	return &Server{
		addr: addr,
		db:   db,
		echo: echo.New(),
	}
}

func (s *Server) Run() error {
	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.Recover())

	s.echo.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	srv := &http.Server{
		Addr:         s.addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		if err := s.echo.StartServer(srv); err != nil {
			log.Fatalf(
				"failed to start the streaming server %v",
				err,
			)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGINT,
	)

	s.bootstrap()

	for _, route := range s.echo.Routes() {
		fmt.Printf(
			"%s \t %s\n",
			route.Method,
			route.Path,
		)
	}

	<-quit
	log.Println("bridge service is shutting down")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	return s.echo.Shutdown(ctx)

}
