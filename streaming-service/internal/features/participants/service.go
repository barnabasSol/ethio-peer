package participants

import (
	"context"
	"encoding/json"
	broker "ep-streaming-service/internal/broker/rabbitmq"
	"ep-streaming-service/internal/features/common"
	"ep-streaming-service/internal/features/common/livekit"
	"ep-streaming-service/internal/models"
	"log"
	"net/http"
	"strconv"
	"time"

	lk_protocol "github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"

	"github.com/labstack/echo/v4"
	"github.com/livekit/protocol/auth"
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
	br     *broker.RabbitMQ
	repo   Repository
}

func NewService(
	r Repository,
	rmq *broker.RabbitMQ,
	cfg *livekit.Config,
) Service {
	return &service{
		repo:   r,
		br:     rmq,
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

	sess_obj_id, err := bson.ObjectIDFromHex(session_id)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid session id",
		)
	}
	session, err := s.repo.GetSession(ctx, sess_obj_id)
	if err != nil {
		return nil, err
	}
	res, err := s.rc.ListParticipants(
		ctx,
		&lk_protocol.ListParticipantsRequest{
			Room: session_id,
		},
	)
	if err != nil {
		log.Println(err)
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed fetching participants",
		)
	}

	participants := make(
		[]Participant,
		0,
		len(res.Participants),
	)
	for _, p := range res.Participants {
		is_main := false
		var metadata struct {
			Name           string `json:"name"`
			ProfilePicture string `json:"profile_picture"`
			IsAdmin        string `json:"is_admin"`
		}
		err := json.Unmarshal(
			[]byte(p.Metadata),
			&metadata,
		)
		if err != nil {
			metadata.Name = ""
			metadata.ProfilePicture = ""
		}

		if session.Owner.Username == p.Identity {
			is_main = true
		}

		participant := Participant{
			Name:           metadata.Name,
			Username:       p.Identity,
			ProfilePicture: metadata.ProfilePicture,
			IsAnonymous:    false,
			IsMain:         is_main,
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
	soid, err := bson.ObjectIDFromHex(req.SessionId)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid id",
		)
	}
	session, err := s.repo.GetSession(ctx, soid)
	if err != nil {
		return nil, err
	}
	if session.StartsAt.After(time.Now().UTC()) {
		return nil, echo.NewHTTPError(
			http.StatusConflict,
			"session isn't on schedule yet",
		)
	}

	if session.EndedAt != nil {
		return nil, echo.NewHTTPError(
			http.StatusConflict,
			"session has ended",
		)
	}

	at := auth.NewAccessToken(
		s.lk_cfg.ApiKey,
		s.lk_cfg.ApiSecret,
	)

	var grant *auth.VideoGrant
	var metadata string

	var u *models.Participant
	if req.Username != session.Owner.Username {
		u = getParticipantByusername(
			session.Participants,
			req.Username,
		)
		if u == nil {
			err := s.repo.Insert(ctx, false, req)
			if err != nil {
				return nil, err
			}
		}
	}

	publish, subscribe, publishdata := true, true, true

	if req.Username == session.Owner.Username {
		grant = &auth.VideoGrant{
			RoomJoin:       true,
			RoomRecord:     true,
			RoomAdmin:      true,
			RoomCreate:     true,
			CanPublishData: &publish,
			Room:           req.SessionId,
		}
		metadata = `{
        "name": "` + session.Owner.Name + `",
        "is_admin": "` + "true" + `",
        "profile_picture": "` + session.Owner.ProfilePicture + `"
    }`
	} else {
		grant = &auth.VideoGrant{
			RoomJoin:       true,
			Room:           req.SessionId,
			CanPublish:     &publish,
			CanSubscribe:   &subscribe,
			CanPublishData: &publishdata,
		}
		metadata = `{
        "name": "` + req.Name + `",
        "is_admin": "` + "false" + `",
        "is_muted": "` + strconv.FormatBool(u.IsMuted) + `",
        "profile_picture": "` + req.ProfilePicture + `"
    }`

	}

	at.SetVideoGrant(grant).
		SetMetadata(metadata).
		SetIdentity(req.Username).
		SetValidFor(5 * time.Minute)

	token, err := at.ToJWT()
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to generate stream token",
		)
	}

	new_member_pl := broker.NewParticipantPayload{
		SessionId: req.SessionId,
		MemberId:  req.UserId,
	}

	new_member, err := json.Marshal(new_member_pl)

	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to marshal new member payload",
		)
	}

	s.br.Publish(broker.Message{
		Exchange: "Session_Exg",
		Topic:    "session.member.joined",
		Data:     new_member,
	})

	return &common.Response[string]{
		Message: "success",
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
