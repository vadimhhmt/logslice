// Package config handles parsing and validation of CLI flags and
// configuration for logslice.
package config

import (
	"errors"
	"flag"
	"fmt"
	"time"
)

// Config holds all runtime configuration derived from CLI flags.
type Config struct {
	// Input
	FilePath string

	// Time range filters
	From time.Time
	To   time.Time

	// Field pattern filter (comma-separated key=pattern pairs)
	Pattern string

	// Output
	Format string // json | pretty | raw
	Fields string // comma-separated list of fields to include

	// Stats
	ShowStats bool
}

// Parse reads flags from os.Args and returns a validated Config.
func Parse() (*Config, error) {
	cfg := &Config{}

	var fromStr, toStr string

	flag.StringVar(&cfg.FilePath, "file", "", "path to log file (default: stdin)")
	flag.StringVar(&fromStr, "from", "", "start of time range (RFC3339)")
	flag.StringVar(&toStr, "to", "", "end of time range (RFC3339)")
	flag.StringVar(&cfg.Pattern, "pattern", "", "field filter pattern, e.g. level=error,service=api")
	flag.StringVar(&cfg.Format, "format", "json", "output format: json, pretty, raw")
	flag.StringVar(&cfg.Fields, "fields", "", "comma-separated fields to include in output")
	flag.BoolVar(&cfg.ShowStats, "stats", false, "print summary statistics after processing")

	flag.Parse()

	return cfg, cfg.validate(fromStr, toStr)
}

func (c *Config) validate(fromStr, toStr string) error {
	var err error

	if fromStr != "" {
		c.From, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return fmt.Errorf("invalid -from value %q: %w", fromStr, err)
		}
	}

	if toStr != "" {
		c.To, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			return fmt.Errorf("invalid -to value %q: %w", toStr, err)
		}
	}

	if !c.From.IsZero() && !c.To.IsZero() && c.To.Before(c.From) {
		return errors.New("-to must not be before -from")
	}

	switch c.Format {
	case "json", "pretty", "raw":
		// valid
	default:
		return fmt.Errorf("unknown -format %q: must be json, pretty, or raw", c.Format)
	}

	return nil
}
