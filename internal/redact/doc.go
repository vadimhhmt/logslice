// Package redact provides field-level redaction for structured log entries.
//
// A Redactor masks values whose field name matches a configured sensitive-key
// list (case-insensitive) or whose string value matches one of the registered
// regular expressions. Redacted values are replaced with the literal string
// "[REDACTED]".
//
// Usage:
//
//	r := redact.New()
//	r.AddKey("ssn")
//	r.AddPattern(`^Bearer\s+`)
//	clean := r.Apply(entry.Fields)
package redact
