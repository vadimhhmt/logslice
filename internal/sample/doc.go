// Package sample implements log entry sampling for logslice.
//
// Two strategies are supported:
//
//   - StrategyRate: deterministic 1-in-N sampling. Every Nth entry
//     seen by Accept is kept. Useful for high-volume streams where
//     a predictable fraction is needed.
//
//   - StrategyReservoir: reservoir (Knuth Algorithm R) sampling.
//     Collect entries as they arrive, then call Flush to retrieve
//     a statistically uniform random sample of at most `size` entries
//     regardless of input length.
//
// Example — rate sampling:
//
//	s, _ := sample.New(sample.StrategyRate, 10, 0, 0)
//	for _, e := range entries {
//		if s.Accept(e) { process(e) }
//	}
//
// Example — reservoir sampling:
//
//	s, _ := sample.New(sample.StrategyReservoir, 0, 1000, time.Now().UnixNano())
//	for _, e := range entries { s.Collect(e) }
//	for _, e := range s.Flush() { process(e) }
package sample
