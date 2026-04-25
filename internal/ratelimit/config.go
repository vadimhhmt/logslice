package ratelimit

import (
	"flag"
	"fmt"
	"time"
)

// Config holds the parsed CLI flags for the rate-limiter.
type Config struct {
	Enabled      bool
	MaxPerWindow int
	Window       time.Duration
}

// RegisterFlags registers rate-limit flags on the provided FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.BoolVar(&cfg.Enabled, "ratelimit", false,
		"enable rate limiting of log entries per time window")
	fs.IntVar(&cfg.MaxPerWindow, "ratelimit-max", 1000,
		"maximum number of entries allowed per window (requires --ratelimit)")
	fs.DurationVar(&cfg.Window, "ratelimit-window", time.Minute,
		"duration of the rate-limit window (requires --ratelimit)")
}

// Build constructs a Limiter from the config, or returns nil if rate
// limiting is disabled. Returns an error if the parameters are invalid.
func (c *Config) Build() (*Limiter, error) {
	if !c.Enabled {
		return nil, nil
	}
	l, err := New(c.MaxPerWindow, c.Window)
	if err != nil {
		return nil, fmt.Errorf("ratelimit config: %w", err)
	}
	return l, nil
}
