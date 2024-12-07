package service

import (
	"context"
	"medods-test/internal/auth/types"
)

type SessionRepo interface {
	CreateSession(ctx context.Context, session types.Session) error
	GetSessionById(ctx context.Context, sessionId string) (*types.Session, error)
	SetUsed(ctx context.Context, sessionId string) error
	CreateAndSetUsed(ctx context.Context, session types.Session, usedSessionId string) error
}
