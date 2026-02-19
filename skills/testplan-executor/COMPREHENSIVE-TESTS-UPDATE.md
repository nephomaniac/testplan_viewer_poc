# testplan-executor: Comprehensive Test Types Update

**Date**: 2026-02-19
**Purpose**: Add execution support for comprehensive test types (negative, input validation, RBAC, load, network)

---

## New Test Types to Support

The testplan-generator and testplan-educator skills now generate these additional test types:

1. **Negative tests** - Expected failures with error validation
2. **Input validation tests** - 5 categories (positive, negative, missing, corrupt, boundary)
3. **RBAC tests** - Permission boundary validation
4. **Load tests** - Resource constraints and performance
5. **Network tests** - Resilience and impairment

The testplan-executor skill needs updates to handle these test types safely and effectively.

---

## Required Updates

### 1. Pre-Execution Filtering

Add ability to filter tests by type before execution:

```bash
# Execute only read-only tests (safe in production)
claude --skill testplan-executor "Execute testplan.json filter: read-only only"

# Execute validation and negative tests (skip load/network)
claude --skill testplan-executor "Execute testplan.json types: validation,negative"

# Execute all tests except impairment types
claude --skill testplan-executor "Execute testplan.json exclude: load,network"
```

**Implementation**: Check `impact_type` and `test_type` attributes before execution

### 2. RBAC Verification

Before executing RBAC tests, verify current user has required permissions:

```bash
# Check RBAC level before test
if [ "$rbac_level" == "cluster-admin" ]; then
  oc auth can-i '*' '*' --all-namespaces || {
    echo "❌ Test requires cluster-admin, but user lacks permission"
    echo "→ Skip test or elevate permissions"
    exit 1
  }
fi
```

### 3. Load Test Safety Checks

Before executing load tests:

```bash
# Verify not in production
if kubectl get namespace | grep -q "production\|prod"; then
  echo "❌ Load tests detected in production environment"
  echo "   Risk level: ${risk_level}"
  echo "   Resource requirements: ${cpu}, ${memory}"
  echo ""
  echo "Require confirmation to proceed? (yes/no)"
  read confirmation
  [ "$confirmation" != "yes" ] && exit 1
fi

# Check resource availability
available_cpu=$(kubectl top nodes | awk 'NR>1 {sum+=$3} END {print sum}')
required_cpu="${resource_requirements.cpu}"

if [ "$available_cpu" -lt "$required_cpu" ]; then
  echo "⚠️ Warning: Insufficient CPU for load test"
  echo "   Available: ${available_cpu}"
  echo "   Required: ${required_cpu}"
fi
```

### 4. Network Test Impairment Tracking

Track network impairments and ensure cleanup:

```bash
# Before network test
echo "⚠️ NETWORK IMPAIRMENT TEST"
echo "   Type: ${impairment_type}"
echo "   Config: ${impairment_config}"
echo "   Blast radius: ${blast_radius}"
echo "   Cleanup required: YES"
echo ""

# Store impairment state for cleanup
echo "${test_id}|${impairment_type}|${impairment_config}" >> /tmp/active-impairments.txt

# Execute test
...

# Verify cleanup
if grep -q "${test_id}" /tmp/active-impairments.txt; then
  echo "❌ Impairment not cleaned up!"
  echo "   Manual cleanup required:"
  echo "   ${cleanup_commands}"
fi
```

### 5. Input Validation Test Execution

For input validation tests, show expected vs actual:

```bash
# Positive input validation
echo "✅ Testing with VALID input: ${input_value}"
result=$(execute_command)
if [ "$result" == "success" ]; then
  echo "   ✓ Accepted as expected"
else
  echo "   ✗ Rejected (should have been accepted)"
fi

# Negative input validation
echo "❌ Testing with INVALID input: ${input_value}"
result=$(execute_command)
if [[ "$result" =~ "Error" ]]; then
  echo "   ✓ Rejected as expected"
  echo "   Error message: $result"
else
  echo "   ✗ Accepted (should have been rejected)"
fi
```

