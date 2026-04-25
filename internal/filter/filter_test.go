package filter_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

var (
	base   = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	before = base.Add(-time.Hour)
	after  = base.Add(time.Hour)
)

func TestFilter_TimeRange(t *testing.T) {
	f, err := filter.New(filter.Options{From: base, To: after})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name string
		ts   time.Time
		want bool
	}{
		{"within range", base, true},
		{"at upper bound", after, true},
		{"before range", before, false},
		{"after range", after.Add(time.Second), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := f.Match(`{"level":"info"}`, tc.ts)
			if got != tc.want {
				t.Errorf("Match() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFilter_Pattern(t *testing.T) {
	f, err := filter.New(filter.Options{Pattern: `"level":"error"`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !f.Match(`{"level":"error","msg":"oops"}`, base) {
		t.Error("expected match for error line")
	}
	if f.Match(`{"level":"info","msg":"ok"}`, base) {
		t.Error("expected no match for info line")
	}
}

func TestFilter_InvalidPattern(t *testing.T) {
	_, err := filter.New(filter.Options{Pattern: `[invalid`})
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}

func TestFilter_NoConstraints(t *testing.T) {
	f, err := filter.New(filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match(`any line`, base) {
		t.Error("expected match with no constraints")
	}
}

func TestFilter_PatternAndTimeRange(t *testing.T) {
	f, err := filter.New(filter.Options{
		Pattern: `"level":"error"`,
		From:    base,
		To:      after,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name string
		line string
		ts   time.Time
		want bool
	}{
		{"matching pattern and in range", `{"level":"error"}`, base, true},
		{"matching pattern but out of range", `{"level":"error"}`, before, false},
		{"in range but pattern mismatch", `{"level":"info"}`, base, false},
		{"both mismatched", `{"level":"info"}`, before, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := f.Match(tc.line, tc.ts)
			if got != tc.want {
				t.Errorf("Match() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestInRange(t *testing.T) {
	if !filter.InRange(base, before, after) {
		t.Error("expected base to be in range [before, after]")
	}
	if filter.InRange(before.Add(-time.Second), before, after) {
		t.Error("expected out-of-range timestamp to fail")
	}
	if !filter.InRange(base, time.Time{}, time.Time{}) {
		t.Error("expected unbounded range to always match")
	}
}
