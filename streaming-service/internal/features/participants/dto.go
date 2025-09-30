package participants

type KickOut struct {
	ParticipantId string `json:"participant_id"`
}

type Join struct {
	Name           string `json:"name"`
	Username       string `json:"username"`
	SessionId      string `json:"session_id"`
	ProfilePicture string `json:"profile_picture"`
	AsAnonymous    bool   `json:"as_anonymous"`
}

type Participant struct {
	Name           string `json:"name"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
	IsAnonymous    bool   `json:"is_anonymous"`
	IsMain         bool   `json:"is_main"`
}

type Flag struct {
	OwnerId    string `json:"owner_id"`
	UserId     string `json:"user_id"`
	FlagStatus int    `json:"flag_status"`
}
