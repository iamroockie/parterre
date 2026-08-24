package main

import (
	"log/slog"
	"os"

	"github.com/iamroockie/parterre/internal/platform/config"
	"github.com/iamroockie/parterre/internal/platform/logger"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed run", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log := logger.New(os.Stdout, cfg.LogLevel).With("env", cfg.AppEnv)
	slog.SetDefault(log)

	return nil
}
