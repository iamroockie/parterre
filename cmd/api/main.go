package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/iamroockie/parterre/internal/platform/config"
	"github.com/iamroockie/parterre/internal/platform/httpx"
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

	logFormat := logFormatFromEnv(cfg.Env)
	log := logger.New(os.Stdout, logFormat, cfg.LogLevel)
	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	addr := net.JoinHostPort(cfg.HTTP.Host, strconv.FormatUint(uint64(cfg.HTTP.Port), 10))
	srv := http.Server{
		Addr:              addr,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
		Handler:           httpx.NewRouter(log),
	}

	srvError := make(chan error, 1)
	go func() {
		log.Info("server running", "addr", srv.Addr)
		srvError <- srv.ListenAndServe()
	}()

	select {
	case err := <-srvError:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
	}

	log.Info("shutdown started")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	log.Info("shutdown finished")

	return nil
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
