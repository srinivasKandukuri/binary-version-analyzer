# GoLand Debug Configuration Guide

This guide provides comprehensive debug configurations for the Binary Version Analyzer Cobra CLI application in GoLand.

## 🚀 Quick Setup

1. **Open the project** in GoLand
2. **Set your API keys** in the environment variables (see configurations below)
3. **Select a debug configuration** from the dropdown
4. **Set breakpoints** in your code
5. **Click Debug** (Shift+F9)

## 📋 Available Debug Configurations

### 1. **Debug Analyze Command** 
**Purpose**: Debug the main binary analysis functionality
- **Command**: `analyze /usr/bin/ls --verbose --show-patterns`
- **Use Case**: Test binary scanning and AI analysis
- **Environment**: Groq provider with default settings
- **Breakpoint Suggestions**:
  - `cmd/analyze.go:runAnalyze()` - Main analysis entry point
  - `internal/analyzer.go:ScanBinary()` - Binary scanning logic
  - `internal/analyzer.go:AnalyzeWithAI()` - AI analysis call

### 2. **Debug Patterns Test**
**Purpose**: Debug pattern testing and interactive mode
- **Command**: `patterns test --interactive --verbose`
- **Use Case**: Test regex patterns and interactive CLI
- **Breakpoint Suggestions**:
  - `cmd/patterns.go:runInteractiveTest()` - Interactive mode logic
  - `cmd/patterns.go:testStringAgainstPatterns()` - Pattern matching
  - `patterns/version_patterns.go:ValidatePattern()` - Pattern validation

### 3. **Debug Patterns Validate**
**Purpose**: Debug pattern validation system
- **Command**: `patterns validate --verbose`
- **Use Case**: Test pattern validation and regex correctness
- **Breakpoint Suggestions**:
  - `cmd/patterns.go:runPatternsValidate()` - Validation entry point
  - `patterns/version_patterns.go:ValidateAllPatterns()` - Validation logic
  - `patterns/version_patterns.go:ValidatePattern()` - Individual pattern tests

### 4. **Debug OpenAI Analysis**
**Purpose**: Debug OpenAI provider integration
- **Command**: `analyze /usr/bin/curl --provider openai --model gpt-4 --temperature 0.2 --verbose`
- **Use Case**: Test OpenAI API integration and different AI models
- **Environment**: OpenAI provider with GPT-4
- **Breakpoint Suggestions**:
  - `providers/openai.go:AnalyzeVersions()` - OpenAI API call
  - `providers/factory.go:CreateProvider()` - Provider creation
  - `providers/config.go:LoadConfigFromEnv()` - Configuration loading

### 5. **Debug Help Commands**
**Purpose**: Debug CLI help system and command structure
- **Command**: `--help`
- **Use Case**: Test Cobra CLI help generation and command structure
- **Breakpoint Suggestions**:
  - `cmd/root.go:Execute()` - Root command execution
  - `cmd/root.go:initConfig()` - Configuration initialization

### 6. **Debug Large Binary**
**Purpose**: Debug large binary file handling (like trivy-scanner)
- **Command**: `analyze /home/sk/Downloads/trivy-scanner --verbose --show-patterns`
- **Use Case**: Test improved binary scanning with large files
- **Environment**: Groq provider with default settings
- **Breakpoint Suggestions**:
  - `internal/analyzer.go:ScanBinary()` - Binary scanning with buffer management
  - `internal/analyzer.go:isPrintable()` - Improved printable character detection
  - Line processing loop to inspect buffer handling

## 🔧 Environment Variables Setup

### Required API Keys
Before debugging, set your API keys in the configurations:

#### For Groq (Default)
```
GROQ_API_KEY=your-actual-groq-api-key-here
```

#### For OpenAI
```
OPENAI_API_KEY=your-actual-openai-api-key-here
```

### Configuration Variables
Each debug configuration includes these environment variables:

| Variable | Purpose | Default Value |
|----------|---------|---------------|
| `AI_PROVIDER` | AI provider selection | `groq` |
| `AI_MODEL` | Model override | Provider default |
| `AI_TEMPERATURE` | Response randomness | `0.1` |
| `AI_MAX_TOKENS` | Maximum response tokens | `50` |
| `AI_TIMEOUT` | Request timeout (seconds) | `30` |

## 🎯 Debugging Strategies

### 1. **Binary Analysis Debugging**
```go
// Set breakpoints in these key functions:
func (ba *BinaryAnalyzer) ScanBinary(path string) ([]string, error) {
    // Debug binary file reading and pattern matching
}

func (ba *BinaryAnalyzer) AnalyzeWithAI(binaryName string, candidates []string) (string, error) {
    // Debug AI provider calls and response handling
}
```

### 2. **Pattern System Debugging**
```go
// Debug pattern matching:
func testStringAgainstPatterns(testStr string) error {
    // Set breakpoint here to inspect pattern matching
    for _, pattern := range patterns.VersionPatterns {
        result := pattern.Pattern.FindStringSubmatch(testStr)
        // Inspect 'result' variable
    }
}
```

