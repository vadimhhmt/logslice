// Package checkpoint tracks the last successfully processed log entry
// so that a subsequent run can resume from where it left off.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

// State holds the persisted position within a log stream.
type State struct {
	LastTimestamp time.Time `json:"last_timestamp"`
	LinesProcessed int64     `json:"lines_processed"`
}

// Manager reads and writes checkpoint state to a file.
type Manager struct {
	mu   sync.Mutex
	path string
	state State
}

// New returns a Manager backed by the given file path.
// If the file exists its contents are loaded immediately.
func New(path string) (*Manager, error) {
	m := &Manager{path: path}
	if err := m.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return m, nil
}

// Record updates the in-memory state with the latest timestamp and
// increments the processed-line counter, then flushes to disk.
func (m *Manager) Record(ts time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if ts.After(m.state.LastTimestamp) {
		m.state.LastTimestamp = ts
	}
	m.state.LinesProcessed++
	return m.flush()
}

// State returns a copy of the current checkpoint state.
func (m *Manager) State() State {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.state
}

// Reset clears the persisted state and removes the checkpoint file.
func (m *Manager) Reset() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = State{}
	return os.Remove(m.path)
}

func (m *Manager) load() error {
	f, err := os.Open(m.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&m.state)
}

func (m *Manager) flush() error {
	f, err := os.CreateTemp("", "checkpoint-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	if err := json.NewEncoder(f).Encode(m.state); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	f.Close()
	return os.Rename(tmp, m.path)
}
