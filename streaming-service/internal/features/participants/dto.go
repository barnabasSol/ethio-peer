package participants

type KickOut struct {
	ParticipantId string `json:"participant_id"`
}

type Join struct {
	UserId      string `json:"user_id"`
	SessionId   string `json:"session_id"`
	AsAnonymous bool   `json:"as_anonymous"`
}

type Particiant struct {
	Name           string `json:"name"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
	IsAnonymous    bool   `json:"is_anonymous"`
}

type Flag struct {
	OwnerId    string `json:"owner_id"`
	UserId     string `json:"user_id"`
	FlagStatus int    `json:"flag_status"`
}
