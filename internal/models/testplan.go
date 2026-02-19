package models

// TestPlan represents the complete test plan structure optimized for learning
type TestPlan struct {
	Metadata         Metadata                   `json:"metadata"`
	LearningPath     LearningPath               `json:"learning_path"`
	Prerequisites    Prerequisites              `json:"prerequisites"`
	Concepts         map[string]Concept         `json:"concepts"`
	TestCases        map[string]TestCase        `json:"testcases"`
	ProgressTracking ProgressTracking           `json:"progress_tracking"`
}

// Metadata contains document-level information
type Metadata struct {
	DocumentTitle      string   `json:"document_title"`
	Version            string   `json:"version"`
	Date               string   `json:"date"`
	Status             string   `json:"status"`
	Epic               string   `json:"epic"`
	Purpose            string   `json:"purpose"`
	TargetAudience     string   `json:"target_audience"`
	EstimatedTotalTime string   `json:"estimated_total_time"`
	LearningObjectives []string `json:"learning_objectives"`
}

// LearningPath defines the recommended sequence for learning
type LearningPath struct {
	Description      string              `json:"description"`
	Sequence         []string            `json:"sequence"`
	AlternativePaths map[string][]string `json:"alternative_paths"`
}

// Prerequisites contains required knowledge, access, and environment setup
type Prerequisites struct {
	RequiredKnowledge []RequiredKnowledge `json:"required_knowledge"`
	RequiredAccess    []RequiredAccess    `json:"required_access"`
	EnvironmentSetup  EnvironmentSetup    `json:"environment_setup"`
}

type RequiredKnowledge struct {
	Topic     string   `json:"topic"`
	Level     string   `json:"level"`
	Resources []string `json:"resources"`
}

type RequiredAccess struct {
	Name         string `json:"name"`
	Purpose      string `json:"purpose"`
	HowToGet     string `json:"how_to_get"`
	Verification string `json:"verification"`
}

type EnvironmentSetup struct {
	Variables map[string]string `json:"variables"`
	Tools     []Tool            `json:"tools"`
}

type Tool struct {
	Name    string `json:"name"`
	Install string `json:"install"`
	Verify  string `json:"verify"`
}

// Concept represents a technical concept explained in the guide
type Concept struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	WhyItMatters string   `json:"why_it_matters"`
	DiagramURL   *string  `json:"diagram_url"`
	RelatedTests []string `json:"related_tests"`
}

// TestCase represents a single test with learning objectives
type TestCase struct {
	Metadata       TestMetadata    `json:"metadata"`
	Learning       Learning        `json:"learning"`
	TestExecution  TestExecution   `json:"test_execution"`
	Troubleshooting Troubleshooting `json:"troubleshooting"`
	NextSteps      NextSteps       `json:"next_steps"`
	References     References      `json:"references"`
}

type TestMetadata struct {
	ID                 string           `json:"id"`
	Name               string           `json:"name"`
	Title              string           `json:"title"`
	Category           string           `json:"category"`
	Difficulty         string           `json:"difficulty"`
	EstimatedTime      string           `json:"estimated_time"`
	HandsOnPercentage  int              `json:"hands_on_percentage"`
	SourceTest         string           `json:"source_test,omitempty"`
	SystemImpact       SystemImpact     `json:"system_impact"`
	StateManagement    StateManagement  `json:"state_management"`
	Safety             Safety           `json:"safety"`
}

// SystemImpact describes what changes a test makes to the system
type SystemImpact struct {
	Type              string   `json:"type"`               // read-only, modifies-state, destructive
	Description       string   `json:"description"`
	AffectedResources []string `json:"affected_resources"`
	Reversible        bool     `json:"reversible"`
	Persistence       string   `json:"persistence"`        // temporary, permanent
	RiskLevel         string   `json:"risk_level"`         // low, medium, high, critical
}

// StateManagement defines backup and cleanup requirements
type StateManagement struct {
	RequiresBackup     bool   `json:"requires_backup"`
	BackupTest         string `json:"backup_test,omitempty"`
	RequiresCleanup    bool   `json:"requires_cleanup"`
	CleanupTest        string `json:"cleanup_test,omitempty"`
	RestoresState      bool   `json:"restores_state"`
	StateVerification  string `json:"state_verification,omitempty"`
}

// Safety defines safety constraints and execution requirements
type Safety struct {
	CanRunInProduction   bool   `json:"can_run_in_production"`
	RequiresConfirmation bool   `json:"requires_confirmation"`
	WarningMessage       string `json:"warning_message,omitempty"`
	SafeToRetry          bool   `json:"safe_to_retry"`
	Idempotent           bool   `json:"idempotent"`
}

type Learning struct {
	Objectives              []string               `json:"objectives"`
	ConceptsCovered         []string               `json:"concepts_covered"`
	SkillsGained            []string               `json:"skills_gained"`
	CommonBeginnerMistakes  []CommonMistake        `json:"common_beginner_mistakes"`
}

type CommonMistake struct {
	Mistake       string `json:"mistake"`
	Consequence   string `json:"consequence"`
	HowToAvoid    string `json:"how_to_avoid"`
}

type TestExecution struct {
	Objective      string       `json:"objective"`
	Duration       string       `json:"duration"`
	Dependencies   []string     `json:"dependencies"`
	Prerequisites  []string     `json:"prerequisites"`
	Steps          []Step       `json:"steps"`
	Validation     Validation   `json:"validation"`
}

type Step struct {
	StepNumber              int           `json:"step_number"`
	Title                   string        `json:"title"`
	LearningNote            string        `json:"learning_note,omitempty"`
	Command                 *string       `json:"command"`                            // nullable for manual-only steps
	ExpectedOutput          *string       `json:"expected_output"`                    // nullable for manual-only steps
	WhyThisStep             string        `json:"why_this_step,omitempty"`
	DocumentationReference  string        `json:"documentation_reference,omitempty"`
	ManualSteps             []string      `json:"manual_steps,omitempty"`             // for visual validation
	CommonErrors            []CommonError `json:"common_errors,omitempty"`
}

type CommonError struct {
	Error     string `json:"error"`
	Solution  string `json:"solution"`
	LearnMore string `json:"learn_more"`
}

type Validation struct {
	SuccessCriteria  []string `json:"success_criteria"`
	HowToVerify      string   `json:"how_to_verify"`
	WhatSuccessMeans string   `json:"what_success_means"`
	WhatFailureMeans string   `json:"what_failure_means"`
}

type Troubleshooting struct {
	CommonFailures []CommonFailure `json:"common_failures"`
	GettingHelp    []string        `json:"getting_help"`
}

type CommonFailure struct {
	Symptom       string   `json:"symptom"`
	Cause         string   `json:"cause"`
	DebugSteps    []string `json:"debug_steps"`
	Fix           string   `json:"fix"`
	Prevention    string   `json:"prevention"`
	Documentation string   `json:"documentation,omitempty"`
}

type NextSteps struct {
	OnSuccess        string   `json:"on_success"`
	BonusExploration []string `json:"bonus_exploration"`
}

type References struct {
	Documentation []string `json:"documentation"`
	RelatedCode   []string `json:"related_code"`
	SOPs          []string `json:"sops"`
}

type ProgressTracking struct {
	Description     string   `json:"description"`
	LocalStorageKey string   `json:"local_storage_key"`
	TrackedItems    []string `json:"tracked_items"`
}
