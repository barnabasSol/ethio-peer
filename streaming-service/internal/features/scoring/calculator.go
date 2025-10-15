package scoring

import (
	"context"
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

	var sum float32 = 0
	for _, score := range scores {
		sum += score
	}
	average := sum / float32(len(scores))

	formatted := fmt.Sprintf("%.1f", average)
	fmt.Println(formatted)

}
