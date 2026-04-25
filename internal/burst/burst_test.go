package burst_test

import (
	"testing"
	"time"

	"logslice/internal/burst"
)

func makeEntry(ts time.Time, fields map[string]any) map[string]any {
	entry := map[string]any{
		"timestamp": ts,
	}
	for k, v := range fields {
		entry[k] = v
	}
	return entry
}

func TestDetector_NoBurstWhenSpread(t *testing.T) {
	det := burst.New(3, 5*time.Second, "level")

	base := time.Now()
	entries := []map[string]any{
		makeEntry(base, map[string]any{"level": "error"}),
		makeEntry(base.Add(2*time.Second), map[string]any{"level": "error"}),
		makeEntry(base.Add(4*time.Second), map[string]any{"level": "error"}),
	}

	var bursts []map[string]any
	for _, e := range entries {
		if b := det.Add(e); b != nil {
			bursts = append(bursts, b)
		}
	}

	if len(bursts) != 0 {
		t.Errorf("expected no burst events, got %d", len(bursts))
	}
}

func TestDetector_BurstDetectedWhenThresholdExceeded(t *testing.T) {
	det := burst.New(3, 5*time.Second, "level")

	base := time.Now()
	entries := []map[string]any{
		makeEntry(base, map[string]any{"level": "error"}),
		makeEntry(base.Add(500*time.Millisecond), map[string]any{"level": "error"}),
		makeEntry(base.Add(1*time.Second), map[string]any{"level": "error"}),
		makeEntry(base.Add(1500*time.Millisecond), map[string]any{"level": "error"}),
	}

	var bursts []map[string]any
	for _, e := range entries {
		if b := det.Add(e); b != nil {
			bursts = append(bursts, b)
		}
	}

	if len(bursts) == 0 {
		t.Fatal("expected at least one burst event")
	}

	b := bursts[0]
	if _, ok := b["burst_count"]; !ok {
		t.Error("burst event missing 'burst_count' field")
	}
	if _, ok := b["burst_window"]; !ok {
		t.Error("burst event missing 'burst_window' field")
	}
	if _, ok := b["burst_key"]; !ok {
		t.Error("burst event missing 'burst_key' field")
	}
}

func TestDetector_BurstResetsAfterWindow(t *testing.T) {
	det := burst.New(2, 1*time.Second, "level")

	base := time.Now()

	// First burst window
	det.Add(makeEntry(base, map[string]any{"level": "error"}))
	det.Add(makeEntry(base.Add(200*time.Millisecond), map[string]any{"level": "error"}))
	det.Add(makeEntry(base.Add(400*time.Millisecond), map[string]any{"level": "error"}))

	// Second window — should be treated independently
	secondBase := base.Add(2 * time.Second)
	var secondBursts []map[string]any
	for _, e := range []map[string]any{
		makeEntry(secondBase, map[string]any{"level": "error"}),
		makeEntry(secondBase.Add(100*time.Millisecond), map[string]any{"level": "error"}),
		makeEntry(secondBase.Add(200*time.Millisecond), map[string]any{"level": "error"}),
	} {
		if b := det.Add(e); b != nil {
			secondBursts = append(secondBursts, b)
		}
	}

	if len(secondBursts) == 0 {
		t.Error("expected burst detection to reset and fire again in second window")
	}
}

func TestDetector_DifferentKeyValuesTrackedSeparately(t *testing.T) {
	det := burst.New(2, 5*time.Second, "level")

	base := time.Now()
	results := map[string]int{}

	entries := []map[string]any{
		makeEntry(base, map[string]any{"level": "error"}),
		makeEntry(base.Add(100*time.Millisecond), map[string]any{"level": "error"}),
		makeEntry(base.Add(200*time.Millisecond), map[string]any{"level": "error"}),
		makeEntry(base.Add(300*time.Millisecond), map[string]any{"level": "warn"}),
		makeEntry(base.Add(400*time.Millisecond), map[string]any{"level": "warn"}),
		makeEntry(base.Add(500*time.Millisecond), map[string]any{"level": "warn"}),
	}

	for _, e := range entries {
		if b := det.Add(e); b != nil {
			key, _ := b["burst_key"].(string)
			results[key]++
		}
	}

	if results["error"] == 0 {
		t.Error("expected burst event for 'error' level")
	}
	if results["warn"] == 0 {
		t.Error("expected burst event for 'warn' level")
	}
}

func TestDetector_MissingKeyFieldIgnored(t *testing.T) {
	det := burst.New(2, 5*time.Second, "level")

	base := time.Now()
	for i := 0; i < 5; i++ {
		e := makeEntry(base.Add(time.Duration(i)*100*time.Millisecond), map[string]any{"msg": "no level here"})
		if b := det.Add(e); b != nil {
			t.Errorf("unexpected burst event for entry without key field: %v", b)
		}
	}
}
