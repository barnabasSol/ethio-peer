package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

const PeerCollection = "peers"

type Peer struct {
	UserId       bson.ObjectID `bson:"_id,omitempty"`
	OverallScore byte          `bson:"overall_score"`
	ProfilePhoto string        `bson:"profile_photo"`
	OnlineStatus bool          `bson:"online_status"`
	Bio          string        `bson:"bio"`
	Interests    []string      `bson:"interests"`
}
