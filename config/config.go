package config

import "github.com/caarlos0/env/v7"

type Config struct {
	PostgresConfig PostgresConfig `envPrefix:"PG_"`
	RabbitMqConfig RabbitMqConfig `envPrefix:"RMQ_"`
}

type PostgresConfig struct {
	Dsn string `env:"DSN" envDefault:"host=localhost user=postgres password=password dbname=balance-service port=5432 sslmode=disable"`
}

type RabbitMqConfig struct {
	Url string `env:"URL" envDefault:"amqp://guest:guest@localhost:5672/"`
}

func Read() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
