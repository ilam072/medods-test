package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig     DBConfig
	ServerConfig ServerConfig
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
