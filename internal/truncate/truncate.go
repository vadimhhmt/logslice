// Package truncate provides utilities for limiting field value lengths
// in log entries before output or further processing.
package truncate

import (
	"unicode/utf8"
)

// Truncator holds configuration for field truncation.
type Truncator struct {
	maxLen  int
	suffix  string
	fields  map[string]bool // if non-empty, only truncate these fields
}

// New creates a Truncator that limits string field values to maxLen runes.
// If maxLen <= 0, no truncation is applied.
func New(maxLen int, suffix string, fields []string) *Truncator {
	set := make(map[string]bool, len(fields))
	for _, f := range fields {
		if f != "" {
			set[f] = true
		}
	}
	return &Truncator{
		maxLen: maxLen,
		suffix: suffix,
		fields: set,
	}
}

// Apply returns a copy of entry with string values truncated according to
// the Truncator's configuration. Non-string values are left unchanged.
func (t *Truncator) Apply(entry map[string]interface{}) map[string]interface{} {
	if t.maxLen <= 0 {
		return entry
	}
	out := make(map[string]interface{}, len(entry))
	for k, v := range entry {
		if t.shouldTruncate(k) {
			if s, ok := v.(string); ok {
				v = t.truncateString(s)
			}
		}
		out[k] = v
	}
	return out
}

func (t *Truncator) shouldTruncate(field string) bool {
	if len(t.fields) == 0 {
		return true
	}
	return t.fields[field]
}

func (t *Truncator) truncateString(s string) string {
	if utf8.RuneCountInString(s) <= t.maxLen {
		return s
	}
	runes := []rune(s)
	return string(runes[:t.maxLen]) + t.suffix
}
