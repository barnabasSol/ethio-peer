package otp

import (
	"context"
	"ep-auth-service/internal/features/jwt"
	"log"
)

type Service interface {
	VerifyOTP(
		ctx context.Context,
		ov OtpVerification,
	) (*OtpSuccess, error)
}

type service struct {
	repo Repository
	t    jwt.Generator
	m    *OTPManager
}

func NewService(
	m *OTPManager,
	repo Repository,
	t jwt.Generator,
) Service {
	return &service{
		repo: repo,
		t:    t,
		m:    m,
	}
}

func (s *service) VerifyOTP(
	ctx context.Context,
	ov OtpVerification,
) (*OtpSuccess, error) {
	s.m.mu.RLock()
	v, found := s.m.collection[ov.SessionId]
	s.m.mu.RUnlock()
	if !found {
		return nil, ErrInvalidOTP
	}
	if v.Value != ov.Code {
		return nil, ErrIncorrectOTP
	}

	s.m.removeOTP(ov.SessionId)

	user, err := s.repo.GetUserById(ctx, v.UserId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	s.repo.UpdateUser(ctx, user.Id, true, true)

	token, err := s.t.GenerateAccessToken(*user)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	refresh, err := s.t.GenerateRefreshToken(32)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &OtpSuccess{
		UserId:       user.Id.Hex(),
		AccessToken:  &token,
		RefreshToken: &refresh,
	}, nil
}
