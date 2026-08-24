package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/iamroockie/parterre/internal/platform/config"
	"github.com/iamroockie/parterre/internal/platform/postgres"
)

func Run(ctx context.Context, cfg config.Config, log *slog.Logger, ln net.Listener) error {
	pool, err := postgres.NewPool(ctx, cfg.Postgres.DSN, cfg.Postgres.MaxConns)
	if err != nil {
		return fmt.Errorf("postgres connect: %w", err)
	}
	defer pool.Close()

	r := newRouter(log)
	routes(r, routesDeps{cfg, pool})

	srv := newServer(r, cfg.HTTP.IdleTimeout)
	srvError := make(chan error, 1)
	go func() {
		log.Info("server running", "addr", ln.Addr().String())
		srvError <- srv.Serve(ln)
	}()

	select {
	case err := <-srvError:
		return fmt.Errorf("server: %w", err)
	case <-ctx.Done():
	}

	if err := Shutdown(srv, srvError, log, cfg.ShutdownTimeout); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	return nil
}

func Shutdown(srv *http.Server, ch chan error, log *slog.Logger, timeout time.Duration) error {
	log.Info("shutdown started")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	select {
	case err := <-ch:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("close server: %w", err)
		}
	case <-shutdownCtx.Done():
	}

	log.Info("shutdown finished")

	return nil
}

func newServer(r chi.Router, idle time.Duration) *http.Server {
	return &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       idle,
		Handler:           r,
	}
}
