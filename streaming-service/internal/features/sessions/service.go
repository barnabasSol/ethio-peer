package sessions

import (
	"context"
	"ep-streaming-service/internal/features/common/livekit"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/livekit/protocol/auth"
	lk_protcol "github.com/livekit/protocol/livekit"

	lksdk "github.com/livekit/server-sdk-go/v2"
)

type Service interface {
	CreteSession(
		ctx context.Context,
		username, user_id string,
		session Create,
	) (string, error)
	EndSession(ctx context.Context, session_id, owner_id string) error
	GetSessions(ctx context.Context, filter string)
}

type service struct {
	rc     *lksdk.RoomServiceClient
	lk_cfg *livekit.Config
	repo   Repository
}

func NewService(
	repo Repository,
	cfg livekit.Config,
) Service {
	return &service{
		repo:   repo,
		lk_cfg: &cfg,
		rc: lksdk.NewRoomServiceClient(
			cfg.Host,
			cfg.ApiKey,
			cfg.ApiSecret,
		),
	}
}

func (s *service) CreteSession(
	ctx context.Context,
	username string,
	user_id string,
	session Create,
) (string, error) {
	sid, err := s.repo.InsertSession(
		ctx,
		session,
		user_id,
		username,
	)
	if err != nil {
		return "", err
	}
	_, err = s.rc.CreateRoom(
		ctx,
		&lk_protcol.CreateRoomRequest{
			Name:         sid,
			EmptyTimeout: 120,
		},
	)
	if err != nil {
		return "", echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to create stream session",
		)
	}
	grant := &auth.VideoGrant{
		RoomJoin:   true,
		Room:       sid,
		RoomCreate: true,
	}

	at := auth.NewAccessToken(
		s.lk_cfg.ApiKey,
		s.lk_cfg.ApiSecret,
	)
	at.SetVideoGrant(grant).
		SetIdentity(username).
		SetValidFor(time.Hour)

	token, err := at.ToJWT()
	if err != nil {
		return "", echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to generate stream token",
		)
	}
	return token, nil
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
