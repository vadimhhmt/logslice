// Package transform provides field renaming and value mapping
// for log entries before output.
package transform

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// Rule describes a single field transformation.
type Rule struct {
	FromField string
	ToField   string
	ValueMap  map[string]string // optional: remap specific values
}

// Transformer applies a set of Rules to log entries.
type Transformer struct {
	rules []Rule
}

// New creates a Transformer from the given rules.
func New(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Apply returns a new LogEntry with all rules applied.
// Original entry is not modified.
func (t *Transformer) Apply(entry parser.LogEntry) parser.LogEntry {
	out := parser.LogEntry{
		Timestamp: entry.Timestamp,
		Raw:       entry.Raw,
		Fields:    make(map[string]interface{}, len(entry.Fields)),
	}
	for k, v := range entry.Fields {
		out.Fields[k] = v
	}
	for _, r := range t.rules {
		applyRule(out.Fields, r)
	}
	return out
}

func applyRule(fields map[string]interface{}, r Rule) {
	v, ok := fields[r.FromField]
	if !ok {
		return
	}
	if r.ToField != "" && r.ToField != r.FromField {
		fields[r.ToField] = v
		delete(fields, r.FromField)
		v = fields[r.ToField]
	}
	if len(r.ValueMap) > 0 {
		key := fmt.Sprintf("%v", v)
		if mapped, found := r.ValueMap[strings.ToLower(key)]; found {
			target := r.ToField
			if target == "" {
				target = r.FromField
			}
			fields[target] = mapped
		}
	}
}
