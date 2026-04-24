package pipeline_test

import (
	"bytes"
	"strings"
	"testing"

	"logslice/internal/filter"
	"logslice/internal/output"
	"logslice/internal/pipeline"
	"logslice/internal/reader"
	"logslice/internal/stats"
)

const sampleLogs = `{"ts":"2024-01-01T10:00:00Z","level":"info","msg":"startup"}
{"ts":"2024-01-01T11:00:00Z","level":"error","msg":"crash"}
{"ts":"2024-01-01T12:00:00Z","level":"info","msg":"shutdown"}
`

func makeConfig(t *testing.T, logs string) pipeline.Config {
	t.Helper()
	r := reader.New(strings.NewReader(logs))
	f, err := filter.New(filter.Options{})
	if err != nil {
		t.Fatalf("filter.New: %v", err)
	}
	fmt := output.New(output.Options{Format: "json"})
	col := stats.New()
	return pipeline.Config{Reader: r, Filter: f, Formatter: fmt, Collector: col}
}

func TestPipeline_AllMatch(t *testing.T) {
	cfg := makeConfig(t, sampleLogs)
	var buf bytes.Buffer
	res, err := pipeline.Run(cfg, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Read != 3 || res.Matched != 3 || res.Dropped != 0 {
		t.Errorf("got %+v, want Read=3 Matched=3 Dropped=0", res)
	}
}

func TestPipeline_InvalidLinesDropped(t *testing.T) {
	logs := "not json\n" + sampleLogs
	cfg := makeConfig(t, logs)
	var buf bytes.Buffer
	res, err := pipeline.Run(cfg, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Dropped != 1 {
		t.Errorf("expected 1 dropped, got %d", res.Dropped)
	}
}

func TestPipeline_EmptyInput(t *testing.T) {
	cfg := makeConfig(t, "")
	var buf bytes.Buffer
	res, err := pipeline.Run(cfg, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Read != 0 {
		t.Errorf("expected 0 read, got %d", res.Read)
	}
}

func TestPipeline_NilCollector(t *testing.T) {
	r := reader.New(strings.NewReader(sampleLogs))
	f, _ := filter.New(filter.Options{})
	fmt := output.New(output.Options{Format: "json"})
	cfg := pipeline.Config{Reader: r, Filter: f, Formatter: fmt, Collector: nil}
	var buf bytes.Buffer
	_, err := pipeline.Run(cfg, &buf)
	if err != nil {
		t.Fatalf("unexpected error with nil collector: %v", err)
	}
}
