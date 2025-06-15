package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"binary-version-analyzer/patterns"
)

var (
	interactive bool
	testString  string
	priority    int
	showDetails bool
)

// patternsCmd represents the patterns command group
var patternsCmd = &cobra.Command{
	Use:   "patterns",
	Short: "Manage and test version detection patterns",
	Long: `The patterns command group provides tools for managing, testing, and 
validating the regex patterns used for version detection.

You can list all patterns, test specific strings against patterns, 
validate pattern correctness, and run interactive testing sessions.`,
	Example: `  # List all patterns
  binary-version-analyzer patterns list

  # Test a string against all patterns
  binary-version-analyzer patterns test "version 1.2.3"

  # Interactive testing mode
  binary-version-analyzer patterns test --interactive

  # Validate all patterns
  binary-version-analyzer patterns validate

  # Show detailed pattern documentation
  binary-version-analyzer patterns docs`,
}

// patternsListCmd lists all available patterns
var patternsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all version detection patterns",
	Long: `List displays all available regex patterns used for version detection,
showing their names, priorities, and basic information.`,
	Example: `  # List all patterns
  binary-version-analyzer patterns list

  # List patterns with details
  binary-version-analyzer patterns list --details

  # List only high priority patterns
  binary-version-analyzer patterns list --priority 1`,
	RunE: runPatternsList,
}

// patternsTestCmd tests strings against patterns
var patternsTestCmd = &cobra.Command{
	Use:   "test [string]",
	Short: "Test a string against all patterns",
	Long: `Test evaluates a given string against all version detection patterns
and shows which patterns match and what they extract.`,
	Example: `  # Test a specific string
  binary-version-analyzer patterns test "version 1.2.3"

  # Interactive testing mode
  binary-version-analyzer patterns test --interactive

  # Test with verbose output
  binary-version-analyzer patterns test "libssl-1.1.1" --verbose`,
	RunE: runPatternsTest,
}

// patternsValidateCmd validates all patterns
var patternsValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate all patterns against their examples",
	Long: `Validate tests all regex patterns against their predefined examples
to ensure they work correctly and extract the expected values.`,
	Example: `  # Validate all patterns
  binary-version-analyzer patterns validate

  # Validate with verbose output
  binary-version-analyzer patterns validate --verbose`,
	RunE: runPatternsValidate,
}

// patternsDocsCmd shows detailed pattern documentation
var patternsDocsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Show detailed pattern documentation",
	Long: `Docs displays comprehensive documentation for all patterns including
descriptions, purposes, examples, and use cases.`,
	Example: `  # Show all pattern documentation
  binary-version-analyzer patterns docs

  # Show documentation for specific priority
  binary-version-analyzer patterns docs --priority 1`,
	RunE: runPatternsDocs,
}

func init() {
	rootCmd.AddCommand(patternsCmd)

	// Add subcommands
	patternsCmd.AddCommand(patternsListCmd)
	patternsCmd.AddCommand(patternsTestCmd)
	patternsCmd.AddCommand(patternsValidateCmd)
	patternsCmd.AddCommand(patternsDocsCmd)

	// Flags for list command
	patternsListCmd.Flags().BoolVar(&showDetails, "details", false, "Show detailed pattern information")
	patternsListCmd.Flags().IntVar(&priority, "priority", 0, "Filter by priority level (1-8)")

	// Flags for test command
	patternsTestCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive testing mode")
	patternsTestCmd.Flags().StringVarP(&testString, "string", "s", "", "String to test (alternative to positional arg)")

	// Flags for docs command
	patternsDocsCmd.Flags().IntVar(&priority, "priority", 0, "Show docs for specific priority level")
}

func runPatternsList(cmd *cobra.Command, args []string) error {
	fmt.Println("📋 Version Pattern Summary")
	fmt.Println(strings.Repeat("=", 30))
	fmt.Println()

	sortedPatterns := patterns.GetPatternsByPriority()

	for i, pattern := range sortedPatterns {
		// Filter by priority if specified
		if priority > 0 && pattern.Priority != priority {
			continue
		}

		fmt.Printf("%2d. %-25s (Priority: %d)\n", i+1, pattern.Name, pattern.Priority)

		if showDetails || verbose {
			fmt.Printf("    %s\n", pattern.Description)
			fmt.Printf("    Pattern: %s\n", pattern.Pattern.String())
			if len(pattern.Examples) > 0 {
				fmt.Printf("    Example: %s", pattern.Examples[0])
				if len(pattern.Expected) > 0 {
					fmt.Printf(" → %s", pattern.Expected[0])
				}
				fmt.Println()
			}
		}
		fmt.Println()
	}

	return nil
}

