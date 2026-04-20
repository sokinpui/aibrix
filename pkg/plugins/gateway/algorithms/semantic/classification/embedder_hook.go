package classification

import (
	"sync"
)

var (
	embedderMu             sync.Mutex
	embedderOverrideActive bool
)

// SetEmbeddingFuncForTests overrides the embedding generator for tests/benchmarks.
// It returns a restore function that must be called to revert to the original implementation.
func SetEmbeddingFuncForTests(fn func(string, string, int) (*InferenceEmbeddingOutput, error)) func() {
	return func() {}
}
