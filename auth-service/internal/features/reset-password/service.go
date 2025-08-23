package changepassword

type Service interface {
	ChangePassword(user_id string, password_hash string)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) ChangePassword(user_id string, password_hash string) {

}
