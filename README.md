# 🔍 Binary Version Analyzer

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![CLI](https://img.shields.io/badge/CLI-Cobra-brightgreen?style=for-the-badge)](https://github.com/spf13/cobra)
[![AI Powered](https://img.shields.io/badge/AI-Powered-purple?style=for-the-badge)](https://groq.com/)

> **AI-powered binary version detection tool with sophisticated regex pattern matching and modern CLI interface**

## ✨ Features

- 🤖 **AI-Powered Analysis** - Uses Groq/OpenAI to intelligently determine the most likely version
- 🎯 **15 Regex Patterns** - Comprehensive pattern library covering all common version formats
- 🚀 **Modern CLI** - Built with Cobra CLI for professional command-line experience
- 🔧 **Multiple AI Providers** - Support for Groq and OpenAI with easy extensibility
- 📊 **Multiple Output Formats** - Text, JSON, and YAML output options
- 🧪 **Interactive Testing** - Built-in pattern testing and validation tools
- ⚡ **High Performance** - Optimized for large binary files with smart buffering
- 🎮 **Developer Friendly** - Comprehensive debug configurations and documentation

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/srinivasKandukuri/binary-version-analyzer.git
cd binary-version-analyzer

# Build the application
go build -o binary-version-analyzer

# Set your API key
export GROQ_API_KEY="your-groq-api-key"

# Analyze a binary
./binary-version-analyzer analyze /usr/bin/ls
```

### Basic Usage

```bash
# Analyze a binary file
binary-version-analyzer analyze /usr/bin/curl

# Use OpenAI instead of Groq
binary-version-analyzer analyze /usr/bin/git --provider openai

# Save results to JSON
binary-version-analyzer analyze /usr/bin/python3 --output json --save results.json

# Interactive pattern testing
binary-version-analyzer patterns test --interactive
```

## 📋 Command Structure

```
binary-version-analyzer
├── analyze [binary_path]           # Main binary analysis
├── patterns                        # Pattern management
│   ├── list                       # List all patterns
│   ├── test [string]              # Test patterns
│   ├── validate                   # Validate patterns
│   └── docs                       # Pattern documentation
├── completion [shell]             # Shell completion
└── help [command]                 # Help system
```

## 🎯 Example Output

```bash
$ binary-version-analyzer analyze /home/sk/go/bin/boltbrowser

🔧 Initializing configuration...
🚀 Starting analysis of: /home/sk/go/bin/boltbrowser
🔍 Analyzing binary: /home/sk/go/bin/boltbrowser
🤖 Using AI Provider: Groq
🔧 AI Configuration:
   Provider: groq
   Model: gemma2-9b-it
   Temperature: 0.10
   Max Tokens: 50
   Base URL: https://api.groq.com/openai/v1
   Timeout: 30s
   API Key: gsk_E4h1***

🧪 Using 15 version detection patterns
💡 Run 'binary-version-analyzer patterns list' to see all patterns

📊 Scanning for version candidates...
✅ Found 4 potential version candidates:
   1. 1.3.1
   2. 0.0.0-20170904143325
   3. 0.0.4
   4. 0.0.0-20180819125858

🧠 Analyzing with Groq AI...
🎯 Most likely version for boltbrowser: 1.3.1
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `GROQ_API_KEY` | Groq API key | - | Yes (if using Groq) |
| `OPENAI_API_KEY` | OpenAI API key | - | Yes (if using OpenAI) |
| `AI_PROVIDER` | AI provider (`groq`, `openai`) | `groq` | No |
| `AI_MODEL` | Override default model | Provider default | No |
| `AI_TEMPERATURE` | Response randomness (0.0-2.0) | `0.1` | No |
| `AI_MAX_TOKENS` | Maximum response tokens | `50` | No |
| `AI_TIMEOUT` | Request timeout (seconds) | `30` | No |

### Command-Line Flags

```bash
# Global flags (available for all commands)
--provider string     # AI provider (groq, openai)
--model string        # AI model override
--temperature float   # AI temperature (0.0-2.0)
--max-tokens int      # Maximum tokens (1-4096)
--timeout int         # Timeout in seconds
--verbose, -v         # Verbose output

# Analyze command flags
--output, -o string   # Output format (text, json, yaml)
--save string         # Save results to file
--show-config         # Display AI configuration
--show-patterns       # Display pattern information
```

## 🧪 Pattern System

The tool uses **15 sophisticated regex patterns** with priority-based matching:

### Pattern Categories
- **Priority 1**: Most reliable (Standard version declarations, Stable releases)
- **Priority 2**: Common formats (V-prefixed, Semantic versions)
- **Priority 3**: System libraries (GLIBC, Library versions)
- **Priority 4**: Build contexts (Build versions, Package versions)
- **Priority 5**: Development tools (Compiler versions, API versions)
- **Priority 6-8**: Fallback patterns (Date-based, Copyright years)

### Pattern Testing

```bash
# List all patterns
binary-version-analyzer patterns list

# Test a string against patterns
binary-version-analyzer patterns test "version 1.2.3"

# Interactive testing mode
binary-version-analyzer patterns test --interactive

# Validate all patterns
binary-version-analyzer patterns validate
```

## 🏗️ Architecture

```
binary-version-analyzer/
├── cmd/                       # Cobra CLI commands
│   ├── root.go               # Root command & global flags
│   ├── analyze.go            # Binary analysis command
│   └── patterns.go           # Pattern management
├── internal/                  # Core application logic
│   └── analyzer.go           # Binary analyzer & results
├── providers/                 # AI provider implementations
│   ├── interface.go          # Provider interface
│   ├── config.go            # Configuration management
│   ├── groq.go              # Groq implementation
│   ├── openai.go            # OpenAI implementation
│   └── factory.go           # Provider factory
├── patterns/                  # Version detection patterns
│   └── version_patterns.go  # Regex patterns with docs
└── .idea/runConfigurations/  # GoLand debug configs
```

## 🛠️ Development

### Prerequisites
- Go 1.21 or higher
- API key for Groq or OpenAI

### Building from Source

```bash
# Clone the repository
git clone https://github.com/srinivasKandukuri/binary-version-analyzer.git
cd binary-version-analyzer

# Download dependencies
go mod tidy

# Build the application
go build -o binary-version-analyzer

# Run tests
go test ./...
```

### GoLand Debug Setup

The project includes comprehensive GoLand debug configurations:

- **Debug Analyze Command** - Main binary analysis
- **Debug Patterns Test** - Pattern testing and validation
- **Debug OpenAI Analysis** - OpenAI provider testing
- **Debug Large Binary** - Large file handling
- **Debug Help Commands** - CLI system testing

See [`GOLAND_DEBUG_SETUP.md`](GOLAND_DEBUG_SETUP.md) for detailed setup instructions.

## 📚 Documentation

- **[Pattern Documentation](PATTERNS.md)** - Detailed regex pattern guide
- **[Cobra CLI Features](COBRA_CLI_FEATURES.md)** - Complete CLI documentation
- **[GoLand Debug Setup](GOLAND_DEBUG_SETUP.md)** - Development environment setup

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Adding New AI Providers

```go
// 1. Implement the AIProvider interface
type YourProvider struct {
    // implementation
}

func (p *YourProvider) AnalyzeVersions(binaryName string, candidates []string) (string, error) {
    // your implementation
}

// 2. Add to factory.go
const ProviderYour AIProviderType = "your-provider"

// 3. Update CreateProvider method
```

### Adding New Patterns

```go
// Add to patterns/version_patterns.go
{
    Name:        "Your Pattern",
    Pattern:     regexp.MustCompile(`your-regex`),
    Description: "What it matches",
    Purpose:     "Why it's useful",
    Examples:    []string{"example1", "example2"},
    Expected:    []string{"result1", "result2"},
    Priority:    5,
}
```

## 📊 Performance

- **Large File Support** - Handles binaries up to 1GB+ with 1MB buffer
- **Smart Filtering** - Skips binary data and long lines automatically
- **Early Exit** - Stops after finding sufficient version candidates
- **Memory Efficient** - Processes files line-by-line without loading entirely

## 🔒 Security

- **API Key Protection** - Keys are never logged or exposed
- **Input Validation** - All inputs are validated and sanitized
- **Safe Regex** - Patterns are tested against ReDoS attacks
- **No Data Persistence** - No sensitive data is stored locally

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra CLI](https://github.com/spf13/cobra) - Excellent CLI framework
- [Groq](https://groq.com/) - Fast AI inference
- [OpenAI](https://openai.com/) - Advanced AI models
- Go community for excellent tooling and libraries

## 📞 Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/srinivasKandukuri/binary-version-analyzer/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/srinivasKandukuri/binary-version-analyzer/discussions)
- 📧 **Contact**: ksrinivas.cse@gmail.com

---

<div align="center">

**⭐ Star this repository if you find it useful! ⭐**

Made with ❤️ by [Srinivas Kandukuri](https://github.com/srinivasKandukuri)

</div> 