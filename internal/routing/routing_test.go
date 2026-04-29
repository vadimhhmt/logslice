package routing_test

import (
	"testing"
	"time"

	"logslice/internal/parser"
	"logslice/internal/routing"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{Timestamp: time.Now(), Fields: fields}
}

func TestRouter_FirstMatchWins(t *testing.T) {
	r := routing.New(nil, "default")
	_ = r.AddRule("level", "error", "errors")
	_ = r.AddRule("level", "warn", "warnings")

	if got := r.Route(makeEntry(map[string]any{"level": "error"})); got != "errors" {
		t.Fatalf("expected errors, got %q", got)
	}
	if got := r.Route(makeEntry(map[string]any{"level": "warn"})); got != "warnings" {
		t.Fatalf("expected warnings, got %q", got)
	}
}

func TestRouter_DefaultBucketWhenNoMatch(t *testing.T) {
	r := routing.New(nil, "catch-all")
	_ = r.AddRule("level", "error", "errors")

	got := r.Route(makeEntry(map[string]any{"level": "info"}))
	if got != "catch-all" {
		t.Fatalf("expected catch-all, got %q", got)
	}
}

func TestRouter_EmptyDefaultDropsEntry(t *testing.T) {
	r := routing.New(nil, "")
	got := r.Route(makeEntry(map[string]any{"level": "debug"}))
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestRouter_MissingFieldSkipsRule(t *testing.T) {
	r := routing.New(nil, "default")
	_ = r.AddRule("service", "auth", "auth-bucket")

	got := r.Route(makeEntry(map[string]any{"level": "info"}))
	if got != "default" {
		t.Fatalf("expected default, got %q", got)
	}
}

func TestRouter_AddRule_InvalidPatternReturnsError(t *testing.T) {
	r := routing.New(nil, "default")
	if err := r.AddRule("level", "[invalid", "bucket"); err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestRouter_Dispatch(t *testing.T) {
	r := routing.New(nil, "other")
	_ = r.AddRule("level", "error", "errors")

	in := make(chan parser.Entry, 4)
	errCh := make(chan parser.Entry, 4)
	otherCh := make(chan parser.Entry, 4)

	buckets := map[string]chan<- parser.Entry{
		"errors": errCh,
		"other":  otherCh,
	}

	in <- makeEntry(map[string]any{"level": "error"})
	in <- makeEntry(map[string]any{"level": "info"})
	in <- makeEntry(map[string]any{"level": "error"})
	close(in)

	r.Dispatch(in, buckets)

	if len(errCh) != 2 {
		t.Fatalf("expected 2 error entries, got %d", len(errCh))
	}
	if len(otherCh) != 1 {
		t.Fatalf("expected 1 other entry, got %d", len(otherCh))
	}
}
