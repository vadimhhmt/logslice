// Package ratelimit implements a time-window-based rate limiter for log
// entries processed by logslice.
//
// A Limiter counts entries within a rolling window defined by a fixed
// duration. Once the maximum number of entries for the current window has
// been reached, subsequent entries are dropped until the window resets.
//
// The window boundary is derived from the entry's own timestamp, making
// the limiter deterministic regardless of wall-clock time. Entries that
// carry no timestamp (zero time.Time) bypass rate limiting entirely.
//
// Usage:
//
//	l, err := ratelimit.New(500, time.Minute)
//	if err != nil { ... }
//	if l.Allow(entry) {
//	    // process entry
//	}
package ratelimit
