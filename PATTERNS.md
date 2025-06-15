# Version Detection Patterns Documentation

This document provides detailed information about all regex patterns used for version detection in the Binary Version Analyzer.

## Overview

The Binary Version Analyzer uses 15 different regex patterns to identify version strings in binary files. Each pattern is designed to capture specific version formats commonly found in software binaries.

## Pattern Priority System

Patterns are assigned priority levels (1-10, where 1 is highest priority):
- **Priority 1**: Most reliable and specific patterns
- **Priority 2-3**: Common version formats
- **Priority 4-5**: Specialized patterns for specific contexts
- **Priority 6-8**: Fallback patterns and temporal indicators
- **Priority 9-10**: Reserved for future patterns

## Pattern Details

### 1. Standard Version Declaration (Priority: 1)
**Pattern**: `(?i)version\s*[=:]?\s*([\d.]+(?:-[\w.]+)?)`

**Description**: Matches explicit version declarations with 'version' keyword followed by version number

**Purpose**: Captures the most common way software declares its version in binaries

**Examples**:
- `version 1.2.3` → `1.2.3`
- `Version: 2.4.1` → `2.4.1`
- `version=3.1.0-beta` → `3.1.0-beta`
- `VERSION 1.0.0` → `1.0.0`

**Use Cases**: 
- Software with explicit version strings
- Configuration files embedded in binaries
- Debug information

---

### 2. Stable Release Version (Priority: 1)
**Pattern**: `(?i)([\d.]+(?:-[\w.]+)?)\s*\(stable\)`

**Description**: Matches version numbers explicitly marked as stable releases

**Purpose**: Identifies stable/production versions vs development versions

**Examples**:
- `2.1.3 (stable)` → `2.1.3`
- `1.0.0-rc1 (stable)` → `1.0.0-rc1`
- `3.2.1(stable)` → `3.2.1`

**Use Cases**:
- Production software releases
- Package management systems
- Release documentation

---

### 3. V-Prefixed Version (Priority: 2)
**Pattern**: `(?i)\bv\s*([\d.]+(?:-[\w.]+)?)\b`

**Description**: Matches version numbers prefixed with 'v' or 'V'

**Purpose**: Common in Git tags and version strings where 'v' prefix is used

**Examples**:
- `v1.2.3` → `1.2.3`
- `V2.0.1` → `2.0.1`
- `v3.1.0-alpha` → `3.1.0-alpha`
- `built with v1.18.5` → `1.18.5`

**Use Cases**:
- Git-based version tags
- Build system outputs
- Compiler version strings

---

### 4. Release Keyword Version (Priority: 2)
**Pattern**: `(?i)([\d.]+(?:-[\w.]+)?)\s*release`

**Description**: Matches version numbers followed by 'release' keyword

**Purpose**: Captures versions in release notes or release-specific contexts

**Examples**:
- `1.4.2 release` → `1.4.2`
- `2.0.0-beta release` → `2.0.0-beta`
- `3.1.1release` → `3.1.1`

**Use Cases**:
- Release announcements
- Changelog entries
- Distribution packages

---

### 5. Semantic Version (Priority: 2)
**Pattern**: `(?i)\b([\d]+\.[\d]+\.[\d]+(?:-[\w.]+)?(?:\+[\w.]+)?)\b`

**Description**: Matches semantic versioning format (MAJOR.MINOR.PATCH with optional pre-release and build metadata)

**Purpose**: Captures standard semantic versions used by most modern software

**Examples**:
- `1.2.3` → `1.2.3`
- `10.15.7` → `10.15.7`
- `2.1.0-alpha.1` → `2.1.0-alpha.1`
- `1.0.0+20220101` → `1.0.0+20220101`
- `3.2.1-beta.2+build.123` → `3.2.1-beta.2+build.123`

**Use Cases**:
- Modern software following SemVer
- NPM packages
- Docker images
- API versions

---

### 6. Stable Keyword Version (Priority: 2)
**Pattern**: `(?i)([\d.]+(?:-[\w.]+)?)\s*stable`

**Description**: Matches version numbers followed by 'stable' keyword (without parentheses)

**Purpose**: Alternative format for stable version identification

**Examples**:
- `1.8.0 stable` → `1.8.0`
- `2.1.3-rc1 stable` → `2.1.3-rc1`
- `3.0.0stable` → `3.0.0`

**Use Cases**:
- Linux distributions
- Package repositories
- Software documentation

---

### 7. GLIBC Version (Priority: 3)
**Pattern**: `(?i)glibc[-_]?([\d.]+)`

**Description**: Matches GNU C Library (glibc) version numbers

**Purpose**: Important for determining system compatibility and C library version

**Examples**:
- `glibc-2.31` → `2.31`
- `GLIBC_2.27` → `2.27`
- `glibc2.35` → `2.35`
- `glibc_2.28` → `2.28`

**Use Cases**:
- System libraries
- Binary compatibility checking
- Linux distribution identification

---

### 8. Library Version (Priority: 3)
**Pattern**: `(?i)lib\w*[-_]?([\d.]+(?:-[\w.]+)?)`

**Description**: Matches library versions with lib prefix (libssl, libcrypto, etc.)

**Purpose**: Identifies versions of linked libraries and dependencies

**Examples**:
- `libssl-1.1.1` → `1.1.1`
- `libcrypto_3.0.2` → `3.0.2`
- `libz1.2.11` → `1.2.11`
- `libpthread-2.31` → `2.31`

**Use Cases**:
- Shared libraries
- Security analysis
- Dependency tracking

---

### 9. GLIBC Context Version (Priority: 4)
**Pattern**: `(?i)([\d.]+)\s*\(glibc`

