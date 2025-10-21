package peer

type PeerResponse struct {
	UserId       string   `json:"user_id"`
	OverallScore byte     `json:"overall_score"`
	ProfilePhoto string   `json:"profile_photo"`
	OnlineStatus bool     `json:"online_status"`
	Bio          string   `json:"bio"`
	Interests    []string `json:"interests"`
}

type TopPeer struct {
	Id     string `json:"id" bson:"_id"`
	Name   string `json:"name" bson:"name"`
	Photo  string `json:"photo" bson:"profile_picture"`
	Rating string `json:"rating" bson:"overall_score"`
}
