// Package mask provides field-value masking for log entries.
// It replaces the value of nominated fields with a fixed placeholder,
// preserving the field key so downstream consumers can still detect its
// presence without seeing the raw data.
package mask

import "fmt"

const defaultPlaceholder = "***"

// Masker replaces the values of configured fields with a placeholder string.
type Masker struct {
	fields      map[string]struct{}
	placeholder string
}

// Option configures a Masker.
type Option func(*Masker)

// WithPlaceholder overrides the default placeholder string.
func WithPlaceholder(p string) Option {
	return func(m *Masker) {
		if p != "" {
			m.placeholder = p
		}
	}
}

// New returns a Masker that will replace the values of the given fields.
// Field names are matched case-sensitively.
func New(fields []string, opts ...Option) (*Masker, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("mask: at least one field name is required")
	}
	m := &Masker{
		fields:      make(map[string]struct{}, len(fields)),
		placeholder: defaultPlaceholder,
	}
	for _, opt := range opts {
		opt(m)
	}
	for _, f := range fields {
		if f == "" {
			return nil, fmt.Errorf("mask: field name must not be empty")
		}
		m.fields[f] = struct{}{}
	}
	return m, nil
}

// Apply returns a copy of entry with targeted field values replaced by the
// placeholder. Fields not present in entry are silently ignored.
func (m *Masker) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		if _, ok := m.fields[k]; ok {
			out[k] = m.placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

// Fields returns the set of field names that will be masked.
func (m *Masker) Fields() []string {
	out := make([]string, 0, len(m.fields))
	for f := range m.fields {
		out = append(out, f)
	}
	return out
}

// IsMasked reports whether the given field name is configured to be masked.
func (m *Masker) IsMasked(field string) bool {
	_, ok := m.fields[field]
	return ok
}
