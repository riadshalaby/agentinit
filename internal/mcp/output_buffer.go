package mcp

import "sync"

// outputBuffer is a goroutine-safe append-only byte buffer that implements
// io.Writer. Adapters stream subprocess output into it; the manager reads
// from it via StringFrom.
type outputBuffer struct {
	mu   sync.Mutex
	data []byte
}

func (b *outputBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data = append(b.data, p...)
	return len(p), nil
}

// StringFrom returns the buffered output starting at byte offset off,
// along with the current total byte count.
func (b *outputBuffer) StringFrom(off int) (chunk string, total int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if off < 0 {
		off = 0
	}
	total = len(b.data)
	if off >= total {
		return "", total
	}
	return string(b.data[off:]), total
}
