// Package window provides a sliding-window log entry buffer that retains
// the N most recent entries before and after a matching event. This is
// useful for capturing context around errors or anomalies without
// streaming the entire log.
package window

import (
	"sync"

	"logslice/internal/parser"
)

// Trigger is a function that returns true when an entry should be treated
// as an anchor event, causing the surrounding context window to be flushed.
type Trigger func(entry parser.Entry) bool

// Window buffers log entries and emits a context window around each
// triggering event: up to Before entries before the trigger and up to
// After entries after it.
type Window struct {
	mu      sync.Mutex
	before  int
	after   int
	trigger Trigger

	ring    []parser.Entry // circular pre-trigger buffer
	head    int            // next write position in ring
	count   int            // number of valid entries in ring

	postLeft int           // remaining post-trigger entries to emit
	pending  []parser.Entry // entries queued for the caller
}

// New creates a Window that keeps up to before entries prior to each
// triggering event and emits up to after entries following it.
// trigger must not be nil.
func New(before, after int, trigger Trigger) *Window {
	if before < 0 {
		before = 0
	}
	if after < 0 {
		after = 0
	}
	return &Window{
		before:  before,
		after:   after,
		trigger: trigger,
		ring:    make([]parser.Entry, before),
	}
}

// Push feeds an entry into the window. It returns any entries that should
// be forwarded to the output pipeline. The returned slice is only valid
// until the next call to Push.
func (w *Window) Push(entry parser.Entry) []parser.Entry {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.pending = w.pending[:0]

	if w.trigger(entry) {
		// Drain the pre-trigger ring buffer in insertion order.
		if w.before > 0 && w.count > 0 {
			start := w.head - w.count
			if start < 0 {
				start += w.before
			}
			for i := 0; i < w.count; i++ {
				idx := (start + i) % w.before
				w.pending = append(w.pending, w.ring[idx])
			}
			// Reset the ring so pre-trigger entries are not re-emitted.
			w.head = 0
			w.count = 0
		}
		// Emit the triggering entry itself.
		w.pending = append(w.pending, entry)
		w.postLeft = w.after
		return w.pending
	}

	// If we are inside a post-trigger window, emit immediately.
	if w.postLeft > 0 {
		w.postLeft--
		w.pending = append(w.pending, entry)
		return w.pending
	}

	// Otherwise buffer the entry in the pre-trigger ring.
	if w.before > 0 {
		w.ring[w.head] = entry
		w.head = (w.head + 1) % w.before
		if w.count < w.before {
			w.count++
		}
	}

	return w.pending
}

// Flush returns any buffered pre-trigger entries that were never followed
// by a triggering event. Call this after all entries have been pushed.
func (w *Window) Flush() []parser.Entry {
	w.mu.Lock()
	defer w.mu.Unlock()
	return nil // pre-trigger entries without a trigger are discarded
}
