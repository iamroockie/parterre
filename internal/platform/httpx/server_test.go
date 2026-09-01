package httpx_test

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/httpx"
)

func TestGracefulShutdown(t *testing.T) {
	shutdownPeriod := 100 * time.Millisecond
	handlerStarted := make(chan struct{}, 1)
	status := make(chan int, 1)
	fastHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		handlerStarted <- struct{}{}
		time.Sleep(shutdownPeriod / 5)
		w.WriteHeader(http.StatusOK)
	})
	slowHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		handlerStarted <- struct{}{}
		time.Sleep(shutdownPeriod * 5)
		w.WriteHeader(http.StatusOK)
	})
	timeoutHandler := http.TimeoutHandler(slowHandler, shutdownPeriod, "timeout")
	http.HandleFunc("/fast", fastHandler)
	http.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		timeoutHandler.ServeHTTP(w, r)
	})

	tests := map[string]struct {
		path            string
		wantShutdownErr error
		wantStatus      int
	}{
		"fast": {"/fast", nil, http.StatusOK},
		"slow": {"/slow", context.DeadlineExceeded, http.StatusServiceUnavailable},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ln := newListener(t)
			srv := httpx.NewServer(nil, time.Minute)
			srv.Start(ln)

			go func() {
				// nolint:noctx
				resp, err := http.Get("http://" + ln.Addr().String() + "/" + tt.path)
				require.NoError(t, err)
				defer func() { _ = resp.Body.Close() }()
				status <- resp.StatusCode
			}()

			<-handlerStarted

			require.ErrorIs(t, srv.Shutdown(shutdownPeriod), tt.wantShutdownErr)
			require.Equal(t, tt.wantStatus, <-status)
		})
	}
}

func newListener(t testing.TB) net.Listener {
	t.Helper()

	var lc net.ListenConfig
	ln, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })

	return ln
}
