// Package retry provides a pipeline stage that re-queues log entries
// that fail downstream processing, up to a configurable maximum attempt count.
package retry

import (
	"time"

	"logslice/internal/parser"
)

// Policy describes how retries are scheduled.
type Policy struct {
	MaxAttempts int
	Delay       time.Duration
}

// Retryer wraps a processing function and retries it on failure.
type Retryer struct {
	policy  Policy
	sleepFn func(time.Duration) // injectable for tests
}

// New returns a Retryer with the given policy.
func New(p Policy) *Retryer {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	return &Retryer{policy: p, sleepFn: time.Sleep}
}

// ProcessFunc is the signature of a function applied to a log entry.
// It returns an error if the entry could not be processed.
type ProcessFunc func(entry parser.Entry) error

// Run passes entry through fn, retrying up to MaxAttempts times.
// It returns the last error if all attempts are exhausted, or nil on success.
func (r *Retryer) Run(entry parser.Entry, fn ProcessFunc) error {
	var err error
	for attempt := 0; attempt < r.policy.MaxAttempts; attempt++ {
		if err = fn(entry); err == nil {
			return nil
		}
		if attempt < r.policy.MaxAttempts-1 && r.policy.Delay > 0 {
			r.sleepFn(r.policy.Delay)
		}
	}
	return err
}

// Counts tracks retry statistics.
type Counts struct {
	Succeeded int
	Failed    int
	Retried   int
}

// RunAll applies fn to every entry in entries, collecting counts.
func (r *Retryer) RunAll(entries []parser.Entry, fn ProcessFunc) Counts {
	var c Counts
	for _, e := range entries {
		var attempts int
		var lastErr error
		for attempts = 0; attempts < r.policy.MaxAttempts; attempts++ {
			if lastErr = fn(e); lastErr == nil {
				break
			}
			if attempts < r.policy.MaxAttempts-1 && r.policy.Delay > 0 {
				r.sleepFn(r.policy.Delay)
			}
		}
		if lastErr != nil {
			c.Failed++
		} else {
			if attempts > 0 {
				c.Retried++
			}
			c.Succeeded++
		}
	}
	return c
}
