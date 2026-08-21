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
	GracefulShutdownTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT"  envDefault:"30s"`
	HTTP                    HTTPConfig    `envPrefix:"HTTP_"`
}

type HTTPConfig struct {
	Host              string        `env:"HOST" envDefault:"0.0.0.0"`
	Port              uint16        `env:"PORT" envDefault:"8080"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"5s"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout       time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`
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

	if c.HTTP.Port == 0 {
		return errors.New("HTTP_PORT can not be 0")
	}

	if c.HTTP.ReadTimeout <= 0 {
		return errors.New("HTTP_READ_TIMEOUT must be greater than 0")
	}

	if c.HTTP.WriteTimeout <= 0 {
		return errors.New("HTTP_WRITE_TIMEOUT must be greater than 0")
	}

	if c.HTTP.ReadHeaderTimeout <= 0 {
		return errors.New("HTTP_READ_HEADER_TIMEOUT must be greater than 0")
	}

	if c.HTTP.IdleTimeout <= 0 {
		return errors.New("HTTP_IDLE_TIMEOUT must be greater than 0")
	}

	if c.GracefulShutdownTimeout <= 0 {
		return errors.New("GRACEFUL_SHUTDOWN_TIMEOUT must be greater than 0")
	}

	return nil
}
