package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/iamroockie/parterre/internal/app/api"
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

	log := logger.New(os.Stdout, logFormatFromEnv(cfg.Env), cfg.LogLevel)
	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", cfg.HTTP.Addr())
	if err != nil {
		return err
	}

	return api.Run(ctx, cfg, log, ln)
}

func logFormatFromEnv(env config.Environment) logger.Format {
	switch env {
	case config.EnvLocal:
		return logger.FormatText
	case config.EnvProd:
		return logger.FormatJSON
	default:
		return logger.FormatJSON
	}
}
