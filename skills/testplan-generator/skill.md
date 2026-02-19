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

**Don't:**
- Skip environment variants to save time (they're critical)
- Ignore edge cases and boundary conditions
- Forget negative testing (invalid inputs, error paths)
- Overlook integration testing
- Create tests without clear validation criteria
- Generate tests for unsupported configurations
- Duplicate tests unnecessarily

## Success Metrics

A successful test plan generation includes:
- **70%+ coverage** of documented features
- **3+ environment variants** for each critical test
- **Negative tests** for all user inputs
- **Integration tests** for all external dependencies
- **Documented gaps** with recommendations
- **Executable tests** with clear validation criteria
- **Realistic time estimates** based on complexity
