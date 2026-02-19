# testplan-generator Skill Update Summary

**Date**: 2026-02-19
**Purpose**: Add comprehensive test generation phases (7-12) to support filterable test cases

---

## Overview

Updated `skills/testplan-generator/skill.md` to generate comprehensive test coverage including negative testing, input validation, RBAC testing, load testing, and network testing using the new ChildTestCase format.

---

## Changes Summary

**File**: `skills/testplan-generator/skill.md`
- **Before**: 834 lines, 6 phases
- **After**: 1,652 lines, 12 phases
- **Added**: ~818 lines, 6 new phases

---

## New Phases Added

### Phase 7: Negative Testing Generation (30-45 minutes)

Generates negative test cases for every positive test to validate error handling:

**Test Categories:**
1. Invalid Configuration
2. Permission Denials
3. Resource Conflicts
4. Missing Dependencies

**Requirements:**
- Generate 2-3 negative tests per positive test
- All negative tests use `test_type: "negative"`
- Set `impact_type: "read-only"` (fails before making changes)
- Define `expected_output` with specific error message
- Production-safe (`can_run_in_production: true`)

**Example:**
```json
{
  "child_tests": [{
    "title": "Attempt deployment with invalid image",
    "test_type": "negative",
    "impact_type": "read-only",
    "input_validation": "negative",
    "expected_output": "Error: ErrImagePull"
  }]
}
```

---

### Phase 8: Input Validation Testing (30-45 minutes)

Generates comprehensive input validation tests covering all input categories:

**Categories:**
1. **Positive** (`input_validation: "positive"`) - Valid inputs
2. **Negative** (`input_validation: "negative"`) - Invalid types/values
3. **Missing** (`input_validation: "missing"`) - Required fields omitted
4. **Corrupt** (`input_validation: "corrupt"`) - Malformed data
5. **Boundary** (`input_validation: "boundary"`) - Min/max limits

**Input Validation Matrix:**

| Input | Positive | Negative | Missing | Corrupt | Boundary |
|-------|----------|----------|---------|---------|----------|
| image | ✓ Valid | ✓ Non-existent | ✓ No field | ✓ Malformed | N/A |
| replicas | ✓ Valid | ✓ Negative | ✓ No field | ✓ String | ✓ 0, 10000 |
| namespace | ✓ Exists | ✓ Non-existent | ✓ No ns | ✓ Invalid chars | ✓ 63-char |

---

### Phase 9: RBAC and Permission Testing (30-45 minutes)

Generates tests to validate role-based access control:

**RBAC Levels:**
1. **cluster-admin** - Cluster-scoped resources (CRDs, ClusterRoles)
2. **namespace-admin** - Namespace-scoped admin (Roles, RoleBindings)
3. **edit** - Standard resource creation/modification
4. **view** - Read-only operations

**Strategy:**
- For each operation, determine minimum required RBAC level
- Generate positive test with correct permissions
- Generate negative test with insufficient permissions
- Validate RBAC enforcement with "Forbidden" errors

**Example:**
```json
{
  "child_tests": [{
    "title": "Create CRD (requires cluster-admin)",
    "test_type": "rbac",
    "rbac_level": "cluster-admin",
    "impact_type": "modifies-state"
  }, {
    "title": "View user attempts CRD creation",
    "test_type": "negative",
    "rbac_level": "view",
    "expected_output": "Error from server (Forbidden)"
  }]
}
```

---

### Phase 10: Load and Performance Testing (45-60 minutes)

Generates tests to validate behavior under load and resource constraints:

**Load Test Categories:**

1. **Resource Limit Tests**
   - Test with insufficient CPU/memory limits
   - Validate error handling when limits exceeded

2. **Load Generation Tests** (`impact_type: "load-generation"`)
   - Deploy stress pods (CPU, memory, network)
   - Validate system continues functioning under load
   - **Must** include cleanup procedures

3. **Performance Tests**
   - Measure and validate performance metrics
   - Establish baselines

**Example:**
```json
{
  "child_tests": [{
    "title": "Deploy CPU stress pods",
    "test_type": "load",
    "impact_type": "load-generation",
    "resource_requirements": {
      "cpu": "2",
      "memory": "512Mi",
      "load_generation": true
    },
    "safety": {
      "can_run_in_production": false,
      "requires_cleanup": true,
      "risk_level": "medium"
    }
  }]
}
```

