// Package offset provides a processor that skips the first N log entries,
// useful for resuming processing from a known position in a log stream.
package offset

import "github.com/logslice/logslice/internal/parser"

// Skipper discards the first N entries passed through it, then forwards
// all subsequent entries unchanged.
type Skipper struct {
	skip    int
	seen    int
	dropped int
}

// New returns a Skipper that will discard the first n entries.
// If n is zero or negative, no entries are skipped.
func New(n int) *Skipper {
	if n < 0 {
		n = 0
	}
	return &Skipper{skip: n}
}

// Process evaluates the entry against the skip counter.
// It returns (entry, true) once the skip threshold has been reached,
// and (zero, false) while entries are still being discarded.
func (s *Skipper) Process(entry parser.Entry) (parser.Entry, bool) {
	if s.seen < s.skip {
		s.seen++
		s.dropped++
		return parser.Entry{}, false
	}
	s.seen++
	return entry, true
}

// Dropped returns the number of entries that have been skipped so far.
func (s *Skipper) Dropped() int {
	return s.dropped
}

// Reset restores the Skipper to its initial state, allowing the skip
// counter to be replayed from zero.
func (s *Skipper) Reset() {
	s.seen = 0
	s.dropped = 0
}
