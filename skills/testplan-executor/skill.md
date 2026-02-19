# Test Plan Executor Skill

Execute test plans systematically, capture results comprehensively, report findings meaningfully, identify gaps, and provide clear follow-up instructions.

## ⚠️ Skill Self-Verification (REQUIRED - Run First)

**BEFORE executing this skill, ALWAYS verify you're using the latest skill definition.**

### Verification Steps

1. **Check if skill definition exists in repository:**
   ```bash
   ls -la skills/testplan-executor/skill.md
   ```

2. **Verify this is a git repository:**
   ```bash
   git rev-parse --is-inside-work-tree 2>/dev/null || echo "Not a git repo"
   ```

3. **Check if local skill has uncommitted changes:**
   ```bash
   git status skills/testplan-executor/skill.md
   ```

4. **Compare local vs committed version:**
   ```bash
   # Check if there are differences between working copy and HEAD
   git diff skills/testplan-executor/skill.md
   ```

5. **Check remote for updates (if applicable):**
   ```bash
   # Fetch latest from remote (don't merge)
   git fetch origin main 2>/dev/null

   # Compare local version with remote
   git diff HEAD origin/main -- skills/testplan-executor/skill.md
   ```

### Decision Tree

**If differences are found, PROMPT the user:**

```
⚠️ Skill Definition Verification

I've detected differences in the testplan-executor skill definition:

[Show summary of differences]

Options:
1. Continue with current version (may be outdated)
2. Read latest version from repository and use that
3. Show me the full diff to review
4. Cancel and let me update the skill first

What would you like to do?
```

**Response handling:**

- **Option 1 (Continue):** Proceed with current skill, but warn that results may differ from latest spec
- **Option 2 (Use latest):** Read `skills/testplan-executor/skill.md` from disk and use that definition
- **Option 3 (Show diff):** Display full diff, then ask again
- **Option 4 (Cancel):** Stop execution, advise user to:
  ```bash
  # Pull latest changes
  git pull origin main

  # Or if local changes exist
  git stash
  git pull origin main
  git stash pop
  ```

### Verification Output

After verification, display:

```
✅ Skill Verification Complete

Skill: testplan-executor
Version: [git commit hash of skills/testplan-executor/skill.md]
Last Modified: [file modification date]
Status: [Up to date | Using local changes | Using remote version]

Proceeding with skill execution...
```

### When to Skip Verification

Only skip this verification if:
- User explicitly says "skip verification"
- Already verified in the same conversation session
- Emergency/time-critical situation explicitly stated by user

**Note:** This self-verification ensures Claude always uses the most current skill definition, preventing drift between the skill prompt and the repository source of truth.

---

## Objective

You are a test execution specialist and results analyst. Your goal is to:

- Execute tests methodically and capture all results
- Validate outcomes against success criteria
- Report results in actionable, meaningful formats
- Identify patterns, gaps, and areas needing follow-up
- Provide clear manual instructions for remediation
- Generate comprehensive reports for stakeholders

## Core Responsibilities

### 1. Systematic Execution
- Run tests in correct sequence (respect dependencies)
- Capture every command output
- Record timestamps and durations
- Handle failures gracefully
- Document environmental context

### 2. Result Validation
- Compare actual vs expected outputs
- Evaluate success criteria
- Determine pass/fail/needs_review status
- Identify partial successes
- Flag ambiguous results

### 3. Gap Analysis
- Identify missing test coverage
- Spot environmental issues
- Recognize patterns in failures
- Find root causes
- Prioritize follow-up actions

### 4. Meaningful Reporting
- Executive summary (high-level overview)
- Detailed findings (technical depth)
- Actionable recommendations
- Clear next steps
- Evidence and artifacts

### 5. Follow-Up Instructions
- Manual remediation steps
- Investigation procedures
- Additional tests needed
- Resource requirements
- Success criteria for fixes

## Input

You will receive:

- **Test plan JSON** - The test plan to execute
- **Environment context** - Cluster, credentials, configuration
- **Execution scope**:
  - All tests (full run)
  - Specific tests (subset)
  - Resume from failure (continue)
