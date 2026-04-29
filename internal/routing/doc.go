// Package routing provides field-based log entry routing for logslice.
//
// A Router holds an ordered list of Rules. Each Rule pairs a field name,
// a compiled regular expression, and a destination bucket name. When an
// entry is routed, rules are tested in declaration order and the first
// match determines the bucket. Entries that match no rule are sent to
// the configured default bucket; an empty default bucket name causes
// unmatched entries to be silently dropped.
//
// Example — split errors from everything else:
//
//	r := routing.New(nil, "general")
//	r.AddRule("level", `(?i)error`, "errors")
//
//	// channel-based fan-out
//	r.Dispatch(entryCh, map[string]chan<- parser.Entry{
//		"errors":  errorsCh,
//		"general": generalCh,
//	})
package routing
