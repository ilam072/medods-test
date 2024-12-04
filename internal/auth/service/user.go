package service

import (
	"context"
	"errors"
	"medods-test/internal/auth/types"
	"medods-test/pkg/hash"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepo interface {
	Create(ctx context.Context, user types.User) error
	GetUserByID(ctx context.Context, userId int) (*types.User, error)
}

type User struct {
	userrepo UserRepo
	hasher   hash.PasswordHasher
}

func (u *User) SignUp(ctx context.Context, input types.UserDTO) error {
	passwordHash, err := u.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	user := types.User{
		Email:    input.Email,
		Password: passwordHash,
	}

	return u.userrepo.Create(ctx, user)
}
