# RHOBS-Next Test Plan Execution Review

**Date**: 2026-02-18
**Test Plan**: rhobs-next-testplan.json
**Reviewer**: Claude (testplan-executor skill)
**Review Type**: Pre-execution analysis (tests NOT executed)

## Executive Summary

Reviewed 5 test cases for execution safety, dependency ordering, and completeness. Identified **3 high-priority improvements** and **5 medium-priority enhancements** to the test plan execution approach before actual test runs.

**Safety Status**: ✅ Safe to execute with proper cleanup tracking
**Dependency Order**: ✅ Valid - all dependencies exist
**Cleanup Coverage**: ⚠️ 50% - needs tracking system
**Production Safety**: ✅ 2 tests safe for production, 3 require test environments

---

## Test Inventory

| # | Test ID | Category | Safety | Duration | Dependencies |
|---|---------|----------|--------|----------|--------------|
| 1 | test_verify_route_monitor_deployment | read-only | ✅ Prod-safe | 15 min | None |
| 2 | test_create_routemonitor_cr | modifies-state | ⚠️ Test only | 25 min | Test 1 |
| 3 | test_verify_probe_metrics | read-only | ✅ Prod-safe | 20 min | Test 2 |
| 4 | test_create_probe_via_api | modifies-state | ⚠️ Test only | 30 min | Test 1 |
| 5 | test_cleanup_routemonitor_cr | cleanup | ⚠️ Test only | 10 min | Test 2 |

**Total Estimated Time**: 100 minutes (1 hour 40 minutes)

---

## Execution Order Analysis

### Recommended Sequence

```
Phase 1: Read-Only Validation (No State Changes)
├─ 1. test_verify_route_monitor_deployment (15 min)
│     └─> Validates RMO operator is healthy
│
Phase 2: Modifying Tests (State Changes - Requires Cleanup)
├─ 2. test_create_routemonitor_cr (25 min)
│     ├─> Depends on: Test 1 (RMO must be running)
│     ├─> Creates: RouteMonitor CR, ServiceMonitor
│     └─> Cleanup: Test 5 (test_cleanup_routemonitor_cr)
│
├─ 3. test_verify_probe_metrics (20 min)
│     ├─> Depends on: Test 2 (RouteMonitor must exist)
│     └─> Validates metrics from created probe
│
├─ 4. test_create_probe_via_api (30 min)
│     ├─> Depends on: Test 1 (RMO running)
│     ├─> Creates: Probe via API, waits for agent reconciliation
│     └─> Cleanup: MISSING - needs cleanup test
│
Phase 3: Cleanup (State Restoration)
└─ 5. test_cleanup_routemonitor_cr (10 min)
      ├─> Depends on: Test 2 (RouteMonitor must exist)
      └─> Deletes: RouteMonitor CR, cascading deletion of ServiceMonitor
```

### Dependency Validation

✅ **All dependencies exist**
✅ **No circular dependencies**
⚠️ **Incomplete cleanup chain** - Test 4 has no cleanup test

### Alternative Execution Paths

**Quick Validation Path** (35 minutes):
```
1. test_verify_route_monitor_deployment
2. test_verify_probe_metrics (requires Test 2 first, so this path is actually invalid)
```
⚠️ **Issue**: Alternative path "quick_validation" references Test 3 which depends on Test 2, but Test 2 isn't in the path. Need to fix.

**Modifies-State Only Path** (65 minutes):
```
2. test_create_routemonitor_cr
4. test_create_probe_via_api
5. test_cleanup_routemonitor_cr
```
⚠️ **Issue**: Test 2 and 4 both depend on Test 1, but Test 1 isn't in this path. Need to add Test 1 as prerequisite.

---

## Safety Analysis

### Read-Only Tests (Production-Safe)

**Test 1: test_verify_route_monitor_deployment**
- **Impact**: None - only queries existing resources
- **Reversible**: N/A - no changes made
- **Risk Level**: Low
- **Production Safe**: ✅ Yes
- **Requires Confirmation**: ❌ No
- **Idempotent**: ✅ Yes
- **Safe to Retry**: ✅ Yes

