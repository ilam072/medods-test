package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

// TODO: поменять user_id на user_GUID и везде поменять логику

func (r *UserRepo) Create(ctx context.Context, user types.User) error {
	query := `INSERT INTO users (user_uuid, email, password) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, query, user.UserUUID, user.Email, user.Password)
	// TODO: unique constraint error handle
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == ErrUniqueViolationCode {
				return ErrUniqueContraintFailed
			}
		}
		return fmt.Errorf("SQL: CreateUser: Exec(): %w", err)
	}
	return err
}

func (r *UserRepo) GetUserByCreds(ctx context.Context, email string, password string) (*types.User, error) {
	user := types.User{}

	query := `SELECT user_uuid, email, password 
			  FROM users
	          WHERE email = $1 AND password = $2`

	if err := r.pool.QueryRow(ctx, query, email, password).Scan(
		&user.UserUUID,
		&user.Email,
		&user.Password,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf(`SQL: GetUserByCreds: Scan(): %w`, err)
	}

	return &user, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, userUUID int) (*types.User, error) {
	user := types.User{}

	query := `SELECT user_uuid, email, password
              FROM users
              WHERE user_uuid = $1`

	if err := r.pool.QueryRow(ctx, query, userUUID).Scan(
		&user.UserUUID,
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
