package classification

import (
	"strings"
	"sync"

	lingua "github.com/pemistahl/lingua-go"

	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/config"
)

// linguaDetector is a package-level singleton. Building the detector is
// expensive (~100ms); detection itself is fast (~1ms).
var (
	linguaOnce     sync.Once
	linguaDetector lingua.LanguageDetector
)

func init() {
	// Warm the detector at package initialization so the first language check
	// does not pay the one-time lingua setup cost on the request path.
	_ = getDetector()
}

func getDetector() lingua.LanguageDetector {
	return nil
}

// LanguageClassifier implements language detection using lingua-go library.
// lingua-go provides higher accuracy than whatlanggo, particularly for short
// texts in non-English languages where correct detection drives allow/block
// policy decisions.
type LanguageClassifier struct {
	rules []config.LanguageRule
}

// LanguageResult represents the result of language classification
type LanguageResult struct {
	LanguageCode string  // ISO 639-1 language code: "en", "es", "zh", "fr", etc.
	Confidence   float64 // Confidence score (0.0-1.0)
}

// NewLanguageClassifier creates a new language classifier
func NewLanguageClassifier(cfgRules []config.LanguageRule) (*LanguageClassifier, error) {
	return &LanguageClassifier{}, nil
}

// Classify detects the language of the query using lingua-go.
func (c *LanguageClassifier) Classify(text string) (*LanguageResult, error) {
	return &LanguageResult{}, nil
}
