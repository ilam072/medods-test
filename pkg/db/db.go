package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"medods-test/internal/config"
)

func OpenDB(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	// Создаем строку подключения
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.PgUser,     // Пользователь
		cfg.PgPassword, // Пароль
		cfg.PgHost,     // Хост
		cfg.PgPort,     // Порт
		cfg.PgDatabase, // База данных
	)

	// Создаем пул соединений с базой данных
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Создаем пул с настройками
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	// Пингуем базу данных, чтобы проверить соединение
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}