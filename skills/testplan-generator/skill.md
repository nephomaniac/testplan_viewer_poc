# Test Plan Generator Skill

Generate comprehensive, well-structured test plans from any combination of artifacts including code repositories, existing tests, documentation, SOPs, and guides. Focus on complete coverage across different environments and configurations.

## ⚠️ Skill Self-Verification (REQUIRED - Run First)

**BEFORE executing this skill, ALWAYS verify you're using the latest skill definition.**

### Verification Steps

1. **Check if skill definition exists in repository:**
   ```bash
   ls -la skills/testplan-generator/skill.md
   ```

2. **Verify this is a git repository:**
   ```bash
   git rev-parse --is-inside-work-tree 2>/dev/null || echo "Not a git repo"
   ```

3. **Check if local skill has uncommitted changes:**
   ```bash
   git status skills/testplan-generator/skill.md
   ```

4. **Compare local vs committed version:**
   ```bash
   # Check if there are differences between working copy and HEAD
   git diff skills/testplan-generator/skill.md
   ```

5. **Check remote for updates (if applicable):**
   ```bash
   # Fetch latest from remote (don't merge)
   git fetch origin main 2>/dev/null

   # Compare local version with remote
   git diff HEAD origin/main -- skills/testplan-generator/skill.md
   ```

### Decision Tree

**If differences are found, PROMPT the user:**

```
⚠️ Skill Definition Verification

I've detected differences in the testplan-generator skill definition:

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
- **Option 2 (Use latest):** Read `skills/testplan-generator/skill.md` from disk and use that definition
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

Skill: testplan-generator
Version: [git commit hash of skills/testplan-generator/skill.md]
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

You are a test plan architect specializing in comprehensive test coverage. Your goal is to analyze provided artifacts and generate complete test plans that cover all aspects of the system under test, including:

- All functional paths (happy path, edge cases, error handling)
- Different environment configurations (cluster types, cloud providers, regions)
- Various input permutations
- Integration points and dependencies
- Performance and scalability scenarios
- Security and compliance requirements

## Inputs

You will receive one or more of the following artifacts:

### Primary Focus Artifact (Required)
Choose ONE as the starting point:
- **Code repository** - The codebase to be tested
- **Documentation** - User guides, API docs, architecture docs
- **Existing test suite** - Unit, integration, or E2E tests in any format

### Supporting Artifacts (Optional but Recommended)
- Git repositories with related code
- Existing test cases (JUnit XML, pytest, Go test, etc.)
- Product documentation and user guides
- SOPs (Standard Operating Procedures)
- Architecture diagrams
- API specifications (OpenAPI, gRPC proto)
- Configuration examples
- Known issues and bug reports

### Context Information
- **Target system**: What is being tested (e.g., "ROSA HCP Monitoring Operator")
- **Environment types**: Configurations to cover (e.g., "AWS/GCP/Azure", "single-AZ/multi-AZ")
- **User personas**: Who will execute tests (e.g., "SRE engineers", "QE team")
- **Coverage goals**: What aspects to emphasize (e.g., "security", "performance", "user workflows")

## Analysis Process

### Phase 1: Artifact Analysis (15-30 minutes)

**1. Explore Primary Artifact**
```bash
# For code repository
find . -name "*.go" -o -name "*.py" -o -name "*.js"
cat README.md docs/architecture.md
tree -L 3 -d

# For documentation
grep -r "install\|configure\|deploy" docs/
find . -name "*.md" | xargs cat

