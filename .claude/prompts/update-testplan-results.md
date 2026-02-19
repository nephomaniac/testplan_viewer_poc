# Update Test Plan Results

## Objective

Update a test plan JSON file with execution results from a test run. Merge execution data (pass/fail status, actual outputs, timestamps, issues) into the original test plan to create a comprehensive test report.

## Your Role

You are a test results aggregator and reporter. Your goal is to:

1. **Load original test plan** - Read the base test plan JSON
2. **Load execution results** - Read results from execute-testplan.md
3. **Merge data** - Combine original test plan with execution results
4. **Add execution metadata** - Timestamp, environment, executor info
5. **Generate updated JSON** - Output test plan with results embedded
6. **Create summary report** - Human-readable test report

## Inputs You'll Receive

The user will provide:

- **Original test plan JSON**: The test plan that was executed
- **Execution results JSON**: Results from execute-testplan.md
- **Optional: Execution logs**: Log files for additional context
- **Optional: Screenshots/artifacts**: Files to reference in results

## Update Process

### 1. Load Input Files (2-3 minutes)

**Read original test plan:**

```bash
# Load and validate original test plan
cat examples/rhobs_test_plan_v2.json | jq . > /dev/null
if [ $? -eq 0 ]; then
  echo "✓ Original test plan loaded"
  original_plan=$(cat examples/rhobs_test_plan_v2.json)
else
  echo "✗ Invalid JSON in original test plan"
  exit 1
fi
```

**Read execution results:**

```bash
# Load execution results
cat test-results/20260218_103000/results.json | jq . > /dev/null
if [ $? -eq 0 ]; then
  echo "✓ Execution results loaded"
  execution_results=$(cat test-results/20260218_103000/results.json)
else
  echo "✗ Invalid JSON in execution results"
  exit 1
fi
```

**Verify compatibility:**

```bash
# Check that test IDs match
original_test_ids=$(echo "$original_plan" | jq -r '.testcases | keys[]' | sort)
results_test_ids=$(echo "$execution_results" | jq -r '.test_results | keys[]' | sort)

if [ "$original_test_ids" = "$results_test_ids" ]; then
  echo "✓ Test IDs match between plan and results"
else
  echo "⚠ Test ID mismatch - some tests may not have results"
fi
```

### 2. Merge Results into Test Plan (10-20 minutes)

**Create updated structure:**

The updated test plan maintains the original structure but adds execution data:

```json
{
  "metadata": {
    // Original metadata
    "document_title": "...",
    "version": "...",
    // NEW: Execution information
    "last_executed": "2026-02-18T10:30:00Z",
    "execution_history": [
      {
        "execution_id": "exec-20260218-103000",
        "date": "2026-02-18T10:30:00Z",
        "executor": "claude-sonnet-4.5",
        "summary": {
          "total": 3,
          "passed": 2,
          "failed": 1,
          "pass_rate": 66.67
        },
        "environment": {
          "cluster": "test-cluster-abc123",
          "ocp_version": "4.15.1"
        }
      }
    ]
  },
  "testcases": {
    "test_1": {
      // Original test case definition
      "metadata": { ... },
      "learning": { ... },
      "test_execution": {
        // Original execution plan
        "objective": "...",
        "steps": [
          {
            "step_number": 1,
            "title": "...",
            "command": "...",
            "expected_output": "...",
            // NEW: Execution results for this step
            "execution_result": {
              "status": "passed",
              "actual_output": "...",
              "exit_code": 0,
              "execution_time_seconds": 2,
              "timestamp": "2026-02-18T10:30:05Z",
              "matched_expected": true
            }
          }
        ],
        "validation": {
          "success_criteria": ["..."],
          // NEW: Validation results
          "validation_result": {
            "criteria_met": 3,
            "criteria_total": 3,
            "all_passed": true,
            "details": [
              {
                "criterion": "Cluster appears in 'rosa list clusters'",
                "met": true,
                "verification_output": "..."
              }
            ]
          }
        }
      },
      // NEW: Overall test result
      "test_result": {
        "execution_id": "exec-20260218-103000",
        "status": "passed",
        "started_at": "2026-02-18T10:30:00Z",
        "completed_at": "2026-02-18T11:12:00Z",
        "execution_time_minutes": 42,
        "executor": "claude-sonnet-4.5",
        "environment": {
          "cluster": "test-cluster-abc123",
          "ocp_version": "4.15.1"
        },
        "issues_encountered": [],
        "notes": "Test completed successfully",
        "artifacts": {
          "logs": ["test-results/20260218_103000/test1.log"],
          "screenshots": [],
          "files": []
        }
      }
    }
  }
}
```

**Merge algorithm:**

