package loggertest

import (
	"bytes"
	"sync"
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
