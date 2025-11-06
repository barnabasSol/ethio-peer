package logout

import (
	"context"
	"ep-auth-service/internal/features/common/cache"
)

type Service interface {
	LogoutUser(ctx context.Context, user_id string)
}

type service struct {
	redis cache.Redis
	r     Repository
}

func NewService() Service {
	return &service{}
}

func (s *service) LogoutUser(ctx context.Context, user_id string) {
	panic("unimplemented")
}
