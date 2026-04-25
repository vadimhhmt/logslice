package ratelimit_test

import (
	"flag"
	"testing"
	"time"

	"logslice/internal/ratelimit"
)

func newFlagSet(t *testing.T) (*flag.FlagSet, *ratelimit.Config) {
	t.Helper()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg := &ratelimit.Config{}
	ratelimit.RegisterFlags(fs, cfg)
	return fs, cfg
}

func TestConfig_DisabledByDefault(t *testing.T) {
	_, cfg := newFlagSet(t)
	l, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l != nil {
		t.Error("expected nil limiter when disabled")
	}
}

func TestConfig_EnabledBuildsLimiter(t *testing.T) {
	fs, cfg := newFlagSet(t)
	if err := fs.Parse([]string{"--ratelimit", "--ratelimit-max=50", "--ratelimit-window=30s"}); err != nil {
		t.Fatalf("parse error: %v", err)
	}
	l, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
}

func TestConfig_InvalidMaxReturnsError(t *testing.T) {
	_, cfg := newFlagSet(t)
	cfg.Enabled = true
	cfg.MaxPerWindow = 0
	cfg.Window = time.Minute
	if _, err := cfg.Build(); err == nil {
		t.Error("expected error for MaxPerWindow=0")
	}
}

func TestConfig_InvalidWindowReturnsError(t *testing.T) {
	_, cfg := newFlagSet(t)
	cfg.Enabled = true
	cfg.MaxPerWindow = 10
	cfg.Window = 0
	if _, err := cfg.Build(); err == nil {
		t.Error("expected error for Window=0")
	}
}
