# testplan-from-jira Skill

**Version**: 1.0.0
**Purpose**: Generate comprehensive test plans directly from Jira tickets (Stories, Epics)

---

## Overview

The testplan-from-jira skill automates test plan generation from Jira tickets, extracting requirements from acceptance criteria, user stories, and descriptions to create comprehensive, executable test plans.

**Use this skill when:**
- You have a Jira Story/Epic describing a feature or work item
- You want to generate tests that align with development work
- You need traceability between Jira tickets and test coverage
- You're implementing Test-Driven Development practices

**Do NOT use this skill when:**
- Generating tests from code/artifacts (use testplan-generator instead)
- Validating existing test plans (use testplan-reviewer instead)
- Executing tests (use testplan-executor instead)

---

## Inputs

### Required
- **Jira ticket ID** - Ticket to generate tests from (e.g., `SREP-3120`, `OHSS-1234`)

### Optional
- **Epic scope** - For Epics, generate tests for all child Stories (default: false)
- **Test type mix** - Target distribution (default: auto-detect from ticket type)
- **Output path** - Where to save generated test plan (default: `testplan.json`)

---

## Workflow

### Phase 1: Jira Ticket Fetching (5 minutes)

**Objective**: Retrieve Jira ticket via API and extract metadata

**Actions**:

**1. Fetch Jira Ticket**

```bash
# Fetch ticket via Jira REST API (using HTTPie)
http GET "https://issues.redhat.com/rest/api/2/issue/SREP-3120" \
  "Authorization: Bearer $(jq -r '.jira' ~/.config/claude/claude.json)"
```

**2. Extract Key Fields**

From the Jira response, extract:
- **Summary**: Ticket title
- **Description**: Full description with acceptance criteria
- **Issue Type**: Story, Epic, Task, Bug
- **Priority**: Critical, High, Normal, etc.
- **Labels**: Tags like "testing", "rbac", "performance"
- **Components**: Affected components
- **Epic Link**: Parent epic (if Story)
- **Acceptance Criteria**: Parsed from description

**Jira Description Parsing**:

Typical SREP Story format:
```
# Context
<Background information>

# User Story
As a [x], I want [y], so that I can [z].

# Impact/Value
<Why this matters>

# What needs to be done
1. Step one
2. Step two
3. Step three

# Acceptance Criteria
1. Clear, testable condition 1
2. Clear, testable condition 2
3. Clear, testable condition n
```

**3. Identify Test Generation Strategy**

Based on issue type:
- **Story**: Generate tests from acceptance criteria
- **Epic**: Generate tests from all child Stories
- **Bug**: Generate regression test + negative tests
- **Task**: Generate validation tests for task completion

**Output**: Structured Jira data ready for test generation

---

### Phase 2: Requirement Analysis (10-15 minutes)

**Objective**: Analyze requirements and identify test scenarios

**Analysis Steps**:

**1. Parse Acceptance Criteria**

Each acceptance criterion becomes a test scenario:

```
Acceptance Criterion: "Probe metrics appear in Prometheus within 60 seconds"

Generates:
- Positive test: Verify probe metrics appear within 60s
- Negative test: Verify error when probe fails
- Boundary test: Verify behavior at exactly 60s threshold
```

**2. Identify RBAC Requirements**

From description, identify permission levels needed:

```
"Create CRD in cluster" → requires cluster-admin
"Create deployment in namespace" → requires namespace-admin or edit
"Read pod status" → requires view
```

**3. Identify Input Validation Requirements**

From "What needs to be done" steps, identify inputs:

```
"Create RouteMonitor CR with target URL"

Inputs to validate:
- URL: positive (valid), negative (invalid), missing, corrupt (malformed), boundary (edge cases)
```

**4. Identify Load/Network Testing Needs**

From labels, components, and description:
- Labels: `performance` → add load tests
- Labels: `distributed` → add network tests
- Components: `api` → add API load tests

**Output**: Test scenario matrix

---

### Phase 3: Test Case Generation (20-30 minutes)

**Objective**: Generate comprehensive test cases from requirements

**Test Generation Pattern**:

For each acceptance criterion, generate:

### 1. Positive/Validation Tests

