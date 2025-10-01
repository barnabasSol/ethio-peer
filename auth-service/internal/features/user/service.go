package user

import (
	"context"
	"ep-auth-service/internal/features/common"
	"ep-auth-service/internal/genproto/user"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service interface {
	GetUser(
		ctx context.Context,
		req *user.GetUserRequest,
	) (*common.Response[user.GetUserResponse], error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}
func (s *service) GetUser(
	ctx context.Context,
	req *user.GetUserRequest,
) (*common.Response[user.GetUserResponse], error) {
	user_obj_id, err := bson.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	user_res, err := s.repo.GetUser(ctx, user_obj_id)
	if err != nil {
		return nil, err
	}

	return &common.Response[user.GetUserResponse]{
		Message: "",
		Data: user.GetUserResponse{
			UserId:         user_res.Id.Hex(),
			Username:       user_res.Username,
			Name:           user_res.Name,
			Roles:          user_res.Roles,
			Email:          user_res.Email,
			InstituteEmail: user_res.InstituteEmail,
			CreatedAt:      timestamppb.New(user_res.CreatedAt),
		},
	}, nil
}
