package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"binary-version-analyzer/patterns"
	"binary-version-analyzer/providers"
)

// BinaryAnalyzer handles binary file analysis
type BinaryAnalyzer struct {
	aiProvider providers.AIProvider
	patterns   []*regexp.Regexp
}

// AnalysisResult represents the result of a binary analysis
type AnalysisResult struct {
	BinaryPath   string    `json:"binary_path" yaml:"binary_path"`
	BinaryName   string    `json:"binary_name" yaml:"binary_name"`
	Version      string    `json:"version" yaml:"version"`
	Candidates   []string  `json:"candidates" yaml:"candidates"`
	Provider     string    `json:"ai_provider" yaml:"ai_provider"`
	Model        string    `json:"ai_model" yaml:"ai_model"`
	PatternCount int       `json:"pattern_count" yaml:"pattern_count"`
	Timestamp    time.Time `json:"timestamp" yaml:"timestamp"`
}

// NewBinaryAnalyzer creates a new binary analyzer
func NewBinaryAnalyzer(aiProvider providers.AIProvider) *BinaryAnalyzer {
	return &BinaryAnalyzer{
		aiProvider: aiProvider,
		patterns:   patterns.GetCompiledPatterns(),
	}
}

// GetPatternCount returns the number of patterns being used
func (ba *BinaryAnalyzer) GetPatternCount() int {
	return len(ba.patterns)
}

// ScanBinary scans a binary file for version candidates
func (ba *BinaryAnalyzer) ScanBinary(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %v", path, err)
	}
	defer file.Close()

	var candidates []string
	candidateSet := make(map[string]bool) // To avoid duplicates

	// Use a much larger buffer for binary files and implement custom split function
	const maxBufferSize = 4 * 1024 * 1024 // 4MB buffer
	scanner := bufio.NewScanner(file)

	// Create a custom buffer with maximum size
	buf := make([]byte, maxBufferSize)
	scanner.Buffer(buf, maxBufferSize)

	// Use a custom split function that handles extremely long lines gracefully
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		// Look for newline
		if i := strings.IndexByte(string(data), '\n'); i >= 0 {
			// If line is too long, truncate it
			if i > 2000 {
				i = 2000
			}
			return i + 1, data[0:i], nil
		}

		// If we're at EOF, return whatever we have (up to reasonable limit)
		if atEOF {
			if len(data) > 2000 {
				return len(data), data[0:2000], nil
			}
			return len(data), data, nil
		}

		// If buffer is getting too full, process what we have
		if len(data) > maxBufferSize-1000 {
			// Find a reasonable break point (space, tab, or just truncate)
			breakPoint := 2000
			for i := 1999; i >= 1000; i-- {
				if data[i] == ' ' || data[i] == '\t' {
					breakPoint = i
					break
				}
			}
			return breakPoint, data[0:breakPoint], nil
		}

		// Request more data
		return 0, nil, nil
	})

	lineCount := 0
	maxLines := 50000 // Limit scanning to prevent excessive processing

	for scanner.Scan() && lineCount < maxLines {
		lineCount++
		line := scanner.Text()

		// Skip very long lines (likely binary data)
		if len(line) > 1000 {
			continue
		}

		// Skip binary data lines (lines with non-printable characters)
		if !isPrintable(line) {
			continue
		}

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		for _, pattern := range ba.patterns {
			matches := pattern.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) > 1 {
					version := strings.TrimSpace(match[1])
					if isValidVersion(version) && !candidateSet[version] {
						candidates = append(candidates, version)
						candidateSet[version] = true
					}
				}
			}
		}

		// Early exit if we found enough candidates
		if len(candidates) >= 20 {
			break
		}
	}

	// Handle scanner errors more gracefully
	if err := scanner.Err(); err != nil {
		// If it's still a "token too long" error, try a different approach
		if strings.Contains(err.Error(), "token too long") {
			return ba.scanBinaryChunked(path)
		}
		return nil, fmt.Errorf("error scanning file: %v", err)
	}

	return candidates, nil
}

