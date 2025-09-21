package participants

type Join struct {
	SessionName string `json:"room_name"`
	IsAnonymous bool   `json:"is_anonymous"`
}

type KickOut struct {
	ParticipantId string `json:"participant_id"`
}

type Particiant struct {
	Name           string `json:"name"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
	IsAnonymous    bool   `json:"is_anonymous"`
}
