package analytics

type SessionAnalytics struct {
	Weekly *[]Weekly `json:"weekly,omitempty"`
}

type Weekly struct {
	TopTopicParticipantCount int    `json:"top_topic_participant_count" bson:"top_topic_participant_count"`
	SessionCount             int    `json:"session_count" bson:"session_count"`
	ParticipantCount         int    `json:"participant_count" bson:"participant_count"`
	TopTopic                 Topic  `json:"top_topic" bson:"top_topic"`
	CreatedAt                string `json:"created_at" bson:"created_at"`
}

type Topic struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}
