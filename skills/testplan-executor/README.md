# Test Plan Executor Skill

**Focus:** Test execution, results reporting, and gap identification

**Purpose:** Execute tests systematically, report meaningfully, and provide clear follow-up actions

## What This Skill Does

The Test Plan Executor skill executes test plans and provides comprehensive analysis:

- ⚙️ **Systematic execution** - Run tests in correct order with dependency handling
- 📊 **Result validation** - Compare actual vs expected, determine pass/fail
- 🔍 **Gap analysis** - Identify patterns, root causes, and missing coverage
- 📝 **Meaningful reports** - Executive summaries and detailed technical findings
- 🛠️ **Follow-up instructions** - Clear manual remediation steps with success criteria

## When to Use

Use this skill when you need to:

- Execute a test plan automatically
- Validate test plan accuracy
- Capture comprehensive test results
- Identify failure patterns and root causes
- Generate reports for stakeholders
- Create remediation action plans

## Inputs Required

- **Test plan JSON** (required) - The test plan to execute
- **Environment context** (required) - Cluster, credentials, region
- **Execution scope** (optional) - Which tests to run
  - All tests (default)
  - Specific tests: `test_1,test_3,test_5`
  - Resume from: `test_3:step_4`
- **Execution mode** (optional) - How to execute
  - `automated` (default) - Run all commands automatically
  - `interactive` - Prompt before each step
  - `dry-run` - Show what would be executed

## Example Usage

### Example 1: Full Automated Execution

```
Execute the complete test plan automatically.

Test plan: examples/observability-operator-testplan.json
Environment:
  Cluster: test-cluster-abc123
  OCP version: 4.15.1
  Cloud: AWS us-east-1
Mode: automated
```

**Claude will:**
1. Verify prerequisites and environment
2. Execute all 5 tests in sequence
3. Capture all outputs and timestamps
4. Validate against success criteria
5. Generate 3 reports (executive, detailed, follow-up)
6. Identify gaps and provide remediation steps
7. Output: `test-results/20260218_143000/`

### Example 2: Interactive Execution

```
Execute test plan with confirmation prompts.

Test plan: examples/sensitive-operations.json
Mode: interactive
```

**Claude will:**
1. Show each step before executing
2. Ask: "Execute this step? (y/n)"
3. Display output after execution
4. Ask: "Does output match expected? (y/n)"
5. Record user's assessment
6. Continue based on responses

### Example 3: Resume from Failure

```
Resume execution from where we left off yesterday.

Test plan: examples/observability-operator-testplan.json
Previous results: test-results/20260217_103000/results.json
Resume from: test_3:step_4
```

**Claude will:**
1. Load previous execution state
2. Skip completed tests (test_1, test_2)
3. Skip completed steps in test_3 (steps 1-3)
4. Resume at test_3, step 4
5. Continue through remaining tests

### Example 4: Dry-Run Mode

```
Show me what would be executed without actually running commands.

Test plan: examples/production-validation.json
Mode: dry-run
```

**Claude will:**
1. Show execution plan
2. Display each command that would run
3. Show expected outputs
4. Estimate total execution time
5. Identify potential issues
6. No actual execution

## Output Structure

### Executive Summary (REPORT_EXECUTIVE.md)
```markdown
# Executive Summary

| Metric | Value | Status |
|--------|-------|--------|
| Total Tests | 5 | - |
| Passed | 3 | ✅ 60% |
| Failed | 1 | ❌ 20% |
| Needs Review | 1 | ⚠️ 20% |

## Key Findings
✅ Strengths: Core functionality validated
❌ Issues: Prerequisites missing
⚠️ Needs Review: Output format differs

## Recommendations
1. Install missing components
2. Re-run failed tests
3. Review output formats
```

