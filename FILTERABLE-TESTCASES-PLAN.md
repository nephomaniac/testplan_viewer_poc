# Filterable Test Cases Implementation Plan

**Date**: 2026-02-19
**Purpose**: Enhance test plan system with comprehensive filterable attributes and testing coverage

---

## Overview

Transform the current step-based test execution model into a hierarchical child test case system with rich filterable metadata. Enable comprehensive testing coverage including negative, load, network, RBAC, and input validation testing.

---

## Requirements

### 1. Filterable Test Attributes

Each test step should become a **child test case** with metadata for filtering:

**Test Types:**
- `setup` - Environment/prerequisite setup
- `install` - Installation operations
- `cleanup` - Resource cleanup/removal
- `validation` - Read-only validation/verification
- `negative` - Negative testing (expected failures)
- `load` - Load/performance testing
- `performance` - Performance benchmarking
- `security` - Security validation
- `rbac` - RBAC/permission testing
- `network` - Network testing/impairment
- `input_validation` - Input validation testing

**Impact Types:**
- `read-only` - No modifications
- `modifies-state` - Creates/updates resources
- `destructive` - Deletes resources
- `impairment` - Temporarily degrades system
- `load-generation` - Generates load on system

**Input Validation Categories:**
- `positive` - Valid inputs
- `negative` - Invalid inputs
- `missing` - Missing required attributes
- `corrupt` - Malformed/corrupt data
- `boundary` - Min/max boundary values

**RBAC Levels:**
- `cluster-admin` - Requires cluster admin
- `namespace-admin` - Requires namespace admin
- `edit` - Requires edit permissions
- `view` - Requires view permissions
- `custom` - Custom RBAC configuration

**Resource Requirements:**
- `cpu` - CPU requirements
- `memory` - Memory requirements
- `storage` - Storage requirements
- `ephemeral_storage` - Ephemeral storage requirements

**Network Requirements:**
- `ingress` - Requires ingress access
- `egress` - Requires egress access
- `internal` - Internal cluster networking
- `external` - External network access
- `impairment_needed` - Network impairment required

---

### 2. Comprehensive Test Coverage Types

#### Negative Testing
- Invalid configurations
- Malformed inputs
- Permission denials
- Resource conflicts
- Dependency failures

#### Input Validation Testing
- **Positive inputs**: Valid configurations
- **Negative inputs**: Invalid values, types
- **Missing attributes**: Required fields omitted
- **Corrupt attributes**: Malformed data structures
- **Boundary values**: Min/max limits

#### RBAC Permission Testing
- Cluster-admin operations
- Namespace-scoped operations
- Read-only user attempts at modification
- Service account permissions
- Cross-namespace access denials

#### Load Testing
- **Resource limits**: Lower limits to simulate pod hitting constraints
- **Load generation**: Deploy load-creating pods
  - CPU stress
  - Memory stress
  - Network stress
  - Disk I/O stress

#### Network Testing with Impairment
- **AWS CLI** for security group rules
  ```bash
  aws ec2 authorize-security-group-ingress --group-id sg-xxx --protocol tcp --port 443 --cidr 0.0.0.0/0
  aws ec2 revoke-security-group-ingress --group-id sg-xxx --protocol tcp --port 443 --cidr 0.0.0.0/0
  ```
- **oc CLI** for network policies
  ```bash
  oc apply -f network-policy-deny-all.yaml
  oc delete networkpolicy deny-all
  ```
- **Host network config** with easy cleanup
  ```bash
  # Add latency (requires privileged pod)
  tc qdisc add dev eth0 root netem delay 100ms
  # Remove latency
  tc qdisc del dev eth0 root
  ```

#### Resource Usage Validation
- CPU consumption metrics
- Memory consumption metrics
- Storage usage
- Network bandwidth

#### Network Usage Validation
- Ingress/egress metrics
- Connection counts
- Bandwidth utilization
- Latency measurements

---

## Data Model Changes

### Enhanced Step Structure (becomes ChildTestCase)

