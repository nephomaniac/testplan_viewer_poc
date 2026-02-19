# testplan-reviewer Skill

**Version**: 1.0.0
**Purpose**: Validate test plans for comprehensive coverage, safety, and quality before execution

---

## Overview

The testplan-reviewer skill performs post-generation quality assurance on test plans, validating that they meet comprehensive testing standards, have complete safety documentation, and are ready for safe execution.

**Use this skill when:**
- A test plan has been generated and needs validation
- You want to verify test coverage before execution
- You need to identify gaps or missing test types
- You want to assess overall test plan quality

**Do NOT use this skill when:**
- Generating a new test plan (use testplan-generator instead)
- Executing tests (use testplan-executor instead)
- Adding educational content (use testplan-educator instead)

---

## Inputs

### Required
- **Test plan JSON file path** - Path to the test plan to review (e.g., `examples/camo/camo-testplan.json`)

### Optional
- **Coverage targets** - Expected test type distribution (default: 30-40% negative, 100% input validation)
- **Safety level** - Minimum safety documentation level (default: all tests classified)
- **Report format** - Output format (markdown, json, console) (default: markdown)

---

## Workflow

### Phase 1: Test Plan Loading and Validation (5 minutes)

**Objective**: Load and validate test plan structure

**Actions**:
1. Read test plan JSON file
2. Validate JSON structure matches schema
3. Verify required sections present (metadata, testcases, concepts)
4. Check for structural errors

**Validation checks**:
```bash
# Read test plan
cat examples/camo/camo-testplan.json | jq .

# Validate structure
jq 'has("metadata") and has("testcases") and has("concepts")' examples/camo/camo-testplan.json

# Check test count
jq '.testcases | length' examples/camo/camo-testplan.json
```

**Output**: Structural validation results

---

### Phase 2: Coverage Validation (15-20 minutes)

**Objective**: Verify comprehensive test coverage across all test types

**Coverage Matrix**:

| Test Type | Expected Ratio | Checks |
|-----------|---------------|---------|
| Validation | 30-50% | Baseline functionality tests |
| Negative | 30-40% | Error handling tests |
| Input Validation | Per-input | All 5 categories covered |
| RBAC | Per-permission | All boundaries tested |
| Load | Per-critical-component | Resource constraints tested |
| Network | Per-distributed-component | Resilience tested |
| Cleanup | 100% of modifying tests | All state changes reversed |

**Analysis**:

**1. Calculate Test Type Distribution**
```bash
# Count tests by type
jq '.testcases | to_entries | map(.value.test_execution.child_tests // [] | map(.test_type)) | flatten | group_by(.) | map({type: .[0], count: length})' testplan.json

# Expected output:
[
  {"type": "validation", "count": 20},
  {"type": "negative", "count": 15},
  {"type": "install", "count": 5},
  {"type": "rbac", "count": 8},
  {"type": "load", "count": 4},
  {"type": "network", "count": 3},
  {"type": "cleanup", "count": 10}
]
```

**2. Verify Negative Test Ratio**
```bash
# Calculate negative test percentage
jq '
  (.testcases | to_entries | map(.value.test_execution.child_tests // [] | map(select(.test_type == "negative")) | length) | add) as $negative |
  (.testcases | to_entries | map(.value.test_execution.child_tests // [] | map(select(.test_type == "validation")) | length) | add) as $validation |
  ($negative / $validation * 100 | floor)
' testplan.json
```

**Target**: 30-40% negative tests relative to positive tests

**3. Check Input Validation Coverage**

For each input parameter, verify tests cover all 5 categories:

| Input | Positive | Negative | Missing | Corrupt | Boundary |
|-------|----------|----------|---------|---------|----------|
| image | ✓ | ✓ | ✓ | ✓ | N/A |
| replicas | ✓ | ✓ | ✓ | ✓ | ✓ |
| namespace | ✓ | ✓ | ✓ | ✓ | ✓ |

```bash
# Find all input validation tests
jq '.testcases | to_entries | map(.value.test_execution.child_tests // [] | map(select(.input_validation != null))) | flatten | group_by(.input_validation) | map({category: .[0].input_validation, count: length})' testplan.json
```

