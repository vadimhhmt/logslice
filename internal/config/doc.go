// Package config provides CLI flag parsing and validation for the logslice
// tool.
//
// Usage:
//
//	cfg, err := config.Parse()
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Config fields map directly to CLI flags:
//
//	-file     path to input log file (omit to read from stdin)
//	-from     RFC3339 start of time window
//	-to       RFC3339 end of time window
//	-pattern  field=value filter pairs, comma-separated
//	-format   output format (json | pretty | raw)
//	-fields   comma-separated list of fields to keep in output
//	-stats    print summary statistics after processing
package config
