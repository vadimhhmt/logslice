package offset_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/offset"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Raw:       `{"msg":"` + msg + `"}`,
		Fields:    map[string]any{"msg": msg},
	}
}

func TestSkipper_ZeroSkipPassesAll(t *testing.T) {
	s := offset.New(0)
	for i := 0; i < 5; i++ {
		_, ok := s.Process(makeEntry("x"))
		if !ok {
			t.Fatalf("entry %d: expected pass-through, got dropped", i)
		}
	}
	if s.Dropped() != 0 {
		t.Fatalf("expected 0 dropped, got %d", s.Dropped())
	}
}

func TestSkipper_SkipsFirstN(t *testing.T) {
	s := offset.New(3)
	passed := 0
	for i := 0; i < 6; i++ {
		_, ok := s.Process(makeEntry("x"))
		if ok {
			passed++
		}
	}
	if passed != 3 {
		t.Fatalf("expected 3 passed, got %d", passed)
	}
	if s.Dropped() != 3 {
		t.Fatalf("expected 3 dropped, got %d", s.Dropped())
	}
}

func TestSkipper_EntryContentPreserved(t *testing.T) {
	s := offset.New(1)
	s.Process(makeEntry("skip"))
	e := makeEntry("keep")
	out, ok := s.Process(e)
	if !ok {
		t.Fatal("expected entry to pass through")
	}
	if out.Fields["msg"] != "keep" {
		t.Fatalf("expected msg=keep, got %v", out.Fields["msg"])
	}
}

func TestSkipper_NegativeSkipTreatedAsZero(t *testing.T) {
	s := offset.New(-5)
	_, ok := s.Process(makeEntry("x"))
	if !ok {
		t.Fatal("negative skip should pass all entries")
	}
}

func TestSkipper_Reset(t *testing.T) {
	s := offset.New(2)
	s.Process(makeEntry("a"))
	s.Process(makeEntry("b"))
	s.Reset()
	if s.Dropped() != 0 {
		t.Fatalf("after reset expected 0 dropped, got %d", s.Dropped())
	}
	_, ok := s.Process(makeEntry("c"))
	if ok {
		t.Fatal("after reset first entry should be dropped again")
	}
}