**4. Verify RBAC Boundary Coverage**

For each operation requiring RBAC, verify tests cover permission boundaries:

```bash
# Find RBAC tests and their levels
jq '.testcases | to_entries | map(.value.test_execution.child_tests // [] | map(select(.rbac_level != null))) | flatten | group_by(.rbac_level) | map({level: .[0].rbac_level, count: length})' testplan.json
```

**Expected**: Tests for view, edit, namespace-admin, cluster-admin levels

**5. Identify Coverage Gaps**

```bash
# Find features without tests
jq '.metadata.features_covered // [] | .[]' testplan.json > features.txt
jq '.testcases | to_entries | map(.value.metadata.feature // "unknown")' testplan.json > tested_features.txt
comm -23 <(sort features.txt) <(sort tested_features.txt)
```

**Output**: Coverage gap report

---

### Phase 3: Test Dependency Analysis (15-20 minutes)

**Objective**: Ensure test dependencies are correct and execution order is valid

**Checks**:

**1. Cleanup Test Pairing**

Every modifying test MUST have a cleanup test:

```bash
# Find tests that modify state without cleanup
jq -r '.testcases | to_entries | map(
  select(.value.safety.impact_type == "modifies-state") |
  select(.value.safety.cleanup_test_id == null or .value.safety.cleanup_test_id == "") |
  .key
) | .[]' testplan.json
```

**Expected**: Empty result (all modifying tests have cleanup)

**2. Dependency Chain Validation**

```bash
# Build dependency graph
jq '.testcases | to_entries | map({
  test_id: .key,
  depends_on: .value.execution_flow.depends_on // []
}) | group_by(.test_id)' testplan.json

# Check for circular dependencies
# (Use topological sort algorithm)
```

**3. Orphaned Tests**

```bash
# Find tests not in learning_path.sequence
jq -r '
  (.learning_path.sequence // []) as $sequence |
  (.testcases | keys) as $all_tests |
  ($all_tests - $sequence) | .[]
' testplan.json
```

**Expected**: Only cleanup tests should be orphaned (not in main sequence)

**Output**: Dependency analysis report with any issues flagged

---

### Phase 4: Safety and Risk Assessment (15-20 minutes)

**Objective**: Validate safety documentation and production readiness

**Safety Checks**:

**1. Impact Type Classification**

All tests MUST have `impact_type` assigned:

```bash
# Find tests without impact_type
jq -r '.testcases | to_entries | map(
  .value.test_execution.child_tests // [] | map(
    select(.impact_type == null or .impact_type == "") | .title
  )
) | flatten | .[]' testplan.json
```

**Expected**: Empty result

**2. Production Safety Validation**

Tests marked `can_run_in_production: true` MUST be read-only:

```bash
# Find production-safe tests that modify state
jq -r '.testcases | to_entries | map(
  .value.test_execution.child_tests // [] | map(
    select(.safety.can_run_in_production == true) |
    select(.impact_type != "read-only") |
    .title
  )
) | flatten | .[]' testplan.json
```

**Expected**: Empty result (production-safe = read-only)

**3. Cleanup Procedures**

High-risk tests MUST have detailed cleanup procedures:

```bash
# Find high-risk tests without cleanup documentation
jq -r '.testcases | to_entries | map(
  select(.value.safety.risk_level == "high" or .value.safety.risk_level == "critical") |
  select(.value.safety.cleanup_procedure == null or .value.safety.cleanup_procedure == "") |
  .key
) | .[]' testplan.json
```

**4. Blast Radius Documentation**

Impairment tests MUST document blast radius:

```bash
# Find impairment tests without blast radius
jq -r '.testcases | to_entries | map(
  .value.test_execution.child_tests // [] | map(
    select(.impact_type == "impairment" or .impact_type == "load-generation") |
    select(.network_requirements.blast_radius == null) |
    select(.resource_requirements.blast_radius == null) |
    .title
  )
) | flatten | .[]' testplan.json
```

**Output**: Safety assessment report with risk level summary

