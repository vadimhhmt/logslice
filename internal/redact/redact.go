// Package redact provides utilities for masking sensitive field values
// in structured log entries before output.
package redact

import (
	"regexp"
	"strings"
)

// defaultSensitiveKeys are field names that are redacted by default.
var defaultSensitiveKeys = []string{
	"password", "passwd", "secret", "token", "api_key", "apikey",
	"authorization", "auth", "credential", "private_key",
}

const redactedValue = "[REDACTED]"

// Redactor masks sensitive values in log entry fields.
type Redactor struct {
	keys     map[string]struct{}
	patterns []*regexp.Regexp
}

// New creates a Redactor with the default sensitive key list.
// Additional keys and value patterns can be added via AddKey and AddPattern.
func New() *Redactor {
	r := &Redactor{
		keys: make(map[string]struct{}),
	}
	for _, k := range defaultSensitiveKeys {
		r.keys[strings.ToLower(k)] = struct{}{}
	}
	return r
}

// AddKey registers an additional field name whose value should be redacted.
func (r *Redactor) AddKey(key string) {
	r.keys[strings.ToLower(strings.TrimSpace(key))] = struct{}{}
}

// AddPattern registers a regexp; any field value matching it will be redacted.
func (r *Redactor) AddPattern(pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	r.patterns = append(r.patterns, re)
	return nil
}

// Apply returns a copy of fields with sensitive values masked.
func (r *Redactor) Apply(fields map[string]any) map[string]any {
	out := make(map[string]any, len(fields))
	for k, v := range fields {
		if _, sensitive := r.keys[strings.ToLower(k)]; sensitive {
			out[k] = redactedValue
			continue
		}
		if r.valueMatchesPattern(v) {
			out[k] = redactedValue
			continue
		}
		out[k] = v
	}
	return out
}

func (r *Redactor) valueMatchesPattern(v any) bool {
	if len(r.patterns) == 0 {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	for _, re := range r.patterns {
		if re.MatchString(s) {
			return true
		}
	}
	return false
}
