package pipeline

import (
	"fmt"
	"os"

	"logslice/internal/config"
	"logslice/internal/filter"
	"logslice/internal/output"
	"logslice/internal/reader"
	"logslice/internal/stats"
)

// Build constructs a ready-to-run pipeline.Config from the application
// config. It opens the input file when cfg.File is non-empty, otherwise
// it reads from os.Stdin.
func Build(cfg *config.Config) (Config, func(), error) {
	var r *reader.Reader
	cleanup := func() {}

	if cfg.File != "" {
		f, err := os.Open(cfg.File)
		if err != nil {
			return Config{}, cleanup, fmt.Errorf("open %q: %w", cfg.File, err)
		}
		cleanup = func() { f.Close() }
		r = reader.New(f)
	} else {
		r = reader.New(os.Stdin)
	}

	fopts := filter.Options{
		From:     cfg.From,
		To:       cfg.To,
		Patterns: cfg.Patterns,
	}
	f, err := filter.New(fopts)
	if err != nil {
		cleanup()
		return Config{}, func() {}, fmt.Errorf("build filter: %w", err)
	}

	fopts2 := output.Options{
		Format: cfg.Format,
		Fields: cfg.Fields,
	}
	fmt := output.New(fopts2)
	col := stats.New()

	return Config{
		Reader:    r,
		Filter:    f,
		Formatter: fmt,
		Collector: col,
	}, cleanup, nil
}
