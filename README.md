# Binary Version Analyzer

This tool analyzes binary files to extract and identify their versions using a combination of regex pattern matching and AI-powered analysis with pluggable AI providers.

## Features

- 🔍 Scans binary files for potential version strings using regex patterns
- 📊 Extracts version candidates from binary content
- 🤖 **Pluggable AI Architecture** - Easily switch between different AI providers
- 🎯 Uses AI to analyze and determine the most likely version
- 🔧 Supports multiple version string patterns
- 📝 Handles various version formats and notations

## Supported AI Providers

- **Groq** (default) - Fast and free AI inference
- **OpenAI** - GPT models for version analysis
- **Extensible** - Easy to add new AI providers

## Prerequisites

- Go 1.18 or later
- API key for your chosen AI provider

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Build the binary:
   ```bash
   go build -o binary-version-analyzer
   ```

## Usage

### Basic Analysis

```bash
# Using Groq (default)
export GROQ_API_KEY="your-groq-api-key"
binary-version-analyzer analyze /usr/bin/ls

# Using OpenAI
export AI_PROVIDER="openai"
export OPENAI_API_KEY="your-openai-api-key"
binary-version-analyzer analyze /usr/bin/curl
```

### Command Structure

The tool uses a modern CLI structure with subcommands:

```bash
binary-version-analyzer [command] [flags]

Available Commands:
  analyze     Analyze a binary file to detect its version
  patterns    Manage and test version detection patterns
  help        Help about any command
```

### Advanced Configuration Examples

#### Using Command-Line Flags
```bash
# Custom model and temperature via flags
binary-version-analyzer analyze /usr/bin/curl \
  --provider groq \
  --model "llama-3.3-70b-versatile" \
  --temperature 0.2

# OpenAI with GPT-4 and custom settings
binary-version-analyzer analyze /usr/bin/python3 \
  --provider openai \
  --model gpt-4 \
  --temperature 0.1 \
  --max-tokens 100 \
  --timeout 60

# Save results to JSON file
binary-version-analyzer analyze /usr/bin/git \
  --output json \
  --save results.json
```

#### Using Environment Variables
```bash
# Custom model and temperature
export GROQ_API_KEY="your-groq-api-key"
export AI_MODEL="llama-3.3-70b-versatile"
export AI_TEMPERATURE="0.2"
binary-version-analyzer analyze /usr/bin/curl

# Custom base URL (for proxies or self-hosted)
export AI_BASE_URL="https://your-proxy.com/v1"
binary-version-analyzer analyze /usr/bin/git
```

### Environment Variables

#### Required Variables
| Variable | Description | Required |
|----------|-------------|----------|
| `GROQ_API_KEY` | API key for Groq | Yes (if using Groq) |
| `OPENAI_API_KEY` | API key for OpenAI | Yes (if using OpenAI) |

#### Configuration Variables
| Variable | Description | Default | Range |
|----------|-------------|---------|-------|
| `AI_PROVIDER` | AI provider to use (`groq`, `openai`) | `groq` | - |
| `AI_MODEL` | Override default model for provider | Provider default | - |
| `AI_TEMPERATURE` | Control randomness in AI responses | `0.1` | 0.0-2.0 |
| `AI_MAX_TOKENS` | Maximum response tokens | `50` | 1-4096 |
| `AI_BASE_URL` | Custom API base URL | Provider default | - |
| `AI_TIMEOUT` | Request timeout in seconds | `30` | 1-300 |

#### Default Models
- **Groq**: `llama-3.1-70b-versatile`
- **OpenAI**: `gpt-3.5-turbo`

## Project Structure

```
binary-version-analyzer/
├── main.go                    # Main application entry point
├── cmd/                       # Cobra CLI commands
│   ├── root.go               # Root command and global flags
│   ├── analyze.go            # Binary analysis command
│   └── patterns.go           # Pattern management commands
├── internal/                  # Internal application logic
│   └── analyzer.go           # Binary analyzer and result types
├── providers/                 # AI provider implementations
│   ├── interface.go          # AIProvider interface definition
│   ├── config.go            # Configuration management
│   ├── groq.go              # Groq API implementation
│   ├── openai.go            # OpenAI API implementation
│   └── factory.go           # Factory pattern for creating providers
├── patterns/                  # Version detection patterns
│   └── version_patterns.go  # Regex patterns with documentation
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
├── README.md                 # Project documentation
└── PATTERNS.md              # Detailed pattern documentation
```

