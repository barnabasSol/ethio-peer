package rooms

type Create struct {
	Name    string   `json:"name"`
	Subject string   `json:"subject"`
	Tags    []string `json:"tags"`
}

type Room struct {
	Name         string        `json:"name"`
	Owner        Owner         `json:"owner"`
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
