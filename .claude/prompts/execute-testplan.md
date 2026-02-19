# Execute Test Plan

## Objective

Execute a test plan by running through each test case step-by-step, validating results, capturing outputs, and determining pass/fail status. This prompt enables Claude to act as a test executor, running manual tests and determining outcomes.

## Your Role

You are a test execution engineer performing manual validation. Your goal is to:

1. **Prepare environment** - Verify prerequisites are met
2. **Execute test steps** - Run each command systematically
3. **Capture outputs** - Record actual results from each step
4. **Validate results** - Compare actual vs expected outcomes
5. **Determine status** - Mark each test as pass/fail/needs_review/skipped
6. **Document issues** - Record errors, blockers, and deviations

## Inputs You'll Receive

The user will provide:

- **Test plan JSON file**: Path to the test plan to execute
- **Environment context**: Cluster details, credentials, environment variables
- **Execution scope**:
  - All tests (full run)
  - Specific tests (comma-separated list)
  - Resume from step (continue failed execution)
- **Execution mode**:
  - Automated (Claude runs all commands)
  - Interactive (Claude prompts for each step)
  - Dry-run (Claude shows what would be executed)

## Pre-Execution Checks

### 1. Validate Test Plan (2-5 minutes)

**Load and validate JSON:**

```bash
# Check JSON syntax
cat examples/testplan.json | jq . > /dev/null && echo "Valid JSON" || echo "Invalid JSON"

# Validate with viewer
./build/testplan-viewer -i examples/testplan.json -o /tmp/test.html

# Extract test case list
jq -r '.testcases | keys[]' examples/testplan.json
```

**Verify test plan structure:**
- Metadata is complete
- All test cases have required fields
- Learning path sequence is valid
- Dependencies are correctly specified

### 2. Verify Prerequisites (5-10 minutes)

**Check required knowledge:**
- Review `prerequisites.required_knowledge`
- Confirm user has necessary background
- Suggest training if gaps exist

**Verify required access:**

```bash
# For each item in prerequisites.required_access
jq -r '.prerequisites.required_access[] | .verification' testplan.json

# Example checks:
oc whoami  # OpenShift authentication
aws sts get-caller-identity  # AWS access
kubectl cluster-info  # Kubernetes access
```

**Setup environment:**

```bash
# Set environment variables from test plan
jq -r '.prerequisites.environment_setup.variables | to_entries[] | "export \(.key)=\(.value)"' testplan.json

# Install required tools
jq -r '.prerequisites.environment_setup.tools[] | .install' testplan.json

# Verify tools are installed
jq -r '.prerequisites.environment_setup.tools[] | .verify' testplan.json
```

**Check dependencies:**

```bash
# Verify test dependencies are satisfied
# If test_2 depends on test_1, ensure test_1 completed successfully
```

### 3. Prepare Execution Environment (5 minutes)

**Create execution log directory:**

```bash
mkdir -p test-results/$(date +%Y%m%d_%H%M%S)
cd test-results/$(date +%Y%m%d_%H%M%S)
```

**Set up logging:**

```bash
# Log all commands and outputs
exec > >(tee execution.log)
exec 2>&1
```

**Initialize results structure:**

```json
{
  "execution_metadata": {
    "test_plan": "examples/rhobs_test_plan_v2.json",
    "execution_date": "2026-02-18T10:30:00Z",
    "executor": "claude",
    "environment": {
      "cluster": "test-cluster-abc123",
      "ocp_version": "4.15.1",
      "region": "us-east-1"
    }
  },
  "test_results": {}
}
```

## Execution Process

### For Each Test Case

#### Step 1: Pre-Test Validation

**Check test prerequisites:**

```bash
# Review test-specific prerequisites
jq -r '.testcases.test_1.test_execution.prerequisites[]' testplan.json

# Verify dependencies completed
# If test depends on test_A, check test_A status = "passed"
```

**Estimate time:**

```bash
# Show user the estimated time
jq -r '.testcases.test_1.metadata.estimated_time' testplan.json
echo "This test will take approximately [time]. Ready to proceed? (y/n)"
```

