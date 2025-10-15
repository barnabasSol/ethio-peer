package worker

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ScoreUpdateJob struct {
	SessionId bson.ObjectID
}

var ScoreUpdateChan = make(chan ScoreUpdateJob, 100)
