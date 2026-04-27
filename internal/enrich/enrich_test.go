package enrich

import (
	"strings"
	"testing"
	"time"

	"logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Fields:    fields,
		Raw:       "{}",
	}
}

func TestEnricher_AddsField(t *testing.T) {
	e := New([]Rule{{From: "level", To: "level_upper", Fn: UpperCase}})
	out := e.Apply(makeEntry(map[string]any{"level": "warn"}))
	if got := out.Fields["level_upper"]; got != "WARN" {
		t.Fatalf("expected WARN, got %v", got)
	}
}

func TestEnricher_MissingSourceSkipped(t *testing.T) {
	e := New([]Rule{{From: "missing", To: "out", Fn: UpperCase}})
	out := e.Apply(makeEntry(map[string]any{"level": "info"}))
	if _, ok := out.Fields["out"]; ok {
		t.Fatal("expected field to be absent")
	}
}

func TestEnricher_OverwritesExistingTarget(t *testing.T) {
	e := New([]Rule{{From: "env", To: "env", Fn: LowerCase}})
	out := e.Apply(makeEntry(map[string]any{"env": "PROD"}))
	if got := out.Fields["env"]; got != "prod" {
		t.Fatalf("expected prod, got %v", got)
	}
}

func TestEnricher_OriginalUnmodified(t *testing.T) {
	e := New([]Rule{{From: "level", To: "level_upper", Fn: UpperCase}})
	orig := makeEntry(map[string]any{"level": "debug"})
	_ = e.Apply(orig)
	if _, ok := orig.Fields["level_upper"]; ok {
		t.Fatal("original entry must not be modified")
	}
}

func TestEnricher_FnReturnsFalseSkips(t *testing.T) {
	neverMatch := func(string) (string, bool) { return "", false }
	e := New([]Rule{{From: "level", To: "out", Fn: neverMatch}})
	out := e.Apply(makeEntry(map[string]any{"level": "info"}))
	if _, ok := out.Fields["out"]; ok {
		t.Fatal("expected field absent when Fn returns false")
	}
}

func TestEnricher_MultipleRulesApplied(t *testing.T) {
	rules := []Rule{
		{From: "level", To: "level_up", Fn: UpperCase},
		{From: "host", To: "host_lo", Fn: LowerCase},
	}
	e := New(rules)
	out := e.Apply(makeEntry(map[string]any{"level": "info", "host": "WEB01"}))
	if got := out.Fields["level_up"]; got != "INFO" {
		t.Fatalf("level_up: want INFO got %v", got)
	}
	if got := out.Fields["host_lo"]; got != "web01" {
		t.Fatalf("host_lo: want web01 got %v", got)
	}
}

func TestEnricher_CustomFn(t *testing.T) {
	prefix := func(v string) (string, bool) {
		return "svc-" + strings.TrimSpace(v), true
	}
	e := New([]Rule{{From: "name", To: "svc_name", Fn: prefix}})
	out := e.Apply(makeEntry(map[string]any{"name": "auth"}))
	if got := out.Fields["svc_name"]; got != "svc-auth" {
		t.Fatalf("expected svc-auth, got %v", got)
	}
}