# For existing tests
find . -name "*test*.go" -o -name "test_*.py"
grep -r "func Test\|def test_" .
```

**2. Identify System Components**
- Main features and capabilities
- External dependencies
- Configuration options
- API endpoints or CLI commands
- Data models and persistence
- Integration points

**3. Map Environment Variations**
List all environment dimensions:
- Cloud providers (AWS, GCP, Azure, bare metal)
- Cluster configurations (single-zone, multi-zone, HCP, classic)
- Deployment methods (Helm, Operator, manual)
- Scale variations (small, medium, large clusters)
- Network configurations (public, private, proxy)
- Security contexts (FIPS, restricted SCC)

### Phase 2: Coverage Matrix Design (20-40 minutes)

**Build a coverage matrix** considering:

**Functional Coverage:**
- Installation/Setup
- Configuration
- Core functionality
- Advanced features
- Error handling
- Upgrade/Migration
- Uninstall/Cleanup

**Environmental Coverage:**
| Test Case | AWS | GCP | Azure | Single-AZ | Multi-AZ | FIPS |
|-----------|-----|-----|-------|-----------|----------|------|
| Install   | ✓   | ✓   | ✓     | ✓         | ✓        | ✓    |
| Config    | ✓   | ✓   | ✓     | ✓         | ✓        | ✓    |
| Monitor   | ✓   | -   | -     | ✓         | ✓        | -    |

**Input Permutations:**
- Valid inputs (happy path)
- Boundary conditions (min/max values)
- Invalid inputs (negative testing)
- Edge cases (empty, null, special characters)

**Integration Coverage:**
- Upstream dependencies
- Downstream consumers
- Third-party services
- Platform APIs

### Phase 3: Test Case Generation (30-60 minutes per test)

For each identified test scenario, generate a complete test case following the template structure:

**Test Case Structure:**
```json
{
  "test_id": {
    "metadata": {
      "id": "test_install_aws_multi_az",
      "title": "Install operator on AWS multi-AZ cluster",
      "category": "installation",
      "difficulty": "intermediate",
      "estimated_time": "30 minutes",
      "environments": ["aws", "multi-az"],
      "coverage_areas": ["installation", "networking", "high-availability"]
    },
    "test_execution": {
      "objective": "Verify operator installs correctly on multi-AZ AWS cluster",
      "prerequisites": ["AWS credentials", "Multi-AZ ROSA cluster"],
      "environment_setup": {
        "cloud": "AWS",
        "cluster_type": "ROSA HCP",
        "availability_zones": ["us-east-1a", "us-east-1b", "us-east-1c"],
        "cluster_size": "medium"
      },
      "steps": [...]
    }
  }
}
```

**Iteration Strategy:**
1. Start with core happy-path test
2. Add environment variations (create variant tests for each cloud/config)
3. Add negative test cases
4. Add boundary condition tests
5. Add integration tests

### Phase 4: Environment Variant Generation (20-30 minutes)

For each test case, automatically generate variants for different environments:

**Base test:** `test_install_operator`

**Generated variants:**
- `test_install_operator_aws_single_az`
- `test_install_operator_aws_multi_az`
- `test_install_operator_gcp_single_region`
- `test_install_operator_azure_availability_zones`
- `test_install_operator_fips_enabled`
- `test_install_operator_restricted_network`

**Variant template:**
```json
{
  "test_install_operator_aws_multi_az": {
    "metadata": {
      "id": "test_install_operator_aws_multi_az",
      "title": "Install operator on AWS multi-AZ cluster",
      "parent_test": "test_install_operator",
      "environment_variant": {
        "cloud": "AWS",
        "zones": "multi-az"
      }
    },
    "environment_setup": {
      "cloud": "AWS",
      "cluster_type": "ROSA HCP",
      "availability_zones": 3,
      "specific_setup": [
        "Verify cluster spans multiple AZs",
        "Check load balancer configuration"
      ]
    }
  }
}
```

### Phase 5: Gap Analysis (10-20 minutes)

Identify and document gaps in coverage:

**Coverage report:**
```json
{
  "coverage_summary": {
    "total_test_cases": 24,
    "environments_covered": {
      "aws": {"tests": 15, "coverage": "high"},
      "gcp": {"tests": 6, "coverage": "medium"},
      "azure": {"tests": 3, "coverage": "low"}
    },
    "functional_coverage": {
      "installation": "complete",
      "configuration": "complete",
      "monitoring": "partial",
      "troubleshooting": "minimal",
      "upgrade": "not_covered"
    },
    "identified_gaps": [
      "No tests for GCP private clusters",
      "Missing upgrade path testing",
      "No performance/scale tests",
      "Limited security testing"
    ],
    "recommendations": [
      "Add GCP private cluster variant tests",
      "Create upgrade test suite",
      "Add load testing scenarios"
    ]
  }
}
```

### Phase 6: System Impact and Safety Analysis (CRITICAL - 30-45 minutes)

**REQUIRED:** Analyze every test case to determine its impact on the system under test and ensure safe execution with proper state management.

#### Step 1: Classify System Impact

For EACH test case, analyze the commands/operations to classify impact type:

**Read-Only Tests** - Makes no changes to system state:
```bash
# Examples of read-only operations
oc get pods
kubectl describe deployment
curl http://service/metrics
promtool query instant 'up'
oc logs deployment/app
```

**Characteristics:**
- Only queries, reads, describes, or observes
- No create, apply, patch, scale, delete operations
- No modification flags (--patch, --replicas, etc.)
- Safe to run anywhere, including production

**Modifies-State Tests** - Changes system state, but reversible:
```bash
# Examples of modifying operations
oc apply -f deployment.yaml
oc scale deployment/app --replicas=3
oc patch configmap/config --type=merge
kubectl create namespace test
helm install myapp ./chart
```

**Characteristics:**
- Creates, updates, or scales resources
- Changes persist after test completes
- Reversible with cleanup/deletion
- Requires cleanup test to restore state

**Destructive Tests** - Irreversible or high-risk changes:
```bash
# Examples of destructive operations
oc delete deployment production-app
oc delete pvc data-volume
kubectl drain node-1 --force --delete-emptydir-data
rosa delete cluster --cluster=prod
oc delete namespace production
```

**Characteristics:**
- Deletes resources permanently
- Uses --force flags
- Cannot be undone without backup
- Requires BOTH backup AND cleanup tests

#### Step 2: Create Backup/Cleanup Pattern

For every **modifies-state** or **destructive** test, create a 4-test pattern:

**Pattern Example: Modifying Cluster Configuration**

**Test 1: Backup Test (read-only)**
```json
{
  "test_backup_cluster_config": {
    "metadata": {
      "id": "test_backup_cluster_config",
      "category": "backup",
      "system_impact": {
        "type": "read-only",
        "description": "Exports cluster config to backup file",
        "affected_resources": [],
        "reversible": true,
        "persistence": "temporary",
        "risk_level": "low"
      },
      "state_management": {
        "requires_backup": false,
        "requires_cleanup": false,
        "restores_state": false
      },
      "safety": {
        "can_run_in_production": true,
        "requires_confirmation": false,
        "safe_to_retry": true,
        "idempotent": true
      }
    },
    "test_execution": {
      "steps": [
        {
          "command": "oc get cm cluster-config -o yaml > backup/config-$(date +%Y%m%d-%H%M%S).yaml"
        }
      ]
    }
  }
}
```

**Test 2: Modify Test (modifies-state)**
```json
{
  "test_modify_cluster_config": {
    "metadata": {
      "id": "test_modify_cluster_config",
      "category": "configuration",
      "system_impact": {
        "type": "modifies-state",
        "description": "Modifies cluster configuration",
        "affected_resources": ["ConfigMap: cluster-config"],
        "reversible": true,
        "persistence": "permanent",
        "risk_level": "medium"
      },
      "state_management": {
        "requires_backup": true,
        "backup_test": "test_backup_cluster_config",
        "requires_cleanup": true,
        "cleanup_test": "test_restore_cluster_config"
      },
      "safety": {
        "can_run_in_production": false,
        "requires_confirmation": true,
        "warning_message": "⚠️ This test modifies cluster configuration. Backup will run first. Cleanup required after.",
        "safe_to_retry": true,
        "idempotent": true
      }
    },
    "test_execution": {
      "dependencies": ["test_backup_cluster_config"],
      "steps": [
        {
          "command": "oc patch cm cluster-config --type=merge -p '{\"data\":{\"new-key\":\"value\"}}'"
        }
      ]
    }
  }
}
```

**Test 3: Cleanup/Restore Test (modifies-state)**
```json
{
  "test_restore_cluster_config": {
    "metadata": {
      "id": "test_restore_cluster_config",
      "category": "cleanup",
      "system_impact": {
        "type": "modifies-state",
        "description": "Restores cluster config from backup",
        "affected_resources": ["ConfigMap: cluster-config"],
        "reversible": false,
        "persistence": "permanent",
        "risk_level": "low"
      },
      "state_management": {
        "requires_backup": false,
        "requires_cleanup": false,
        "restores_state": true,
        "state_verification": "test_verify_config_restored"
      },
      "safety": {
        "can_run_in_production": false,
        "requires_confirmation": false,
        "safe_to_retry": true,
        "idempotent": true
      }
    },
    "test_execution": {
      "dependencies": ["test_backup_cluster_config"],
      "steps": [
        {
          "command": "oc apply -f backup/config-*.yaml"
        }
      ]
    }
  }
}
```

**Test 4: Verification Test (read-only)**
```json
{
  "test_verify_config_restored": {
    "metadata": {
      "id": "test_verify_config_restored",
      "category": "validation",
      "system_impact": {
        "type": "read-only",
        "description": "Verifies config matches original backup",
        "affected_resources": [],
        "reversible": true,
        "persistence": "temporary",
        "risk_level": "low"
      },
      "state_management": {
        "requires_backup": false,
        "requires_cleanup": false,
        "restores_state": false
      },
      "safety": {
        "can_run_in_production": true,
        "requires_confirmation": false,
        "safe_to_retry": true,
        "idempotent": true
      }
    },
    "test_execution": {
      "dependencies": ["test_restore_cluster_config"],
      "steps": [
        {
          "command": "diff <(oc get cm cluster-config -o yaml) backup/config-*.yaml"
        }
      ]
    }
  }
}
```

#### Step 3: Assign Safety Attributes

For EVERY test, include complete safety metadata:

**Safety Decision Matrix:**

| Impact Type | Backup Required | Cleanup Required | Production Safe | Confirmation | Risk Level |
|-------------|-----------------|------------------|-----------------|--------------|------------|
| read-only | No | No | Yes | No | low |
| modifies-state (reversible) | No | Yes | No | Yes | medium |
| modifies-state (temp) | No | Yes | No | Yes | low-medium |
| destructive | Yes | Yes | Never | Yes | high-critical |

**Field Guidance:**

- **system_impact.type**: Auto-detect from commands (get/describe = read-only, apply/create = modifies-state, delete/drain = destructive)
- **system_impact.description**: Human-readable explanation of what changes
- **system_impact.affected_resources**: List specific resources (e.g., "Deployment: my-app", "Namespace: test-ns")
- **system_impact.reversible**: Can it be undone? (delete operations = false)
- **system_impact.persistence**: "temporary" if auto-cleaned, "permanent" if persists
- **system_impact.risk_level**: Impact severity (low/medium/high/critical)
- **state_management.requires_backup**: true for destructive tests
- **state_management.backup_test**: test_id of backup test
- **state_management.requires_cleanup**: true for modifies-state and destructive
- **state_management.cleanup_test**: test_id of cleanup/restore test
- **state_management.restores_state**: true ONLY for cleanup tests themselves
- **state_management.state_verification**: Optional verification test_id
- **safety.can_run_in_production**: false for modifies-state and destructive
- **safety.requires_confirmation**: true for high-risk tests
- **safety.warning_message**: Clear warning for risky tests
- **safety.safe_to_retry**: Can test be retried if it fails?
- **safety.idempotent**: Running multiple times has same effect?

#### Step 4: Create Test Dependencies

Link tests in execution order:

**For modifying tests:**
```
1. backup_test (if destructive)
2. modify_test (depends on backup)
3. cleanup_test (user runs manually or automatically)
4. verify_test (depends on cleanup)
```

**Update test_execution.dependencies:**
```json
{
  "test_modify_cluster_config": {
    "test_execution": {
      "dependencies": ["test_backup_cluster_config"]
    }
  },
  "test_restore_cluster_config": {
    "test_execution": {
      "dependencies": ["test_backup_cluster_config"]
    }
  },
  "test_verify_config_restored": {
    "test_execution": {
      "dependencies": ["test_restore_cluster_config"]
    }
  }
}
```

#### Safety Analysis Checklist

Before finalizing test plan, verify EVERY test has:

- [ ] system_impact.type classified correctly (read-only/modifies-state/destructive)
- [ ] system_impact.affected_resources listed specifically
- [ ] system_impact.risk_level assigned appropriately
- [ ] state_management.requires_backup set to true for destructive tests
- [ ] state_management.backup_test links to actual backup test ID
- [ ] state_management.requires_cleanup set to true for modifying tests
- [ ] state_management.cleanup_test links to actual cleanup test ID
- [ ] safety.can_run_in_production set to false for unsafe tests
- [ ] safety.warning_message written for high-risk tests
- [ ] Backup test exists and is read-only
- [ ] Cleanup test exists and restores state
- [ ] Dependencies properly linked (backup → modify → cleanup → verify)
- [ ] Verification test validates restoration (optional but recommended)

**CRITICAL REQUIREMENT:** No modifies-state or destructive test should exist without its corresponding backup/cleanup tests. If you cannot create a safe cleanup test, mark the test as "manual-cleanup-required" and document the manual steps.

---

### Phase 7: Negative Testing Generation (30-45 minutes)

**REQUIRED:** For every positive/happy-path test, generate corresponding negative test cases to validate error handling and system resilience.

#### What is Negative Testing?

Negative testing validates that the system correctly handles invalid inputs, error conditions, and unexpected scenarios. These tests should **expect failures** and verify appropriate error messages and system behavior.

#### Generate Negative Tests For:

**1. Invalid Configuration**
For each configuration test, create a negative variant with invalid values:
```json
{
  "test_create_deployment_invalid_image": {
    "metadata": {
      "id": "test_create_deployment_invalid_image",
      "title": "Attempt deployment with non-existent image",
      "category": "negative",
      "difficulty": "beginner"
    },
    "test_execution": {
      "objective": "Validate error handling when invalid image is specified",
      "child_tests": [{
        "id": "test_create_deployment_invalid_image_001",
        "parent_test_id": "test_create_deployment_invalid_image",
        "step_number": 1,
        "title": "Apply deployment with invalid image reference",
        "test_type": "negative",
        "impact_type": "read-only",
        "input_validation": "negative",
        "rbac_level": "edit",
        "command": "oc apply -f deployment-invalid-image.yaml",
        "expected_output": "Error: ErrImagePull",
        "learning_note": "Validates that Kubernetes correctly rejects invalid image references and provides clear error messages",
        "why_this_step": "Invalid images are a common deployment error. This test ensures proper error detection.",
        "safety": {
          "can_run_in_production": true,
          "requires_cleanup": false,
          "risk_level": "low"
        },
        "validation": {
          "success_criteria": ["Error message contains 'ErrImagePull' or 'ImagePullBackOff'", "Pod enters error state"],
          "failure_criteria": ["Deployment succeeds", "No error message shown"],
          "metrics_to_collect": []
        }
      }]
    }
  }
}
```

**2. Permission Denials**
Test operations without required permissions:
```json
{
  "test_create_secret_insufficient_permissions": {
    "metadata": {
      "id": "test_create_secret_insufficient_permissions",
      "title": "Attempt secret creation with view-only permissions",
      "category": "negative"
    },
    "test_execution": {
      "child_tests": [{
        "id": "test_create_secret_insufficient_permissions_001",
        "parent_test_id": "test_create_secret_insufficient_permissions",
        "step_number": 1,
        "title": "View user attempts to create secret",
        "test_type": "negative",
        "impact_type": "read-only",
        "input_validation": "negative",
        "rbac_level": "view",
        "command": "oc create secret generic test --from-literal=key=value --as=view-user",
        "expected_output": "Error from server (Forbidden): secrets is forbidden",
        "learning_note": "Validates RBAC prevents view-only users from creating resources",
        "safety": {
          "can_run_in_production": true,
          "risk_level": "low"
        }
      }]
    }
  }
}
```

**3. Resource Conflicts**
Test creation of duplicate resources:
```json
{
  "child_tests": [{
    "title": "Attempt to create duplicate resource",
    "test_type": "negative",
    "impact_type": "read-only",
    "input_validation": "negative",
    "command": "oc create deployment duplicate-app --image=nginx",
    "expected_output": "Error from server (AlreadyExists)",
    "learning_note": "Validates proper error handling for duplicate resource creation"
  }]
}
```

**4. Missing Dependencies**
Test scenarios where required resources don't exist:
```json
{
  "child_tests": [{
    "title": "Create custom resource before CRD exists",
    "test_type": "negative",
    "impact_type": "read-only",
    "expected_output": "Error: the server doesn't have a resource type",
    "learning_note": "Validates dependency checking and clear error messages"
  }]
}
```

#### Negative Test Generation Strategy

For EACH positive test, generate 2-3 negative counterparts:

**Example:**
- Positive: `test_create_pagerduty_secret` (creates secret successfully)
- Negative 1: `test_create_pagerduty_secret_missing_key` (missing required key)
- Negative 2: `test_create_pagerduty_secret_invalid_namespace` (non-existent namespace)
- Negative 3: `test_create_pagerduty_secret_no_permissions` (RBAC denial)

#### Negative Test Requirements

All negative tests MUST:
- Set `test_type: "negative"`
- Set `impact_type: "read-only"` (should fail before making changes)
- Define `expected_output` with specific error message
- Include `learning_note` explaining why failure is expected
- Be production-safe (`can_run_in_production: true`)

---

### Phase 8: Input Validation Testing (30-45 minutes)

**REQUIRED:** Generate comprehensive input validation tests covering positive, negative, missing, corrupt, and boundary cases.

#### Input Validation Categories

**1. Positive Input Tests (`input_validation: "positive"`)**
Valid inputs expected to succeed:
```json
{
  "child_tests": [{
    "title": "Create secret with valid key format",
    "test_type": "validation",
    "impact_type": "modifies-state",
    "input_validation": "positive",
    "rbac_level": "edit",
    "command": "oc create secret generic pd-secret --from-literal=PAGERDUTY_KEY=abc123",
    "expected_output": "secret/pd-secret created",
    "learning_note": "Validates that properly formatted inputs are accepted",
    "safety": {
      "can_run_in_production": false,
      "requires_cleanup": true,
      "cleanup_procedure": "oc delete secret pd-secret"
    }
  }]
}
```

**2. Negative Input Tests (`input_validation: "negative"`)**
Invalid types, formats, or values:
```json
{
  "child_tests": [{
    "title": "Attempt secret creation with empty value",
    "test_type": "negative",
    "impact_type": "read-only",
    "input_validation": "negative",
    "command": "oc create secret generic pd-secret --from-literal=PAGERDUTY_KEY=",
    "expected_output": "Error: data[PAGERDUTY_KEY]: Invalid value",
    "learning_note": "Validates that empty secret values are rejected with clear error"
  }]
}
```

**3. Missing Attribute Tests (`input_validation: "missing"`)**
Required fields/attributes omitted:
```json
{
  "child_tests": [{
    "title": "Deploy without required image field",
    "test_type": "negative",
    "impact_type": "read-only",
    "input_validation": "missing",
    "command": "oc apply -f deployment-no-image.yaml",
    "expected_output": "Error: spec.template.spec.containers[0].image: Required value",
    "learning_note": "Validates that missing required fields are caught with helpful error messages"
  }]
}
```

**4. Corrupt Data Tests (`input_validation: "corrupt"`)**
Malformed YAML/JSON, invalid structure:
```json
{
  "child_tests": [{
    "title": "Apply malformed YAML configuration",
    "test_type": "negative",
    "impact_type": "read-only",
    "input_validation": "corrupt",
    "command": "oc apply -f corrupt-syntax.yaml",
    "expected_output": "error: error parsing corrupt-syntax.yaml",
    "learning_note": "Validates that malformed YAML is rejected with parse errors"
  }]
}
```

**5. Boundary Value Tests (`input_validation: "boundary"`)**
Min/max values, limits:
```json
{
  "child_tests": [{
    "title": "Create deployment with 0 replicas (minimum boundary)",
    "test_type": "validation",
    "impact_type": "modifies-state",
    "input_validation": "boundary",
    "command": "oc create deployment test --image=nginx --replicas=0",
    "learning_note": "Tests minimum boundary - 0 replicas should be accepted"
  }, {
    "title": "Create deployment with 10000 replicas (maximum boundary)",
    "test_type": "validation",
    "impact_type": "read-only",
    "input_validation": "boundary",
    "command": "oc create deployment test --image=nginx --replicas=10000",
    "expected_output": "Error: replicas exceeds maximum",
    "learning_note": "Tests maximum boundary - excessively high replicas should be rejected"
  }]
}
```

#### Input Validation Matrix

For EACH input field/parameter, generate tests covering:

| Input | Positive | Negative | Missing | Corrupt | Boundary |
|-------|----------|----------|---------|---------|----------|
| image | ✓ Valid ref | ✓ Non-existent | ✓ No field | ✓ Malformed | N/A |
| replicas | ✓ Valid num | ✓ Negative | ✓ No replicas | ✓ String | ✓ 0, 10000 |
| namespace | ✓ Exists | ✓ Non-existent | ✓ No ns | ✓ Invalid chars | ✓ 63-char |

---

### Phase 9: RBAC and Permission Testing (30-45 minutes)

**REQUIRED:** Generate tests to validate role-based access control and security boundaries.

#### RBAC Test Generation Rules

For EACH operation, determine minimum required RBAC level and generate:
1. Positive test with correct permissions
2. Negative test with insufficient permissions

#### RBAC Level Categories

**1. Cluster-Admin Operations (`rbac_level: "cluster-admin"`)**
Cluster-scoped resources (CRDs, ClusterRoles, Nodes):
```json
{
  "child_tests": [{
    "title": "Create CustomResourceDefinition",
    "test_type": "rbac",
    "impact_type": "modifies-state",
    "rbac_level": "cluster-admin",
    "command": "oc apply -f mycrd.yaml",
    "expected_output": "customresourcedefinition.apiextensions.k8s.io/mycrds.example.com created",
    "learning_note": "CRD creation requires cluster-admin permissions as it's cluster-scoped",
    "safety": {
      "can_run_in_production": false,
      "requires_cleanup": true,
      "cleanup_procedure": "oc delete crd mycrds.example.com",
      "risk_level": "high"
    }
  }]
}
```

**2. Namespace-Admin Operations (`rbac_level: "namespace-admin"`)**
Namespace-scoped admin operations (Roles, RoleBindings):
```json
{
  "child_tests": [{
    "title": "Create Role in namespace",
    "test_type": "rbac",
    "impact_type": "modifies-state",
    "rbac_level": "namespace-admin",
    "command": "oc create role pod-reader --verb=get,list --resource=pods -n test-ns",
    "safety": {
      "can_run_in_production": false,
      "requires_cleanup": true
    }
  }]
}
```

**3. Edit Permissions (`rbac_level: "edit"`)**
Standard resource creation/modification:
```json
{
  "child_tests": [{
    "title": "Create deployment",
    "test_type": "install",
    "impact_type": "modifies-state",
    "rbac_level": "edit",
    "command": "oc create deployment nginx --image=nginx -n test-ns"
  }]
}
```

**4. View-Only Operations (`rbac_level: "view"`)**
Read-only operations:
```json
{
  "child_tests": [{
    "title": "List pods in namespace",
    "test_type": "validation",
    "impact_type": "read-only",
    "rbac_level": "view",
    "command": "oc get pods -n test-ns",
    "safety": {
      "can_run_in_production": true,
      "risk_level": "low"
    }
  }]
}
```

**5. Permission Denial Tests**
Validate RBAC enforcement:
```json
{
  "test_create_deployment_as_viewer": {
    "metadata": {
      "title": "Attempt deployment creation as view-only user",
      "category": "negative"
    },
    "test_execution": {
      "child_tests": [{
        "title": "View user attempts deployment creation",
        "test_type": "negative",
        "impact_type": "read-only",
        "rbac_level": "view",
        "command": "oc create deployment nginx --image=nginx --as=view-user",
        "expected_output": "Error from server (Forbidden)",
        "learning_note": "Validates RBAC prevents view users from creating resources"
      }]
    }
  }
}
```

#### RBAC Test Requirements

All RBAC tests MUST:
- Set `rbac_level` to minimum required permission
- For negative tests: use lower permission level + expect Forbidden error
- Document why specific permission level is needed
- Test service account permissions if operator uses them

---

### Phase 10: Load and Performance Testing (45-60 minutes)

**REQUIRED:** Generate tests to validate system behavior under load and resource constraints.

#### Load Testing Categories

**1. Resource Limit Tests (`test_type: "load"`)**
Test behavior when hitting resource limits:
```json
{
  "test_pod_with_low_memory": {
    "metadata": {
      "title": "Deploy pod with insufficient memory limits",
      "category": "load"
    },
    "test_execution": {
      "child_tests": [{
        "title": "Create pod with memory limit too low for application",
        "test_type": "load",
        "impact_type": "modifies-state",
        "resource_requirements": {
          "cpu": "100m",
          "memory": "50Mi"
        },
        "command": "oc apply -f pod-low-memory.yaml",
        "learning_note": "Tests pod behavior when memory limit is insufficient - pod should be OOMKilled",
        "safety": {
          "can_run_in_production": false,
          "requires_cleanup": true,
          "cleanup_procedure": "oc delete pod low-memory-pod",
          "risk_level": "medium"
        }
      }]
    }
  }
}
```

**2. Load Generation Tests (`impact_type: "load-generation"`)**
Deploy pods that generate load:
```json
{
  "test_operator_under_cpu_stress": {
    "metadata": {
      "title": "Test operator reconciliation under CPU stress",
      "category": "load"
    },
    "test_execution": {
      "child_tests": [{
        "id": "cpu_stress_deploy",
        "title": "Deploy CPU stress pods",
        "test_type": "load",
        "impact_type": "load-generation",
        "resource_requirements": {
          "cpu": "2",
          "memory": "512Mi",
          "load_generation": true
        },
        "command": "oc apply -f cpu-stress-deployment.yaml",
        "learning_note": "Deploys pods using stress-ng to consume CPU. Tests operator behavior under resource pressure.",
        "why_this_step": "Operators should continue functioning even when cluster resources are constrained",
        "safety": {
          "can_run_in_production": false,
          "requires_cleanup": true,
          "cleanup_procedure": "oc delete -f cpu-stress-deployment.yaml",
          "risk_level": "medium"
        }
      }, {
        "id": "cpu_stress_verify",
        "title": "Verify operator reconciles under load",
        "test_type": "validation",
        "impact_type": "read-only",
        "command": "oc get customresource -w",
        "duration": "5 minutes",
        "learning_note": "Validates operator continues processing reconciliation loops despite CPU stress"
      }, {
        "id": "cpu_stress_cleanup",
        "title": "Remove stress pods",
        "test_type": "cleanup",
        "impact_type": "modifies-state",
        "command": "oc delete -f cpu-stress-deployment.yaml"
      }]
    }
  }
}
```

**3. Performance Testing (`test_type: "performance"`)**
Measure and validate performance metrics:
```json
{
  "child_tests": [{
    "title": "Measure API response time",
    "test_type": "performance",
    "impact_type": "read-only",
    "command": "time curl -s http://api-service/health",
    "learning_note": "Measures API response latency to establish performance baseline",
    "validation": {
      "success_criteria": ["Response time < 100ms", "Status code 200"],
      "metrics_to_collect": ["response_time_ms", "status_code"]
    }
  }]
}
```

#### Load Test Deployment Examples

**CPU Stress:**
```yaml
# cpu-stress-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cpu-stress
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: stress
        image: polinux/stress
        args: ["--cpu", "2", "--timeout", "300s"]
        resources:
          requests:
            cpu: "2"
            memory: "512Mi"
