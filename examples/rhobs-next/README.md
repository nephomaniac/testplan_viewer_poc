# RHOBS-Next Test Plan Example

**Created**: 2026-02-18
**Epic**: [SREP-3109](https://issues.redhat.com/browse/SREP-3109) - RHOBS Next and Synthetic Monitoring E2E Testing
**Methodology**: testplan-generator → testplan-educator → HTML generation → testplan-executor review

---

## What This Example Demonstrates

This comprehensive example showcases the **complete test plan workflow** for the RHOBS Next synthetic monitoring stack:

1. ✅ **Repository Analysis** - Explored rhobs-synthetics-api, rhobs-synthetics-agent, and route-monitor-operator codebases
2. ✅ **JIRA Integration** - Incorporated requirements from SREP-3109 epic and child tickets
3. ✅ **Test Generation** - Created 5 comprehensive test cases with safety classifications
4. ✅ **Educational Enhancement** - Added hyperlinks, learning notes, and conceptual overviews
5. ✅ **HTML Generation** - Produced interactive web-based test plan viewer
6. ✅ **Execution Review** - Analyzed safety, dependencies, and execution approach (tests NOT executed)
7. ✅ **Continuous Improvement** - Identified 15 improvements to the repository tooling

---

## Files in This Example

| File | Purpose | Lines | Description |
|------|---------|-------|-------------|
| **rhobs-next-testplan.json** | Test Plan | 3,847 | Complete test plan with 5 tests, 10 concepts, safety attributes, educational content |
| **rhobs-next-testplan.html** | Interactive Viewer | Generated | Web-based viewer with progress tracking, hyperlinks, collapsible sections |
| **EXECUTION-REVIEW.md** | Execution Analysis | 870 | Comprehensive review of test execution approach, safety analysis, gaps identified |
| **IMPROVEMENTS.md** | Enhancement Documentation | 467 | Specific improvements identified during test plan creation |
| **ENHANCEMENT_SUMMARY.md** | Educational Enhancements | 274 | Statistics and details of educational content additions |
| **STRUCTURE_FIXES.md** | JSON Structure Fixes | 312 | Documentation of Go model alignment fixes |
| **README.md** | This File | - | Overview and guide |

---

## Test Plan Contents

### Metadata
- **Title**: RHOBS Next and Synthetic Monitoring E2E Test Plan
- **Target Audience**: SRE Engineers, QE Team, Platform Engineers, New Engineers
- **Estimated Time**: 8-12 hours
- **Learning Objectives**: 7 comprehensive objectives covering RHOBS architecture, probe lifecycle, multi-tenancy, metrics forwarding

### Prerequisites Defined
- **Required Access**: 5 access requirements (cluster admin, stage environment, backplane, GitHub, Prow/CI)
- **Required Tools**: 6 tools with installation links and verification commands
- **Required Knowledge**: 5 knowledge areas with difficulty levels and learning resources
- **Environment Setup**: 3 environment variables, 6 tools configured

### Concepts Explained (10 Total)
Each concept includes:
- Title and comprehensive description
- "Why it matters" section
- Related tests
- Further reading links (documentation, tutorials)

**Key Concepts**:
1. Synthetic Monitoring
2. Route Monitor Operator (RMO)
3. RHOBS Synthetics API
4. RHOBS Synthetics Agent
5. Probe Custom Resource
6. Blackbox Exporter
7. HyperShift HostedControlPlane
8. osde2e Testing Framework
9. Prometheus ServiceMonitor
10. Metrics Forwarding to RHOBS

### Test Cases (5 Total)

| # | Test ID | Category | Difficulty | Duration | Safety | Status |
|---|---------|----------|------------|----------|--------|--------|
| 1 | test_verify_route_monitor_deployment | validation | beginner | 15 min | ✅ Read-only | ✅ Complete |
| 2 | test_create_routemonitor_cr | configuration | intermediate | 25 min | ⚠️ Modifies-state | ✅ Complete |
| 3 | test_verify_probe_metrics | validation | intermediate | 20 min | ✅ Read-only | ✅ Complete |
| 4 | test_create_probe_via_api | integration | advanced | 30 min | ⚠️ Modifies-state | ⚠️ Missing cleanup |
| 5 | test_cleanup_routemonitor_cr | cleanup | beginner | 10 min | ⚠️ Cleanup | ✅ Complete |

**Safety Breakdown**:
- Read-only tests: 2 (safe for production)
- Modifies-state tests: 2 (test environments only)
- Cleanup tests: 1 (restores state)
- **Gap Identified**: Test 4 has no cleanup test (see EXECUTION-REVIEW.md)

### Educational Content Statistics

**Per Test** (averages):
- Learning objectives: 4-5 per test
- Conceptual overview: Yes (for intermediate/advanced tests)
- Test steps: 9-10 steps per test
- Learning notes: 100% step coverage
- Common errors documented: 78% step coverage
- Troubleshooting issues: 3-5 per test
- Reference links: 4-6 per test

**Technical Terms Hyperlinked**: 14 unique terms
- First occurrence includes inline definition and "learn more" link
- Examples: Custom Resource, ServiceMonitor, Operator, CRD, blackbox-exporter, PromQL, HyperShift, multi-tenancy

**Total Educational Enhancement**:
- 141 KB of JSON (comprehensive detail)
- 46 steps with learning notes
- 58 common errors documented
- 14 concepts with hyperlinked explanations
- Production-relevant context throughout

---

## How This Example Was Created

### Phase 1: Repository Exploration (30 minutes)
Using the Explore agent with "medium" thoroughness:
- Analyzed 3 RHOBS repositories
- Identified key components and their responsibilities
- Mapped end-to-end data flow (RouteMonitor → RMO → API → Agent → Probe → blackbox-exporter → Prometheus)
- Located E2E tests and documentation
- Extracted integration points and dependencies

**Key Findings**:
- RMO handles both standard clusters (ServiceMonitor) and HyperShift (API integration)
- Agent reconciles from multiple API endpoints simultaneously with deduplication
- Probe lifecycle: pending → active → running with status feedback loop
- URL validation via HEAD requests prevents creating broken probes

### Phase 2: Test Case Generation (45 minutes)
Following testplan-generator skill methodology:
- Created coverage matrix (functional, environmental, integration)
- Generated 5 test cases covering validation, configuration, integration, cleanup
- Classified system impact for each test (read-only vs modifies-state)
- Auto-linked safety attributes (can_run_in_production, requires_confirmation, cleanup requirements)
- Documented affected resources and risk levels

**Challenges**:
- Initial JSON structure didn't match Go data model (prerequisites, references)
- Learning path referenced tests not yet created
- **Improvement identified**: Generator needs schema validation before writing JSON

### Phase 3: Educational Enhancement (60 minutes)
Following testplan-educator skill methodology:
- Added hyperlinks for all technical terms (first occurrence only)
- Enhanced learning_note for every step with "why" and "what happens internally"
- Added conceptual_overview to intermediate/advanced tests
- Expanded common_errors with cause, solution, and learning_point
- Integrated production context ("when would you use this in practice?")

**Statistics**:
- 14 unique technical terms hyperlinked
- 100% step coverage for learning notes
- 78% step coverage for common errors
- 3 tests with conceptual overviews

**Improvements Identified**:
- HTML template needs markdown rendering for hyperlinks
- Go model needs ConceptualOverview field
- Skill needs parameter to control enhancement level

### Phase 4: Structure Alignment (30 minutes)
Fixed JSON structure to match Go data model:
- **Prerequisites**: Converted arrays to structured objects (RequiredAccess, RequiredKnowledge, EnvironmentSetup)
- **References**: Converted array of objects to categorized arrays (documentation, related_code, sops)
- **Metadata**: Renamed last_updated → date
- **Removed**: conceptual_overview (not in Go model yet - improvement #14)

**Lesson Learned**: Schema validation before generation would catch these early

### Phase 5: HTML Generation (5 minutes)
After structure fixes:
```bash
./build/testplan-viewer -i examples/rhobs-next/rhobs-next-testplan.json \
  -o examples/rhobs-next/rhobs-next-testplan.html
```

**Result**: ✅ Success - 141 KB HTML file with 5 tests, 10 concepts, interactive features

**Limitations Discovered**:
- Hyperlinks render as plain text (markdown not parsed) - improvement #4
- No inline conceptual overview rendering (field not in model) - improvement #14
- localStorage key is generic "testplan_progress" (could conflict) - minor issue

### Phase 6: Execution Review (90 minutes)
Following testplan-executor skill methodology:
- Analyzed dependencies and execution order (topological sort)
- Validated safety classifications
- Identified cleanup coverage gaps (Test 4 missing cleanup)
- Reviewed alternative paths (found invalid dependencies)
- Created pre-execution checklists
- Documented execution flow walkthrough for each test
- Identified 10 improvements to execution approach

**Critical Findings**:
- ❌ Test 4 has no cleanup test - creates resource leak
- ⚠️ Alternative paths include tests without their dependencies
- ⚠️ No automated cleanup tracking system
- ⚠️ No production environment detection

**Recommendations**:
- **DO NOT execute Test 4** until cleanup test created
- Read-only tests (1, 3) safe to run anytime
- Modifying tests (2, 4) require test environments and cleanup tracking

---

## Key Learnings

### What Worked Well ✅

1. **Comprehensive Repository Analysis** - Explore agent provided excellent understanding of RHOBS architecture
2. **Safety-First Approach** - Every test classified, cleanup requirements clear
3. **Educational Value** - Hyperlinks and learning notes transform tests into learning resources
4. **Incremental Validation** - Found structure issues before execution, not during
5. **Skills Provide Structure** - Following skill methodology ensured completeness

### Challenges Encountered ⚠️

1. **Schema Drift** - Generated JSON didn't match Go model (prerequisites, references)
2. **Manual Structure Fixes** - Significant work to align with model after generation
3. **Missing Features** - Hyperlinks don't render (markdown support needed)
4. **Incomplete Cleanup** - Generator created modifying test without cleanup pair
5. **Alternative Path Validation** - No check for dependency completeness

### Improvements Identified (15 Total)

**Critical (3)**:
- Schema validation in testplan-generator
- Auto-generate cleanup tests for modifying tests
- Validate alternative paths include dependencies

**High Priority (7)**:
- HTML markdown rendering for hyperlinks
- Cleanup tracking system
- Production environment detection
- Enhanced cleanup validation
- JSON schema definition
- Evidence collection on failure
- Validator checks cleanup test properties

**Medium Priority (5)**:
- Test duration tracking
- Partial success handling
- Test report template
- ConceptualOverview field in Go model
- Parameterized skills

**See**: [REPOSITORY-IMPROVEMENTS.md](/Users/maclark/sandbox/testplan_tools_poc/REPOSITORY-IMPROVEMENTS.md) for full details

---

## How to Use This Example

### Viewing the Test Plan

**Open the HTML viewer**:
```bash
open examples/rhobs-next/rhobs-next-testplan.html
```

**Features**:
- Progress tracking (saved in browser localStorage)
- Collapsible test sections
- Search functionality
- Difficulty filtering
- Copy-to-clipboard for commands
- Safety badges (Read-Only, Modifies-State)

### Executing the Tests (NOT RECOMMENDED YET)

**⚠️ WARNING**: Do not execute Test 4 (test_create_probe_via_api) until cleanup test is created. See EXECUTION-REVIEW.md for details.

**Safe to execute**:
- Test 1: test_verify_route_monitor_deployment (read-only)
- Test 3: test_verify_probe_metrics (read-only, requires Test 2 first)

**Requires test environment + cleanup**:
- Test 2: test_create_routemonitor_cr + Test 5: test_cleanup_routemonitor_cr

**Before executing ANY test**:
1. Read [EXECUTION-REVIEW.md](EXECUTION-REVIEW.md)
2. Verify environment is NON-PRODUCTION
3. Review pre-execution checklist
4. Ensure cleanup tests exist for all modifying tests

### Learning from This Example

**For Test Plan Authors**:
- Study the comprehensive structure (metadata, prerequisites, concepts, test cases)
- Note the safety classification pattern (read-only, modifies-state, cleanup)
- Observe how educational content is integrated (learning notes, hyperlinks, conceptual overviews)
- Review the cleanup pairing pattern (Test 2 ← → Test 5)

**For Repository Contributors**:
- Review the 15 improvements identified
- Understand the schema alignment challenges
- See the gap between generated JSON and Go model
- Learn from the execution safety analysis

**For Educators**:
- Examine how technical terms are hyperlinked and defined
- Study the "why this matters" and "real-world context" patterns
- Note the common errors documentation approach
- Review the conceptual overview structure

---

## Next Steps

### Immediate (Before Execution)

1. ✅ Review EXECUTION-REVIEW.md
2. ❌ Create test_cleanup_probe_via_api (critical gap)
3. ❌ Fix alternative_paths dependencies
4. ❌ Add production environment detection script

### Short Term (Repository Improvements)

1. ❌ Implement schema validation in testplan-generator skill
2. ❌ Add markdown rendering to HTML template
3. ❌ Create cleanup tracking system in testplan-executor
4. ❌ Add ConceptualOverview field to Go model

### Medium Term (Additional Test Coverage)

The current test plan covers basic validation and integration. Additional tests needed:

5. test_verify_agent_deployment - Verify synthetics agent is running
6. test_agent_probe_reconciliation - Verify agent creates Probe CRs from API
7. test_url_validation - Verify agent validates URLs before creating Probes
8. test_hostedcontrolplane_integration - Verify RMO creates probes for HyperShift clusters
9. test_osde2e_cluster_creation - Set up osde2e test harness
10. test_osde2e_smoke_suite - Run full smoke test suite in osde2e

These were in the original learning_path.sequence but not yet created.

---

## References

### JIRA
- [SREP-3109](https://issues.redhat.com/browse/SREP-3109) - Parent Epic
- [SREP-3117](https://issues.redhat.com/browse/SREP-3117) - Set up osde2e test harness
- [SREP-3118](https://issues.redhat.com/browse/SREP-3118) - Implement metrics verification smoke test
- [SREP-3119](https://issues.redhat.com/browse/SREP-3119) - Implement logs verification smoke test
- [SREP-3120](https://issues.redhat.com/browse/SREP-3120) - Implement synthetics verification smoke test

### Repositories Analyzed
- [rhobs/api](https://github.com/rhobs/api) - RHOBS Synthetics API
- [rhobs/rhobs-synthetics-agent](https://github.com/rhobs/rhobs-synthetics-agent) - Synthetics Agent
- [openshift/route-monitor-operator](https://github.com/openshift/route-monitor-operator) - Route Monitor Operator

### Documentation
- [Prometheus Blackbox Exporter](https://github.com/prometheus/blackbox_exporter)
- [Prometheus Operator](https://prometheus-operator.dev/)
- [HyperShift Documentation](https://hypershift-docs.netlify.app/)
- [osde2e Repository](https://github.com/openshift/osde2e)

---

## Metrics

- **Time to Create**: ~4 hours (exploration, generation, enhancement, review, documentation)
- **Test Cases**: 5 (demonstrating read-only, modifies-state, cleanup patterns)
- **Concepts**: 10 (comprehensive RHOBS architecture coverage)
- **Learning Objectives**: 7 top-level + 20+ per-test objectives
- **Technical Terms Hyperlinked**: 14
- **Total JSON Size**: 141 KB
- **Educational Content**: 46 learning notes, 58 common errors, 3 conceptual overviews
- **Improvements Identified**: 15 (3 critical, 7 high, 5 medium priority)

---

## Conclusion

This example successfully demonstrates the **complete test plan workflow** from repository analysis to execution review, producing a comprehensive, educational test plan with strong safety controls. The process identified significant improvements to repository tooling, particularly around schema validation, cleanup automation, and educational content rendering.

**Status**: ✅ Test plan ready for review | ⚠️ Not ready for execution (missing cleanup test for Test 4)

**Recommended Action**: Implement critical improvements #1-3 before creating additional test plans or executing tests.

🎉 **Ready for the next operator codebase!**