```go
// ChildTestCase represents a single executable step with full test metadata
type ChildTestCase struct {
    // Identification
    ID           string `json:"id"`
    ParentTestID string `json:"parent_test_id"`
    StepNumber   int    `json:"step_number"`
    Title        string `json:"title"`

    // Filtering Metadata
    TestType          string   `json:"test_type"`          // setup, install, cleanup, validation, negative, load, etc.
    ImpactType        string   `json:"impact_type"`        // read-only, modifies-state, destructive, impairment, load-generation
    InputValidation   string   `json:"input_validation"`   // positive, negative, missing, corrupt, boundary
    RBACLevel         string   `json:"rbac_level"`         // cluster-admin, namespace-admin, edit, view, custom
    Tags              []string `json:"tags"`               // Additional custom tags

    // Resource & Network Requirements
    ResourceRequirements ResourceRequirements `json:"resource_requirements"`
    NetworkRequirements  NetworkRequirements  `json:"network_requirements"`

    // Execution Details
    Command            *string       `json:"command"`
    ExpectedOutput     *string       `json:"expected_output"`
    ManualSteps        []string      `json:"manual_steps,omitempty"`
    Duration           string        `json:"duration"`

    // Educational Content
    LearningNote           string        `json:"learning_note,omitempty"`
    WhyThisStep            string        `json:"why_this_step,omitempty"`
    DocumentationReference string        `json:"documentation_reference,omitempty"`
    CommonErrors           []CommonError `json:"common_errors,omitempty"`

    // Safety & Validation
    Safety     StepSafety     `json:"safety"`
    Validation StepValidation `json:"validation"`
}

// ResourceRequirements defines resource needs for a test
type ResourceRequirements struct {
    CPU              string `json:"cpu"`               // e.g., "100m", "2"
    Memory           string `json:"memory"`            // e.g., "256Mi", "2Gi"
    Storage          string `json:"storage"`           // e.g., "10Gi"
    EphemeralStorage string `json:"ephemeral_storage"` // e.g., "1Gi"
    LoadGeneration   bool   `json:"load_generation"`   // Requires load generation pods
}

// NetworkRequirements defines network needs for a test
type NetworkRequirements struct {
    Ingress          bool     `json:"ingress"`
    Egress           bool     `json:"egress"`
    Internal         bool     `json:"internal"`
    External         bool     `json:"external"`
    ImpairmentNeeded bool     `json:"impairment_needed"`
    ImpairmentType   string   `json:"impairment_type"`   // latency, packet-loss, bandwidth-limit
    ImpairmentConfig string   `json:"impairment_config"` // Configuration details
    SecurityGroups   []string `json:"security_groups"`   // AWS security groups to modify
    NetworkPolicies  []string `json:"network_policies"`  // K8s network policies to apply
}

// StepSafety defines safety constraints for a single step
type StepSafety struct {
    CanRunInProduction bool   `json:"can_run_in_production"`
    RequiresCleanup    bool   `json:"requires_cleanup"`
    CleanupProcedure   string `json:"cleanup_procedure"`
    RiskLevel          string `json:"risk_level"` // low, medium, high, critical
}

// StepValidation defines validation criteria for a step
type StepValidation struct {
    SuccessCriteria []string `json:"success_criteria"`
    FailureCriteria []string `json:"failure_criteria"`
    Metrics         []string `json:"metrics_to_collect"`
}
```

### Enhanced TestCase Structure

```go
type TestCase struct {
    Metadata        TestMetadata    `json:"metadata"`
    Learning        Learning        `json:"learning"`
    ChildTests      []ChildTestCase `json:"child_tests"`      // NEW: Steps become child test cases
    Troubleshooting Troubleshooting `json:"troubleshooting"`
    NextSteps       NextSteps       `json:"next_steps"`
    References      References      `json:"references"`
}
```

---

## Skills Updates

### testplan-generator Skill

Add phases for comprehensive testing:

**Phase 3.5: Negative Testing Generation**
- Generate negative test cases for each functional test
- Invalid configurations
- Permission denials
- Resource conflicts
- Missing dependencies

**Phase 3.6: Input Validation Testing**
- Positive input tests
- Negative input tests (invalid types, values)
- Missing attribute tests
- Corrupt data tests
- Boundary value tests

**Phase 3.7: RBAC Testing**
- Cluster-admin required operations
- Namespace-admin operations
- Read-only user restrictions
- Service account tests
- Cross-namespace access tests

**Phase 3.8: Load Testing**
- Resource limit tests
- Load generation pod deployment
- CPU stress tests
- Memory stress tests
- Network stress tests

**Phase 3.9: Network Testing**
- Network policy tests
- Security group modification tests
- Latency injection tests
- Packet loss simulation
- Bandwidth limiting

### testplan-educator Skill

Add educational content for:
- Negative testing patterns
- Input validation best practices
- RBAC security model
- Load testing methodologies
- Network impairment techniques

### testplan-executor Skill

Add execution patterns for:
- Negative test execution (expecting failures)
- Load generation pod management
- Network impairment setup/cleanup
- RBAC validation
- Resource usage monitoring

