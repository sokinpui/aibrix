package classification

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/config"
)

// getEmbeddingWithModelType is a package-level variable for computing single embeddings.
// It exists so tests can override it.
var getEmbeddingWithModelType = GetEmbeddingWithModelType

// getMultiModalTextEmbedding computes a text embedding via the multimodal model.
// Package-level var so tests can override it.
var getMultiModalTextEmbedding = func(text string, targetDim int) ([]float32, error) {
	if MultiModalEncodeText == nil {
		return nil, fmt.Errorf("MultiModalEncodeText hook not initialized")
	}
	output, err := MultiModalEncodeText(text, targetDim)
	if err != nil {
		return nil, err
	}
	return output.Embedding, nil
}

// getMultiModalImageEmbedding computes an image embedding from a base64-encoded
// image (raw base64 or data-URI) via the multimodal model.
// Also supports local file paths for preloading knowledge-base image candidates.
// Package-level var so tests can override it.
var getMultiModalImageEmbedding = func(imageRef string, targetDim int) ([]float32, error) {
	if imageRef == "" {
		return nil, fmt.Errorf("imageRef cannot be empty")
	}

	payload := imageRef

	// If imageRef is a local file path, read and base64-encode it
	if strings.HasPrefix(imageRef, "/") || strings.HasPrefix(imageRef, "./") {
		data, err := os.ReadFile(imageRef)
		if err != nil {
			return nil, fmt.Errorf("failed to read image file %q: %w", imageRef, err)
		}
		payload = base64.StdEncoding.EncodeToString(data)
	} else if idx := strings.Index(imageRef, ";base64,"); idx >= 0 {
		// Strip data-URI prefix if present (e.g. "data:image/png;base64,...")
		payload = imageRef[idx+len(";base64,"):]
	}
	
	if MultiModalEncodeImage == nil {
		return nil, fmt.Errorf("MultiModalEncodeImage hook not initialized")
	}
	output, err := MultiModalEncodeImage(payload, targetDim)
	if err != nil {
		return nil, err
	}
	return output.Embedding, nil
}

// initMultiModalModel is a package-level var for initializing the multimodal model.
var initMultiModalModel = InitMultiModalModel

// EmbeddingClassifierInitializer initializes KeywordEmbeddingClassifier for embedding based classification
type EmbeddingClassifierInitializer interface {
	Init(qwen3ModelPath string, gemmaModelPath string, mmBertModelPath string, useCPU bool) error
}

type ExternalModelBasedEmbeddingInitializer struct{}

func (c *ExternalModelBasedEmbeddingInitializer) Init(qwen3ModelPath string, gemmaModelPath string, mmBertModelPath string, useCPU bool) error {
	return nil
}

// createEmbeddingInitializer creates the appropriate keyword embedding initializer based on configuration
func createEmbeddingInitializer() EmbeddingClassifierInitializer {
	return nil
}

// EmbeddingClassifier performs embedding-based similarity classification.
// When preloading is enabled, candidate embeddings are computed once at initialization
// and reused for all classification requests, significantly improving performance.
type EmbeddingClassifier struct {
	rules []config.EmbeddingRule

	// Optimization: preloaded candidate embeddings
	candidateEmbeddings map[string][]float32 // candidate text -> embedding vector
	rulePrototypeBanks  map[string]*prototypeBank

	// Configuration
	optimizationConfig config.HNSWConfig
	preloadEnabled     bool
	modelType          string // Model type to use for embeddings ("qwen3" or "gemma")
}

// NewEmbeddingClassifier creates a new EmbeddingClassifier.
// If optimization config has PreloadEmbeddings enabled, candidate embeddings
// will be precomputed at initialization time for better runtime performance.
func NewEmbeddingClassifier(cfgRules []config.EmbeddingRule, optConfig config.HNSWConfig) (*EmbeddingClassifier, error) {
	return &EmbeddingClassifier{}, nil
}

// preloadCandidateEmbeddings computes embeddings for all unique candidates across all rules
// Uses concurrent processing for better performance
func (c *EmbeddingClassifier) preloadCandidateEmbeddings() error {
	return nil
}

// getModelType returns the model type to use for embeddings
func (c *EmbeddingClassifier) getModelType() string {
	return ""
}

// IsKeywordEmbeddingClassifierEnabled checks if Keyword embedding classification rules are properly configured
func (c *Classifier) IsKeywordEmbeddingClassifierEnabled() bool {
	return false
}

// initializeKeywordEmbeddingClassifier initializes the KeywordEmbedding classification model
func (c *Classifier) initializeKeywordEmbeddingClassifier() error {
	return nil
}

// Classify performs Embedding similarity classification on the given text.
// Returns the single best matching rule. Wraps ClassifyAll internally.
func (c *EmbeddingClassifier) Classify(text string) (string, float64, error) {
	return "", 0.0, nil
}

// ClassifyAll performs embedding similarity classification on the given text.
// Returns the highest-ranking matched rules, limited by embedding_config.top_k
// (default 1, 0 disables truncation). When top_k is increased, the decision
// engine can compose multiple embedding matches together.
func (c *EmbeddingClassifier) ClassifyAll(text string) ([]MatchedRule, error) {
	return nil, nil
}

// ClassifyDetailed performs full label scoring and returns the complete score
// distribution plus all accepted matches before top-k output shaping.
func (c *EmbeddingClassifier) ClassifyDetailed(text string) (*EmbeddingClassificationResult, error) {
	return &EmbeddingClassificationResult{}, nil
}

// MatchedRule holds the result for a matched embedding rule
type MatchedRule struct {
	RuleName string
	Score    float64
	Method   string // "hard" or "soft"
}

type EmbeddingRuleScore struct {
	Name           string
	Score          float64
	Best           float64
	Support        float64
	Threshold      float64
	PrototypeCount int
}

type EmbeddingClassificationResult struct {
	Scores  []EmbeddingRuleScore
	Matches []MatchedRule
}

// findAllMatchedRules aggregates candidate similarities per rule and returns all
// accepted matches before final top-k output shaping.
func (c *EmbeddingClassifier) findAllMatchedRules(scoredRules []EmbeddingRuleScore) []MatchedRule {
	return nil
}

func (c *EmbeddingClassifier) scoreRules(queryEmbedding []float32) ([]EmbeddingRuleScore, error) {
	return nil, nil
}

func (c *EmbeddingClassifier) sortMatches(matches []MatchedRule) []MatchedRule {
	return nil
}

func (c *EmbeddingClassifier) sortAndLimitMatches(matches []MatchedRule) []MatchedRule {
	return nil
}

func (c *EmbeddingClassifier) embeddingAggregationOptions(rule config.EmbeddingRule) prototypeScoreOptions {
	return prototypeScoreOptions{}
}

// cosineSimilarity computes cosine similarity between two vectors.
// Assumes vectors are normalized (which they should be from BERT-style models).
func cosineSimilarity(a, b []float32) float32 {
	return 0
}

// GetPreloadStats returns statistics about preloaded embeddings
func (c *EmbeddingClassifier) GetPreloadStats() int {
	return 0
}

func (c *EmbeddingClassifier) rebuildRulePrototypeBanks() {
}
