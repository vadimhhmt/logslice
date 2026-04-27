// Package label provides a pipeline stage that attaches a static key/value
// label to every log entry. Labels are commonly used to tag entries with
// metadata such as the originating file, host, or deployment environment
// before entries from multiple sources are merged into a single stream.
//
// Usage:
//
//	lbl, err := label.New("source", "api-server", false)
//	if err != nil { ... }
//
//	for _, e := range entries {
//		tagged := lbl.Apply(e)
//	}
//
// The overwrite flag controls whether an existing field with the same key
// is replaced (true) or left untouched (false).
package label