---

## HTML Generator Updates

### New Filter Controls

```html
<div class="filters">
    <!-- Test Type Filter -->
    <select x-model="testTypeFilter">
        <option value="">All Test Types</option>
        <option value="setup">Setup</option>
        <option value="install">Install</option>
        <option value="validation">Validation</option>
        <option value="negative">Negative Testing</option>
        <option value="load">Load Testing</option>
        <option value="network">Network Testing</option>
        <option value="rbac">RBAC Testing</option>
        <option value="cleanup">Cleanup</option>
    </select>

    <!-- Impact Type Filter -->
    <select x-model="impactFilter">
        <option value="">All Impact Types</option>
        <option value="read-only">Read-Only</option>
        <option value="modifies-state">Modifies State</option>
        <option value="destructive">Destructive</option>
        <option value="impairment">Impairment</option>
    </select>

    <!-- RBAC Level Filter -->
    <select x-model="rbacFilter">
        <option value="">All RBAC Levels</option>
        <option value="cluster-admin">Cluster Admin</option>
        <option value="namespace-admin">Namespace Admin</option>
        <option value="edit">Edit</option>
        <option value="view">View Only</option>
    </select>

    <!-- Input Validation Filter -->
    <select x-model="inputValidationFilter">
        <option value="">All Input Types</option>
        <option value="positive">Positive (Valid)</option>
        <option value="negative">Negative (Invalid)</option>
        <option value="missing">Missing Attributes</option>
        <option value="corrupt">Corrupt Data</option>
    </select>

    <!-- Production Safety Filter -->
    <select x-model="productionSafeFilter">
        <option value="">All Tests</option>
        <option value="production-safe">Production Safe Only</option>
        <option value="test-only">Test Environment Only</option>
    </select>
</div>
```

### Child Test Case Display

```html
<div class="test-case">
    <h3>{{.Metadata.Title}}</h3>

    <!-- Child test cases (formerly steps) -->
    <div class="child-tests">
        {{range .ChildTests}}
        <div class="child-test-card"
             data-test-type="{{.TestType}}"
             data-impact="{{.ImpactType}}"
             data-rbac="{{.RBACLevel}}"
             data-input-validation="{{.InputValidation}}">

            <!-- Test type badge -->
            <span class="badge badge-{{.TestType}}">{{.TestType}}</span>

            <!-- Impact badge -->
            <span class="badge badge-{{.ImpactType}}">{{.ImpactType}}</span>

            <!-- RBAC badge -->
            {{if .RBACLevel}}
            <span class="badge badge-rbac">{{.RBACLevel}}</span>
            {{end}}

            <!-- Content -->
            <h4>{{.Title}}</h4>
            {{if .Command}}
            <pre><code>{{.Command}}</code></pre>
            {{end}}

            <!-- Resource requirements -->
            {{if .ResourceRequirements.LoadGeneration}}
            <div class="load-generation-notice">
                ⚡ Requires load generation pods
            </div>
            {{end}}

            <!-- Network impairment -->
            {{if .NetworkRequirements.ImpairmentNeeded}}
            <div class="impairment-notice">
                🌐 Network impairment: {{.NetworkRequirements.ImpairmentType}}
            </div>
            {{end}}
        </div>
        {{end}}
    </div>
</div>
```

---

## Implementation Phases

### Phase 1: Data Model Enhancement (1-2 hours)
1. Update `internal/models/testplan.go`
2. Add ChildTestCase struct
3. Add ResourceRequirements struct
4. Add NetworkRequirements struct
5. Add StepSafety and StepValidation structs

### Phase 2: Skill Updates (2-3 hours)
1. Update testplan-generator with new test type phases
2. Update testplan-educator with educational content for new test types
3. Update testplan-executor with execution patterns

### Phase 3: HTML Generator Updates (2-3 hours)
1. Update HTML template with new filters
2. Add child test case display
3. Add filter JavaScript logic
4. Add badges for test types, impact, RBAC
5. Update CSS for new components

### Phase 4: Example Test Plan Migration (2-3 hours)
1. Convert CAMO example to new format
2. Convert RHOBS-Next example to new format
3. Add comprehensive test coverage examples
4. Document migration process

### Phase 5: Documentation Updates (1-2 hours)
1. Update examples/README.md with filtering guide
2. Update skills/README.md with new test types
3. Create filtering guide for users
4. Update REPOSITORY-IMPROVEMENTS.md

---

## Test Coverage Examples