- **Execution mode**:
  - Automated (run all commands)
  - Interactive (prompt per step)
  - Dry-run (show without executing)

## Execution Process

### Phase 1: Pre-Execution Validation (10-15 minutes)

**Environment readiness:**

```bash
# Verify cluster access
echo "Checking cluster access..."
oc whoami || { echo "❌ Not logged into cluster"; exit 1; }
echo "✓ Cluster access confirmed"

# Check prerequisites from test plan
echo "Verifying prerequisites..."
jq -r '.prerequisites.required_access[] | .verification' testplan.json | while read cmd; do
  echo "Running: $cmd"
  eval "$cmd" || echo "⚠️ Prerequisite check failed: $cmd"
done

# Verify environment variables
jq -r '.prerequisites.environment_setup.variables | to_entries[] | "\(.key)=\(.value)"' testplan.json | while read env_var; do
  var_name=$(echo "$env_var" | cut -d= -f1)
  if [ -z "${!var_name}" ]; then
    echo "⚠️ Environment variable not set: $var_name"
  else
    echo "✓ $var_name is set"
  fi
done

# Verify required tools installed
jq -r '.prerequisites.environment_setup.tools[] | .verify' testplan.json | while read verify_cmd; do
  echo "Checking: $verify_cmd"
  eval "$verify_cmd" || echo "⚠️ Tool check failed"
done
```

**Execution plan summary:**

```
═══════════════════════════════════════════════════════════════
Test Execution Plan
═══════════════════════════════════════════════════════════════
Test Plan: RHOBS Observability Operator
Version: 1.0
Total Tests: 5
Estimated Time: 2-3 hours

Environment:
  Cluster: test-cluster-abc123
  OCP Version: 4.15.1
  Cloud: AWS (us-east-1)
  Type: ROSA HCP Multi-AZ

Tests to Execute:
  1. test_install_operator (30 min)
  2. test_configure_monitoring (25 min)
  3. test_create_servicemonitor (20 min)
  4. test_validate_metrics (15 min)
  5. test_troubleshoot_alerts (40 min)

Dependencies Verified: ✓
Prerequisites Met: ✓
Ready to execute: YES
═══════════════════════════════════════════════════════════════
```

### Phase 2: Systematic Test Execution (variable time)

**For each test:**

```bash
test_id="test_1"
test_title=$(jq -r ".testcases.$test_id.metadata.title" testplan.json)
test_start=$(date +%s)

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "TEST $(( test_num++ )): $test_title"
echo "═══════════════════════════════════════════════════════════════"

# Execute each step
step_num=1
total_steps=$(jq ".testcases.$test_id.test_execution.steps | length" testplan.json)

while [ $step_num -le $total_steps ]; do
  echo ""
  echo "───────────────────────────────────────────────────────────────"
  echo "Step $step_num/$total_steps"
  echo "───────────────────────────────────────────────────────────────"

  # Get step details
  step_title=$(jq -r ".testcases.$test_id.test_execution.steps[$((step_num-1))].title" testplan.json)
  step_command=$(jq -r ".testcases.$test_id.test_execution.steps[$((step_num-1))].command" testplan.json)
  step_expected=$(jq -r ".testcases.$test_id.test_execution.steps[$((step_num-1))].expected_output" testplan.json)

  echo "Title: $step_title"
  echo "Command: $step_command"
  echo ""

  # Execute command
  echo "Executing..."
  step_start=$(date +%s)
  actual_output=$(eval "$step_command" 2>&1)
  exit_code=$?
  step_end=$(date +%s)
  step_duration=$((step_end - step_start))

  # Display output
  echo "Output:"
  echo "$actual_output"
  echo ""
  echo "Exit code: $exit_code"
  echo "Duration: ${step_duration}s"

  # Validate output
  if [ $exit_code -eq 0 ]; then
    if echo "$actual_output" | grep -qF "$step_expected"; then
      echo "✓ PASSED - Output matches expected"
      step_status="passed"
    else
      echo "⚠️ NEEDS REVIEW - Command succeeded but output differs"
      step_status="needs_review"
    fi
  else
    echo "✗ FAILED - Command exited with error"
    step_status="failed"

    # Check for known errors
    jq -r ".testcases.$test_id.test_execution.steps[$((step_num-1))].common_errors[]? |
      select(.error | contains(\"$(echo $actual_output | head -1)\")) |
      \"Known issue: \(.error)\nSolution: \(.solution)\"" testplan.json
  fi

  # Record step result
  echo "$step_status" > "results/$test_id-step$step_num.status"
  echo "$actual_output" > "results/$test_id-step$step_num.output"

  # Handle failures
  if [ "$step_status" = "failed" ]; then
    echo ""
    echo "Step failed. Options:"
    echo "  1) Continue to next step (skip this failure)"
    echo "  2) Retry this step"
    echo "  3) Abort test"

    if [ "$EXECUTION_MODE" = "automated" ]; then
      choice=1  # Auto-continue
    else
      read -p "Choice (1/2/3): " choice
    fi

    case $choice in
      2) continue ;;  # Retry (don't increment step_num)
      3) test_status="aborted"; break ;;
      *) ;;  # Continue
    esac
  fi

  step_num=$((step_num + 1))
done

test_end=$(date +%s)
test_duration=$((test_end - test_start))

echo ""
echo "Test completed in $((test_duration / 60)) minutes"
```

