package label

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Raw:       `{"msg":"test"}`,
		Fields:    fields,
	}
}

func TestLabeler_AddsField(t *testing.T) {
	l, err := New("source", "app1", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := makeEntry(map[string]interface{}{"msg": "hello"})
	out := l.Apply(e)
	if got := out.Fields["source"]; got != "app1" {
		t.Errorf("expected source=app1, got %v", got)
	}
}

func TestLabeler_NoOverwriteWhenFieldExists(t *testing.T) {
	l, _ := New("source", "app1", false)
	e := makeEntry(map[string]interface{}{"source": "original"})
	out := l.Apply(e)
	if got := out.Fields["source"]; got != "original" {
		t.Errorf("expected original, got %v", got)
	}
}

func TestLabeler_OverwriteWhenEnabled(t *testing.T) {
	l, _ := New("source", "app1", true)
	e := makeEntry(map[string]interface{}{"source": "old"})
	out := l.Apply(e)
	if got := out.Fields["source"]; got != "app1" {
		t.Errorf("expected app1, got %v", got)
	}
}

func TestLabeler_OriginalUnmodified(t *testing.T) {
	l, _ := New("env", "prod", true)
	e := makeEntry(map[string]interface{}{"msg": "hi"})
	l.Apply(e)
	if _, ok := e.Fields["env"]; ok {
		t.Error("original entry should not be modified")
	}
}

func TestLabeler_EmptyKeyReturnsError(t *testing.T) {
	_, err := New("", "val", false)
	if err == nil {
		t.Error("expected error for empty key")
	}
}

func TestLabeler_Process(t *testing.T) {
	l, _ := New("region", "eu-west", false)
	in := make(chan parser.Entry, 3)
	out := make(chan parser.Entry, 3)

	for i := 0; i < 3; i++ {
		in <- makeEntry(map[string]interface{}{"i": i})
	}
	close(in)

	l.Process(in, out)
	close(out)

	count := 0
	for e := range out {
		count++
		if e.Fields["region"] != "eu-west" {
			t.Errorf("entry %d missing region label", count)
		}
	}
	if count != 3 {
		t.Errorf("expected 3 entries, got %d", count)
	}
}