#### Step 2: Execute Test Steps

**For each step in test_execution.steps:**

```json
{
  "step_number": 1,
  "title": "Verify AWS credentials",
  "command": "aws sts get-caller-identity",
  "expected_output": "{\n  \"UserId\": \"AIDACKCEVSQ6C2EXAMPLE\",\n  \"Account\": \"123456789012\"\n}"
}
```

**Execution pattern:**

```bash
# 1. Display step information
echo "═══════════════════════════════════════════════════════════════"
echo "Step 1: Verify AWS credentials"
echo "───────────────────────────────────────────────────────────────"
echo "Learning note: Understanding AWS authentication is crucial for ROSA"
echo "Why this step: ROSA needs AWS credentials to provision resources"
echo ""

# 2. Show command to be executed
echo "Command:"
echo "  aws sts get-caller-identity"
echo ""

# 3. Execute command and capture output
echo "Executing..."
actual_output=$(aws sts get-caller-identity 2>&1)
exit_code=$?

# 4. Display actual output
echo "Actual output:"
echo "$actual_output"
echo ""

# 5. Compare with expected output
expected_output='{"UserId": "AIDACKCEVSQ6C2EXAMPLE"}'
if echo "$actual_output" | jq -e . > /dev/null 2>&1; then
  # Output is valid JSON, can do structured comparison
  if echo "$actual_output" | jq -e '.Account' > /dev/null 2>&1; then
    echo "✓ Output contains expected fields"
    step_status="passed"
  else
    echo "✗ Output missing expected fields"
    step_status="failed"
  fi
else
  # Not JSON or command failed
  if [ $exit_code -eq 0 ]; then
    echo "⚠ Command succeeded but output format unexpected"
    step_status="needs_review"
  else
    echo "✗ Command failed with exit code: $exit_code"
    step_status="failed"
  fi
fi

# 6. Check for common errors
if echo "$actual_output" | grep -q "Unable to locate credentials"; then
  echo ""
  echo "⚠ Common error detected: Unable to locate credentials"
  echo "Solution: Run 'aws configure' to set up credentials"
  echo "Learn more: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html"
fi

# 7. Record step result
jq --arg step "1" --arg status "$step_status" --arg output "$actual_output" \
  '.test_results.test_1.steps[$step] = {status: $status, output: $output}' \
  results.json > tmp.json && mv tmp.json results.json

# 8. Decide whether to continue
if [ "$step_status" = "failed" ]; then
  echo ""
  echo "Step failed. Options:"
  echo "  1) Skip this step and continue"
  echo "  2) Retry this step"
  echo "  3) Abort test"
  read -p "Choice (1/2/3): " choice

  case $choice in
    1) echo "Skipping step..." ;;
    2) echo "Retrying step..."; # Re-run step logic ;;
    3) echo "Aborting test"; test_status="aborted"; break ;;
  esac
fi
```

#### Step 3: Validate Test Results

**After all steps complete, validate success criteria:**

```bash
# Get validation criteria from test plan
jq -r '.testcases.test_1.test_execution.validation.success_criteria[]' testplan.json

# For each criterion, perform check
# Example criteria: "Cluster appears in 'rosa list clusters'"
rosa list clusters | grep "$CLUSTER_NAME"
if [ $? -eq 0 ]; then
  echo "✓ Criterion met: Cluster appears in rosa list clusters"
  criteria_met=$((criteria_met + 1))
else
  echo "✗ Criterion not met: Cluster not in rosa list clusters"
fi

# Calculate overall test status
total_criteria=$(jq '.testcases.test_1.test_execution.validation.success_criteria | length' testplan.json)
if [ $criteria_met -eq $total_criteria ]; then
  test_status="passed"
elif [ $criteria_met -eq 0 ]; then
  test_status="failed"
else
  test_status="needs_review"  # Partial success
fi
```

#### Step 4: Record Test Results

**Create test result entry:**