**Stress Deployment Templates Provided:**
- CPU stress (using polinux/stress image)
- Memory stress
- Network stress

---

### Phase 11: Network Testing with Impairment (45-60 minutes)

Generates tests to validate behavior under network stress:

**Network Test Categories:**

1. **Network Policy Tests**
   - Deny-all policies
   - Specific ingress/egress restrictions

2. **Latency Injection**
   - Add artificial network latency (using tc qdisc)
   - Test timeout handling

3. **Packet Loss**
   - Simulate packet loss scenarios
   - Validate retry logic

4. **Security Group Tests** (AWS-specific)
   - Temporarily block traffic via security groups
   - Test network partition scenarios

**Network Requirements Structure:**
```json
{
  "network_requirements": {
    "impairment_needed": true,
    "impairment_type": "latency",
    "impairment_config": "100ms delay on eth0",
    "network_policies": ["deny-all-policy.yaml"],
    "security_groups": ["sg-abc123"]
  }
}
```

**Safety Requirements:**
- ALL network tests: `can_run_in_production: false`
- Risk level: `"high"` or `"critical"`
- Detailed cleanup procedures required
- Test cleanup before marking complete

**Examples Provided:**
- Deny-all NetworkPolicy YAML
- tc qdisc commands for latency/packet-loss
- AWS security group modification commands

---

### Phase 12: Test Generation Summary

Generates a comprehensive summary showing coverage distribution:

```json
{
  "test_generation_summary": {
    "total_tests_generated": 87,
    "by_test_type": {
      "validation": 25,
      "negative": 18,
      "install": 12,
      "rbac": 10,
      "load": 8,
      "network": 6,
      "cleanup": 8
    },
    "by_impact_type": {
      "read-only": 45,
      "modifies-state": 32,
      "destructive": 2,
      "impairment": 6,
      "load-generation": 2
    },
    "production_safe_tests": 45,
    "cleanup_coverage": "100%"
  }
}
```

---

## Updated Sections

### Quality Checklist

Added new checkboxes for comprehensive testing:

**Comprehensive Testing (NEW):**
- [ ] Negative tests generated (30-40% ratio)
- [ ] Input validation tests cover all categories
- [ ] RBAC tests validate all permission boundaries
- [ ] Load tests for performance-critical components
- [ ] Network tests for distributed systems
- [ ] All tests use child_tests format

**Test Type Distribution:**
- [ ] Validation, Negative, RBAC, Load, Network, Setup/Install/Cleanup tests

**Filterable Attributes:**
- [ ] Every child test has test_type assigned
- [ ] Every child test has impact_type assigned
- [ ] RBAC-sensitive tests have rbac_level assigned
- [ ] Input validation tests have input_validation assigned
- [ ] Load tests have resource_requirements defined
- [ ] Network tests have network_requirements defined

### Success Metrics

Updated to include comprehensive testing metrics:

**Coverage Metrics:**
- 70%+ feature coverage
- 30-40% negative tests
- 100% input validation coverage

**Test Type Distribution (50 tests):**
- 20 validation (40%)
- 10 negative (20%)
- 8 input validation (16%)
- 5 RBAC (10%)
- 4 load (8%)
- 3 network (6%)
- 5 cleanup (10%)

**Quality Indicators:**
- All modifying tests have cleanup tests
- All negative tests define expected errors
- All RBAC tests specify minimum permission level
- All load tests define resource requirements
- All network tests define cleanup procedures
- 100% production safety classification

**New Format Adoption:**
- 100% use child_tests format
- Every child test has filterable attributes

### Tips for Comprehensive Coverage

Added guidance for new test types:

**Do:**
- Use child_tests format (not legacy steps)
- Generate negative tests (30-40% ratio)
- Validate all inputs (positive/negative/missing/corrupt/boundary)
- Test RBAC boundaries
- Include load tests for critical components
- Test network resilience

**Don't:**
- Use legacy steps format
- Skip negative tests
- Forget cleanup tests
- Omit filterable attributes
- Create load/network tests without cleanup

---

## Complete Phase Sequence

The skill now follows this comprehensive workflow:

