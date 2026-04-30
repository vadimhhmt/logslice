// Package retry provides a configurable retry stage for the logslice pipeline.
//
// When a downstream processing function returns an error (for example, a
// transient write failure), the Retryer will re-invoke the function up to
// MaxAttempts times, sleeping for Delay between each attempt.
//
// Basic usage:
//
//	r := retry.New(retry.Policy{
//		MaxAttempts: 3,
//		Delay:       50 * time.Millisecond,
//	})
//
//	err := r.Run(entry, func(e parser.Entry) error {
//		return writeToSink(e)
//	})
//
// RunAll is a convenience wrapper that processes a slice of entries and
// returns aggregate Counts of successes, retried successes, and failures.
package retry
