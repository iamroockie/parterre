package config

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env                     Environment   `env:"APP_ENV,notEmpty"`
	LogLevel                slog.Level    `env:"LOG_LEVEL" envDefault:"info"`
	PostgresDSN             string        `env:"POSTGRES_DSN,notEmpty"`
	GracefulShutdownTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT" envDefault:"30s"`
	HTTP                    HTTPConfig    `envPrefix:"HTTP_"`
}

func Load() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return Config{}, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}

func (c Config) validate() error {
	if err := c.Env.Validate(); err != nil {
		return err
	}

	if err := c.HTTP.validate(); err != nil {
		return err
	}

	if c.GracefulShutdownTimeout <= 0 {
		return errors.New("GRACEFUL_SHUTDOWN_TIMEOUT must be greater than 0")
	}

	return nil
}
