package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rhobs/testplan-viewer/generator"
	"github.com/rhobs/testplan-viewer/models"
	"github.com/spf13/cobra"
)

var (
	inputFile  string
	outputFile string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "testplan-viewer",
		Short: "Generate interactive HTML from RHOBS test plan JSON",
		Long: `testplan-viewer is a CLI tool that generates interactive,
educational HTML documentation from RHOBS test plan JSON files.

The generated HTML includes:
  - Progress tracking with localStorage
  - Interactive collapsible sections
  - Syntax highlighting for code blocks
  - Copy-to-clipboard for commands
  - Search and filter functionality
  - Difficulty-based filtering
  - Learning path visualization`,
		RunE: runGenerate,
	}

	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "rhobs_test_plan_v2.json",
		"Input JSON file path")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "testplan.html",
		"Output HTML file path")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Read JSON file
	fmt.Printf("📖 Reading test plan from: %s\n", inputFile)
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Parse JSON
	fmt.Println("🔍 Parsing test plan JSON...")
	var testPlan models.TestPlan
	if err := json.Unmarshal(data, &testPlan); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate test plan
	fmt.Println("✅ Validating test plan structure...")
	if err := validateTestPlan(&testPlan); err != nil {
		return fmt.Errorf("invalid test plan: %w", err)
	}

	// Generate HTML
	fmt.Printf("🎨 Generating interactive HTML to: %s\n", outputFile)
	if err := generator.GenerateHTML(&testPlan, outputFile); err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	// Print summary
	fmt.Println("\n✨ Success! Generated test plan HTML with:")
	fmt.Printf("   • %d test cases\n", len(testPlan.TestCases))
	fmt.Printf("   • %d concepts explained\n", len(testPlan.Concepts))
	fmt.Printf("   • Learning path: %s\n", testPlan.LearningPath.Description)
	fmt.Printf("   • Target audience: %s\n", testPlan.Metadata.TargetAudience)
	fmt.Printf("   • Estimated time: %s\n", testPlan.Metadata.EstimatedTotalTime)
	fmt.Printf("\n🚀 Open %s in your browser to get started!\n", outputFile)

	return nil
}

func validateTestPlan(tp *models.TestPlan) error {
	if tp.Metadata.DocumentTitle == "" {
		return fmt.Errorf("metadata.document_title is required")
	}

	if len(tp.TestCases) == 0 {
		return fmt.Errorf("at least one test case is required")
	}

	// Validate learning path references existing tests
	for _, testID := range tp.LearningPath.Sequence {
		if _, exists := tp.TestCases[testID]; !exists {
			return fmt.Errorf("learning path references non-existent test: %s", testID)
		}
	}

	// Validate test dependencies
	for testID, testCase := range tp.TestCases {
		for _, depID := range testCase.TestExecution.Dependencies {
			if _, exists := tp.TestCases[depID]; !exists {
				return fmt.Errorf("test %s depends on non-existent test: %s", testID, depID)
			}
		}
	}

	fmt.Println("   ✓ Metadata valid")
	fmt.Println("   ✓ Test cases valid")
	fmt.Println("   ✓ Dependencies valid")
	fmt.Println("   ✓ Learning path valid")

	return nil
}
