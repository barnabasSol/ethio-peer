package server

import (
	"ep-auth-service/internal/features/user"
	user_proto "ep-auth-service/internal/genproto/user"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr string
	db   *mongo.Client
}

func NewGrpcServer(
	addr string,
	db *mongo.Client,
) *gRPCServer {
	return &gRPCServer{
		addr,
		db,
	}
}

func (g *gRPCServer) Run() error {
	lis, err := net.Listen("tcp", g.addr)
	if err != nil {
		return err
	}
	gs := grpc.NewServer()
	ur := user.NewRepository(g.db)
	us := user.NewService(ur)
	uh := user.NewGrpcHandler(us)
	user_proto.RegisterUserServiceServer(gs, uh)

	log.Println("user grpc started")
	return gs.Serve(lis)
}
