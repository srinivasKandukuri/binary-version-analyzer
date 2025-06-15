package patterns

import (
	"fmt"
	"regexp"
	"strings"
)

// VersionPattern represents a regex pattern with metadata
type VersionPattern struct {
	Name        string         // Human-readable name for the pattern
	Pattern     *regexp.Regexp // Compiled regex pattern
	Description string         // Detailed description of what this pattern matches
	Purpose     string         // Why we use this pattern
	Examples    []string       // Example strings that would match
	Expected    []string       // Expected extracted versions from examples
	Priority    int            // Priority level (1=highest, 10=lowest)
}

// VersionPatterns contains all regex patterns for version detection
var VersionPatterns = []VersionPattern{
	{
		Name:        "Standard Version Declaration",
		Pattern:     regexp.MustCompile(`(?i)version\s*[=:]?\s*([\d.]+(?:-[\w.]+)?)`),
		Description: "Matches explicit version declarations with 'version' keyword followed by version number",
		Purpose:     "Captures the most common way software declares its version in binaries",
		Examples: []string{
			"version 1.2.3",
			"Version: 2.4.1",
			"version=3.1.0-beta",
			"VERSION 1.0.0",
		},
		Expected: []string{"1.2.3", "2.4.1", "3.1.0-beta", "1.0.0"},
		Priority: 1,
	},
	{
		Name:        "V-Prefixed Version",
		Pattern:     regexp.MustCompile(`(?i)\bv\s*([\d.]+(?:-[\w.]+)?)\b`),
		Description: "Matches version numbers prefixed with 'v' or 'V'",
		Purpose:     "Common in Git tags and version strings where 'v' prefix is used",
		Examples: []string{
			"v1.2.3",
			"V2.0.1",
			"v3.1.0-alpha",
			"built with v1.18.5",
		},
		Expected: []string{"1.2.3", "2.0.1", "3.1.0-alpha", "1.18.5"},
		Priority: 2,
	},
	{
		Name:        "Stable Release Version",
		Pattern:     regexp.MustCompile(`(?i)([\d.]+(?:-[\w.]+)?)\s*\(stable\)`),
		Description: "Matches version numbers explicitly marked as stable releases",
		Purpose:     "Identifies stable/production versions vs development versions",
		Examples: []string{
			"2.1.3 (stable)",
			"1.0.0-rc1 (stable)",
			"3.2.1(stable)",
		},
		Expected: []string{"2.1.3", "1.0.0-rc1", "3.2.1"},
		Priority: 1,
	},
	{
		Name:        "Release Keyword Version",
		Pattern:     regexp.MustCompile(`(?i)([\d.]+(?:-[\w.]+)?)\s*release`),
		Description: "Matches version numbers followed by 'release' keyword",
		Purpose:     "Captures versions in release notes or release-specific contexts",
		Examples: []string{
			"1.4.2 release",
			"2.0.0-beta release",
			"3.1.1release",
		},
		Expected: []string{"1.4.2", "2.0.0-beta", "3.1.1"},
		Priority: 2,
	},
	{
		Name:        "GLIBC Version",
		Pattern:     regexp.MustCompile(`(?i)glibc[-_]?([\d.]+)`),
		Description: "Matches GNU C Library (glibc) version numbers",
		Purpose:     "Important for determining system compatibility and C library version",
		Examples: []string{
			"glibc-2.31",
			"GLIBC_2.27",
			"glibc2.35",
			"glibc_2.28",
		},
		Expected: []string{"2.31", "2.27", "2.35", "2.28"},
		Priority: 3,
	},
	{
		Name:        "GLIBC Context Version",
		Pattern:     regexp.MustCompile(`(?i)([\d.]+)\s*\(glibc`),
		Description: "Matches version numbers in glibc context (version before glibc reference)",
		Purpose:     "Alternative pattern for glibc version detection in different formats",
		Examples: []string{
			"2.31 (glibc)",
			"2.27 (GLIBC compatible)",
			"2.35(glibc",
		},
		Expected: []string{"2.31", "2.27", "2.35"},
		Priority: 4,
	},
	{
		Name:        "Library Version",
		Pattern:     regexp.MustCompile(`(?i)lib\w*[-_]([\d.]+(?:-[\w.]+)?)`),
		Description: "Matches library versions with lib prefix (libssl, libcrypto, etc.)",
		Purpose:     "Identifies versions of linked libraries and dependencies",
		Examples: []string{
			"libssl-1.1.1",
			"libcrypto_3.0.2",
			"libz-1.2.11",
			"libpthread-2.31",
		},
		Expected: []string{"1.1.1", "3.0.2", "1.2.11", "2.31"},
		Priority: 3,
	},
	{
		Name:        "Stable Keyword Version",
		Pattern:     regexp.MustCompile(`(?i)([\d.]+(?:-[\w.]+)?)\s*stable`),
		Description: "Matches version numbers followed by 'stable' keyword (without parentheses)",
		Purpose:     "Alternative format for stable version identification",
		Examples: []string{
			"1.8.0 stable",
			"2.1.3-rc1 stable",
			"3.0.0stable",
		},
		Expected: []string{"1.8.0", "2.1.3-rc1", "3.0.0"},
		Priority: 2,
	},
	{
		Name:        "Build Version",
		Pattern:     regexp.MustCompile(`(?i)build\s*[#:]?\s*([\d.]+(?:-[\w.]+)?)`),
		Description: "Matches build version numbers with 'build' keyword",
		Purpose:     "Captures build-specific version information and build numbers",
		Examples: []string{
			"build 1.2.3",
			"Build: 2.0.1",
			"build#1.5.0-beta",
			"BUILD 3.1.0",
		},
		Expected: []string{"1.2.3", "2.0.1", "1.5.0-beta", "3.1.0"},
		Priority: 4,
	},
	{
		Name:        "Semantic Version",
		Pattern:     regexp.MustCompile(`(?i)\b([\d]+\.[\d]+\.[\d]+(?:-[\w.]+)?(?:\+[\w.]+)?)\b`),
		Description: "Matches semantic versioning format (MAJOR.MINOR.PATCH with optional pre-release and build metadata)",
		Purpose:     "Captures standard semantic versions used by most modern software",
		Examples: []string{
			"1.2.3",
			"10.15.7",
			"2.1.0-alpha.1",
			"1.0.0+20220101",
			"3.2.1-beta.2+build.123",
		},
		Expected: []string{"1.2.3", "10.15.7", "2.1.0-alpha.1", "1.0.0+20220101", "3.2.1-beta.2+build.123"},
		Priority: 2,
	},
	{
		Name:        "Compiler Version",
		Pattern:     regexp.MustCompile(`(?i)(?:gcc|clang|msvc)(?:[-_\s]+version\s+|[-_\s]*)([\d.]+)`),
		Description: "Matches compiler version numbers (GCC, Clang, MSVC)",
		Purpose:     "Identifies the compiler version used to build the binary",
		Examples: []string{
			"gcc-9.4.0",
			"clang 13.0.1",
			"MSVC_19.29",
			"gcc version 11.2.0",
		},
		Expected: []string{"9.4.0", "13.0.1", "19.29", "11.2.0"},
		Priority: 5,
	},
	{
		Name:        "Package Version",
		Pattern:     regexp.MustCompile(`(?i)(?:pkg|package)[-_\s]*([\d.]+(?:-[\w.]+)?)`),
		Description: "Matches package version numbers with pkg/package prefix",
		Purpose:     "Captures package management system version information",
		Examples: []string{
			"pkg-1.2.3",
			"package 2.0.1",
			"PKG_3.1.0-beta",
		},
		Expected: []string{"1.2.3", "2.0.1", "3.1.0-beta"},
		Priority: 4,
	},
	{
		Name:        "Copyright Year Version",
		Pattern:     regexp.MustCompile(`(?i)copyright.*?(\d{4})`),
		Description: "Matches copyright years which can indicate software age/version era",
		Purpose:     "Provides temporal context when explicit version numbers are not available",
		Examples: []string{
			"Copyright (c) 2023",
			"Copyright 2022 Company",
			"(C) Copyright 2021",
		},
		Expected: []string{"2023", "2022", "2021"},
		Priority: 8,
	},
	{
		Name:        "Date-based Version",
		Pattern:     regexp.MustCompile(`(?i)\b(20\d{2}[.\-_]?(?:0[1-9]|1[0-2])[.\-_]?(?:0[1-9]|[12]\d|3[01]))\b`),
		Description: "Matches date-based version numbers (YYYY.MM.DD, YYYY-MM-DD, YYYYMMDD)",
		Purpose:     "Captures date-based versioning schemes used by some software",
		Examples: []string{
			"2023.12.15",
			"2022-03-21",
			"20231201",
			"2023_11_30",
		},
		Expected: []string{"2023.12.15", "2022-03-21", "20231201", "2023_11_30"},
		Priority: 6,
	},
	{
		Name:        "API Version",
		Pattern:     regexp.MustCompile(`(?i)api.*?([\d.]+)`),
		Description: "Matches API version numbers",
		Purpose:     "Identifies API versions for libraries and services",
		Examples: []string{
			"API v1.2",
			"api_version_2.1",
			"API-3.0",
		},
		Expected: []string{"1.2", "2.1", "3.0"},
		Priority: 5,
	},
}