**Description**: Matches version numbers in glibc context (version before glibc reference)

**Purpose**: Alternative pattern for glibc version detection in different formats

**Examples**:
- `2.31 (glibc)` → `2.31`
- `2.27 (GLIBC compatible)` → `2.27`
- `2.35(glibc` → `2.35`

**Use Cases**:
- Compatibility strings
- System information
- Build configurations

---

### 10. Build Version (Priority: 4)
**Pattern**: `(?i)build\s*[#:]?\s*([\d.]+(?:-[\w.]+)?)`

**Description**: Matches build version numbers with 'build' keyword

**Purpose**: Captures build-specific version information and build numbers

**Examples**:
- `build 1.2.3` → `1.2.3`
- `Build: 2.0.1` → `2.0.1`
- `build#1.5.0-beta` → `1.5.0-beta`
- `BUILD 3.1.0` → `3.1.0`

**Use Cases**:
- CI/CD systems
- Build artifacts
- Development versions

---

### 11. Package Version (Priority: 4)
**Pattern**: `(?i)(?:pkg|package)[-_\s]*([\d.]+(?:-[\w.]+)?)`

**Description**: Matches package version numbers with pkg/package prefix

**Purpose**: Captures package management system version information

**Examples**:
- `pkg-1.2.3` → `1.2.3`
- `package 2.0.1` → `2.0.1`
- `PKG_3.1.0-beta` → `3.1.0-beta`

**Use Cases**:
- Package managers
- Distribution systems
- Software repositories

---

### 12. Compiler Version (Priority: 5)
**Pattern**: `(?i)(?:gcc|clang|msvc)[-_\s]*([\d.]+)`

**Description**: Matches compiler version numbers (GCC, Clang, MSVC)

**Purpose**: Identifies the compiler version used to build the binary

**Examples**:
- `gcc-9.4.0` → `9.4.0`
- `clang 13.0.1` → `13.0.1`
- `MSVC_19.29` → `19.29`
- `gcc version 11.2.0` → `11.2.0`

**Use Cases**:
- Build environment analysis
- Compatibility checking
- Security vulnerability assessment

---

### 13. API Version (Priority: 5)
**Pattern**: `(?i)api[-_\s]*v?([\d.]+)`

**Description**: Matches API version numbers

**Purpose**: Identifies API versions for libraries and services

**Examples**:
- `API v1.2` → `1.2`
- `api_version_2.1` → `2.1`
- `API-3.0` → `3.0`

**Use Cases**:
- Web services
- Library APIs
- Protocol versions

---

### 14. Date-based Version (Priority: 6)
**Pattern**: `(?i)\b(20\d{2}[.\-_]?(?:0[1-9]|1[0-2])[.\-_]?(?:0[1-9]|[12]\d|3[01]))\b`

**Description**: Matches date-based version numbers (YYYY.MM.DD, YYYY-MM-DD, YYYYMMDD)

**Purpose**: Captures date-based versioning schemes used by some software

**Examples**:
- `2023.12.15` → `2023.12.15`
- `2022-03-21` → `2022-03-21`
- `20231201` → `20231201`
- `2023_11_30` → `2023_11_30`

**Use Cases**:
- Snapshot releases
- Daily builds
- Time-based versioning

---

### 15. Copyright Year Version (Priority: 8)
**Pattern**: `(?i)copyright.*?(\d{4})`

**Description**: Matches copyright years which can indicate software age/version era

**Purpose**: Provides temporal context when explicit version numbers are not available

**Examples**:
- `Copyright (c) 2023` → `2023`
- `Copyright 2022 Company` → `2022`
- `(C) Copyright 2021` → `2021`

**Use Cases**:
- Fallback version indication
- Software age estimation
- Legal information

## Testing Patterns

### Using the Cobra CLI
```bash
# List all patterns
binary-version-analyzer patterns list

# List patterns with details
binary-version-analyzer patterns list --details

# Show patterns by priority
binary-version-analyzer patterns list --priority 1

# Test a specific string
binary-version-analyzer patterns test "version 1.2.3"

# Interactive testing mode
binary-version-analyzer patterns test --interactive

# Validate all patterns
binary-version-analyzer patterns validate

# Show detailed pattern documentation
binary-version-analyzer patterns docs
```

## Adding New Patterns

To add a new pattern:

1. Edit `patterns/version_patterns.go`
2. Add your pattern to the `VersionPatterns` slice:
```go
{
    Name:        "Your Pattern Name",
    Pattern:     regexp.MustCompile(`your-regex-pattern`),
    Description: "What this pattern matches",
    Purpose:     "Why we use this pattern",
    Examples:    []string{"example1", "example2"},
    Expected:    []string{"result1", "result2"},
    Priority:    5, // Choose appropriate priority
},
```
3. Test your pattern:
```bash
binary-version-analyzer patterns validate
```

## Pattern Performance

Patterns are ordered by priority for optimal performance:
- High-priority patterns are checked first
- More specific patterns take precedence
- Generic patterns serve as fallbacks

## Common Issues

### False Positives
- **Issue**: Pattern matches unrelated numbers
- **Solution**: Use word boundaries (`\b`) and context

### Missing Matches
- **Issue**: Valid versions not detected
- **Solution**: Add specific patterns or broaden existing ones

### Performance
- **Issue**: Too many patterns slow down scanning
- **Solution**: Optimize regex and use appropriate priorities

## Best Practices

1. **Specificity**: Make patterns as specific as possible
2. **Context**: Include surrounding context when possible
3. **Testing**: Always test patterns with real examples
4. **Documentation**: Provide clear examples and use cases
5. **Priority**: Assign appropriate priority levels 