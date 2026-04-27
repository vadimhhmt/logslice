// Package merge provides a Merger that combines multiple sorted log entry
// streams into a single time-ordered stream.
package merge

import (
	"container/heap"
	"time"

	"github.com/your-org/logslice/internal/parser"
)

// entrySource pairs a parsed log entry with the index of the source stream
// it came from.
type entrySource struct {
	entry  parser.Entry
	source int
}

// minHeap implements heap.Interface over entrySource, ordered by timestamp.
type minHeap []entrySource

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].entry.Time.Before(h[j].entry.Time) }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minHeap) Push(x any) {
	*h = append(*h, x.(entrySource))
}

func (h *minHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Merger merges multiple channels of log entries into one time-ordered stream.
type Merger struct {
	sources []<-chan parser.Entry
}

// New creates a Merger over the provided source channels.
func New(sources []<-chan parser.Entry) *Merger {
	return &Merger{sources: sources}
}

// Merge reads from all source channels and emits entries in ascending
// timestamp order. Entries with a zero timestamp are passed through
// immediately in arrival order. The returned channel is closed once all
// sources are exhausted.
func (m *Merger) Merge() <-chan parser.Entry {
	out := make(chan parser.Entry, 64)

	go func() {
		defer close(out)

		h := &minHeap{}
		heap.Init(h)

		channels := make([]<-chan parser.Entry, len(m.sources))
		copy(channels, m.sources)

		// Prime the heap with one entry from each source.
		active := make([]bool, len(channels))
		for i, ch := range channels {
			if e, ok := <-ch; ok {
				heap.Push(h, entrySource{entry: e, source: i})
				active[i] = true
			}
		}

		for h.Len() > 0 {
			es := heap.Pop(h).(entrySource)
			out <- es.entry

			// Refill from the same source.
			if e, ok := <-channels[es.source]; ok {
				heap.Push(h, entrySource{entry: e, source: es.source})
			}
		}
	}()

	return out
}

// ZeroTime is a sentinel used for entries whose timestamp could not be parsed.
var ZeroTime = time.Time{}
