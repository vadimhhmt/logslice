package redact_test

import (
	"flag"
	"testing"

	"logslice/internal/redact"
)

func newFlagSet() *flag.FlagSet {
	return flag.NewFlagSet("test", flag.ContinueOnError)
}

func TestConfig_DisabledByDefault(t *testing.T) {
	var cfg redact.Config
	fs := newFlagSet()
	redact.RegisterFlags(fs, &cfg)
	_ = fs.Parse([]string{})
	cfg.Finalise()
	r, err := cfg.BuildRedactor()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r != nil {
		t.Error("expected nil redactor when disabled")
	}
}

func TestConfig_EnabledNoExtras(t *testing.T) {
	var cfg redact.Config
	fs := newFlagSet()
	redact.RegisterFlags(fs, &cfg)
	_ = fs.Parse([]string{"-redact"})
	cfg.Finalise()
	r, err := cfg.BuildRedactor()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil redactor")
	}
	// Default keys should still redact password.
	out := r.Apply(map[string]any{"password": "x", "msg": "hi"})
	if out["password"] != "[REDACTED]" {
		t.Errorf("expected password redacted, got %v", out["password"])
	}
}

func TestConfig_ExtraKeys(t *testing.T) {
	var cfg redact.Config
	fs := newFlagSet()
	redact.RegisterFlags(fs, &cfg)
	_ = fs.Parse([]string{"-redact", "-redact-keys", "ssn, dob"})
	cfg.Finalise()
	r, err := cfg.BuildRedactor()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := r.Apply(map[string]any{"ssn": "123", "dob": "1990", "name": "Eve"})
	if out["ssn"] != "[REDACTED]" {
		t.Errorf("expected ssn redacted")
	}
	if out["dob"] != "[REDACTED]" {
		t.Errorf("expected dob redacted")
	}
	if out["name"] != "Eve" {
		t.Errorf("expected name unchanged")
	}
}

func TestConfig_InvalidPattern(t *testing.T) {
	var cfg redact.Config
	fs := newFlagSet()
	redact.RegisterFlags(fs, &cfg)
	_ = fs.Parse([]string{"-redact", "-redact-patterns", "[bad"})
	cfg.Finalise()
	_, err := cfg.BuildRedactor()
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestConfig_DisabledIgnoresKeys(t *testing.T) {
	// When redaction is disabled, BuildRedactor should return nil even if
	// extra keys or patterns are provided.
	var cfg redact.Config
	fs := newFlagSet()
	redact.RegisterFlags(fs, &cfg)
	_ = fs.Parse([]string{"-redact-keys", "ssn"})
	cfg.Finalise()
	r, err := cfg.BuildRedactor()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r != nil {
		t.Error("expected nil redactor when disabled, even with extra keys set")
	}
}
