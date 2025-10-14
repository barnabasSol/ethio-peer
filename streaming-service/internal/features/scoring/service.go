package scoring

import (
	"context"
	broker "ep-streaming-service/internal/broker/rabbitmq"
)

type Service interface {
	ScoreSession(context.Context, Score, string) error
}

type service struct {
	repo Repository
	br   *broker.RabbitMQ
}

func NewService(repo Repository, rmq *broker.RabbitMQ) Service {
	return &service{
		repo: repo,
		br:   rmq,
	}
}

func (s *service) ScoreSession(
	ctx context.Context,
	sc Score,
	user_id string,
) error {
	err := s.repo.InsertScore(ctx, sc, user_id)
	if err != nil {
		return err
	}
	return nil
}
