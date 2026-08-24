package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"
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

	log := logger.New(os.Stdout, cfg.LogLevel).With("env", cfg.AppEnv)
	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		stop()
	}()

	var lc net.ListenConfig
	addr := net.JoinHostPort("0.0.0.0", strconv.FormatUint(uint64(cfg.HTTP.Port), 10))
	ln, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return err
	}

	return api.Run(ctx, cfg, log, ln)
}