func runPatternsTest(cmd *cobra.Command, args []string) error {
	if interactive {
		return runInteractiveTest()
	}

	// Get test string from args or flag
	var testStr string
	if len(args) > 0 {
		testStr = strings.Join(args, " ")
	} else if testString != "" {
		testStr = testString
	} else {
		return fmt.Errorf("please provide a string to test or use --interactive mode")
	}

	return testStringAgainstPatterns(testStr)
}

func runPatternsValidate(cmd *cobra.Command, args []string) error {
	fmt.Println("🧪 Validating Version Patterns")
	fmt.Println(strings.Repeat("=", 35))
	fmt.Println()

	if patterns.ValidateAllPatterns() {
		fmt.Println("\n🎉 All patterns are working correctly!")
		return nil
	} else {
		fmt.Println("\n❌ Some patterns have issues!")
		return fmt.Errorf("pattern validation failed")
	}
}

func runPatternsDocs(cmd *cobra.Command, args []string) error {
	if priority > 0 {
		// Show docs for specific priority
		fmt.Printf("📋 Pattern Documentation - Priority %d\n", priority)
		fmt.Println(strings.Repeat("=", 40))
		fmt.Println()

		sortedPatterns := patterns.GetPatternsByPriority()
		found := false

		for _, pattern := range sortedPatterns {
			if pattern.Priority == priority {
				found = true
				printPatternDetails(pattern)
			}
		}

		if !found {
			return fmt.Errorf("no patterns found with priority %d", priority)
		}
	} else {
		// Show all documentation
		patterns.PrintPatternInfo()
	}

	return nil
}

func runInteractiveTest() error {
	fmt.Println("🎮 Interactive Pattern Testing Mode")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Println()
	fmt.Println("Enter strings to test against all patterns.")
	fmt.Println("Commands:")
	fmt.Println("  'quit' or 'exit' - Exit interactive mode")
	fmt.Println("  'help' - Show available commands")
	fmt.Println("  'list' - List all patterns")
	fmt.Println("  'validate' - Validate all patterns")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("🔍 Enter test string: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		switch strings.ToLower(input) {
		case "quit", "exit":
			fmt.Println("👋 Goodbye!")
			return nil
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("  'quit' or 'exit' - Exit")
			fmt.Println("  'help' - Show this help")
			fmt.Println("  'list' - List all patterns")
			fmt.Println("  'validate' - Validate patterns")
			fmt.Println("  Or enter any string to test")
			continue
		case "list":
			runPatternsList(nil, nil)
			continue
		case "validate":
			runPatternsValidate(nil, nil)
			continue
		}

		fmt.Println()
		testStringAgainstPatterns(input)
		fmt.Println(strings.Repeat("-", 50))
	}

	return nil
}

func testStringAgainstPatterns(testStr string) error {
	fmt.Printf("🔍 Testing string: \"%s\"\n", testStr)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	matches := 0
	for _, pattern := range patterns.VersionPatterns {
		result := pattern.Pattern.FindStringSubmatch(testStr)
		if len(result) > 1 {
			matches++
			fmt.Printf("✅ %s (Priority: %d)\n", pattern.Name, pattern.Priority)
			fmt.Printf("   Extracted: \"%s\"\n", result[1])
			if verbose {
				fmt.Printf("   Pattern: %s\n", pattern.Pattern.String())
				fmt.Printf("   Purpose: %s\n", pattern.Purpose)
			}
			fmt.Println()
		}
	}

	if matches == 0 {
		fmt.Println("❌ No patterns matched this string")
		fmt.Println()
		fmt.Println("💡 Suggestions:")
		fmt.Println("   • Make sure the string contains a version number")
		fmt.Println("   • Check if the version format is supported")
		fmt.Println("   • Run 'binary-version-analyzer patterns list' to see all patterns")
	} else {
		fmt.Printf("📊 Total matches: %d\n", matches)
	}

	return nil
}

func printPatternDetails(pattern patterns.VersionPattern) {
	fmt.Printf("### %s (Priority: %d)\n", pattern.Name, pattern.Priority)
	fmt.Printf("**Pattern**: `%s`\n\n", pattern.Pattern.String())
	fmt.Printf("**Description**: %s\n\n", pattern.Description)
	fmt.Printf("**Purpose**: %s\n\n", pattern.Purpose)

	if len(pattern.Examples) > 0 {
		fmt.Println("**Examples**:")
		for i, example := range pattern.Examples {
			expected := ""
			if i < len(pattern.Expected) {
				expected = fmt.Sprintf(" → `%s`", pattern.Expected[i])
			}
			fmt.Printf("- `%s`%s\n", example, expected)
		}
		fmt.Println()
	}

	fmt.Println("---")
	fmt.Println()
}