```json
{
  "test_results": {
    "test_1": {
      "test_id": "test_1",
      "test_title": "Create ROSA HCP Cluster",
      "status": "passed",
      "execution_time": "42 minutes",
      "started_at": "2026-02-18T10:30:00Z",
      "completed_at": "2026-02-18T11:12:00Z",
      "steps": [
        {
          "step_number": 1,
          "title": "Verify AWS credentials",
          "status": "passed",
          "actual_output": "{\n  \"UserId\": \"AIDAI...\",\n  \"Account\": \"123456789012\"\n}",
          "execution_time": "2 seconds"
        }
      ],
      "validation": {
        "criteria_met": 3,
        "criteria_total": 3,
        "details": [
          {
            "criterion": "Cluster appears in 'rosa list clusters'",
            "met": true
          }
        ]
      },
      "notes": "Test completed successfully with no issues",
      "issues_encountered": [],
      "screenshots": []
    }
  }
}
```

## Execution Modes

### 1. Automated Execution

**Full automation mode - Claude runs all commands:**

```bash
# User specifies automated mode
execute_testplan --mode automated --tests all

# Claude executes:
# - Verifies prerequisites automatically
# - Runs all test steps in sequence
# - Captures all outputs
# - Validates results programmatically
# - Records results to JSON
# - Generates summary report
```

**Use for:**
- Regression testing
- Nightly validation
- CI/CD integration
- Tests with deterministic outcomes

### 2. Interactive Execution

**Claude prompts user before each action:**

```bash
# User specifies interactive mode
execute_testplan --mode interactive --tests test_1,test_2

# Claude behavior:
# 1. Shows test overview
# 2. Asks: "Ready to start test_1?"
# 3. Shows each step
# 4. Asks: "Execute this step?"
# 5. Runs command when approved
# 6. Shows output
# 7. Asks: "Does this match expected output?"
# 8. Records user's determination
```

**Use for:**
- Learning/training scenarios
- Tests requiring manual verification
- Visual validation steps
- Tests with human decision points

### 3. Dry-Run Mode

**Show what would be executed without running:**

```bash
# User specifies dry-run mode
execute_testplan --mode dry-run --tests all

# Claude behavior:
# - Loads test plan
# - Shows each test and step
# - Displays commands that would run
# - Shows expected outputs
# - Identifies potential issues
# - Estimates total execution time
# - No actual execution
```

**Use for:**
- Previewing test plan before running
- Verifying test plan correctness
- Time estimation
- Risk assessment

## Status Determination

### Test Step Status

**passed**: Command executed successfully, output matches expected
```bash
# Exit code 0
# Output contains expected patterns/values
# No errors in stderr
```

**failed**: Command failed or output incorrect
```bash
# Exit code non-zero
# Output missing expected values
# Error messages in stderr
# Validation criteria not met
```

**needs_review**: Ambiguous result requiring human judgment
```bash
# Command succeeded but output format differs
# Partial match of expected output
# Non-critical warnings present
# Success criteria partially met
```

**skipped**: Step not executed due to blocker
```bash
# Previous step failed critically
# Prerequisite not met
# Environmental blocker (no cluster access)
# User chose to skip
```

### Test Case Status

**Overall test status based on step statuses:**

```
All steps passed → test status: PASSED
Any step failed → test status: FAILED
Any step needs_review → test status: NEEDS_REVIEW
All steps skipped → test status: SKIPPED
Execution aborted → test status: ABORTED
Blocker prevented execution → test status: BLOCKED
```

## Error Handling

### Command Execution Errors

**Capture detailed error information:**

```bash
# Run command with error capture
output=$(command 2>&1)
exit_code=$?

if [ $exit_code -ne 0 ]; then
  error_info=$(cat <<EOF
{
  "command": "command here",
  "exit_code": $exit_code,
  "stdout": "...",
  "stderr": "...",
  "timestamp": "$(date -Iseconds)"
}
EOF
)

  # Check for known errors in test plan
  jq -r '.testcases.test_1.test_execution.steps[0].common_errors[] |
    select(.error | contains("pattern")) |
    "Known issue: \(.error)\nSolution: \(.solution)"' testplan.json

  # Consult troubleshooting section
  jq -r '.testcases.test_1.troubleshooting.common_failures[] |
    select(.symptom | contains("error pattern")) |
    "Symptom: \(.symptom)\nCause: \(.cause)\nFix: \(.fix)"' testplan.json
fi
```

