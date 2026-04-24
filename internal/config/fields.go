package config

import "strings"

// FieldList parses a comma-separated fields string into a deduplicated slice.
// An empty string returns nil, indicating all fields should be included.
func FieldList(raw string) []string {
	if raw == "" {
		return nil
	}

	seen := make(map[string]struct{})
	var result []string

	for _, part := range strings.Split(raw, ",") {
		f := strings.TrimSpace(part)
		if f == "" {
			continue
		}
		if _, ok := seen[f]; !ok {
			seen[f] = struct{}{}
			result = append(result, f)
		}
	}

	if len(result) == 0 {
		return nil
	}
	return result
}

// PatternPairs parses a comma-separated "key=value" pattern string into a map.
// Values are treated as regular-expression patterns by the filter package.
func PatternPairs(raw string) map[string]string {
	result := make(map[string]string)
	if raw == "" {
		return result
	}

	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		idx := strings.IndexByte(part, '=')
		if idx <= 0 {
			continue
		}
		key := strings.TrimSpace(part[:idx])
		val := strings.TrimSpace(part[idx+1:])
		if key != "" {
			result[key] = val
		}
	}

	return result
}
