package monitoring

import (
	"context"
	"ep-streaming-service/internal/features/common/livekit"
	"ep-streaming-service/internal/features/participants"
	"net/http"

	"github.com/labstack/echo/v4"
	lk_protocol "github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service interface {
	ToggleAudio(
		ctx context.Context,
		req ToggleAudio,
		by string,
	) error
}

type service struct {
	pr     participants.Repository
	lk_cfg livekit.Config
	rc     *lksdk.RoomServiceClient
}

func NewService(
	cfg *livekit.Config,
	pr participants.Repository,
) Service {

	return &service{
		pr:     pr,
		lk_cfg: *cfg,
		rc: lksdk.NewRoomServiceClient(
			cfg.Host,
			cfg.ApiKey,
			cfg.ApiSecret,
		),
	}
}

func (s *service) ToggleAudio(
	ctx context.Context,
	req ToggleAudio,
	by string,
) error {
	soid, err := bson.ObjectIDFromHex(req.SessionId)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid session id",
		)
	}
	sess, err := s.pr.GetSession(ctx, soid)
	if err != nil {
		return err
	}
	participant, err := s.pr.GetParticipantByUsername(ctx, req.Username)
	if err != nil {
		return err
	}
	if by != req.Username && sess.Owner.Username != by {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"only main streamer is allowed that action",
		)
	}
	p, err := s.rc.GetParticipant(
		ctx,
		&lk_protocol.RoomParticipantIdentity{
			Room:     req.SessionId,
			Identity: req.Username,
		},
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed fetching participant",
		)
	}

	_, err = s.rc.MutePublishedTrack(
		ctx,
		&lk_protocol.MuteRoomTrackRequest{
			Room:     req.SessionId,
			Identity: req.Username,
			TrackSid: p.Sid,
			Muted:    !participant.IsMuted,
		},
	)

	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to mute",
		)
	}
	return nil
}