---

### Phase 5: Attribute Completeness Check (10-15 minutes)

**Objective**: Verify all filterable attributes are properly assigned

**Attribute Matrix**:

| Attribute | Required For | Validation |
|-----------|-------------|------------|
| test_type | ALL tests | Not null or empty |
| impact_type | ALL tests | Valid value (read-only, modifies-state, etc.) |
| rbac_level | RBAC tests | Valid level (view, edit, namespace-admin, cluster-admin) |
| input_validation | Input validation tests | Valid category (positive, negative, missing, corrupt, boundary) |
| resource_requirements | Load tests | Defined with CPU/memory |
| network_requirements | Network tests | Defined with impairment type |

**Validation Queries**:

```bash
# Check test_type completeness
jq '[.testcases | to_entries | map(.value.test_execution.child_tests // []) | flatten | map(select(.test_type == null or .test_type == ""))] | length' testplan.json

# Check impact_type completeness
jq '[.testcases | to_entries | map(.value.test_execution.child_tests // []) | flatten | map(select(.impact_type == null or .impact_type == ""))] | length' testplan.json

# Check RBAC-sensitive operations have rbac_level
jq '[.testcases | to_entries | map(.value.test_execution.child_tests // []) | flatten | map(select(.test_type == "rbac") | select(.rbac_level == null or .rbac_level == ""))] | length' testplan.json

# Check load tests have resource_requirements
jq '[.testcases | to_entries | map(.value.test_execution.child_tests // []) | flatten | map(select(.test_type == "load") | select(.resource_requirements == null))] | length' testplan.json

# Check network tests have network_requirements
jq '[.testcases | to_entries | map(.value.test_execution.child_tests // []) | flatten | map(select(.test_type == "network") | select(.network_requirements == null))] | length' testplan.json
```

**Expected**: All queries return 0 (complete attribute coverage)

**Output**: Attribute completeness report

---

### Phase 6: Quality Metrics and Scoring (10-15 minutes)

**Objective**: Calculate comprehensive quality score and provide recommendations

**Metrics to Calculate**:

**1. Coverage Score (0-100)**

```
Coverage Score = (
  (test_type_distribution_score * 0.3) +
  (negative_test_ratio_score * 0.2) +
  (input_validation_coverage_score * 0.2) +
  (rbac_coverage_score * 0.15) +
  (cleanup_coverage_score * 0.15)
)
```

**2. Safety Score (0-100)**

```
Safety Score = (
  (impact_classification_completeness * 0.3) +
  (production_safety_validation * 0.25) +
  (cleanup_documentation_quality * 0.25) +
  (risk_documentation_quality * 0.20)
)
```

**3. Attribute Completeness Score (0-100)**

```
Attribute Score = (
  (test_type_completeness * 0.25) +
  (impact_type_completeness * 0.25) +
  (rbac_level_completeness * 0.20) +
  (input_validation_completeness * 0.15) +
  (resource_requirements_completeness * 0.10) +
  (network_requirements_completeness * 0.05)
)
```

**4. Overall Quality Score (0-100)**

```
Overall Score = (
  (Coverage Score * 0.40) +
  (Safety Score * 0.35) +
  (Attribute Score * 0.25)
)
```

**Quality Tiers**:
- 90-100: Excellent - Ready for execution
- 75-89: Good - Minor improvements recommended
- 60-74: Fair - Significant gaps to address
- 0-59: Poor - Major improvements required

**Output**: Quality scorecard with detailed metrics

---

### Phase 7: Generate Review Report (10 minutes)

**Objective**: Create comprehensive review report with actionable recommendations

**Report Structure**:

```markdown
# Test Plan Review Report

**Test Plan**: <name>
**Reviewed**: <timestamp>
**Overall Quality Score**: <score>/100 (<tier>)

---

## Executive Summary

- **Total Tests**: <count>
- **Coverage Score**: <score>/100
- **Safety Score**: <score>/100
- **Attribute Completeness**: <score>/100
- **Critical Issues**: <count>
- **Warnings**: <count>

---

## Coverage Analysis

### Test Type Distribution

| Type | Count | Percentage | Target | Status |
|------|-------|------------|--------|--------|
| Validation | 20 | 40% | 30-50% | ✅ Pass |
| Negative | 15 | 30% | 30-40% | ✅ Pass |
| RBAC | 8 | 16% | 10-20% | ✅ Pass |
| Load | 4 | 8% | 5-10% | ✅ Pass |
| Network | 3 | 6% | 3-8% | ✅ Pass |

### Coverage Gaps

- ❌ Missing input validation for `replicas` parameter (boundary tests)
- ⚠️ No RBAC tests for cluster-admin operations
- ⚠️ Load tests missing for critical component X

---

## Safety Assessment

### Risk Level Distribution

| Level | Count | Tests |
|-------|-------|-------|
| Low | 30 | Read-only validation tests |
| Medium | 15 | Modifying tests with cleanup |
| High | 5 | Load/network impairment tests |
| Critical | 2 | Cluster-wide operations |

### Safety Issues

- ❌ Test `test_create_probe_via_api` has no cleanup test (CRITICAL)
- ⚠️ Test `test_inject_latency` missing blast radius documentation
- ✅ All production-safe tests are read-only

---

## Dependency Analysis

### Execution Order

Tests can be executed in the following order:
1. test_verify_deployment (no dependencies)
2. test_create_resource (depends on #1)
3. test_verify_metrics (depends on #2)
4. test_cleanup_resource (cleanup for #2)

### Dependency Issues

- ❌ Circular dependency: test_a → test_b → test_a
- ⚠️ Orphaned test: test_xyz not in learning path

---

## Attribute Completeness

| Attribute | Complete | Missing | Percentage |
|-----------|----------|---------|------------|
| test_type | 50 | 0 | 100% ✅ |
| impact_type | 48 | 2 | 96% ⚠️ |
| rbac_level | 8/8 | 0 | 100% ✅ |
| input_validation | 12/15 | 3 | 80% ⚠️ |
| resource_requirements | 4/4 | 0 | 100% ✅ |
| network_requirements | 2/3 | 1 | 67% ❌ |

---

## Recommendations

### Critical (Must Fix Before Execution)

1. **Create cleanup test for `test_create_probe_via_api`**
   - Impact: Resource leak, state corruption
   - Action: Generate cleanup test following pattern of other cleanup tests
   - Estimated effort: 15 minutes

2. **Fix circular dependency: test_a ↔ test_b**
   - Impact: Impossible to execute
   - Action: Restructure dependencies
   - Estimated effort: 10 minutes

### High Priority (Recommended Before Execution)

3. **Add missing input validation tests for `replicas` boundary cases**
   - Impact: Incomplete coverage
   - Action: Add tests for 0 replicas and maximum replicas
   - Estimated effort: 20 minutes

4. **Document blast radius for `test_inject_latency`**
   - Impact: Unknown risk during execution
   - Action: Add network_requirements.blast_radius field
   - Estimated effort: 5 minutes

### Medium Priority (Quality Improvements)

5. **Add RBAC tests for cluster-admin operations**
   - Impact: Security boundary not validated
   - Action: Add positive and negative RBAC tests
   - Estimated effort: 30 minutes

6. **Complete impact_type for 2 tests**
   - Impact: Cannot filter tests properly
   - Action: Assign appropriate impact_type values
   - Estimated effort: 5 minutes

---

## Conclusion

**Overall Assessment**: <TIER>

<Tier-specific summary>

**Ready for Execution**: <YES/NO>
- If YES: "This test plan meets quality standards and can be executed safely."
- If NO: "Address <count> critical issues before execution."

**Next Steps**:
1. Fix critical issues
2. Re-run testplan-reviewer
3. Generate updated HTML (make html-<example>)
4. Proceed with testplan-executor

---

**Generated by**: testplan-reviewer skill v1.0.0
**Report timestamp**: <ISO 8601 timestamp>
```

**Output**: Save report as `<testplan>-review-report.md`

---

## Output Format

### Review Report (Markdown)

