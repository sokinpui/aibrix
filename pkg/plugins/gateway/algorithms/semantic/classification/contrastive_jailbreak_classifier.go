package classification

import (
	"time"

	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/config"
)

// ContrastiveJailbreakResult holds the analysis outcome for a single rule.
type ContrastiveJailbreakResult struct {
	MaxScore       float32 // Highest contrastive score across analysed messages
	WorstMessage   string  // The message that produced MaxScore
	WorstMsgIndex  int     // Index of that message in the input slice
	JailbreakSim   float32 // max_sim(worstMsg, jailbreak_kb) for the worst message
	BenignSim      float32 // max_sim(worstMsg, benign_kb) for the worst message
	TotalMessages  int     // Number of messages analysed
	ProcessingTime time.Duration
}

// ContrastiveJailbreakClassifier implements contrastive embedding similarity
// for jailbreak detection. It mirrors the ComplexityClassifier pattern:
// pre-computed KB embeddings at init, fast cosine scoring at request time.
//
// Score = max_sim(msg, jailbreak_kb) − max_sim(msg, benign_kb)
// When include_history is true, the maximum score across all user messages
// in the conversation is used (multi-turn chain detection).
type ContrastiveJailbreakClassifier struct {
	rule config.JailbreakRule

	// Pre-computed embeddings for the two knowledge bases
	jailbreakEmbeddings map[string][]float32 // pattern text → embedding
	benignEmbeddings    map[string][]float32 // pattern text → embedding

	modelType string
}

// NewContrastiveJailbreakClassifier creates and initialises a classifier for a
// single contrastive JailbreakRule. KB embeddings are computed eagerly using a
// worker pool (same approach as ComplexityClassifier).
func NewContrastiveJailbreakClassifier(rule config.JailbreakRule, defaultModelType string) (*ContrastiveJailbreakClassifier, error) {
	return &ContrastiveJailbreakClassifier{}, nil
}

// AnalyzeMessages computes the contrastive score for each message and returns
// the result with the maximum score (multi-turn chain detection).
// If messages is empty the returned MaxScore is -1.
func (c *ContrastiveJailbreakClassifier) AnalyzeMessages(messages []string) ContrastiveJailbreakResult {
	return ContrastiveJailbreakResult{}
}

// preloadKBEmbeddings concurrently computes embeddings for jailbreak and benign
// pattern knowledge bases, following the same worker-pool approach as
// ComplexityClassifier.preloadCandidateEmbeddings.
func (c *ContrastiveJailbreakClassifier) preloadKBEmbeddings() error {
	return nil
}
