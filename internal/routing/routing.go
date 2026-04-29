// Package routing provides field-based log entry routing,
// directing entries to named output buckets based on field value matches.
package routing

import (
	"fmt"
	"regexp"

	"logslice/internal/parser"
)

// Rule maps a compiled pattern to a bucket name.
type Rule struct {
	Field   string
	Pattern *regexp.Regexp
	Bucket  string
}

// Router dispatches log entries to named buckets.
type Router struct {
	rules   []Rule
	default_ string
}

// New creates a Router with the given rules and a fallback bucket name.
// An empty defaultBucket causes unmatched entries to be dropped.
func New(rules []Rule, defaultBucket string) *Router {
	return &Router{rules: rules, default_: defaultBucket}
}

// AddRule appends a routing rule. pattern is a regular expression.
func (r *Router) AddRule(field, pattern, bucket string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("routing: invalid pattern %q: %w", pattern, err)
	}
	r.rules = append(r.rules, Rule{Field: field, Pattern: re, Bucket: bucket})
	return nil
}

// Route returns the bucket name for the given entry.
// Rules are evaluated in order; the first match wins.
// If no rule matches, the default bucket is returned.
// An empty string means the entry should be dropped.
func (r *Router) Route(entry parser.Entry) string {
	for _, rule := range r.rules {
		v, ok := entry.Fields[rule.Field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		if rule.Pattern.MatchString(s) {
			return rule.Bucket
		}
	}
	return r.default_
}

// Dispatch reads entries from in, routes each one, and sends it to the
// matching channel in buckets. Unknown bucket names are silently dropped.
// Dispatch returns when in is closed; it does not close any bucket channel.
func (r *Router) Dispatch(in <-chan parser.Entry, buckets map[string]chan<- parser.Entry) {
	for entry := range in {
		bucket := r.Route(entry)
		if ch, ok := buckets[bucket]; ok {
			ch <- entry
		}
	}
}
