package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"medods-test/internal/auth/repo/postgres"
	"medods-test/internal/auth/types"
	"medods-test/pkg/auth"
	"medods-test/pkg/email"
	"medods-test/pkg/hash"
	"time"
)

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrSessionNotFound         = errors.New("session not found")
	ErrRefreshTokenExpired     = errors.New("refresh token is expired")
	ErrInvalidRefreshToken     = errors.New("invalid refresh token")
	ErrRefreshTokenAlreadyUsed = errors.New("refresh token already used")
)

type UserRepo interface {
	Create(ctx context.Context, user types.User) error
	GetUserByCreds(ctx context.Context, email string, password string) (*types.User, error)
	GetUserByID(ctx context.Context, userId string) (*types.User, error)
}

type User struct {
	userrepo    UserRepo
	sessionrepo SessionRepo

	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
	stmp         email.Sender

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

	if err = u.userrepo.Create(ctx, user); err != nil {
		if errors.Is(err, postgres.ErrUniqueContraintFailed) {
			return ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (u *User) SingIn(ctx context.Context, input types.UserDTO, IP string) (types.Tokens, error) {
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
	return u.CreateSession(ctx, user.UserUUID, IP) // todo: refactor
}

func (u *User) CreateSession(ctx context.Context, userId string, IP string) (types.Tokens, error) {
	var (
		tokens types.Tokens
		err    error
	)

	sessionId := uuid.NewString()

	tokens.AccessToken, err = u.tokenManager.NewJWT(sessionId, userId, IP, u.accessTokenTTL) //todo: refactor
	if err != nil {
		return tokens, err
	}

	tokens.RefreshToken, err = u.tokenManager.NewRefreshToken()
	if err != nil {
		return tokens, err
	}

	hashToken, err := u.tokenManager.HashToken(tokens.RefreshToken)
	if err != nil {
		return tokens, err
	}
	session := types.Session{
		SessionId:    sessionId,
		UserId:       userId,
		RefreshToken: hashToken,
		ExpiresAt:    time.Now().Add(u.refreshTokenTTL),
	}

	err = u.sessionrepo.CreateSession(ctx, session)
	return tokens, err
}

func (u *User) CreateSessionAndSetUsed(ctx context.Context, userId string, IP string, usedSessionId string) (types.Tokens, error) {
	var (
		tokens types.Tokens
		err    error
	)

	sessionId := uuid.NewString()

	tokens.AccessToken, err = u.tokenManager.NewJWT(sessionId, userId, IP, u.accessTokenTTL) //todo: refactor
	if err != nil {
		return tokens, err
	}

	tokens.RefreshToken, err = u.tokenManager.NewRefreshToken()
	if err != nil {
		return tokens, err
	}

	hashToken, err := u.tokenManager.HashToken(tokens.RefreshToken)
	if err != nil {
		return tokens, err
	}
	session := types.Session{
		SessionId:    sessionId,
		UserId:       userId,
		RefreshToken: hashToken,
		ExpiresAt:    time.Now().Add(u.refreshTokenTTL),
	}

	err = u.sessionrepo.CreateAndSetUsed(ctx, session, usedSessionId)
	return tokens, err
}

func (u *User) RefreshTokens(ctx context.Context, newClientIP string, accessToken, refreshToken string) (types.Tokens, error) {
	// Парсим AccessToken и получаем SessionId
	// По этому SessionId получаем Session с захешированным RefreshToken из БД
	// Декодируем RefreshToken и сравниваем его с RefreshToken из RequestBody, для проверки, что токены взаимосвязаны
	// Проверяю RefreshToken на "был ли он использован" (used bool)
	// Если да, то отправляю ошибку ErrRefreshTokenNotValid
	// Если нет, то ставлю user = true и продолжаю код
	// Проверяю RefreshToken на "истек ли он"
	// Сравниваем новый и старый IP – при необоходимости отправляем email warning
	// Все норм: создаю новый RefreshToken в базе

	// HashAndCompare
	//todo: ПРОВЕРИТЬ ВСЕ
	sessionId, userId, oldClientIP, err := u.tokenManager.ParseToken(accessToken)
	if err != nil {
		return types.Tokens{}, err
	}

	session, err := u.sessionrepo.GetSessionById(ctx, sessionId)
	if err != nil {
		return types.Tokens{}, err
	}
	if session == nil {
		return types.Tokens{}, ErrSessionNotFound
	}

	if err = bcrypt.CompareHashAndPassword([]byte(session.RefreshToken), []byte(refreshToken)); err != nil {
		return types.Tokens{}, ErrInvalidRefreshToken
	}

	if !session.Used {
		return types.Tokens{}, ErrRefreshTokenAlreadyUsed
	}

	if session.IsRefreshTokenExpired() {
		return types.Tokens{}, ErrRefreshTokenExpired
	}

	// todo: в конце?
	if err = u.sessionrepo.SetUsed(ctx, session.SessionId); err != nil {
		return types.Tokens{}, err
	}

	tokens, err := u.CreateSessionAndSetUsed(ctx, userId, newClientIP, session.SessionId)
	if err != nil {
		return types.Tokens{}, err
	}

	if oldClientIP != newClientIP {
		user, err := u.userrepo.GetUserByID(ctx, userId)
		if err != nil {
			return types.Tokens{}, err
		}

		send := email.Send{
			Recipient: user.Email,
			Subject:   "Внимание!",
			Body: `<h1>Смена IP-адреса</h1>
<p>Мы заметили, что ваш IP-адрес изменился.</p>
<p>Если вы не осуществляли вход с этого IP, пожалуйста, свяжитесь с нашей службой поддержки или смените пароль</p>
<p>С уважением,<br>Команда поддержки</p>`,
		}
		if err = u.stmp.Send(send); err != nil {
			return types.Tokens{}, err
		}
	}

	return tokens, nil
}
