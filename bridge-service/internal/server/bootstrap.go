package server

import (
	"ep-bridge-service/internal/features/common/transport"
	"ep-bridge-service/internal/features/peer"
	"os"
)

func (s *Server) bootstrap() error {
	aggr_group := s.echo.Group("/aggregate")

	peer_port := os.Getenv("PEER_SERVICE_GRPC_PORT")
	s.peerClient = transport.NewGrpcClient(peer_port)

	ps := peer.NewService(s.peerClient)
	peer.InitHandler(ps, aggr_group)

	return nil
}