```

**Memory Stress:**
```yaml
# memory-stress-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: memory-stress
spec:
  replicas: 2
  template:
    spec:
      containers:
      - name: stress
        image: polinux/stress
        args: ["--vm", "2", "--vm-bytes", "1G", "--timeout", "300s"]
        resources:
          requests:
            memory: "2Gi"
```

---

### Phase 11: Network Testing with Impairment (45-60 minutes)

**REQUIRED:** Generate tests to validate system behavior under network stress and impairment.

#### Network Testing Categories

**1. Network Policy Tests (`test_type: "network", `impairment_type: "network-policy"`)**
Restrict traffic using Kubernetes NetworkPolicy:
```json
{
  "test_app_with_deny_all_policy": {
    "metadata": {
      "title": "Test application with deny-all network policy",
      "category": "network"
    },
    "test_execution": {
      "child_tests": [{
        "id": "network_apply_deny",
        "title": "Apply deny-all network policy",
        "test_type": "network",
        "impact_type": "impairment",
        "network_requirements": {
          "impairment_needed": true,
          "impairment_type": "network-policy",
          "impairment_config": "Deny all ingress and egress traffic",
          "network_policies": ["deny-all-policy.yaml"]
        },
        "command": "oc apply -f deny-all-network-policy.yaml",
        "learning_note": "Applies NetworkPolicy denying all traffic to test isolation behavior",
        "why_this_step": "Network policies can accidentally block critical traffic - tests resilience",
        "safety": {
          "can_run_in_production": false,
          "requires_cleanup": true,
          "cleanup_procedure": "oc delete networkpolicy deny-all",
          "risk_level": "high"
        }
      }, {
        "id": "network_verify_blocked",
        "title": "Verify traffic is blocked",
        "test_type": "validation",
        "impact_type": "read-only",
        "command": "oc exec test-pod -- curl --max-time 5 http://service",
        "expected_output": "curl: (28) Connection timed out",
        "learning_note": "Confirms network policy is enforced - connection should timeout"
      }, {
        "id": "network_cleanup",
        "title": "Remove network policy",
        "test_type": "cleanup",
        "impact_type": "modifies-state",
        "command": "oc delete networkpolicy deny-all",
        "learning_note": "Restores normal network connectivity"
      }]
    }
  }
}
```

