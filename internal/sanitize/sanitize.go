// Package sanitize provides helpers for cleaning and normalising raw log
// lines before they are parsed or forwarded to the output formatter.
package sanitize

import (
	"strings"
	"unicode"
)

// MaxLineBytes is the maximum number of bytes a single log line may occupy
// after trimming. Lines that exceed this limit are truncated and annotated
// with a "…" suffix so that downstream consumers are aware of the loss.
const MaxLineBytes = 16 * 1024 // 16 KiB

// Line trims leading/trailing whitespace, strips non-printable control
// characters (except for regular ASCII space), and truncates the result to
// MaxLineBytes. It returns an empty string for blank or whitespace-only
// input so that callers can skip such lines cheaply.
func Line(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}

	// Remove non-printable runes that are not ordinary whitespace.
	s = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\t' {
			return -1
		}
		return r
	}, s)

	if len(s) > MaxLineBytes {
		s = s[:MaxLineBytes] + "…"
	}

	return s
}

// FieldName normalises a JSON field name by lower-casing it and trimming
// surrounding whitespace. This makes field lookups case-insensitive without
// requiring callers to remember the exact casing used in log files.
func FieldName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
