package analytics

type SessionAnalytics struct {
	Weekly *[]Weekly `json:"weekly,omitempty"`
	Hourly *[]Hourly `json:"hourly,omitempty"`
}

type Hourly struct {
	SessionCount     int    `json:"session_count" bson:"sessions_created"`
	ParticipantCount int    `json:"participant_count" bson:"total_participants"`
	Hour             string `json:"hour" bson:"hour"`
}

type Weekly struct {
	SessionCount     int    `json:"session_count" bson:"sessions_created"`
	ParticipantCount int    `json:"participant_count" bson:"total_particpants"`
	Date             string `json:"created_at" bson:"date"`
}

// {
//   "date": "2025-10-19",
//   "total_participants": 3,
//   "sessions_created": 4
// }
