package classification

import (
	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/config"
	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/utils/entropy"
)

// ModalityClassificationResult holds the result of modality signal classification
type ModalityClassificationResult struct {
	Modality   string  // "AR", "DIFFUSION", or "BOTH"
	Confidence float32 // Confidence score (0.0-1.0)
	Method     string  // Detection method used: "classifier", "keyword", or "hybrid"
}

// classifyModality determines the response modality for a text prompt.
// It supports three configurable methods via ModalityDetectionConfig:
//   - "classifier": ML-based (mmBERT-32K) — errors if model not loaded
//   - "keyword":    Configurable keyword matching — requires keywords in config
//   - "hybrid":     Classifier when available + keyword confirmation/fallback (default)
func (c *Classifier) classifyModality(text string, detectionConfig *config.ModalityDetectionConfig) ModalityClassificationResult {
	return ModalityClassificationResult{}
}

// classifyModalityByClassifier uses the mmBERT-32K ML classifier exclusively.
func (c *Classifier) classifyModalityByClassifier(text string, cfg *config.ModalityDetectionConfig) ModalityClassificationResult {
	return ModalityClassificationResult{}
}

// classifyModalityByKeyword uses keyword patterns from config to detect modality.
func (c *Classifier) classifyModalityByKeyword(text string, cfg *config.ModalityDetectionConfig) ModalityClassificationResult {
	return ModalityClassificationResult{}
}

// classifyModalityHybrid uses the ML classifier as primary, with keyword matching as
// fallback (when classifier is unavailable) or confirmation (when classifier confidence is low).
func (c *Classifier) classifyModalityHybrid(text string, cfg *config.ModalityDetectionConfig) ModalityClassificationResult {
	return ModalityClassificationResult{}
}

// ClassifyCategoryWithEntropy performs category classification with entropy-based reasoning decision
func (c *Classifier) ClassifyCategoryWithEntropy(text string) (string, float64, entropy.ReasoningDecision, error) {
	return "", 0.0, entropy.ReasoningDecision{}, nil
}

// tryKeywordBasedClassification attempts classification via keyword and embedding classifiers.
// Returns matched=true if a classifier produced a result.
func (c *Classifier) tryKeywordBasedClassification(text string) (string, float64, entropy.ReasoningDecision, bool, error) {
	return "", 0.0, entropy.ReasoningDecision{}, false, nil
}

// makeReasoningDecisionForKeywordCategory creates a reasoning decision for keyword-matched categories
func (c *Classifier) makeReasoningDecisionForKeywordCategory(category string) entropy.ReasoningDecision {
	return entropy.ReasoningDecision{}
}

// classifyCategoryWithEntropyInTree performs category classification with entropy using in-tree model
func (c *Classifier) classifyCategoryWithEntropyInTree(text string) (string, float64, entropy.ReasoningDecision, error) {
	return "", 0.0, entropy.ReasoningDecision{}, nil
}

func (c *Classifier) recordEntropyMetrics(probabilities []float32, reasoningDecision entropy.ReasoningDecision, entropyLatency float64) {
}
