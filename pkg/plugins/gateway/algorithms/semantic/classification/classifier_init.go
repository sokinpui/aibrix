package classification

import (
	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/config"
)

type CategoryInitializer interface {
	Init(modelID string, useCPU bool, numClasses ...int) error
}

type CategoryInitializerImpl struct {
	usedModernBERT bool // Track which init path succeeded for inference routing
}

func (c *CategoryInitializerImpl) Init(modelID string, useCPU bool, numClasses ...int) error {
	return nil
}

// MmBERT32KCategoryInitializerImpl uses mmBERT-32K (YaRN RoPE, 32K context) for intent classification
type MmBERT32KCategoryInitializerImpl struct {
	usedMmBERT32K bool
}

func (c *MmBERT32KCategoryInitializerImpl) Init(modelID string, useCPU bool, numClasses ...int) error {
	return nil
}

// createCategoryInitializer creates the category initializer (auto-detecting)
func createCategoryInitializer() CategoryInitializer {
	return &CategoryInitializerImpl{}
}

// createMmBERT32KCategoryInitializer creates an mmBERT-32K category initializer
func createMmBERT32KCategoryInitializer() CategoryInitializer {
	return &MmBERT32KCategoryInitializerImpl{}
}

type CategoryInference interface {
	Classify(text string) (InferenceClassResult, error)
	ClassifyWithProbabilities(text string) (InferenceClassResultWithProbs, error)
}

type CategoryInferenceImpl struct{}

func (c *CategoryInferenceImpl) Classify(text string) (InferenceClassResult, error) {
	return InferenceClassResult{}, nil
}

func (c *CategoryInferenceImpl) ClassifyWithProbabilities(text string) (InferenceClassResultWithProbs, error) {
	return InferenceClassResultWithProbs{}, nil
}

// createCategoryInference creates the category inference (auto-detecting)
func createCategoryInference() CategoryInference {
	return &CategoryInferenceImpl{}
}

// MmBERT32KCategoryInferenceImpl uses mmBERT-32K for intent classification
type MmBERT32KCategoryInferenceImpl struct{}

func (c *MmBERT32KCategoryInferenceImpl) Classify(text string) (InferenceClassResult, error) {
	return InferenceClassResult{}, nil
}

func (c *MmBERT32KCategoryInferenceImpl) ClassifyWithProbabilities(text string) (InferenceClassResultWithProbs, error) {
	return InferenceClassResultWithProbs{}, nil
}

// createMmBERT32KCategoryInference creates mmBERT-32K category inference
func createMmBERT32KCategoryInference() CategoryInference {
	return &MmBERT32KCategoryInferenceImpl{}
}

type JailbreakInitializer interface {
	Init(modelID string, useCPU bool, numClasses ...int) error
}

type JailbreakInitializerImpl struct {
	usedModernBERT bool // Track which init path succeeded for inference routing
}

func (c *JailbreakInitializerImpl) Init(modelID string, useCPU bool, numClasses ...int) error {
	return nil
}

// createJailbreakInitializer creates the jailbreak initializer (auto-detecting)
func createJailbreakInitializer() JailbreakInitializer {
	return &JailbreakInitializerImpl{}
}

// MmBERT32KJailbreakInitializerImpl uses mmBERT-32K (YaRN RoPE, 32K context) for jailbreak detection
type MmBERT32KJailbreakInitializerImpl struct {
	usedMmBERT32K bool
}

func (c *MmBERT32KJailbreakInitializerImpl) Init(modelID string, useCPU bool, numClasses ...int) error {
	return nil
}

// createMmBERT32KJailbreakInitializer creates an mmBERT-32K jailbreak initializer
func createMmBERT32KJailbreakInitializer() JailbreakInitializer {
	return &MmBERT32KJailbreakInitializerImpl{}
}

type JailbreakInference interface {
	Classify(text string) (InferenceClassResult, error)
}

type JailbreakInferenceImpl struct{}

func (c *JailbreakInferenceImpl) Classify(text string) (InferenceClassResult, error) {
	return InferenceClassResult{}, nil
}

// createJailbreakInferenceCandle creates Candle-based jailbreak inference (auto-detecting)
func createJailbreakInferenceCandle() JailbreakInference {
	return &JailbreakInferenceImpl{}
}

