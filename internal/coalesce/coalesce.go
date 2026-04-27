// Package coalesce provides a processor that merges a prioritised list of
// source fields into a single destination field, using the first non-empty
// value found.
package coalesce

import "maps"

// Entry is the log entry type shared across logslice packages.
type Entry map[string]any

// Coalescer writes the first non-empty value from Sources into Dest.
type Coalescer struct {
	dest    string
	sources []string
}

// New creates a Coalescer that reads from sources (in order) and writes the
// first non-empty value to dest. At least one source must be provided.
func New(dest string, sources []string) (*Coalescer, error) {
	if dest == "" {
		return nil, errorf("dest field name must not be empty")
	}
	if len(sources) == 0 {
		return nil, errorf("at least one source field is required")
	}
	filtered := make([]string, 0, len(sources))
	seen := make(map[string]struct{}, len(sources))
	for _, s := range sources {
		if s == "" {
			continue
		}
		if _, dup := seen[s]; dup {
			continue
		}
		seen[s] = struct{}{}
		filtered = append(filtered, s)
	}
	if len(filtered) == 0 {
		return nil, errorf("all provided source field names were empty or duplicate")
	}
	return &Coalescer{dest: dest, sources: filtered}, nil
}

// Process evaluates the entry and returns a (possibly new) entry with the
// coalesced field set. If no source field contains a non-empty value the entry
// is returned unchanged.
func (c *Coalescer) Process(e Entry) Entry {
	for _, src := range c.sources {
		v, ok := e[src]
		if !ok {
			continue
		}
		if s, isStr := v.(string); isStr && s == "" {
			continue
		}
		if v == nil {
			continue
		}
		out := cloneEntry(e)
		out[c.dest] = v
		return out
	}
	return e
}

func cloneEntry(e Entry) Entry {
	return maps.Clone(e)
}

func errorf(msg string) error {
	return &coalesceError{msg: msg}
}

type coalesceError struct{ msg string }

func (e *coalesceError) Error() string { return "coalesce: " + e.msg }
