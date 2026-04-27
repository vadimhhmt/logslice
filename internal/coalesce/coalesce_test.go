package coalesce

import (
	"testing"
)

func makeEntry(pairs ...any) Entry {
	e := make(Entry, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		e[pairs[i].(string)] = pairs[i+1]
	}
	return e
}

func TestCoalescer_FirstSourceWins(t *testing.T) {
	c, _ := New("result", []string{"a", "b", "c"})
	out := c.Process(makeEntry("a", "alpha", "b", "beta"))
	if got := out["result"]; got != "alpha" {
		t.Fatalf("expected alpha, got %v", got)
	}
}

func TestCoalescer_FallsBackToSecondSource(t *testing.T) {
	c, _ := New("result", []string{"a", "b"})
	out := c.Process(makeEntry("b", "beta"))
	if got := out["result"]; got != "beta" {
		t.Fatalf("expected beta, got %v", got)
	}
}

func TestCoalescer_EmptyStringSkipped(t *testing.T) {
	c, _ := New("result", []string{"a", "b"})
	out := c.Process(makeEntry("a", "", "b", "fallback"))
	if got := out["result"]; got != "fallback" {
		t.Fatalf("expected fallback, got %v", got)
	}
}

func TestCoalescer_NilValueSkipped(t *testing.T) {
	c, _ := New("result", []string{"a", "b"})
	out := c.Process(makeEntry("a", nil, "b", 42))
	if got := out["result"]; got != 42 {
		t.Fatalf("expected 42, got %v", got)
	}
}

func TestCoalescer_NoMatchReturnsOriginal(t *testing.T) {
	c, _ := New("result", []string{"x", "y"})
	e := makeEntry("a", "hello")
	out := c.Process(e)
	if _, exists := out["result"]; exists {
		t.Fatal("expected result field to be absent")
	}
}

func TestCoalescer_OriginalUnmodified(t *testing.T) {
	c, _ := New("result", []string{"a", "b"})
	e := makeEntry("a", "alpha")
	_ = c.Process(e)
	if _, exists := e["result"]; exists {
		t.Fatal("original entry should not be mutated")
	}
}

func TestNew_EmptyDestReturnsError(t *testing.T) {
	_, err := New("", []string{"a"})
	if err == nil {
		t.Fatal("expected error for empty dest")
	}
}

func TestNew_NoSourcesReturnsError(t *testing.T) {
	_, err := New("result", []string{})
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestNew_DeduplicatesSources(t *testing.T) {
	c, err := New("result", []string{"a", "a", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.sources) != 2 {
		t.Fatalf("expected 2 unique sources, got %d", len(c.sources))
	}
}

func TestCoalescer_NonStringValuePreserved(t *testing.T) {
	c, _ := New("result", []string{"score"})
	out := c.Process(makeEntry("score", 3.14))
	if got := out["result"]; got != 3.14 {
		t.Fatalf("expected 3.14, got %v", got)
	}
}
