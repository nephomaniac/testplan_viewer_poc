# Test Plan Tools Repository - Comprehensive Improvements

**Date**: 2026-02-18
**Context**: Improvements identified while creating RHOBS-next test plan
**Workflow**: testplan-generator → testplan-educator → HTML generation → testplan-executor review

---

## Executive Summary

During the creation of the RHOBS-next test plan, we identified **15 improvements** across skills, JSON schema, parsers, HTML template, and tooling. These improvements fall into three categories:

- 🔴 **Critical** (3): Block test plan generation/execution
- 🟡 **High Priority** (7): Significantly improve quality and safety
- 🟢 **Medium Priority** (5): Quality of life and user experience

---

## Critical Improvements (Must Fix)

### 1. Testplan-Generator Skill Must Follow Go Data Model

**Problem**: The testplan-generator skill creates JSON structures that don't match `/Users/maclark/sandbox/testplan_tools_poc/internal/models/testplan.go`

**Evidence**:
- Generated `prerequisites.required_access` as array of strings, but model expects array of `RequiredAccess` objects
- Generated `prerequisites.required_tools` as custom format, but model expects `environment_setup.tools`
- Generated `references` as array of objects, but model expects object with `{documentation: [], related_code: [], sops: []}`
- Included `conceptual_overview` field which doesn't exist in Go model

**Impact**: Test plans fail to parse, generation workflow broken

**Root Cause**: Testplan-generator skill doesn't validate against actual Go struct definitions

**Solution**:
1. Add schema validation step to testplan-generator skill
2. Before generating JSON, read `internal/models/testplan.go` and extract struct definitions
3. Use json-schema-generator to create JSON schema from Go structs
4. Validate generated JSON against schema before writing file
5. Update skill templates to match current model structure

**Implementation**:
```bash
# Add to testplan-generator skill Phase 7: Schema Validation
go install github.com/a-h/generate/...@latest
generate -f internal/models/testplan.go -o schema/testplan-schema.json
jq . generated-testplan.json --schema schema/testplan-schema.json
```

**Priority**: 🔴 **CRITICAL** - Blocks all test plan generation

---

### 2. Missing Cleanup Tests for Modifies-State Tests

**Problem**: Test cases that modify state don't always have corresponding cleanup tests

**Evidence**:
- `test_create_probe_via_api` creates probe in RHOBS API but has no cleanup test
- Alternative generation workflows may create modifying tests without cleanup pairs

**Impact**: Resource leaks, test environment pollution, failed subsequent runs

**Solution**:
1. Update testplan-generator skill Phase 6 to REQUIRE cleanup test for every modifies-state test
2. Auto-generate cleanup test skeleton when modifying test is created
3. Validate during parsing that every modifies-state test has cleanup_test field populated
4. Add validator check: `if system_impact.type == "modifies-state" && !state_management.cleanup_test { error }`

**Implementation in validator.go**:
```go
func (v *Validator) validateCleanupCoverage(tp *models.TestPlan) error {
    for id, test := range tp.TestCases {
        if test.Metadata.SystemImpact.Type == "modifies-state" {
            if !test.Metadata.StateManagement.RequiresCleanup {
                return fmt.Errorf("test %s modifies state but requires_cleanup=false", id)
            }
            if test.Metadata.StateManagement.CleanupTest == "" {
                return fmt.Errorf("test %s modifies state but cleanup_test not specified", id)
            }
            // Verify cleanup test exists
            if _, exists := tp.TestCases[test.Metadata.StateManagement.CleanupTest]; !exists {
                return fmt.Errorf("test %s cleanup_test %s does not exist", id, test.Metadata.StateManagement.CleanupTest)
            }
        }
    }
    return nil
}
```

**Priority**: 🔴 **CRITICAL** - Resource leaks, broken test workflows

---

### 3. Alternative Paths Don't Validate Dependencies

**Problem**: Alternative paths can include tests without including their dependencies