// GetPatternsByPriority returns patterns sorted by priority (highest first)
func GetPatternsByPriority() []VersionPattern {
	patterns := make([]VersionPattern, len(VersionPatterns))
	copy(patterns, VersionPatterns)

	// Sort by priority (lower number = higher priority)
	for i := 0; i < len(patterns)-1; i++ {
		for j := i + 1; j < len(patterns); j++ {
			if patterns[i].Priority > patterns[j].Priority {
				patterns[i], patterns[j] = patterns[j], patterns[i]
			}
		}
	}

	return patterns
}

// GetCompiledPatterns returns just the compiled regex patterns for scanning
func GetCompiledPatterns() []*regexp.Regexp {
	patterns := make([]*regexp.Regexp, len(VersionPatterns))
	for i, pattern := range VersionPatterns {
		patterns[i] = pattern.Pattern
	}
	return patterns
}

// PrintPatternInfo prints detailed information about all patterns
func PrintPatternInfo() {
	fmt.Println("📋 Version Pattern Documentation")
	fmt.Println(strings.Repeat("=", 50))

	for i, pattern := range GetPatternsByPriority() {
		fmt.Printf("\n%d. %s (Priority: %d)\n", i+1, pattern.Name, pattern.Priority)
		fmt.Printf("   Pattern: %s\n", pattern.Pattern.String())
		fmt.Printf("   Description: %s\n", pattern.Description)
		fmt.Printf("   Purpose: %s\n", pattern.Purpose)

		fmt.Printf("   Examples:\n")
		for j, example := range pattern.Examples {
			expected := ""
			if j < len(pattern.Expected) {
				expected = fmt.Sprintf(" → %s", pattern.Expected[j])
			}
			fmt.Printf("     • %s%s\n", example, expected)
		}
	}
}

// ValidatePattern tests a pattern against its examples
func ValidatePattern(pattern VersionPattern) bool {
	for i, example := range pattern.Examples {
		matches := pattern.Pattern.FindStringSubmatch(example)
		if len(matches) < 2 {
			fmt.Printf("❌ Pattern '%s' failed to match example: %s\n", pattern.Name, example)
			return false
		}

		extracted := matches[1]
		if i < len(pattern.Expected) && extracted != pattern.Expected[i] {
			fmt.Printf("❌ Pattern '%s' extracted '%s' but expected '%s' from: %s\n",
				pattern.Name, extracted, pattern.Expected[i], example)
			return false
		}
	}
	return true
}

// ValidateAllPatterns tests all patterns against their examples
func ValidateAllPatterns() bool {
	fmt.Println("🧪 Validating all version patterns...")
	allValid := true

	for _, pattern := range VersionPatterns {
		if !ValidatePattern(pattern) {
			allValid = false
		} else {
			fmt.Printf("✅ Pattern '%s' validated successfully\n", pattern.Name)
		}
	}

	if allValid {
		fmt.Println("\n🎉 All patterns validated successfully!")
	} else {
		fmt.Println("\n❌ Some patterns failed validation!")
	}

	return allValid
}
