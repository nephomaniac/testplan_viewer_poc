# Test Safety and State Management

Comprehensive guide for managing system state during test execution, ensuring tests can be safely run and reversed.

## Overview

Tests fall into three categories based on their impact on the system under test:

1. **Read-only** - No changes made, safe to run anywhere
2. **Modifies-state** - Makes changes that persist, requires cleanup
3. **Destructive** - Irreversible or high-risk changes, requires backup

## Test Case Attributes

### system_impact

Describes what changes the test makes to the system.

```json
{
  "system_impact": {
    "type": "read-only|modifies-state|destructive",
    "description": "Explanation of what changes this test makes",
    "affected_resources": [
      "List of resources that will be modified"
    ],
    "reversible": true,
    "persistence": "temporary|permanent",
    "risk_level": "low|medium|high|critical"
  }
}
```

**Fields:**
- `type` - Category of system impact
- `description` - Human-readable explanation
- `affected_resources` - Specific resources modified (deployments, configs, data, etc.)
- `reversible` - Can changes be undone?
- `persistence` - Do changes last beyond test execution?
- `risk_level` - Severity if test fails or corrupts state

### state_management

Defines backup and cleanup requirements.

```json
{
  "state_management": {
    "requires_backup": false,
    "backup_test": "test_backup_id",
    "requires_cleanup": false,
    "cleanup_test": "test_cleanup_id",
    "restores_state": false,
    "state_verification": "test_verify_id"
  }
}
```

**Fields:**
- `requires_backup` - Must backup before running?
- `backup_test` - Test ID that performs backup (if requires_backup=true)
- `requires_cleanup` - Must cleanup after running?
- `cleanup_test` - Test ID that performs cleanup/restoration
- `restores_state` - Does THIS test restore state? (for cleanup tests)
- `state_verification` - Optional test to verify state was restored

### safety

Safety constraints and execution requirements.

```json
{
  "safety": {
    "can_run_in_production": false,
    "requires_confirmation": true,
    "warning_message": "Warning to display before running",
    "safe_to_retry": true,
    "idempotent": false
  }
}
```

**Fields:**
- `can_run_in_production` - Safe for production environments?
- `requires_confirmation` - Prompt user before running?
- `warning_message` - Message to show if requires_confirmation=true
- `safe_to_retry` - Can be retried if fails?
- `idempotent` - Running multiple times has same effect?

## Impact Types

### Read-Only Tests

**Characteristics:**
- No changes to system state
- Only reads/observes existing state
- Safe to run in production
- No cleanup required

**Example commands:**
```bash
oc get pods
kubectl describe deployment
curl http://service/metrics
promtool query instant 'up'
```

**Template:**
```json
{
  "system_impact": {
    "type": "read-only",
    "description": "Only reads cluster information, makes no changes",
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
}
```

### Modifies-State Tests

**Characteristics:**
- Changes system state
- Changes persist after test
- Reversible with cleanup
- Requires cleanup test

**Example commands:**
```bash
oc apply -f deployment.yaml
oc scale deployment/app --replicas=3
oc patch configmap/config --type=merge
kubectl create namespace test
```

**Template:**
```json
{
  "system_impact": {
    "type": "modifies-state",
    "description": "Creates deployment and service in test namespace",
    "affected_resources": [
      "Namespace: test-ns",
      "Deployment: test-app",
      "Service: test-svc"
    ],
    "reversible": true,
    "persistence": "permanent",
    "risk_level": "medium"
  },
  "state_management": {
    "requires_backup": false,
    "requires_cleanup": true,
    "cleanup_test": "test_cleanup_deployment",
    "restores_state": false
  },
  "safety": {
    "can_run_in_production": false,
    "requires_confirmation": true,
    "warning_message": "⚠️ This test creates resources that will persist until cleanup.",
    "safe_to_retry": true,
    "idempotent": true
  }
}
```

### Destructive Tests

**Characteristics:**
- Deletes or irreversibly modifies data
- Cannot be undone without backup
- High risk if run incorrectly
- Requires backup AND cleanup