**Evidence**:
- "quick_validation" path includes `test_verify_probe_metrics` which depends on `test_create_routemonitor_cr`, but that test isn't in the path
- Parser validates `learning_path.sequence` but not `learning_path.alternative_paths`

**Impact**: Alternative paths fail when executed, confusing user experience

**Solution**:
1. Extend validator to check alternative paths
2. For each alternative path, verify all test dependencies are either:
   - Included in the path earlier in the sequence, OR
   - Listed as external prerequisites for the path
3. Add validation error if dependency missing

**Implementation in validator.go**:
```go
func (v *Validator) validateAlternativePaths(tp *models.TestPlan) error {
    for pathName, pathTests := range tp.LearningPath.AlternativePaths {
        for _, testID := range pathTests {
            test := tp.TestCases[testID]
            for _, depID := range test.TestExecution.Dependencies {
                // Check if dependency is earlier in this path
                depFound := false
                for _, earlierTest := range pathTests {
                    if earlierTest == depID {
                        depFound = true
                        break
                    }
                    if earlierTest == testID {
                        break // Reached current test, dep must come before
                    }
                }
                if !depFound {
                    return fmt.Errorf("alternative path %s includes test %s which depends on %s, but %s is not in the path", pathName, testID, depID, depID)
                }
            }
        }
    }
    return nil
}
```

**Priority**: 🔴 **CRITICAL** - Breaks alternative execution paths

---

## High Priority Improvements

### 4. HTML Template Doesn't Render Hyperlinks

**Problem**: Educational enhancements add hyperlinks in markdown format `[text](url)`, but HTML template renders them as plain text

**Evidence**:
- learning_notes contain `[Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)`
- HTML shows literal brackets instead of clickable links

**Impact**: User can't click "learn more" links, defeats purpose of educational hyperlinks

**Solution**:
1. Add markdown parsing to HTML template
2. Use Go's `blackfriday` or `goldmark` markdown library
3. Parse `learning_note` fields as markdown before rendering
4. Sanitize HTML to prevent XSS (use `bluemonday` library)

**Implementation in html.go**:
```go
import (
    "github.com/gomarkdown/markdown"
    "github.com/gomarkdown/markdown/html"
    "github.com/gomarkdown/markdown/parser"
    "github.com/microcosm-cc/bluemonday"
)

funcMap := template.FuncMap{
    // ... existing functions ...
    "markdown": func(s string) template.HTML {
        md := parser.NewWithExtensions(parser.CommonExtensions)
        renderer := html.NewRenderer(html.RendererOptions{})
        unsafe := markdown.ToHTML([]byte(s), md, renderer)
        safe := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
        return template.HTML(safe)
    },
}
```

**In testplan.html template**:
```html
<!-- Before -->
<p class="text-gray-700">{{ $step.LearningNote }}</p>

<!-- After -->
<div class="text-gray-700 prose prose-sm max-w-none">
  {{ markdown $step.LearningNote }}
</div>
```

**Priority**: 🟡 **HIGH** - Core educational feature not working

---

### 5. No Cleanup Tracking System

**Problem**: When tests run, there's no tracking of which cleanup tests are needed

**Evidence**:
- Test 2 runs, user cancels execution, Test 5 cleanup never runs
- Resources left in cluster, next test run fails due to conflicting resources

**Impact**: Cluster pollution, manual cleanup burden, test unreliability

**Solution**:
1. Create cleanup tracking in testplan-executor skill
2. Before running modifying test, record cleanup requirement in `cleanup_needed.json`
3. After test completes (success or failure), prompt user to run cleanup
4. At end of test run, show summary: "3 tests require cleanup. Run now?"
5. Save state for later if user chooses "later"

