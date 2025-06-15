# Cobra CLI Implementation

This document showcases the comprehensive Cobra CLI implementation for the Binary Version Analyzer.

## Overview

The application has been completely refactored to use **Cobra CLI**, providing a modern, professional command-line interface with:

- **Hierarchical commands** with subcommands
- **Rich help system** with examples and descriptions
- **Global and local flags** with validation
- **Interactive modes** for pattern testing
- **Multiple output formats** (text, JSON, YAML)
- **Environment variable integration** with CLI flag overrides
- **Auto-completion support** (built-in Cobra feature)

## Command Structure

```
binary-version-analyzer
├── analyze [binary_path]           # Main binary analysis command
├── patterns                        # Pattern management command group
│   ├── list                       # List all patterns
│   ├── test [string]              # Test patterns against strings
│   ├── validate                   # Validate all patterns
│   └── docs                       # Show pattern documentation
├── completion [shell]             # Generate shell completion scripts
└── help [command]                 # Help about any command
```

## Key Features

### 1. Rich Help System

Every command includes comprehensive help with examples:

```bash
$ binary-version-analyzer --help
Binary Version Analyzer is a sophisticated tool that scans binary files 
to extract and identify their versions using regex pattern matching combined 
with AI-powered analysis.

The tool supports multiple AI providers (Groq, OpenAI) and uses 15 different 
regex patterns to detect version strings in various formats.

Usage:
  binary-version-analyzer [command]

Examples:
  # Analyze a binary file
  binary-version-analyzer analyze /usr/bin/ls

  # Use OpenAI with custom settings
  binary-version-analyzer analyze /usr/bin/curl --provider openai --model gpt-4

  # Show all available patterns
  binary-version-analyzer patterns list

  # Test a pattern interactively
  binary-version-analyzer patterns test --interactive
```

### 2. Global and Local Flags

**Global Flags** (available for all commands):
- `--provider` - AI provider selection (groq, openai)
- `--model` - AI model override
- `--temperature` - AI temperature control (0.0-2.0)
- `--max-tokens` - Maximum response tokens (1-4096)
- `--base-url` - Custom API base URL
- `--timeout` - Request timeout in seconds (1-300)
- `--verbose, -v` - Enable verbose output
- `--config` - Config file path

**Local Flags** (command-specific):
- `analyze` command:
  - `--show-config` - Display AI configuration (default: true)
  - `--show-patterns` - Display pattern information
  - `--output, -o` - Output format (text, json, yaml)
  - `--save` - Save results to file
- `patterns list` command:
  - `--details` - Show detailed pattern information
  - `--priority` - Filter by priority level
- `patterns test` command:
  - `--interactive, -i` - Interactive testing mode
  - `--string, -s` - String to test (alternative to positional arg)

### 3. Analyze Command

The main binary analysis command with rich output options:

```bash
# Basic analysis
binary-version-analyzer analyze /usr/bin/ls

# With custom AI settings
binary-version-analyzer analyze /usr/bin/curl \
  --provider openai \
  --model gpt-4 \
  --temperature 0.2

# Save results to JSON
binary-version-analyzer analyze /usr/bin/git \
  --output json \
  --save results.json \
  --show-patterns

# Verbose analysis with all information
binary-version-analyzer analyze /usr/bin/python3 \
  --verbose \
  --show-config \
  --show-patterns
```

### 4. Pattern Management Commands

#### List Patterns
```bash
# List all patterns
binary-version-analyzer patterns list

# List with detailed information
binary-version-analyzer patterns list --details

# Filter by priority
binary-version-analyzer patterns list --priority 1

# Verbose output
binary-version-analyzer patterns list --verbose
```

#### Test Patterns
```bash
# Test a specific string
binary-version-analyzer patterns test "version 1.2.3"

# Interactive testing mode
binary-version-analyzer patterns test --interactive

# Test with verbose output
binary-version-analyzer patterns test "libssl-1.1.1" --verbose
```

#### Validate Patterns
```bash
# Validate all patterns against their examples
binary-version-analyzer patterns validate

# Validate with verbose output
binary-version-analyzer patterns validate --verbose
```

#### Pattern Documentation
```bash
# Show all pattern documentation
binary-version-analyzer patterns docs

# Show documentation for specific priority
binary-version-analyzer patterns docs --priority 1
```

### 5. Interactive Pattern Testing

The interactive mode provides a REPL-like experience:

```bash
$ binary-version-analyzer patterns test --interactive
🎮 Interactive Pattern Testing Mode
========================================

Enter strings to test against all patterns.
Commands:
  'quit' or 'exit' - Exit interactive mode
  'help' - Show available commands
  'list' - List all patterns
  'validate' - Validate all patterns

🔍 Enter test string: version 1.2.3
🔍 Testing string: "version 1.2.3"
==================================================

✅ Standard Version Declaration (Priority: 1)
   Extracted: "1.2.3"

✅ Semantic Version (Priority: 2)
   Extracted: "1.2.3"

📊 Total matches: 2
--------------------------------------------------
🔍 Enter test string: quit
👋 Goodbye!
```

