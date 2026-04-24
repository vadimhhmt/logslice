package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Options configures the log reader.
type Options struct {
	// MaxLineBytes limits the maximum size of a single log line.
	// Defaults to 1MB if zero.
	MaxLineBytes int
}

// Reader reads log lines from a source.
type Reader struct {
	scanner *bufio.Scanner
	lines   chan string
	errs    chan error
}

// New creates a Reader from an io.Reader.
func New(r io.Reader, opts Options) *Reader {
	maxBytes := opts.MaxLineBytes
	if maxBytes <= 0 {
		maxBytes = 1024 * 1024 // 1MB default
	}

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, maxBytes), maxBytes)

	return &Reader{
		scanner: scanner,
		lines:   make(chan string, 64),
		errs:    make(chan error, 1),
	}
}

// NewFromFile opens a file and returns a Reader for it.
func NewFromFile(path string, opts Options) (*Reader, io.Closer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("reader: open file: %w", err)
	}
	return New(f, opts), f, nil
}

// Lines returns the channel on which scanned lines are delivered.
func (r *Reader) Lines() <-chan string {
	return r.lines
}

// Err returns the first non-EOF error encountered during scanning.
func (r *Reader) Err() error {
	select {
	case err := <-r.errs:
		return err
	default:
		return nil
	}
}

// Start begins reading lines asynchronously and closes Lines() when done.
func (r *Reader) Start() {
	go func() {
		defer close(r.lines)
		for r.scanner.Scan() {
			line := r.scanner.Text()
			if line == "" {
				continue
			}
			r.lines <- line
		}
		if err := r.scanner.Err(); err != nil {
			r.errs <- fmt.Errorf("reader: scan: %w", err)
		}
	}()
}
