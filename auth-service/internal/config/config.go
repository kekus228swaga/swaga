package config

import "github.com/caarlos0/env/v10"

type Config struct {
	Port      string `env:"PORT" envDefault:"8081"`
	DBDSN     string `env:"DB_DSN" envDefault:"postgres://postgres:postgres@localhost:5432/orderflow?sslmode=disable"`
	JWTSecret string `env:"JWT_SECRET" envDefault:"dev-secret-change-me"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
