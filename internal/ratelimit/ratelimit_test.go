package ratelimit_test

import (
	"testing"
	"time"

	"logslice/internal/parser"
	"logslice/internal/ratelimit"
)

func makeEntry(ts time.Time) parser.Entry {
	return parser.Entry{
		Timestamp: ts,
		Raw:       `{"ts":"" }`,
		Fields:    map[string]interface{}{},
	}
}

func TestLimiter_AllowsUpToMax(t *testing.T) {
	l, err := ratelimit.New(3, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 3; i++ {
		if !l.Allow(makeEntry(base.Add(time.Duration(i) * time.Second))) {
			t.Errorf("entry %d should be allowed", i)
		}
	}
	if l.Dropped != 0 {
		t.Errorf("expected 0 dropped, got %d", l.Dropped)
	}
}

func TestLimiter_DropsOverMax(t *testing.T) {
	l, _ := ratelimit.New(2, time.Minute)
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	l.Allow(makeEntry(base))
	l.Allow(makeEntry(base.Add(10 * time.Second)))
	got := l.Allow(makeEntry(base.Add(20 * time.Second)))
	if got {
		t.Error("third entry should be dropped")
	}
	if l.Dropped != 1 {
		t.Errorf("expected 1 dropped, got %d", l.Dropped)
	}
}

func TestLimiter_ResetsAfterWindow(t *testing.T) {
	l, _ := ratelimit.New(1, time.Minute)
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if !l.Allow(makeEntry(base)) {
		t.Fatal("first entry should be allowed")
	}
	if l.Allow(makeEntry(base.Add(10 * time.Second))) {
		t.Fatal("second entry in same window should be dropped")
	}
	// New window.
	if !l.Allow(makeEntry(base.Add(61 * time.Second))) {
		t.Error("entry in new window should be allowed")
	}
}

func TestLimiter_ZeroTimestampAlwaysAllowed(t *testing.T) {
	l, _ := ratelimit.New(1, time.Minute)
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	l.Allow(makeEntry(base)) // fill bucket
	if !l.Allow(makeEntry(time.Time{})) {
		t.Error("zero-timestamp entry should always be allowed")
	}
}

func TestLimiter_InvalidArgs(t *testing.T) {
	if _, err := ratelimit.New(0, time.Minute); err == nil {
		t.Error("expected error for maxPerWindow=0")
	}
	if _, err := ratelimit.New(1, 0); err == nil {
		t.Error("expected error for window=0")
	}
}
