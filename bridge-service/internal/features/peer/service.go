package peer

import (
	"context"
	"ep-bridge-service/internal/features/common/transport"
	"ep-bridge-service/internal/genproto/peer"
	"ep-bridge-service/internal/genproto/user"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Service interface {
	GetPeer(ctx context.Context, user_id string) (*PeerResponse, error)
	GetTopPeers(ctx context.Context) (*[]TopPeer, error)
}

type service struct {
	ugc *transport.GrpcClient
	pgc *transport.GrpcClient
}

func NewService(
	ugc *transport.GrpcClient,
	pgc *transport.GrpcClient,
) Service {
	return service{
		pgc: pgc,
		ugc: ugc,
	}
}

func (s service) GetTopPeers(ctx context.Context) (*[]TopPeer, error) {
	pc := peer.NewPeerServiceClient(s.pgc.Conn)
	uc := user.NewUserServiceClient(s.ugc.Conn)
	resp, err := pc.GetTopPeers(ctx, &peer.Empty{})

	var top_peers []TopPeer

	for _, peer := range resp.TopPeers {
		user, err := uc.GetUser(
			ctx,
			&user.GetUserRequest{UserId: peer.UserId},
		)
		if err != nil {
			log.Println(err)
			continue
		}
		top_peers = append(top_peers, TopPeer{
			Id:       peer.UserId,
			Rating:   peer.OverallScore,
			Name:     user.Name,
			Username: user.Username,
			Photo:    peer.ProfilePhoto,
		})
	}

	if err != nil {
		log.Println(err)
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed fetching top peers",
		)
	}
	log.Println(resp.TopPeers)
	return &top_peers, nil
}

func (s service) GetPeer(
	ctx context.Context,
	user_id string,
) (*PeerResponse, error) {
	c := peer.NewPeerServiceClient(s.pgc.Conn)
	req := &peer.GetPeerRequest{UserId: user_id}

	resp, err := c.GetPeer(ctx, req)

	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to fetch peer")
	}
	return &PeerResponse{
		UserId:       resp.UserId,
		OverallScore: resp.OverallScore,
		ProfilePhoto: resp.ProfilePhoto,
		OnlineStatus: resp.OnlineStatus,
		Bio:          resp.Bio,
		Interests:    resp.Interests,
	}, nil
}