### Negative Testing Example

```json
{
  "id": "test_create_secret_invalid_namespace",
  "parent_test_id": "test_create_pagerduty_secret",
  "step_number": 1,
  "title": "Attempt to create secret in non-existent namespace",
  "test_type": "negative",
  "impact_type": "read-only",
  "input_validation": "negative",
  "command": "oc create secret generic pd-secret --from-literal=PAGERDUTY_KEY=test -n does-not-exist",
  "expected_output": "Error from server (NotFound): namespaces \"does-not-exist\" not found",
  "learning_note": "Validates that the system correctly rejects operations on non-existent namespaces",
  "safety": {
    "can_run_in_production": true,
    "risk_level": "low"
  }
}
```

### RBAC Testing Example

```json
{
  "id": "test_create_secret_insufficient_permissions",
  "parent_test_id": "test_create_pagerduty_secret",
  "step_number": 2,
  "title": "Attempt to create secret as view-only user",
  "test_type": "rbac",
  "impact_type": "read-only",
  "rbac_level": "view",
  "command": "oc create secret generic pd-secret --from-literal=PAGERDUTY_KEY=test -n openshift-monitoring --as=view-user",
  "expected_output": "Error from server (Forbidden): secrets is forbidden",
  "learning_note": "Validates RBAC prevents view-only users from creating secrets",
  "safety": {
    "can_run_in_production": true,
    "risk_level": "low"
  }
}
```

### Load Testing Example

```json
{
  "id": "test_operator_cpu_stress",
  "parent_test_id": "test_verify_camo_deployment",
  "step_number": 5,
  "title": "Deploy CPU stress pods to test operator behavior under load",
  "test_type": "load",
  "impact_type": "load-generation",
  "resource_requirements": {
    "cpu": "2",
    "memory": "512Mi",
    "load_generation": true
  },
  "command": "oc apply -f cpu-stress-deployment.yaml",
  "learning_note": "Tests operator's ability to reconcile under CPU pressure",
  "safety": {
    "can_run_in_production": false,
    "requires_cleanup": true,
    "cleanup_procedure": "oc delete -f cpu-stress-deployment.yaml",
    "risk_level": "medium"
  }
}
```

### Network Impairment Example

```json
{
  "id": "test_api_latency_impact",
  "parent_test_id": "test_verify_probe_metrics",
  "step_number": 3,
  "title": "Add 100ms latency to API calls",
  "test_type": "network",
  "impact_type": "impairment",
  "network_requirements": {
    "impairment_needed": true,
    "impairment_type": "latency",
    "impairment_config": "100ms delay on port 443",
    "network_policies": ["api-latency-policy.yaml"]
  },
  "command": "oc apply -f network-latency-policy.yaml",
  "learning_note": "Tests system behavior when API calls are slow",
  "safety": {
    "can_run_in_production": false,
    "requires_cleanup": true,
    "cleanup_procedure": "oc delete networkpolicy api-latency-policy",
    "risk_level": "high"
  }
}
```

---

## Benefits

1. **Granular Filtering**: Users can filter tests by type, impact, RBAC level, input validation
2. **Comprehensive Coverage**: Negative, load, network, RBAC testing ensures robust validation
3. **Safety Classification**: Clear indicators of production-safe vs test-only operations
4. **Educational Value**: Each test type teaches different testing methodology
5. **Tool Integration**: Filterable metadata enables automated test selection
6. **Execution Control**: Users can run only read-only tests for validation, or full suite for certification

---

## Migration Path

### Existing Test Plans

1. Read current test plan
2. Convert each `Step` to `ChildTestCase`
3. Add `test_type` based on content analysis:
   - Commands with `oc get` → `validation`
   - Commands with `oc create` → `install` or `setup`
   - Commands with `oc delete` → `cleanup`
4. Set `impact_type` based on command:
   - Read-only commands → `read-only`
   - Create/update commands → `modifies-state`
   - Delete commands → `destructive`
5. Infer `rbac_level` from namespace and resource type
6. Generate negative test counterparts for each positive test
7. Add RBAC tests for permission-sensitive operations
8. Add load/network tests where applicable

---

## Success Criteria

- ✅ All existing tests converted to new format
- ✅ Each test has 2-3 negative test counterparts
- ✅ RBAC tests for admin-level operations
- ✅ Load tests for performance-critical components
- ✅ Network tests for distributed components
- ✅ HTML viewer supports filtering by all attributes
- ✅ Skills generate comprehensive test coverage
- ✅ Documentation updated with filtering guide
