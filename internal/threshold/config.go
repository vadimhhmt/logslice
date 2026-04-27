package threshold

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// Config holds raw flag values for threshold checking.
type Config struct {
	Enabled bool
	Field   string
	Min     string
	Max     string
}

// RegisterFlags attaches threshold flags to fs.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.BoolVar(&cfg.Enabled, "threshold", false, "enable numeric field threshold filtering")
	fs.StringVar(&cfg.Field, "threshold-field", "", "field name to evaluate")
	fs.StringVar(&cfg.Min, "threshold-min", "", "minimum allowed value (inclusive)")
	fs.StringVar(&cfg.Max, "threshold-max", "", "maximum allowed value (inclusive)")
}

// Build constructs a Checker from the config, or returns nil when disabled.
func (c *Config) Build() (*Checker, error) {
	if !c.Enabled {
		return nil, nil
	}
	field := strings.TrimSpace(c.Field)
	if field == "" {
		return nil, fmt.Errorf("threshold: --threshold-field is required when threshold is enabled")
	}
	var minPtr, maxPtr *float64
	if c.Min != "" {
		v, err := strconv.ParseFloat(strings.TrimSpace(c.Min), 64)
		if err != nil {
			return nil, fmt.Errorf("threshold: invalid --threshold-min %q: %w", c.Min, err)
		}
		minPtr = &v
	}
	if c.Max != "" {
		v, err := strconv.ParseFloat(strings.TrimSpace(c.Max), 64)
		if err != nil {
			return nil, fmt.Errorf("threshold: invalid --threshold-max %q: %w", c.Max, err)
		}
		maxPtr = &v
	}
	return New(field, minPtr, maxPtr)
}
