package providers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// AIConfig represents configuration for AI providers
type AIConfig struct {
	Provider    AIProviderType `json:"provider"`
	APIKey      string         `json:"api_key"`
	Model       string         `json:"model"`
	Temperature float64        `json:"temperature"`
	MaxTokens   int            `json:"max_tokens"`
	BaseURL     string         `json:"base_url,omitempty"`
	Timeout     int            `json:"timeout,omitempty"` // in seconds
}

// DefaultConfigs provides default configurations for each provider
var DefaultConfigs = map[AIProviderType]AIConfig{
	ProviderGroq: {
		Provider:    ProviderGroq,
		Model:       "llama-3.1-70b-versatile",
		Temperature: 0.1,
		MaxTokens:   50,
		BaseURL:     "https://api.groq.com/openai/v1",
		Timeout:     30,
	},
	ProviderOpenAI: {
		Provider:    ProviderOpenAI,
		Model:       "gpt-3.5-turbo",
		Temperature: 0.1,
		MaxTokens:   50,
		BaseURL:     "https://api.openai.com/v1",
		Timeout:     30,
	},
}

// LoadConfigFromEnv loads AI configuration from environment variables
func LoadConfigFromEnv() (*AIConfig, error) {
	// Determine provider
	providerStr := strings.ToLower(os.Getenv("AI_PROVIDER"))
	if providerStr == "" {
		providerStr = "groq" // default
	}

	var providerType AIProviderType
	switch providerStr {
	case "groq":
		providerType = ProviderGroq
	case "openai":
		providerType = ProviderOpenAI
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", providerStr)
	}

	// Start with default config
	config := DefaultConfigs[providerType]

	// Load API key
	var apiKey string
	switch providerType {
	case ProviderGroq:
		apiKey = os.Getenv("GROQ_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("GROQ_API_KEY environment variable is required")
		}
	case ProviderOpenAI:
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
		}
	}
	config.APIKey = apiKey

	// Override with environment variables if present
	if model := os.Getenv("AI_MODEL"); model != "" {
		config.Model = model
	}

	if tempStr := os.Getenv("AI_TEMPERATURE"); tempStr != "" {
		if temp, err := strconv.ParseFloat(tempStr, 64); err == nil {
			if temp >= 0.0 && temp <= 2.0 {
				config.Temperature = temp
			} else {
				return nil, fmt.Errorf("AI_TEMPERATURE must be between 0.0 and 2.0, got: %f", temp)
			}
		} else {
			return nil, fmt.Errorf("invalid AI_TEMPERATURE value: %s", tempStr)
		}
	}

	if tokensStr := os.Getenv("AI_MAX_TOKENS"); tokensStr != "" {
		if tokens, err := strconv.Atoi(tokensStr); err == nil {
			if tokens > 0 && tokens <= 4096 {
				config.MaxTokens = tokens
			} else {
				return nil, fmt.Errorf("AI_MAX_TOKENS must be between 1 and 4096, got: %d", tokens)
			}
		} else {
			return nil, fmt.Errorf("invalid AI_MAX_TOKENS value: %s", tokensStr)
		}
	}

	if baseURL := os.Getenv("AI_BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	if timeoutStr := os.Getenv("AI_TIMEOUT"); timeoutStr != "" {
		if timeout, err := strconv.Atoi(timeoutStr); err == nil {
			if timeout > 0 && timeout <= 300 {
				config.Timeout = timeout
			} else {
				return nil, fmt.Errorf("AI_TIMEOUT must be between 1 and 300 seconds, got: %d", timeout)
			}
		} else {
			return nil, fmt.Errorf("invalid AI_TIMEOUT value: %s", timeoutStr)
		}
	}

	return &config, nil
}

// GetProviderSpecificEnvVars returns provider-specific environment variable names
func GetProviderSpecificEnvVars(providerType AIProviderType) map[string]string {
	switch providerType {
	case ProviderGroq:
		return map[string]string{
			"API_KEY": "GROQ_API_KEY",
			"MODEL":   "GROQ_MODEL",
		}
	case ProviderOpenAI:
		return map[string]string{
			"API_KEY": "OPENAI_API_KEY",
			"MODEL":   "OPENAI_MODEL",
		}
	default:
		return map[string]string{}
	}
}

// ValidateConfig validates the AI configuration
func ValidateConfig(config *AIConfig) error {
	if config.APIKey == "" {
		return fmt.Errorf("API key is required")
	}

	if config.Model == "" {
		return fmt.Errorf("model is required")
	}

	if config.Temperature < 0.0 || config.Temperature > 2.0 {
		return fmt.Errorf("temperature must be between 0.0 and 2.0")
	}

	if config.MaxTokens <= 0 || config.MaxTokens > 4096 {
		return fmt.Errorf("max tokens must be between 1 and 4096")
	}

	if config.Timeout <= 0 || config.Timeout > 300 {
		return fmt.Errorf("timeout must be between 1 and 300 seconds")
	}

	return nil
}

// PrintConfigInfo prints the current configuration (without sensitive data)
func PrintConfigInfo(config *AIConfig) {
	fmt.Printf("🔧 AI Configuration:\n")
	fmt.Printf("   Provider: %s\n", config.Provider)
	fmt.Printf("   Model: %s\n", config.Model)
	fmt.Printf("   Temperature: %.2f\n", config.Temperature)
	fmt.Printf("   Max Tokens: %d\n", config.MaxTokens)
	if config.BaseURL != "" {
		fmt.Printf("   Base URL: %s\n", config.BaseURL)
	}
	fmt.Printf("   Timeout: %ds\n", config.Timeout)
	fmt.Printf("   API Key: %s***\n", config.APIKey[:min(8, len(config.APIKey))])
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
