package participants

import (
	"context"
	"ep-streaming-service/internal/features/common"
	"ep-streaming-service/internal/features/common/livekit"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/livekit/protocol/auth"
)

type Service interface {
	Join(ctx context.Context, req Join) (*common.Response[string], error)
}

type service struct {
	lk_cfg livekit.Config
	repo   Repository
}

// Join implements Service.
func (s *service) Join(
	ctx context.Context,
	req Join,
) (*common.Response[string], error) {
	session, err := s.repo.GetSession(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}

	at := auth.NewAccessToken(
		s.lk_cfg.ApiKey,
		s.lk_cfg.ApiSecret,
	)

	var grant *auth.VideoGrant

	if session.OwnerId == req.UserId {
		at.SetVideoGrant(grant).
			SetMetadata(`{"role":"main_streamer"}`).
			SetIdentity(req.UserId).
			SetValidFor(20 * time.Minute)
		grant = &auth.VideoGrant{
			RoomJoin:   true,
			RoomRecord: true,
			RoomAdmin:  true,
			RoomCreate: true,
			Room:       req.SessionId,
		}
	} else {
		at.SetVideoGrant(grant).
			SetMetadata(`{"role":"audience"}`).
			SetIdentity(req.UserId).
			SetValidFor(20 * time.Minute)
		grant = &auth.VideoGrant{
			RoomJoin: true,
			Room:     req.SessionId,
		}
	}

	token, err := at.ToJWT()
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to generate stream token",
		)
	}

	return &common.Response[string]{
		Message: "",
		Data:    token,
	}, nil
}

func NewService(r Repository, cfg *livekit.Config) Service {
	return &service{
		repo:   r,
		lk_cfg: *cfg,
	}
}
