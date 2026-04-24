package parser

import (
	"encoding/json"
	"fmt"
	"time"
)

// LogEntry represents a single parsed log line.
type LogEntry struct {
	Timestamp time.Time
	Fields    map[string]interface{}
	Raw       string
}

// TimeFields is the ordered list of JSON keys tried when extracting a timestamp.
var TimeFields = []string{"time", "timestamp", "ts", "@timestamp"}

// TimeFormats is the ordered list of time layouts attempted during parsing.
var TimeFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
}

// ParseLine parses a single JSON log line into a LogEntry.
// Returns an error if the line is not valid JSON or no timestamp field is found.
func ParseLine(line string) (*LogEntry, error) {
	if len(line) == 0 {
		return nil, fmt.Errorf("empty line")
	}

	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return nil, fmt.Errorf("json parse error: %w", err)
	}

	ts, err := extractTimestamp(fields)
	if err != nil {
		return nil, err
	}

	return &LogEntry{
		Timestamp: ts,
		Fields:    fields,
		Raw:       line,
	}, nil
}

// Field returns the value of a named field from the log entry, along with a
// boolean indicating whether the field was present.
func (e *LogEntry) Field(name string) (interface{}, bool) {
	v, ok := e.Fields[name]
	return v, ok
}

// FieldString returns the string value of a named field. If the field is
// absent or not a string, it returns ("", false).
func (e *LogEntry) FieldString(name string) (string, bool) {
	v, ok := e.Fields[name]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func extractTimestamp(fields map[string]interface{}) (time.Time, error) {
	for _, key := range TimeFields {
		val, ok := fields[key]
		if !ok {
			continue
		}
		switch v := val.(type) {
		case string:
			for _, layout := range TimeFormats {
				if t, err := time.Parse(layout, v); err == nil {
					return t, nil
				}
			}
			return time.Time{}, fmt.Errorf("unrecognised time format for field %q: %s", key, v)
		case float64:
			// Unix epoch seconds (possibly fractional)
			sec := int64(v)
			nsec := int64((v - float64(sec)) * 1e9)
			return time.Unix(sec, nsec).UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("no recognised timestamp field found")
}
