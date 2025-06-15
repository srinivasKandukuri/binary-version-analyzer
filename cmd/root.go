package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	aiProvider    string
	aiModel       string
	aiTemperature float64
	aiMaxTokens   int
	aiBaseURL     string
	aiTimeout     int
	verbose       bool
	configFile    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "binary-version-analyzer",
	Short: "AI-powered binary version detection tool",
	Long: `Binary Version Analyzer is a sophisticated tool that scans binary files 
to extract and identify their versions using regex pattern matching combined 
with AI-powered analysis.

The tool supports multiple AI providers (Groq, OpenAI) and uses 15 different 
regex patterns to detect version strings in various formats.`,
	Example: `  # Analyze a binary file
  binary-version-analyzer analyze /usr/bin/ls

  # Use OpenAI with custom settings
  binary-version-analyzer analyze /usr/bin/curl --provider openai --model gpt-4

  # Show all available patterns
  binary-version-analyzer patterns list

  # Test a pattern interactively
  binary-version-analyzer patterns test --interactive`,
	Version: "1.0.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&aiProvider, "provider", "", "AI provider to use (groq, openai)")
	rootCmd.PersistentFlags().StringVar(&aiModel, "model", "", "AI model to use (overrides provider default)")
	rootCmd.PersistentFlags().Float64Var(&aiTemperature, "temperature", -1, "AI temperature (0.0-2.0)")
	rootCmd.PersistentFlags().IntVar(&aiMaxTokens, "max-tokens", -1, "Maximum AI response tokens (1-4096)")
	rootCmd.PersistentFlags().StringVar(&aiBaseURL, "base-url", "", "Custom AI API base URL")
	rootCmd.PersistentFlags().IntVar(&aiTimeout, "timeout", -1, "Request timeout in seconds (1-300)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file (default is $HOME/.binary-version-analyzer.yaml)")

	// Mark some flags as mutually exclusive or required by specific commands
	rootCmd.MarkFlagsMutuallyExclusive("config", "provider")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if verbose {
		fmt.Println("🔧 Initializing configuration...")
	}

	// Override environment variables with command line flags if provided
	if aiProvider != "" {
		os.Setenv("AI_PROVIDER", aiProvider)
	}
	if aiModel != "" {
		os.Setenv("AI_MODEL", aiModel)
	}
	if aiTemperature >= 0 {
		os.Setenv("AI_TEMPERATURE", fmt.Sprintf("%.2f", aiTemperature))
	}
	if aiMaxTokens > 0 {
		os.Setenv("AI_MAX_TOKENS", fmt.Sprintf("%d", aiMaxTokens))
	}
	if aiBaseURL != "" {
		os.Setenv("AI_BASE_URL", aiBaseURL)
	}
	if aiTimeout > 0 {
		os.Setenv("AI_TIMEOUT", fmt.Sprintf("%d", aiTimeout))
	}
}