## Architecture

The project uses a clean, modular architecture with the following components:

### Core Interface
```go
type AIProvider interface {
    AnalyzeVersions(binaryName string, candidates []string) (string, error)
    GetProviderName() string
}
```

### Components
- **`providers/interface.go`** - Defines the AIProvider interface and common types
- **`providers/groq.go`** - Groq API implementation
- **`providers/openai.go`** - OpenAI API implementation  
- **`providers/factory.go`** - Factory pattern for creating providers
- **`main.go`** - Main application logic and binary scanning

## Adding New AI Providers

To add a new AI provider:

1. Create a new file in the `providers/` directory (e.g., `claude.go`)
2. Implement the `AIProvider` interface:
   ```go
   package providers
   
   type ClaudeProvider struct {
       // your implementation
   }
   
   func NewClaudeProvider(apiKey string) *ClaudeProvider {
       // your constructor
   }
   
   func (c *ClaudeProvider) AnalyzeVersions(binaryName string, candidates []string) (string, error) {
       // your implementation
   }
   
   func (c *ClaudeProvider) GetProviderName() string {
       return "Claude"
   }
   ```
3. Update `providers/factory.go` to include your new provider:
   ```go
   const (
       ProviderGroq   AIProviderType = "groq"
       ProviderOpenAI AIProviderType = "openai"
       ProviderClaude AIProviderType = "claude"  // Add this
   )
   ```
4. Add the case in the factory's `CreateProvider` method
5. Add environment variable handling in `CreateProviderFromEnv`

## Version Patterns

The tool uses 15 sophisticated regex patterns to detect version strings:

### Pattern Categories
- **Priority 1**: Most reliable patterns (Standard version declarations, Stable releases)
- **Priority 2**: Common formats (V-prefixed, Semantic versions, Release keywords)
- **Priority 3**: System libraries (GLIBC, Library versions)
- **Priority 4**: Build contexts (Build versions, Package versions)
- **Priority 5**: Development tools (Compiler versions, API versions)
- **Priority 6-8**: Fallback patterns (Date-based, Copyright years)

### Pattern Management

```bash
# List all patterns
binary-version-analyzer patterns list

# List patterns with details
binary-version-analyzer patterns list --details

# Show patterns by priority
binary-version-analyzer patterns list --priority 1

# Test a string against all patterns
binary-version-analyzer patterns test "version 1.2.3"

# Interactive pattern testing
binary-version-analyzer patterns test --interactive

# Validate all patterns
binary-version-analyzer patterns validate

# Show detailed pattern documentation
binary-version-analyzer patterns docs
```

For detailed pattern documentation, see [PATTERNS.md](PATTERNS.md).

## Example Output

```
🔍 Analyzing binary: /usr/bin/ls
🤖 Using AI Provider: Groq
🔧 AI Configuration:
   Provider: groq
   Model: llama-3.1-70b-versatile
   Temperature: 0.10
   Max Tokens: 50
   Base URL: https://api.groq.com/openai/v1
   Timeout: 30s
   API Key: gsk_1234***

📊 Scanning for version candidates...

✅ Found 3 potential version candidates:
   1. 8.32
   2. 2.31
   3. 1.3.2

🧠 Analyzing with Groq AI...

🎯 Most likely version for ls: 8.32
```

## Error Handling

The tool provides detailed error messages for:
- Missing API keys
- Invalid binary files
- Network connectivity issues
- AI provider errors

## Testing

Test the tool with a system binary:
```bash
# Set your API key
export GROQ_API_KEY="your-groq-api-key"

# Test with a common binary
binary-version-analyzer analyze /usr/bin/ls

# Test with different provider
export AI_PROVIDER="openai"
export OPENAI_API_KEY="your-openai-key"
binary-version-analyzer analyze /usr/bin/curl
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add your AI provider implementation in the `providers/` directory
4. Update documentation
5. Submit a pull request

## License

MIT 