```bash
# Use jq to merge results into original plan
jq --slurpfile results test-results/results.json '
  # Add execution metadata
  .metadata.last_executed = $results[0].execution_metadata.execution_date |
  .metadata.execution_history += [{
    execution_id: $results[0].execution_metadata.execution_id,
    date: $results[0].execution_metadata.execution_date,
    executor: $results[0].execution_metadata.executor,
    summary: $results[0].summary,
    environment: $results[0].execution_metadata.environment
  }] |

  # Merge test results
  .testcases |= with_entries(
    .value.test_result = $results[0].test_results[.key] |

    # Merge step results
    .value.test_execution.steps |= [
      .[] | . + {
        execution_result: (
          $results[0].test_results[.key].steps[] |
          select(.step_number == .step_number)
        )
      }
    ] |

    # Merge validation results
    .value.test_execution.validation.validation_result =
      $results[0].test_results[.key].validation
  )
' examples/rhobs_test_plan_v2.json > examples/rhobs_test_plan_v2_results.json
```

### 3. Add Execution Summary (5-10 minutes)

**Create execution summary section:**

```json
{
  "execution_summary": {
    "latest_execution": {
      "execution_id": "exec-20260218-103000",
      "date": "2026-02-18T10:30:00Z",
      "environment": "test-cluster-abc123 (OCP 4.15.1)",
      "total_time_minutes": 83,
      "results": {
        "total_tests": 3,
        "passed": 2,
        "failed": 1,
        "needs_review": 0,
        "skipped": 0,
        "blocked": 0,
        "pass_rate_percentage": 66.67
      }
    },
    "test_status_summary": {
      "test_1": {
        "title": "Create ROSA HCP Cluster",
        "status": "passed",
        "time": "42 minutes"
      },
      "test_2": {
        "title": "Configure Monitoring",
        "status": "passed",
        "time": "25 minutes"
      },
      "test_3": {
        "title": "Validate Metrics Collection",
        "status": "failed",
        "time": "16 minutes",
        "failure_reason": "Metrics not appearing in Prometheus"
      }
    },
    "issues_summary": [
      {
        "test_id": "test_3",
        "step": 4,
        "issue": "Prometheus not scraping ServiceMonitor targets",
        "impact": "high",
        "next_steps": "Investigate ServiceMonitor label selector mismatch"
      }
    ]
  }
}
```

### 4. Generate Human-Readable Report (10-15 minutes)

**Create markdown report:**

```markdown
# Test Execution Report

**Test Plan:** RHOBS Test Plan v2
**Execution Date:** 2026-02-18 10:30:00 EST
**Executor:** Claude Sonnet 4.5
**Environment:** test-cluster-abc123 (OpenShift 4.15.1, us-east-1)

## Summary

| Metric | Value |
|--------|-------|
| Total Tests | 3 |
| Passed | 2 ✅ |
| Failed | 1 ❌ |
| Needs Review | 0 ⚠️ |
| Skipped | 0 ⏭️ |
| Blocked | 0 🚫 |
| Pass Rate | 66.67% |
| Total Time | 1h 23min |

## Test Results

### ✅ Test 1: Create ROSA HCP Cluster
**Status:** PASSED
**Time:** 42 minutes
**Completed:** 2026-02-18 11:12:00

All 5 steps completed successfully. All validation criteria met (3/3).

#### Steps:
1. ✅ Verify AWS credentials (2s)
2. ✅ Create ROSA cluster (98s)
3. ✅ Wait for cluster ready (35min)
4. ✅ Configure cluster authentication (45s)
5. ✅ Verify cluster access (3s)

**Notes:** Test completed without issues. Cluster creation time was within expected range.

---

### ✅ Test 2: Configure Monitoring
**Status:** PASSED
**Time:** 25 minutes
**Completed:** 2026-02-18 11:37:00

All steps completed successfully.

---

### ❌ Test 3: Validate Metrics Collection
**Status:** FAILED
**Time:** 16 minutes
**Completed:** 2026-02-18 11:53:00

Test failed at step 4: "Query Prometheus for application metrics"

#### Failed Step Details:
**Step 4:** Query Prometheus for application metrics
**Command:** `oc exec -n openshift-monitoring prometheus-0 -- promtool query instant http://localhost:9090 'up{job="myapp"}'`
**Expected Output:** `up{job="myapp"} => 1`
**Actual Output:** `Error: no data for query`
**Exit Code:** 1

#### Root Cause Analysis:
Prometheus is not scraping the application's metrics endpoint. Investigation shows:
- ServiceMonitor exists and is valid ✓
- Application pods are running ✓
- Metrics endpoint is accessible ✓
- **Issue:** ServiceMonitor label selector doesn't match service labels ❌

