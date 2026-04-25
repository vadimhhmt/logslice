package burst

import (
	"flag"
	"testing"
	"time"
)

func newFlagSet() *flag.FlagSet {
	return flag.NewFlagSet("test", flag.ContinueOnError)
}

func TestConfig_DisabledByDefault(t *testing.T) {
	fs := newFlagSet()
	cfg := RegisterFlags(fs)
	_ = fs.Parse([]string{})

	if cfg.Enabled {
		t.Fatal("expected burst to be disabled by default")
	}

	d, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != nil {
		t.Fatal("expected nil detector when disabled")
	}
}

func TestConfig_EnabledBuildsDetector(t *testing.T) {
	fs := newFlagSet()
	cfg := RegisterFlags(fs)
	_ = fs.Parse([]string{"--burst", "--burst-window=2s", "--burst-threshold=5"})

	if !cfg.Enabled {
		t.Fatal("expected burst to be enabled")
	}
	if cfg.Window != 2*time.Second {
		t.Fatalf("expected window 2s, got %s", cfg.Window)
	}
	if cfg.Threshold != 5 {
		t.Fatalf("expected threshold 5, got %d", cfg.Threshold)
	}

	d, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil detector")
	}
}

func TestConfig_InvalidThresholdReturnsError(t *testing.T) {
	fs := newFlagSet()
	cfg := RegisterFlags(fs)
	_ = fs.Parse([]string{"--burst", "--burst-threshold=0"})

	_, err := cfg.Build()
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestConfig_InvalidWindowReturnsError(t *testing.T) {
	fs := newFlagSet()
	cfg := RegisterFlags(fs)
	_ = fs.Parse([]string{"--burst"})
	cfg.Window = 0

	_, err := cfg.Build()
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestConfig_GroupByField(t *testing.T) {
	fs := newFlagSet()
	cfg := RegisterFlags(fs)
	_ = fs.Parse([]string{"--burst", "--burst-group-by=service"})

	if cfg.GroupBy != "service" {
		t.Fatalf("expected group-by 'service', got %q", cfg.GroupBy)
	}

	d, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil detector")
	}
}