**2. Latency Injection Tests (`impairment_type: "latency"`)**
Add artificial network latency:
```json
{
  "test_api_with_latency": {
    "metadata": {
      "title": "Test API with 100ms network latency",
      "category": "network"
    },
    "test_execution": {
      "child_tests": [{
        "title": "Inject 100ms latency",
        "test_type": "network",
        "impact_type": "impairment",
        "network_requirements": {
          "impairment_needed": true,
          "impairment_type": "latency",
          "impairment_config": "100ms delay on eth0"
        },
        "manual_steps": [
          "Deploy privileged debug pod with tc utilities",
          "Run: tc qdisc add dev eth0 root netem delay 100ms",
          "Verify: ping shows ~100ms increase in RTT"
        ],
        "learning_note": "Uses Linux tc (traffic control) to add network latency",
        "safety": {
          "can_run_in_production": false,
          "requires_cleanup": true,
          "cleanup_procedure": "tc qdisc del dev eth0 root",
          "risk_level": "critical"
        }
      }, {
        "title": "Verify application handles latency",
        "test_type": "validation",
        "impact_type": "read-only",
        "command": "curl http://api-service/health",
        "learning_note": "Tests if application remains functional with degraded network"
      }]
    }
  }
}
```

**3. Packet Loss Tests (`impairment_type: "packet-loss"`)**
Simulate packet loss:
```json
{
  "network_requirements": {
    "impairment_needed": true,
    "impairment_type": "packet-loss",
    "impairment_config": "10% packet loss"
  },
  "manual_steps": [
    "tc qdisc add dev eth0 root netem loss 10%",
    "Verify with: ping -c 100 <target> (should show ~10% loss)"
  ]
}
```

