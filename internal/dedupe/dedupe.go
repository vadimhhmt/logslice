// Package dedupe provides log entry deduplication based on a configurable
// fingerprint of selected fields. Consecutive duplicate entries are dropped;
// a count of suppressed lines is optionally injected into the surviving entry.
package dedupe

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"logslice/internal/parser"
)

// Deduplicator tracks the last-seen fingerprint and drops consecutive
// duplicate log entries.
type Deduplicator struct {
	fields      []string // fields used to build the fingerprint
	injectCount bool     // add "_suppressed" field to surviving entry
	lastKey     string
	suppressed  int
}

// New returns a Deduplicator that fingerprints entries using the given fields.
// When fields is empty every field in the entry is used.
// If injectCount is true, a "_suppressed" integer field is added to the entry
// that survives after a run of duplicates.
func New(fields []string, injectCount bool) *Deduplicator {
	return &Deduplicator{
		fields:      fields,
		injectCount: injectCount,
	}
}

// Process evaluates entry against the last-seen fingerprint.
// It returns (entry, true) when the entry should be forwarded and
// (zero, false) when it is a duplicate and should be dropped.
// Callers must flush any pending suppressed count by calling Flush at EOF.
func (d *Deduplicator) Process(entry parser.Entry) (parser.Entry, bool) {
	key := d.fingerprint(entry)

	if key == d.lastKey {
		d.suppressed++
		return parser.Entry{}, false
	}

	// New unique entry — annotate the previous surviving entry's count before
	// we return this one.  We can only annotate the *current* outgoing entry
	// with the count of duplicates that preceded it (i.e. duplicates of the
	// previous key), so we embed the count on the new entry instead.
	out := entry
	if d.injectCount && d.suppressed > 0 {
		out = cloneEntry(entry)
		out.Fields["_suppressed"] = d.suppressed
	}

	d.lastKey = key
	d.suppressed = 0
	return out, true
}

// Flush returns a synthetic summary entry if there are suppressed duplicates
// still pending at end-of-stream, otherwise returns false.
func (d *Deduplicator) Flush() (int, bool) {
	if d.suppressed == 0 {
		return 0, false
	}
	n := d.suppressed
	d.suppressed = 0
	return n, true
}

// fingerprint builds a stable hash key from the selected fields of entry.
func (d *Deduplicator) fingerprint(entry parser.Entry) string {
	subset := make(map[string]any)
	if len(d.fields) == 0 {
		subset = entry.Fields
	} else {
		for _, f := range d.fields {
			if v, ok := entry.Fields[f]; ok {
				subset[f] = v
			}
		}
	}
	b, err := json.Marshal(subset)
	if err != nil {
		return fmt.Sprintf("%v", subset)
	}
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

// cloneEntry returns a shallow copy of entry with a new Fields map.
func cloneEntry(e parser.Entry) parser.Entry {
	fields := make(map[string]any, len(e.Fields)+1)
	for k, v := range e.Fields {
		fields[k] = v
	}
	e.Fields = fields
	return e
}
