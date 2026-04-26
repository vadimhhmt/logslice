// Package tail implements a fixed-capacity ring buffer that retains the
// last N parsed log entries, providing behaviour analogous to `tail -n`.
//
// Usage:
//
//	tr := tail.New(20)   // keep last 20 entries
//	for _, e := range entries {
//		tr.Push(e)
//	}
//	for _, e := range tr.Entries() {
//		fmt.Println(e.Raw)
//	}
//
// The ring buffer uses O(n) memory regardless of how many entries are
// pushed, making it safe to use on arbitrarily large log streams.
package tail
