package config

import (
	"flag"
	"os"
	"testing"
	"time"
)

// resetFlags resets flag.CommandLine so each sub-test gets a clean slate.
func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

func TestConfig_Defaults(t *testing.T) {
	resetFlags()
	cfg := &Config{Format: "json"}
	if err := cfg.validate("", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.From != (time.Time{}) {
		t.Errorf("expected zero From, got %v", cfg.From)
	}
}

func TestConfig_ValidTimeRange(t *testing.T) {
	resetFlags()
	cfg := &Config{Format: "json"}
	err := cfg.validate("2024-01-01T00:00:00Z", "2024-01-02T00:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.From.IsZero() || cfg.To.IsZero() {
		t.Error("expected non-zero From and To")
	}
}

func TestConfig_ToBeforeFrom(t *testing.T) {
	resetFlags()
	cfg := &Config{Format: "json"}
	err := cfg.validate("2024-01-02T00:00:00Z", "2024-01-01T00:00:00Z")
	if err == nil {
		t.Fatal("expected error for to < from")
	}
}

func TestConfig_InvalidFrom(t *testing.T) {
	resetFlags()
	cfg := &Config{Format: "json"}
	err := cfg.validate("not-a-date", "")
	if err == nil {
		t.Fatal("expected error for invalid from")
	}
}

func TestConfig_InvalidFormat(t *testing.T) {
	resetFlags()
	cfg := &Config{Format: "xml"}
	err := cfg.validate("", "")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestConfig_AllFormats(t *testing.T) {
	for _, f := range []string{"json", "pretty", "raw"} {
		resetFlags()
		cfg := &Config{Format: f}
		if err := cfg.validate("", ""); err != nil {
			t.Errorf("format %q should be valid, got: %v", f, err)
		}
	}
}