### 3. **CLI Command Debugging**
```go
// Debug Cobra command execution:
func runAnalyze(cmd *cobra.Command, args []string) error {
    // Set breakpoint to inspect command arguments and flags
    binaryPath := args[0]
    // Debug flag parsing and validation
}
```

### 4. **AI Provider Debugging**
```go
// Debug API calls:
func (g *GroqProvider) AnalyzeVersions(binaryName string, candidates []string) (string, error) {
    // Set breakpoint to inspect API request/response
    resp, err := g.client.CreateChatCompletion(ctx, req)
    // Inspect 'resp' and 'err'
}
```

## 🔍 Common Debugging Scenarios

### Scenario 1: Pattern Not Matching
1. Use **Debug Patterns Test** configuration
2. Set breakpoint in `testStringAgainstPatterns()`
3. Inspect `pattern.Pattern.FindStringSubmatch(testStr)` result
4. Check regex pattern and test string compatibility

### Scenario 2: AI API Issues
1. Use **Debug Analyze Command** or **Debug OpenAI Analysis**
2. Set breakpoint in provider's `AnalyzeVersions()` method
3. Inspect API request payload and response
4. Check environment variables and API key validity

### Scenario 3: CLI Flag Issues
1. Use **Debug Help Commands** configuration
2. Set breakpoint in `cmd/root.go:initConfig()`
3. Inspect flag parsing and environment variable override logic
4. Check `cobra.Command` flag definitions

### Scenario 4: Binary Scanning Issues
1. Use **Debug Analyze Command** or **Debug Large Binary** configuration
2. Set breakpoint in `internal/analyzer.go:ScanBinary()`
3. Inspect file reading and line processing
4. Check `isPrintable()` and `isValidVersion()` functions

### Scenario 5: Large Binary File Issues
1. Use **Debug Large Binary** configuration
2. Set breakpoint at buffer allocation: `scanner.Buffer(buf, maxBufferSize)`
3. Monitor `lineCount` and `len(candidates)` variables
4. Check line length filtering: `if len(line) > 1000`
5. Inspect early exit conditions for performance

## 🛠️ Advanced Debugging Tips

### 1. **Custom Debug Configurations**
Create additional configurations for specific test cases:

```xml
<!-- Example: Debug with custom binary -->
<parameters value="analyze /path/to/your/binary --verbose --show-config" />
```

### 2. **Environment File Support**
Create a `.env` file in your project root:
```bash
GROQ_API_KEY=your-groq-key
OPENAI_API_KEY=your-openai-key
AI_TEMPERATURE=0.2
AI_MAX_TOKENS=100
```

### 3. **Conditional Breakpoints**
Set conditional breakpoints for specific scenarios:
```go
// Break only when specific patterns match
if pattern.Name == "Semantic Version" && len(result) > 1 {
    // Breakpoint here
}
```

### 4. **Watch Variables**
Add these expressions to your watch list:
- `len(candidates)` - Number of version candidates found
- `config.Provider` - Current AI provider
- `pattern.Priority` - Current pattern priority
- `result[1]` - Extracted version string

## 📝 Debugging Checklist

Before debugging, ensure:
- [ ] API keys are set correctly
- [ ] Binary file exists and is readable
- [ ] Go modules are downloaded (`go mod tidy`)
- [ ] Project builds successfully (`go build`)
- [ ] Breakpoints are set in relevant functions
- [ ] Watch variables are configured

## 🚨 Common Issues and Solutions

### Issue 1: "API key not found"
**Solution**: Check environment variables in debug configuration

### Issue 2: "Binary file not found"
**Solution**: Use absolute paths or ensure file exists: `/usr/bin/ls`

### Issue 3: "No patterns match"
**Solution**: Debug pattern validation first with **Debug Patterns Validate**

### Issue 4: "AI request timeout"
**Solution**: Increase `AI_TIMEOUT` environment variable

### Issue 5: "bufio.Scanner: token too long"
**Solution**: Use the improved binary scanning logic (fixed in latest version)
- The scanner now uses a 1MB buffer
- Long lines (>1000 chars) are automatically skipped
- Processing is limited to 50,000 lines for performance

### Issue 6: "Large binary files take too long"
**Solution**: The scanner now includes performance optimizations:
- Early exit after finding 20 version candidates
- Line length filtering to skip binary data
- Improved printable character detection

## 🎓 Learning Resources

### Key Files to Understand
1. **`cmd/root.go`** - Cobra CLI setup and global configuration
2. **`cmd/analyze.go`** - Main analysis command logic
3. **`cmd/patterns.go`** - Pattern management commands
4. **`internal/analyzer.go`** - Core binary analysis logic
5. **`providers/`** - AI provider implementations
6. **`patterns/version_patterns.go`** - Regex pattern definitions

### Debugging Flow
1. **CLI Entry** → `main.go:main()` → `cmd.Execute()`
2. **Command Parsing** → `cmd/root.go` → Specific command file
3. **Business Logic** → `internal/analyzer.go` or `patterns/`
4. **AI Integration** → `providers/` directory
5. **Output** → Console or file output

This comprehensive debug setup will help you efficiently develop and troubleshoot the Binary Version Analyzer application in GoLand! 🚀 