```json
{
  "test_id": "test_verify_probe_metrics",
  "title": "Verify probe metrics appear in Prometheus",
  "category": "validation",
  "difficulty": "intermediate",
  "estimated_duration": "15 minutes",
  "test_execution": {
    "child_tests": [{
      "id": "test_verify_probe_metrics_step_001",
      "title": "Create probe and wait for metrics",
      "test_type": "validation",
      "impact_type": "modifies-state",
      "rbac_level": "edit",
      "command": "oc create -f probe.yaml && sleep 65 && oc exec prometheus-0 -- curl http://localhost:9090/api/v1/query?query=probe_success",
      "expected_output": "\"result\":[{\"value\":[1]}]",
      "safety": {
        "can_run_in_production": false,
        "requires_cleanup": true,
        "cleanup_test_id": "test_cleanup_probe",
        "risk_level": "medium"
      }
    }]
  },
  "jira_traceability": {
    "jira_ticket": "SREP-3120",
    "jira_url": "https://issues.redhat.com/browse/SREP-3120",
    "acceptance_criterion": "Probe metrics appear in Prometheus within 60 seconds",
    "requirement_type": "functional"
  }
}
```

### 2. Negative Tests (30-40% ratio)

For each positive test, generate negative scenarios:

```json
{
  "test_id": "test_probe_invalid_url",
  "title": "Attempt to create probe with invalid URL",
  "category": "negative",
  "test_execution": {
    "child_tests": [{
      "title": "Create probe with invalid URL",
      "test_type": "negative",
      "impact_type": "read-only",
      "input_validation": "negative",
      "command": "oc create -f probe-invalid-url.yaml",
      "expected_output": "Error: invalid URL format"
    }]
  }
}
```

### 3. Input Validation Tests

For each input, generate 5-category matrix:

| Input | Positive | Negative | Missing | Corrupt | Boundary |
|-------|----------|----------|---------|---------|----------|
| URL | ✓ valid | ✓ invalid | ✓ no field | ✓ malformed | ✓ max length |

### 4. RBAC Tests

For operations requiring specific permissions:

```json
{
  "child_tests": [{
    "title": "Create probe as cluster-admin (should succeed)",
    "test_type": "rbac",
    "rbac_level": "cluster-admin",
    "impact_type": "modifies-state"
  }, {
    "title": "View user attempts probe creation (should fail)",
    "test_type": "negative",
    "rbac_level": "view",
    "expected_output": "Error from server (Forbidden)"
  }]
}
```

### 5. Cleanup Tests

For each modifying test, generate cleanup:

```json
{
  "test_id": "test_cleanup_probe",
  "title": "Cleanup probe and verify removal",
  "category": "cleanup",
  "test_execution": {
    "child_tests": [{
      "title": "Delete probe CR",
      "test_type": "cleanup",
      "impact_type": "modifies-state",
      "command": "oc delete probe test-probe",
      "safety": {
        "restores_state": true,
        "risk_level": "low"
      }
    }]
  }
}
```

**Output**: Complete set of test cases

---

### Phase 4: Traceability Matrix Creation (10-15 minutes)

**Objective**: Link each test to Jira requirements for bidirectional traceability

**Traceability Structure**:

```json
{
  "traceability_matrix": {
    "jira_epic": "SREP-3109",
    "jira_story": "SREP-3120",
    "total_acceptance_criteria": 5,
    "tests_per_criterion": {
      "Probe metrics appear in Prometheus within 60 seconds": [
        "test_verify_probe_metrics",
        "test_probe_metrics_timeout",
        "test_probe_invalid_url"
      ],
      "Probe status updates reflect blackbox-exporter health": [
        "test_verify_probe_status",
        "test_probe_status_failure"
      ]
    },
    "coverage": {
      "total_criteria": 5,
      "criteria_with_tests": 5,
      "coverage_percentage": 100
    }
  }
}
```

**Bi-directional Links**:
- **Test → Jira**: Each test has `jira_traceability` section
- **Jira → Tests**: Traceability matrix maps criteria to tests

**Output**: Traceability matrix

---

### Phase 5: Test Plan Metadata Generation (10 minutes)

**Objective**: Generate test plan metadata from Jira ticket

**Metadata Mapping**:

```json
{
  "metadata": {
    "title": "<Jira Summary> - Test Plan",
    "description": "Comprehensive test plan for <Jira Ticket>",
    "version": "1.0.0",
    "date": "<current date>",
    "authors": ["<Jira Assignee>", "Claude Sonnet 4.5"],
    "target_audience": "SRE Engineers, QE Team",
    "estimated_total_time": "<calculated from test durations>",
    "learning_objectives": [
      "<Extracted from Impact/Value section>",
      "Understand <feature> implementation",
      "Master testing strategies for <component>"
    ],
    "jira_integration": {
      "epic": "<Epic Link>",
      "story": "<Story ID>",
      "priority": "<Jira Priority>",
      "labels": ["<Jira Labels>"],
      "components": ["<Jira Components>"]
    }
  }
}
```