**Example commands:**
```bash
oc delete deployment production-app
oc delete pvc data-volume
kubectl drain node-1 --force --delete-emptydir-data
rosa delete cluster --cluster=prod
```

**Template:**
```json
{
  "system_impact": {
    "type": "destructive",
    "description": "Deletes production deployment - CANNOT be undone without backup",
    "affected_resources": [
      "Deployment: production-app",
      "All associated pods",
      "Service endpoints"
    ],
    "reversible": false,
    "persistence": "permanent",
    "risk_level": "critical"
  },
  "state_management": {
    "requires_backup": true,
    "backup_test": "test_backup_production_app",
    "requires_cleanup": true,
    "cleanup_test": "test_restore_production_app",
    "restores_state": false
  },
  "safety": {
    "can_run_in_production": false,
    "requires_confirmation": true,
    "warning_message": "🚨 DESTRUCTIVE: This test will DELETE production resources. Backup MUST complete first. Continue only if you understand the risk.",
    "safe_to_retry": false,
    "idempotent": false
  }
}
```

## Backup/Cleanup/Restore Pattern

### The Four-Test Pattern

For any test that modifies state:

1. **Backup Test** - Save current state
2. **Modify Test** - Make changes
3. **Cleanup Test** - Restore original state
4. **Verify Test** - Confirm restoration

### Example: Modifying Cluster Configuration

**Test 1: Backup**
```json
{
  "test_backup_cluster_config": {
    "metadata": {
      "category": "backup",
      "system_impact": {
        "type": "read-only",
        "description": "Exports cluster config to backup file"
      },
      "state_management": {
        "requires_backup": false,
        "requires_cleanup": false,
        "restores_state": false
      }
    },
    "test_execution": {
      "steps": [
        {
          "step_number": 1,
          "command": "oc get cm cluster-config -o yaml > backup/config-$(date +%Y%m%d-%H%M%S).yaml"
        }
      ]
    }
  }
}
```

**Test 2: Modify**
```json
{
  "test_modify_cluster_config": {
    "metadata": {
      "category": "configuration",
      "system_impact": {
        "type": "modifies-state",
        "description": "Modifies cluster configuration",
        "affected_resources": ["cluster-config ConfigMap"]
      },
      "state_management": {
        "requires_backup": true,
        "backup_test": "test_backup_cluster_config",
        "requires_cleanup": true,
        "cleanup_test": "test_restore_cluster_config"
      }
    },
    "test_execution": {
      "dependencies": ["test_backup_cluster_config"],
      "steps": [
        {
          "step_number": 1,
          "command": "oc patch cm cluster-config --type=merge -p '{...}'"
        }
      ]
    }
  }
}
```

**Test 3: Cleanup/Restore**
```json
{
  "test_restore_cluster_config": {
    "metadata": {
      "category": "cleanup",
      "system_impact": {
        "type": "modifies-state",
        "description": "Restores cluster config from backup"
      },
      "state_management": {
        "requires_backup": false,
        "requires_cleanup": false,
        "restores_state": true,
        "state_verification": "test_verify_config_restored"
      }
    },
    "test_execution": {
      "dependencies": ["test_backup_cluster_config"],
      "steps": [
        {
          "step_number": 1,
          "command": "oc apply -f backup/config-*.yaml"
        }
      ]
    }
  }
}
```

**Test 4: Verify**
```json
{
  "test_verify_config_restored": {
    "metadata": {
      "category": "validation",
      "system_impact": {
        "type": "read-only",
        "description": "Verifies config matches original"
      },
      "state_management": {
        "requires_backup": false,
        "requires_cleanup": false,
        "restores_state": false
      }
    },
    "test_execution": {
      "dependencies": ["test_restore_cluster_config"],
      "steps": [
        {
          "step_number": 1,
          "command": "diff <(oc get cm cluster-config -o yaml) backup/config-*.yaml"
        }
      ]
    }
  }
}
```

## Execution Logic

### Pre-Execution Checks (testplan-executor)

Before running any test:

