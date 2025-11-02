package scoring

import (
	"context"
	"encoding/json"
	broker "ep-streaming-service/internal/broker/rabbitmq"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *service) calculate_and_que(
	ctx context.Context,
	soid bson.ObjectID,
) {
	scores, err := s.repo.GetScores(ctx, soid)
	if err != nil {
		log.Println(err)
		return
	}

	average := avg(scores)

	formatted := fmt.Sprintf("%.1f", average)

	err = s.repo.UpdateComputedScore(ctx, soid, formatted)

	if err != nil {
		log.Println(err)
		return
	}

	owner, err := s.repo.GetOwnerId(ctx, soid)
	if err != nil {
		log.Println(err)
		return
	}

	peer_score, err := s.repo.GetAverageSessionScore(ctx, owner)
	if err != nil {
		log.Println(err)
		return
	}

	if peer_score == nil {
		log.Println("Peer score is nil")
		return
	}
	psc := float32(*peer_score)

	score_json := broker.ScorePayload{
		UserId: owner,
		Score:  psc,
	}
	score, err := json.Marshal(&score_json)
	if err != nil {
		log.Println(err)
		return
	}

	s.br.Publish(broker.Message{
		Exchange: "score_exchange",
		Topic:    "score.new",
		Data:     score,
	})

}

func avg(scores []float32) float32 {
	var sum float32 = 0
	for _, score := range scores {
		sum += score
	}
	average := sum / float32(len(scores))
	return average
}
