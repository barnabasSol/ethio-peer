package refreshtoken

import (
	"context"
	broker "ep-auth-service/internal/broker/rabbitmq"
	"ep-auth-service/internal/features/common"
	"ep-auth-service/internal/features/common/jwt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Service interface {
	Refresh(
		ctx context.Context,
		req Request,
	) (*common.Response[Response], error)
}

type service struct {
	token_service jwt.Generator
	broker        *broker.RabbitMQ
	rep           Repository
}

func (s *service) Refresh(
	ctx context.Context,
	req Request,
) (*common.Response[Response], error) {
	rt, err := s.rep.GetRefreshToken(ctx, req)
	if err != nil {
		return nil, err
	}
	if rt.RefreshToken != req.RefreshToken {
		err := s.rep.DeleteRefreshToken(ctx, req)
		if err != nil {
			return nil, echo.NewHTTPError(
				http.StatusInternalServerError,
				"failed to delete refresh token",
			)
		}
		return nil, echo.NewHTTPError(
			http.StatusUnauthorized,
			"invalid session",
		)
	}
	u, err := s.rep.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}

	token, err := s.token_service.GenerateAccessToken(*u)
	if err != nil {
		log.Println(err)
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to generate access token",
		)
	}

	refresh, err := s.token_service.GenerateRefreshToken(32)
	if err != nil {
		log.Println(err)
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to generate refresh token",
		)
	}

	if err := s.rep.InsertRefreshToken(ctx, refresh, req); err != nil {
		return nil, err
	}
	return &common.Response[Response]{
		Message: "successfully refreshed",
		Data: Response{
			AccessToken:  token,
			RefreshToken: refresh,
		},
	}, nil
}

func NewService(
	rep Repository,
	rmq *broker.RabbitMQ,
	tg jwt.Generator,
) Service {
	return &service{
		rep:           rep,
		token_service: tg,
		broker:        rmq,
	}
}
