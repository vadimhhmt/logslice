// Package stats provides aggregation and summary statistics
// for processed log entries.
package stats

import (
	"fmt"
	"io"
	"time"

	"github.com/user/logslice/internal/parser"
)

// Collector accumulates statistics about log entries processed
// during a logslice run.
type Collector struct {
	Total     int
	Matched   int
	Skipped   int
	Earliest  time.Time
	Latest    time.Time
	Levels    map[string]int
	hasTime   bool
}

// New returns an initialised Collector.
func New() *Collector {
	return &Collector{
		Levels: make(map[string]int),
	}
}

// Record updates the collector with a parsed log entry.
// matched indicates whether the entry passed all filters.
func (c *Collector) Record(entry parser.Entry, matched bool) {
	c.Total++
	if !matched {
		c.Skipped++
		return
	}
	c.Matched++

	if !entry.Timestamp.IsZero() {
		if !c.hasTime || entry.Timestamp.Before(c.Earliest) {
			c.Earliest = entry.Timestamp
		}
		if !c.hasTime || entry.Timestamp.After(c.Latest) {
			c.Latest = entry.Timestamp
		}
		c.hasTime = true
	}

	if lvl, ok := entry.Fields["level"]; ok {
		if s, ok := lvl.(string); ok {
			c.Levels[s]++
		}
	}
}

// Print writes a human-readable summary to w.
func (c *Collector) Print(w io.Writer) {
	fmt.Fprintf(w, "--- summary ---\n")
	fmt.Fprintf(w, "total lines : %d\n", c.Total)
	fmt.Fprintf(w, "matched     : %d\n", c.Matched)
	fmt.Fprintf(w, "skipped     : %d\n", c.Skipped)

	if c.hasTime {
		fmt.Fprintf(w, "earliest    : %s\n", c.Earliest.Format(time.RFC3339))
		fmt.Fprintf(w, "latest      : %s\n", c.Latest.Format(time.RFC3339))
	}

	if len(c.Levels) > 0 {
		fmt.Fprintf(w, "levels      :\n")
		for lvl, n := range c.Levels {
			fmt.Fprintf(w, "  %-10s %d\n", lvl, n)
		}
	}
}
