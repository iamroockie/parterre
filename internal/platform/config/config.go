package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppEnv          Env            `env:"APP_ENV,notEmpty"`
	LogLevel        slog.Level     `env:"LOG_LEVEL" envDefault:"info"`
	ShutdownTimeout time.Duration  `env:"SHUTDOWN_TIMEOUT" envDefault:"30s"`
	HTTP            HTTPConfig     `envPrefix:"HTTP_"`
	Postgres        PostgresConfig `envPrefix:"POSTGRES_"`
}

type HTTPConfig struct {
	Port         uint16        `env:"PORT" envDefault:"8080"`
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`
	ReadyTimeout time.Duration `env:"READY_TIMEOUT" envDefault:"2s"`
}

type PostgresConfig struct {
	DSN      string `env:"DSN,notEmpty"`
	MaxConns int32  `env:"MAX_CONNS" envDefault:"6"`
}

func Load() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}
