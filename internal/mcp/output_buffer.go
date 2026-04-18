package mcp

import "sync"

// outputBuffer is a goroutine-safe append-only byte buffer that implements
// io.Writer. Adapters stream subprocess output into it; the manager reads
// from it via StringFrom or StringFromLimit.
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
	return b.StringFromLimit(off, 0)
}

// StringFromLimit returns the buffered output starting at byte offset off,
// capped to at most limit bytes when limit is positive, along with the
// current total byte count.
func (b *outputBuffer) StringFromLimit(off, limit int) (chunk string, total int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if off < 0 {
		off = 0
	}
	total = len(b.data)
	if off >= total {
		return "", total
	}
	end := total
	if limit > 0 && off+limit < end {
		end = off + limit
	}
	return string(b.data[off:end]), total
}

// Tail returns the last n bytes of the buffered output.
func (b *outputBuffer) Tail(n int) string {
	b.mu.Lock()
	defer b.mu.Unlock()
	if n <= 0 || len(b.data) == 0 {
		return ""
	}
	start := 0
	if n < len(b.data) {
		start = len(b.data) - n
	}
	return string(b.data[start:])
}
