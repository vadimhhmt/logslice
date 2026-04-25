package aggregate_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/aggregate"
	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Raw:       "{}",
		Fields:    fields,
	}
}

func TestAggregator_GroupsByField(t *testing.T) {
	a := aggregate.New("level")
	a.Add(makeEntry(map[string]any{"level": "info"}))
	a.Add(makeEntry(map[string]any{"level": "info"}))
	a.Add(makeEntry(map[string]any{"level": "error"}))

	results := a.Results()
	if len(results) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(results))
	}
	if results[0].Value != "info" || results[0].Count != 2 {
		t.Errorf("expected info=2, got %s=%d", results[0].Value, results[0].Count)
	}
	if results[1].Value != "error" || results[1].Count != 1 {
		t.Errorf("expected error=1, got %s=%d", results[1].Value, results[1].Count)
	}
}

func TestAggregator_MissingField(t *testing.T) {
	a := aggregate.New("level")
	a.Add(makeEntry(map[string]any{"msg": "hello"}))

	results := a.Results()
	if len(results) != 1 || results[0].Value != "<missing>" {
		t.Errorf("expected <missing> bucket, got %+v", results)
	}
}

func TestAggregator_Total(t *testing.T) {
	a := aggregate.New("level")
	for i := 0; i < 5; i++ {
		a.Add(makeEntry(map[string]any{"level": "debug"}))
	}
	if a.Total() != 5 {
		t.Errorf("expected total 5, got %d", a.Total())
	}
}

func TestAggregator_SortedByCountDesc(t *testing.T) {
	a := aggregate.New("status")
	for i := 0; i < 3; i++ {
		a.Add(makeEntry(map[string]any{"status": "200"}))
	}
	a.Add(makeEntry(map[string]any{"status": "500"}))
	a.Add(makeEntry(map[string]any{"status": "500"}))
	a.Add(makeEntry(map[string]any{"status": "404"}))

	results := a.Results()
	if results[0].Value != "200" {
		t.Errorf("expected 200 first, got %s", results[0].Value)
	}
}

func TestAggregator_Print(t *testing.T) {
	a := aggregate.New("level")
	a.Add(makeEntry(map[string]any{"level": "warn"}))

	var buf bytes.Buffer
	a.Print(&buf)
	out := buf.String()

	if !strings.Contains(out, "level") {
		t.Error("expected field name in output")
	}
	if !strings.Contains(out, "warn") {
		t.Error("expected value 'warn' in output")
	}
	if !strings.Contains(out, "1") {
		t.Error("expected count in output")
	}
}
