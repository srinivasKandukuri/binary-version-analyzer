package providers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// OpenAIProvider implements the AIProvider interface for OpenAI API
type OpenAIProvider struct {
	config *AIConfig
	client *openai.Client
}

// NewOpenAIProvider creates a new OpenAI provider with configuration
func NewOpenAIProvider(config *AIConfig) *OpenAIProvider {
	clientConfig := openai.DefaultConfig(config.APIKey)

	// Set custom base URL if provided
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	}

	// Set timeout
	clientConfig.HTTPClient.Timeout = time.Duration(config.Timeout) * time.Second

	client := openai.NewClientWithConfig(clientConfig)

	return &OpenAIProvider{
		config: config,
		client: client,
	}
}

// GetConfig returns the current configuration
func (o *OpenAIProvider) GetConfig() *AIConfig {
	return o.config
}

// UpdateConfig updates the provider configuration
func (o *OpenAIProvider) UpdateConfig(config *AIConfig) error {
	if err := ValidateConfig(config); err != nil {
		return err
	}

	// Create new client with updated config
	clientConfig := openai.DefaultConfig(config.APIKey)
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	}
	clientConfig.HTTPClient.Timeout = time.Duration(config.Timeout) * time.Second

	o.config = config
	o.client = openai.NewClientWithConfig(clientConfig)
	return nil
}

// SetModel allows changing the model used by OpenAI
func (o *OpenAIProvider) SetModel(model string) {
	o.config.Model = model
}

// SetTemperature allows changing the temperature
func (o *OpenAIProvider) SetTemperature(temp float64) {
	o.config.Temperature = temp
}

// SetMaxTokens allows changing the max tokens
func (o *OpenAIProvider) SetMaxTokens(tokens int) {
	o.config.MaxTokens = tokens
}

// AnalyzeVersions implements the AIProvider interface
func (o *OpenAIProvider) AnalyzeVersions(binaryName string, candidates []string) (string, error) {
	if len(candidates) == 0 {
		return "", fmt.Errorf("no version candidates provided")
	}

	prompt := o.buildPrompt(binaryName, candidates)

	req := openai.ChatCompletionRequest{
		Model: o.config.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a version number analyzer. Your task is to identify the most likely semantic version from a list of candidates. Respond with only the version number, nothing else.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens:   o.config.MaxTokens,
		Temperature: float32(o.config.Temperature),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(o.config.Timeout)*time.Second)
	defer cancel()

	resp, err := o.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("error calling OpenAI API: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI API")
	}

	version := strings.TrimSpace(resp.Choices[0].Message.Content)
	return version, nil
}

// GetProviderName returns the name of the provider
func (o *OpenAIProvider) GetProviderName() string {
	return "OpenAI"
}

// buildPrompt creates the prompt for version analysis
func (o *OpenAIProvider) buildPrompt(binaryName string, candidates []string) string {
	return fmt.Sprintf(`Given the following candidate strings, identify the most likely semantic version for the %s binary. Ignore unrelated floats or library dependencies.

Candidates:
%s

Please provide only the most likely version number in your response, nothing else.`, binaryName, "- "+strings.Join(candidates, "\n- "))
}
