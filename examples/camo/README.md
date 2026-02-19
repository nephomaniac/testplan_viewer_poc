# Configure AlertManager Operator (CAMO) Test Plan Example

## Overview

This directory contains a comprehensive test plan for the Configure AlertManager Operator (CAMO), demonstrating the full test plan generation, education, and execution workflow using the skills in this repository.

## What is CAMO?

CAMO is a Kubernetes operator that dynamically configures Prometheus Alertmanager based on secrets and configmaps. It watches for:
- `pd-secret` - Triggers PagerDuty integration
- `goalert-secret` - Triggers GoAlert integration
- `dms-secret` - Triggers Dead Man's Snitch integration
- `managed-namespaces` ConfigMap - Defines which namespaces to monitor
- `ocp-namespaces` ConfigMap - Defines OCP platform namespaces

When these resources are created/updated/deleted, CAMO automatically reconfigures Alertmanager routing.

## Files in This Directory

- **camo-testplan.json** - Complete test plan with 4 test cases covering:
  - Deployment validation (read-only)
  - Alertmanager config inspection (read-only)
  - PagerDuty secret creation (modifies-state, with cleanup)
  - PagerDuty secret cleanup (restores state)

- **camo-testplan.html** - Interactive HTML viewer generated from the JSON

- **EXECUTION-REVIEW.md** - Detailed review of test execution approach with improvements identified

## Test Plan Structure

### Metadata
- **Title**: Configure AlertManager Operator (CAMO) Test Plan
- **Target Audience**: SRE Engineers, QE Team, Platform Engineers
- **Estimated Time**: 6-8 hours (full plan)
- **Learning Objectives**:
  - Understand CAMO's dynamic configuration model
  - Learn secret-based alert routing
  - Master troubleshooting techniques

### Concepts Explained (7)
1. CAMO Operator architecture
2. Alertmanager Configuration Secret
3. Config Validation
4. PagerDuty Integration
5. GoAlert Integration
6. Dead Man's Snitch Integration
7. Cluster Readiness Checks

### Test Cases (4)

#### 1. test_verify_camo_deployment (Read-Only)
- **Purpose**: Verify operator is installed and running
- **Duration**: 10 minutes
- **Safety**: ✅ Safe to run in production
- **Steps**: Check deployment, pod, CSV status

#### 2. test_read_alertmanager_config (Read-Only)
- **Purpose**: Inspect Alertmanager configuration
- **Duration**: 15 minutes
- **Safety**: ✅ Safe to run in production
- **Steps**: Extract config from secret, parse YAML, list receivers

#### 3. test_create_pagerduty_secret (Modifies-State)
- **Purpose**: Create pd-secret and verify CAMO reconfigures Alertmanager
- **Duration**: 20 minutes
- **Safety**: ⚠️ Modifies state, requires cleanup
- **Cleanup Test**: test_cleanup_pagerduty_secret
- **Steps**: Create secret, wait for reconciliation, verify receiver added, check metrics

#### 4. test_cleanup_pagerduty_secret (Cleanup)
- **Purpose**: Remove pd-secret and restore original Alertmanager config
- **Duration**: 10 minutes
- **Safety**: ✅ Restores state to pre-test condition
- **Steps**: Delete secret, verify receiver removed, check metrics

## How This Example Was Created

### 1. Test Plan Generation (testplan-generator skill)
Analyzed CAMO codebase:
```bash
# Explored code structure
find /path/to/camo -name "*.go" -type f

# Read E2E tests
cat test/e2e/configure_alertmanager_operator_tests.go

# Read controller code
cat controllers/secret_controller.go

# Read documentation
cat README.md LOCAL_TESTING.md
```

Generated test plan JSON following testplan-generator methodology:
- Identified all watched resources (secrets, configmaps)
- Created test coverage matrix (deployment, configuration, metrics, troubleshooting)
- Classified tests by system impact (read-only vs modifies-state)
- Created backup/cleanup patterns for state-modifying tests
- Added safety attributes to all tests

### 2. Educational Enhancement (testplan-educator skill)
Added educational content:
- Learning objectives for each test
- Concept explanations with real-world analogies
- Step-by-step learning notes
- Common beginner mistakes and how to avoid them
- Troubleshooting guidance

### 3. HTML Generation
```bash
./build/testplan-viewer -i examples/camo/camo-testplan.json -o examples/camo/camo-testplan.html
```

Result: Interactive HTML with:
- Safety badges (read-only, modifies-state)
- Expandable test cards
- Copy-to-clipboard commands
- Progress tracking
- Search and filtering

### 4. Execution Review (testplan-executor skill)
Reviewed execution flow:
- Pre-execution safety validation
- Production environment checks
- Backup requirement verification
- User confirmation for risky tests
- Post-execution cleanup prompts
- State restoration verification

