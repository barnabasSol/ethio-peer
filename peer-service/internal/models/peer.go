package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Peer struct {
	UserId       bson.ObjectID `bson:"_id,omitempty"`
	OverallScore byte          `bson:"overall_score"`
	OnlineStatus bool          `bson:"online_status"`
	Bio          string        `bson:"bio"`
	Interests    []string      `bson:"interests"`
}
