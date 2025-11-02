package server

import (
	"ep-bridge-service/internal/features/common/cache"
	"ep-bridge-service/internal/features/common/transport"
	"ep-bridge-service/internal/features/peer"
	"ep-bridge-service/internal/features/room"
	"ep-bridge-service/internal/features/user"
	"log"
	"os"
	"time"
)

func (s *Server) bootstrap() error {
	aggr_group := s.echo.Group("/aggregate")

	peer_port := os.Getenv("PEER_SERVICE_GRPC_PORT")
	s.peerClient = transport.NewGrpcClient("peer-service" + peer_port)

	auth_port := os.Getenv("AUTH_SERVICE_GRPC_PORT")
	s.userClient = transport.NewGrpcClient("auth-service" + auth_port)

	ps := peer.NewService(s.userClient, s.peerClient)
	peer.InitHandler(ps, aggr_group)

	resource_port := os.Getenv("RESOURCE_SERVICE_GRPC_PORT")
	s.resourceClient = transport.NewGrpcClient("resource-service" + resource_port)

	redis_addr := os.Getenv("REDIS_ADDR")

	cache, err := cache.New(
		redis_addr,
		30*time.Minute,
	)

	if err != nil {
		log.Println("failed to init redis")
	}

	auth_group := s.echo.Group("/auth")
	us := user.NewService(
		s.userClient,
		s.peerClient,
		cache,
	)
	user.InitHandler(us, auth_group)

	resource_group := s.echo.Group("/resource")
	rs := room.NewService(s.resourceClient, cache)
	room.InitHandler(rs, resource_group)

	return nil
}
