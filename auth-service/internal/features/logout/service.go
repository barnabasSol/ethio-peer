package logout

import "context"

type Service interface {
	LogoutUser(ctx context.Context, user_id string)
}

type service struct {
	r Repository
}

func NewService() Service {
	return &service{}
}

func (s *service) LogoutUser(ctx context.Context, user_id string) {
	panic("unimplemented")
}
