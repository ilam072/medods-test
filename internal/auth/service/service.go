package service

type Repository struct {
	UserRepo UserRepo
}

type Service struct {
	repository *Repository
}

func New(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) User() *User {
	return &User{
		userrepo: s.repository.UserRepo,
	}
}
