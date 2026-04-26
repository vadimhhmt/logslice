package window_test

import (
	"testing"
	"time"

	"logslice/internal/window"
)

func makeEntry(ts time.Time, fields map[string]any) map[string]any {
	e := make(map[string]any, len(fields)+1)
	for k, v := range fields {
		e[k] = v
	}
	if !ts.IsZero() {
		e["_time"] = ts
	}
	return e
}

func TestWindow_EmitsWhenFull(t *testing.T) {
	w := window.New(3)

	base := time.Now()
	emitted := 0

	for i := 0; i < 3; i++ {
		e := makeEntry(base.Add(time.Duration(i)*time.Second), map[string]any{"msg": "hello"})
		got := w.Push(e)
		if i < 2 && got != nil {
			t.Fatalf("entry %d: expected nil before window full, got %v", i, got)
		}
		if i == 2 {
			if got == nil {
				t.Fatal("expected window to emit on third push, got nil")
			}
			emitted = len(got)
		}
	}

	if emitted != 3 {
		t.Errorf("expected 3 entries in emitted window, got %d", emitted)
	}
}

func TestWindow_SlidingDropsOldest(t *testing.T) {
	w := window.New(2)

	base := time.Now()

	e1 := makeEntry(base, map[string]any{"msg": "first"})
	e2 := makeEntry(base.Add(time.Second), map[string]any{"msg": "second"})
	e3 := makeEntry(base.Add(2*time.Second), map[string]any{"msg": "third"})

	if w.Push(e1) != nil {
		t.Fatal("unexpected emit on first push")
	}
	got := w.Push(e2)
	if got == nil {
		t.Fatal("expected emit on second push (window size 2)")
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}

	got = w.Push(e3)
	if got == nil {
		t.Fatal("expected emit on third push")
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries after slide, got %d", len(got))
	}
	if got[0]["msg"] != "second" {
		t.Errorf("expected oldest to be dropped; first entry msg = %v", got[0]["msg"])
	}
	if got[1]["msg"] != "third" {
		t.Errorf("expected newest to be third; second entry msg = %v", got[1]["msg"])
	}
}

func TestWindow_FlushReturnsBuffered(t *testing.T) {
	w := window.New(5)

	base := time.Now()
	for i := 0; i < 3; i++ {
		e := makeEntry(base.Add(time.Duration(i)*time.Second), map[string]any{"i": i})
		w.Push(e)
	}

	flushed := w.Flush()
	if len(flushed) != 3 {
		t.Errorf("expected 3 flushed entries, got %d", len(flushed))
	}

	// After flush the buffer should be empty.
	again := w.Flush()
	if len(again) != 0 {
		t.Errorf("expected empty flush after drain, got %d entries", len(again))
	}
}

func TestWindow_ZeroSizeAlwaysEmits(t *testing.T) {
	w := window.New(0)

	e := makeEntry(time.Now(), map[string]any{"msg": "immediate"})
	got := w.Push(e)
	if got == nil {
		t.Fatal("size-0 window should emit on every push")
	}
	if len(got) != 1 {
		t.Errorf("expected 1 entry, got %d", len(got))
	}
}

func TestWindow_Len(t *testing.T) {
	w := window.New(4)

	if w.Len() != 0 {
		t.Errorf("expected initial length 0, got %d", w.Len())
	}

	base := time.Now()
	w.Push(makeEntry(base, map[string]any{"x": 1}))
	if w.Len() != 1 {
		t.Errorf("expected length 1 after one push, got %d", w.Len())
	}

	w.Push(makeEntry(base.Add(time.Second), map[string]any{"x": 2}))
	if w.Len() != 2 {
		t.Errorf("expected length 2 after two pushes, got %d", w.Len())
	}
}
