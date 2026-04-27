package merge_test

import (
	"testing"
	"time"

	"github.com/your-org/logslice/internal/merge"
	"github.com/your-org/logslice/internal/parser"
)

func feed(entries []parser.Entry) <-chan parser.Entry {
	ch := make(chan parser.Entry, len(entries))
	for _, e := range entries {
		ch <- e
	}
	close(ch)
	return ch
}

func makeEntry(ts time.Time, msg string) parser.Entry {
	return parser.Entry{
		Time: ts,
		Fields: map[string]any{"msg": msg},
		Raw:    `{"msg":"` + msg + `"}`,
	}
}

func collect(ch <-chan parser.Entry) []parser.Entry {
	var out []parser.Entry
	for e := range ch {
		out = append(out, e)
	}
	return out
}

func TestMerge_TwoSortedStreams(t *testing.T) {
	t0 := time.Unix(1000, 0)
	s1 := feed([]parser.Entry{
		makeEntry(t0.Add(0), "a"),
		makeEntry(t0.Add(2*time.Second), "c"),
		makeEntry(t0.Add(4*time.Second), "e"),
	})
	s2 := feed([]parser.Entry{
		makeEntry(t0.Add(1*time.Second), "b"),
		makeEntry(t0.Add(3*time.Second), "d"),
	})

	m := merge.New([]<-chan parser.Entry{s1, s2})
	result := collect(m.Merge())

	if len(result) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(result))
	}
	expected := []string{"a", "b", "c", "d", "e"}
	for i, e := range result {
		if e.Fields["msg"] != expected[i] {
			t.Errorf("position %d: got %q, want %q", i, e.Fields["msg"], expected[i])
		}
	}
}

func TestMerge_SingleSource(t *testing.T) {
	t0 := time.Unix(2000, 0)
	s1 := feed([]parser.Entry{
		makeEntry(t0, "x"),
		makeEntry(t0.Add(time.Second), "y"),
	})

	m := merge.New([]<-chan parser.Entry{s1})
	result := collect(m.Merge())

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestMerge_EmptySources(t *testing.T) {
	s1 := feed(nil)
	s2 := feed(nil)

	m := merge.New([]<-chan parser.Entry{s1, s2})
	result := collect(m.Merge())

	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}

func TestMerge_OneEmptyOneNot(t *testing.T) {
	t0 := time.Unix(3000, 0)
	s1 := feed(nil)
	s2 := feed([]parser.Entry{
		makeEntry(t0, "only"),
	})

	m := merge.New([]<-chan parser.Entry{s1, s2})
	result := collect(m.Merge())

	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Fields["msg"] != "only" {
		t.Errorf("unexpected entry: %v", result[0].Fields)
	}
}
