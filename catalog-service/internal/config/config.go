package config

import "github.com/caarlos0/env/v10"

type Config struct {
	Port     string `env:"PORT" envDefault:"8083"`
	MongoURI string `env:"MONGO_URI" envDefault:"mongodb://mongodb:27017/"`
	DBName   string `env:"MONGO_DB" envDefault:"catalog"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
