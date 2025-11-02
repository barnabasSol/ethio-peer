package peer

type PeerResponse struct {
	UserId       string   `json:"user_id"`
	OverallScore string   `json:"overall_score"`
	ProfilePhoto string   `json:"profile_photo"`
	OnlineStatus bool     `json:"online_status"`
	Bio          string   `json:"bio"`
	Interests    []string `json:"interests"`
}

type TopPeer struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Rating   string `json:"rating"`
	Photo    string `json:"photo"`
}
