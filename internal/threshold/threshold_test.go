package threshold_test

import (
	"testing"
	"time"

	"logslice/internal/parser"
	"logslice/internal/threshold"
)

func makeEntry(field string, value any) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Raw:       "{}",
		Fields:    map[string]any{field: value},
	}
}

func ptr(f float64) *float64 { return &f }

func TestChecker_AboveMin(t *testing.T) {
	c, _ := threshold.New("latency", ptr(10), nil)
	if !c.Allow(makeEntry("latency", float64(20))) {
		t.Fatal("expected entry above min to be allowed")
	}
}

func TestChecker_BelowMin(t *testing.T) {
	c, _ := threshold.New("latency", ptr(10), nil)
	if c.Allow(makeEntry("latency", float64(5))) {
		t.Fatal("expected entry below min to be dropped")
	}
}

func TestChecker_AboveMax(t *testing.T) {
	c, _ := threshold.New("latency", nil, ptr(100))
	if c.Allow(makeEntry("latency", float64(200))) {
		t.Fatal("expected entry above max to be dropped")
	}
}

func TestChecker_WithinRange(t *testing.T) {
	c, _ := threshold.New("latency", ptr(10), ptr(100))
	if !c.Allow(makeEntry("latency", float64(50))) {
		t.Fatal("expected entry within range to be allowed")
	}
}

func TestChecker_MissingFieldPassesThrough(t *testing.T) {
	c, _ := threshold.New("latency", ptr(10), ptr(100))
	e := parser.Entry{Timestamp: time.Now(), Raw: "{}", Fields: map[string]any{}}
	if !c.Allow(e) {
		t.Fatal("expected missing field to pass through")
	}
}

func TestChecker_NonNumericFieldPassesThrough(t *testing.T) {
	c, _ := threshold.New("latency", ptr(10), ptr(100))
	if !c.Allow(makeEntry("latency", "fast")) {
		t.Fatal("expected non-numeric string to pass through")
	}
}

func TestChecker_StringNumericValue(t *testing.T) {
	c, _ := threshold.New("latency", ptr(10), ptr(100))
	if !c.Allow(makeEntry("latency", "50")) {
		t.Fatal("expected parseable string number to be evaluated")
	}
	if c.Allow(makeEntry("latency", "200")) {
		t.Fatal("expected out-of-range string number to be dropped")
	}
}

func TestNew_EmptyFieldReturnsError(t *testing.T) {
	_, err := threshold.New("", ptr(0), nil)
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNew_MinGreaterThanMaxReturnsError(t *testing.T) {
	_, err := threshold.New("latency", ptr(100), ptr(10))
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}
