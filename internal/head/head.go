// Package head provides a processor that passes through only the first N
// log entries, then signals the pipeline to stop.
package head

import "github.com/logslice/logslice/internal/parser"

// Taker holds state for the head processor.
type Taker struct {
	max   int
	seen  int
}

// New returns a Taker that will pass through at most n entries.
// If n is zero or negative, all entries are passed through.
func New(n int) *Taker {
	return &Taker{max: n}
}

// Done reports whether the Taker has already collected its quota.
// When Done returns true the caller should stop feeding entries.
func (t *Taker) Done() bool {
	if t.max <= 0 {
		return false
	}
	return t.seen >= t.max
}

// Process accepts an entry and returns (entry, true) if it should be
// forwarded, or (zero, false) once the quota has been reached.
func (t *Taker) Process(e parser.Entry) (parser.Entry, bool) {
	if t.max <= 0 {
		return e, true
	}
	if t.seen >= t.max {
		return parser.Entry{}, false
	}
	t.seen++
	return e, true
}

// Reset resets the internal counter so the Taker can be reused.
func (t *Taker) Reset() {
	t.seen = 0
}
