# Comprehensive Testing Phases for testplan-generator Skill

**NOTE**: This content should be added to `skills/testplan-generator/skill.md` as Phase 7-11, inserted before the "Quality Checklist" section (line 737).

---

## Phase 7: Negative Testing Generation (30-45 minutes)

**REQUIRED:** For every positive/happy-path test, generate corresponding negative test cases to validate error handling and system resilience.

### What is Negative Testing?

Negative testing validates that the system correctly handles invalid inputs, error conditions, and unexpected scenarios. These tests should **expect failures** and verify appropriate error messages and system behavior.

### Negative Test Categories

**1. Invalid Configuration**
**2. Permission Denials**
**3. Resource Conflicts**
**4. Missing Dependencies**

### Child Test Case Format

All negative tests use the new `child_tests` format with:
- `test_type: "negative"`
- `impact_type`: usually "read-only" (fails before making changes)
- `input_validation`: "negative", "missing", "corrupt", or "boundary"
- `expected_output`: Specific error message expected

[Full examples in FILTERABLE-TESTCASES-IMPLEMENTATION.md]

---

## Phase 8: Input Validation Testing (30-45 minutes)

**REQUIRED:** Generate comprehensive input validation tests covering positive, negative, missing, corrupt, and boundary cases.

### Input Validation Categories

1. **Positive** (`input_validation: "positive"`) - Valid inputs
2. **Negative** (`input_validation: "negative"`) - Invalid types/values
3. **Missing** (`input_validation: "missing"`) - Required fields omitted
4. **Corrupt** (`input_validation: "corrupt"`) - Malformed data
5. **Boundary** (`input_validation: "boundary"`) - Min/max limits

[Full examples in FILTERABLE-TESTCASES-IMPLEMENTATION.md]

---

## Phase 9: RBAC and Permission Testing (30-45 minutes)

**REQUIRED:** Generate tests to validate role-based access control and security boundaries.

### RBAC Test Categories

1. **Cluster-Admin** (`rbac_level: "cluster-admin"`) - Cluster-wide operations
2. **Namespace-Admin** (`rbac_level: "namespace-admin"`) - Namespace-scoped admin
3. **Edit** (`rbac_level: "edit"`) - Standard edit operations
4. **View-Only** (`rbac_level: "view"`) - Read-only operations
5. **Permission Denials** - Validate RBAC denials

[Full examples in FILTERABLE-TESTCASES-IMPLEMENTATION.md]

---

## Phase 10: Load and Performance Testing (45-60 minutes)

**REQUIRED:** Generate tests to validate system behavior under load and resource constraints.

### Load Testing Categories

1. **Resource Limit Tests** - Test behavior when hitting resource limits
2. **Load Generation Tests** - Deploy load-generating pods
3. **Performance Testing** - Measure performance metrics

### Resource Requirements

```json
{
  "resource_requirements": {
    "cpu": "2",
    "memory": "512Mi",
    "load_generation": true
  }
}
```

[Full examples in FILTERABLE-TESTCASES-IMPLEMENTATION.md]

---

## Phase 11: Network Testing with Impairment (45-60 minutes)

**REQUIRED:** Generate tests to validate system behavior under network stress and impairment.

### Network Testing Categories

1. **Network Policy Tests** - Restrict traffic via NetworkPolicy
2. **Latency Injection** - Add artificial latency
3. **Packet Loss** - Simulate packet loss
4. **Security Group Tests** - AWS security group modifications

### Network Requirements

```json
{
  "network_requirements": {
    "impairment_needed": true,
    "impairment_type": "latency",
    "impairment_config": "100ms delay on port 443",
    "network_policies": ["deny-all-policy.yaml"]
  }
}
```

[Full examples in FILTERABLE-TESTCASES-IMPLEMENTATION.md]

---

## Integration with Existing Phases

These new phases (7-11) should be executed AFTER Phase 6 (Safety Analysis) and BEFORE Quality Checklist:

**Complete Phase Sequence:**

1. Artifact Analysis
2. Coverage Matrix Design
3. Test Case Generation
4. Environment Variant Generation
5. Gap Analysis
6. System Impact and Safety Analysis
7. **Negative Testing Generation** (NEW)
8. **Input Validation Testing** (NEW)
9. **RBAC and Permission Testing** (NEW)
10. **Load and Performance Testing** (NEW)
11. **Network Testing with Impairment** (NEW)
12. Quality Checklist
13. Final JSON Output

---

## Child Test Case vs Legacy Steps

All new test types should use the `child_tests` format instead of legacy `steps`:

### Legacy Format (Deprecated):
```json
{
  "test_execution": {
    "steps": [{"step_number": 1, "title": "...", "command": "..."}]
  }
}
```

### New Format (Preferred):
```json
{
  "test_execution": {
    "child_tests": [{
      "id": "test_001_step_001",
      "parent_test_id": "test_001",
      "step_number": 1,
      "title": "...",
      "test_type": "validation",
      "impact_type": "read-only",
      "rbac_level": "view",
      "input_validation": "positive",
      "resource_requirements": {},
      "network_requirements": {},
      "command": "...",
      "safety": {},
      "validation": {}
    }]
  }
}
```

---

## Test Generation Targets

For a comprehensive test plan, aim to generate:

- **30-40% negative tests** (relative to positive tests)
- **100% input validation coverage** for all inputs
- **RBAC tests for all permission boundaries**
- **Load tests for performance-critical components**
- **Network tests for distributed systems**

**Example Distribution** (50 total tests):
- 20 positive/validation tests
- 10 negative tests
- 8 input validation tests
- 5 RBAC tests
- 4 load tests
- 3 network tests
- 5 cleanup tests

---

## File: skills/testplan-generator/skill.md Update Location

Insert this content at **line 737**, before:
```markdown
## Quality Checklist

Before finalizing test plan:
```

---

## Implementation Status

- ✅ Data model supports all test types (Phase 1 complete)
- ✅ HTML viewer filters all test types (Phase 2 complete)
- ⏳ Skill generates comprehensive tests (Phase 3 in progress)
- ⏳ Examples migrated to new format (Phase 4 pending)

---

## References

- **Implementation Summary**: FILTERABLE-TESTCASES-IMPLEMENTATION.md
- **Planning Document**: FILTERABLE-TESTCASES-PLAN.md
- **Data Model**: internal/models/testplan.go
- **HTML Template**: internal/generator/templates/testplan.html
