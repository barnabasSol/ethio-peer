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

func (g *GrpcHandler) GetTopPeers(
	ctx context.Context,
	req *peer.Empty,
) (*peer.TopPeersResponse, error) {
	result, err := g.s.GetTopPeers(ctx)
	if err != nil {
		return nil, err
	}
	var top_peers []*peer.TopPeer
	for _, p := range *result.Data {
		top_peers = append(
			top_peers,
			&peer.TopPeer{
				UserId:       p.Id,
				Name:         "",
				ProfilePhoto: p.Photo,
				OverallScore: p.Rating,
			},
		)
	}
	return &peer.TopPeersResponse{
		TopPeers: top_peers,
	}, nil
}