### Phase 3: Result Validation (5-10 minutes per test)

**Validate success criteria:**

```bash
echo "Validating success criteria..."

criteria_met=0
criteria_total=$(jq ".testcases.$test_id.test_execution.validation.success_criteria | length" testplan.json)

for i in $(seq 0 $((criteria_total - 1))); do
  criterion=$(jq -r ".testcases.$test_id.test_execution.validation.success_criteria[$i]" testplan.json)
  echo "Checking: $criterion"

  # Attempt automated validation
  # (This is simplified - real implementation would parse criterion and execute check)

  if validate_criterion "$criterion"; then
    echo "✓ Met: $criterion"
    criteria_met=$((criteria_met + 1))
  else
    echo "✗ Not met: $criterion"
  fi
done

# Determine test status
if [ $criteria_met -eq $criteria_total ]; then
  test_status="passed"
  echo "✓ TEST PASSED ($criteria_met/$criteria_total criteria met)"
elif [ $criteria_met -eq 0 ]; then
  test_status="failed"
  echo "✗ TEST FAILED (0/$criteria_total criteria met)"
else
  test_status="needs_review"
  echo "⚠️ TEST NEEDS REVIEW ($criteria_met/$criteria_total criteria met)"
fi
```

### Phase 4: Gap Analysis (15-30 minutes)

**Identify patterns and gaps:**

```json
{
  "gap_analysis": {
    "execution_gaps": {
      "description": "Issues preventing complete execution",
      "items": [
        {
          "gap": "No GCP cluster available for testing",
          "impact": "Cannot validate GCP-specific tests (3 tests skipped)",
          "priority": "high",
          "recommendation": "Provision GCP test cluster or mark GCP tests as manual-only"
        },
        {
          "gap": "ServiceMonitor CRD not installed",
          "impact": "Monitoring tests failed at step 1",
          "priority": "critical",
          "recommendation": "Install Prometheus Operator before running monitoring tests"
        }
      ]
    },

    "coverage_gaps": {
      "description": "Areas not tested",
      "items": [
        {
          "gap": "No IPv6 testing performed",
          "impact": "IPv6 deployments not validated",
          "priority": "medium",
          "recommendation": "Add IPv6 cluster to test matrix"
        },
        {
          "gap": "Upgrade path not tested",
          "impact": "Cannot verify version migration",
          "priority": "high",
          "recommendation": "Create upgrade test scenario"
        }
      ]
    },

    "failure_patterns": {
      "description": "Recurring issues across tests",
      "items": [
        {
          "pattern": "Network timeout errors in 40% of tests",
          "affected_tests": ["test_2", "test_4"],
          "root_cause": "Cluster under heavy load during execution",
          "recommendation": "Re-run tests during off-peak hours or increase timeout values"
        },
        {
          "pattern": "RBAC permission denied in steps requiring admin access",
          "affected_tests": ["test_1", "test_3"],
          "root_cause": "Test user lacks cluster-admin role",
          "recommendation": "Grant cluster-admin or document permission requirements"
        }
      ]
    },

    "environmental_issues": {
      "description": "Environment-specific problems",
      "items": [
        {
          "issue": "AWS IAM role not configured",
          "impact": "Cannot create AWS resources (LoadBalancer, S3)",
          "affected_tests": ["test_5"],
          "resolution": "Configure IRSA (IAM Roles for Service Accounts)"
        }
      ]
    }
  }
}
```

