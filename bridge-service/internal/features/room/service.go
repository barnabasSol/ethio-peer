package room

import (
	"context"
	"ep-bridge-service/internal/features/common/cache"
	"ep-bridge-service/internal/features/common/transport"
)

type Service interface {
	GetRoomContent(ctx context.Context)
}

type service struct {
	gc    *transport.GrpcClient
	cache *cache.Redis
}

func NewService(
	gc *transport.GrpcClient,
	cache *cache.Redis,
) Service {
	return &service{
		gc:    gc,
		cache: cache,
	}
}

func (s *service) GetRoomContent(ctx context.Context) {
	panic("unimplemented")
}