**4. AWS Security Group Tests (platform-specific)**
Temporarily block traffic via security groups:
```json
{
  "child_tests": [{
    "title": "Block ingress via security group",
    "test_type": "network",
    "impact_type": "impairment",
    "network_requirements": {
      "impairment_needed": true,
      "impairment_type": "security-group-deny",
      "security_groups": ["sg-abc123"]
    },
    "command": "aws ec2 revoke-security-group-ingress --group-id sg-abc123 --protocol tcp --port 443 --cidr 0.0.0.0/0",
    "learning_note": "Removes security group rule to simulate network partition",
    "safety": {
      "can_run_in_production": false,
      "requires_cleanup": true,
      "cleanup_procedure": "aws ec2 authorize-security-group-ingress --group-id sg-abc123 --protocol tcp --port 443 --cidr 0.0.0.0/0",
      "risk_level": "critical"
    }
  }]
}
```

#### Network Policy Example

**Deny-All NetworkPolicy:**
```yaml
# deny-all-network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all
  namespace: test-namespace
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
```

#### Network Impairment Commands

**Latency:**
```bash
# Add 100ms latency
tc qdisc add dev eth0 root netem delay 100ms

# Add 50ms ± 10ms variance
tc qdisc add dev eth0 root netem delay 50ms 10ms

# Remove
tc qdisc del dev eth0 root
```

