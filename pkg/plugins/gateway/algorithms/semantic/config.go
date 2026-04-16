package semantic

import "os"

// Constants for Routing Signals and Roles
const (
	SignalTypeKeyword   = "keyword"
	SignalTypeEmbedding = "embedding"
	SignalTypeDomain    = "domain"

	ModelRoleClassification = "classification"
	ModelRoleScoring        = "scoring"

	APIFormatOpenAI    = "openai"
	APIFormatAnthropic = "anthropic"
)

// RouterConfig represents the core configuration for the LLM Router.
type RouterConfig struct {
	Version      string                     `yaml:"version"`
	DefaultModel string                     `yaml:"default_model"`
	ModelConfig  map[string]ModelParams     `yaml:"model_config"`
	Decisions    []Decision                 `yaml:"decisions"`
	Signals      Signals                    `yaml:"signals"`
	Endpoints    []VLLMEndpoint             `yaml:"vllm_endpoints"`
	Profiles     map[string]ProviderProfile `yaml:"provider_profiles,omitempty"`
	Global       GlobalOptions              `yaml:"global"`
}

// GlobalOptions captures engine-level control knobs.
type GlobalOptions struct {
	AutoModelName   string               `yaml:"auto_model_name,omitempty"`
	ClearRouteCache bool                 `yaml:"clear_route_cache"`
	ModelSelection  ModelSelectionConfig `yaml:"model_selection"`
}

// ModelSelectionConfig defines global defaults for the selection algorithm.
type ModelSelectionConfig struct {
	Method   string                  `yaml:"method"` // "static" or "router_dc"
	RouterDC RouterDCSelectionConfig `yaml:"router_dc,omitempty"`
}

// RouterDCSelectionConfig specific to Dual Contrastive learning.
type RouterDCSelectionConfig struct {
	Temperature         float64 `yaml:"temperature"`
	MinSimilarity       float64 `yaml:"min_similarity"`
	UseCapabilities     bool    `yaml:"use_capabilities"`
	RequireDescriptions bool    `yaml:"require_descriptions"`
}

// Decision represents a single routing intent.
type Decision struct {
	Name        string           `yaml:"name"`
	Description string           `yaml:"description,omitempty"`
	Priority    int              `yaml:"priority,omitempty"`
	Rules       RuleNode         `yaml:"rules"`
	ModelRefs   []ModelRef       `yaml:"modelRefs,omitempty"`
	Algorithm   *AlgorithmConfig `yaml:"algorithm,omitempty"`
}

// AlgorithmConfig allows per-decision algorithm overrides.
type AlgorithmConfig struct {
	Type     string                   `yaml:"type"`
	RouterDC *RouterDCSelectionConfig `yaml:"router_dc,omitempty"`
}

// ModelRef references a specific model in the catalog.
type ModelRef struct {
	Model    string  `yaml:"model"`
	LoRAName string  `yaml:"lora_name,omitempty"`
	Weight   float64 `yaml:"weight,omitempty"`
}

// RuleNode is a boolean expression tree.
type RuleNode struct {
	Type       string     `yaml:"type,omitempty"`
	Name       string     `yaml:"name,omitempty"`
	Operator   string     `yaml:"operator,omitempty"`
	Conditions []RuleNode `yaml:"conditions,omitempty"`
}

// ModelParams defines capabilities and metadata for a specific LLM.
type ModelParams struct {
	Description        string            `yaml:"description,omitempty"`
	Capabilities       []string          `yaml:"capabilities,omitempty"`
	QualityScore       float64           `yaml:"quality_score,omitempty"`
	Pricing            ModelPricing      `yaml:"pricing,omitempty"`
	PreferredEndpoints []string          `yaml:"preferred_endpoints,omitempty"`
	APIFormat          string            `yaml:"api_format,omitempty"`
	AccessKey          string            `yaml:"access_key,omitempty"`
	ExternalModelIDs   map[string]string `yaml:"external_model_ids,omitempty"`
}

type ModelPricing struct {
	Currency        string  `yaml:"currency,omitempty"`
	PromptPer1M     float64 `yaml:"prompt_per_1m,omitempty"`
	CompletionPer1M float64 `yaml:"completion_per_1m,omitempty"`
}

// Signals contains the definitions for rule-based matching.
type Signals struct {
	KeywordRules   []KeywordRule   `yaml:"keyword_rules,omitempty"`
	EmbeddingRules []EmbeddingRule `yaml:"embedding_rules,omitempty"`
	Categories     []Category      `yaml:"categories,omitempty"`
}

type KeywordRule struct {
	Name     string   `yaml:"name"`
	Keywords []string `yaml:"keywords"`
	Operator string   `yaml:"operator"` // "or", "and"
}

type EmbeddingRule struct {
	Name       string   `yaml:"name"`
	Threshold  float32  `yaml:"threshold"`
	Candidates []string `yaml:"candidates"`
}

type Category struct {
	Name        string       `yaml:"name"`
	Description string       `yaml:"description,omitempty"`
	ModelScores []ModelScore `yaml:"model_scores,omitempty"`
}

type ModelScore struct {
	Model string  `yaml:"model"`
	Score float64 `yaml:"score"`
}

// VLLMEndpoint represents a physical backend.
type VLLMEndpoint struct {
	Name                string `yaml:"name"`
	Address             string `yaml:"address"`
	Port                int    `yaml:"port"`
	Weight              int    `yaml:"weight,omitempty"`
	Type                string `yaml:"type,omitempty"`
	ProviderProfileName string `yaml:"provider_profile,omitempty"`
}

// ProviderProfile for cloud LLM providers (OpenAI, etc.).
type ProviderProfile struct {
	Type    string `yaml:"type"`
	BaseURL string `yaml:"base_url,omitempty"`
	APIKey  string `yaml:"api_key,omitempty"`
}

// GetAccessKey expands environment variables for API keys.
func (p ModelParams) GetAccessKey() string {
	if p.AccessKey == "" {
		return ""
	}
	return os.ExpandEnv(p.AccessKey)
}
