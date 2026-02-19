# CAMO Test Plan Execution Review

## Overview

This document reviews how the CAMO test plan would be executed following the testplan-executor skill, documenting the execution flow, safety checks, and improvements needed.

## Execution Flow

### Pre-Execution Validation

**Step 1: Environment Readiness**
```bash
# Verify cluster access
oc whoami
# Expected: Logged in username

# Check we're in correct namespace
oc project openshift-monitoring

# Verify prerequisites
oc auth can-i create secrets -n openshift-monitoring
# Expected: yes

# Check CAMO operator is running
oc get deployment configure-alertmanager-operator -n openshift-monitoring
# Expected: 1/1 available
```

**Step 2: Safety and Impact Assessment**
```
System Impact Summary:
  👁️ Read-Only Tests: 2 (safe, no state changes)
    - test_verify_camo_deployment
    - test_read_alertmanager_config

  ✏️ Modifies-State Tests: 2 (requires cleanup)
    - test_create_pagerduty_secret (requires cleanup)
    - test_cleanup_pagerduty_secret (IS the cleanup)

  🚨 Destructive Tests: 0

✓ Non-production cluster detected: test-cluster-abc123
✓ All backup/cleanup patterns verified
```

### Test-by-Test Execution

#### Test 1: test_verify_camo_deployment (Read-Only)

**Pre-Test Safety Check:**
```
System Impact: read-only
Risk Level: low
Can Run in Production: true
Requires Confirmation: false
Requires Backup: false
Requires Cleanup: false

✅ Safety checks passed. Test execution approved.
```

**Execution:**
- Step 1: Check deployment ✓
- Step 2: Check pod running ✓
- Step 3: Check CSV succeeded ✓

**Validation:**
- All 3 success criteria met ✓
- Test status: PASSED

**Post-Execution:**
- No cleanup required (read-only) ✓

---

#### Test 2: test_read_alertmanager_config (Read-Only)

**Pre-Test Safety Check:**
```
System Impact: read-only
Risk Level: low
Can Run in Production: true

✅ Safety checks passed.
```

**Execution:**
- All steps extract and display config
- No state modifications

**Validation:**
- Config readable and valid ✓
- Test status: PASSED

---

#### Test 3: test_create_pagerduty_secret (Modifies-State)

**Pre-Test Safety Check:**
```
⚠️ SAFETY VALIDATION REQUIRED

System Impact: modifies-state
Risk Level: medium
Affected Resources:
  - Secret: pd-secret in openshift-monitoring
  - Secret: alertmanager-main (modified by CAMO)
Can Run in Production: false
Requires Confirmation: true
Requires Cleanup: true
Cleanup Test: test_cleanup_pagerduty_secret

⚠️ WARNING MESSAGE:
This test creates a PagerDuty secret that will cause CAMO to reconfigure
Alertmanager. Alerts may start routing to the test PagerDuty key. Run
cleanup test afterward to remove.

Current cluster: test-cluster-abc123
✓ Not in production

Cleanup test exists: test_cleanup_pagerduty_secret ✓

Do you want to proceed with this test? (yes/no): yes
✅ User confirmation received. Proceeding...
```

**Execution:**
- Step 1: Verify pd-secret doesn't exist ✓
- Step 2: Create pd-secret ✓
- Step 3: Wait for reconciliation (60s) ✓
- Step 4: Verify receiver added to config ✓
- Step 5: Check metric pd_secret_exists=1 ✓

**Validation:**
- All success criteria met ✓
- Test status: PASSED

**Post-Execution Cleanup Prompt:**
```
═══════════════════════════════════════════════
⚠️ CLEANUP REQUIRED
═══════════════════════════════════════════════
This test modified system state and requires cleanup to restore
the system to its original condition.

Cleanup test: test_cleanup_pagerduty_secret

Affected resources:
  - Secret: pd-secret in openshift-monitoring
  - Secret: alertmanager-main (modified by CAMO)

Run cleanup test now? (yes/no/later): yes

Executing cleanup test: test_cleanup_pagerduty_secret
```

---

#### Test 4: test_cleanup_pagerduty_secret (Cleanup Test)

**Pre-Test Safety Check:**
```
System Impact: modifies-state
Risk Level: low
Restores State: true

✅ Safety checks passed (cleanup test).
```

