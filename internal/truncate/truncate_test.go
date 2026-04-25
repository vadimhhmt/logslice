package truncate

import (
	"strings"
	"testing"
)

func entry(pairs ...interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestTruncator_NoopWhenZeroMaxLen(t *testing.T) {
	tr := New(0, "...", nil)
	in := entry("msg", strings.Repeat("a", 200))
	out := tr.Apply(in)
	if out["msg"] != in["msg"] {
		t.Errorf("expected no truncation, got %v", out["msg"])
	}
}

func TestTruncator_ShortValueUnchanged(t *testing.T) {
	tr := New(50, "...", nil)
	in := entry("msg", "hello")
	out := tr.Apply(in)
	if out["msg"] != "hello" {
		t.Errorf("expected 'hello', got %v", out["msg"])
	}
}

func TestTruncator_LongValueTruncated(t *testing.T) {
	tr := New(10, "...", nil)
	in := entry("msg", strings.Repeat("x", 20))
	out := tr.Apply(in)
	got, ok := out["msg"].(string)
	if !ok {
		t.Fatal("expected string")
	}
	if got != strings.Repeat("x", 10)+"..." {
		t.Errorf("unexpected truncated value: %q", got)
	}
}

func TestTruncator_NonStringUnchanged(t *testing.T) {
	tr := New(5, "...", nil)
	in := entry("count", 42)
	out := tr.Apply(in)
	if out["count"] != 42 {
		t.Errorf("expected 42, got %v", out["count"])
	}
}

func TestTruncator_FieldScopedTruncation(t *testing.T) {
	tr := New(5, "~", []string{"msg"})
	long := strings.Repeat("z", 20)
	in := entry("msg", long, "detail", long)
	out := tr.Apply(in)

	msg, _ := out["msg"].(string)
	if msg != "zzzzz~" {
		t.Errorf("msg: expected truncated, got %q", msg)
	}
	detail, _ := out["detail"].(string)
	if detail != long {
		t.Errorf("detail: expected unchanged, got %q", detail)
	}
}

func TestTruncator_OriginalUnmodified(t *testing.T) {
	tr := New(3, "!", nil)
	in := entry("msg", "abcdef")
	_ = tr.Apply(in)
	if in["msg"] != "abcdef" {
		t.Error("original entry was modified")
	}
}

func TestTruncator_MultiByte(t *testing.T) {
	tr := New(3, "…", nil)
	// each character is 3 bytes in UTF-8
	in := entry("msg", "日本語テスト")
	out := tr.Apply(in)
	got, _ := out["msg"].(string)
	if got != "日本語…" {
		t.Errorf("expected rune-based truncation, got %q", got)
	}
}
