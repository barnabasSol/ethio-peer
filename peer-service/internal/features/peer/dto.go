package peer

type PeerResponse struct {
	UserId       string   `json:"user_id"`
	OverallScore byte     `json:"overall_score"`
	ProfilePhoto string   `json:"profile_photo"`
	OnlineStatus bool     `json:"online_status"`
	Bio          string   `json:"bio"`
	Interests    []string `json:"interests"`
}
