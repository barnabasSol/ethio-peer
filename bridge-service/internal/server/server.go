package server

import (
	"context"
	"ep-bridge-service/internal/features/common/transport"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	addr       string
	echo       *echo.Echo
	peerClient *transport.GrpcClient
	userClient *transport.GrpcClient
}

func New(
	addr string,
) *Server {
	return &Server{
		addr: addr,
		echo: echo.New(),
	}
}

func (s *Server) Run() error {

	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.ContextTimeout(60 * time.Second))

	s.echo.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "I'm OK")
	})

	srv := &http.Server{
		Addr:         s.addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		if err := s.echo.StartServer(srv); err != nil {
			log.Fatalf("failed to start the bridge server %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGINT)

	s.bootstrap()

	for _, route := range s.echo.Routes() {
		fmt.Printf("%s \t %s\n", route.Method, route.Path)
	}

	<-quit
	log.Println("bridge service is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if s.peerClient != nil && s.peerClient.Conn != nil {
		s.peerClient.Conn.Close()
	}
	return s.echo.Shutdown(ctx)

}
