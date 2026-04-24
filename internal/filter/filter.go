package filter

import (
	"regexp"
	"time"
)

// Options holds the filtering criteria for log lines.
type Options struct {
	From    time.Time
	To      time.Time
	Pattern string
}

// Filter applies time range and field pattern filtering to parsed log entries.
type Filter struct {
	opts    Options
	pattern *regexp.Regexp
}

// New creates a new Filter from the given Options.
// Returns an error if the pattern is an invalid regular expression.
func New(opts Options) (*Filter, error) {
	f := &Filter{opts: opts}
	if opts.Pattern != "" {
		re, err := regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
		f.pattern = re
	}
	return f, nil
}

// Match reports whether the given log entry (represented as a raw line and
// its parsed timestamp) passes all active filter criteria.
func (f *Filter) Match(line string, ts time.Time) bool {
	if !f.opts.From.IsZero() && ts.Before(f.opts.From) {
		return false
	}
	if !f.opts.To.IsZero() && ts.After(f.opts.To) {
		return false
	}
	if f.pattern != nil && !f.pattern.MatchString(line) {
		return false
	}
	return true
}

// InRange reports whether ts falls within [from, to].
// A zero value for from or to means that bound is unbounded.
func InRange(ts, from, to time.Time) bool {
	if !from.IsZero() && ts.Before(from) {
		return false
	}
	if !to.IsZero() && ts.After(to) {
		return false
	}
	return true
}
