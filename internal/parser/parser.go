package parser

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rhobs/testplan-viewer/internal/models"
)

// Parser handles test plan parsing and validation
type Parser struct {
	validator *Validator
}

// New creates a new Parser instance with default validator
func New() *Parser {
	return &Parser{
		validator: NewValidator(),
	}
}

// NewWithValidator creates a Parser with a custom validator
func NewWithValidator(v *Validator) *Parser {
	return &Parser{
		validator: v,
	}
}

// ParseFile reads and parses a test plan JSON file
func (p *Parser) ParseFile(filePath string) (*models.TestPlan, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return p.Parse(data)
}

// Parse parses test plan JSON data and validates it
func (p *Parser) Parse(data []byte) (*models.TestPlan, error) {
	var testPlan models.TestPlan
	if err := json.Unmarshal(data, &testPlan); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if err := p.validator.Validate(&testPlan); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &testPlan, nil
}

// SetVerbose controls whether the validator outputs success messages
func (p *Parser) SetVerbose(verbose bool) {
	p.validator.SetVerbose(verbose)
}
