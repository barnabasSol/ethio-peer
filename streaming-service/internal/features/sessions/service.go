package sessions

import (
	"context"
	broker "ep-streaming-service/internal/broker/rabbitmq"
	"ep-streaming-service/internal/features/common"
	"ep-streaming-service/internal/features/common/livekit"
	"net/http"

	"github.com/labstack/echo/v4"
	lk_protcol "github.com/livekit/protocol/livekit"

	lksdk "github.com/livekit/server-sdk-go/v2"
)

type Service interface {
	CreteSession(
		ctx context.Context,
		username, user_id string,
		session Create,
	) (*common.Response[CreateResponse], error)
	EndSession(ctx context.Context, session_id, owner_id string) error
	UpdateSession(ctx context.Context, req Update, username string) error
	GetSessions(ctx context.Context, filter string)
}

type service struct {
	rc     *lksdk.RoomServiceClient
	lk_cfg *livekit.Config
	rmq    *broker.RabbitMQ
	repo   Repository
}

func NewService(
	repo Repository,
	rmq *broker.RabbitMQ,
	cfg livekit.Config,
) Service {
	return &service{
		repo:   repo,
		lk_cfg: &cfg,
		rmq:    rmq,
		rc: lksdk.NewRoomServiceClient(
			cfg.Host,
			cfg.ApiKey,
			cfg.ApiSecret,
		),
	}
}

func (s *service) UpdateSession(
	ctx context.Context,
	req Update,
	username string,
) error {
	ok, err := s.repo.IsOwner(ctx, req.SessionId, username)
	if err != nil {
		return err
	}
	if !ok {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"you don't own the room",
		)
	}
	err = s.repo.UpdateSession(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) CreteSession(
	ctx context.Context,
	username, user_id string,
	session Create,
) (*common.Response[CreateResponse], error) {
	sid, err := s.repo.InsertSession(
		ctx,
		session,
		username,
	)
	if err != nil {
		return nil, err
	}
	_, err = s.rc.CreateRoom(
		ctx,
		&lk_protcol.CreateRoomRequest{
			Name:         sid,
			EmptyTimeout: 120,
		},
	)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to create stream session",
		)
	}
	return &common.Response[CreateResponse]{
		Message: "session created",
		Data: CreateResponse{
			RoomId: sid,
		},
	}, nil
}

func (s *service) GetSessions(
	ctx context.Context,
	filter string,
) {
	panic("unimplemented")
}

func (s *service) EndSession(
	ctx context.Context,
	owner_id, session_id string,
) error {

	panic("unimplemented")
}
