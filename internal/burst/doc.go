// Package burst provides burst detection for structured log streams.
//
// A burst is defined as a high concentration of log entries sharing a common
// field value (e.g. the same "level", "service", or "error" key) within a
// sliding time window. When the number of matching entries exceeds a
// configurable threshold, the detector emits a synthetic summary entry
// annotating the burst before continuing to forward the original entries.
//
// # Usage
//
//	detector := burst.New(burst.Config{
//		Field:     "level",
//		Threshold: 10,
//		Window:    30 * time.Second,
//	})
//
//	for _, entry := range entries {
//		results := detector.Push(entry)
//		for _, r := range results {
//			fmt.Println(r)
//		}
//	}
//
//	// Flush any remaining buffered entries at end of stream.
//	for _, r := range detector.Flush() {
//		fmt.Println(r)
//	}
//
// # Burst Summary Entry
//
// When a burst is detected the injected summary entry contains:
//
//	{
//		"_burst": true,
//		"_burst_field": "<field>",
//		"_burst_value": "<value>",
//		"_burst_count": <n>,
//		"_burst_window_seconds": <w>
//	}
//
// The summary is prepended to the slice returned by Push so consumers see
// the annotation before the entries that triggered it.
package burst
