package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"medods-test/internal/auth/types"
)

// TODO: CreateSession

type SessionRepo struct {
	pool *pgxpool.Pool
}

func NewSessionRepo(db *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{
		pool: db,
	}
}

func (s *SessionRepo) CreateSession(ctx context.Context, session types.Session) error {
	query := `INSERT INTO sessions (id, user_id, refresh_token, expires_at)
					VALUES ($1, $2, $3, $4)`
	_, err := s.pool.Exec(ctx, query, session.SessionId, session.UserId, session.RefreshToken, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("SQL: CreateSession: Exec(): %w", err)
	}
	return err
}
