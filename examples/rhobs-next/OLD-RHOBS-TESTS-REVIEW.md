# Review of Old RHOBS Tests (from examples/rhobs/)

**Date**: 2026-02-19
**Purpose**: Document tests from old examples/rhobs/ directory before removal

---

## Tests in Old RHOBS Directory

The `examples/rhobs/` directory contained three test plans that are NOT in rhobs-next:

### 1. test_1_cluster_creation
**Purpose**: Create osde2e test harness cluster
**Coverage**: Cluster setup, osde2e framework initialization

### 2. test_2_logs_flow
**Purpose**: Verify logs verification smoke test
**Coverage**: Logging pipeline, log collection, log forwarding

### 3. test_3_synthetic_monitoring
**Purpose**: Verify synthetic monitoring smoke test
**Coverage**: Synthetic probe creation, metrics collection

---

## Tests in RHOBS-Next

The `examples/rhobs-next/rhobs-next-testplan.json` contains different tests focused on Route Monitor Operator and Synthetics API:

1. **test_verify_route_monitor_deployment** - Validate RMO deployment
2. **test_create_routemonitor_cr** - Create RouteMonitor custom resource
3. **test_verify_probe_metrics** - Verify probe metrics collection
4. **test_create_probe_via_api** - Create probe via API
5. **test_cleanup_routemonitor_cr** - Cleanup RouteMonitor CR

---

## Gap Analysis

### Tests Not Yet in RHOBS-Next

#### Missing: Cluster Creation Test
**From**: test_1_cluster_creation
**Coverage Gap**: Setting up osde2e test harness
**JIRA Reference**: [SREP-3117](https://issues.redhat.com/browse/SREP-3117) - Set up osde2e test harness
**Recommendation**: Should be added to rhobs-next as it's part of the SREP-3109 epic

#### Missing: Logs Flow Test
**From**: test_2_logs_flow
**Coverage Gap**: Logging pipeline validation
**JIRA Reference**: [SREP-3119](https://issues.redhat.com/browse/SREP-3119) - Implement logs verification smoke test
**Recommendation**: Should be added to rhobs-next as it's part of the SREP-3109 epic

#### Missing: General Synthetic Monitoring Test
**From**: test_3_synthetic_monitoring
**Coverage Gap**: End-to-end synthetic monitoring smoke test
**JIRA Reference**: [SREP-3120](https://issues.redhat.com/browse/SREP-3120) - Implement synthetics verification smoke test
**Recommendation**: Should be added to rhobs-next as it's part of the SREP-3109 epic

---

## Recommendation

The old rhobs tests represent test cases that were planned but not yet implemented in rhobs-next. According to the JIRA epic (SREP-3109), the following tests should be added to rhobs-next:

### Tests to Add (from JIRA SREP-3109 child tickets):

1. **test_osde2e_cluster_creation** (SREP-3117)
   - Set up osde2e test harness
   - Create HCP cluster for testing
   - Configure test environment

2. **test_logs_verification_smoke** (SREP-3119)
   - Verify log collection from pods
   - Validate log forwarding to centralized logging
   - Check log query functionality

3. **test_synthetics_verification_smoke** (SREP-3120)
   - End-to-end synthetic monitoring test
   - Validate RouteMonitor → API → Agent → Probe → Metrics pipeline
   - Verify metrics appear in Prometheus/Thanos

4. **test_metrics_verification_smoke** (SREP-3118)
   - Verify metrics collection
   - Validate metrics forwarding
   - Check metrics query functionality

---

## Files in Old examples/rhobs/

1. **rhobs_test_plan.json** (54KB)
   - 3 test cases
   - Basic structure without educational content
   - v1 format

2. **rhobs_test_plan_v2.json** (34KB)
   - 3 test cases (same as v1)
   - Slightly different structure
   - v2 format

3. **RHOBS_Manual_Test_Plan.md** (52KB)
   - Manual markdown test plan
   - Original source document
   - Not in JSON format

---

## Action Items for RHOBS-Next

### Immediate
- ✅ Remove examples/rhobs/ directory (old v1 format)
- ⏳ Document missing tests (this file)

### Short Term (Add Missing Tests)
1. ⏳ Add test_osde2e_cluster_creation (SREP-3117)
2. ⏳ Add test_logs_verification_smoke (SREP-3119)
3. ⏳ Add test_synthetics_verification_smoke (SREP-3120)
4. ⏳ Add test_metrics_verification_smoke (SREP-3118)

### Medium Term (Enhance Existing Tests)
1. ⏳ Convert rhobs-next tests to new ChildTestCase format
2. ⏳ Add negative tests for each positive test
3. ⏳ Add RBAC tests
4. ⏳ Add load tests
5. ⏳ Add network tests

---

## Conclusion

The old examples/rhobs/ directory can be safely removed as:
1. The tests it contained represent gaps, not complete implementations
2. Those gaps are documented here and in JIRA
3. The rhobs-next directory is the canonical location for RHOBS test plans
4. The missing tests should be added to rhobs-next based on JIRA requirements

---

**Status**: Old rhobs directory removed
**Next**: Add missing tests to rhobs-next per JIRA epic SREP-3109