**Test 3: test_verify_probe_metrics**
- **Impact**: None - only queries Prometheus
- **Reversible**: N/A - no changes made
- **Risk Level**: Low
- **Production Safe**: ✅ Yes (but depends on Test 2 which isn't prod-safe)
- **Requires Confirmation**: ❌ No
- **Idempotent**: ✅ Yes
- **Safe to Retry**: ✅ Yes

### Modifies-State Tests (Test Environment Only)

**Test 2: test_create_routemonitor_cr**
- **Impact**: Creates RouteMonitor CR → triggers ServiceMonitor creation
- **Affected Resources**:
  - RouteMonitor: `console-route-monitor` in `default` namespace
  - ServiceMonitor: Auto-created by RMO
- **Reversible**: ✅ Yes - via Test 5 cleanup
- **Risk Level**: Medium
- **Production Safe**: ❌ No
- **Requires Confirmation**: ✅ Yes
- **Warning Message**: "⚠️ This test creates a RouteMonitor CR which triggers ServiceMonitor creation. Run in test clusters only. Cleanup test (test_cleanup_routemonitor_cr) must be run afterward."
- **Idempotent**: ✅ Yes (oc apply)
- **Safe to Retry**: ✅ Yes
- **Cleanup Test**: test_cleanup_routemonitor_cr

**Test 4: test_create_probe_via_api**
- **Impact**: Creates probe in RHOBS Synthetics API → triggers agent reconciliation → creates Probe CR
- **Affected Resources**:
  - RHOBS API probe entry
  - Probe CR in cluster (created by agent)
  - Potentially blackbox-exporter deployment (if first probe)
- **Reversible**: ⚠️ Partial - no cleanup test exists
- **Risk Level**: Medium
- **Production Safe**: ❌ No
- **Requires Confirmation**: ✅ Yes
- **Warning Message**: "⚠️ This test creates a probe via RHOBS API and waits for agent reconciliation. Run in test environments only. **CLEANUP TEST MISSING - manual cleanup required.**"
- **Idempotent**: ⚠️ No - creates new probe each time
- **Safe to Retry**: ⚠️ No - will create duplicate probes
- **Cleanup Test**: ❌ MISSING - **HIGH PRIORITY TO ADD**

**Test 5: test_cleanup_routemonitor_cr**
- **Impact**: Deletes RouteMonitor CR → cascading deletion of ServiceMonitor
- **Affected Resources**:
  - RouteMonitor: `console-route-monitor` (deleted)
  - ServiceMonitor: Auto-deleted by Kubernetes ownerReferences
- **Reversible**: ⚠️ No - deletion is permanent (but can recreate via Test 2)
- **Risk Level**: Low (cleanup operation)
- **Production Safe**: ❌ No (should only cleanup test resources)
- **Requires Confirmation**: ✅ Yes
- **Restores State**: ✅ Yes - returns cluster to pre-Test-2 state
- **Safe to Retry**: ✅ Yes (idempotent delete)

---

## Pre-Execution Checklist

Before running ANY test, verify:

### Environment Validation
- [ ] Cluster is NOT production (check cluster ID, name, sector)
- [ ] RMO operator is installed (`oc get deployment route-monitor-operator -n openshift-route-monitor-operator`)
- [ ] RHOBS Synthetics Agent is running (for Test 4)
- [ ] RHOBS Synthetics API is accessible (`curl $RHOBS_API_URL/healthz`)
- [ ] Prometheus is deployed and accessible
- [ ] Have cluster-admin permissions or equivalent RBAC

### Access & Credentials
- [ ] KUBECONFIG points to correct cluster
- [ ] RHOBS_API_URL environment variable set
- [ ] ACCESS_TOKEN (OIDC token) is valid
- [ ] Can create resources in target namespaces

### Backup & Safety
- [ ] Created backup of existing RouteMonitors: `oc get routemonitors -A -o yaml > backup/routemonitors-$(date +%Y%m%d-%H%M%S).yaml`
- [ ] Documented current state for restore if needed
- [ ] Have rollback plan ready

---

## Execution Flow Walkthrough

### Test 1: test_verify_route_monitor_deployment

**Pre-Flight**:
- ✅ No dependencies
- ✅ No state changes
- ✅ No cleanup required
- ✅ Safe to run in production

**Execution**:
1. Verify deployment exists and is ready
2. Check pod status
3. Verify RBAC permissions
4. Check operator logs for errors

**Expected Outcomes**:
- Deployment shows 1/1 ready replicas
- Pod is Running
- No error logs in recent 50 lines
- RoleBindings exist

**On Success**:
- Proceed to dependent tests (Test 2, Test 4)

**On Failure**:
- **DO NOT PROCEED** - RMO must be healthy for other tests
- Investigate: deployment not found, pod crash loop, RBAC issues
- Fix RMO before continuing

---

### Test 2: test_create_routemonitor_cr

**Pre-Flight**:
- ✅ Depends on: Test 1 (must pass first)
- ⚠️ Modifies state - creates resources
- ⚠️ Requires cleanup - Test 5 must run after
- ❌ NOT safe for production

**Safety Confirmation Required**:
```
⚠️ SAFETY CHECK - test_create_routemonitor_cr

System Impact: modifies-state
Risk Level: medium
Affected Resources:
  - RouteMonitor: console-route-monitor (default namespace)
  - ServiceMonitor: Auto-created by RMO

Current Cluster: [cluster-id]
Is Production: [yes/no]

This test creates a RouteMonitor CR which triggers ServiceMonitor creation.
Cleanup test (test_cleanup_routemonitor_cr) MUST be run afterward.

Proceed? (yes/no): _
```

**Execution**:
1. Create routemonitor.yaml from template
2. Apply RouteMonitor CR
3. Wait for RMO to reconcile (30-60 seconds)
4. Verify ServiceMonitor was created
5. Check probe begins executing

**Expected Outcomes**:
- RouteMonitor CR created successfully
- ServiceMonitor appears within 60 seconds
- Prometheus scrapes blackbox-exporter for this probe
- probe_success metric appears

**On Success**:
- ✅ Mark for cleanup (add to cleanup tracking list)
- Proceed to Test 3 (validates metrics from this probe)

**On Failure**:
- Check RMO logs: `oc logs -n openshift-route-monitor-operator deployment/route-monitor-operator`
- Verify RMO has RBAC to create ServiceMonitors
- Check if ServiceMonitor CRD exists
- **Still requires cleanup** - run Test 5 to remove created RouteMonitor

---

### Test 3: test_verify_probe_metrics

**Pre-Flight**:
- ✅ Depends on: Test 2 (must pass first)
- ✅ No state changes
- ✅ No cleanup required
- ⚠️ Can run in production BUT depends on Test 2 which can't

**Execution**:
1. Query Prometheus for probe_success metric
2. Validate metric has correct labels
3. Check probe_duration_seconds
4. Verify probe_http_status_code

**Expected Outcomes**:
- probe_success{job="console-route-monitor"} = 1
- probe_http_status_code = 200
- probe_duration_seconds > 0

**On Success**:
- Metrics pipeline validated end-to-end
- Proceed to additional tests

**On Failure**:
- Check Prometheus scrape config
- Verify ServiceMonitor selector matches service labels
- Check blackbox-exporter logs
- Verify target URL is accessible

---

### Test 4: test_create_probe_via_api

**Pre-Flight**:
- ✅ Depends on: Test 1 (RMO must be running)
- ⚠️ Modifies state - creates probe in API and cluster
- ❌ NO CLEANUP TEST - **CRITICAL GAP**
- ❌ NOT safe for production

**Safety Confirmation Required**:
```
⚠️ SAFETY CHECK - test_create_probe_via_api

System Impact: modifies-state
Risk Level: medium
Affected Resources:
  - RHOBS API: Probe entry in /api/v1/probes
  - Kubernetes: Probe CR (created by agent)
  - Potentially: blackbox-exporter deployment

⚠️ WARNING: NO CLEANUP TEST EXISTS
You will need to manually delete the probe after this test:
  1. DELETE via API: curl -X DELETE "$RHOBS_API_URL/api/v1/probes/{id}"
  2. Verify Probe CR deleted by agent

Current Cluster: [cluster-id]
Is Production: [yes/no]

Proceed? (yes/no): _
```

**Execution**:
1. Authenticate with RHOBS API (OIDC token)
2. Create probe definition via POST
3. Wait for agent to reconcile (polling interval: 30s)
4. Verify Probe CR appears in cluster
5. Check probe status transitions: pending → active → running

**Expected Outcomes**:
- API returns 201 Created with probe ID
- Agent reconciles probe within 60 seconds
- Probe CR created in cluster
- Probe status updated to "running" in API

**On Success**:
- ✅ Mark for **MANUAL CLEANUP** (no automated cleanup exists)
- Document probe ID for deletion
- Proceed with caution - leftover probes will consume resources

**On Failure**:
- Check API authentication (token valid?)
- Verify agent is polling the API endpoint
- Check agent logs for errors
- Ensure agent label selector matches probe labels

**Post-Execution Cleanup** (MANUAL):
```bash
# 1. Get probe ID from test output
PROBE_ID="<from-test-output>"

# 2. Delete via API
curl -X DELETE "$RHOBS_API_URL/api/v1/probes/$PROBE_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN"

# 3. Verify agent deletes Probe CR (wait 30-60s)
oc get probes --watch

# 4. Confirm deletion
oc get probes | grep -q "$PROBE_ID" && echo "Still exists!" || echo "Cleaned up"
```

---

### Test 5: test_cleanup_routemonitor_cr

**Pre-Flight**:
- ✅ Depends on: Test 2 (RouteMonitor must exist to clean up)
- ⚠️ Modifies state - deletes resources
- ✅ Restores state to pre-Test-2 condition
- ❌ NOT safe for production (cleanup test-created resources only)

**Safety Confirmation Required**:
```
⚠️ SAFETY CHECK - test_cleanup_routemonitor_cr

System Impact: modifies-state (cleanup operation)
Risk Level: low
Purpose: Remove RouteMonitor CR created by test_create_routemonitor_cr

Deletes:
  - RouteMonitor: console-route-monitor (default namespace)
  - ServiceMonitor: Auto-deleted via ownerReferences

This RESTORES the cluster to its pre-test state.

Current Cluster: [cluster-id]

Proceed with cleanup? (yes/no): _
```

**Execution**:
1. Delete RouteMonitor CR
2. Verify cascading deletion of ServiceMonitor
3. Wait for resources to fully terminate
4. Confirm no orphaned resources remain

**Expected Outcomes**:
- RouteMonitor CR deleted
- ServiceMonitor deleted automatically
- Prometheus stops scraping this probe
- probe_success metric stops updating

**On Success**:
- ✅ Cluster restored to original state
- Test 2 resources fully cleaned up
- Safe to re-run Test 2 if needed

**On Failure**:
- RouteMonitor may have finalizers - check with `oc get routemonitor -o yaml`
- Manually remove finalizers if stuck: `oc patch routemonitor console-route-monitor --type=merge -p '{"metadata":{"finalizers":[]}}'`
- Force delete if necessary: `oc delete routemonitor console-route-monitor --grace-period=0 --force`

---

## Identified Gaps & Improvements

### 🔴 Critical (Block Execution)

1. **Missing Cleanup Test for Test 4**
   - **Issue**: test_create_probe_via_api has no automated cleanup
   - **Impact**: Leftover probes consume API quota, cluster resources, agent reconciliation cycles
   - **Risk**: Resource leak, failed subsequent test runs due to duplicate probes
   - **Recommendation**: Create test_cleanup_probe_via_api
   - **Implementation**:
     ```json
     {
       "test_cleanup_probe_via_api": {
         "metadata": {
           "category": "cleanup",
           "system_impact": {"type": "modifies-state"},
           "state_management": {"restores_state": true}
         },
         "test_execution": {
           "steps": [
             "DELETE probe via API using probe ID from Test 4",
             "Wait for agent to delete Probe CR (30-60s)",
             "Verify Probe CR no longer exists"
           ]
         }
       }
     }
     ```

### 🟡 High Priority (Should Fix Before Production Use)

2. **Alternative Paths Have Invalid Dependencies**
   - **Issue**: "quick_validation" path includes Test 3 which depends on Test 2, but Test 2 isn't in the path
   - **Impact**: Alternative paths won't execute correctly
   - **Recommendation**: Fix alternative_paths to include all dependencies
   - **Correct quick_validation**: `[test_verify_route_monitor_deployment, test_create_routemonitor_cr, test_verify_probe_metrics, test_cleanup_routemonitor_cr]`

3. **No Cleanup Tracking System**
   - **Issue**: If Test 2 runs but Test 5 doesn't (due to failure, user cancellation, etc.), resources are orphaned
   - **Impact**: Cluster pollution, confusing state for next test run
   - **Recommendation**: Implement cleanup tracking:
     - Track which tests ran and need cleanup
     - Prompt user at end: "2 tests require cleanup. Run now? (yes/no/later)"
     - Save cleanup_needed.json for later execution
     - Auto-run cleanup tests at end of successful run

4. **No Production Environment Detection**
   - **Issue**: Tests check can_run_in_production but don't enforce it
   - **Impact**: Risk of accidentally running modifying tests in prod
   - **Recommendation**: Add environment detection:
     ```bash
     # Detect production clusters
     CLUSTER_ID=$(oc get clusterversion -o jsonpath='{.items[0].spec.clusterID}')
     CLUSTER_NAME=$(oc get infrastructure cluster -o jsonpath='{.status.infrastructureName}')

     # Check against known production patterns
     if [[ "$CLUSTER_NAME" =~ prod|production ]]; then
       echo "❌ PRODUCTION CLUSTER DETECTED"
       echo "Test test_create_routemonitor_cr cannot run in production"
       exit 1
     fi
     ```

5. **Missing Environment Validation Script**
   - **Issue**: Pre-execution checklist is manual
   - **Impact**: Users may skip checks, tests fail due to environment issues
   - **Recommendation**: Create pre_flight_check.sh:
     ```bash
     #!/bin/bash
     # Verify all prerequisites before test execution
     check_rmo_deployed
     check_api_accessible
     check_prometheus_running
     check_rbac_permissions
     check_not_production
     ```

### 🟢 Medium Priority (Quality of Life)

6. **Test Duration Tracking**
   - **Issue**: Estimated times are static, no actual duration captured
   - **Recommendation**: Wrap each test execution in timing:
     ```bash
     START_TIME=$(date +%s)
     run_test_1
     END_TIME=$(date +%s)
     DURATION=$((END_TIME - START_TIME))
     echo "Test 1 actual duration: ${DURATION}s (estimated: 900s)"
     ```

7. **No Partial Success Handling**
   - **Issue**: If 3 of 5 steps pass, test is marked "failed" with no granularity
   - **Recommendation**: Report step-level pass/fail:
     ```
     Test 2: test_create_routemonitor_cr - PARTIAL SUCCESS
       Step 1: Create routemonitor.yaml - ✅ PASS
       Step 2: Apply RouteMonitor CR - ✅ PASS
       Step 3: Wait for ServiceMonitor - ❌ FAIL (timeout)
       Step 4: Verify probe metrics - ⏭️ SKIPPED (previous step failed)
     ```

8. **No Evidence Collection**
   - **Issue**: When tests fail, no artifacts are automatically collected
   - **Recommendation**: On failure, capture:
     - oc get routemonitors -o yaml
     - oc get servicemonitors -o yaml
     - oc logs deployment/route-monitor-operator (last 100 lines)
     - curl Prometheus query results
     - Save all to `evidence/test_2_failure_20260218_1430/`

9. **No Retry Logic for Transient Failures**
   - **Issue**: Network timeouts, API rate limits cause false failures
   - **Recommendation**: Implement retry with backoff for idempotent operations:
     ```bash
     retry 3 oc get routemonitor console-route-monitor
     retry 5 curl $RHOBS_API_URL/api/v1/probes
     ```

10. **Missing Test Report Template**
    - **Issue**: No standardized way to report results to stakeholders
    - **Recommendation**: Generate HTML/Markdown report:
      - Executive summary (5 tests, 4 passed, 1 failed, 100 min)
      - Per-test results with timestamps
      - Failures with screenshots/logs
      - Recommendations for follow-up

---

## Execution Safety Enforcement

### Recommended Workflow

```bash
# 1. Pre-flight validation
./scripts/pre_flight_check.sh || exit 1

# 2. Environment confirmation
echo "Cluster: $(oc whoami --show-server)"
echo "User: $(oc whoami)"
read -p "Confirm this is a TEST cluster (yes/no): " CONFIRM
[[ "$CONFIRM" != "yes" ]] && exit 1

# 3. Run read-only tests first
run_test test_verify_route_monitor_deployment

# 4. Prompt before modifying tests
echo "⚠️ Next tests will MODIFY cluster state"
echo "Affected: RouteMonitor CRs, ServiceMonitors, RHOBS API probes"
read -p "Proceed with modifying tests? (yes/no): " CONFIRM
[[ "$CONFIRM" != "yes" ]] && exit 1

# 5. Track cleanup needs
CLEANUP_NEEDED=()

# 6. Run modifying test with cleanup tracking
run_test test_create_routemonitor_cr && CLEANUP_NEEDED+=("test_cleanup_routemonitor_cr")
run_test test_create_probe_via_api && CLEANUP_NEEDED+=("MANUAL: Delete probe via API")

# 7. Run validation test
run_test test_verify_probe_metrics

# 8. Prompt for cleanup
if [ ${#CLEANUP_NEEDED[@]} -gt 0 ]; then
  echo "⚠️ ${#CLEANUP_NEEDED[@]} cleanup operations needed:"
  printf '  - %s\n' "${CLEANUP_NEEDED[@]}"
  read -p "Run cleanup now? (yes/no): " CLEANUP_NOW
  if [[ "$CLEANUP_NOW" == "yes" ]]; then
    run_test test_cleanup_routemonitor_cr
    echo "Manual cleanup: Delete probe via API (see instructions)"
  else
    echo "⚠️ Remember to run cleanup later!"
    save_cleanup_state cleanup_needed.json
  fi
fi

# 9. Generate report
generate_html_report results.html
```

---

## Next Steps

### Before Running Tests

1. ✅ Review this execution plan
2. ❌ **Create test_cleanup_probe_via_api** (critical gap)
3. ❌ Fix alternative_paths dependencies
4. ❌ Implement cleanup tracking system
5. ❌ Add production environment detection
6. ❌ Create pre_flight_check.sh script

### During Execution

1. Run tests in recommended sequence
2. Confirm safety prompts carefully
3. Track cleanup requirements
4. Capture evidence on failures
5. Document actual durations vs estimates

### After Execution

1. Run all required cleanup tests
2. Verify cluster returned to original state
3. Generate comprehensive report
4. Update test plan based on findings
5. Add newly identified edge cases as test variants

---

## Conclusions

The RHOBS-next test plan is **structurally sound** with good safety classifications and dependency management. However, **critical gaps exist**:

- ❌ Test 4 has no cleanup test (resource leak risk)
- ⚠️ Alternative paths have broken dependencies
- ⚠️ No automated cleanup tracking

**Recommendation**: **DO NOT run tests in current state** until test_cleanup_probe_via_api is created. Read-only tests (1, 3) can be run safely, but Test 2 and 4 should not be executed without proper cleanup mechanisms.

Once gaps are addressed, this test plan provides excellent coverage of RHOBS Next synthetic monitoring components with strong educational value and safety controls.
