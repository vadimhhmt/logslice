package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"logslice/internal/filter"
	"logslice/internal/output"
	"logslice/internal/parser"
	"logslice/internal/reader"
)

const timeLayout = time.RFC3339

func main() {
	var (
		start   = flag.String("start", "", "start time (RFC3339), e.g. 2024-01-01T00:00:00Z")
		end     = flag.String("end", "", "end time (RFC3339), e.g. 2024-01-02T00:00:00Z")
		pattern = flag.String("pattern", "", "field pattern filter, e.g. level=error")
		format  = flag.String("format", "json", "output format: json, pretty, raw")
		fields  = flag.String("fields", "", "comma-separated fields to include in output")
		input   = flag.String("file", "", "input log file (default: stdin)")
	)
	flag.Parse()

	var startTime, endTime time.Time
	var parseErr error

	if *start != "" {
		startTime, parseErr = time.Parse(timeLayout, *start)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "invalid --start time: %v\n", parseErr)
			os.Exit(1)
		}
	}

	if *end != "" {
		endTime, parseErr = time.Parse(timeLayout, *end)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "invalid --end time: %v\n", parseErr)
			os.Exit(1)
		}
	}

	f, err := filter.New(startTime, endTime, *pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid filter: %v\n", err)
		os.Exit(1)
	}

	fmt := output.New(*format, *fields)

	var r *reader.Reader
	if *input != "" {
		r, err = reader.NewFromFile(*input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot open file: %v\n", err)
			os.Exit(1)
		}
	} else {
		r = reader.New(os.Stdin)
	}

	for r.Next() {
		line := r.Line()
		entry, parseErr := parser.ParseLine(line)
		if parseErr != nil {
			continue
		}
		if !f.InRange(entry) {
			continue
		}
		if out, fmtErr := fmt.Format(entry); fmtErr == nil {
			os.Stdout.WriteString(out + "\n")
		}
	}

	if err := r.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "read error: %v\n", err)
		os.Exit(1)
	}
}
