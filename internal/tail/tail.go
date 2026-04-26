// Package tail provides a log entry emitter that follows a stream
// and emits only the last N entries — similar to `tail -n`.
package tail

import "github.com/logslice/logslice/internal/parser"

// Tailer holds a fixed-size ring buffer and emits the last N entries.
type Tailer struct {
	n      int
	buf    []parser.Entry
	head   int
	count  int
}

// New returns a Tailer that retains the last n entries.
// If n <= 0 it defaults to 10.
func New(n int) *Tailer {
	if n <= 0 {
		n = 10
	}
	return &Tailer{
		n:   n,
		buf: make([]parser.Entry, n),
	}
}

// Push adds an entry to the ring buffer, overwriting the oldest when full.
func (t *Tailer) Push(e parser.Entry) {
	t.buf[t.head] = e
	t.head = (t.head + 1) % t.n
	if t.count < t.n {
		t.count++
	}
}

// Entries returns the retained entries in chronological order (oldest first).
func (t *Tailer) Entries() []parser.Entry {
	out := make([]parser.Entry, t.count)
	start := 0
	if t.count == t.n {
		start = t.head
	}
	for i := 0; i < t.count; i++ {
		out[i] = t.buf[(start+i)%t.n]
	}
	return out
}

// Len returns the number of entries currently held.
func (t *Tailer) Len() int { return t.count }
