package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"medods-test/internal/auth/types"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		pool: db,
	}
}

func (r *UserRepo) Create(ctx context.Context, user types.User) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2)`
	_, err := r.pool.Exec(ctx, query, user.Email, user.Password)
	// TODO: unique constraint error handle
	if err != nil {
		return fmt.Errorf("SQL: CreateUser: Exec(): %w", err)
	}
	return err
}

func (r *UserRepo) GetUserByID(ctx context.Context, userId int) (*types.User, error) {
	user := types.User{}
	query := `SELECT email, password
              FROM users
              WHERE user_id = $1`

	if err := r.pool.QueryRow(ctx, query, userId).Scan(
		&user.Email,
		&user.Password,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf(`SQL: GetUserById: Scan(): %w`, err)
	}
	return &user, nil
}
