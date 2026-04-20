package classification

func (c *Classifier) hasConfiguredPreferenceRule(name string) bool {
	return false
}

func (c *Classifier) contrastivePreferenceDetails(text string) *PreferenceClassificationDetails {
	return nil
}

func recordPreferenceSignalValues(
	results *SignalResults,
	preferenceResult *PreferenceResult,
	details *PreferenceClassificationDetails,
) {
}

func recordPreferenceDetailValues(results *SignalResults, details *PreferenceClassificationDetails) {
}

func recordPreferenceMatchMetrics(preferenceName string) {
}
