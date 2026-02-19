package parser

import (
	"fmt"

	"github.com/rhobs/testplan-viewer/internal/models"
)

// Validator validates test plan structure and relationships
type Validator struct {
	verbose bool
}

// NewValidator creates a new Validator with verbose output enabled
func NewValidator() *Validator {
	return &Validator{verbose: true}
}

// SetVerbose controls whether validation outputs success messages
func (v *Validator) SetVerbose(verbose bool) {
	v.verbose = verbose
}

// Validate performs comprehensive validation on a test plan
func (v *Validator) Validate(tp *models.TestPlan) error {
	if err := v.validateMetadata(tp); err != nil {
		return err
	}

	if err := v.validateTestCases(tp); err != nil {
		return err
	}

	if err := v.validateLearningPath(tp); err != nil {
		return err
	}

	if err := v.validateDependencies(tp); err != nil {
		return err
	}

	if v.verbose {
		fmt.Println("   ✓ Metadata valid")
		fmt.Println("   ✓ Test cases valid")
		fmt.Println("   ✓ Dependencies valid")
		fmt.Println("   ✓ Learning path valid")
	}

	return nil
}

// validateMetadata checks that required metadata fields are present
func (v *Validator) validateMetadata(tp *models.TestPlan) error {
	if tp.Metadata.DocumentTitle == "" {
		return fmt.Errorf("metadata.document_title is required")
	}
	return nil
}

// validateTestCases ensures at least one test case exists
func (v *Validator) validateTestCases(tp *models.TestPlan) error {
	if len(tp.TestCases) == 0 {
		return fmt.Errorf("at least one test case is required")
	}
	return nil
}

// validateLearningPath ensures all learning path test references exist
func (v *Validator) validateLearningPath(tp *models.TestPlan) error {
	for _, testID := range tp.LearningPath.Sequence {
		if _, exists := tp.TestCases[testID]; !exists {
			return fmt.Errorf("learning path references non-existent test: %s", testID)
		}
	}
	return nil
}

// validateDependencies ensures all test dependencies reference existing tests
func (v *Validator) validateDependencies(tp *models.TestPlan) error {
	for testID, testCase := range tp.TestCases {
		for _, depID := range testCase.TestExecution.Dependencies {
			if _, exists := tp.TestCases[depID]; !exists {
				return fmt.Errorf("test %s depends on non-existent test: %s", testID, depID)
			}
		}
	}
	return nil
}
