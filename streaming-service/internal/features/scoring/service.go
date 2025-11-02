package scoring

import (
	"context"
	broker "ep-streaming-service/internal/broker/rabbitmq"
	"ep-streaming-service/internal/worker"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service interface {
	ScoreSession(context.Context, Score, string, string) error
}

type service struct {
	repo Repository
	br   *broker.RabbitMQ
}

func NewService(
	repo Repository,
	rmq *broker.RabbitMQ,
) Service {
	s := &service{
		repo: repo,
		br:   rmq,
	}

	ctx := context.Background()

	go func() {
		log.Println("score worker up")
		for job := range worker.ScoreUpdateChan {
			s.calculate_and_que(
				ctx,
				job.SessionId,
			)
		}
	}()

	return s
}

func (s *service) ScoreSession(
	ctx context.Context,
	sc Score,
	username string,
	user_id string,
) error {
	soid, err := bson.ObjectIDFromHex(sc.SessionId)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid session id",
		)
	}

	is_scored, err := s.repo.IsScored(ctx, soid, user_id)
	if err != nil {
		return err
	}

	if !is_scored {
		err := s.repo.InsertScore(ctx, sc, username, user_id)
		if err != nil {
			return err
		}
		worker.ScoreUpdateChan <- worker.ScoreUpdateJob{
			UserId:    user_id,
			SessionId: soid,
		}
		return nil
	}
	err = s.repo.UpdateScore(ctx, sc, user_id)
	if err != nil {
		return err
	}
	worker.ScoreUpdateChan <- worker.ScoreUpdateJob{
		UserId:    user_id,
		SessionId: soid,
	}
	return nil
}
