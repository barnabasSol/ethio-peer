package resetpassword

type Service interface {
	VerifyCredential(VerifyRequest) (string, error)
	ChangePassword(ChangePasswordRequest) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) VerifyCredential(req VerifyRequest) (string, error) {
	return "", nil
}

func (s *service) ChangePassword(req ChangePasswordRequest) error {
	return nil
}