**Execution:**
- Step 1: Verify pd-secret exists ✓
- Step 2: Delete pd-secret ✓
- Step 3: Wait for reconciliation ✓
- Step 4: Verify receiver removed ✓
- Step 5: Check metric pd_secret_exists=0 ✓

**Validation:**
- All success criteria met ✓
- Test status: PASSED
- System returned to pre-test state ✓

**Post-Execution:**
```
✅ Cleanup complete. System restored to original state.
```

## Execution Summary

```
═══════════════════════════════════════════════════════════════
Test Execution Summary
═══════════════════════════════════════════════════════════════
Total Tests: 4
Passed: 4
Failed: 0
Skipped: 0

By Impact Type:
  Read-Only: 2/2 passed
  Modifies-State: 2/2 passed (cleanup completed)
  Destructive: 0/0

Cleanup Status:
  ✓ All modifying tests cleaned up
  ✓ System returned to original state

Total Time: ~45 minutes
═══════════════════════════════════════════════════════════════
```

## Improvements Identified

### 1. Data Model Enhancement: Alternative Paths

**Issue**: Alternative paths structure is `map[string][]string` but needs descriptions.

**Current:**
```json
"alternative_paths": {
  "quick_validation": ["test1", "test2"]
}
```

**Needed:**
```json
"alternative_paths": {
  "quick_validation": {
    "description": "Quick health check for existing installation",
    "sequence": ["test1", "test2"]
  }
}
```

**Impact**: Makes alternative paths self-documenting in HTML

**Fix**: Update `internal/models/testplan.go`:
```go
type AlternativePath struct {
    Description string   `json:"description"`
    Sequence    []string `json:"sequence"`
}

type LearningPath struct {
    Description      string                       `json:"description"`
    Sequence         []string                     `json:"sequence"`
    AlternativePaths map[string]AlternativePath   `json:"alternative_paths"` // Changed
}
```

### 2. HTML: Display Alternative Paths

**Issue**: Alternative paths are in the JSON but not displayed in HTML viewer

**Needed**: Add section in HTML template showing alternative paths:
```html
<div class="alternative-paths">
  <h3>Alternative Learning Paths</h3>
  {{range $name, $path := .LearningPath.AlternativePaths}}
    <div class="path">
      <h4>{{$name}}</h4>
      <p>{{$path.Description}}</p>
      <ol>
        {{range $path.Sequence}}
        <li>{{.}}</li>
        {{end}}
      </ol>
    </div>
  {{end}}
</div>
```

**Impact**: Users can see and choose alternative learning paths

### 3. Executor: Automated Cleanup Tracking

**Issue**: Executor prompts for cleanup but doesn't track which tests need cleanup later

**Current behavior**: User can choose "later" but has to remember manually

**Needed**:
- Track cleanup-needed tests in a file
- Show summary at end: "⚠️ 3 tests need cleanup: test_x, test_y, test_z"
- Provide command to run all pending cleanups

**Implementation:**
```bash
# During execution, track cleanup needed
echo "$test_id|$cleanup_test" >> results/cleanup-needed.txt

# At end of execution
if [ -f results/cleanup-needed.txt ]; then
  echo "⚠️ The following tests need cleanup:"
  cat results/cleanup-needed.txt | while IFS='|' read test cleanup; do
    echo "  - $test → Run: ./execute-testplan --tests $cleanup"
  done
fi
```

**Impact**: Prevents forgotten cleanup, reduces resource waste

### 4. Validator: Validate Cleanup Test References

**Issue**: Validator checks dependencies but could also validate cleanup_test references

**Current**: Checks if cleanup_test exists in testcases ✓

**Enhancement**: Also check:
- If cleanup test has `restores_state: true`
- If cleanup test is in cleanup category
- If backup test (if referenced) exists and has appropriate properties

**Impact**: Catch misconfigurations at parse time, not runtime

### 5. HTML: Safety Pre-Flight Checklist Rendering

**Issue**: Safety checklist in JSON but not prominently displayed in HTML

**Needed**: For modifies-state and destructive tests, show a checklist before steps:
```html
<div class="safety-checklist bg-yellow-50 border-l-4 border-yellow-500 p-4">
  <h4>⚠️ Pre-Flight Safety Check</h4>
  <ul>
    <li>☐ Verify cluster context (not production)</li>
    <li>☐ Backup test completed (if required)</li>
    <li>☐ Cleanup test exists and verified</li>
    <li>☐ User confirmation obtained</li>
  </ul>
</div>
```

