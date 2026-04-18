package config

import "github.com/caarlos0/env/v10"

type Config struct {
	KafkaBrokers string `env:"KAFKA_BROKERS" envDefault:"kafka:9092"`
	Topic        string `env:"KAFKA_TOPIC" envDefault:"order.events"`
	GroupID      string `env:"GROUP_ID" envDefault:"analytics-group"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