### Environmental Blockers

**Detect and report blockers:**

```bash
# No cluster access
if ! oc whoami &> /dev/null; then
  echo "BLOCKER: No cluster access"
  echo "Cannot proceed with tests requiring oc commands"
  test_status="blocked"
  blocker="No authenticated cluster connection"
fi

# Missing credentials
if ! aws sts get-caller-identity &> /dev/null; then
  echo "BLOCKER: No AWS credentials"
  test_status="blocked"
  blocker="AWS credentials not configured"
fi

# Resource already exists
if oc get cluster "$CLUSTER_NAME" &> /dev/null; then
  echo "BLOCKER: Cluster $CLUSTER_NAME already exists"
  echo "Options: 1) Use different name, 2) Delete existing cluster"
  test_status="blocked"
  blocker="Resource name conflict"
fi
```

## Output Formats

### Execution Log (execution.log)

```
═══════════════════════════════════════════════════════════════
Test Plan Execution: RHOBS Test Plan v2
Started: 2026-02-18 10:30:00
Environment: test-cluster-abc123 (OCP 4.15.1)
═══════════════════════════════════════════════════════════════

[10:30:01] Loading test plan: examples/rhobs_test_plan_v2.json
[10:30:02] ✓ Test plan valid
[10:30:03] ✓ Prerequisites verified
[10:30:04] Starting execution of 3 test cases

───────────────────────────────────────────────────────────────
Test 1 of 3: Create ROSA HCP Cluster
Estimated time: 45 minutes
Difficulty: beginner
───────────────────────────────────────────────────────────────

[10:30:05] Step 1/5: Verify AWS credentials
[10:30:05]   Command: aws sts get-caller-identity
[10:30:06]   ✓ Passed
[10:30:06]   Output: {"UserId": "AIDAI...", "Account": "123456789012"}

[10:30:07] Step 2/5: Create ROSA cluster
[10:30:07]   Command: rosa create cluster --cluster-name test-cluster
[10:31:45]   ✓ Passed
[10:31:45]   Output: Cluster 'test-cluster' created successfully

...

[11:12:00] ✓ Test 1 PASSED (42 minutes)
[11:12:00]   3/3 validation criteria met

═══════════════════════════════════════════════════════════════
Execution Summary
═══════════════════════════════════════════════════════════════
Total tests: 3
Passed: 2
Failed: 1
Needs Review: 0
Skipped: 0
Blocked: 0

Total time: 1 hour 23 minutes
═══════════════════════════════════════════════════════════════
```

### Results JSON (results.json)

**Complete results structure for update-testplan-results.md:**

```json
{
  "execution_metadata": {
    "test_plan_file": "examples/rhobs_test_plan_v2.json",
    "test_plan_version": "1.0",
    "execution_id": "exec-20260218-103000",
    "execution_date": "2026-02-18T10:30:00Z",
    "completed_date": "2026-02-18T11:53:00Z",
    "executor": "claude-sonnet-4.5",
    "execution_mode": "automated",
    "environment": {
      "cluster_name": "test-cluster-abc123",
      "cluster_id": "abc123...",
      "ocp_version": "4.15.1",
      "rosa_version": "1.2.35",
      "aws_region": "us-east-1",
      "tools": {
        "rosa": "1.2.35",
        "oc": "4.15",
        "aws": "2.15.1"
      }
    }
  },
  "test_results": {
    "test_1": {
      "test_id": "test_1",
      "test_title": "Create ROSA HCP Cluster",
      "status": "passed",
      "execution_time_minutes": 42,
      "started_at": "2026-02-18T10:30:00Z",
      "completed_at": "2026-02-18T11:12:00Z",
      "steps": [
        {
          "step_number": 1,
          "title": "Verify AWS credentials",
          "status": "passed",
          "command_executed": "aws sts get-caller-identity",
          "actual_output": "...",
          "exit_code": 0,
          "execution_time_seconds": 2,
          "timestamp": "2026-02-18T10:30:05Z"
        }
      ],
      "validation": {
        "criteria_met": 3,
        "criteria_total": 3,
        "details": [
          {
            "criterion": "Cluster appears in 'rosa list clusters'",
            "met": true,
            "verification_command": "rosa list clusters | grep test-cluster",
            "verification_output": "test-cluster ready ..."
          }
        ]
      },
      "issues_encountered": [],
      "notes": "Test completed successfully",
      "artifacts": {
        "logs": ["step1-output.log"],
        "screenshots": [],
        "files": []
      }
    }
  },
  "summary": {
    "total_tests": 3,
    "passed": 2,
    "failed": 1,
    "needs_review": 0,
    "skipped": 0,
    "blocked": 0,
    "total_execution_time_minutes": 83,
    "pass_rate": 66.67
  }
}
```

