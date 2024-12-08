package service

import (
	"medods-test/pkg/auth"
	"medods-test/pkg/email"
	"medods-test/pkg/hash"
	"time"
)

const (
	salt = "asqlasaj"
)

type Repository struct {
	UserRepo    UserRepo
	SessionRepo SessionRepo
}

type Service struct {
	repository *Repository
}

func New(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) User(manager auth.TokenManager, smtp email.Sender, accessTokenTTL, refreshTokenTTL time.Duration) *User {
	return &User{
		userrepo:        s.repository.UserRepo,
		sessionrepo:     s.repository.SessionRepo,
		hasher:          hash.NewSHA1Hasher(salt),
		tokenManager:    manager,
		smtp:            smtp,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}