**Packet Loss:**
```bash
# Add 10% packet loss
tc qdisc add dev eth0 root netem loss 10%

# Remove
tc qdisc del dev eth0 root
```

#### Network Testing Safety Requirements

ALL network impairment tests MUST:
- Set `can_run_in_production: false`
- Set `risk_level: "high"` or `"critical"`
- Define `cleanup_procedure` in detail
- Test cleanup procedure before marking complete
- Document affected services/pods to limit blast radius
- Include verification step after cleanup

---

### Phase 12: Test Generation Summary

After generating all test types, create a summary showing coverage:

```json
{
  "test_generation_summary": {
    "total_tests_generated": 87,
    "by_test_type": {
      "validation": 25,
      "negative": 18,
      "install": 12,
      "rbac": 10,
      "load": 8,
      "network": 6,
      "cleanup": 8
    },
    "by_impact_type": {
      "read-only": 45,
      "modifies-state": 32,
      "destructive": 2,
      "impairment": 6,
      "load-generation": 2
    },
    "by_rbac_level": {
      "view": 30,
      "edit": 25,
      "namespace-admin": 15,
      "cluster-admin": 10,
      "custom": 7
    },
    "input_validation_coverage": {
      "positive": 40,
      "negative": 18,
      "missing": 12,
      "corrupt": 8,
      "boundary": 9
    },
    "production_safe_tests": 45,
    "test_environment_only": 42,
    "cleanup_coverage": "100%"
  }
}
```