// scanBinaryChunked is a fallback method for extremely problematic binary files
func (ba *BinaryAnalyzer) scanBinaryChunked(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %v", path, err)
	}
	defer file.Close()

	var candidates []string
	candidateSet := make(map[string]bool)

	// Read file in chunks and process byte by byte
	const chunkSize = 64 * 1024 // 64KB chunks
	buffer := make([]byte, chunkSize)
	var lineBuffer strings.Builder
	processedBytes := 0
	maxBytes := 100 * 1024 * 1024 // Process max 100MB

	for processedBytes < maxBytes {
		n, err := file.Read(buffer)
		if n == 0 {
			break
		}
		processedBytes += n

		for i := 0; i < n; i++ {
			b := buffer[i]

			// If we hit a newline or the line gets too long, process it
			if b == '\n' || lineBuffer.Len() > 1000 {
				line := lineBuffer.String()
				lineBuffer.Reset()

				// Process the line if it looks printable
				if len(line) > 0 && len(line) <= 1000 && isPrintable(line) {
					for _, pattern := range ba.patterns {
						matches := pattern.FindAllStringSubmatch(line, -1)
						for _, match := range matches {
							if len(match) > 1 {
								version := strings.TrimSpace(match[1])
								if isValidVersion(version) && !candidateSet[version] {
									candidates = append(candidates, version)
									candidateSet[version] = true
								}
							}
						}
					}
				}

				// Early exit if we found enough candidates
				if len(candidates) >= 20 {
					return candidates, nil
				}
			} else if b >= 32 && b <= 126 {
				// Only add printable ASCII characters
				lineBuffer.WriteByte(b)
			}
		}

		if err != nil {
			break
		}
	}

	// Process any remaining line
	if lineBuffer.Len() > 0 {
		line := lineBuffer.String()
		if isPrintable(line) {
			for _, pattern := range ba.patterns {
				matches := pattern.FindAllStringSubmatch(line, -1)
				for _, match := range matches {
					if len(match) > 1 {
						version := strings.TrimSpace(match[1])
						if isValidVersion(version) && !candidateSet[version] {
							candidates = append(candidates, version)
							candidateSet[version] = true
						}
					}
				}
			}
		}
	}

	return candidates, nil
}

// AnalyzeWithAI uses AI to determine the most likely version from candidates
func (ba *BinaryAnalyzer) AnalyzeWithAI(binaryName string, candidates []string) (string, error) {
	return ba.aiProvider.AnalyzeVersions(binaryName, candidates)
}

// Helper functions
func isPrintable(s string) bool {
	// Quick check for empty strings
	if len(s) == 0 {
		return false
	}

	// For performance, only check first 200 characters for very long strings
	checkLen := len(s)
	if checkLen > 200 {
		checkLen = 200
	}

	nonPrintableCount := 0
	for i, r := range s {
		if i >= checkLen {
			break
		}

		// Allow common whitespace characters
		if r == '\t' || r == '\n' || r == '\r' {
			continue
		}

		// Count non-printable characters
		if r < 32 || r > 126 {
			nonPrintableCount++
			// If more than 10% of checked characters are non-printable, consider it binary
			if nonPrintableCount > checkLen/10 {
				return false
			}
		}
	}

	return true
}

func isValidVersion(version string) bool {
	// Basic validation for version strings
	if len(version) == 0 || len(version) > 20 {
		return false
	}

	// Should contain at least one digit and one dot
	hasDigit := false
	hasDot := false
	for _, r := range version {
		if r >= '0' && r <= '9' {
			hasDigit = true
		} else if r == '.' {
			hasDot = true
		} else if r != '-' && r != '_' {
			// Allow only digits, dots, hyphens, and underscores
			return false
		}
	}

	return hasDigit && hasDot
}

// SaveAsJSON saves the analysis result as JSON
func (ar *AnalysisResult) SaveAsJSON(filename string) error {
	ar.Timestamp = time.Now()

	data, err := json.MarshalIndent(ar, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling to JSON: %v", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON file: %v", err)
	}

	fmt.Printf("💾 Results saved to %s\n", filename)
	return nil
}

// SaveAsYAML saves the analysis result as YAML
func (ar *AnalysisResult) SaveAsYAML(filename string) error {
	ar.Timestamp = time.Now()

	data, err := yaml.Marshal(ar)
	if err != nil {
		return fmt.Errorf("error marshaling to YAML: %v", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing YAML file: %v", err)
	}

	fmt.Printf("💾 Results saved to %s\n", filename)
	return nil
}

// SaveAsText saves the analysis result as plain text
func (ar *AnalysisResult) SaveAsText(filename string) error {
	ar.Timestamp = time.Now()

	var sb strings.Builder
	sb.WriteString("Binary Version Analysis Report\n")
	sb.WriteString("==============================\n\n")
	sb.WriteString(fmt.Sprintf("Binary Path: %s\n", ar.BinaryPath))
	sb.WriteString(fmt.Sprintf("Binary Name: %s\n", ar.BinaryName))
	sb.WriteString(fmt.Sprintf("Detected Version: %s\n", ar.Version))
	sb.WriteString(fmt.Sprintf("AI Provider: %s\n", ar.Provider))
	sb.WriteString(fmt.Sprintf("AI Model: %s\n", ar.Model))
	sb.WriteString(fmt.Sprintf("Patterns Used: %d\n", ar.PatternCount))
	sb.WriteString(fmt.Sprintf("Analysis Time: %s\n\n", ar.Timestamp.Format(time.RFC3339)))

	sb.WriteString("Version Candidates Found:\n")
	for i, candidate := range ar.Candidates {
		sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, candidate))
	}

	err := os.WriteFile(filename, []byte(sb.String()), 0644)
	if err != nil {
		return fmt.Errorf("error writing text file: %v", err)
	}

	fmt.Printf("💾 Results saved to %s\n", filename)
	return nil
}
