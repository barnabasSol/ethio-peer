package resetpassword

import "ep-auth-service/internal/features/common/cache"

type Service interface {
	VerifyCredential(VerifyRequest) (string, error)
	ChangePassword(ChangePasswordRequest) error
}

type service struct {
	cache *cache.Redis
	repo  Repository
}

func NewService(repo Repository, c *cache.Redis) Service {
	return &service{
		repo:  repo,
		cache: c,
	}
}

func (s *service) VerifyCredential(req VerifyRequest) (string, error) {
	return "", nil
}

func (s *service) ChangePassword(req ChangePasswordRequest) error {
	return nil
}
