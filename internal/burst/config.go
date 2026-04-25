package burst

import (
	"flag"
	"fmt"
	"time"
)

// Config holds the parsed flag values for burst detection.
type Config struct {
	Enabled   bool
	Window    time.Duration
	Threshold int
	GroupBy   string
}

// RegisterFlags registers burst-detection flags on the given FlagSet.
func RegisterFlags(fs *flag.FlagSet) *Config {
	c := &Config{}
	fs.BoolVar(&c.Enabled, "burst", false, "enable burst detection")
	fs.DurationVar(&c.Window, "burst-window", 5*time.Second, "sliding window size for burst detection")
	fs.IntVar(&c.Threshold, "burst-threshold", 10, "number of events within window that triggers a burst")
	fs.StringVar(&c.GroupBy, "burst-group-by", "", "field to group burst counts by (empty = global)")
	return c
}

// Build constructs a Detector from the config, or returns nil if disabled.
func (c *Config) Build() (*Detector, error) {
	if !c.Enabled {
		return nil, nil
	}
	if c.Threshold <= 0 {
		return nil, fmt.Errorf("burst-threshold must be greater than zero, got %d", c.Threshold)
	}
	if c.Window <= 0 {
		return nil, fmt.Errorf("burst-window must be a positive duration, got %s", c.Window)
	}
	return New(c.Window, c.Threshold, c.GroupBy), nil
}
