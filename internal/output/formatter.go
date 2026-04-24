package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Format defines the output format for log lines.
type Format string

const (
	FormatJSON    Format = "json"
	FormatPretty  Format = "pretty"
	FormatRaw     Format = "raw"
)

// Formatter writes parsed log entries to an output stream.
type Formatter struct {
	w      io.Writer
	format Format
	fields []string
}

// New creates a new Formatter with the given writer, format, and optional field selection.
func New(w io.Writer, format Format, fields []string) *Formatter {
	return &Formatter{
		w:      w,
		format: format,
		fields: fields,
	}
}

// Write formats and writes a single parsed log entry.
func (f *Formatter) Write(entry map[string]interface{}) error {
	if len(f.fields) > 0 {
		entry = selectFields(entry, f.fields)
	}

	switch f.format {
	case FormatPretty:
		return f.writePretty(entry)
	case FormatRaw:
		return f.writeRaw(entry)
	default:
		return f.writeJSON(entry)
	}
}

func (f *Formatter) writeJSON(entry map[string]interface{}) error {
	b, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("output: marshal error: %w", err)
	}
	_, err = fmt.Fprintln(f.w, string(b))
	return err
}

func (f *Formatter) writePretty(entry map[string]interface{}) error {
	b, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("output: marshal error: %w", err)
	}
	_, err = fmt.Fprintln(f.w, string(b))
	return err
}

func (f *Formatter) writeRaw(entry map[string]interface{}) error {
	parts := make([]string, 0, len(entry))
	for k, v := range entry {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	_, err = fmt.Fprintln(f.w, strings.Join(parts, " "))
	return err
}

func selectFields(entry map[string]interface{}, fields []string) map[string]interface{} {
	selected := make(map[string]interface{}, len(fields))
	for _, field := range fields {
		if val, ok := entry[field]; ok {
			selected[field] = val
		}
	}
	return selected
}
