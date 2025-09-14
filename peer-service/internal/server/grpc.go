package server

import (
	broker "ep-peer-service/internal/broker/rabbitmq"
	"ep-peer-service/internal/features/peer"
	proto_peer "ep-peer-service/internal/genproto/peer"
	"net"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr   string
	db     *mongo.Client
	broker *broker.RabbitMQ
}

func NewGrpcServer(
	addr string,
	db *mongo.Client,
	broker *broker.RabbitMQ,
) *gRPCServer {
	return &gRPCServer{
		addr,
		db,
		broker,
	}
}

func (g *gRPCServer) Run() error {
	lis, err := net.Listen("tcp", g.addr)
	if err != nil {
		return err
	}
	gs := grpc.NewServer()
	pr := peer.NewRepository(g.db)
	ps := peer.NewService(pr)
	ph := peer.NewGrpcHandler(ps)
	proto_peer.RegisterPeerServiceServer(gs, ph)

	return gs.Serve(lis)
}