### 6. Negative Test Expected Failure Handling

Negative tests SHOULD fail - that's success:

```bash
# Execute negative test
exit_code=0
output=$(execute_command) || exit_code=$?

# For negative tests, failure is success
if [ "$test_type" == "negative" ]; then
  if [ $exit_code -ne 0 ] && [[ "$output" =~ "$expected_output" ]]; then
    echo "✅ Test PASSED (failed as expected with correct error)"
  else
    echo "❌ Test FAILED (did not fail as expected)"
  fi
else
  # Normal test - failure is failure
  if [ $exit_code -eq 0 ]; then
    echo "✅ Test PASSED"
  else
    echo "❌ Test FAILED"
  fi
fi
```

---

## Execution Order for Comprehensive Tests

### Recommended Order:

1. **Setup/Install tests** - Prepare environment
2. **Validation tests** - Verify baseline functionality (read-only first)
3. **Input validation tests** - Positive → Negative → Missing → Corrupt → Boundary
4. **RBAC tests** - View → Edit → Namespace-Admin → Cluster-Admin
5. **Negative tests** - Error handling validation
6. **Load tests** - Performance and resource constraints
7. **Network tests** - Resilience (highest risk last)
8. **Cleanup tests** - Restore state

### Dependency-Based Execution:

```bash
# Build dependency graph
tests_with_deps=$(jq -r '.testcases | to_entries | map({id: .key, deps: .value.execution_flow.depends_on // []})' testplan.json)

# Topological sort for execution order
execution_order=$(topological_sort "$tests_with_deps")

# Execute in order
for test_id in $execution_order; do
  execute_test "$test_id"
done
```

---

## Safety Gates for High-Risk Tests

### Pre-Execution Checklist for Impairment Tests:

```
High-Risk Test Detected: ${test_id}

Risk Level: ${risk_level}
Impact Type: ${impact_type}
Blast Radius: ${blast_radius}

Pre-Flight Checklist:
[ ] Verified environment is NON-PRODUCTION
[ ] Verified backup/snapshot exists
[ ] Verified cleanup procedure documented
[ ] Verified monitoring/alerting in place
[ ] Notified team of impairment test
[ ] Confirmed rollback plan ready

Proceed with execution? (yes/no)
```

### Require Explicit Confirmation:

```bash
if [ "$risk_level" == "critical" ] || [ "$impact_type" == "impairment" ]; then
  echo "This test requires explicit confirmation."
  echo "Type the test ID to confirm: ${test_id}"
  read user_input

  if [ "$user_input" != "$test_id" ]; then
    echo "❌ Confirmation failed. Test skipped."
    return 1
  fi
fi
```

---

## Parallel Execution for Independent Tests

Tests without dependencies can run in parallel:

```bash
# Find independent tests
independent_tests=$(jq -r '.testcases | to_entries | map(select(.value.execution_flow.depends_on == null or .value.execution_flow.depends_on | length == 0)) | map(.key) | .[]' testplan.json)

# Execute in parallel (up to 4 concurrent)
echo "$independent_tests" | xargs -n 1 -P 4 -I {} execute_test {}

# Wait for all to complete
wait

# Check results
check_parallel_results
```

**Safety note**: Only parallelize read-only tests. Never parallelize:
- Modifying tests
- Tests with shared resources
- Load/network impairment tests

---

## Rollback on Failure

If a test fails and leaves system in bad state:

```bash
execute_test() {
  local test_id="$1"
  local cleanup_test_id=$(jq -r ".testcases.\"$test_id\".safety.cleanup_test_id" testplan.json)

  # Execute test
  if ! run_test "$test_id"; then
    echo "❌ Test failed: $test_id"

    # If cleanup test exists, run it immediately
    if [ -n "$cleanup_test_id" ] && [ "$cleanup_test_id" != "null" ]; then
      echo "🧹 Running cleanup test: $cleanup_test_id"
      run_test "$cleanup_test_id"
    fi

    return 1
  fi
}
```

