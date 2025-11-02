package sessions

import (
	"context"
	"encoding/json"
	broker "ep-streaming-service/internal/broker/rabbitmq"
	"ep-streaming-service/internal/features/common"
	"ep-streaming-service/internal/features/common/livekit"
	"ep-streaming-service/internal/features/common/pagination"
	"log"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	lk_protocol "github.com/livekit/protocol/livekit"

	lksdk "github.com/livekit/server-sdk-go/v2"
)

type Service interface {
	CreateSession(
		ctx context.Context,
		username, user_id string,
		session Create,
	) (*common.Response[CreateResponse], error)
	EndSession(
		ctx context.Context,
		session_id, owner_username string,
	) error
	UpdateSession(
		ctx context.Context,
		req Update,
		username string,
	) error
	GetSessions(
		ctx context.Context,
		pagination pagination.Pagination,
		filter string,
	) (*common.Response[[]Session], error)
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
			"you don't own the session",
		)
	}
	err = s.repo.UpdateSession(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) CreateSession(
	ctx context.Context,
	username, user_id string,
	session Create,
) (*common.Response[CreateResponse], error) {
	sid, err := s.repo.InsertSession(
		ctx,
		session,
		username,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	_, err = s.rc.CreateRoom(
		ctx,
		&lk_protocol.CreateRoomRequest{
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
	new_room := broker.NewSessionPayload{
		OwnerId:   user_id,
		SessionId: sid,
		UserName:  username,
		TopicId:   session.Topic.Id,
	}

	new_room_json, err := json.Marshal(new_room)

	if err != nil {
		log.Println(err)
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to marshal session payload",
		)
	}

	err = s.rmq.Publish(broker.Message{
		Exchange: "Session_Exg",
		Topic:    "session.created",
		Data:     new_room_json,
	})

	if err != nil {
		log.Println("failed to publish event to resource service")
		log.Println(err)
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
	pagination pagination.Pagination,
	req string,
) (*common.Response[[]Session], error) {

	res, err := s.repo.GetSessions(ctx, pagination, req)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	for i := range *res {
		(*res)[i].SessionId = (*res)[i].Id.Hex()
		sess_id := (*res)[i].SessionId
		if (*res)[i].EndedAt != nil {
			(*res)[i].Duration = (*res)[i].EndedAt.Sub((*res)[i].StartsAt).String()
		}
		i := i
		wg.Go(func() {
			lk_p, err := s.rc.ListParticipants(
				ctx,
				&lk_protocol.ListParticipantsRequest{
					Room: sess_id,
				},
			)
			if err != nil {
				log.Println("ListParticipants error:", err)
				return
			}

			(*res)[i].Participants = []Participant{}

			unique_participants := make(map[string]Participant)

			isLive := false

			for _, lkp := range lk_p.Participants {
				var metadata struct {
					Name           string `json:"name"`
					ProfilePicture string `json:"profile_picture"`
					IsAdmin        string `json:"is_admin"`
				}

				err := json.Unmarshal([]byte(lkp.Metadata), &metadata)
				if err != nil {
					log.Println("Metadata unmarshal error:", err)
					continue
				}

				if _, exists := unique_participants[lkp.Identity]; !exists {
					unique_participants[lkp.Identity] = Participant{
						Username:       lkp.Identity,
						ProfilePicture: metadata.ProfilePicture,
					}
					if metadata.IsAdmin == "true" {
						isLive = true
					}
				}
			}

			for _, p := range unique_participants {
				(*res)[i].Participants = append((*res)[i].Participants, p)
			}

			(*res)[i].IsLive = isLive
		})
	}
	wg.Wait()

	return &common.Response[[]Session]{
		Data: *res,
	}, nil
}

func (s *service) EndSession(
	ctx context.Context,
	session_id, owner_username string,
) error {
	ok, err := s.repo.IsOwner(
		ctx,
		session_id,
		owner_username,
	)
	if err != nil {
		return err
	}
	if !ok {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"you can't end this session",
		)
	}
	_, err = s.rc.DeleteRoom(ctx, &lk_protocol.DeleteRoomRequest{
		Room: session_id,
	})
	if err != nil {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"failed to end session please try again later",
		)
	}
	yes := true
	err = s.repo.UpdateSession(ctx, Update{
		SessionId: session_id,
		IsEnded:   &yes,
	})
	if err != nil {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"failed to end session please try again later",
		)
	}
	return nil
}
