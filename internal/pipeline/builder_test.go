package pipeline_test

import (
	"os"
	"path/filepath"
	"testing"

	"logslice/internal/config"
	"logslice/internal/pipeline"
)

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.log")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp log: %v", err)
	}
	return p
}

func TestBuild_FromFile(t *testing.T) {
	p := writeTempLog(t, sampleLogs)
	cfg := &config.Config{File: p, Format: "json"}
	pipe, cleanup, err := pipeline.Build(cfg)
	defer cleanup()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if pipe.Reader == nil || pipe.Filter == nil || pipe.Formatter == nil || pipe.Collector == nil {
		t.Error("expected all pipeline components to be non-nil")
	}
}

func TestBuild_MissingFile(t *testing.T) {
	cfg := &config.Config{File: "/no/such/file.log", Format: "json"}
	_, cleanup, err := pipeline.Build(cfg)
	defer cleanup()
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestBuild_InvalidPattern(t *testing.T) {
	p := writeTempLog(t, sampleLogs)
	cfg := &config.Config{
		File:     p,
		Format:   "json",
		Patterns: map[string]string{"level": "[invalid"},
	}
	_, cleanup, err := pipeline.Build(cfg)
	defer cleanup()
	if err == nil {
		t.Fatal("expected error for invalid regex pattern, got nil")
	}
}
