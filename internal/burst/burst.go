// Package burst provides a burst-detection stage for the log pipeline.
// It identifies time windows where log volume exceeds a configurable threshold
// and can either flag or drop entries that fall outside burst windows.
package burst

import (
	"fmt"
	"time"

	"logslice/internal/parser"
)

// Mode controls what the Detector does when a burst is detected.
type Mode int

const (
	// ModeFlag annotates entries inside a burst window with a "_burst" field.
	ModeFlag Mode = iota
	// ModeDrop discards entries that are NOT inside a burst window.
	ModeDrop
)

// Detector tracks log entry rates over a sliding window and either flags or
// drops entries depending on whether a burst threshold has been exceeded.
type Detector struct {
	window    time.Duration
	threshold int
	mode      Mode

	// ring buffer of timestamps within the current window
	bucket []time.Time
}

// New creates a Detector that considers a burst to have occurred when more than
// threshold entries arrive within window. mode controls the action taken.
func New(window time.Duration, threshold int, mode Mode) (*Detector, error) {
	if window <= 0 {
		return nil, fmt.Errorf("burst: window must be positive, got %s", window)
	}
	if threshold <= 0 {
		return nil, fmt.Errorf("burst: threshold must be positive, got %d", threshold)
	}
	return &Detector{
		window:    window,
		threshold: threshold,
		mode:      mode,
		bucket:    make([]time.Time, 0, threshold+1),
	}, nil
}

// Process evaluates entry against the current burst state.
// It returns the (possibly annotated) entry and whether it should be kept.
// A nil entry is returned when the entry should be dropped.
func (d *Detector) Process(entry parser.Entry) (*parser.Entry, bool) {
	ts := entry.Timestamp
	if ts.IsZero() {
		// Entries without a timestamp are always passed through unchanged.
		return &entry, true
	}

	d.evict(ts)
	d.bucket = append(d.bucket, ts)

	inBurst := len(d.bucket) > d.threshold

	switch d.mode {
	case ModeFlag:
		if inBurst {
			cloned := cloneEntry(entry)
			cloned.Fields["_burst"] = true
			return &cloned, true
		}
		return &entry, true

	case ModeDrop:
		if !inBurst {
			return nil, false
		}
		return &entry, true
	}

	return &entry, true
}

// evict removes timestamps from the bucket that fall outside the window
// relative to the given reference time.
func (d *Detector) evict(ref time.Time) {
	cutoff := ref.Add(-d.window)
	keep := 0
	for _, t := range d.bucket {
		if t.After(cutoff) {
			d.bucket[keep] = t
			keep++
		}
	}
	d.bucket = d.bucket[:keep]
}

// WindowCount returns the number of entries currently tracked within the
// active window. Useful for testing and diagnostics.
func (d *Detector) WindowCount() int {
	return len(d.bucket)
}

func cloneEntry(e parser.Entry) parser.Entry {
	fields := make(map[string]any, len(e.Fields))
	for k, v := range e.Fields {
		fields[k] = v
	}
	return parser.Entry{
		Timestamp: e.Timestamp,
		Raw:       e.Raw,
		Fields:    fields,
	}
}