### Detailed Technical Report (REPORT_DETAILED.md)
```markdown
# Detailed Findings

## Test 1: Install Operator ✅ PASSED

### Steps Executed
| Step | Status | Duration |
|------|--------|----------|
| 1 | ✅ Passed | 2s |
| 2 | ✅ Passed | 1s |

### Validation Results
✅ All 4 criteria met

### Artifacts
- Logs: results/test_1_logs.txt

---

## Test 3: Create ServiceMonitor ❌ FAILED

### Failure Details
- Command: `oc apply -f servicemonitor.yaml`
- Exit code: 1
- Error: CRD not found

### Root Cause Analysis
ServiceMonitor CRD not installed

### Remediation Steps
[Detailed manual steps]
```

### Follow-Up Action Plan (REPORT_FOLLOWUP.md)
```markdown
# Follow-Up Action Plan

## Priority 1: Critical Blockers

### Action 1.1: Install Prometheus Operator

**Steps:**
1. Download bundle:
   ```bash
   curl -LO https://example.com/bundle.yaml
   ```

2. Apply bundle:
   ```bash
   oc create -f bundle.yaml
   ```

3. Verify:
   ```bash
   oc get crd | grep servicemonitors
   ```

**Success Criteria:**
- [ ] CRD exists
- [ ] Operator pod running

**Time:** 10 minutes
```

### Results JSON (results.json)
```json
{
  "execution_metadata": {
    "test_plan": "examples/observability.json",
    "execution_date": "2026-02-18T14:30:00Z",
    "environment": {
      "cluster": "test-cluster-abc123",
      "ocp_version": "4.15.1"
    }
  },
  "summary": {
    "total": 5,
    "passed": 3,
    "failed": 1,
    "needs_review": 1,
    "pass_rate": 60.0
  },
  "test_results": {
    "test_1": {
      "status": "passed",
      "duration_minutes": 28,
      "steps": [...]
    }
  },
  "gap_analysis": {
    "execution_gaps": [...],
    "coverage_gaps": [...],
    "failure_patterns": [...]
  }
}
```

## Integration with Other Skills

**Works well with:**
- **testplan-generator** - Execute generated comprehensive test plans
- **testplan-educator** - Execute educational test plans with learning context

**Typical workflow:**
1. **testplan-generator** → Generate comprehensive coverage
2. **testplan-educator** → Add educational content
3. **testplan-executor** → Execute and validate (this skill)

## Best Practices

### Do:
- Verify prerequisites before starting
- Capture all outputs (even on success)
- Record environmental context
- Identify patterns in failures
- Provide specific remediation steps
- Include time estimates
- Generate multiple report formats

### Don't:
- Skip prerequisite checks
- Continue blindly after critical failures
- Assume output matches without validation
- Generate reports without gap analysis
- Forget to preserve artifacts
- Provide vague recommendations

## Success Criteria

Successful execution includes:
- ✅ All tests executed or explicitly skipped
- ✅ 100% output capture
- ✅ Clear pass/fail/needs_review determination
- ✅ Root cause analysis for failures
- ✅ Prioritized gap analysis
- ✅ Actionable remediation plans
- ✅ Multiple report formats
- ✅ Preserved artifacts

## Files Generated

After execution:
```
test-results/YYYYMMDD_HHMMSS/
├── execution.log              # Complete log
├── results.json               # Machine-readable
├── REPORT_EXECUTIVE.md        # For management
├── REPORT_DETAILED.md         # For engineers
├── REPORT_FOLLOWUP.md         # Action plan
├── gap-analysis.json          # Identified gaps
└── artifacts/                 # Logs, screenshots
    ├── test_1_logs.txt
    └── test_3_error.txt
```

## Status Meanings

### Test Status
- **passed** - All steps successful, criteria met
- **failed** - Critical failure or criteria not met
- **needs_review** - Ambiguous result, human judgment needed
- **skipped** - Not executed due to blocker
- **aborted** - Stopped by user/critical error

### Step Status
- **passed** - Command succeeded, output matched
- **failed** - Command failed or output incorrect
- **needs_review** - Succeeded but output differs
- **skipped** - Not executed

## Next Steps

After using this skill:
1. Review executive summary
2. Address critical blockers
3. Investigate failures with detailed report
4. Follow remediation action plan
5. Re-execute failed tests
6. Update test plan based on learnings
