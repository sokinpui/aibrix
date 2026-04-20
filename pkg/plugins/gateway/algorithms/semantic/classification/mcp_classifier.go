package classification

import (
	"context"

	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/config"

	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/utils/entropy"
)

const (
	// DefaultMCPThreshold is the default confidence threshold for MCP classification.
	// For multi-class classification, a value of 0.5 means that a predicted class must have at least 50% confidence
	// to be selected. Adjust this threshold as needed for your use case.
	DefaultMCPThreshold = 0.5
)

// MCPClassificationResult holds the classification result with routing information from MCP server
type MCPClassificationResult struct {
	Class        int
	Confidence   float32
	CategoryName string
	Model        string // Model recommended by MCP server
	UseReasoning *bool  // Whether to use reasoning (nil means use default)
}

// MCPCategoryInitializer initializes MCP connection for category classification
type MCPCategoryInitializer interface {
	Init(cfg *config.RouterConfig) error
	Close() error
}

// MCPCategoryInference performs classification via MCP
type MCPCategoryInference interface {
	Classify(ctx context.Context, text string) (InferenceClassResult, error)
	ClassifyWithProbabilities(ctx context.Context, text string) (InferenceClassResultWithProbs, error)
	ListCategories(ctx context.Context) (*CategoryMapping, error)
}

// MCPCategoryClassifier implements both MCPCategoryInitializer and MCPCategoryInference.
//
// Protocol Contract:
// This client relies on the MCP server to respect the protocol defined in the
// github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/mcp/api package.
//
// The MCP server must implement these tools:
//  1. list_categories - Returns api.ListCategoriesResponse
//  2. classify_text - Returns api.ClassifyResponse or api.ClassifyWithProbabilitiesResponse
//
// The MCP server controls both classification AND routing decisions. When the server returns
// "model" and "use_reasoning" in the classification response, the router will use those values.
// If not provided, the router falls back to the default_model configuration.
//
// For detailed type definitions and examples, see the api package documentation.
type MCPCategoryClassifier struct {
	toolName string
	config   *config.RouterConfig
}

// Init initializes the MCP client connection
func (m *MCPCategoryClassifier) Init(cfg *config.RouterConfig) error {
	return nil
}

// discoverClassificationTool finds the appropriate classification tool from available MCP tools
func (m *MCPCategoryClassifier) discoverClassificationTool() error {
	return nil
}

// Close closes the MCP client connection
func (m *MCPCategoryClassifier) Close() error {
	return nil
}

// Classify performs category classification via MCP
func (m *MCPCategoryClassifier) Classify(ctx context.Context, text string) (InferenceClassResult, error) {
	return InferenceClassResult{}, nil
}

// ClassifyWithProbabilities performs category classification with full probability distribution via MCP
func (m *MCPCategoryClassifier) ClassifyWithProbabilities(ctx context.Context, text string) (InferenceClassResultWithProbs, error) {
	return InferenceClassResultWithProbs{}, nil
}

// ListCategories retrieves the category mapping from the MCP server
func (m *MCPCategoryClassifier) ListCategories(ctx context.Context) (*CategoryMapping, error) {
	return nil, nil
}

// createMCPCategoryInitializer creates an MCP category initializer
func createMCPCategoryInitializer() MCPCategoryInitializer {
	return &MCPCategoryClassifier{}
}

// createMCPCategoryInference creates an MCP category inference from the initializer
func createMCPCategoryInference(initializer MCPCategoryInitializer) MCPCategoryInference {
	if inf, ok := initializer.(MCPCategoryInference); ok {
		return inf
	}
	return &MCPCategoryClassifier{}
}

// IsMCPCategoryEnabled checks if MCP-based category classification is properly configured.
// Note: tool_name is optional and will be auto-discovered during initialization if not specified.
func (c *Classifier) IsMCPCategoryEnabled() bool {
	return false
}

// initializeMCPCategoryClassifier initializes the MCP category classification model
func (c *Classifier) initializeMCPCategoryClassifier() error {
	return nil
}

// classifyCategoryWithEntropyMCP performs category classification with entropy using MCP
func (c *Classifier) classifyCategoryWithEntropyMCP(text string) (string, float64, entropy.ReasoningDecision, error) {
	return "", 0.0, entropy.ReasoningDecision{}, nil
}

// withMCPCategory creates an option function for MCP category classifier
func withMCPCategory(mcpInitializer MCPCategoryInitializer, mcpInference MCPCategoryInference) option {
	return func(c *Classifier) {}
}
