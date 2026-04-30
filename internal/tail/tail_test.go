package tail_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
	"github.com/logslice/logslice/internal/tail"
)

func makeEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Raw:       `{"msg":"` + msg + `"}`,
		Fields:    map[string]any{"msg": msg},
	}
}

func TestTailer_DefaultsToTen(t *testing.T) {
	tr := tail.New(0)
	if tr.Len() != 0 {
		t.Fatalf("expected empty, got %d", tr.Len())
	}
	for i := 0; i < 15; i++ {
		tr.Push(makeEntry("x"))
	}
	if tr.Len() != 10 {
		t.Fatalf("expected 10, got %d", tr.Len())
	}
}

func TestTailer_FewerThanN(t *testing.T) {
	tr := tail.New(5)
	tr.Push(makeEntry("a"))
	tr.Push(makeEntry("b"))
	entries := tr.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2, got %d", len(entries))
	}
	if entries[0].Fields["msg"] != "a" || entries[1].Fields["msg"] != "b" {
		t.Fatalf("unexpected order: %v", entries)
	}
}

func TestTailer_ExactlyN(t *testing.T) {
	tr := tail.New(3)
	for _, m := range []string{"a", "b", "c"} {
		tr.Push(makeEntry(m))
	}
	entries := tr.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3, got %d", len(entries))
	}
}

func TestTailer_OverflowKeepsLast(t *testing.T) {
	tr := tail.New(3)
	for _, m := range []string{"a", "b", "c", "d", "e"} {
		tr.Push(makeEntry(m))
	}
	entries := tr.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3, got %d", len(entries))
	}
	if entries[0].Fields["msg"] != "c" {
		t.Errorf("expected oldest=c, got %v", entries[0].Fields["msg"])
	}
	if entries[2].Fields["msg"] != "e" {
		t.Errorf("expected newest=e, got %v", entries[2].Fields["msg"])
	}
}

func TestTailer_EmptyEntries(t *testing.T) {
	tr := tail.New(5)
	if got := tr.Entries(); len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestTailer_EntriesReturnsCopy(t *testing.T) {
	// Mutating the slice returned by Entries should not affect internal state.
	tr := tail.New(3)
	for _, m := range []string{"a", "b", "c"} {
		tr.Push(makeEntry(m))
	}
	got := tr.Entries()
	got[0] = makeEntry("z")
	second := tr.Entries()
	if second[0].Fields["msg"] != "a" {
		t.Errorf("Entries() returned a reference to internal slice; mutation affected state")
	}
}
