package user

import (
	"context"
	"ep-auth-service/internal/genproto/user"
)

type GrpcHandler struct {
	user.UnimplementedUserServiceServer
	s Service
}

func NewGrpcHandler(s Service) *GrpcHandler {
	return &GrpcHandler{
		s: s,
	}
}

func (g *GrpcHandler) GetUser(
	ctx context.Context,
	req *user.GetUserRequest,
) (*user.GetUserResponse, error) {
	result, err := g.s.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}
