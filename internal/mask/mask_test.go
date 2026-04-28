package mask_test

import (
	"testing"

	"github.com/logslice/logslice/internal/mask"
)

func makeEntry(pairs ...any) map[string]any {
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestMasker_TargetedFieldReplaced(t *testing.T) {
	m, err := mask.New([]string{"password"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := m.Apply(makeEntry("user", "alice", "password", "s3cr3t"))
	if out["password"] != "***" {
		t.Errorf("expected placeholder, got %q", out["password"])
	}
	if out["user"] != "alice" {
		t.Errorf("expected user unchanged, got %q", out["user"])
	}
}

func TestMasker_UnrelatedFieldsUntouched(t *testing.T) {
	m, _ := mask.New([]string{"token"})
	out := m.Apply(makeEntry("msg", "hello", "level", "info"))
	if out["msg"] != "hello" || out["level"] != "info" {
		t.Errorf("unrelated fields should be unchanged: %v", out)
	}
}

func TestMasker_MissingFieldIgnored(t *testing.T) {
	m, _ := mask.New([]string{"secret"})
	out := m.Apply(makeEntry("msg", "hi"))
	if _, exists := out["secret"]; exists {
		t.Error("missing field should not be injected")
	}
}

func TestMasker_CustomPlaceholder(t *testing.T) {
	m, _ := mask.New([]string{"pin"}, mask.WithPlaceholder("[REDACTED]"))
	out := m.Apply(makeEntry("pin", "1234"))
	if out["pin"] != "[REDACTED]" {
		t.Errorf("expected custom placeholder, got %q", out["pin"])
	}
}

func TestMasker_OriginalUnmodified(t *testing.T) {
	m, _ := mask.New([]string{"key"})
	orig := makeEntry("key", "value", "other", 42)
	m.Apply(orig)
	if orig["key"] != "value" {
		t.Error("Apply must not mutate the original entry")
	}
}

func TestMasker_MultipleFields(t *testing.T) {
	m, _ := mask.New([]string{"a", "b"})
	out := m.Apply(makeEntry("a", 1, "b", 2, "c", 3))
	if out["a"] != "***" || out["b"] != "***" {
		t.Errorf("both fields should be masked: %v", out)
	}
	if out["c"] != 3 {
		t.Error("c should be unchanged")
	}
}

func TestNew_NoFieldsReturnsError(t *testing.T) {
	_, err := mask.New(nil)
	if err == nil {
		t.Error("expected error for empty field list")
	}
}

func TestNew_EmptyFieldNameReturnsError(t *testing.T) {
	_, err := mask.New([]string{""})
	if err == nil {
		t.Error("expected error for empty field name")
	}
}
