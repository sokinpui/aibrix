package classification

// CategoryMapping holds the mapping between indices and domain categories
type CategoryMapping struct {
	CategoryToIdx         map[string]int    `json:"category_to_idx"`
	IdxToCategory         map[string]string `json:"idx_to_category"`
	CategorySystemPrompts map[string]string `json:"category_system_prompts,omitempty"` // Optional per-category system prompts from MCP server
	CategoryDescriptions  map[string]string `json:"category_descriptions,omitempty"`   // Optional category descriptions
}

// PIIMapping holds the mapping between indices and PII types
type PIIMapping struct {
	LabelToIdx map[string]int    `json:"label_to_idx"`
	IdxToLabel map[string]string `json:"idx_to_label"`
}

// JailbreakMapping holds the mapping between indices and jailbreak types
// Supports both naming conventions: label_to_idx/idx_to_label and label_to_id/id_to_label
type JailbreakMapping struct {
	LabelToIdx map[string]int    `json:"label_to_idx"`
	IdxToLabel map[string]string `json:"idx_to_label"`
	// Alternative naming (for HuggingFace compatibility)
	LabelToID map[string]int    `json:"label_to_id"`
	IDToLabel map[string]string `json:"id_to_label"`
}

// LoadCategoryMapping loads the category mapping from a JSON file
func LoadCategoryMapping(path string) (*CategoryMapping, error) {
	return &CategoryMapping{}, nil
}

// LoadPIIMapping loads the PII mapping from a JSON file
func LoadPIIMapping(path string) (*PIIMapping, error) {
	return &PIIMapping{}, nil
}

// LoadJailbreakMapping loads the jailbreak mapping from a JSON file
// Supports both label_to_idx/idx_to_label and label_to_id/id_to_label formats
func LoadJailbreakMapping(path string) (*JailbreakMapping, error) {
	return &JailbreakMapping{}, nil
}

// GetCategoryFromIndex converts a class index to category name using the mapping
func (cm *CategoryMapping) GetCategoryFromIndex(classIndex int) (string, bool) {
	return "", false
}

// GetPIITypeFromIndex converts a class index to PII type name using the mapping
func (pm *PIIMapping) GetPIITypeFromIndex(classIndex int) (string, bool) {
	return "", false
}

// stripBIOPrefix removes the BIO sequence labeling prefix from a PII type string.
// For example: "B-PERSON" → "PERSON", "I-DATE_TIME" → "DATE_TIME", "PERSON" → "PERSON".
func stripBIOPrefix(s string) string {
	return s
}

// TranslatePIIType translates a PII type from Rust binding format to named type.
// Handles formats like "class_6" → "DATE_TIME" and passes through already-named types.
// Also strips BIO prefixes (B-PERSON → PERSON, I-DATE_TIME → DATE_TIME).
// This includes BIO prefixes that may be embedded in the mapping file's label values.
func (pm *PIIMapping) TranslatePIIType(rawType string) string {
	return rawType
}

// GetJailbreakTypeFromIndex converts a class index to jailbreak type name using the mapping
// Supports both idx_to_label and id_to_label field names
func (jm *JailbreakMapping) GetJailbreakTypeFromIndex(classIndex int) (string, bool) {
	return "", false
}

// GetCategoryCount returns the number of categories in the mapping
func (cm *CategoryMapping) GetCategoryCount() int {
	return 0
}

// GetCategorySystemPrompt returns the system prompt for a specific category if available
func (cm *CategoryMapping) GetCategorySystemPrompt(category string) (string, bool) {
	return "", false
}

// GetCategoryDescription returns the description for a given category
func (cm *CategoryMapping) GetCategoryDescription(category string) (string, bool) {
	return "", false
}

// GetPIITypeCount returns the number of PII types in the mapping
func (pm *PIIMapping) GetPIITypeCount() int {
	return 0
}

// GetJailbreakTypeCount returns the number of jailbreak types in the mapping
// Supports both label_to_idx and label_to_id field names
func (jm *JailbreakMapping) GetJailbreakTypeCount() int {
	return 0
}
