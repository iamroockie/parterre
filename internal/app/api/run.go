package api

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/iamroockie/parterre/internal/platform/config"
	"github.com/iamroockie/parterre/internal/platform/httpx"
)

func Run(ctx context.Context, cfg config.Config, log *slog.Logger, ln net.Listener) error {
	srv := &http.Server{
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
		Handler:           httpx.NewRouter(log),
	}

	srvError := make(chan error, 1)
	go func() {
		log.Info("server running", "addr", ln.Addr())
		srvError <- srv.Serve(ln)
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
