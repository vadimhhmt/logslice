package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/logslice/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestCheckpoint_NewMissingFileIsOK(t *testing.T) {
	_, err := checkpoint.New(tempPath(t))
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
}

func TestCheckpoint_RecordUpdatesTimestamp(t *testing.T) {
	m, _ := checkpoint.New(tempPath(t))
	ts := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	if err := m.Record(ts); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if got := m.State().LastTimestamp; !got.Equal(ts) {
		t.Errorf("LastTimestamp = %v, want %v", got, ts)
	}
}

func TestCheckpoint_RecordIncrementsCount(t *testing.T) {
	m, _ := checkpoint.New(tempPath(t))
	ts := time.Now()
	for i := 0; i < 5; i++ {
		_ = m.Record(ts)
	}
	if got := m.State().LinesProcessed; got != 5 {
		t.Errorf("LinesProcessed = %d, want 5", got)
	}
}

func TestCheckpoint_PersistedAndReloaded(t *testing.T) {
	path := tempPath(t)
	m1, _ := checkpoint.New(path)
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	_ = m1.Record(ts)
	_ = m1.Record(ts)

	m2, err := checkpoint.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if got := m2.State().LinesProcessed; got != 2 {
		t.Errorf("reloaded LinesProcessed = %d, want 2", got)
	}
	if got := m2.State().LastTimestamp; !got.Equal(ts) {
		t.Errorf("reloaded LastTimestamp = %v, want %v", got, ts)
	}
}

func TestCheckpoint_OnlyAdvancesTimestamp(t *testing.T) {
	m, _ := checkpoint.New(tempPath(t))
	newer := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	older := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_ = m.Record(newer)
	_ = m.Record(older)
	if got := m.State().LastTimestamp; !got.Equal(newer) {
		t.Errorf("LastTimestamp regressed to %v", got)
	}
}

func TestCheckpoint_Reset(t *testing.T) {
	path := tempPath(t)
	m, _ := checkpoint.New(path)
	_ = m.Record(time.Now())
	if err := m.Reset(); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected checkpoint file to be removed after Reset")
	}
	if got := m.State().LinesProcessed; got != 0 {
		t.Errorf("LinesProcessed after Reset = %d, want 0", got)
	}
}
