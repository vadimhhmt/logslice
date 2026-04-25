// Package transform provides field-level transformations for log entries.
//
// Transformations include:
//   - Renaming fields (e.g. "msg" → "message")
//   - Remapping field values (e.g. level "warn" → "warning")
//
// Rules are applied in order; each rule operates on the result of the
// previous one. The original LogEntry is never mutated.
//
// CLI flags:
//
//	-rename old=new    rename a field before output
//	-remap field:v=w   replace a specific field value before output
package transform
