// Package enrich adds derived fields to log entries based on existing field values.
package enrich

import (
	"fmt"
	"strings"

	"logslice/internal/parser"
)

// Rule describes a single enrichment: read From, apply Fn, write result to To.
type Rule struct {
	From string
	To   string
	Fn   func(value string) (string, bool)
}

// Enricher applies a set of rules to each log entry.
type Enricher struct {
	rules []Rule
}

// New returns an Enricher that will apply the given rules in order.
func New(rules []Rule) *Enricher {
	return &Enricher{rules: rules}
}

// Apply returns a new entry with enriched fields added. The original is not
// modified. If a target field already exists it is overwritten.
func (e *Enricher) Apply(entry parser.Entry) parser.Entry {
	out := cloneEntry(entry)
	for _, r := range e.rules {
		raw, ok := out.Fields[r.From]
		if !ok {
			continue
		}
		val := fmt.Sprintf("%v", raw)
		if result, matched := r.Fn(val); matched {
			out.Fields[r.To] = result
		}
	}
	return out
}

// cloneEntry performs a shallow copy of the entry's field map.
func cloneEntry(e parser.Entry) parser.Entry {
	fields := make(map[string]any, len(e.Fields))
	for k, v := range e.Fields {
		fields[k] = v
	}
	return parser.Entry{Timestamp: e.Timestamp, Fields: fields, Raw: e.Raw}
}

// UpperCase is a convenience Fn that upper-cases the source value.
func UpperCase(v string) (string, bool) { return strings.ToUpper(v), true }

// LowerCase is a convenience Fn that lower-cases the source value.
func LowerCase(v string) (string, bool) { return strings.ToLower(v), true }