```bash
# 1. Check system_impact.type
if [ "$type" = "modifies-state" ] || [ "$type" = "destructive" ]; then

  # 2. Check requires_backup
  if [ "$requires_backup" = "true" ]; then
    # Verify backup test completed successfully
    if [ "$backup_test_status" != "passed" ]; then
      echo "❌ Cannot run: Backup test ($backup_test) must complete first"
      exit 1
    fi
  fi

  # 3. Check requires_confirmation
  if [ "$requires_confirmation" = "true" ]; then
    echo "$warning_message"
    read -p "Proceed? (yes/no): " response
    if [ "$response" != "yes" ]; then
      echo "Test execution cancelled by user"
      exit 0
    fi
  fi

  # 4. Check production environment
  if [ "$ENVIRONMENT" = "production" ] && [ "$can_run_in_production" = "false" ]; then
    echo "❌ Cannot run: Test not approved for production"
    exit 1
  fi
fi
```

### Post-Execution Cleanup

After test completes (success or failure):

```bash
# 1. Check if cleanup is required
if [ "$requires_cleanup" = "true" ]; then

  # 2. Prompt for cleanup
  echo "Test complete. Cleanup required to restore system state."
  echo "Cleanup test: $cleanup_test"
  read -p "Run cleanup now? (yes/no): " response

  if [ "$response" = "yes" ]; then
    # 3. Run cleanup test
    execute_test "$cleanup_test"

    # 4. Run verification if defined
    if [ -n "$state_verification" ]; then
      execute_test "$state_verification"
    fi
  else
    echo "⚠️ WARNING: System state modified. Run cleanup manually:"
    echo "  ./execute-testplan --tests $cleanup_test"
  fi
fi
```

## Automated Analysis

The testplan-executor skill analyzes commands to determine intrusiveness:

### Read-Only Detection

**Indicators:**
- Commands: `get`, `describe`, `list`, `show`, `cat`, `curl` (GET)
- No write operations
- No modification flags

**Example:**
```bash
# Read-only
oc get pods
kubectl describe deployment
curl http://api/status
```

### Modifies-State Detection

**Indicators:**
- Commands: `apply`, `create`, `patch`, `scale`, `set`
- Write operations
- State modification flags

**Example:**
```bash
# Modifies state
oc apply -f deployment.yaml
kubectl scale deployment/app --replicas=3
oc patch configmap --type=merge
```

### Destructive Detection

**Indicators:**
- Commands: `delete`, `drain`, `evict`, `purge`
- `--force` flags
- Irreversible operations

**Example:**
```bash
# Destructive
oc delete deployment prod-app
kubectl drain node --force
rosa delete cluster
```

## Best Practices

### For Test Authors

**Do:**
- Always classify tests correctly (read-only vs modifies-state vs destructive)
- Create backup tests for any modifying test
- Create cleanup tests to restore state
- Set appropriate risk_level
- Write clear warning_message for risky tests
- Make tests idempotent when possible

**Don't:**
- Skip state_management fields
- Mark modifying tests as read-only
- Forget to link cleanup_test
- Assume tests can be retried
- Mix read and write operations in same test

### For Test Executors

**Do:**
- Always run backup tests first
- Respect requires_confirmation
- Run cleanup tests after modifying tests
- Verify state was restored
- Log all state changes

**Don't:**
- Skip pre-execution checks
- Run modifying tests in production
- Ignore cleanup requirements
- Retry destructive tests without investigation
- Proceed if backup failed

## Safety Checklist

Before running tests, verify:

- [ ] Read-only tests: Can run freely
- [ ] Modifies-state tests:
  - [ ] Backup test exists and passed
  - [ ] Cleanup test exists
  - [ ] Not running in production
  - [ ] User confirmed execution
- [ ] Destructive tests:
  - [ ] Backup test exists and passed
  - [ ] Restore test exists and tested
  - [ ] Double confirmation from user
  - [ ] Never run in production
  - [ ] Full understanding of impact

## Examples

See complete examples in:
- `.claude/templates/backup-cleanup-pattern.json`
- `.claude/templates/testplan-template.json` (test_1 and test_2)
- `examples/rhobs_test_plan_v2.json` (when updated)
