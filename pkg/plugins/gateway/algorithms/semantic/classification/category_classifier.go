package classification

import (
	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/utils/entropy"
)

// matchDomainCategories returns the domain categories that exceed the configured
// threshold, using entropy analysis to decide between top-1 and multi-category output.
func (c *Classifier) matchDomainCategories(
	domainResult InferenceClassResultWithProbs,
	topCategoryName string,
) []entropy.CategoryProbability {
	return nil
}

func (c *Classifier) buildCategoryNameMappings() {
}

// translateMMLUToGeneric translates an MMLU-Pro category to a generic category if mapping exists
func (c *Classifier) translateMMLUToGeneric(mmluCategory string) string {
	return ""
}
