package redact

import (
	"flag"
	"strings"
)

// Config holds redaction settings parsed from CLI flags.
type Config struct {
	// ExtraKeys is the list of additional field names to redact.
	ExtraKeys []string
	// ValuePatterns is the list of regexp patterns whose matching values are redacted.
	ValuePatterns []string
	// Enabled indicates whether redaction is active at all.
	Enabled bool
}

// RegisterFlags binds redaction CLI flags to the provided FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	var keys, patterns string
	fs.BoolVar(&cfg.Enabled, "redact", false, "enable field redaction for sensitive values")
	fs.StringVar(&keys, "redact-keys", "", "comma-separated extra field names to redact")
	fs.StringVar(&patterns, "redact-patterns", "", "comma-separated regexp patterns; matching values are redacted")
	// Post-parse hooks are not possible with flag.FlagSet directly, so callers
	// must call Finalise after parsing.
	cfg.rawKeys = &keys
	cfg.rawPatterns = &patterns
}

// Finalise splits the raw comma-separated strings into slices.
// Must be called after flag.FlagSet.Parse.
func (c *Config) Finalise() {
	if c.rawKeys != nil && *c.rawKeys != "" {
		for _, k := range strings.Split(*c.rawKeys, ",") {
			if t := strings.TrimSpace(k); t != "" {
				c.ExtraKeys = append(c.ExtraKeys, t)
			}
		}
	}
	if c.rawPatterns != nil && *c.rawPatterns != "" {
		for _, p := range strings.Split(*c.rawPatterns, ",") {
			if t := strings.TrimSpace(p); t != "" {
				c.ValuePatterns = append(c.ValuePatterns, t)
			}
		}
	}
}

// BuildRedactor constructs a Redactor from the Config.
// Returns nil if redaction is not enabled.
func (c *Config) BuildRedactor() (*Redactor, error) {
	if !c.Enabled {
		return nil, nil
	}
	r := New()
	for _, k := range c.ExtraKeys {
		r.AddKey(k)
	}
	for _, p := range c.ValuePatterns {
		if err := r.AddPattern(p); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// unexported fields used to bridge flag parsing.
type configInternal = Config

func init() {
	// Ensure configInternal alias is used; no-op.
	_ = configInternal{}
}
