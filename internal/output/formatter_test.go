package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func sampleEntry() map[string]interface{} {
	return map[string]interface{}{
		"timestamp": "2024-01-15T10:00:00Z",
		"level":     "info",
		"message":   "service started",
		"service":   "api",
	}
}

func TestFormatter_JSON(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatJSON, nil)

	if err := f.Write(sampleEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if out["level"] != "info" {
		t.Errorf("expected level=info, got %v", out["level"])
	}
}

func TestFormatter_Pretty(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatPretty, nil)

	if err := f.Write(sampleEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "\n") {
		t.Error("expected pretty output to contain newlines")
	}
	var out map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &out); err != nil {
		t.Fatalf("pretty output is not valid JSON: %v", err)
	}
}

func TestFormatter_Raw(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatRaw, nil)

	if err := f.Write(sampleEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "=") {
		t.Error("expected raw output to contain key=value pairs")
	}
}

func TestFormatter_FieldSelection(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatJSON, []string{"level", "message"})

	if err := f.Write(sampleEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := out["timestamp"]; ok {
		t.Error("timestamp should have been excluded by field selection")
	}
	if _, ok := out["level"]; !ok {
		t.Error("level should be present in field selection")
	}
}

func TestFormatter_EmptyEntry(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatJSON, nil)

	if err := f.Write(map[string]interface{}{}); err != nil {
		t.Fatalf("unexpected error on empty entry: %v", err)
	}

	if strings.TrimSpace(buf.String()) != "{}" {
		t.Errorf("expected '{}', got %q", buf.String())
	}
}
