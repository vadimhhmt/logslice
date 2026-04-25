// Package truncate provides field-value length limiting for log entries.
//
// A Truncator can be configured with a maximum rune count and an optional
// suffix (e.g. "...") appended to values that exceed the limit. Truncation
// can be applied to all string fields or scoped to a named subset.
//
// Example:
//
//	tr := truncate.New(120, "...", []string{"message", "error"})
//	processed := tr.Apply(entry)
//
// Non-string field values are passed through unchanged. The original entry
// map is never modified; Apply always returns a new map.
package truncate
