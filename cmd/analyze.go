package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"binary-version-analyzer/internal"
	"binary-version-analyzer/providers"
)

var (
	showConfig   bool
	showPatterns bool
	outputFormat string
	saveResults  string
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze [binary_path]",
	Short: "Analyze a binary file to detect its version",
	Long: `Analyze scans a binary file using regex patterns to find potential version 
strings, then uses AI to determine the most likely version.

The command supports various output formats and can save results to a file.`,
	Example: `  # Basic analysis
  binary-version-analyzer analyze /usr/bin/ls

  # Analyze with custom AI settings
  binary-version-analyzer analyze /usr/bin/curl --provider openai --model gpt-4

  # Show configuration and patterns used
  binary-version-analyzer analyze /usr/bin/python3 --show-config --show-patterns

  # Save results to JSON file
  binary-version-analyzer analyze /usr/bin/git --output json --save results.json`,
	Args: cobra.ExactArgs(1),
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	// Local flags for analyze command
	analyzeCmd.Flags().BoolVar(&showConfig, "show-config", true, "Display AI configuration")
	analyzeCmd.Flags().BoolVar(&showPatterns, "show-patterns", false, "Display pattern information")
	analyzeCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text, json, yaml)")
	analyzeCmd.Flags().StringVar(&saveResults, "save", "", "Save results to file")

	// Mark binary path as required
	analyzeCmd.MarkFlagRequired("binary_path")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	binaryPath := args[0]

	// Check if file exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("❌ Error: File %s does not exist", binaryPath)
	}

	if verbose {
		fmt.Printf("🔍 Starting analysis of: %s\n", binaryPath)
	}

	// Load configuration from environment (with CLI overrides)
	config, err := providers.LoadConfigFromEnv()
	if err != nil {
		return fmt.Errorf("❌ Error loading configuration: %v", err)
	}

	// Create AI provider using factory
	factory := providers.NewAIFactory()
	aiProvider, err := factory.CreateProvider(config)
	if err != nil {
		return fmt.Errorf("❌ Error creating AI provider: %v", err)
	}

	// Create analyzer
	analyzer := internal.NewBinaryAnalyzer(aiProvider)

	// Display information
	fmt.Printf("🔍 Analyzing binary: %s\n", binaryPath)
	fmt.Printf("🤖 Using AI Provider: %s\n", aiProvider.GetProviderName())

	if showConfig {
		providers.PrintConfigInfo(config)
		fmt.Println()
	}

	if showPatterns {
		fmt.Printf("📋 Using %d version detection patterns\n", analyzer.GetPatternCount())
		if verbose {
			fmt.Println("💡 Run 'binary-version-analyzer patterns list' to see all patterns")
		}
		fmt.Println()
	}

	fmt.Println("📊 Scanning for version candidates...")

	// Scan the binary for version candidates
	candidates, err := analyzer.ScanBinary(binaryPath)
	if err != nil {
		return fmt.Errorf("❌ Error scanning binary: %v", err)
	}

	if len(candidates) == 0 {
		fmt.Println("❌ No version candidates found in the binary.")
		fmt.Println("💡 Try running 'binary-version-analyzer patterns list' to see what patterns are used")
		return nil
	}

	fmt.Printf("\n✅ Found %d potential version candidates:\n", len(candidates))
	for i, candidate := range candidates {
		fmt.Printf("   %d. %s\n", i+1, candidate)
	}

	fmt.Printf("\n🧠 Analyzing with %s AI...\n", aiProvider.GetProviderName())

	// Analyze with AI
	binaryName := filepath.Base(binaryPath)
	version, err := analyzer.AnalyzeWithAI(binaryName, candidates)
	if err != nil {
		return fmt.Errorf("❌ Error analyzing with AI: %v", err)
	}

	// Create result
	result := &internal.AnalysisResult{
		BinaryPath:   binaryPath,
		BinaryName:   binaryName,
		Version:      version,
		Candidates:   candidates,
		Provider:     aiProvider.GetProviderName(),
		Model:        config.Model,
		PatternCount: analyzer.GetPatternCount(),
	}

	// Output result
	if err := outputResult(result, outputFormat, saveResults); err != nil {
		return fmt.Errorf("❌ Error outputting result: %v", err)
	}

	fmt.Printf("\n🎯 Most likely version for %s: %s\n", binaryName, version)
	return nil
}

func outputResult(result *internal.AnalysisResult, format, saveFile string) error {
	if saveFile == "" {
		return nil // No saving required
	}

	switch format {
	case "json":
		return result.SaveAsJSON(saveFile)
	case "yaml":
		return result.SaveAsYAML(saveFile)
	case "text":
		return result.SaveAsText(saveFile)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}
