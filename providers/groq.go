package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GroqProvider implements the AIProvider interface for Groq API
type GroqProvider struct {
	config *AIConfig
	client *http.Client
}

// GroqRequest represents the request structure for Groq API
type GroqRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GroqResponse represents the response from Groq API
type GroqResponse struct {
	Choices []Choice `json:"choices"`
}

// Choice represents a choice in the response
type Choice struct {
	Message Message `json:"message"`
}

// NewGroqProvider creates a new Groq AI provider with configuration
func NewGroqProvider(config *AIConfig) *GroqProvider {
	return &GroqProvider{
		config: config,
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}
}

// GetConfig returns the current configuration
func (g *GroqProvider) GetConfig() *AIConfig {
	return g.config
}

// UpdateConfig updates the provider configuration
func (g *GroqProvider) UpdateConfig(config *AIConfig) error {
	if err := ValidateConfig(config); err != nil {
		return err
	}
	g.config = config
	g.client.Timeout = time.Duration(config.Timeout) * time.Second
	return nil
}

// SetModel allows changing the model used by Groq
func (g *GroqProvider) SetModel(model string) {
	g.config.Model = model
}

// SetTemperature allows changing the temperature
func (g *GroqProvider) SetTemperature(temp float64) {
	g.config.Temperature = temp
}

// SetMaxTokens allows changing the max tokens
func (g *GroqProvider) SetMaxTokens(tokens int) {
	g.config.MaxTokens = tokens
}

// AnalyzeVersions implements the AIProvider interface
func (g *GroqProvider) AnalyzeVersions(binaryName string, candidates []string) (string, error) {
	if len(candidates) == 0 {
		return "", fmt.Errorf("no version candidates provided")
	}

	prompt := g.buildPrompt(binaryName, candidates)

	reqBody := GroqRequest{
		Model: g.config.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a version number analyzer. Your task is to identify the most likely semantic version from a list of candidates. Respond with only the version number, nothing else.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   g.config.MaxTokens,
		Temperature: g.config.Temperature,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	// Use base URL from config
	url := fmt.Sprintf("%s/chat/completions", g.config.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.config.APIKey)

	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var groqResp GroqResponse
	if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("no response from Groq API")
	}

	version := strings.TrimSpace(groqResp.Choices[0].Message.Content)
	return version, nil
}

// GetProviderName returns the name of the provider
func (g *GroqProvider) GetProviderName() string {
	return "Groq"
}

// buildPrompt creates the prompt for version analysis
func (g *GroqProvider) buildPrompt(binaryName string, candidates []string) string {
	return fmt.Sprintf(`Given the following candidate strings, identify the most likely semantic version for the %s binary. Ignore unrelated floats or library dependencies.

Candidates:
%s

Please provide only the most likely version number in your response, nothing else.`, binaryName, "- "+strings.Join(candidates, "\n- "))
}