**Output**: Complete metadata section

---

### Phase 6: Output Generation (5 minutes)

**Objective**: Generate final test plan JSON with all components

**Test Plan Structure**:

```json
{
  "metadata": { /* From Phase 5 */ },
  "prerequisites": { /* Auto-generated based on test requirements */ },
  "concepts": { /* Generated from Jira context */ },
  "testcases": { /* From Phase 3 */ },
  "learning_path": {
    "title": "Progressive Learning Path",
    "description": "From basic validation to advanced integration",
    "sequence": ["test_1", "test_2", "test_3"],
    "alternative_paths": []
  },
  "traceability_matrix": { /* From Phase 4 */ }
}
```

**File Naming Convention**:
- For Story: `<component>-<ticket-id>-testplan.json` (e.g., `synthetics-srep-3120-testplan.json`)
- For Epic: `<epic-name>-testplan.json` (e.g., `rhobs-e2e-smoke-tests-testplan.json`)

**Save Location**:
- Default: `testplans/<ticket-id>-testplan.json`
- Custom: User-specified path

**Output**: Complete test plan JSON file ready for review/execution

---

## Test Type Distribution

Based on Jira ticket type, generate different test mixes:

### Story (Feature Implementation)
- 40% Validation tests (happy path)
- 30% Negative tests (error handling)
- 15% Input validation tests
- 10% RBAC tests
- 5% Cleanup tests

### Epic (Large Feature)
- Generate tests from all child Stories
- Add integration tests between Stories
- Add E2E smoke test suite

### Bug (Regression)
- 50% Regression test (reproduce bug)
- 30% Negative tests (similar scenarios)
- 20% Input validation (what caused bug)

### Task (Infrastructure/Tooling)
- 60% Validation tests (task completed)
- 30% Negative tests (failure modes)
- 10% Cleanup tests

---

## Jira Field Mapping

| Jira Field | Test Plan Field | Notes |
|------------|----------------|--------|
| Summary | metadata.title | Add "- Test Plan" suffix |
| Description | Various | Parse for context, acceptance criteria |
| Acceptance Criteria | Test cases | One test per criterion |
| Priority | Test priority | Map to test execution order |
| Components | metadata.components | Identify which components to test |
| Labels | Test types | performance → load tests, security → RBAC |
| Epic Link | metadata.jira_integration.epic | Parent epic reference |
| Assignee | metadata.authors | Original assignee + Claude |

---

## Example Usage

### Generate tests from Story

```bash
# Generate test plan from SREP-3120
claude --skill testplan-from-jira "Generate test plan for SREP-3120"

# Output: synthetics-srep-3120-testplan.json
```

### Generate tests from Epic (with child Stories)

```bash
# Generate comprehensive test plan from Epic and all children
claude --skill testplan-from-jira "Generate test plan for SREP-3109 including all child Stories"

# Output: rhobs-e2e-smoke-tests-testplan.json
# Includes tests from SREP-3117, SREP-3118, SREP-3119, SREP-3120
```

### Generate tests with custom output path

```bash
# Save to specific location
claude --skill testplan-from-jira "Generate test plan for SREP-3120 save to examples/synthetics/testplan.json"

# Output: examples/synthetics/testplan.json
```

---

## Output Format

### Test Plan JSON

Complete test plan following the standard structure:
- Metadata with Jira integration
- Prerequisites auto-generated from test requirements
- Concepts extracted from Jira context
- Test cases with traceability links
- Learning path ordered by dependencies
- Traceability matrix mapping criteria to tests

### Traceability Report

Optionally generate markdown traceability report:

```markdown
# Traceability Matrix: SREP-3120

## Jira Ticket
- **ID**: SREP-3120
- **Title**: Implement synthetics verification smoke test
- **Type**: Story
- **Priority**: High

## Acceptance Criteria Coverage

### AC1: Probe metrics appear in Prometheus within 60 seconds
- test_verify_probe_metrics ✅
- test_probe_metrics_timeout ✅
- test_probe_invalid_url ✅

### AC2: Probe status updates reflect blackbox-exporter health
- test_verify_probe_status ✅
- test_probe_status_failure ✅

## Coverage Summary
- Total Acceptance Criteria: 5
- Criteria with Tests: 5
- Coverage: 100%

## Test Distribution
- Validation: 8 (40%)
- Negative: 6 (30%)
- Input Validation: 3 (15%)
- RBAC: 2 (10%)
- Cleanup: 1 (5%)
```

