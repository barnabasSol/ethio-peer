package peer

import (
	"context"
	"ep-peer-service/internal/genproto/peer"
)

type GrpcHandler struct {
	peer.UnimplementedPeerServiceServer
	s Service
}

func NewGrpcHandler(s Service) *GrpcHandler {
	return &GrpcHandler{
		s: s,
	}
}

func (g *GrpcHandler) GetPeer(
	ctx context.Context,
	req *peer.GetPeerRequest,
) (*peer.GetPeerResponse, error) {
	result, err := g.s.GetPeer(ctx, req)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}
