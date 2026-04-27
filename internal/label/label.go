// Package label attaches a static or derived tag to every log entry
// that passes through the pipeline. This is useful when merging streams
// from multiple sources and you need to identify the origin of each entry.
package label

import (
	"fmt"

	"github.com/logslice/logslice/internal/parser"
)

// Labeler adds a fixed key/value pair to every entry it processes.
type Labeler struct {
	key       string
	value     string
	overwrite bool
}

// New returns a Labeler that will set entry.Fields[key] = value.
// When overwrite is false the field is only written if it is not already
// present in the entry.
func New(key, value string, overwrite bool) (*Labeler, error) {
	if key == "" {
		return nil, fmt.Errorf("label: key must not be empty")
	}
	return &Labeler{key: key, value: value, overwrite: overwrite}, nil
}

// Apply attaches the label to e, returning a shallow-cloned entry so the
// original is never mutated.
func (l *Labeler) Apply(e parser.Entry) parser.Entry {
	_, exists := e.Fields[l.key]
	if exists && !l.overwrite {
		return e
	}
	out := cloneEntry(e)
	out.Fields[l.key] = l.value
	return out
}

// Process reads entries from in, labels each one, and writes to out.
// It returns when in is closed.
func (l *Labeler) Process(in <-chan parser.Entry, out chan<- parser.Entry) {
	for e := range in {
		out <- l.Apply(e)
	}
}

func cloneEntry(e parser.Entry) parser.Entry {
	fields := make(map[string]interface{}, len(e.Fields))
	for k, v := range e.Fields {
		fields[k] = v
	}
	e.Fields = fields
	return e
}
