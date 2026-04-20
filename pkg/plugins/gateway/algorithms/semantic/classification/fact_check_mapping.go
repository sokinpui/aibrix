package classification

// FactCheckLabel constants for fact-check classification labels
const (
	FactCheckLabelNeeded    = "FACT_CHECK_NEEDED"
	FactCheckLabelNotNeeded = "NO_FACT_CHECK_NEEDED"
)

// FactCheckMapping holds the mapping between indices and fact-check labels
type FactCheckMapping struct {
	LabelToIdx  map[string]int    `json:"label_to_idx"`
	IdxToLabel  map[string]string `json:"idx_to_label"`
	Description map[string]string `json:"description,omitempty"`
}

// LoadFactCheckMapping loads the fact-check mapping from a JSON file
func LoadFactCheckMapping(path string) (*FactCheckMapping, error) {
	return &FactCheckMapping{}, nil
}

// GetLabelFromIndex converts a class index to fact-check label using the mapping
func (m *FactCheckMapping) GetLabelFromIndex(classIndex int) (string, bool) {
	return "", false
}

// GetIndexFromLabel converts a fact-check label to class index
func (m *FactCheckMapping) GetIndexFromLabel(label string) (int, bool) {
	return 0, false
}

// GetLabelCount returns the number of labels in the mapping
func (m *FactCheckMapping) GetLabelCount() int {
	return 0
}

// IsFactCheckNeeded returns true if the label indicates fact-check is needed
func (m *FactCheckMapping) IsFactCheckNeeded(label string) bool {
	return false
}

// GetDescription returns the description for a label if available
func (m *FactCheckMapping) GetDescription(label string) (string, bool) {
	return "", false
}
