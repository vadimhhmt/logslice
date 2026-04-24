// Package stats collects and reports aggregate statistics for a
// logslice processing run.
//
// Usage:
//
//	c := stats.New()
//	for _, entry := range entries {
//		matched := filter.Matches(entry)
//		c.Record(entry, matched)
//	}
//	c.Print(os.Stderr)
//
// The Collector tracks total/matched/skipped line counts, the
// earliest and latest timestamps seen in matched entries, and a
// breakdown of log levels found in the "level" field.
package stats
