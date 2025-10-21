package peer

import (
	"context"
	"ep-peer-service/internal/features/common"
	"ep-peer-service/internal/genproto/peer"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service interface {
	GetPeer(
		ctx context.Context,
		req *peer.GetPeerRequest,
	) (*common.Response[peer.GetPeerResponse], error)
	GetTopPeers(
		ctx context.Context,
	) (*[]TopPeer, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}
func (s *service) GetPeer(
	ctx context.Context,
	req *peer.GetPeerRequest,
) (*common.Response[peer.GetPeerResponse], error) {
	user_obj_id, err := bson.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, common.ErrInvalidPeerId
	}
	peer_res, err := s.repo.GetPeer(ctx, user_obj_id)
	if err != nil {
		return nil, err
	}

	return &common.Response[peer.GetPeerResponse]{
		Message: "",
		Data: peer.GetPeerResponse{
			UserId:       peer_res.UserId.Hex(),
			OverallScore: peer_res.OverallScore,
			ProfilePhoto: peer_res.ProfilePhoto,
			OnlineStatus: peer_res.OnlineStatus,
			Bio:          peer_res.Bio,
			Interests:    peer_res.Interests,
		},
	}, nil
}

func (s *service) GetTopPeers(ctx context.Context) (*[]TopPeer, error) {
	return s.repo.GetTopPeers(ctx)
}