**Implementation**:
```bash
# In test execution wrapper
run_test() {
    local test_id=$1
    local cleanup_test=$(jq -r ".testcases[\"$test_id\"].metadata.state_management.cleanup_test" testplan.json)

    if [[ -n "$cleanup_test" && "$cleanup_test" != "null" ]]; then
        echo "$cleanup_test" >> cleanup_needed.txt
    fi

    # Run the actual test
    execute_test_steps "$test_id"
}

# At end of test run
if [ -f cleanup_needed.txt ]; then
    echo "⚠️ $(wc -l < cleanup_needed.txt) cleanup tests needed:"
    cat cleanup_needed.txt
    read -p "Run cleanup now? (yes/no/later): " CLEANUP
    if [[ "$CLEANUP" == "yes" ]]; then
        while read cleanup_test; do
            run_test "$cleanup_test"
        done < cleanup_needed.txt
        rm cleanup_needed.txt
    elif [[ "$CLEANUP" == "later" ]]; then
        mv cleanup_needed.txt cleanup_needed_$(date +%Y%m%d_%H%M%S).txt
        echo "Saved for later. Run: ./run_cleanup.sh cleanup_needed_*.txt"
    fi
fi
```

**Priority**: 🟡 **HIGH** - Test reliability and environment cleanliness

---

### 6. No Production Environment Detection

**Problem**: Tests check `can_run_in_production` but don't enforce it

**Evidence**:
- Test metadata has `can_run_in_production: false`, but nothing prevents execution in prod
- User could accidentally run destructive test in production cluster

**Impact**: **CATASTROPHIC** if modifying test runs in production

**Solution**:
1. Add environment detection to testplan-executor skill pre-flight checks
2. Query cluster for production indicators:
   - Cluster name contains "prod", "production"
   - Infrastructure name matches known production clusters
   - Specific label `environment=production`
3. Block execution of tests with `can_run_in_production: false`
4. Require explicit override flag: `--allow-production` (log and audit this)

**Implementation**:
```bash
detect_production() {
    local cluster_id=$(oc get clusterversion -o jsonpath='{.items[0].spec.clusterID}')
    local cluster_name=$(oc get infrastructure cluster -o jsonpath='{.status.infrastructureName}')
    local env_label=$(oc get infrastructure cluster -o jsonpath='{.metadata.labels.environment}')

    # Check patterns
    if [[ "$cluster_name" =~ prod|production ]]; then
        echo "PRODUCTION"
        return 0
    fi

    if [[ "$env_label" == "production" ]]; then
        echo "PRODUCTION"
        return 0
    fi

    # Check against known production cluster IDs (from OCM)
    if grep -q "$cluster_id" production_clusters.txt; then
        echo "PRODUCTION"
        return 0
    fi

    echo "NON-PRODUCTION"
    return 1
}

check_production_safety() {
    local test_id=$1
    local can_run_prod=$(jq -r ".testcases[\"$test_id\"].metadata.safety.can_run_in_production" testplan.json)
    local env=$(detect_production)

    if [[ "$env" == "PRODUCTION" && "$can_run_prod" == "false" ]]; then
        echo "❌ BLOCKED: Test $test_id cannot run in production"
        echo "Current cluster: $(oc whoami --show-server)"
        echo "Environment: $env"
        exit 1
    fi
}
```

**Priority**: 🟡 **HIGH** - Safety critical for production clusters

---

### 7. Testplan-Generator Doesn't Auto-Generate Cleanup Tests

**Problem**: When generator creates modifies-state test, it doesn't create matching cleanup test

**Evidence**:
- Generator created `test_create_probe_via_api` but left cleanup_test field empty
- Manual creation of cleanup test required

**Impact**: Incomplete test plans, manual work, easy to forget cleanup

**Solution**:
1. Update testplan-generator skill Phase 6 to auto-generate cleanup test
2. When creating modifies-state test, generate paired cleanup test automatically
3. Cleanup test template:
   - Category: "cleanup"
   - System impact: same affected_resources but deletion operations
   - State management: restores_state: true
   - Safety: same safety constraints as modify test
   - Steps: reverse operations of modify test

