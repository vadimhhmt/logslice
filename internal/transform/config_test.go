package transform_test

import (
	"flag"
	"testing"

	"github.com/yourorg/logslice/internal/transform"
)

func newFlagSet() (*flag.FlagSet, *transform.RuleFlags) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	rf := &transform.RuleFlags{}
	transform.RegisterFlags(fs, rf)
	return fs, rf
}

func TestConfig_NoFlags(t *testing.T) {
	_, rf := newFlagSet()
	rules, err := rf.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected 0 rules, got %d", len(rules))
	}
}

func TestConfig_RenameFlag(t *testing.T) {
	fs, rf := newFlagSet()
	if err := fs.Parse([]string{"-rename", "msg=message"}); err != nil {
		t.Fatal(err)
	}
	rules, err := rf.Build()
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 || rules[0].FromField != "msg" || rules[0].ToField != "message" {
		t.Errorf("unexpected rules: %+v", rules)
	}
}

func TestConfig_RemapFlag(t *testing.T) {
	fs, rf := newFlagSet()
	if err := fs.Parse([]string{"-remap", "level:warn=warning"}); err != nil {
		t.Fatal(err)
	}
	rules, err := rf.Build()
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].ValueMap["warn"] != "warning" {
		t.Errorf("unexpected value map: %v", rules[0].ValueMap)
	}
}

func TestConfig_InvalidRename(t *testing.T) {
	fs, _ := newFlagSet()
	if err := fs.Parse([]string{"-rename", "badvalue"}); err == nil {
		t.Error("expected error for missing '=' in rename")
	}
}

func TestConfig_InvalidRemap(t *testing.T) {
	fs, _ := newFlagSet()
	if err := fs.Parse([]string{"-remap", "nocodon"}); err == nil {
		t.Error("expected error for malformed remap")
	}
}
