package providers

// AIProvider defines the interface for AI providers
type AIProvider interface {
	AnalyzeVersions(binaryName string, candidates []string) (string, error)
	GetProviderName() string
}

// AIRequest represents a common request structure for AI analysis
type AIRequest struct {
	BinaryName  string   `json:"binary_name"`
	Candidates  []string `json:"candidates"`
	Temperature float64  `json:"temperature,omitempty"`
	MaxTokens   int      `json:"max_tokens,omitempty"`
}

// AIResponse represents a common response structure from AI providers
type AIResponse struct {
	Version      string  `json:"version"`
	Confidence   float64 `json:"confidence,omitempty"`
	ProviderName string  `json:"provider_name"`
}
