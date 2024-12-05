package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

var (
	ErrEmptySigningKey = errors.New("signing key is empty")
)

type TokenClaims struct {
	jwt.RegisteredClaims
	IP        string `json:"ip"`
	SessionId string
}

type TokenManager interface {
	NewJWT(sessionId, userId, userIP string, ttl time.Duration) (string, error)
	ParseToken(accessToken string) (string, string, string, error)
	NewRefreshToken() (string, error)
	HashToken(refreshToken string) (string, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, ErrEmptySigningKey
	}
	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(sessionId, userId, userIP string, ttl time.Duration) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userId,
		},
		IP:        userIP,
		SessionId: sessionId,
	})

	return jwtToken.SignedString([]byte(m.signingKey))
}

func (m *Manager) ParseToken(accessToken string) (string, string, string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})

	if err != nil {
		return "", "", "", err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return "", "", "", fmt.Errorf("error get token claims")
	}

	return claims.SessionId, claims.Subject, claims.IP, err
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	refreshToken := base64.StdEncoding.EncodeToString(b)
	return refreshToken, nil
}

func (m *Manager) HashToken(refreshToken string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
