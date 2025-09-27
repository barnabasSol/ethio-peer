package participants

import (
	"context"
	"encoding/json"
	"ep-streaming-service/internal/features/common"
	"ep-streaming-service/internal/features/common/livekit"
	"ep-streaming-service/internal/models"
	"net/http"
	"time"

	lk "github.com/livekit/protocol/livekit"

	"github.com/labstack/echo/v4"
	"github.com/livekit/protocol/auth"
	lksdk "github.com/livekit/server-sdk-go/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service interface {
	Join(ctx context.Context, req Join) (*common.Response[string], error)
	GetParticipants(
		ctx context.Context,
		session_id string,
	) (*[]common.Response[[]Participant], error)
}

type service struct {
	lk_cfg livekit.Config
	rc     *lksdk.RoomServiceClient
	repo   Repository
}

func NewService(r Repository, cfg *livekit.Config) Service {

	return &service{
		repo:   r,
		lk_cfg: *cfg,
		rc: lksdk.NewRoomServiceClient(
			cfg.Host,
			cfg.ApiKey,
			cfg.ApiSecret,
		),
	}
}
func (s *service) GetParticipants(
	ctx context.Context,
	session_id string,
) (
	*[]common.Response[[]Participant],
	error,
) {
	res, err := s.rc.ListParticipants(ctx, &lk.ListParticipantsRequest{
		Room: session_id,
	})
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed fetching participants",
		)
	}

	participants := make([]Participant, 0, len(res.Participants))
	for _, p := range res.Participants {
		var metadata struct {
			Name           string `json:"name"`
			ProfilePicture string `json:"profile_picture"`
		}
		err := json.Unmarshal([]byte(p.Metadata), &metadata)
		if err != nil {
			metadata.Name = ""
			metadata.ProfilePicture = ""
		}

		participant := Participant{
			Name:           metadata.Name,
			Username:       p.Identity,
			ProfilePicture: metadata.ProfilePicture,
			IsAnonymous:    false,
		}
		participants = append(participants, participant)
	}

	var response []common.Response[[]Participant]
	response = append(
		response,
		common.Response[[]Participant]{Data: participants},
	)
	return &response, nil
}

func (s *service) Join(
	ctx context.Context,
	req Join,
) (*common.Response[string], error) {
	sess_obj_id, err := bson.ObjectIDFromHex(req.SessionId)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid id",
		)
	}
	session, err := s.repo.GetSession(ctx, sess_obj_id)
	if err != nil {
		return nil, err
	}

	at := auth.NewAccessToken(
		s.lk_cfg.ApiKey,
		s.lk_cfg.ApiSecret,
	)

	var grant *auth.VideoGrant
	var metadata string

	u := getParticipantByusername(session.Participants, req.Username)
	if u == nil {
		err := s.repo.Insert(ctx, false, req)
		if err != nil {
			return nil, err
		}
	}

	if u != nil && u.IsOwner {
		grant = &auth.VideoGrant{
			RoomJoin:   true,
			RoomRecord: true,
			RoomAdmin:  true,
			RoomCreate: true,
			Room:       req.SessionId,
		}
		metadata = `{
        "name": "` + req.Name + `",
        "profile_picture": "` + u.ProfilePicture + `"
    }`
	} else {
		grant = &auth.VideoGrant{
			RoomJoin: true,
			Room:     req.SessionId,
		}
		metadata = `{
        "name": "` + req.Name + `",
        "profile_picture": "` + req.ProfilePicture + `"
    }`

	}

	at.SetVideoGrant(grant).
		SetMetadata(metadata).
		SetIdentity(req.Username).
		SetValidFor(20 * time.Minute)

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

func getParticipantByusername(
	participants []models.Participant,
	username string,
) *models.Participant {
	for i := range participants {
		if participants[i].Username == username {
			return &participants[i]
		}
	}
	return nil
}
