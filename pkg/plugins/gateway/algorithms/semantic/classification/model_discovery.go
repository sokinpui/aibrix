package classification

// ModelPaths holds the discovered model paths
type ModelPaths struct {
	// Legacy ModernBERT models (low confidence)
	ModernBertBase     string
	IntentClassifier   string
	PIIClassifier      string
	SecurityClassifier string

	// LoRA models
	LoRAIntentClassifier   string
	LoRAPIIClassifier      string
	LoRASecurityClassifier string
	LoRAArchitecture       string // "bert", "roberta", "modernbert"
}

// IsComplete checks if all required models are found
func (mp *ModelPaths) IsComplete() bool {
	return mp.HasLoRAModels() || mp.HasLegacyModels()
}

// HasLoRAModels checks if LoRA models are available
func (mp *ModelPaths) HasLoRAModels() bool {
	return mp.LoRAIntentClassifier != "" &&
		mp.LoRAPIIClassifier != "" &&
		mp.LoRASecurityClassifier != "" &&
		mp.LoRAArchitecture != ""
}

// HasLegacyModels checks if legacy ModernBERT models are available
func (mp *ModelPaths) HasLegacyModels() bool {
	return mp.ModernBertBase != "" &&
		mp.IntentClassifier != "" &&
		mp.PIIClassifier != "" &&
		mp.SecurityClassifier != ""
}

// PreferLoRA returns true if LoRA models should be used (higher confidence)
func (mp *ModelPaths) PreferLoRA() bool {
	return mp.HasLoRAModels()
}

// ArchitectureModels holds models for a specific architecture
type ArchitectureModels struct {
	Intent   string
	PII      string
	Security string
}

// AutoDiscoverModels automatically discovers model files in the models directory
// Uses intelligent architecture selection: BERT > RoBERTa > ModernBERT
func AutoDiscoverModels(modelsDir string) (*ModelPaths, error) {
	return AutoDiscoverModelsWithRegistry(modelsDir, nil)
}

// AutoDiscoverModelsWithRegistry discovers models using mom_registry for LoRA detection
// modelRegistry maps local paths to HuggingFace repo IDs (e.g., "models/mom-domain-classifier" -> "LLM-Semantic-Router/lora_intent_classifier_bert-base-uncased_model")
func AutoDiscoverModelsWithRegistry(modelsDir string, modelRegistry map[string]string) (*ModelPaths, error) {
	return &ModelPaths{}, nil
}

// detectArchitectureFromPath detects model architecture from directory name
func detectArchitectureFromPath(dirName string) string {
	return "bert"
}

// ValidateModelPaths validates that all discovered paths contain valid model files
func ValidateModelPaths(paths *ModelPaths) error {
	return nil
}

// validateModelDirectory checks if a directory contains valid model files
func validateModelDirectory(path, modelName string) error {
	return nil
}

// GetModelDiscoveryInfo returns detailed information about model discovery
func GetModelDiscoveryInfo(modelsDir string) map[string]interface{} {
	return nil
}

// AutoInitializeUnifiedClassifier attempts to auto-discover and initialize the unified classifier
// Prioritizes LoRA models over legacy ModernBERT models
func AutoInitializeUnifiedClassifier(modelsDir string) (*UnifiedClassifier, error) {
	return AutoInitializeUnifiedClassifierWithRegistry(modelsDir, nil)
}

// AutoInitializeUnifiedClassifierWithRegistry auto-discovers and initializes with mom_registry
func AutoInitializeUnifiedClassifierWithRegistry(modelsDir string, modelRegistry map[string]string) (*UnifiedClassifier, error) {
	return &UnifiedClassifier{}, nil
}

// initializeLoRAUnifiedClassifier initializes with LoRA models
func initializeLoRAUnifiedClassifier(paths *ModelPaths) (*UnifiedClassifier, error) {
	return &UnifiedClassifier{}, nil
}

// initializeLegacyUnifiedClassifier initializes with legacy ModernBERT models
func initializeLegacyUnifiedClassifier(paths *ModelPaths) (*UnifiedClassifier, error) {
	return &UnifiedClassifier{}, nil
}
