// Package aggregate groups parsed log entries by the value of a chosen field
// and reports per-value counts.
//
// Usage:
//
//	a := aggregate.New("level")
//	for _, entry := range entries {
//		a.Add(entry)
//	}
//	a.Print(os.Stdout)
//
// Results are sorted by count descending, with ties broken alphabetically.
// Entries that do not contain the target field are counted under "<missing>".
package aggregate