Generated markdown file with:
- Executive summary with quality score
- Detailed coverage analysis
- Safety assessment
- Dependency validation
- Attribute completeness check
- Prioritized recommendations
- Clear YES/NO execution readiness decision

### Console Summary

```
╔════════════════════════════════════════════════════════════╗
║         Test Plan Quality Review Summary                   ║
╠════════════════════════════════════════════════════════════╣
║ Test Plan: RHOBS-Next Synthetic Monitoring                ║
║ Overall Score: 87/100 (GOOD)                              ║
╠════════════════════════════════════════════════════════════╣
║ Coverage Score:      92/100 ✅                            ║
║ Safety Score:        85/100 ⚠️                             ║
║ Attribute Score:     84/100 ⚠️                             ║
╠════════════════════════════════════════════════════════════╣
║ Critical Issues:     1 ❌                                  ║
║ Warnings:            4 ⚠️                                   ║
║ Info:                2 ℹ️                                   ║
╠════════════════════════════════════════════════════════════╣
║ Ready for Execution: NO                                    ║
║ Reason: 1 critical issue (missing cleanup test)           ║
╠════════════════════════════════════════════════════════════╣
║ Report saved: examples/rhobs-next/review-report.md        ║
╚════════════════════════════════════════════════════════════╝
```

---

## Success Criteria

A test plan passes review when:

**Coverage** ✅
- 30-40% negative tests relative to positive tests
- 100% input validation coverage for all inputs
- RBAC tests for all permission boundaries
- Load tests for performance-critical components
- Network tests for distributed components
- Cleanup tests for 100% of modifying tests

**Safety** ✅
- All tests have impact_type classified
- Production-safe tests are read-only
- High-risk tests have cleanup documentation
- Impairment tests document blast radius
- Risk levels assigned appropriately

**Attributes** ✅
- Every child test has test_type
- Every child test has impact_type
- RBAC tests have rbac_level
- Input validation tests have input_validation category
- Load tests have resource_requirements
- Network tests have network_requirements

**Dependencies** ✅
- No circular dependencies
- No orphaned tests (except cleanup)
- All modifying tests have cleanup pairs
- Execution order is valid

**Quality Score** ✅
- Overall score ≥ 75/100
- No critical issues blocking execution
- Report generated with actionable recommendations

---

## Tips for Using This Skill

**Best Practices**:
1. **Run after generation**: Review immediately after testplan-generator completes
2. **Fix critical issues first**: Address blocking issues before execution
3. **Iterate**: Re-run reviewer after fixes to verify
4. **Compare scores**: Track quality improvements over time
5. **Use as gate**: Require minimum score before execution approval

**Common Issues Found**:
- Missing cleanup tests for modifying operations
- Incomplete input validation coverage
- Production-safe tests that modify state
- Circular dependencies in test execution order
- Missing filterable attributes (test_type, impact_type)

**Integration with Other Skills**:
- **After testplan-generator**: Validate generated test plan
- **Before testplan-executor**: Safety gate before execution
- **With testplan-educator**: Verify educational content quality
- **Standalone**: Quality check existing test plans

---

## Example Usage

```bash
# Review a test plan
claude --skill testplan-reviewer "Review examples/rhobs-next/rhobs-next-testplan.json"

# Review with custom coverage targets
claude --skill testplan-reviewer "Review examples/camo/camo-testplan.json with 40% negative test target"

# Generate JSON report instead of markdown
claude --skill testplan-reviewer "Review examples/rhobs-next/rhobs-next-testplan.json output format json"
```

---

## Version History

**v1.0.0** (2026-02-19)
- Initial release
- Coverage validation (test type distribution, negative ratio, input validation)
- Dependency analysis (cleanup pairing, circular dependencies, orphans)
- Safety assessment (impact classification, production safety, cleanup docs)
- Attribute completeness (test_type, impact_type, rbac_level, etc.)
- Quality scoring (coverage, safety, attributes, overall)
- Review report generation (markdown, JSON, console)

---

**Maintained by**: SREP Observability Team
**Related Skills**: testplan-generator, testplan-educator, testplan-executor, testplan-from-jira
