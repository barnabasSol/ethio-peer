package monitoring

type ToggleAudio struct {
	Username  string `json:"username"`
	SessionId string `json:"session_id"`
	IsMuted   bool   `json:"is_muted"`
}
