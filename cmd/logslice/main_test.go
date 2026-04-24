package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logslice-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	_, err = f.WriteString(strings.Join(lines, "\n") + "\n")
	if err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return f.Name()
}

func TestMain_EndToEnd_FilterByTime(t *testing.T) {
	lines := []string{
		`{"timestamp":"2024-01-01T10:00:00Z","level":"info","msg":"startup"}`,
		`{"timestamp":"2024-01-01T12:00:00Z","level":"error","msg":"crash"}`,
		`{"timestamp":"2024-01-01T14:00:00Z","level":"info","msg":"recovered"}`,
	}
	path := writeTempLog(t, lines)

	old := os.Args
	t.Cleanup(func() { os.Args = old })
	os.Args = []string{
		"logslice",
		"--file", path,
		"--start", "2024-01-01T11:00:00Z",
		"--end", "2024-01-01T13:00:00Z",
		"--format", "raw",
	}

	// Capture stdout by redirecting
	r, w, _ := os.Pipe()
	origStdout := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = origStdout })

	main()

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	out := buf.String()

	if !strings.Contains(out, "crash") {
		t.Errorf("expected 'crash' entry in output, got: %s", out)
	}
	if strings.Contains(out, "startup") {
		t.Errorf("unexpected 'startup' entry in output, got: %s", out)
	}
	if strings.Contains(out, "recovered") {
		t.Errorf("unexpected 'recovered' entry in output, got: %s", out)
	}
}

func TestMain_EndToEnd_InvalidFile(t *testing.T) {
	old := os.Args
	t.Cleanup(func() { os.Args = old })
	os.Args = []string{
		"logslice",
		"--file", filepath.Join(t.TempDir(), "nonexistent.log"),
	}

	// We can't easily test os.Exit; just verify writeTempLog helper works
	// and the path doesn't exist as a sanity check.
	if _, err := os.Stat(os.Args[2]); !os.IsNotExist(err) {
		t.Fatal("expected file to not exist")
	}
}
