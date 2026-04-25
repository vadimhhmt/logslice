package highlight_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/user/logslice/internal/highlight"
)

func TestHighlighter_DisabledReturnsUnchanged(t *testing.T) {
	re := regexp.MustCompile(`error`)
	h := highlight.New(highlight.Red, false, re)
	input := "level=error msg=something failed"
	if got := h.Apply(input); got != input {
		t.Errorf("expected unchanged string, got %q", got)
	}
}

func TestHighlighter_NoPatternsReturnsUnchanged(t *testing.T) {
	h := highlight.New(highlight.Cyan, true)
	input := "level=info msg=hello"
	if got := h.Apply(input); got != input {
		t.Errorf("expected unchanged string, got %q", got)
	}
}

func TestHighlighter_MatchWrappedInColour(t *testing.T) {
	re := regexp.MustCompile(`error`)
	h := highlight.New(highlight.Red, true, re)
	got := h.Apply("level=error msg=fatal error occurred")
	if !strings.Contains(got, highlight.Red+"error"+highlight.Reset) {
		t.Errorf("expected ANSI colour wrap, got %q", got)
	}
}

func TestHighlighter_MultiplePatterns(t *testing.T) {
	r1 := regexp.MustCompile(`warn`)
	r2 := regexp.MustCompile(`timeout`)
	h := highlight.New(highlight.Yellow, true, r1, r2)
	got := h.Apply("warn: connection timeout")
	if !strings.Contains(got, highlight.Yellow+"warn"+highlight.Reset) {
		t.Errorf("expected warn highlighted, got %q", got)
	}
	if !strings.Contains(got, highlight.Yellow+"timeout"+highlight.Reset) {
		t.Errorf("expected timeout highlighted, got %q", got)
	}
}

func TestHighlighter_ApplyToFields_MatchedField(t *testing.T) {
	re := regexp.MustCompile(`ERROR`)
	h := highlight.New(highlight.Red, true, re)
	line := "time=2024-01-01 level=ERROR msg=boom"
	got := h.ApplyToFields(line, []string{"level"})
	if !strings.Contains(got, highlight.Red+"ERROR"+highlight.Reset) {
		t.Errorf("expected field value highlighted, got %q", got)
	}
}

func TestHighlighter_ApplyToFields_FieldAbsent(t *testing.T) {
	re := regexp.MustCompile(`ERROR`)
	h := highlight.New(highlight.Red, true, re)
	line := "time=2024-01-01 level=ERROR msg=boom"
	got := h.ApplyToFields(line, []string{"missing"})
	if got != line {
		t.Errorf("expected unchanged line when field absent, got %q", got)
	}
}

func TestHighlighter_ApplyToFields_DisabledNoChange(t *testing.T) {
	re := regexp.MustCompile(`ERROR`)
	h := highlight.New(highlight.Red, false, re)
	line := "level=ERROR msg=boom"
	got := h.ApplyToFields(line, []string{"level"})
	if got != line {
		t.Errorf("expected unchanged when disabled, got %q", got)
	}
}
