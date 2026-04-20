package classification

// InferenceClassResult represents a basic classification result.
type InferenceClassResult struct {
	Class      int
	Confidence float32
}

// InferenceClassResultWithProbs represents a classification result with full probability distribution.
type InferenceClassResultWithProbs struct {
	Class         int
	Confidence    float32
	Probabilities []float32
}

// InferenceTokenResult represents the result of token-level classification (e.g., PII).
type InferenceTokenResult struct {
	Tokens []InferenceToken
}

// InferenceToken represents a single classified token.
type InferenceToken struct {
	Text       string
	Label      string
	Start      int
	End        int
	Confidence float32
}

// InferenceEmbeddingOutput represents a generated embedding vector.
type InferenceEmbeddingOutput struct {
	Embedding []float32
}

// InferenceNLILabel represents Natural Language Inference labels.
type InferenceNLILabel int

const (
	InferenceNLIEntailment InferenceNLILabel = iota
	InferenceNLINeutral
	InferenceNLIContradiction
	InferenceNLIError
)

// Inference Hooks - These package-level variables allow hollowing the implementation.
// They should be initialized by the runtime or during tests.
var (
	GetEmbedding               func(text string, targetDim int) ([]float32, error)
	GetEmbeddingWithModelType  func(text string, modelType string, targetDim int) (*InferenceEmbeddingOutput, error)
	MultiModalEncodeText       func(text string, targetDim int) (*InferenceEmbeddingOutput, error)
	MultiModalEncodeImage      func(base64Payload string, targetDim int) (*InferenceEmbeddingOutput, error)
	InitMultiModalModel        func(path string, useCPU bool) error
)
