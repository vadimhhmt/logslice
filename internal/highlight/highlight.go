// Package highlight provides ANSI colour highlighting for matched
// field values in log output.
package highlight

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// Highlighter wraps matched substrings in an ANSI colour sequence.
type Highlighter struct {
	colour   string
	patterns []*regexp.Regexp
	enabled  bool
}

// New returns a Highlighter that applies colour to substrings matching
// any of the provided patterns. When enabled is false all methods
// return their input unchanged so callers need not branch.
func New(colour string, enabled bool, patterns ...*regexp.Regexp) *Highlighter {
	return &Highlighter{
		colour:   colour,
		patterns: patterns,
		enabled:  enabled,
	}
}

// Apply wraps every match found in s with the configured ANSI colour.
func (h *Highlighter) Apply(s string) string {
	if !h.enabled || len(h.patterns) == 0 {
		return s
	}
	for _, re := range h.patterns {
		s = re.ReplaceAllStringFunc(s, func(m string) string {
			return fmt.Sprintf("%s%s%s", h.colour, m, Reset)
		})
	}
	return s
}

// ApplyToFields highlights values of the named fields inside a JSON-like
// key:value string produced by the pretty formatter. Only string values
// that contain a pattern match are coloured.
func (h *Highlighter) ApplyToFields(line string, fields []string) string {
	if !h.enabled || len(h.patterns) == 0 {
		return line
	}
	for _, f := range fields {
		prefix := f + "="
		idx := strings.Index(line, prefix)
		if idx == -1 {
			continue
		}
		valStart := idx + len(prefix)
		valEnd := strings.IndexAny(line[valStart:], " \t\n")
		var val string
		if valEnd == -1 {
			val = line[valStart:]
		} else {
			val = line[valStart : valStart+valEnd]
		}
		highlighted := h.Apply(val)
		if highlighted != val {
			line = line[:valStart] + highlighted + line[valStart+len(val):]
		}
	}
	return line
}
