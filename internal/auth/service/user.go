package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"medods-test/internal/auth/types"
	"medods-test/pkg/auth"
	"medods-test/pkg/hash"
	"strconv"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepo interface {
	Create(ctx context.Context, user types.User) error
	GetUserByCreds(ctx context.Context, email string, password string) (*types.User, error)
	GetUserByID(ctx context.Context, userId int) (*types.User, error)
}

type User struct {
	userrepo    UserRepo
	sessionrepo SessionRepo

	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func (u *User) SignUp(ctx context.Context, input types.UserDTO) error {
	passwordHash, err := u.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	userUUID := uuid.NewString()
	user := types.User{
		UserUUID: userUUID,
		Email:    input.Email,
		Password: passwordHash,
	}

	return u.userrepo.Create(ctx, user)
}

func (u *User) SingIn(ctx context.Context, input types.UserSignInDTO, IP string) (types.Tokens, error) {
	password, err := u.hasher.Hash(input.Password)
	if err != nil {
		return types.Tokens{}, err
	}

	user, err := u.userrepo.GetUserByCreds(ctx, input.Email, password)
	if err != nil {
		return types.Tokens{}, err
	}

	if user == nil {
		return types.Tokens{}, ErrUserNotFound
	}
	return u.CreateSession(ctx, user.UserId, IP)
}

func (u *User) CreateSession(ctx context.Context, userId int, IP string) (types.Tokens, error) {
	var (
		tokens types.Tokens
		err    error
	)

	sessionId := uuid.NewString()

	tokens.AccessToken, err = u.tokenManager.NewJWT(sessionId, strconv.Itoa(userId), IP, u.accessTokenTTL)
	if err != nil {
		return tokens, err
	}

	tokens.RefreshToken, err = u.tokenManager.NewRefreshToken()
	if err != nil {
		return tokens, err
	}

	hashToken, err := u.tokenManager.HashToken(tokens.RefreshToken)
	session := types.Session{
		SessionId:    sessionId,
		UserId:       userId,
		RefreshToken: hashToken,
		ExpiresAt:    time.Now().Add(u.refreshTokenTTL),
	}

	err = u.sessionrepo.CreateSession(ctx, session)
	return tokens, err
}
