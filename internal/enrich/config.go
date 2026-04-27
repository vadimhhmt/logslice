package enrich

import (
	"flag"
	"fmt"
	"strings"
)

// ruleSpec holds the raw flag value for a single enrichment rule.
// Format: "from_field:to_field:fn" where fn is one of upper, lower.
type ruleSpec struct {
	specs []string
}

func (r *ruleSpec) String() string { return strings.Join(r.specs, ",") }
func (r *ruleSpec) Set(v string) error {
	r.specs = append(r.specs, v)
	return nil
}

// RegisterFlags adds enrich-related flags to the given FlagSet and returns a
// builder func that constructs an *Enricher from the parsed flags.
func RegisterFlags(fs *flag.FlagSet) func() (*Enricher, error) {
	specs := &ruleSpec{}
	fs.Var(specs, "enrich", "enrichment rule: from_field:to_field:fn (fn: upper|lower); repeatable")

	return func() (*Enricher, error) {
		if len(specs.specs) == 0 {
			return nil, nil
		}
		rules, err := parseSpecs(specs.specs)
		if err != nil {
			return nil, err
		}
		return New(rules), nil
	}
}

func parseSpecs(specs []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(specs))
	for _, s := range specs {
		parts := strings.SplitN(s, ":", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("enrich: invalid rule %q (want from:to:fn)", s)
		}
		from, to, fnName := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2])
		if from == "" || to == "" {
			return nil, fmt.Errorf("enrich: from and to fields must not be empty in %q", s)
		}
		var fn func(string) (string, bool)
		switch strings.ToLower(fnName) {
		case "upper":
			fn = UpperCase
		case "lower":
			fn = LowerCase
		default:
			return nil, fmt.Errorf("enrich: unknown fn %q in rule %q", fnName, s)
		}
		rules = append(rules, Rule{From: from, To: to, Fn: fn})
	}
	return rules, nil
}
