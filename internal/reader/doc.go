// Package reader provides an asynchronous line-by-line reader for log sources.
//
// It wraps a standard io.Reader (or file path) and delivers non-empty lines
// over a channel, making it straightforward to integrate with the parser and
// filter pipeline:
//
//	r, closer, err := reader.NewFromFile("app.log", reader.Options{})
//	if err != nil { ... }
//	defer closer.Close()
//
//	r.Start()
//	for line := range r.Lines() {
//		entry, err := parser.ParseLine(line)
//		...
//	}
//	if err := r.Err(); err != nil { ... }
package reader