---

## Output Format

Generate test plan JSON following the schema in `.claude/templates/testplan-template.json` with these key sections:

### 1. Metadata with Coverage Info
```json
{
  "metadata": {
    "document_title": "Test Plan: [System Name]",
    "coverage_strategy": "Comprehensive coverage across AWS, GCP, Azure with focus on multi-AZ deployments",
    "environments_tested": ["aws-single-az", "aws-multi-az", "gcp", "azure"],
    "test_count_by_category": {
      "installation": 8,
      "configuration": 12,
      "validation": 15,
      "troubleshooting": 6
    }
  }
}
```

### 2. Environment Configuration Matrix
```json
{
  "test_environments": {
    "aws_single_az": {
      "description": "AWS ROSA HCP single availability zone",
      "setup_requirements": ["AWS credentials", "Single-AZ cluster"],
      "applicable_tests": ["test_1", "test_2", "test_5"]
    },
    "aws_multi_az": {
      "description": "AWS ROSA HCP multiple availability zones",
      "setup_requirements": ["AWS credentials", "Multi-AZ cluster"],
      "applicable_tests": ["test_1", "test_2", "test_3", "test_5", "test_6"]
    }
  }
}
```

### 3. Test Cases with Variants
Each test includes environment-specific considerations:
```json
{
  "testcases": {
    "test_install_operator": {
      "metadata": {
        "environment_variants": ["aws", "gcp", "azure"],
        "variant_considerations": {
          "aws": "Verify IAM roles created correctly",
          "gcp": "Check service account permissions",
          "azure": "Validate managed identity setup"
        }
      }
    }
  }
}
```

### 4. Coverage Summary Section
```json
{
  "coverage_summary": {
    "total_tests": 41,
    "environment_coverage": {
      "aws": 20,
      "gcp": 12,
      "azure": 9
    },
    "functional_coverage": {
      "installation": "complete",
      "configuration": "complete",
      "operations": "partial",
      "troubleshooting": "basic"
    },
    "input_permutations_covered": {
      "valid_inputs": 25,
      "invalid_inputs": 10,
      "boundary_conditions": 6
    },
    "gaps_identified": [
      "No ARM architecture testing",
      "Limited IPv6 coverage",
      "Missing disaster recovery scenarios"
    ]
  }
}
```

