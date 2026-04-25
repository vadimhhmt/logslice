package transform

import (
	"flag"
	"fmt"
	"strings"
)

// RuleFlags holds raw CLI values for transform rules.
type RuleFlags struct {
	Renames  []string // "old=new"
	ValueMap []string // "field:old=new"
}

// RegisterFlags attaches transform flags to the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, rf *RuleFlags) {
	fs.Func("rename", "rename field: old=new (repeatable)", func(s string) error {
		if !strings.Contains(s, "=") {
			return fmt.Errorf("rename: expected old=new, got %q", s)
		}
		rf.Renames = append(rf.Renames, s)
		return nil
	})
	fs.Func("remap", "remap value: field:old=new (repeatable)", func(s string) error {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 || !strings.Contains(parts[1], "=") {
			return fmt.Errorf("remap: expected field:old=new, got %q", s)
		}
		rf.ValueMap = append(rf.ValueMap, s)
		return nil
	})
}

// Build converts RuleFlags into a slice of Rules.
func (rf *RuleFlags) Build() ([]Rule, error) {
	var rules []Rule
	for _, r := range rf.Renames {
		parts := strings.SplitN(r, "=", 2)
		rules = append(rules, Rule{FromField: parts[0], ToField: parts[1]})
	}
	remaps := map[string]map[string]string{}
	for _, r := range rf.ValueMap {
		colon := strings.Index(r, ":")
		field := r[:colon]
		pair := r[colon+1:]
		kv := strings.SplitN(pair, "=", 2)
		if remaps[field] == nil {
			remaps[field] = map[string]string{}
		}
		remaps[field][strings.ToLower(kv[0])] = kv[1]
	}
	for field, vm := range remaps {
		rules = append(rules, Rule{FromField: field, ValueMap: vm})
	}
	return rules, nil
}
