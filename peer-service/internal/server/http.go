package server

import (
	"context"
	broker "ep-peer-service/internal/broker/rabbitmq"
	"ep-peer-service/internal/features/common"
	"ep-peer-service/internal/features/peer"
	"ep-peer-service/internal/features/profile"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type server struct {
	addr   string
	echo   *echo.Echo
	db     *mongo.Client
	broker *broker.RabbitMQ
}

func NewHttpServer(
	addr string,
	db *mongo.Client,
	broker *broker.RabbitMQ,
) *server {
	return &server{
		addr:   addr,
		echo:   echo.New(),
		db:     db,
		broker: broker,
	}
}

func (s *server) Run() error {

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
			log.Fatalf("failed to start the peer-server %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGINT)

	cfg := common.GetMinioConfig()

	minio_client, err := minio.New(cfg.URL, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Key, cfg.Secret, ""),
		Secure: true,
	})

	pr := profile.NewRepository(s.db)
	ps := profile.NewService(minio_client, pr)
	profile.InitHandler(ps, s.echo.Group("/profile"))

	peer_repo := peer.NewRepository(s.db)
	peer_service := peer.NewService(peer_repo)
	peer.NewHandler(peer_service, *s.echo.Group("top"))

	if err != nil {
		log.Fatalln(err)
	}

	for _, route := range s.echo.Routes() {
		fmt.Printf("%s \t %s\n", route.Method, route.Path)
	}

	<-quit

	log.Println("peer-service has gracefully shutdown")
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return s.echo.Shutdown(ctx)
}
