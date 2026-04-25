package sample_test

import (
	"testing"
	"time"

	"logslice/internal/parser"
	"logslice/internal/sample"
)

func makeEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Raw:       []byte(`{"msg":"` + msg + `"}`),
		Fields:    map[string]any{"msg": msg},
	}
}

func TestSampler_RateKeepsEveryNth(t *testing.T) {
	s, err := sample.New(sample.StrategyRate, 3, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var kept int
	for i := 0; i < 9; i++ {
		if s.Accept(makeEntry("x")) {
			kept++
		}
	}
	if kept != 3 {
		t.Errorf("expected 3 kept entries, got %d", kept)
	}
}

func TestSampler_RateInvalidRate(t *testing.T) {
	_, err := sample.New(sample.StrategyRate, 0, 0, 0)
	if err == nil {
		t.Error("expected error for rate=0")
	}
}

func TestSampler_ReservoirFillsToCapacity(t *testing.T) {
	s, err := sample.New(sample.StrategyReservoir, 0, 5, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 3; i++ {
		s.Collect(makeEntry("a"))
	}
	out := s.Flush()
	if len(out) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out))
	}
}

func TestSampler_ReservoirCapsAtSize(t *testing.T) {
	s, err := sample.New(sample.StrategyReservoir, 0, 5, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 100; i++ {
		s.Collect(makeEntry("b"))
	}
	out := s.Flush()
	if len(out) != 5 {
		t.Errorf("expected reservoir size 5, got %d", len(out))
	}
}

func TestSampler_ReservoirFlushResetsState(t *testing.T) {
	s, err := sample.New(sample.StrategyReservoir, 0, 4, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 10; i++ {
		s.Collect(makeEntry("c"))
	}
	s.Flush()
	out := s.Flush()
	if len(out) != 0 {
		t.Errorf("expected empty reservoir after double flush, got %d", len(out))
	}
}

func TestSampler_ReservoirInvalidSize(t *testing.T) {
	_, err := sample.New(sample.StrategyReservoir, 0, 0, 0)
	if err == nil {
		t.Error("expected error for size=0")
	}
}

func TestSampler_AcceptIgnoredForReservoir(t *testing.T) {
	s, err := sample.New(sample.StrategyReservoir, 0, 5, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Accept(makeEntry("x")) {
		t.Error("Accept should return false for reservoir strategy")
	}
}
