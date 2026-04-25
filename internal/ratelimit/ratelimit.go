// Package ratelimit provides a token-bucket style log entry rate limiter
// that drops entries exceeding a maximum count per time window.
package ratelimit

import (
	"fmt"
	"time"

	"logslice/internal/parser"
)

// Limiter drops log entries that exceed a maximum count within a rolling
// time window. Entries are counted per window bucket; once the bucket is
// full, subsequent entries are dropped until the window resets.
type Limiter struct {
	maxPerWindow int
	window       time.Duration
	bucketStart  time.Time
	bucketCount  int
	Dropped      int
}

// New returns a Limiter that allows at most maxPerWindow entries per window
// duration. Returns an error if either argument is non-positive.
func New(maxPerWindow int, window time.Duration) (*Limiter, error) {
	if maxPerWindow <= 0 {
		return nil, fmt.Errorf("ratelimit: maxPerWindow must be > 0, got %d", maxPerWindow)
	}
	if window <= 0 {
		return nil, fmt.Errorf("ratelimit: window must be > 0, got %s", window)
	}
	return &Limiter{
		maxPerWindow: maxPerWindow,
		window:       window,
	}, nil
}

// Allow returns true if the entry should be passed through, false if it
// should be dropped. The entry's timestamp is used to determine the current
// window; entries with a zero timestamp are always allowed.
func (l *Limiter) Allow(entry parser.Entry) bool {
	ts := entry.Timestamp
	if ts.IsZero() {
		return true
	}

	// Reset bucket when the window has elapsed.
	if l.bucketStart.IsZero() || ts.Sub(l.bucketStart) >= l.window {
		l.bucketStart = ts.Truncate(l.window)
		l.bucketCount = 0
	}

	l.bucketCount++
	if l.bucketCount > l.maxPerWindow {
		l.Dropped++
		return false
	}
	return true
}
