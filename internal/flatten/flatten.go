// Package flatten provides a processor that flattens nested JSON objects
// in log entries into dot-separated top-level keys.
package flatten

import (
	"fmt"
	"strings"

	"logslice/internal/parser"
)

// Flattener collapses nested map fields into dot-notation keys.
type Flattener struct {
	separator string
	maxDepth  int
}

// New returns a Flattener using the given separator (e.g. ".") and maximum
// nesting depth. A maxDepth of 0 means unlimited.
func New(separator string, maxDepth int) *Flattener {
	if separator == "" {
		separator = "."
	}
	return &Flattener{separator: separator, maxDepth: maxDepth}
}

// Process returns a copy of entry with all nested map values promoted to
// top-level dot-notation keys. The original entry is never mutated.
func (f *Flattener) Process(entry parser.Entry) parser.Entry {
	out := make(parser.Entry, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for k, v := range entry {
		if nested, ok := v.(map[string]interface{}); ok {
			delete(out, k)
			f.flatten(out, k, nested, 1)
		}
	}
	return out
}

func (f *Flattener) flatten(dst parser.Entry, prefix string, src map[string]interface{}, depth int) {
	for k, v := range src {
		key := strings.Join([]string{prefix, k}, f.separator)
		if nested, ok := v.(map[string]interface{}); ok && (f.maxDepth == 0 || depth < f.maxDepth) {
			f.flatten(dst, key, nested, depth+1)
		} else {
			if nested, ok := v.(map[string]interface{}); ok {
				// max depth reached — store as formatted string
				dst[key] = fmt.Sprintf("%v", nested)
			} else {
				dst[key] = v
			}
		}
	}
}