---

## Test Result Enhancements

For comprehensive test types, capture additional metadata:

```json
{
  "test_id": "test_verify_rbac_cluster_admin",
  "test_type": "rbac",
  "rbac_level": "cluster-admin",
  "result": "PASS",
  "execution_time": "5s",
  "timestamp": "2026-02-19T10:30:00Z",
  "actual_user": "maclark",
  "user_rbac_level": "cluster-admin",
  "rbac_verification": {
    "required_level": "cluster-admin",
    "actual_level": "cluster-admin",
    "match": true
  }
}
```

For negative tests:

```json
{
  "test_id": "test_invalid_image",
  "test_type": "negative",
  "input_validation": "negative",
  "result": "PASS",
  "expected_failure": true,
  "actual_failure": true,
  "expected_error": "Error: ErrImagePull",
  "actual_error": "Error: ErrImagePull: image not found",
  "error_match": true
}
```

For load tests:

```json
{
  "test_id": "test_cpu_stress",
  "test_type": "load",
  "result": "PASS",
  "resource_usage": {
    "cpu_before": "20%",
    "cpu_during": "95%",
    "cpu_after": "22%",
    "memory_before": "2Gi",
    "memory_during": "6Gi",
    "memory_after": "2.1Gi"
  },
  "performance_impact": "High CPU during test, returned to baseline after cleanup"
}
```

---

## Integration with testplan-reviewer

Before execution, run review:

```bash
# 1. Review test plan quality
claude --skill testplan-reviewer "Review testplan.json"

# 2. Check review score
quality_score=$(jq -r '.overall_score' testplan-review-report.json)

if [ "$quality_score" -lt 75 ]; then
  echo "❌ Test plan quality score too low: $quality_score/100"
  echo "   Fix critical issues before execution"
  exit 1
fi

# 3. Check for blocking issues
critical_issues=$(jq -r '.critical_issues | length' testplan-review-report.json)

if [ "$critical_issues" -gt 0 ]; then
  echo "❌ ${critical_issues} critical issue(s) found"
  echo "   Review: testplan-review-report.md"
  exit 1
fi

# 4. Proceed with execution
echo "✅ Quality check passed. Proceeding with execution."
```

---

## Summary of Required Changes

### testplan-executor/skill.md Updates:

1. **Add Pre-Execution Filtering section** - Filter by test_type, impact_type
2. **Add RBAC Verification section** - Check user permissions before RBAC tests
3. **Add Load Test Safety section** - Environment checks, resource verification
4. **Add Network Impairment Tracking section** - Track and verify cleanup
5. **Update Test Execution Logic** - Handle negative tests (expected failures)
6. **Add Parallel Execution section** - Execute independent read-only tests concurrently
7. **Add Rollback on Failure section** - Auto-cleanup on test failure
8. **Enhance Result Capture** - Additional metadata for new test types
9. **Add Integration with testplan-reviewer** - Quality gate before execution

### New Skills Integration:

- Call testplan-reviewer before execution
- Support test plans generated by testplan-from-jira (with Jira traceability)
- Report results back to Jira (optional)

---

## Implementation Priority

### High Priority (Required for Comprehensive Tests)
1. ✅ Negative test handling (expected failures = success)
2. ✅ RBAC verification
3. ✅ Safety gates for impairment tests
4. ✅ Cleanup tracking and verification

### Medium Priority (Quality of Life)
5. ⏳ Pre-execution filtering by test type
6. ⏳ Parallel execution for independent tests
7. ⏳ Integration with testplan-reviewer
8. ⏳ Enhanced result capture

### Lower Priority (Nice to Have)
9. ⏳ Rollback on failure automation
10. ⏳ Jira result reporting
11. ⏳ Performance metrics collection

---

**Status**: Update plan documented, ready for implementation
**Next**: Update skills/testplan-executor/skill.md with these patterns
