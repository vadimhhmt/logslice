// Package merge implements multi-stream log merging for logslice.
//
// When logslice is given multiple input files (e.g. log shards from different
// hosts), each file is parsed into an independent channel of parser.Entry
// values. The Merger combines those channels into a single channel whose
// entries are emitted in ascending timestamp order using a min-heap.
//
// Usage:
//
//	sources := []<-chan parser.Entry{chanA, chanB, chanC}
//	m := merge.New(sources)
//	for entry := range m.Merge() {
//		// process entry in time order
//	}
//
// Each source channel must itself be sorted (or approximately sorted) for the
// global output to be strictly ordered. If a source emits entries out of
// order, those entries will appear in the output in the order they are
// received from that source relative to the current heap state.
package merge
