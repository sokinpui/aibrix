package classification

import (
	"time"

	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/config"
)

// VLLMJailbreakInference implements JailbreakInference using vLLM REST API
type VLLMJailbreakInference struct {
	client     *VLLMClient
	modelName  string
	threshold  float32
	timeout    time.Duration
	parserType string // Parser type: "qwen3guard", "json", "simple", "auto"
}

// NewVLLMJailbreakInference creates a new vLLM-based jailbreak inference instance
// Takes ExternalModelConfig directly
func NewVLLMJailbreakInference(cfg *config.ExternalModelConfig, defaultThreshold float32) (*VLLMJailbreakInference, error) {
	return &VLLMJailbreakInference{}, nil
}

// Classify implements the JailbreakInference interface
func (v *VLLMJailbreakInference) Classify(text string) (InferenceClassResult, error) {
	return InferenceClassResult{}, nil
}

// parseSafetyOutput parses safety model output - uses parser type or auto-detection
func (v *VLLMJailbreakInference) parseSafetyOutput(output string) (bool, float32, []string) {
	return false, 0.0, nil
}

// determineParserType determines which parser to use based on config or model name
func (v *VLLMJailbreakInference) determineParserType() string {
	return ""
}

// parseQwen3GuardFormat parses Qwen3Guard structured output
func (v *VLLMJailbreakInference) parseQwen3GuardFormat(output string) (bool, float32, []string) {
	return false, 0.0, nil
}

// extractCategories extracts violation categories from Qwen3Guard output
// Returns empty slice if no categories found or if "None" is specified
func (v *VLLMJailbreakInference) extractCategories(output string) []string {
	return nil
}

// parseJSONFormat parses JSON output
func (v *VLLMJailbreakInference) parseJSONFormat(output string) (bool, float32) {
	return false, 0.0
}

// parseSimpleFormat parses simple keyword-based output
func (v *VLLMJailbreakInference) parseSimpleFormat(output string) (bool, float32) {
	return false, 0.0
}
