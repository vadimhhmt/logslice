package dedupe_test

import (
	"testing"
	"time"

	"logslice/internal/dedupe"
	"logslice/internal/parser"
)

func makeEntry(msg, level string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Raw:       `{"msg":"` + msg + `"}`,
		Fields:    map[string]any{"msg": msg, "level": level},
	}
}

func TestDeduplicator_UniqueEntriesPassThrough(t *testing.T) {
	d := dedupe.New([]string{"msg"}, false)
	e1, ok1 := d.Process(makeEntry("hello", "info"))
	e2, ok2 := d.Process(makeEntry("world", "info"))

	if !ok1 || !ok2 {
		t.Fatal("expected both unique entries to pass through")
	}
	if e1.Fields["msg"] != "hello" || e2.Fields["msg"] != "world" {
		t.Error("unexpected field values in passed-through entries")
	}
}

func TestDeduplicator_ConsecutiveDuplicatesDropped(t *testing.T) {
	d := dedupe.New([]string{"msg"}, false)
	d.Process(makeEntry("dup", "warn")) // first — passes
	_, ok2 := d.Process(makeEntry("dup", "warn"))
	_, ok3 := d.Process(makeEntry("dup", "warn"))

	if ok2 || ok3 {
		t.Error("expected consecutive duplicates to be dropped")
	}
}

func TestDeduplicator_NonConsecutiveDuplicatesAllowed(t *testing.T) {
	d := dedupe.New([]string{"msg"}, false)
	d.Process(makeEntry("a", "info"))
	d.Process(makeEntry("b", "info"))
	_, ok := d.Process(makeEntry("a", "info")) // not consecutive

	if !ok {
		t.Error("expected non-consecutive duplicate to pass through")
	}
}

func TestDeduplicator_InjectCount(t *testing.T) {
	d := dedupe.New([]string{"msg"}, true)
	d.Process(makeEntry("x", "info")) // first
	d.Process(makeEntry("x", "info")) // dup 1
	d.Process(makeEntry("x", "info")) // dup 2
	out, ok := d.Process(makeEntry("y", "info")) // new — should carry count=2

	if !ok {
		t.Fatal("expected new entry to pass through")
	}
	count, exists := out.Fields["_suppressed"]
	if !exists {
		t.Fatal("expected _suppressed field to be injected")
	}
	if count != 2 {
		t.Errorf("expected _suppressed=2, got %v", count)
	}
}

func TestDeduplicator_FlushReturnsPending(t *testing.T) {
	d := dedupe.New(nil, false)
	d.Process(makeEntry("z", "error"))
	d.Process(makeEntry("z", "error"))
	d.Process(makeEntry("z", "error"))

	n, ok := d.Flush()
	if !ok {
		t.Fatal("expected Flush to signal pending suppressed entries")
	}
	if n != 2 {
		t.Errorf("expected 2 suppressed, got %d", n)
	}
}

func TestDeduplicator_FlushNopWhenEmpty(t *testing.T) {
	d := dedupe.New(nil, false)
	_, ok := d.Flush()
	if ok {
		t.Error("expected Flush to return false when nothing is pending")
	}
}

func TestDeduplicator_AllFieldsFingerprint(t *testing.T) {
	d := dedupe.New(nil, false) // nil → use all fields
	d.Process(makeEntry("same", "info"))
	_, ok := d.Process(makeEntry("same", "info"))
	if ok {
		t.Error("expected duplicate (all fields) to be dropped")
	}
}
