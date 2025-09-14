package sessions

type Service interface {
	CreteSession(username, user_id string, session Create)
	EndSession(session_id, owner_id string)
}

type service struct {
}
