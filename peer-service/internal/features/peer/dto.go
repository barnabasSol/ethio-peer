package peer

type PeerResponse struct {
	UserId       string   `json:"user_id"`
	OverallScore byte     `bson:"overall_score"`
	ProfilePhoto string   `bson:"profile_photo"`
	OnlineStatus bool     `bson:"online_status"`
	Bio          string   `bson:"bio"`
	Interests    []string `bson:"interests"`
}
