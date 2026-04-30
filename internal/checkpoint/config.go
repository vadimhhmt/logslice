package checkpoint

import (
	"errors"
	"flag"
)

// Config holds CLI-parsed options for the checkpoint feature.
type Config struct {
	Enabled bool
	Path    string
	Reset   bool
}

// RegisterFlags attaches checkpoint flags to the provided FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.BoolVar(&cfg.Enabled, "checkpoint", false,
		"resume processing from the last saved position")
	fs.StringVar(&cfg.Path, "checkpoint-file", ".logslice_checkpoint",
		"path to the checkpoint state file")
	fs.BoolVar(&cfg.Reset, "checkpoint-reset", false,
		"clear any existing checkpoint before starting")
}

// Build returns a Manager configured from cfg, or nil when the feature
// is disabled. If Reset is set the existing checkpoint is cleared first.
func (cfg *Config) Build() (*Manager, error) {
	if !cfg.Enabled {
		return nil, nil
	}
	if cfg.Path == "" {
		return nil, errors.New("checkpoint: file path must not be empty")
	}
	m, err := New(cfg.Path)
	if err != nil {
		return nil, err
	}
	if cfg.Reset {
		if err := m.Reset(); err != nil && !errors.Is(err, errNotExist(err)) {
			return nil, err
		}
		// Re-initialise after reset so the manager is in a clean state.
		m, err = New(cfg.Path)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

// errNotExist is a small helper so we avoid importing os in this file.
func errNotExist(err error) error { return err }
