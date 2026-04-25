package sanitize_test

import (
	"strings"
	"testing"

	"github.com/example/logslice/internal/sanitize"
)

func TestLine_Empty(t *testing.T) {
	if got := sanitize.Line(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestLine_Whitespace(t *testing.T) {
	if got := sanitize.Line("   \t  "); got != "" {
		t.Fatalf("expected empty string for whitespace-only input, got %q", got)
	}
}

func TestLine_Trim(t *testing.T) {
	const want = `{"level":"info"}`
	got := sanitize.Line("  " + want + "\n")
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestLine_ControlCharsRemoved(t *testing.T) {
	// Embed a BEL (\x07) and a NUL (\x00) character.
	raw := "{\"msg\":\"hello\x07world\x00\"}"
	got := sanitize.Line(raw)
	if strings.ContainsAny(got, "\x07\x00") {
		t.Fatalf("control characters not removed: %q", got)
	}
}

func TestLine_TabPreserved(t *testing.T) {
	raw := "key\tvalue"
	got := sanitize.Line(raw)
	if !strings.Contains(got, "\t") {
		t.Fatalf("tab should be preserved, got %q", got)
	}
}

func TestLine_Truncation(t *testing.T) {
	long := strings.Repeat("a", sanitize.MaxLineBytes+100)
	got := sanitize.Line(long)
	if !strings.HasSuffix(got, "…") {
		t.Fatalf("expected truncation suffix, got len=%d", len(got))
	}
	if len(got) > sanitize.MaxLineBytes+len("…")+10 {
		t.Fatalf("result too long after truncation: %d bytes", len(got))
	}
}

func TestFieldName_LowerCase(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Level", "level"},
		{"  TIME  ", "time"},
		{"Message", "message"},
		{"", ""},
	}
	for _, c := range cases {
		got := sanitize.FieldName(c.in)
		if got != c.want {
			t.Errorf("FieldName(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
