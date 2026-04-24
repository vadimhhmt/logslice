package stats_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/stats"
)

func makeEntry(ts time.Time, level string) parser.Entry {
	fields := map[string]interface{}{"msg": "hello"}
	if level != "" {
		fields["level"] = level
	}
	return parser.Entry{Timestamp: ts, Fields: fields}
}

func TestCollector_Counts(t *testing.T) {
	c := stats.New()
	now := time.Now()

	c.Record(makeEntry(now, "info"), true)
	c.Record(makeEntry(now.Add(time.Second), "warn"), true)
	c.Record(makeEntry(now.Add(2*time.Second), "info"), false)

	if c.Total != 3 {
		t.Errorf("expected Total=3, got %d", c.Total)
	}
	if c.Matched != 2 {
		t.Errorf("expected Matched=2, got %d", c.Matched)
	}
	if c.Skipped != 1 {
		t.Errorf("expected Skipped=1, got %d", c.Skipped)
	}
}

func TestCollector_TimeRange(t *testing.T) {
	c := stats.New()
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	c.Record(makeEntry(base.Add(5*time.Minute), "info"), true)
	c.Record(makeEntry(base, "info"), true)
	c.Record(makeEntry(base.Add(10*time.Minute), "info"), true)

	if !c.Earliest.Equal(base) {
		t.Errorf("expected Earliest=%v, got %v", base, c.Earliest)
	}
	if !c.Latest.Equal(base.Add(10 * time.Minute)) {
		t.Errorf("expected Latest=%v, got %v", base.Add(10*time.Minute), c.Latest)
	}
}

func TestCollector_LevelCounts(t *testing.T) {
	c := stats.New()
	now := time.Now()

	for i := 0; i < 3; i++ {
		c.Record(makeEntry(now, "info"), true)
	}
	c.Record(makeEntry(now, "error"), true)

	if c.Levels["info"] != 3 {
		t.Errorf("expected info=3, got %d", c.Levels["info"])
	}
	if c.Levels["error"] != 1 {
		t.Errorf("expected error=1, got %d", c.Levels["error"])
	}
}

func TestCollector_Print(t *testing.T) {
	c := stats.New()
	now := time.Now()
	c.Record(makeEntry(now, "info"), true)
	c.Record(makeEntry(now.Add(time.Second), "warn"), false)

	var buf bytes.Buffer
	c.Print(&buf)
	out := buf.String()

	for _, want := range []string{"total", "matched", "skipped", "earliest", "latest"} {
		if !strings.Contains(out, want) {
			t.Errorf("Print output missing %q", want)
		}
	}
}