---

## Success Criteria

A generated test plan is complete when:

**Traceability** ✅
- Every acceptance criterion has at least one test
- Every test links back to Jira ticket
- Traceability matrix shows 100% coverage

**Comprehensive Coverage** ✅
- 30-40% negative tests
- Input validation for all inputs
- RBAC tests for permission boundaries
- Cleanup tests for all modifying operations

**Quality** ✅
- All tests have complete safety attributes
- All tests have filterable attributes (test_type, impact_type)
- Execution order is valid (no circular dependencies)
- Passes testplan-reviewer validation

**Jira Integration** ✅
- Metadata populated from Jira fields
- Traceability links in every test
- Priority mapped to execution order
- Labels/components mapped to test types

---

## Integration with Other Skills

### Workflow: Jira → Test Plan → Review → Execute

```bash
# 1. Generate from Jira
claude --skill testplan-from-jira "Generate test plan for SREP-3120"

# 2. Review quality
claude --skill testplan-reviewer "Review synthetics-srep-3120-testplan.json"

# 3. Add educational content (optional)
claude --skill testplan-educator "Enhance synthetics-srep-3120-testplan.json"

# 4. Generate HTML
make html # or ./build/testplan-viewer -i synthetics-srep-3120-testplan.json -o testplan.html

# 5. Execute tests
claude --skill testplan-executor "Execute synthetics-srep-3120-testplan.json"
```

### Update Jira with Test Results

After execution, update Jira ticket with:
- Link to test plan file
- Link to test execution results
- Test coverage metrics
- Pass/fail status

```bash
# Add comment to Jira ticket with test results
http --ignore-stdin POST "https://issues.redhat.com/rest/api/2/issue/SREP-3120/comment" \
  "Authorization: Bearer $(jq -r '.jira' ~/.config/claude/claude.json)" \
  body:="$(cat <<'EOF'
{
  "body": "Test Plan Generated\n\n* Test Plan: synthetics-srep-3120-testplan.json\n* Total Tests: 20\n* Coverage: 100% of acceptance criteria\n* Quality Score: 92/100 (Excellent)\n\nReady for execution."
}
EOF
  )"
```

---

## Tips for Using This Skill

**Best Practices**:
1. **Well-written acceptance criteria**: Clear, testable criteria generate better tests
2. **Use Jira templates**: SREP Story/Epic templates have structured sections
3. **Review generated plan**: Always review with testplan-reviewer before execution
4. **Iterate**: Re-generate if Jira ticket updated with new requirements
5. **Link back**: Add test plan link to Jira ticket for traceability

**Common Issues**:
- **Vague acceptance criteria**: "It works" → Generate generic tests, manual refinement needed
- **Missing context**: No "What needs to be done" → May miss input validation scenarios
- **Complex Epics**: Many child Stories → Review each Story's tests individually first

**When to Use**:
- ✅ Feature development (TDD approach)
- ✅ Bug fixes (regression tests)
- ✅ Sprint planning (estimate test effort)
- ✅ Traceability requirements (link tests to work items)

**When NOT to Use**:
- ❌ No Jira ticket exists (use testplan-generator with artifacts instead)
- ❌ Jira ticket has no acceptance criteria (add criteria first)
- ❌ Testing existing implementation without ticket (use testplan-generator)

---

## Jira Ticket Quality Requirements

For best results, Jira tickets should have:

**Required**:
- Clear summary/title
- At least one acceptance criterion
- Defined component or epic link

**Recommended**:
- User story format (As a..., I want..., So that...)
- Context section explaining background
- What needs to be done section with steps
- Appropriate labels (testing, performance, security)

**Nice to Have**:
- References to design docs
- Examples or screenshots
- Known edge cases
- RBAC requirements documented

---

## Version History

**v1.0.0** (2026-02-19)
- Initial release
- Jira ticket fetching via REST API
- Requirement analysis from acceptance criteria
- Test generation (positive, negative, input validation, RBAC, cleanup)
- Traceability matrix creation
- Metadata generation from Jira fields
- Integration with testplan-reviewer, testplan-educator, testplan-executor

---

**Maintained by**: SREP Observability Team
**Related Skills**: testplan-generator, testplan-reviewer, testplan-educator, testplan-executor
