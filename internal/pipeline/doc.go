// Package pipeline provides a single-pass processing engine for logslice.
//
// It connects the reader, parser, filter, stats collector, and output
// formatter into one coordinated Run call. Callers supply a Config
// containing pre-constructed dependencies; Run iterates every log line,
// drops lines that cannot be parsed or do not match the active filter,
// and writes matching entries to the supplied io.Writer via the formatter.
//
// A Result value is returned so callers can report how many lines were
// read, matched, and dropped without needing access to internal state.
package pipeline