## Safety Features Demonstrated

### Read-Only Tests
- ✅ Marked with "👁️ Read-Only" badge
- ✅ `can_run_in_production: true`
- ✅ No cleanup required
- ✅ Safe to retry

### Modifies-State Tests
- ⚠️ Marked with "✏️ Modifies State" badge
- ⚠️ `can_run_in_production: false`
- ⚠️ `requires_confirmation: true`
- ⚠️ Warning message displayed
- ⚠️ Cleanup test required
- ⚠️ Lists affected resources

### Cleanup Tests
- 🧹 Category: "cleanup"
- 🧹 `restores_state: true`
- 🧹 Linked from modifying test
- 🧹 Verifies restoration completed

## Viewing the Test Plan

### Option 1: HTML Viewer
```bash
open examples/camo/camo-testplan.html
```

Interactive features:
- Click test cards to expand/collapse
- Check off completed tests (saved in browser)
- Copy commands with one click
- Search and filter tests
- Track progress

### Option 2: Raw JSON
```bash
cat examples/camo/camo-testplan.json | jq .
```

## Running the Tests

**⚠️ IMPORTANT**: These tests are examples and have NOT been executed. Review the EXECUTION-REVIEW.md before running.

To execute (conceptual):
```bash
# 1. Verify prerequisites
oc auth can-i create secrets -n openshift-monitoring

# 2. Run read-only tests safely
# Follow test_verify_camo_deployment steps
# Follow test_read_alertmanager_config steps

# 3. Run modifying tests with cleanup
# Follow test_create_pagerduty_secret steps
# ⚠️ IMPORTANT: Follow test_cleanup_pagerduty_secret afterward
```

## Improvements Identified

During this exercise, we identified 10 improvements to the test plan system:

### High Priority
1. **Alternative paths data model** - Add descriptions to alternative learning paths
2. **Cleanup tracking** - Track which tests need cleanup and remind at end
3. **Safety checklist rendering** - Show pre-flight checklist in HTML

### Medium Priority
4. **Command intrusiveness detection** - Auto-detect if commands modify state
5. **Auto-generate cleanup tests** - Generate cleanup tests automatically
6. **Circular dependency detection** - Prevent impossible execution orders

### Lower Priority
7. **Alternative paths display** - Show alternative paths in HTML
8. **Enhanced validation** - Validate cleanup test properties
9. **Progress tracking** - Track timestamps and notes
10. **Skill interdependencies** - Document skill relationships

See EXECUTION-REVIEW.md for details on each improvement.

## Lessons Learned

### What Worked Well
✅ Safety classification (read-only vs modifies-state) is clear and effective
✅ Backup/cleanup pattern prevents state corruption
✅ Validator catches configuration errors early (missing tests, wrong structure)
✅ HTML viewer makes tests accessible and interactive
✅ Concept explanations help beginners understand why steps matter

### Areas for Improvement
⚠️ Alternative paths need richer structure (descriptions)
⚠️ Cleanup tracking could be automated
⚠️ Command analysis could auto-detect intrusiveness
⚠️ Cleanup test generation could be automated

### Best Practices Demonstrated
1. **Start with read-only tests** - Build confidence before modifying state
2. **Every modifying test needs cleanup** - No exceptions
3. **Explain why, not just what** - Learning notes add context
4. **Safety first** - Classify impact, warn users, require confirmation
5. **Validate early** - Catch errors at parse time, not runtime

## Next Steps

To expand this example:

1. **Add more test cases**:
   - RBAC validation
   - Metrics service verification
   - GoAlert secret creation/cleanup
   - DMS secret creation/cleanup
   - Config validation metric testing
   - Cluster readiness checks

2. **Implement improvements**:
   - Enhanced alternative paths model
   - Automated cleanup tracking
   - Command intrusiveness detection

3. **Create execution automation**:
   - Shell script that follows testplan-executor logic
   - Automated safety checks
   - Cleanup tracking and reminders

## Related Documentation

- CAMO Repository: https://github.com/openshift/configure-alertmanager-operator
- Test Plan Viewer: ../../README.md
- Skills: ../../skills/README.md
- Safety Guide: ../../docs/TEST-SAFETY-AND-STATE-MANAGEMENT.md

## Contributing

To add more CAMO tests:

1. Follow testplan-generator methodology
2. Classify system impact correctly
3. Create cleanup tests for modifying operations
4. Add educational content
5. Update learning_path.sequence
6. Regenerate HTML

```bash
# After editing JSON
./build/testplan-viewer -i examples/camo/camo-testplan.json -o examples/camo/camo-testplan.html
```
