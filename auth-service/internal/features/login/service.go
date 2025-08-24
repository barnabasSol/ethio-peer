package login

import "ep-auth-service/internal/features/jwt"

type Service interface {
	LoginUser(req LoginRequest) (*LoginReponse, error)
}

type service struct {
	token_service jwt.Generator
	rep           Repository
}

func NewService(rep Repository, tg jwt.Generator) Service {
	return &service{
		rep:           rep,
		token_service: tg,
	}
}

func (s *service) LoginUser(user LoginRequest) (*LoginReponse, error) {
	if user.Username != nil {

	}
	return &LoginReponse{}, nil

}
