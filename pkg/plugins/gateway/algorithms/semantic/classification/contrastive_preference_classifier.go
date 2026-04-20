package classification

import (
	"errors"
	"sync"

	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/config"
)

var ErrPreferenceBelowThreshold = errors.New("preference below threshold")

type preferenceEmbeddingTask struct {
	ruleName string
	text     string
}

type preferenceEmbeddingResult struct {
	ruleName  string
	embedding []float32
	err       error
}

// ContrastivePreferenceClassifier performs few-shot preference routing using embeddings.
// It preloads embeddings for each preference rule's examples/description and selects
// the route whose support set is most similar to the incoming query.
type ContrastivePreferenceClassifier struct {
	modelType string

	rules []config.PreferenceRule

	// ruleEmbeddings maps rule name to its support embeddings
	ruleEmbeddings map[string][][]float32
	ruleBanks      map[string]*prototypeBank
	// ruleThresholds stores per-preference similarity thresholds
	ruleThresholds map[string]float32
	prototypeCfg   config.PrototypeScoringConfig

	mu sync.RWMutex
}

type PreferenceRuleScore struct {
	Name           string
	Score          float32
	Best           float32
	Support        float32
	Threshold      float32
	PrototypeCount int
}

type PreferenceClassificationDetails struct {
	Scores       []PreferenceRuleScore
	BestRule     string
	BestScore    float32
	RunnerUpRule string
	RunnerUp     float32
	Margin       float32
}

// NewContrastivePreferenceClassifier builds a contrastive preference classifier.
// modelType follows GetEmbeddingWithModelType (e.g. "qwen3", "gemma", "mmbert").
func NewContrastivePreferenceClassifier(rules []config.PreferenceRule, modelType string) (*ContrastivePreferenceClassifier, error) {
	return &ContrastivePreferenceClassifier{}, nil
}

func NewContrastivePreferenceClassifierWithConfig(
	rules []config.PreferenceRule,
	modelType string,
	prototypeCfg config.PrototypeScoringConfig,
) (*ContrastivePreferenceClassifier, error) {
	return &ContrastivePreferenceClassifier{}, nil
}

// preloadRuleEmbeddings computes embeddings for all rule examples concurrently.
func (c *ContrastivePreferenceClassifier) preloadRuleEmbeddings() error {
	return nil
}

func (c *ContrastivePreferenceClassifier) collectEmbeddingTasks() ([]preferenceEmbeddingTask, error) {
	return nil, nil
}

func (c *ContrastivePreferenceClassifier) embedRuleExamples(
	tasks []preferenceEmbeddingTask,
) <-chan preferenceEmbeddingResult {
	return nil
}

func (c *ContrastivePreferenceClassifier) collectEmbeddedResults(
	resultCh <-chan preferenceEmbeddingResult,
) (int, error) {
	return 0, nil
}

func (c *ContrastivePreferenceClassifier) embeddingWorkerCount(taskCount int) int {
	return 0
}

// Classify picks the preference with the highest similarity to the query.
func (c *ContrastivePreferenceClassifier) Classify(text string) (*PreferenceResult, error) {
	return &PreferenceResult{}, nil
}

func (c *ContrastivePreferenceClassifier) ClassifyDetailed(text string) (*PreferenceClassificationDetails, error) {
	return &PreferenceClassificationDetails{}, nil
}

func (c *ContrastivePreferenceClassifier) collectExamples(rule config.PreferenceRule) []string {
	return nil
}

func (c *ContrastivePreferenceClassifier) rebuildRuleBanks() {
}
