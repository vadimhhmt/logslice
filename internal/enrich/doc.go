// Package enrich provides field enrichment for structured log entries.
//
// An Enricher holds a list of Rules. Each Rule reads a source field, applies a
// transformation function, and writes the result to a destination field. The
// original entry is never mutated; Apply always returns a new copy.
//
// Built-in transformation functions:
//
//	UpperCase  – converts the source value to upper-case
//	LowerCase  – converts the source value to lower-case
//
// Custom functions with the signature func(string) (string, bool) can be
// supplied directly when constructing a Rule. Returning false from the
// function causes the destination field to be left untouched.
//
// CLI integration:
//
//	-enrich from_field:to_field:fn
//
// The flag may be repeated to register multiple rules. Supported fn values
// are "upper" and "lower".
package enrich