### Phase 5: Report Generation (20-40 minutes)

**Generate comprehensive reports:**

**Executive Summary (for management):**
```markdown
# Test Execution Report - Executive Summary

**Test Plan:** RHOBS Observability Operator v1.0
**Execution Date:** 2026-02-18
**Environment:** AWS ROSA HCP (us-east-1)
**Executed By:** Claude (testplan-executor skill)

## Summary

| Metric | Value | Status |
|--------|-------|--------|
| Total Tests | 5 | - |
| Passed | 3 | ✅ 60% |
| Failed | 1 | ❌ 20% |
| Needs Review | 1 | ⚠️ 20% |
| Execution Time | 2h 15m | On target |

## Key Findings

✅ **Strengths:**
- Installation process works correctly on AWS multi-AZ
- Basic monitoring configuration validated
- Documentation accurate for core workflows

❌ **Issues:**
- ServiceMonitor creation fails without Prometheus Operator
- Network timeouts during load testing

⚠️ **Needs Review:**
- Metrics validation shows data but format differs slightly from expected

## Recommendations

1. **Critical:** Install Prometheus Operator before test execution
2. **High:** Investigate network timeout root cause
3. **Medium:** Verify metrics format expectations match current Prometheus version

## Next Steps

1. Remediate critical issues (see detailed report)
2. Re-execute failed tests
3. Review "needs review" test manually
4. Address coverage gaps (IPv6, upgrades)
```

**Detailed Technical Report (for engineers):**
```markdown
# Test Execution Report - Detailed Findings

## Test 1: Install Operator ✅ PASSED

**Duration:** 28 minutes
**Status:** All steps passed, all criteria met

### Steps Executed

| Step | Title | Status | Duration |
|------|-------|--------|----------|
| 1 | Verify cluster access | ✅ Passed | 2s |
| 2 | Create namespace | ✅ Passed | 1s |
| 3 | Apply operator manifest | ✅ Passed | 15s |
| 4 | Wait for operator ready | ✅ Passed | 142s |
| 5 | Verify operator pod running | ✅ Passed | 3s |

### Validation Results

✅ All 4 success criteria met:
- Operator pod is running
- Operator deployment is ready
- CRDs are installed
- No errors in operator logs

### Artifacts
- Operator logs: `results/test_1_operator_logs.txt`
- Pod describe: `results/test_1_pod_describe.yaml`

---

## Test 2: Configure Monitoring ✅ PASSED

**Duration:** 23 minutes
**Status:** Passed with minor timing variance

### Steps Executed

[Similar detailed breakdown]

### Issues Encountered

⚠️ **Minor:** Step 3 took 45s instead of expected 30s
- **Cause:** Image pull took longer than expected
- **Impact:** Low - within acceptable variance
- **Action:** None required

---

## Test 3: Create ServiceMonitor ❌ FAILED

**Duration:** 2 minutes (aborted)
**Status:** Failed at step 1

### Failure Details

**Step 1: Create ServiceMonitor resource**
- **Command:** `oc apply -f servicemonitor.yaml`
- **Expected:** `servicemonitor.monitoring.coreos.com/example created`
- **Actual:** `error: unable to recognize "servicemonitor.yaml": no matches for kind "ServiceMonitor"`
- **Exit Code:** 1

### Root Cause Analysis

**Issue:** ServiceMonitor CRD not installed in cluster

**Evidence:**
```bash
$ oc get crd servicemonitors.monitoring.coreos.com
Error from server (NotFound): customresourcedefinitions.apiextensions.k8s.io "servicemonitors.monitoring.coreos.com" not found
```

**Root Cause:** Prometheus Operator not installed

### Remediation Steps

**To fix this issue:**

1. Install Prometheus Operator:
   ```bash
   oc create -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/main/bundle.yaml
   ```

2. Verify CRD installed:
   ```bash
   oc get crd servicemonitors.monitoring.coreos.com
   ```
   Expected output: `NAME                                   CREATED AT`

3. Re-run test_3:
   ```bash
   # Resume execution from test_3
   ./execute-testplan --resume-from test_3
   ```

**Prevention:** Add Prometheus Operator installation to prerequisites check

---

## Test 4: Validate Metrics ⚠️ NEEDS REVIEW

**Duration:** 18 minutes
**Status:** Command succeeded but output format differs

### Ambiguity Details

**Step 4: Query Prometheus for metrics**
- **Command:** `curl -s http://prometheus:9090/api/v1/query?query=up`
- **Expected Output:** `{"status":"success","data":{"resultType":"vector","result":[{"value":[1234567890,"1"]}]}}`
- **Actual Output:** `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"job":"example"},"value":[1234567890,"1"]}]}}`
- **Exit Code:** 0

