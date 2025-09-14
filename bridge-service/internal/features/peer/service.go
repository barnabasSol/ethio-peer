package peer

import (
	"context"
	"ep-bridge-service/internal/features/common/transport"
	"ep-bridge-service/internal/genproto/peer"
	"errors"
	"log"
)

type Service interface {
	GetPeer(ctx context.Context, user_id string) (*PeerResponse, error)
}

type service struct {
	gc *transport.GrpcClient
}

func NewService(gc *transport.GrpcClient) Service {
	return service{
		gc: gc,
	}
}

func (s service) GetPeer(ctx context.Context, user_id string) (*PeerResponse, error) {
	c := peer.NewPeerServiceClient(s.gc.Conn)
	req := &peer.GetPeerRequest{UserId: user_id}

	resp, err := c.GetPeer(ctx, req)
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to fetch peer")

	}
	return &PeerResponse{
		UserId:       resp.UserId,
		OverallScore: byte(resp.OverallScore),
		ProfilePhoto: resp.ProfilePhoto,
		OnlineStatus: resp.OnlineStatus,
		Bio:          resp.Bio,
		Interests:    resp.Interests,
	}, nil
}