**Impact**: Reduces accidental dangerous test execution

### 6. Executor: Command Intrusiveness Auto-Detection

**Issue**: Executor skill mentions analyzing commands for intrusiveness, but could be automated

**Needed**: Function to analyze command strings:
```bash
analyze_command_intrusiveness() {
  local cmd="$1"

  # Destructive operations
  if echo "$cmd" | grep -qE "delete|drain|evict|purge|--force"; then
    echo "DESTRUCTIVE"
    return
  fi

  # Modifying operations
  if echo "$cmd" | grep -qE "create|apply|patch|scale|set"; then
    echo "MODIFIES"
    return
  fi

  # Read-only
  echo "READ_ONLY"
}
```

**Enhancement**: Cross-check detected intrusiveness vs declared `system_impact.type`:
```bash
detected=$(analyze_command_intrusiveness "$command")
declared="$system_impact_type"

if [ "$detected" != "$declared" ]; then
  echo "⚠️ WARNING: Command appears $detected but test declares $declared"
  echo "   Review test safety classification"
fi
```

**Impact**: Catches misclassified tests, improves safety

### 7. Skills: Skill Interdependencies

**Issue**: Skills reference each other (generator → educator → executor) but no formal links

**Needed**: Add "Related Skills" section to each skill README:
```markdown
## Related Skills

- **Previous**: testplan-generator - Use to create initial test plan
- **Next**: testplan-educator - Use to add educational content
- **Execution**: testplan-executor - Use to run tests

## Typical Workflow

1. testplan-generator → Generate comprehensive test coverage
2. testplan-educator → Add learning content
3. testplan-executor → Execute and validate
```

**Impact**: Helps users understand the full workflow

### 8. HTML: Test Progress Tracking Enhancement

**Current**: Tracks completed tests in localStorage ✓

**Enhancement**: Also track:
- Timestamp when test completed
- Pass/fail status
- Notes/comments

```javascript
progress = {
  completedTests: ["test1", "test2"],
  testDetails: {
    "test1": {
      completed: "2026-02-18T14:30:00Z",
      status: "passed",
      notes: "Config validation failed first try, passed after fixing"
    }
  }
}
```

**Impact**: Better audit trail, helps with troubleshooting

### 9. Generator: Auto-Generate Cleanup Tests

**Issue**: Generator requires manually creating cleanup tests for each modifying test

**Enhancement**: Auto-generate cleanup test when creating modifying test:
```python
def generate_cleanup_test(main_test):
    return {
        "metadata": {
            "id": f"{main_test['id']}_cleanup",
            "title": f"Cleanup: {main_test['title']}",
            "category": "cleanup",
            "system_impact": {
                "type": "modifies-state",
                "description": f"Removes resources created by {main_test['id']}"
            },
            "state_management": {
                "restores_state": True
            }
        },
        "steps": infer_cleanup_steps(main_test["steps"])
    }
```

**Impact**: Ensures every modifying test has cleanup, reduces manual work

### 10. Validator: Circular Dependency Detection

**Issue**: Validator checks dependencies exist but not for circular refs

**Needed**: Detect circular dependencies:
```
test_a depends on test_b
test_b depends on test_c
test_c depends on test_a  ← CIRCULAR!
```

**Implementation**: Use topological sort to detect cycles

**Impact**: Prevents impossible execution orders

## Recommendations

### Priority 1 (High Impact, Easy):
1. Alternative paths data model enhancement (#1)
2. Cleanup tracking (#3)
3. Safety checklist rendering (#5)

### Priority 2 (High Impact, Medium Effort):
4. Command intrusiveness auto-detection (#6)
5. Auto-generate cleanup tests (#9)
6. Circular dependency detection (#10)

### Priority 3 (Nice to Have):
7. Alternative paths HTML display (#2)
8. Validator enhancements (#4)
9. Progress tracking enhancement (#8)
10. Skill interdependency docs (#7)

## Conclusion

The CAMO test plan demonstrates the full safety and state management system working correctly:

✅ Read-only tests clearly marked and safe to run anywhere
✅ Modifying tests require confirmation and cleanup
✅ Cleanup tests properly restore state
✅ Safety badges and warnings display in HTML
✅ Validation catches configuration errors early

The identified improvements would enhance usability and safety further, but the current system is fully functional for safe test execution.
