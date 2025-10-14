package scoring

type Score struct {
	SessionId string  `json:"session_id"`
	Score     float32 `json:"score"`
	Comment   string  `json:"comment"`
}