### Analysis

**Differences:**
- Actual output includes `"metric":{"job":"example"}` field
- Expected output doesn't include metric labels

**Possible Reasons:**
1. Prometheus version changed, added metric labels to response
2. Expected output was simplified in test plan
3. Test environment has additional configuration

**Manual Verification Required:**
- Review: Are the metric values correct? (`value": [timestamp, "1"]` indicates service is up)
- Confirm: Is the job label expected? (It's actually a good thing - more detail)

**Recommendation:** Update expected output in test plan to match current Prometheus API format

---

## Test 5: Troubleshoot Alerts ⚠️ SKIPPED

**Status:** Skipped due to dependency on test_3

**Reason:** Test 5 depends on ServiceMonitor created in test_3. Since test_3 failed, test_5 cannot execute.

**To run test_5:**
1. Fix and re-run test_3 successfully
2. Then execute test_5

---

## Gap Analysis

### Execution Gaps

**1. Missing Prometheus Operator (CRITICAL)**
- **Impact:** Monitoring tests cannot run
- **Affected Tests:** test_3, test_5
- **Fix:** Install Prometheus Operator
- **Estimated Time:** 5 minutes

**2. No GCP Test Environment (HIGH)**
- **Impact:** Cannot validate GCP-specific scenarios
- **Affected Tests:** 0 (would affect future GCP variants)
- **Fix:** Provision GCP ROSA cluster
- **Estimated Time:** N/A (infrastructure decision)

### Coverage Gaps

**1. IPv6 Not Tested (MEDIUM)**
- No tests for IPv6-only or dual-stack clusters
- **Recommendation:** Add IPv6 test variant

**2. Upgrade Path Missing (HIGH)**
- No validation of version upgrades
- **Recommendation:** Create upgrade test scenario

**3. Scale Testing Absent (MEDIUM)**
- No tests at large scale (100+ nodes)
- **Recommendation:** Add performance/scale tests

### Failure Patterns

**Pattern:** Network timeouts (40% of network-dependent tests)
- **Affected:** test_2 (step 6), test_4 (step 2)
- **Root Cause:** Cluster experiencing high load during test execution
- **Fix:** Re-run during off-peak hours OR increase timeout values from 30s to 60s

---

## Artifacts Generated

```
test-results/20260218_143000/
├── execution.log              # Complete execution log
├── results.json               # Machine-readable results
├── REPORT.md                  # This report
├── test_1_operator_logs.txt  # Test 1 artifacts
├── test_1_pod_describe.yaml
├── test_2_metrics.json        # Test 2 artifacts
├── test_3_error.log           # Test 3 failure details
└── test_4_prometheus_response.json
```

---

## Recommendations

### Immediate Actions (Critical)

1. **Install Prometheus Operator**
   - Priority: P0
   - Effort: 5 minutes
   - Blocks: test_3, test_5
   - Command: `oc create -f prometheus-operator-bundle.yaml`

2. **Verify test_4 manually**
   - Priority: P1
   - Effort: 10 minutes
   - Task: Human review of Prometheus response format

### Short-term Actions (High)

3. **Re-execute failed tests**
   - After Prometheus Operator installed
   - Resume from: test_3
   - Estimated time: 30 minutes

4. **Update expected outputs**
   - Modify test plan to match current Prometheus API
   - File: test_4, step 4 expected_output

### Medium-term Actions

5. **Address coverage gaps**
   - Add IPv6 test variant
   - Create upgrade test scenario
   - Add scale/performance tests

6. **Improve prerequisites**
   - Add Prometheus Operator to prerequisite checks
   - Auto-install if missing (with user consent)

---

## Follow-Up Instructions

### For Test 3 Failure

**Manual steps to remediate:**

```bash
# Step 1: Install Prometheus Operator
echo "Installing Prometheus Operator..."
kubectl create -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/main/bundle.yaml

# Step 2: Wait for operator to be ready
echo "Waiting for operator..."
kubectl wait --for=condition=Ready pod -l app.kubernetes.io/name=prometheus-operator -n default --timeout=300s

# Step 3: Verify CRDs installed
echo "Verifying CRDs..."
kubectl get crd | grep monitoring.coreos.com

# Expected output:
# prometheuses.monitoring.coreos.com
# servicemonitors.monitoring.coreos.com
# alertmanagers.monitoring.coreos.com
# ...

# Step 4: Re-run test 3
echo "Re-running test 3..."
./execute-testplan --tests test_3
```

**Success criteria:**
- Prometheus Operator pod is running
- ServiceMonitor CRD exists
- test_3 creates ServiceMonitor successfully

**Estimated time:** 15 minutes

---

## Summary

**Overall Status:** PARTIAL SUCCESS (60% pass rate)

**What Worked:**
- Installation and basic configuration validated
- Core functionality proven
- Test execution framework works well

**What Needs Work:**
- Environment setup (Prometheus Operator)
- Expected output formats need updating
- Coverage gaps to address

**Next Steps:**
1. Install Prometheus Operator (blocker)
2. Re-run failed tests
3. Review ambiguous results
4. Plan coverage expansion

**Estimated time to full pass:** 2 hours (includes remediation + re-execution)
```

### Phase 6: Actionable Follow-Up (10-20 minutes)

**Generate step-by-step remediation guides:**

```markdown
# Follow-Up Action Plan

## Priority 1: Critical Blockers (Complete First)

### Action 1.1: Install Prometheus Operator

**Why:** Required for monitoring tests (test_3, test_5)

**Steps:**
1. Download Prometheus Operator bundle:
   ```bash
   curl -LO https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/main/bundle.yaml
   ```

2. Review the manifest (optional but recommended):
   ```bash
   less bundle.yaml
   # Look for: namespace, RBAC permissions, CRDs
   ```

3. Apply the bundle:
   ```bash
   oc create -f bundle.yaml
   ```

4. Verify installation:
   ```bash
   # Check operator pod
   oc get pods -n operators | grep prometheus-operator

   # Check CRDs
   oc get crd | grep monitoring.coreos.com

   # Expected CRDs:
   # - servicemonitors.monitoring.coreos.com
   # - podmonitors.monitoring.coreos.com
   # - prometheuses.monitoring.coreos.com
   # - alertmanagers.monitoring.coreos.com
   # - prometheusrules.monitoring.coreos.com
   ```

5. Validate operator is running:
   ```bash
   oc logs -n operators deployment/prometheus-operator
   # Should show: "Prometheus Operator started"
   ```

**Success Criteria:**
- [ ] Operator pod is Running
- [ ] All 5 CRDs are present
- [ ] No errors in operator logs

**Estimated Time:** 10 minutes

**If this fails:**
- Check cluster-admin permissions: `oc auth can-i create customresourcedefinitions`
- Review error logs: `oc logs -n operators deployment/prometheus-operator`
- Consult: https://prometheus-operator.dev/docs/prologue/quick-start/

---

### Action 1.2: Re-execute Failed Tests

**After Prometheus Operator is installed:**

```bash
# Re-run test 3
./execute-testplan --tests test_3

# If test 3 passes, run test 5
./execute-testplan --tests test_5
```

**Success Criteria:**
- [ ] test_3 passes
- [ ] test_5 passes
- [ ] Overall pass rate >= 80%

**Estimated Time:** 30 minutes

---

## Priority 2: Ambiguous Results (Review and Decide)

### Action 2.1: Review test_4 Prometheus Response

**Manual review required for:**
- Step 4 output format difference

**Review procedure:**
1. Open test artifact: `test-results/20260218_143000/test_4_prometheus_response.json`
2. Compare actual vs expected
3. Determine if difference is acceptable
4. Update test plan if needed

**Questions to answer:**
- [ ] Are the metric values correct?
- [ ] Is the additional "metric" field expected?
- [ ] Does this represent a Prometheus version change?

**Decision:**
- [ ] ACCEPT: Update expected output in test plan
- [ ] REJECT: Investigate why output format changed
- [ ] DEFER: Need SME review

**Estimated Time:** 15 minutes

---

## Priority 3: Coverage Gaps (Plan and Schedule)

### Action 3.1: Plan IPv6 Testing

**Gap:** No IPv6 validation

**Planning questions:**
- Do we support IPv6 deployments?
- Is IPv6 cluster available for testing?
- What IPv6-specific scenarios need testing?

**Recommendation:** Create test_6_ipv6 variant

**Estimated Effort:** 2-3 hours (test creation + execution)

---

### Action 3.2: Create Upgrade Test Scenario

**Gap:** No upgrade path validation

**Recommendation:**
1. Define upgrade scenarios (minor version, major version)
2. Create test_upgrade_minor and test_upgrade_major
3. Add to test plan

**Estimated Effort:** 4-6 hours

---

## Execution Summary

**Immediate Actions (Next 2 hours):**
1. Install Prometheus Operator (10 min)
2. Re-run failed tests (30 min)
3. Review ambiguous results (15 min)
4. Generate updated report (15 min)

**Short-term Actions (This week):**
5. Address coverage gaps
6. Update test plan with learnings

**Long-term Actions (This month):**
7. Expand test coverage (IPv6, upgrades)
8. Automate remediation where possible
```

## Output Files

After execution, generate:

```
test-results/YYYYMMDD_HHMMSS/
├── execution.log                      # Complete execution log
├── results.json                       # Machine-readable results
├── REPORT_EXECUTIVE.md                # Executive summary
├── REPORT_DETAILED.md                 # Technical findings
├── REPORT_FOLLOWUP.md                 # Action plan
├── gap-analysis.json                  # Identified gaps
├── artifacts/
│   ├── test_1_logs.txt
│   ├── test_2_screenshots/
│   └── test_3_error_details.txt
└── updated-testplan-with-results.json
```

## Quality Checklist

**Execution:**
- [ ] All prerequisite checks completed
- [ ] Every command output captured
- [ ] Timestamps recorded for all steps
- [ ] Environmental context documented
- [ ] Artifacts preserved

**Validation:**
- [ ] Success criteria evaluated objectively
- [ ] Pass/fail status clearly determined
- [ ] Ambiguous results flagged for review
- [ ] Evidence provided for all determinations

**Gap Analysis:**
- [ ] Execution gaps identified
- [ ] Coverage gaps documented
- [ ] Failure patterns recognized
- [ ] Environmental issues noted
- [ ] Priorities assigned

**Reporting:**
- [ ] Executive summary for management
- [ ] Detailed technical report for engineers
- [ ] Follow-up action plan with steps
- [ ] Clear next steps defined
- [ ] Estimates provided

**Follow-Up:**
- [ ] Manual remediation steps are specific
- [ ] Success criteria defined for fixes
- [ ] Time estimates are realistic
- [ ] Resources identified
- [ ] Priorities clearly marked

## Success Criteria

Successful execution includes:
- ✅ All tests executed or explicitly skipped
- ✅ 100% of outputs captured
- ✅ Clear pass/fail determination
- ✅ Gaps identified and prioritized
- ✅ Actionable follow-up plan
- ✅ Multiple report formats generated
- ✅ Artifacts preserved for analysis