1. **Artifact Analysis** - Explore primary artifact
2. **Coverage Matrix Design** - Build comprehensive coverage matrix
3. **Test Case Generation** - Generate core test cases
4. **Environment Variant Generation** - Create cloud/config variants
5. **Gap Analysis** - Identify coverage gaps
6. **System Impact and Safety Analysis** - Classify and create backup/cleanup tests
7. **Negative Testing Generation** ✨ NEW - Error handling validation
8. **Input Validation Testing** ✨ NEW - Comprehensive input coverage
9. **RBAC and Permission Testing** ✨ NEW - Security boundaries
10. **Load and Performance Testing** ✨ NEW - Resource constraints
11. **Network Testing with Impairment** ✨ NEW - Network resilience
12. **Test Generation Summary** ✨ NEW - Coverage report
13. **Output Format** - Generate JSON
14. **Quality Checklist** - Final validation

---

## Integration with Existing Workflow

The new phases integrate seamlessly:

**Execution Flow:**
1. Phases 1-6 generate base tests (happy path, environment variants, safety)
2. **Phase 7-8** add comprehensive negative and input validation tests
3. **Phase 9** adds RBAC security tests
4. **Phase 10-11** add load and network resilience tests
5. **Phase 12** summarizes coverage
6. Quality checklist validates completeness

**Backward Compatibility:**
- Skill can still generate tests without Phases 7-11 if needed
- Each phase is independent and can be skipped if not applicable
- Legacy test plans continue to work

---

## Example Test Generation

For a typical operator with 10 base features:

**Without new phases (old approach):**
- 10 positive tests
- 2-3 negative tests
- Total: ~13 tests

**With new phases (comprehensive approach):**
- 10 positive tests (validation)
- 8 negative tests (Phase 7: ~3 per feature × 3 features)
- 12 input validation tests (Phase 8: positive/negative/missing/corrupt/boundary × 2-3 features)
- 6 RBAC tests (Phase 9: cluster-admin/namespace-admin/edit/view operations)
- 4 load tests (Phase 10: resource limits, stress testing)
- 3 network tests (Phase 11: network policies, latency)
- 8 cleanup tests (for modifying operations)
- **Total: ~51 tests** (4x increase with comprehensive coverage)

**Coverage improvement:**
- Negative testing: 0-20% → 30-40%
- Input validation: Ad-hoc → 100% coverage
- RBAC testing: Minimal → All boundaries tested
- Load testing: None → Performance-critical components
- Network testing: None → Distributed systems validated

---

## Benefits

### For Test Plan Authors
- **Clear guidance** on generating comprehensive tests
- **Templates and examples** for all test types
- **Safety-first approach** with mandatory cleanup
- **Consistent structure** using child_tests format

### For Test Executors
- **Filterable tests** by type, impact, RBAC, validation
- **Production safety** clearly indicated
- **Resource requirements** defined upfront
- **Network impairment** documented with cleanup

### For System Quality
- **Error handling** validated through negative tests
- **Security boundaries** validated through RBAC tests
- **Performance** validated through load tests
- **Resilience** validated through network tests

---

## Files Modified

1. **skills/testplan-generator/skill.md**
   - Added 818 lines
   - Added Phases 7-12
   - Updated Quality Checklist
   - Updated Success Metrics
   - Updated Tips section

---

## Next Steps

### Immediate
1. ✅ Update testplan-generator skill (complete)
2. ⏳ Update testplan-educator skill (add guidance for new test types)
3. ⏳ Update testplan-executor skill (add execution patterns)

### Short Term
1. ⏳ Migrate CAMO example to new format
2. ⏳ Migrate RHOBS-Next example to new format
3. ⏳ Generate comprehensive tests for both examples

### Documentation
1. ⏳ Update examples/README.md with filtering examples
2. ⏳ Update skills/README.md with new test type descriptions
3. ⏳ Create comprehensive test generation guide

---

## Verification

To verify the skill update:

1. **Test with existing example:**
   ```
   Use testplan-generator skill to analyze examples/camo/
   Verify it generates comprehensive test coverage
   Check all test types are present
   ```

2. **Check test distribution:**
   - Should generate ~30-40% negative tests
   - Should have input validation for all inputs
   - Should have RBAC tests for permission boundaries
   - Should suggest load tests for critical components
   - Should suggest network tests for distributed components

3. **Validate format:**
   - All tests use child_tests (not steps)
   - All tests have filterable attributes
   - All modifying tests have cleanup tests
   - All tests have safety classification

---

**Status**: testplan-generator skill updated with comprehensive test generation ✅
**Next**: Update testplan-educator and testplan-executor skills
