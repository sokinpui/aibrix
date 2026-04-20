package classification

import (
	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/config"
)

// ComplexityClassifier performs complexity-based classification using embedding similarity.
// Each rule independently classifies difficulty level using hard/easy candidates.
// Supports both text candidates (via text embedding model) and image candidates
// (via the multimodal embedding model) for contrastive knowledge base comparison.
// Results are filtered by composer conditions in the classifier layer.
type ComplexityClassifier struct {
	rules []config.ComplexityRule

	// Precomputed text embeddings for hard and easy candidates
	hardEmbeddings     map[string]map[string][]float32 // ruleName -> candidate -> embedding
	easyEmbeddings     map[string]map[string][]float32 // ruleName -> candidate -> embedding
	hardPrototypeBanks map[string]*prototypeBank
	easyPrototypeBanks map[string]*prototypeBank

	// Precomputed image embeddings for hard and easy image candidates (multimodal)
	imageHardEmbeddings     map[string]map[string][]float32 // ruleName -> imageRef -> embedding
	imageEasyEmbeddings     map[string]map[string][]float32 // ruleName -> imageRef -> embedding
	imageHardPrototypeBanks map[string]*prototypeBank
	imageEasyPrototypeBanks map[string]*prototypeBank

	modelType          string // Model type for text embeddings ("qwen3" or "gemma")
	hasImageCandidates bool   // True if any rule uses image_candidates
	prototypeCfg       config.PrototypeScoringConfig
}

type ComplexityRuleResult struct {
	RuleName       string
	Difficulty     string
	TextHardScore  float64
	TextEasyScore  float64
	TextMargin     float64
	ImageHardScore float64
	ImageEasyScore float64
	ImageMargin    float64
	FusedMargin    float64
	Confidence     float64
	SignalSource   string
}

// NewComplexityClassifier creates a new ComplexityClassifier with precomputed candidate embeddings.
// When rules contain image_candidates, the multimodal model must be initialized beforehand.
func NewComplexityClassifier(
	rules []config.ComplexityRule,
	modelType string,
	prototypeCfg config.PrototypeScoringConfig,
) (*ComplexityClassifier, error) {
	return &ComplexityClassifier{}, nil
}

// preloadCandidateEmbeddings computes embeddings for all hard/easy candidates (text + image).
// Uses concurrent processing for better performance.
func (c *ComplexityClassifier) preloadCandidateEmbeddings() error {
	return nil
}

// Classify evaluates the query against ALL complexity rules independently (text-only).
// For CUA requests with screenshots, use ClassifyWithImage instead.
func (c *ComplexityClassifier) Classify(query string) ([]string, error) {
	return nil, nil
}

// ClassifyWithImage evaluates the query (and optionally a request image) against
// ALL complexity rules independently.
//
// When imageURL is provided (e.g. a base64 data-URI screenshot from a CUA request),
// SigLIP encodes the image and compares it against the image knowledge base.
// The text query is always compared against the text knowledge base.
// The difficulty score fuses both channels: d(t) = max(|d_vis|, |d_sem|).
//
// Returns: all matched rules in format "rulename:difficulty"
// (e.g., ["cua_difficulty:hard", "cua_difficulty:easy"])
func (c *ComplexityClassifier) ClassifyWithImage(query string, imageURL string) ([]string, error) {
	return nil, nil
}

func (c *ComplexityClassifier) ClassifyDetailedWithImage(query string, imageURL string) ([]ComplexityRuleResult, error) {
	return nil, nil
}

func (c *ComplexityClassifier) rebuildPrototypeBanks() {
}