## Special Cases

### Visual Validation Steps

**For steps requiring visual confirmation:**

```bash
# Example: "Verify dashboard shows metrics"
echo "Step 5: Verify Grafana dashboard shows metrics"
echo ""
echo "This step requires visual confirmation in the Grafana UI:"
echo "1. Open Grafana: https://grafana-route-openshift-monitoring.apps.cluster.example.com"
echo "2. Navigate to 'RHOBS Monitoring' dashboard"
echo "3. Verify panels show data (not 'No data')"
echo "4. Check timestamp is recent (last 5 minutes)"
echo ""
read -p "Do metrics appear correctly in the dashboard? (y/n): " response

if [ "$response" = "y" ]; then
  step_status="passed"
  echo "✓ Visual validation passed"
else
  step_status="needs_review"
  read -p "Describe what you see: " observation
  echo "Observation recorded: $observation"
fi
```

### Long-Running Steps

**For steps that take >5 minutes:**

```bash
# Example: Cluster creation (30-45 minutes)
echo "Step 2: Create ROSA cluster (estimated: 45 minutes)"
echo "Command: rosa create cluster --cluster-name $CLUSTER_NAME --wait"
echo ""
echo "This is a long-running operation. Starting..."

# Run in background with progress monitoring
rosa create cluster --cluster-name "$CLUSTER_NAME" > step2.log 2>&1 &
pid=$!

# Monitor progress
while kill -0 $pid 2>/dev/null; do
  # Check rosa describe for status
  status=$(rosa describe cluster -c "$CLUSTER_NAME" --output json 2>/dev/null | jq -r '.status.state')
  echo "[$(date +%H:%M:%S)] Cluster status: $status"
  sleep 60
done

# Get result
wait $pid
exit_code=$?
output=$(cat step2.log)

echo "Cluster creation completed with exit code: $exit_code"
```

### Retry Logic

**For flaky operations:**

```bash
# Retry pattern for network-dependent commands
max_attempts=3
attempt=1
success=false

while [ $attempt -le $max_attempts ] && [ "$success" = "false" ]; do
  echo "Attempt $attempt/$max_attempts: $command"

  output=$(eval "$command" 2>&1)
  exit_code=$?

  if [ $exit_code -eq 0 ]; then
    success=true
    echo "✓ Command succeeded on attempt $attempt"
  else
    echo "✗ Attempt $attempt failed"
    if [ $attempt -lt $max_attempts ]; then
      echo "Retrying in 10 seconds..."
      sleep 10
    fi
  fi

  attempt=$((attempt + 1))
done

if [ "$success" = "false" ]; then
  echo "✗ Command failed after $max_attempts attempts"
  step_status="failed"
fi
```

## Best Practices

**Do:**
- Verify prerequisites before starting
- Capture all command outputs
- Log timestamps for each step
- Compare actual vs expected outputs
- Consult troubleshooting guides on errors
- Record environmental context
- Take screenshots for visual steps
- Ask for help when status is ambiguous

**Don't:**
- Skip prerequisite checks
- Continue after critical failures
- Assume partial output means success
- Ignore warnings and non-critical errors
- Execute destructive commands without confirmation
- Forget to record actual outputs
- Mark tests as passed without validation
