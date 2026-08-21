package api_test

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/app/api"
	"github.com/iamroockie/parterre/internal/platform/config/configtest"
	"github.com/iamroockie/parterre/internal/platform/logger"
)

func newListener(t testing.TB) net.Listener {
	t.Helper()

	var lc net.ListenConfig
	ln, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })

	return ln
}

func getStatus(t *testing.T, client *http.Client, url string) int {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	return resp.StatusCode
}

func TestRun(t *testing.T) {
	ln := newListener(t)
	cfg := configtest.Config(t)
	log := logger.New(os.Stdout, logger.FormatText, slog.LevelDebug)
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	runErr := make(chan error, 1)
	go func() {
		runErr <- api.Run(ctx, cfg, log, ln)
	}()

	baseURL := "http://" + ln.Addr().String()

	for _, route := range []string{"healthz", "readyz"} {
		t.Run(route, func(t *testing.T) {
			require.Equal(t, http.StatusOK, getStatus(t, client, baseURL+"/"+route))
		})
	}

	cancel()

	select {
	case err := <-runErr:
		require.NoError(t, err)
	case <-time.After(cfg.GracefulShutdownTimeout + time.Second):
		t.Fatal("Run не завершился после отмены контекста")
	}
}
