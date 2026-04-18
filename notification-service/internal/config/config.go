package config

import "github.com/caarlos0/env/v10"

type Config struct {
	RabbitMQURL string `env:"RABBITMQ_URL" envDefault:"amqp://guest:guest@rabbitmq:5672/"`
	QueueName   string `env:"QUEUE_NAME" envDefault:"order.created"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
