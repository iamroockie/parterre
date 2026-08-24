package loggertest

import (
	"bufio"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/logger"
)

func NewLogger(t testing.TB, level slog.Level) (*slog.Logger, *Buffer) {
	t.Helper()
	var b Buffer
	log := logger.New(&b, level)

	return log, &b
}

func Logs(t testing.TB, buf *Buffer) []map[string]any {
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
