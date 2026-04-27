// Package flatten implements a log entry processor that collapses nested
// JSON objects into dot-separated (or custom-separator) top-level keys.
//
// Given an entry such as:
//
//	{"user": {"id": "42", "name": "alice"}, "level": "info"}
//
// the Flattener produces:
//
//	{"user.id": "42", "user.name": "alice", "level": "info"}
//
// This is useful when downstream processors or output formatters expect a
// flat key space, or when filtering on deeply nested fields.
//
// Usage:
//
//	f := flatten.New(".", 0)   // separator=".", unlimited depth
//	out := f.Process(entry)
package flatten
