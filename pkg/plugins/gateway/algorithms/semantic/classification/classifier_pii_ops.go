package classification

// ClassifyPII performs PII token classification on the given text and returns detected PII types
func (c *Classifier) ClassifyPII(text string) ([]string, error) {
	return nil, nil
}

// ClassifyPIIWithThreshold performs PII token classification with a custom threshold
func (c *Classifier) ClassifyPIIWithThreshold(text string, threshold float32) ([]string, error) {
	return nil, nil
}

// ClassifyPIIWithDetails performs PII token classification and returns full entity details including confidence scores
func (c *Classifier) ClassifyPIIWithDetails(text string) ([]PIIDetection, error) {
	return nil, nil
}

// ClassifyPIIWithDetailsAndThreshold performs PII token classification with a custom threshold and returns full entity details
func (c *Classifier) ClassifyPIIWithDetailsAndThreshold(text string, threshold float32) ([]PIIDetection, error) {
	return nil, nil
}

// DetectPIIInContent performs PII classification on all provided content
func (c *Classifier) DetectPIIInContent(allContent []string) []string {
	return nil
}

// AnalyzeContentForPII performs detailed PII analysis on multiple content pieces
func (c *Classifier) AnalyzeContentForPII(contentList []string) (bool, []PIIAnalysisResult, error) {
	return false, nil, nil
}

// AnalyzeContentForPIIWithThreshold performs detailed PII analysis with a custom threshold
func (c *Classifier) AnalyzeContentForPIIWithThreshold(contentList []string, threshold float32) (bool, []PIIAnalysisResult, error) {
	return false, nil, nil
}

// collectPIIRuleContents builds the list of text contents to analyze for a PII rule.
func collectPIIRuleContents(piiText string, nonUserMessages []string, includeHistory bool) []string {
	return nil
}

// collectPIIEntityTypes extracts entity types from cached PII results that meet the threshold.
func (c *Classifier) collectPIIEntityTypes(ruleContents []string, ruleName string, threshold float32, piiCache map[string]cachedPIIResult) map[string]bool {
	return nil
}

// findDeniedEntities returns entity types not covered by the allow-list.
func findDeniedEntities(entityTypes map[string]bool, allowedTypes []string) []string {
	return nil
}
