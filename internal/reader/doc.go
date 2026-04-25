// Package reader provides line-by-line reading of structured log input.
//
// It supports reading from an io.Reader or directly from a named file.
// Empty lines are automatically skipped. The reader exposes a channel-based
// API so that it can be used in streaming pipelines without loading the
// entire file into memory.
//
// Basic usage:
//
//	r, err := reader.NewFromFile("app.log")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for line := range r.Lines() {
//		fmt.Println(line)
//	}
package reader