#### Fix Applied:
Updated ServiceMonitor selector from `app: my-app` to `app: myapp` to match actual service labels.

#### Recommendation:
Retry test after ServiceMonitor update. Consider adding label validation step to prevent this issue in future runs.

---

## Environment Details

### Cluster Information
- **Name:** test-cluster-abc123
- **OpenShift Version:** 4.15.1
- **ROSA Version:** 1.2.35
- **AWS Region:** us-east-1
- **Cluster Type:** HCP (Hosted Control Plane)

### Tools Used
- `rosa` CLI: 1.2.35
- `oc` CLI: 4.15
- `aws` CLI: 2.15.1

## Issues Encountered

### High Priority
1. **ServiceMonitor Label Mismatch** (Test 3, Step 4)
   - Impact: Metrics not collected
   - Status: Fix applied
   - Next Steps: Re-run test to verify fix

### Medium Priority
None

### Low Priority
None

## Recommendations

1. **Add Pre-Flight Checks:** Create a preliminary test to validate ServiceMonitor selectors match existing services before attempting metrics validation
2. **Update Test Plan:** Add a troubleshooting step for label selector mismatches
3. **Retry Test 3:** After ServiceMonitor fix is applied

## Artifacts

### Logs
- [test1.log](test-results/20260218_103000/test1.log)
- [test2.log](test-results/20260218_103000/test2.log)
- [test3.log](test-results/20260218_103000/test3.log)

### Screenshots
None

### Generated Files
- Updated test plan with results: `examples/rhobs_test_plan_v2_results.json`
- Execution log: `test-results/20260218_103000/execution.log`

## Sign-Off

**Executed By:** Claude Sonnet 4.5
**Reviewed By:** _[Pending Review]_
**Approved By:** _[Pending Review]_

**Execution Complete:** 2026-02-18 11:53:00
```

**Save markdown report:**

```bash
# Generate and save report
cat << 'EOF' > test-results/20260218_103000/REPORT.md
[markdown content above]
EOF
```

### 5. Validate Updated Test Plan (5 minutes)

**Validate JSON syntax:**

```bash
# Check JSON is valid
cat examples/rhobs_test_plan_v2_results.json | jq . > /dev/null
if [ $? -eq 0 ]; then
  echo "✓ Updated test plan JSON is valid"
else
  echo "✗ Invalid JSON in updated test plan"
  exit 1
fi
```

**Verify completeness:**

```bash
# Check all executed tests have results
executed_tests=$(jq -r '.execution_summary.latest_execution.results.total_tests' \
  examples/rhobs_test_plan_v2_results.json)

tests_with_results=$(jq '[.testcases | to_entries[] | select(.value.test_result != null)] | length' \
  examples/rhobs_test_plan_v2_results.json)

if [ "$executed_tests" -eq "$tests_with_results" ]; then
  echo "✓ All executed tests have results attached"
else
  echo "⚠ Some tests missing results ($tests_with_results/$executed_tests)"
fi
```

**Test with viewer:**

```bash
# Verify updated test plan renders in HTML
./build/testplan-viewer -i examples/rhobs_test_plan_v2_results.json \
  -o test-results/20260218_103000/testplan-results.html

if [ $? -eq 0 ]; then
  echo "✓ Updated test plan renders successfully"
  echo "View at: test-results/20260218_103000/testplan-results.html"
else
  echo "✗ Failed to render updated test plan"
