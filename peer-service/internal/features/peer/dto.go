package peer

import "go.mongodb.org/mongo-driver/v2/bson"

type PeerResponse struct {
	UserId       string   `json:"user_id"`
	OverallScore byte     `json:"overall_score"`
	ProfilePhoto string   `json:"profile_photo"`
	OnlineStatus bool     `json:"online_status"`
	Bio          string   `json:"bio"`
	Interests    []string `json:"interests"`
}

type TopPeer struct {
	OID    bson.ObjectID `json:"-" bson:"_id"`
	Id     string        `json:"id"`
	Rating string        `json:"rating" bson:"overall_score"`
	Photo  string        `json:"photo" bson:"profile_photo"`
}