// MmBERT32KJailbreakInferenceImpl uses mmBERT-32K for jailbreak detection
type MmBERT32KJailbreakInferenceImpl struct{}

func (c *MmBERT32KJailbreakInferenceImpl) Classify(text string) (InferenceClassResult, error) {
	return InferenceClassResult{}, nil
}

// createMmBERT32KJailbreakInference creates mmBERT-32K jailbreak inference
func createMmBERT32KJailbreakInference() JailbreakInference {
	return &MmBERT32KJailbreakInferenceImpl{}
}

// createJailbreakInference creates the appropriate jailbreak inference based on configuration
// Checks UseMmBERT32K and UseVLLM flags to decide between mmBERT-32K, vLLM, or Candle implementation
// When UseMmBERT32K is true, uses mmBERT-32K (32K context, YaRN RoPE, multilingual)
// When UseVLLM is true, it will try to find external model config with role="guardrail"
func createJailbreakInference(promptGuardCfg *config.PromptGuardConfig, routerCfg *config.RouterConfig) (JailbreakInference, error) {
	return createJailbreakInferenceCandle(), nil
}

type PIIInitializer interface {
	Init(modelID string, useCPU bool, numClasses int) error
}

type PIIInitializerImpl struct {
	usedModernBERT bool // Track which init path succeeded for inference routing
}

func (c *PIIInitializerImpl) Init(modelID string, useCPU bool, numClasses int) error {
	return nil
}

// createPIIInitializer creates the PII initializer (auto-detecting)
func createPIIInitializer() PIIInitializer {
	return &PIIInitializerImpl{}
}

// MmBERT32KPIIInitializerImpl uses mmBERT-32K (YaRN RoPE, 32K context) for PII detection
type MmBERT32KPIIInitializerImpl struct {
	usedMmBERT32K bool
}

func (c *MmBERT32KPIIInitializerImpl) Init(modelID string, useCPU bool, numClasses int) error {
	return nil
}

// createMmBERT32KPIIInitializer creates an mmBERT-32K PII initializer
func createMmBERT32KPIIInitializer() PIIInitializer {
	return &MmBERT32KPIIInitializerImpl{}
}

type PIIInference interface {
	ClassifyTokens(text string) (InferenceTokenResult, error)
}

type PIIInferenceImpl struct{}

func (c *PIIInferenceImpl) ClassifyTokens(text string) (InferenceTokenResult, error) {
	return InferenceTokenResult{}, nil
}

// createPIIInference creates the PII inference (auto-detecting)
func createPIIInference() PIIInference {
	return &PIIInferenceImpl{}
}

// MmBERT32KPIIInferenceImpl uses mmBERT-32K for PII token classification.
// Entity types are returned as "LABEL_{class_id}" by Rust and translated Go-side via PIIMapping.
type MmBERT32KPIIInferenceImpl struct{}

func (c *MmBERT32KPIIInferenceImpl) ClassifyTokens(text string) (InferenceTokenResult, error) {
	return InferenceTokenResult{}, nil
}

// createMmBERT32KPIIInference creates mmBERT-32K PII inference
func createMmBERT32KPIIInference() PIIInference {
	return &MmBERT32KPIIInferenceImpl{}
}

// JailbreakDetection represents the result of jailbreak analysis for a piece of content
type JailbreakDetection struct {
	Content       string  `json:"content"`
	IsJailbreak   bool    `json:"is_jailbreak"`
	JailbreakType string  `json:"jailbreak_type"`
	Confidence    float32 `json:"confidence"`
	ContentIndex  int     `json:"content_index"`
}

// PIIDetection represents detected PII entities in content
type PIIDetection struct {
	EntityType string  `json:"entity_type"` // Type of PII entity (e.g., "PERSON", "EMAIL", "PHONE")
	Start      int     `json:"start"`       // Start character position in original text
	End        int     `json:"end"`         // End character position in original text
	Text       string  `json:"text"`        // Actual entity text
	Confidence float32 `json:"confidence"`  // Confidence score (0.0 to 1.0)
}

// PIIAnalysisResult represents the result of PII analysis for content
type PIIAnalysisResult struct {
	Content      string         `json:"content"`
	HasPII       bool           `json:"has_pii"`
	Entities     []PIIDetection `json:"entities"`
	ContentIndex int            `json:"content_index"`
}
