package sessions

type Create struct {
	Name        string   `json:"name"`
	Subject     string   `json:"subject"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type Session struct {
	Name         string        `json:"name"`
	Owner        Owner         `json:"owner"`
	Duration     string        `json:"duration"`
	Description  string        `json:"description"`
	IsLive       bool          `json:"is_live"`
	Participants []Participant `json:"participants"`
}

type Owner struct {
	Name           string `json:"name"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
}

type Participant struct {
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
}
