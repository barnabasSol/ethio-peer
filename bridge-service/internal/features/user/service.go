package user

import (
	"context"
	"ep-bridge-service/internal/features/common/cache"
	"ep-bridge-service/internal/features/common/transport"
	"ep-bridge-service/internal/genproto/peer"
	"ep-bridge-service/internal/genproto/user"
	"errors"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

type Service interface {
	GetCurrentUser(
		ctx context.Context,
		user_id string,
	) (*CurrentUser, error)
}

type service struct {
	userGrpcClient *transport.GrpcClient
	peerGrpcClient *transport.GrpcClient
	cache          *cache.Redis
}

func NewService(
	ugrpc *transport.GrpcClient,
	pgrpc *transport.GrpcClient,
	cache *cache.Redis,
) Service {
	return &service{
		userGrpcClient: ugrpc,
		peerGrpcClient: pgrpc,
		cache:          cache,
	}
}

func (s *service) GetCurrentUser(
	ctx context.Context,
	user_id string,
) (*CurrentUser, error) {
	cache_key := "user:" + user_id
	var cached_user CurrentUser
	err := s.cache.Get(
		ctx,
		cache_key,
		&cached_user,
	)
	if err != nil {
		log.Println("failed fetching from cache", err)
	}
	if cached_user.UserId != "" {
		return &cached_user, nil
	}

	c := peer.NewPeerServiceClient(s.peerGrpcClient.Conn)
	u := user.NewUserServiceClient(s.userGrpcClient.Conn)

	peer_req := &peer.GetPeerRequest{
		UserId: user_id,
	}
	user_req := &user.GetUserRequest{
		UserId: user_id,
	}

	var (
		user_resp *user.GetUserResponse
		peer_resp *peer.GetPeerResponse
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		resp, err := u.GetUser(ctx, user_req)
		if err != nil {
			log.Println(err)
			return errors.New("failed to fetch user")
		}
		user_resp = resp
		return nil
	})

	g.Go(func() error {
		resp, err := c.GetPeer(ctx, peer_req)
		if err != nil {
			log.Println(err)
			return errors.New("failed to fetch peer")
		}
		peer_resp = resp
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	result := &CurrentUser{
		UserId:         user_id,
		Username:       user_resp.Username,
		Name:           user_resp.Name,
		InstituteEmail: user_resp.InstituteEmail,
		Email:          user_resp.Email,
		OverallScore:   byte(peer_resp.OverallScore),
		ProfilePhoto:   peer_resp.ProfilePhoto,
		OnlineStatus:   peer_resp.OnlineStatus,
		Bio:            peer_resp.Bio,
		Roles:          user_resp.Roles,
		Interests:      peer_resp.Interests,
		CreatedAt:      user_resp.CreatedAt.AsTime(),
	}

	go func(u *CurrentUser) {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			1*time.Second,
		)
		defer cancel()

		if err := s.cache.SetWithTTL(
			ctx,
			cache_key,
			u,
			10*time.Minute,
		); err != nil {
			log.Printf("failed to cache user %s: %v", u.UserId, err)
		}
	}(result)

	return result, nil
}
