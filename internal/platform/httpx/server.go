package httpx

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

const MaxBodyBytes = 1 << 20 // 1MB

func NewServer(h http.Handler, idle time.Duration) *Server {
	return &Server{
		// nolint:exhaustruct_v5
		srv: &http.Server{
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       idle,
			Handler:           h,
		},
	}
}

type Server struct {
	srv *http.Server
}

func (s *Server) Start(ln net.Listener) chan error {
	servErr := make(chan error, 1)
	go func() {
		if err := s.srv.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
			servErr <- err
		}
	}()

	return servErr
}

func (s *Server) Shutdown(timeout time.Duration) error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	return nil
}
