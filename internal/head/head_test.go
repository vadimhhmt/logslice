package head_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/head"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Raw:       `{"msg":"` + msg + `"}`,
		Fields:    map[string]any{"msg": msg},
	}
}

func TestTaker_DefaultsToUnlimited(t *testing.T) {
	tk := head.New(0)
	for i := 0; i < 1000; i++ {
		_, ok := tk.Process(makeEntry("x"))
		if !ok {
			t.Fatalf("expected all entries to pass through when max=0, failed at %d", i)
		}
	}
}

func TestTaker_StopsAfterN(t *testing.T) {
	const n = 5
	tk := head.New(n)
	passed := 0
	for i := 0; i < 10; i++ {
		_, ok := tk.Process(makeEntry("x"))
		if ok {
			passed++
		}
	}
	if passed != n {
		t.Fatalf("expected %d entries, got %d", n, passed)
	}
}

func TestTaker_DoneSignal(t *testing.T) {
	tk := head.New(3)
	if tk.Done() {
		t.Fatal("Done should be false before any entries")
	}
	for i := 0; i < 3; i++ {
		tk.Process(makeEntry("x")) //nolint:errcheck
	}
	if !tk.Done() {
		t.Fatal("Done should be true after quota reached")
	}
}

func TestTaker_Reset(t *testing.T) {
	tk := head.New(2)
	tk.Process(makeEntry("a"))
	tk.Process(makeEntry("b"))
	if !tk.Done() {
		t.Fatal("expected Done after 2 entries")
	}
	tk.Reset()
	if tk.Done() {
		t.Fatal("expected Done=false after Reset")
	}
	_, ok := tk.Process(makeEntry("c"))
	if !ok {
		t.Fatal("expected entry to pass through after Reset")
	}
}

func TestTaker_NegativeMaxUnlimited(t *testing.T) {
	tk := head.New(-1)
	for i := 0; i < 50; i++ {
		_, ok := tk.Process(makeEntry("x"))
		if !ok {
			t.Fatalf("negative max should behave as unlimited, failed at %d", i)
		}
	}
}
