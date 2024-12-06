package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"time"
)

type Config struct {
	DBConfig     DBConfig
	ServerConfig ServerConfig
	AuthConfig   AuthConfig
	SMTPConfig   SMTPConfig
}

type DBConfig struct {
	PgUser     string `env:"PGUSER"`
	PgPassword string `env:"PGPASSWORD"`
	PgHost     string `env:"PGHOST"`
	PgPort     uint16 `env:"PGPORT"`
	PgDatabase string `env:"PGDATABASE"`
	PgSSLMode  string `env:"PGSSLMODE"`
}

type ServerConfig struct {
	HTTPPort string `env:"HTTP_PORT"`
}

type AuthConfig struct {
	AccessTokenTTL  time.Duration `env:"ACCESS_TOKEN_TTL"`
	RefreshTokenTTL time.Duration `env:"REFRESH_TOKEN_TTL"`
	SigningKey      string        `env:"SIGNING_KEY"`
}

type SMTPConfig struct {
	Host string `env:"STMP_HOST"`
	Port string `env:"STMP_PORT"`
	From string `env:"STMP_FROM"`
}

func (s *ServerConfig) Address() string {
	return fmt.Sprintf("localhost:%s", s.HTTPPort)
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from enviroment variables: %w", err)
	}

	return cfg, nil
}
