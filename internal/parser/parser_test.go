package parser

import (
	"testing"
	"time"
)

func TestParseLine_RFC3339(t *testing.T) {
	line := `{"time":"2024-03-15T10:00:00Z","level":"info","msg":"started"}`
	entry, err := ParseLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)
	if !entry.Timestamp.Equal(expected) {
		t.Errorf("timestamp: got %v, want %v", entry.Timestamp, expected)
	}
	if entry.Fields["level"] != "info" {
		t.Errorf("level field: got %v", entry.Fields["level"])
	}
}

func TestParseLine_UnixEpoch(t *testing.T) {
	line := `{"ts":1710496800.5,"msg":"heartbeat"}`
	entry, err := ParseLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Timestamp.Unix() != 1710496800 {
		t.Errorf("unix timestamp mismatch: got %d", entry.Timestamp.Unix())
	}
}

func TestParseLine_AtTimestamp(t *testing.T) {
	line := `{"@timestamp":"2024-03-15T12:30:00.123456789Z","service":"api"}`
	entry, err := ParseLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestParseLine_InvalidJSON(t *testing.T) {
	_, err := ParseLine("not json at all")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseLine_EmptyLine(t *testing.T) {
	_, err := ParseLine("")
	if err == nil {
		t.Fatal("expected error for empty line")
	}
}

func TestParseLine_NoTimestamp(t *testing.T) {
	line := `{"level":"warn","msg":"no time here"}`
	_, err := ParseLine(line)
	if err == nil {
		t.Fatal("expected error when no timestamp field present")
	}
}

func TestParseLine_BadTimeFormat(t *testing.T) {
	line := `{"time":"March 15 2024","msg":"bad format"}`
	_, err := ParseLine(line)
	if err == nil {
		t.Fatal("expected error for unrecognised time format")
	}
}

func TestParseLine_RawPreserved(t *testing.T) {
	line := `{"time":"2024-03-15T10:00:00Z","x":1}`
	entry, err := ParseLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Raw != line {
		t.Errorf("raw mismatch: got %q, want %q", entry.Raw, line)
	}
}
