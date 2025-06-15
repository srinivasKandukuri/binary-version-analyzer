package providers

import (
	"fmt"
)

// AIProviderType represents the type of AI provider
type AIProviderType string

const (
	ProviderGroq   AIProviderType = "groq"
	ProviderOpenAI AIProviderType = "openai"
)

// AIFactory creates AI providers based on configuration
type AIFactory struct{}

// NewAIFactory creates a new AI factory
func NewAIFactory() *AIFactory {
	return &AIFactory{}
}

// CreateProvider creates an AI provider based on the specified type and configuration
func (f *AIFactory) CreateProvider(config *AIConfig) (AIProvider, error) {
	// Validate configuration
	if err := ValidateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %v", err)
	}

	switch config.Provider {
	case ProviderGroq:
		return NewGroqProvider(config), nil
	case ProviderOpenAI:
		return NewOpenAIProvider(config), nil
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", config.Provider)
	}
}

// CreateProviderFromEnv creates an AI provider from environment variables
func (f *AIFactory) CreateProviderFromEnv() (AIProvider, error) {
	// Load configuration from environment
	config, err := LoadConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	// Create provider with configuration
	return f.CreateProvider(config)
}

// CreateProviderWithDefaults creates an AI provider with default configuration
func (f *AIFactory) CreateProviderWithDefaults(providerType AIProviderType, apiKey string) (AIProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	// Get default config and set API key
	config := DefaultConfigs[providerType]
	config.APIKey = apiKey

	return f.CreateProvider(&config)
}

// GetSupportedProviders returns a list of supported AI providers
func (f *AIFactory) GetSupportedProviders() []AIProviderType {
	return []AIProviderType{ProviderGroq, ProviderOpenAI}
}

// GetDefaultConfig returns the default configuration for a provider
func (f *AIFactory) GetDefaultConfig(providerType AIProviderType) (AIConfig, error) {
	config, exists := DefaultConfigs[providerType]
	if !exists {
		return AIConfig{}, fmt.Errorf("unsupported provider: %s", providerType)
	}
	return config, nil
}
