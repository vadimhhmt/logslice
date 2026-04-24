// Package pipeline wires together reader, filter, stats, and output
// into a single processing pass over a structured log stream.
package pipeline

import (
	"io"

	"logslice/internal/filter"
	"logslice/internal/output"
	"logslice/internal/parser"
	"logslice/internal/reader"
	"logslice/internal/stats"
)

// Config holds all dependencies needed to run the pipeline.
type Config struct {
	Reader    *reader.Reader
	Filter    *filter.Filter
	Formatter *output.Formatter
	Collector *stats.Collector
}

// Result summarises what happened during a pipeline run.
type Result struct {
	Read    int
	Matched int
	Dropped int
}

// Run reads every line from the reader, parses it, applies the filter,
// feeds matching entries to the formatter and the stats collector.
// It returns a Result and the first non-EOF error encountered.
func Run(cfg Config, w io.Writer) (Result, error) {
	var res Result

	for {
		line, err := cfg.Reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return res, err
		}

		entry, parseErr := parser.ParseLine(line)
		if parseErr != nil {
			res.Dropped++
			continue
		}

		res.Read++

		if !cfg.Filter.Matches(entry) {
			res.Dropped++
			continue
		}

		res.Matched++

		if cfg.Collector != nil {
			cfg.Collector.Add(entry)
		}

		if writeErr := cfg.Formatter.Write(w, entry); writeErr != nil {
			return res, writeErr
		}
	}

	return res, nil
}
