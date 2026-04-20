package classification

import (
	"sync"

	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/config"
)

// Default feedback type labels (used as fallback if config.json doesn't have id2label)
const (
	FeedbackLabelSatisfied         = "satisfied"
	FeedbackLabelNeedClarification = "need_clarification"
	FeedbackLabelWrongAnswer       = "wrong_answer"
	FeedbackLabelWantDifferent     = "want_different"
)

// FeedbackResult represents the result of user feedback classification
type FeedbackResult struct {
	FeedbackType string  `json:"feedback_type"` // feedback type label from model's id2label
	Confidence   float32 `json:"confidence"`
	Class        int     `json:"class"` // class index from model
}

// FeedbackMapping maps feedback types to class indices
type FeedbackMapping struct {
	LabelToIdx map[string]int
	IdxToLabel map[string]string
}

// FeedbackDetector handles user feedback classification from follow-up messages
type FeedbackDetector struct {
	config       *config.FeedbackDetectorConfig
	mapping      *FeedbackMapping
	initialized  bool
	useMmBERT32K bool // Track if mmBERT-32K is used for inference
	mu           sync.RWMutex
}

// NewFeedbackDetector creates a new feedback detector
func NewFeedbackDetector(cfg *config.FeedbackDetectorConfig) (*FeedbackDetector, error) {
	return &FeedbackDetector{}, nil
}

// loadMappingFromConfig loads the id2label mapping from the model's config.json
func (d *FeedbackDetector) loadMappingFromConfig(modelPath string) error {
	return nil
}

// normalizeFeedbackLabel converts model labels (e.g., "SAT", "NEED_CLARIFICATION") to standard form
func normalizeFeedbackLabel(label string) string {
	return ""
}

// Initialize initializes the feedback detector with the ModernBERT model
func (d *FeedbackDetector) Initialize() error {
	return nil
}

// Classify determines user feedback type from follow-up message using the ML model
func (d *FeedbackDetector) Classify(text string) (*FeedbackResult, error) {
	return &FeedbackResult{}, nil
}

// IsInitialized returns whether the detector is initialized
func (d *FeedbackDetector) IsInitialized() bool {
	return false
}

// GetMapping returns the feedback mapping
func (d *FeedbackDetector) GetMapping() *FeedbackMapping {
	return nil
}
