// Package sample provides log entry sampling strategies for logslice.
// It supports rate-based (1-in-N) and reservoir sampling.
package sample

import (
	"fmt"
	"math/rand"
	"sync"

	"logslice/internal/parser"
)

// Strategy defines how entries are sampled.
type Strategy int

const (
	StrategyRate     Strategy = iota // keep every Nth entry
	StrategyReservoir               // reservoir sampling of fixed size
)

// Sampler filters a stream of log entries according to a sampling strategy.
type Sampler struct {
	strategy  Strategy
	rate      int
	size      int
	counter   int
	reservoir []parser.Entry
	mu        sync.Mutex
	rng       *rand.Rand
}

// New creates a Sampler. For StrategyRate, rate is the keep interval (e.g. 3 = keep 1-in-3).
// For StrategyReservoir, size is the reservoir capacity.
func New(strategy Strategy, rate, size int, seed int64) (*Sampler, error) {
	if strategy == StrategyRate && rate < 1 {
		return nil, fmt.Errorf("sample: rate must be >= 1, got %d", rate)
	}
	if strategy == StrategyReservoir && size < 1 {
		return nil, fmt.Errorf("sample: reservoir size must be >= 1, got %d", size)
	}
	return &Sampler{
		strategy: strategy,
		rate:     rate,
		size:     size,
		rng:      rand.New(rand.NewSource(seed)),
	}, nil
}

// Accept returns true if the entry should be kept under StrategyRate.
// For StrategyReservoir, use Collect + Flush instead.
func (s *Sampler) Accept(e parser.Entry) bool {
	if s.strategy != StrategyRate {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	return s.counter%s.rate == 0
}

// Collect adds an entry to the reservoir (StrategyReservoir only).
func (s *Sampler) Collect(e parser.Entry) {
	if s.strategy != StrategyReservoir {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	if len(s.reservoir) < s.size {
		s.reservoir = append(s.reservoir, e)
		return
	}
	j := s.rng.Intn(s.counter)
	if j < s.size {
		s.reservoir[j] = e
	}
}

// Flush returns the collected reservoir entries and resets state.
func (s *Sampler) Flush() []parser.Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]parser.Entry, len(s.reservoir))
	copy(out, s.reservoir)
	s.reservoir = s.reservoir[:0]
	s.counter = 0
	return out
}
