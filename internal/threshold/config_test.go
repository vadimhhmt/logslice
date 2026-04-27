package threshold_test

import (
	"flag"
	"testing"

	"logslice/internal/threshold"
)

func newFlagSet() (*flag.FlagSet, *threshold.Config) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg := &threshold.Config{}
	threshold.RegisterFlags(fs, cfg)
	return fs, cfg
}

func TestConfig_DisabledByDefault(t *testing.T) {
	_, cfg := newFlagSet()
	checker, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if checker != nil {
		t.Fatal("expected nil checker when disabled")
	}
}

func TestConfig_EnabledBuildsChecker(t *testing.T) {
	fs, cfg := newFlagSet()
	_ = fs.Parse([]string{"--threshold", "--threshold-field", "latency", "--threshold-min", "5", "--threshold-max", "200"})
	checker, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if checker == nil {
		t.Fatal("expected non-nil checker")
	}
}

func TestConfig_MissingFieldReturnsError(t *testing.T) {
	fs, cfg := newFlagSet()
	_ = fs.Parse([]string{"--threshold", "--threshold-min", "5"})
	_, err := cfg.Build()
	if err == nil {
		t.Fatal("expected error when field is missing")
	}
}

func TestConfig_InvalidMinReturnsError(t *testing.T) {
	fs, cfg := newFlagSet()
	_ = fs.Parse([]string{"--threshold", "--threshold-field", "latency", "--threshold-min", "abc"})
	_, err := cfg.Build()
	if err == nil {
		t.Fatal("expected error for non-numeric min")
	}
}

func TestConfig_InvalidMaxReturnsError(t *testing.T) {
	fs, cfg := newFlagSet()
	_ = fs.Parse([]string{"--threshold", "--threshold-field", "latency", "--threshold-max", "xyz"})
	_, err := cfg.Build()
	if err == nil {
		t.Fatal("expected error for non-numeric max")
	}
}
