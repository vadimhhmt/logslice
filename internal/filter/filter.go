// Package filter provides time-range and field-pattern filtering for log entries.
package filter

import (
	"fmt"
	"regexp"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Filter holds the criteria used to include or exclude log entries.
type Filter struct {
	From     *time.Time
	To       *time.Time
	Patterns []*regexp.Regexp
}

// Options configures a Filter.
type Options struct {
	// From is the inclusive start of the time window. Nil means no lower bound.
	From *time.Time
	// To is the inclusive end of the time window. Nil means no upper bound.
	To *time.Time
	// Patterns is a list of regular expression strings that must ALL match at
	// least one field value in an entry for it to be included.
	Patterns []string
}

// New constructs a Filter from Options, compiling all pattern strings.
// Returns an error if any pattern fails to compile.
func New(opts Options) (*Filter, error) {
	f := &Filter{
		From: opts.From,
		To:   opts.To,
	}
	for _, raw := range opts.Patterns {
		re, err := regexp.Compile(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %q: %w", raw, err)
		}
		f.Patterns = append(f.Patterns, re)
	}
	return f, nil
}

// Match reports whether entry satisfies all filter criteria.
//
// Time-range check: if From or To is set the entry's timestamp must fall
// within [From, To] (both bounds inclusive). Entries without a parseable
// timestamp are excluded when a time bound is active.
//
// Pattern check: every compiled pattern must match at least one string value
// found anywhere in the entry's Fields map.
func (f *Filter) Match(entry parser.Entry) bool {
	if f.From != nil || f.To != nil {
		if entry.Timestamp.IsZero() {
			return false
		}
		if !InRange(entry.Timestamp, f.From, f.To) {
			return false
		}
	}

	for _, re := range f.Patterns {
		if !matchesAnyField(re, entry.Fields) {
			return false
		}
	}
	return true
}

// InRange reports whether t falls within the inclusive interval [from, to].
// A nil from or to pointer means that bound is unbounded.
func InRange(t time.Time, from, to *time.Time) bool {
	if from != nil && t.Before(*from) {
		return false
	}
	if to != nil && t.After(*to) {
		return false
	}
	return true
}

// matchesAnyField returns true when re matches the string representation of
// at least one value in fields.
func matchesAnyField(re *regexp.Regexp, fields map[string]interface{}) bool {
	for _, v := range fields {
		if re.MatchString(fmt.Sprintf("%v", v)) {
			return true
		}
	}
	return false
}
