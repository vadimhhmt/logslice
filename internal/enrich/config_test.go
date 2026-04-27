package enrich

import (
	"flag"
	"testing"
)

func newFlagSet() *flag.FlagSet {
	return flag.NewFlagSet("test", flag.ContinueOnError)
}

func TestConfig_DisabledByDefault(t *testing.T) {
	fs := newFlagSet()
	build := RegisterFlags(fs)
	_ = fs.Parse([]string{})
	e, err := build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e != nil {
		t.Fatal("expected nil enricher when no flags set")
	}
}

func TestConfig_ValidUpperRule(t *testing.T) {
	fs := newFlagSet()
	build := RegisterFlags(fs)
	_ = fs.Parse([]string{"-enrich", "level:level_upper:upper"})
	e, err := build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil enricher")
	}
}

func TestConfig_ValidLowerRule(t *testing.T) {
	fs := newFlagSet()
	build := RegisterFlags(fs)
	_ = fs.Parse([]string{"-enrich", "host:host_lo:lower"})
	e, err := build()
	if err != nil || e == nil {
		t.Fatalf("unexpected result: err=%v enricher=%v", err, e)
	}
}

func TestConfig_InvalidRuleMissingParts(t *testing.T) {
	fs := newFlagSet()
	build := RegisterFlags(fs)
	_ = fs.Parse([]string{"-enrich", "level:upper"})
	_, err := build()
	if err == nil {
		t.Fatal("expected error for malformed rule")
	}
}

func TestConfig_UnknownFnReturnsError(t *testing.T) {
	fs := newFlagSet()
	build := RegisterFlags(fs)
	_ = fs.Parse([]string{"-enrich", "level:level_x:rot13"})
	_, err := build()
	if err == nil {
		t.Fatal("expected error for unknown fn")
	}
}

func TestConfig_MultipleRules(t *testing.T) {
	fs := newFlagSet()
	build := RegisterFlags(fs)
	_ = fs.Parse([]string{"-enrich", "level:level_up:upper", "-enrich", "host:host_lo:lower"})
	e, err := build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil enricher")
	}
	if len(e.rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(e.rules))
	}
}
