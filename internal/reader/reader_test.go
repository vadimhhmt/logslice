package reader

import (
	"strings"
	"testing"
)

func collectLines(r *Reader) []string {
	r.Start()
	var out []string
	for line := range r.Lines() {
		out = append(out, line)
	}
	return out
}

func TestReader_BasicLines(t *testing.T) {
	input := `{"ts":"2024-01-01T00:00:00Z","msg":"hello"}
{"ts":"2024-01-01T00:00:01Z","msg":"world"}`

	r := New(strings.NewReader(input), Options{})
	lines := collectLines(r)

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if err := r.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReader_SkipsEmptyLines(t *testing.T) {
	input := "{\"msg\":\"a\"}\n\n{\"msg\":\"b\"}\n\n"

	r := New(strings.NewReader(input), Options{})
	lines := collectLines(r)

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d; lines: %v", len(lines), lines)
	}
}

func TestReader_EmptyInput(t *testing.T) {
	r := New(strings.NewReader(""), Options{})
	lines := collectLines(r)

	if len(lines) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(lines))
	}
	if err := r.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReader_SingleLine(t *testing.T) {
	input := `{"level":"info","msg":"only one"}`

	r := New(strings.NewReader(input), Options{})
	lines := collectLines(r)

	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0] != input {
		t.Fatalf("unexpected line content: %q", lines[0])
	}
}

func TestReader_CustomMaxLineBytes(t *testing.T) {
	input := `{"msg":"short"}`

	r := New(strings.NewReader(input), Options{MaxLineBytes: 256})
	lines := collectLines(r)

	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
}

func TestReader_LineTooLong(t *testing.T) {
	long := strings.Repeat("x", 200)
	r := New(strings.NewReader(long), Options{MaxLineBytes: 64})
	collectLines(r)

	if err := r.Err(); err == nil {
		t.Fatal("expected error for line exceeding MaxLineBytes, got nil")
	}
}

func TestReader_MultipleCallsToErrReturnSameError(t *testing.T) {
	long := strings.Repeat("x", 200)
	r := New(strings.NewReader(long), Options{MaxLineBytes: 64})
	collectLines(r)

	err1 := r.Err()
	err2 := r.Err()
	if err1 == nil {
		t.Fatal("expected error for line exceeding MaxLineBytes, got nil")
	}
	if err1 != err2 {
		t.Fatalf("expected Err() to return the same error on repeated calls, got %v and %v", err1, err2)
	}
}
