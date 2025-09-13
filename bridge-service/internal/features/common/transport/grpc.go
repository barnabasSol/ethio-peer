package transport

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	Conn *grpc.ClientConn
}

func NewGrpcClient(addr string) *GrpcClient {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	return &GrpcClient{
		Conn: conn,
	}
}