## Iteration Guidelines

### First Pass: Core Coverage
- Essential happy path tests
- Basic environment variants (1-2 per cloud)
- Critical failure scenarios

### Second Pass: Comprehensive Coverage
- All environment combinations
- Edge cases and boundary conditions
- Integration test scenarios
- Performance considerations

### Third Pass: Advanced Scenarios
- Complex failure modes
- Upgrade paths
- Security hardening tests
- Compliance validation

## Quality Checklist

Before finalizing test plan:

**Coverage Completeness:**
- [ ] All major features have test coverage
- [ ] Each supported environment has representative tests
- [ ] Both positive and negative test cases included
- [ ] Integration points tested
- [ ] Configuration options validated

**Comprehensive Testing (NEW):**
- [ ] Negative tests generated for each positive test (30-40% ratio)
- [ ] Input validation tests cover: positive, negative, missing, corrupt, boundary
- [ ] RBAC tests validate all permission boundaries
- [ ] Load tests for performance-critical components
- [ ] Network tests for distributed systems
- [ ] All tests use child_tests format (not legacy steps)

**Test Type Distribution:**
- [ ] Validation tests (read-only verification)
- [ ] Negative tests (error handling)
- [ ] RBAC tests (permission boundaries)
- [ ] Load tests (resource constraints)
- [ ] Network tests (impairment scenarios)
- [ ] Setup/Install/Cleanup tests (lifecycle)

**Filterable Attributes:**
- [ ] Every child test has test_type assigned
- [ ] Every child test has impact_type assigned
- [ ] RBAC-sensitive tests have rbac_level assigned
- [ ] Input validation tests have input_validation assigned
- [ ] Load tests have resource_requirements defined
- [ ] Network tests have network_requirements defined

**Environment Variants:**
- [ ] Cloud provider variations documented
- [ ] Cluster configuration variants included
- [ ] Network configuration scenarios covered
- [ ] Security context variations tested

**Input Permutations:**
- [ ] Valid input happy paths tested
- [ ] Invalid input handling verified
- [ ] Boundary conditions identified
- [ ] Edge cases documented

**Gap Analysis:**
- [ ] Coverage gaps identified and documented
- [ ] Recommendations for additional tests provided
- [ ] Priority assigned to missing coverage areas

**Test Quality:**
- [ ] Each test has clear objective
- [ ] Steps are specific and executable
- [ ] Expected outputs defined
- [ ] Validation criteria measurable

## Example Usage

**User provides:**
```
Focus artifact: observability-operator repository
Supporting artifacts:
  - Existing E2E tests in e2e/tests/
  - OpenShift monitoring documentation
  - SOPs for troubleshooting alerts

Environment types:
  - AWS ROSA (single-AZ and multi-AZ)
  - GCP (regional clusters)
  - Azure ARO

Coverage goals:
  - Complete installation coverage
  - Monitoring configuration
  - Alert validation
  - Troubleshooting workflows
```

**You should:**
1. Analyze repository structure
2. Read E2E tests to understand existing coverage
3. Review documentation for feature list
4. Examine SOPs for operational scenarios
5. Create coverage matrix for AWS/GCP/Azure variants
6. Generate 15-20 core test cases
7. Create environment-specific variants (45-60 total tests)
8. Document coverage gaps
9. Output complete test plan JSON

## Tips for Comprehensive Coverage

**Do:**
- Start with a coverage matrix before writing tests
- Generate variants systematically for each environment
- Use existing tests as inspiration for manual tests
- Document what's NOT covered (gaps are valuable info)
- Include both functional and non-functional tests
- Consider failure scenarios, not just happy paths
- Think about different user personas and their workflows
- Include performance and scale considerations
- **Use child_tests format** for all new tests (not legacy steps)
- **Generate negative tests** for every positive test (30-40% ratio)
- **Validate all inputs** with positive/negative/missing/corrupt/boundary tests
- **Test RBAC boundaries** with both allowed and denied operations
- **Include load tests** for performance-critical components
- **Test network resilience** with impairment scenarios

**Don't:**
- Skip environment variants to save time (they're critical)
- Ignore edge cases and boundary conditions
- Forget negative testing (invalid inputs, error paths)
- Overlook integration testing
- Create tests without clear validation criteria
- Generate tests for unsupported configurations
- Duplicate tests unnecessarily
- **Use legacy steps format** (use child_tests instead)
- **Skip negative tests** (they're required, not optional)
- **Forget cleanup tests** for modifying operations
- **Omit filterable attributes** (test_type, impact_type, rbac_level, etc.)
- **Create load/network tests** without cleanup procedures

## Success Metrics

A successful test plan generation includes:

**Coverage Metrics:**
- **70%+ coverage** of documented features
- **3+ environment variants** for each critical test
- **30-40% negative tests** relative to positive tests
- **100% input validation coverage** for all inputs
- **Integration tests** for all external dependencies
- **Documented gaps** with recommendations
- **Executable tests** with clear validation criteria
- **Realistic time estimates** based on complexity

**Test Type Distribution (Target for 50 tests):**
- 20 validation tests (40%)
- 10 negative tests (20%)
- 8 input validation tests (16%)
- 5 RBAC tests (10%)
- 4 load tests (8%)
- 3 network tests (6%)
- 5 cleanup tests (10%)

**Quality Indicators:**
- **All modifying tests** have cleanup tests
- **All negative tests** define expected error messages
- **All RBAC tests** specify minimum required permission level
- **All load tests** define resource requirements
- **All network tests** define impairment type and cleanup procedure
- **100% production safety** classification (can_run_in_production flag set)

**New Format Adoption:**
- **100% of tests** use child_tests format (not legacy steps)
- **Every child test** has filterable attributes (test_type, impact_type, etc.)
- **All resource-intensive tests** have resource_requirements defined
- **All network tests** have network_requirements defined
