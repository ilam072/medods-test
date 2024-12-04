package service

import (
	"context"
	"medods-test/internal/auth/types"
)

type SessionRepo interface {
	CreateSession(ctx context.Context, session types.Session) error
}