**Implementation Pattern**:
```json
// When generating test_create_probe_via_api
{
  "test_create_probe_via_api": {
    "metadata": {
      "system_impact": {"type": "modifies-state"},
      "state_management": {
        "cleanup_test": "test_cleanup_probe_via_api"  // Auto-populate
      }
    }
  }
}

// Auto-generate paired cleanup test
{
  "test_cleanup_probe_via_api": {
    "metadata": {
      "category": "cleanup",
      "system_impact": {
        "type": "modifies-state",
        "description": "Deletes probe created by test_create_probe_via_api",
        "affected_resources": ["Same as test_create_probe_via_api"]
      },
      "state_management": {
        "restores_state": true
      }
    },
    "test_execution": {
      "steps": [
        {"title": "Delete probe via API", "command": "curl -X DELETE ..."},
        {"title": "Verify Probe CR deleted", "command": "oc get probes | grep ..."},
        {"title": "Confirm cleanup complete", "expected_output": "No resources found"}
      ]
    }
  }
}
```

**Priority**: 🟡 **HIGH** - Automation and completeness

---

### 8. JSON Schema Not Defined

**Problem**: No formal JSON schema exists for test plan format

**Evidence**:
- Test plan structure defined only in Go structs
- No validation before attempting to parse
- Errors discovered late (after generation complete)

**Impact**: Wasted time generating invalid JSON, confusing error messages

**Solution**:
1. Generate JSON Schema from Go models
2. Add schema validation step to generator and parser
3. Provide clear error messages when validation fails
4. Include schema in repository at `schema/testplan-v1.schema.json`

**Implementation**:
```bash
# Generate schema from Go
go install github.com/invopop/jsonschema/cmd/jsonschema@latest
jsonschema -output schema/testplan-v1.schema.json github.com/nephomaniac/testplan_tools_poc/internal/models.TestPlan

# Validate in generator
jq --schema schema/testplan-v1.schema.json . generated-testplan.json || {
    echo "❌ Generated JSON doesn't match schema"
    echo "Run: jq --schema schema/testplan-v1.schema.json . generated-testplan.json"
    exit 1
}

# Validate in parser (before unmarshaling)
validateSchema(jsonBytes, schemaPath) error
```

**Priority**: 🟡 **HIGH** - Fail fast, better DX

---

### 9. No Evidence Collection on Test Failure

**Problem**: When tests fail, no artifacts are automatically captured

**Evidence**:
- Test 2 fails, no logs saved
- User has to manually re-run commands to investigate
- Evidence lost if resources are deleted

**Impact**: Difficult troubleshooting, time wasted reproducing failures

**Solution**:
1. Add evidence collection to testplan-executor skill
2. On test failure, automatically capture:
   - All oc get commands (YAML output)
   - Relevant pod logs
   - API responses
   - Prometheus query results
3. Save to `evidence/<test-id>-<timestamp>/` directory
4. Include evidence collection in test report

**Implementation**:
```bash
collect_evidence() {
    local test_id=$1
    local evidence_dir="evidence/${test_id}-$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$evidence_dir"

    echo "📸 Collecting evidence for failed test: $test_id"

    # Capture all relevant resources
    oc get all -A -o yaml > "$evidence_dir/all-resources.yaml"
    oc get events -A --sort-by='.lastTimestamp' > "$evidence_dir/events.txt"

    # Test-specific artifacts
    case "$test_id" in
        test_create_routemonitor_cr)
            oc get routemonitors -A -o yaml > "$evidence_dir/routemonitors.yaml"
            oc get servicemonitors -A -o yaml > "$evidence_dir/servicemonitors.yaml"
            oc logs -n openshift-route-monitor-operator deployment/route-monitor-operator --tail=200 > "$evidence_dir/rmo-logs.txt"
            ;;
        test_verify_probe_metrics)
            curl "http://prometheus:9090/api/v1/query?query=probe_success" > "$evidence_dir/probe-metrics.json"
            ;;
    esac

    echo "✅ Evidence saved to: $evidence_dir"
    echo "$evidence_dir" >> test-failures.log
}

# In test execution wrapper
run_test() {
    local test_id=$1
    if ! execute_test_steps "$test_id"; then
        collect_evidence "$test_id"
        return 1
    fi
}
```

