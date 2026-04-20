package classification

type embeddingPreloadResult struct {
	candidate string
	embedding []float32
	err       error
}

func (c *EmbeddingClassifier) collectUniqueCandidates() []string {
	return nil
}

func (c *EmbeddingClassifier) preloadWorkerCount(candidateCount int) int {
	return 0
}

func (c *EmbeddingClassifier) startCandidateEmbeddingWorkers(
	candidates []string,
	modelType string,
	numWorkers int,
) <-chan embeddingPreloadResult {
	return nil
}

func (c *EmbeddingClassifier) collectCandidateEmbeddingResults(
	resultChan <-chan embeddingPreloadResult,
) (int, error) {
	return 0, nil
}
