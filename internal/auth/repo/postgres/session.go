package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"medods-test/internal/auth/types"
)

type SessionRepo struct {
	pool *pgxpool.Pool
}

func NewSessionRepo(db *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{
		pool: db,
	}
}

func (s *SessionRepo) CreateSession(ctx context.Context, session types.Session) error {
	query := `INSERT INTO sessions (id, user_uuid, refresh_token, expires_at, used)
					VALUES ($1, $2, $3, $4, $5)`
	_, err := s.pool.Exec(ctx, query, session.SessionId, session.UserId, session.RefreshToken, session.ExpiresAt, session.Used)
	if err != nil {
		return fmt.Errorf("SQL: CreateSession: Exec(): %w", err)
	}
	return nil
}

func (s *SessionRepo) GetSessionById(ctx context.Context, sessionId string) (*types.Session, error) {
	session := types.Session{}
	query := `SELECT id, user_uuid, refresh_token, expires_at, used
			  FROM sessions
	          WHERE id = $1`

	if err := s.pool.QueryRow(ctx, query, sessionId).Scan(
		&session.SessionId,
		&session.UserId,
		&session.RefreshToken,
		&session.ExpiresAt,
		&session.Used,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf(`SQL: GetUserByCreds: Scan(): %w`, err)
	}

	return &session, nil
}

func (s *SessionRepo) SetUsed(ctx context.Context, sessionId string) error {
	query := `UPDATE sessions
				SET used = true
				WHERE id = $1`
	_, err := s.pool.Exec(ctx, query, sessionId)
	if err != nil {
		return fmt.Errorf(`SQL: SetUsed: Exec(): %w`, err)
	}

	return nil
}

func (s *SessionRepo) CreateAndSetUsed(ctx context.Context, session types.Session, usedSessionId string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf(`SQL: CreateAndSetUsed: Begin(): %w`, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	setUsedQuery := `UPDATE sessions
				SET used = true
				WHERE id = $1`
	if _, err = tx.Exec(ctx, setUsedQuery, usedSessionId); err != nil {
		return fmt.Errorf(`SQL: CreateAndSetUsed: Exec(): %w`, err)
	}

	createQuery := `INSERT INTO sessions (id, user_uuid, refresh_token, expires_at, used)
					VALUES ($1, $2, $3, $4, $5)`

	_, err = tx.Exec(ctx, createQuery, session.SessionId, session.UserId, session.RefreshToken, session.ExpiresAt, session.Used)
	if err != nil {
		return fmt.Errorf(`SQL: CreateAndSetUsed: Exec(): %w`, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf(`SQL: CreateAndSetUsed: Commit(): %w`, err)
	}

	return nil
}