**Priority**: 🟡 **HIGH** - Faster debugging, better troubleshooting

---

### 10. Validator Doesn't Check Cleanup Test Properties

**Problem**: Validator checks cleanup_test exists, but doesn't verify it's actually a cleanup test

**Evidence**:
- Could set `cleanup_test: "test_verify_route_monitor_deployment"` (wrong category)
- Could reference modifying test as cleanup test
- No validation that cleanup test has `restores_state: true`

**Impact**: Broken cleanup workflows, resources not restored

**Solution**:
1. Extend validator to check cleanup test properties
2. Verify cleanup test has:
   - Category: "cleanup" OR system_impact.type matching parent test
   - state_management.restores_state: true
   - No cleanup_test of its own (cleanup tests don't need cleanup)
3. Warn if cleanup test has dependencies (should be standalone)

**Implementation in validator.go**:
```go
func (v *Validator) validateCleanupTest(tp *models.TestPlan, testID string, cleanupTestID string) error {
    cleanupTest, exists := tp.TestCases[cleanupTestID]
    if !exists {
        return fmt.Errorf("cleanup test %s not found", cleanupTestID)
    }

    // Verify it's actually a cleanup test
    if cleanupTest.Metadata.Category != "cleanup" && cleanupTest.Metadata.SystemImpact.Type != "modifies-state" {
        return fmt.Errorf("test %s specifies cleanup_test %s, but %s is category '%s' (expected 'cleanup')",
            testID, cleanupTestID, cleanupTestID, cleanupTest.Metadata.Category)
    }

    // Verify it restores state
    if !cleanupTest.Metadata.StateManagement.RestoresState {
        return fmt.Errorf("cleanup test %s does not set restores_state=true", cleanupTestID)
    }

    // Warn if cleanup test has its own cleanup test (recursive cleanup)
    if cleanupTest.Metadata.StateManagement.CleanupTest != "" {
        if v.verbose {
            fmt.Printf("   ⚠ Cleanup test %s has its own cleanup_test (recursive cleanup)\n", cleanupTestID)
        }
    }

    return nil
}
```

**Priority**: 🟡 **HIGH** - Data integrity, workflow correctness

---

## Medium Priority Improvements

### 11. No Test Duration Tracking

**Problem**: Estimated times are static, no actual duration captured

**Impact**: Can't improve estimates, can't identify slow tests

**Solution**:
1. Wrap test execution in timing logic
2. Record actual duration vs estimated
3. Save to test-results.json
4. Generate duration report: "Test 2 took 35min (estimated 25min, +40%)"
5. Use historical data to improve estimates over time

**Implementation**:
```bash
time_test() {
    local test_id=$1
    local estimated=$(jq -r ".testcases[\"$test_id\"].metadata.estimated_time" testplan.json)

    local start=$(date +%s)
    run_test "$test_id"
    local result=$?
    local end=$(date +%s)
    local actual=$((end - start))

    echo "{\"test_id\": \"$test_id\", \"estimated\": \"$estimated\", \"actual\": \"${actual}s\", \"status\": $result}" >> durations.json

    return $result
}
```

**Priority**: 🟢 **MEDIUM** - Metrics and continuous improvement

---

### 12. No Partial Success Handling

**Problem**: Test is pass/fail, no granularity for partially completed tests

**Impact**: Step 3 of 5 fails, but steps 1-2 succeeded - lose that information

**Solution**:
1. Report step-level pass/fail
2. Mark test as "PARTIAL SUCCESS" if >50% steps pass
3. Allow resuming from failed step instead of restarting test
4. Save step completion state: `test-2-progress.json`

**Implementation**:
```json
{
  "test_id": "test_create_routemonitor_cr",
  "status": "partial_success",
  "total_steps": 5,
  "completed_steps": 3,
  "failed_step": 4,
  "step_results": [
    {"step": 1, "status": "pass", "duration": "5s"},
    {"step": 2, "status": "pass", "duration": "10s"},
    {"step": 3, "status": "pass", "duration": "30s"},
    {"step": 4, "status": "fail", "error": "ServiceMonitor not found"},
    {"step": 5, "status": "skipped"}
  ]
}
```

**Priority**: 🟢 **MEDIUM** - Better failure analysis

---

### 13. Missing Test Report Template

**Problem**: No standardized way to report results

**Impact**: Inconsistent reporting, manual formatting work

**Solution**:
1. Create HTML report template
2. Generate after test run: `results-<timestamp>.html`
3. Include:
   - Executive summary (5 tests, 4 pass, 1 fail)
   - Per-test details with timestamps
   - Failures with logs and screenshots
   - Cleanup status
   - Recommendations

**Template**:
```html
<!DOCTYPE html>
<html>
<head><title>Test Results - RHOBS-next</title></head>
<body>
  <h1>Test Execution Report</h1>
  <div class="summary">
    <h2>Summary</h2>
    <p>Executed: 5 tests</p>
    <p>Passed: 4 (80%)</p>
    <p>Failed: 1 (20%)</p>
    <p>Duration: 95 minutes (estimated 100 minutes)</p>
  </div>

  <div class="results">
    {{#each tests}}
    <div class="test {{ status }}">
      <h3>{{ id }}: {{ title }}</h3>
      <p>Status: {{ status }}</p>
      <p>Duration: {{ actual_duration }} (estimated {{ estimated_duration }})</p>
      {{#if failed}}
      <div class="failure">
        <h4>Failure Details</h4>
        <pre>{{ error_log }}</pre>
        <p>Evidence: <a href="evidence/{{ id }}">View artifacts</a></p>
      </div>
      {{/if}}
    </div>
    {{/each}}
  </div>
</body>
</html>
```

**Priority**: 🟢 **MEDIUM** - Professional reporting

---

### 14. No Conceptual Overview Field in Go Model

**Problem**: testplan-educator skill adds `conceptual_overview` to tests, but it's not in Go model

**Evidence**:
- Enhancement agent added conceptual_overview to test cases
- Parser fails because field doesn't exist in TestCase struct

**Impact**: Educational enhancements don't persist

**Solution**:
1. Add `ConceptualOverview` field to `TestCase` struct in `internal/models/testplan.go`
2. Update HTML template to render conceptual overview before test steps
3. Update testplan-educator skill to use this field

**Implementation in testplan.go**:
```go
type TestCase struct {
    Metadata            TestMetadata    `json:"metadata"`
    ConceptualOverview  *string         `json:"conceptual_overview,omitempty"`  // ADD THIS
    Learning            Learning        `json:"learning"`
    TestExecution       TestExecution   `json:"test_execution"`
    Troubleshooting     Troubleshooting `json:"troubleshooting"`
    NextSteps           NextSteps       `json:"next_steps"`
    References          References      `json:"references"`
}
```

**In testplan.html**:
```html
{{if $test.ConceptualOverview}}
<div class="conceptual-overview bg-blue-50 border-l-4 border-blue-500 p-4 mb-6">
    <h4 class="font-bold text-blue-900 mb-2">📖 Conceptual Overview</h4>
    <div class="text-blue-800 prose prose-sm">
        {{ markdown $test.ConceptualOverview }}
    </div>
</div>
{{end}}
```

**Priority**: 🟢 **MEDIUM** - Educational enhancement support

---

### 15. Skills Should Accept Parameters for Customization

**Problem**: Skills are invoked with fixed behavior, no customization options

**Evidence**:
- testplan-generator always generates full test plans, can't generate only cleanup tests
- testplan-educator can't be told to focus only on hyperlinks
- testplan-executor can't be configured for dry-run vs actual execution

**Impact**: Inflexible, requires editing skill definitions for variations

**Solution**:
1. Add parameter parsing to skills
2. Support flags like:
   - `testplan-generator --tests-only` (skip concepts/metadata)
   - `testplan-educator --hyperlinks-only` (don't add other enhancements)
   - `testplan-executor --dry-run` (simulate execution, don't run)
3. Document parameters in skill.md files

**Example Usage**:
```bash
# Generate only test cases, reuse existing concepts
claude-code /skill testplan-generator --tests-only --output tests-only.json

# Add hyperlinks without changing existing learning notes
claude-code /skill testplan-educator --hyperlinks-only --input testplan.json

# Review execution without running
claude-code /skill testplan-executor --dry-run --report-only
```

**Priority**: 🟢 **MEDIUM** - Flexibility and reusability

---

## Implementation Priority

### Sprint 1 (Week 1): Critical Fixes
- ✅ #1: Schema validation in testplan-generator
- ✅ #2: Auto-generate cleanup tests
- ✅ #3: Validate alternative paths

**Goal**: Prevent broken test plan generation

### Sprint 2 (Week 2): Safety & Quality
- ✅ #4: Markdown rendering in HTML
- ✅ #5: Cleanup tracking system
- ✅ #6: Production environment detection
- ✅ #7: Enhanced cleanup validation

**Goal**: Safe, reliable test execution

### Sprint 3 (Week 3): User Experience
- ✅ #8: JSON schema definition
- ✅ #9: Evidence collection
- ✅ #10: Duration tracking
- ✅ #11: Partial success handling

**Goal**: Better debugging and observability

### Sprint 4 (Week 4): Polish
- ✅ #12: Test report template
- ✅ #13: Conceptual overview support
- ✅ #14: Parameterized skills

**Goal**: Professional, flexible tooling

---

## Testing Strategy for Improvements

### For Each Improvement:

1. **Unit Test**: Test the specific function/validation
2. **Integration Test**: Test with example test plans
3. **Regression Test**: Ensure existing test plans still work
4. **Documentation**: Update README, skill docs, examples

### Example Test Cases:

**For #2 (Auto-generate cleanup):**
```bash
# Test: Generator creates cleanup test automatically
./test_generator_cleanup.sh
# Expected: modifies-state test + paired cleanup test both exist
# Expected: cleanup_test field populated correctly
# Expected: cleanup test has restores_state=true
```

**For #4 (Markdown rendering):**
```bash
# Test: Hyperlinks render as clickable links
./test_html_markdown.sh
# Expected: [text](url) becomes <a href="url">text</a>
# Expected: Bold **text** becomes <strong>text</strong>
# Expected: XSS prevented (script tags removed)
```

**For #6 (Production detection):**
```bash
# Test: Blocks modifying tests in production
export MOCK_CLUSTER_NAME="production-cluster-abc"
./test_production_detection.sh
# Expected: test_create_routemonitor_cr BLOCKS with error
# Expected: test_verify_route_monitor_deployment ALLOWS (read-only)
```

---

## Conclusion

These 15 improvements address critical gaps in test plan generation, validation, safety, and execution. Implementing them will:

- ✅ Prevent invalid test plan generation
- ✅ Ensure complete cleanup coverage
- ✅ Protect production environments
- ✅ Improve debugging and troubleshooting
- ✅ Enhance educational value
- ✅ Provide professional reporting

**Estimated Implementation Time**: 4 weeks (1 sprint per priority tier)

**Risk**: Low - improvements are additive, don't break existing functionality

**ROI**: High - significantly improves tool quality, safety, and usability
