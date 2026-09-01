package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/iamroockie/parterre/internal/platform/config"
	"github.com/iamroockie/parterre/internal/platform/httpx"
	"github.com/iamroockie/parterre/internal/platform/logger"
	"github.com/iamroockie/parterre/internal/platform/pg"
)

func main() {
	if err := run(); err != nil {
		slog.Error("run", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config failure: %w", err)
	}

	log := logger.New(os.Stdout, cfg.LogLevel).With("env", cfg.AppEnv)
	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var lc net.ListenConfig
	addr := net.JoinHostPort("0.0.0.0", strconv.FormatUint(uint64(cfg.HTTP.Port), 10))
	ln, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("listen failure: %w", err)
	}

	pool, err := pg.NewPool(ctx, cfg.Postgres.DSN, cfg.Postgres.MaxConns)
	if err != nil {
		return fmt.Errorf("postgres connect failure: %w", err)
	}
	defer pool.Close()

	deps := dependencies{
		cfg:  cfg,
		log:  log,
		pool: pool,
	}

	r := router(deps)
	srv := httpx.NewServer(r, cfg.HTTP.IdleTimeout)
	servErr := srv.Start(ln)
	log.Info("server started", "addr", ln.Addr().String())

	select {
	case err := <-servErr:
		return fmt.Errorf("server failure: %w", err)
	case <-ctx.Done():
	}

	log.Info("shutdown started")
	if err := srv.Shutdown(cfg.HTTP.ShutdownTimeout); err != nil {
		return fmt.Errorf("shutdown server failure: %w", err)
	}
	log.Info("shutdown finished")

	return nil
}
