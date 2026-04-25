// Package aggregate provides field-based grouping and counting of log entries.
package aggregate

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/user/logslice/internal/parser"
)

// Aggregator groups log entries by the value of a specified field.
type Aggregator struct {
	field  string
	counts map[string]int
	total  int
}

// New returns an Aggregator that groups entries by the given field name.
func New(field string) *Aggregator {
	return &Aggregator{
		field:  field,
		counts: make(map[string]int),
	}
}

// Add records the value of the configured field from entry into the aggregation.
// Entries missing the field are counted under the key "<missing>".
func (a *Aggregator) Add(entry parser.Entry) {
	a.total++
	val, ok := entry.Fields[a.field]
	if !ok {
		a.counts["<missing>"]++
		return
	}
	key := fmt.Sprintf("%v", val)
	a.counts[key]++
}

// Results returns a sorted slice of (value, count) pairs.
func (a *Aggregator) Results() []Result {
	results := make([]Result, 0, len(a.counts))
	for k, v := range a.counts {
		results = append(results, Result{Value: k, Count: v})
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].Count != results[j].Count {
			return results[i].Count > results[j].Count
		}
		return results[i].Value < results[j].Value
	})
	return results
}

// Total returns the total number of entries processed.
func (a *Aggregator) Total() int { return a.total }

// Print writes a human-readable aggregation table to w.
func (a *Aggregator) Print(w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "FIELD\t%s\n", a.field)
	fmt.Fprintf(tw, "TOTAL\t%d\n", a.total)
	fmt.Fprintln(tw, "---\t---")
	for _, r := range a.Results() {
		fmt.Fprintf(tw, "%s\t%d\n", r.Value, r.Count)
	}
	tw.Flush()
}

// Result holds a single aggregation bucket.
type Result struct {
	Value string
	Count int
}