### 6. Output Formats

The analyze command supports multiple output formats:

#### Text Output (Default)
```
🔍 Analyzing binary: /usr/bin/ls
🤖 Using AI Provider: Groq
📊 Scanning for version candidates...
✅ Found 3 potential version candidates:
   1. 8.32
   2. 2.31
   3. 1.3.2
🎯 Most likely version for ls: 8.32
```

#### JSON Output
```bash
binary-version-analyzer analyze /usr/bin/ls --output json --save results.json
```

```json
{
  "binary_path": "/usr/bin/ls",
  "binary_name": "ls",
  "version": "8.32",
  "candidates": ["8.32", "2.31", "1.3.2"],
  "ai_provider": "Groq",
  "ai_model": "llama-3.1-70b-versatile",
  "pattern_count": 15,
  "timestamp": "2024-01-15T10:30:45Z"
}
```

#### YAML Output
```bash
binary-version-analyzer analyze /usr/bin/ls --output yaml --save results.yaml
```

```yaml
binary_path: /usr/bin/ls
binary_name: ls
version: "8.32"
candidates:
  - "8.32"
  - "2.31"
  - "1.3.2"
ai_provider: Groq
ai_model: llama-3.1-70b-versatile
pattern_count: 15
timestamp: 2024-01-15T10:30:45Z
```

### 7. Configuration Hierarchy

The application supports a flexible configuration hierarchy:

1. **Command-line flags** (highest priority)
2. **Environment variables** (medium priority)
3. **Config file** (lowest priority)
4. **Default values** (fallback)

Example of flag overriding environment variables:
```bash
export AI_PROVIDER="groq"
export AI_MODEL="llama-3.1-70b-versatile"

# This will use OpenAI with GPT-4, overriding the environment
binary-version-analyzer analyze /usr/bin/ls \
  --provider openai \
  --model gpt-4
```

### 8. Error Handling and Validation

The CLI provides comprehensive error handling:

- **File existence validation** for binary paths
- **Flag value validation** (ranges, enums)
- **API key validation** before making requests
- **Pattern validation** with detailed error messages
- **Graceful error messages** with suggestions

### 9. Auto-completion Support

Cobra provides built-in shell completion:

```bash
# Generate completion script for bash
binary-version-analyzer completion bash > /etc/bash_completion.d/binary-version-analyzer

# Generate for zsh
binary-version-analyzer completion zsh > "${fpath[1]}/_binary-version-analyzer"

# Generate for fish
binary-version-analyzer completion fish > ~/.config/fish/completions/binary-version-analyzer.fish
```

## Implementation Benefits

### 1. Professional CLI Experience
- Consistent command structure following Unix conventions
- Rich help system with examples and descriptions
- Proper flag handling with short and long forms
- Hierarchical command organization

### 2. Developer-Friendly
- Easy to extend with new commands
- Clean separation of concerns
- Comprehensive error handling
- Built-in testing support

### 3. User-Friendly
- Interactive modes for exploration
- Multiple output formats for integration
- Flexible configuration options
- Helpful error messages and suggestions

### 4. Production-Ready
- Robust flag validation
- Environment variable integration
- Configuration file support
- Shell completion support

## Code Organization

### cmd/root.go
- Root command definition
- Global flags and configuration
- Initialization logic

### cmd/analyze.go
- Binary analysis command
- Output format handling
- Result saving functionality

### cmd/patterns.go
- Pattern management command group
- Interactive testing mode
- Pattern validation and documentation

### internal/analyzer.go
- Core analysis logic
- Result structures with JSON/YAML tags
- File I/O operations

## Migration from Simple CLI

The migration from a simple argument-based CLI to Cobra provides:

1. **Better UX**: Clear command structure vs. confusing flags
2. **Extensibility**: Easy to add new commands and features
3. **Documentation**: Built-in help system vs. manual documentation
4. **Validation**: Automatic flag validation vs. manual checks
5. **Standards**: Following CLI best practices and conventions

## Future Enhancements

The Cobra structure makes it easy to add:

- **Config management commands** (`config set`, `config get`)
- **Batch processing commands** (`batch analyze`)
- **Export/import commands** for patterns
- **Plugin system** for custom analyzers
- **Daemon mode** for continuous monitoring
- **Web interface** integration commands

This Cobra CLI implementation transforms the Binary Version Analyzer from a simple tool into a professional, extensible command-line application suitable for both interactive use and automation. 