fi
```

## Result Status Meanings

### Test-Level Status

**passed**: All steps completed successfully, all validation criteria met
```json
{
  "status": "passed",
  "validation_result": {
    "all_passed": true,
    "criteria_met": 3,
    "criteria_total": 3
  }
}
```

**failed**: One or more critical steps failed or validation criteria not met
```json
{
  "status": "failed",
  "validation_result": {
    "all_passed": false,
    "criteria_met": 1,
    "criteria_total": 3
  },
  "failure_details": {
    "failed_step": 4,
    "failure_reason": "Command exited with code 1"
  }
}
```

**needs_review**: Ambiguous results requiring human judgment
```json
{
  "status": "needs_review",
  "review_reason": "Output format differs from expected but command succeeded",
  "reviewer_action_required": true
}
```

**skipped**: Test not executed due to dependency or blocker
```json
{
  "status": "skipped",
  "skip_reason": "Prerequisite test_1 failed",
  "dependencies_not_met": ["test_1"]
}
```

**blocked**: Environmental issue prevented execution
```json
{
  "status": "blocked",
  "blocker": "No cluster access - authentication failed",
  "blocker_type": "environmental"
}
```

### Step-Level Status

Each step within a test case has its own status:

```json
{
  "step_number": 1,
  "execution_result": {
    "status": "passed",
    "exit_code": 0,
    "actual_output": "...",
    "matched_expected": true
  }
}
```

## Special Cases

### Partial Test Completion

**When a test is partially completed:**

```json
{
  "test_result": {
    "status": "failed",
    "completed_steps": 3,
    "total_steps": 5,
    "aborted_at_step": 4,
    "abort_reason": "Critical failure in step 3 prevented continuation",
    "steps": [
      {"step_number": 1, "execution_result": {"status": "passed"}},
      {"step_number": 2, "execution_result": {"status": "passed"}},
      {"step_number": 3, "execution_result": {"status": "failed"}},
      {"step_number": 4, "execution_result": {"status": "not_executed"}},
      {"step_number": 5, "execution_result": {"status": "not_executed"}}
    ]
  }
}
```

### Multiple Execution History

**Track multiple test runs over time:**

```json
{
  "metadata": {
    "execution_history": [
      {
        "execution_id": "exec-20260215-143000",
        "date": "2026-02-15T14:30:00Z",
        "summary": {"total": 3, "passed": 1, "failed": 2}
      },
      {
        "execution_id": "exec-20260218-103000",
        "date": "2026-02-18T10:30:00Z",
        "summary": {"total": 3, "passed": 2, "failed": 1}
      }
    ]
  },
  "testcases": {
    "test_1": {
      "test_result": {
        // Latest execution result
        "execution_id": "exec-20260218-103000",
        "status": "passed"
      },
      "execution_history": [
        {
          "execution_id": "exec-20260215-143000",
          "date": "2026-02-15T14:30:00Z",
          "status": "failed"
        },
        {
          "execution_id": "exec-20260218-103000",
          "date": "2026-02-18T10:30:00Z",
          "status": "passed"
        }
      ]
    }
  }
}
```

### Artifacts and Attachments

**Reference external artifacts:**

```json
{
  "test_result": {
    "artifacts": {
      "logs": [
        "test-results/20260218_103000/test1.log",
        "test-results/20260218_103000/operator-logs.txt"
      ],
      "screenshots": [
        {
          "step": 5,
          "description": "Grafana dashboard showing metrics",
          "file": "test-results/20260218_103000/grafana-dashboard.png",
          "timestamp": "2026-02-18T11:12:00Z"
        }
      ],
      "files": [
        {
          "name": "servicemonitor.yaml",
          "purpose": "ServiceMonitor configuration used in test",
          "path": "test-results/20260218_103000/servicemonitor.yaml"
        }
      ]
    }
  }
}
```

## Quality Checklist

Before finalizing updated test plan:

- [ ] All test results merged correctly
- [ ] Execution metadata added to metadata section
- [ ] Step-level results attached to each step
- [ ] Validation results included
- [ ] Test-level status determined correctly
- [ ] Execution summary generated
- [ ] Issues documented with next steps
- [ ] Artifacts referenced with valid paths
- [ ] JSON is syntactically valid
- [ ] HTML renders successfully
- [ ] Markdown report generated
- [ ] All timestamps in ISO 8601 format

## Output Files

After updating test plan with results, you should have:

```
examples/
  rhobs_test_plan_v2.json           # Original test plan
  rhobs_test_plan_v2_results.json   # Updated with results

test-results/
  20260218_103000/
    execution.log                    # Full execution log
    results.json                     # Raw execution results
    REPORT.md                        # Human-readable report
    testplan-results.html            # HTML with results
    test1.log                        # Individual test logs
    test2.log
    test3.log
    screenshots/                     # Visual artifacts
    configs/                         # Config files used
```

## Example Usage

**User provides:**
```
Original test plan: examples/rhobs_test_plan_v2.json
Execution results: test-results/20260218_103000/results.json
```

**You should:**

1. Load both JSON files
2. Validate structure compatibility
3. Merge execution results into test plan
4. Add execution summary
5. Generate markdown report
6. Save updated files:
   - `examples/rhobs_test_plan_v2_results.json`
   - `test-results/20260218_103000/REPORT.md`
7. Render HTML with results
8. Print summary to console

## Best Practices

**Do:**
- Preserve original test plan structure
- Add results as new fields, don't overwrite original data
- Track execution history for trending
- Include enough context in reports for future reference
- Reference artifacts with relative paths
- Use ISO 8601 for all timestamps
- Calculate pass rates and time metrics
- Document issues with clear next steps

**Don't:**
- Overwrite original test plan file
- Lose execution data during merge
- Skip validation after merge
- Forget to update execution_history
- Use absolute paths for artifacts
- Omit environmental context
- Ignore partial failures in summary
- Forget to generate human-readable report
