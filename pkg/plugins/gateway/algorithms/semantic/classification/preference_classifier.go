package classification

import (
	"time"

	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/config"
)

// PreferenceResult represents the result of preference classification
type PreferenceResult struct {
	Preference string  `json:"route"` // The matched route name
	Confidence float32 `json:"confidence,omitempty"`
	Margin     float32 `json:"margin,omitempty"`
}

// PreferenceClassifier handles route preference matching via external LLM
type PreferenceClassifier struct {
	client             *VLLMClient
	modelName          string
	timeout            time.Duration
	preferenceRules    []config.PreferenceRule
	systemPrompt       string
	userPromptTemplate string
	contrastive        *ContrastivePreferenceClassifier

	useContrastive bool
}

// NewPreferenceClassifier creates a new preference classifier
func NewPreferenceClassifier(externalCfg *config.ExternalModelConfig, rules []config.PreferenceRule, localCfg *config.PreferenceModelConfig) (*PreferenceClassifier, error) {
	return &PreferenceClassifier{}, nil
}

// Classify determines the best route preference for the given conversation
func (p *PreferenceClassifier) Classify(conversationJSON string) (*PreferenceResult, error) {
	return &PreferenceResult{}, nil
}

// buildRoutesJSON builds the routes JSON array from preference rules
func (p *PreferenceClassifier) buildRoutesJSON() (string, error) {
	return "", nil
}

// parsePreferenceOutput parses the JSON output from LLM
func (p *PreferenceClassifier) parsePreferenceOutput(output string) (*PreferenceResult, error) {
	return nil, nil
}

// IsInitialized returns true if the classifier is initialized
func (p *PreferenceClassifier) IsInitialized() bool {
	return false
}

// classifyContrastive runs few-shot contrastive routing using embeddings.
func (p *PreferenceClassifier) classifyContrastive(conversationJSON string) (*PreferenceResult, error) {
	return nil, nil
}
