package mwtest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/logger"
)

type Buffer struct {
	mu  sync.RWMutex
	buf bytes.Buffer
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buf.Write(p)
}

func (b *Buffer) String() string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.buf.String()
}

func NewTestLogger(t testing.TB) (*slog.Logger, *Buffer) {
	t.Helper()
	var b Buffer
	log := logger.New(&b, logger.FormatJSON, slog.LevelDebug)

	return log, &b
}

func LogLines(t testing.TB, buf *Buffer) []map[string]any {
	t.Helper()
	result := make([]map[string]any, 0)

	s := bufio.NewScanner(strings.NewReader(buf.String()))
	for s.Scan() {
		var rec map[string]any
		err := json.Unmarshal([]byte(s.Text()), &rec)
		require.NoError(t, err)
		result = append(result, rec)
	}

	return result
}
