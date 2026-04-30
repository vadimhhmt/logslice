package checkpoint_test

import (
	"flag"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/logslice/internal/checkpoint"
)

func newFlagSet(t *testing.T) (*flag.FlagSet, *checkpoint.Config) {
	t.Helper()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var cfg checkpoint.Config
	checkpoint.RegisterFlags(fs, &cfg)
	return fs, &cfg
}

func TestConfig_DisabledByDefault(t *testing.T) {
	_, cfg := newFlagSet(t)
	if cfg.Enabled {
		t.Error("expected checkpoint disabled by default")
	}
	m, err := cfg.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if m != nil {
		t.Error("expected nil Manager when disabled")
	}
}

func TestConfig_EnabledBuildsManager(t *testing.T) {
	fs, cfg := newFlagSet(t)
	path := filepath.Join(t.TempDir(), "cp.json")
	_ = fs.Parse([]string{"-checkpoint", "-checkpoint-file", path})

	m, err := cfg.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil Manager")
	}
}

func TestConfig_ResetClearsState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cp.json")

	// Seed some state.
	seed, _ := checkpoint.New(path)
	_ = seed.Record(time.Now())

	fs, cfg := newFlagSet(t)
	_ = fs.Parse([]string{"-checkpoint", "-checkpoint-file", path, "-checkpoint-reset"})

	m, err := cfg.Build()
	if err != nil {
		t.Fatalf("Build with reset: %v", err)
	}
	if got := m.State().LinesProcessed; got != 0 {
		t.Errorf("LinesProcessed after reset = %d, want 0", got)
	}
}

func TestConfig_EmptyPathReturnsError(t *testing.T) {
	_, cfg := newFlagSet(t)
	cfg.Enabled = true
	cfg.Path = ""
	if _, err := cfg.Build(); err == nil {
		t.Error("expected error for empty checkpoint path")
	}
}
