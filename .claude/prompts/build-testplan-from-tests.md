# Build Test Plan from Existing Tests

## Objective

Generate an educational test plan in JSON format by analyzing existing unit tests, integration tests, and E2E tests. Transform automated test code into manual test procedures that educate engineers about what the tests validate and why.

## Your Role

You are reverse-engineering test coverage into educational material. Your goal is to:

1. **Analyze test code** - Understand what each test validates
2. **Extract test scenarios** - Identify distinct test cases and workflows
3. **Map to manual steps** - Convert automated assertions into manual validation procedures
4. **Add educational context** - Explain why each test matters and what it teaches
5. **Generate valid JSON** - Create properly formatted test plan

## Inputs You'll Receive

The user will provide:

- **Test file paths**: Location of test files to analyze
- **Test framework**: pytest, go test, jest, mocha, junit, etc.
- **Service context**: What the code being tested does
- **Optional: Coverage reports**: To identify gaps in testing

## Analysis Steps

### 1. Test Discovery (5-10 minutes)

**Find all test files:**

```bash
# Go tests
find . -name "*_test.go"
find . -path "*/test/*" -name "*.go"

# Python tests
find . -name "test_*.py" -o -name "*_test.py"
find . -type d -name "tests"

# JavaScript/TypeScript tests
find . -name "*.test.ts" -o -name "*.spec.ts"
find . -name "*.test.js" -o -name "*.spec.js"

# E2E tests
find . -path "*/e2e/*" -o -path "*/integration/*"
```

**Categorize tests:**

- **Unit tests**: Test individual functions/methods
- **Integration tests**: Test component interactions
- **E2E tests**: Test complete user workflows
- **Performance tests**: Test scalability and benchmarks
- **Regression tests**: Test specific bug fixes

### 2. Test Code Analysis (10-20 minutes per test file)

**Read test file structure:**

```python
# Example: Python pytest
def test_create_cluster():
    """Test that cluster creation succeeds with valid params."""
    # Setup
    cluster = ClusterBuilder().with_name("test").build()

    # Execute
    result = api.create_cluster(cluster)

    # Assert
    assert result.status == "success"
    assert result.cluster_id is not None
```

**Identify key elements:**

- **Test name**: What is being tested?
- **Test description/docstring**: Why is this test important?
- **Setup**: What prerequisites are configured?
- **Execution**: What action is performed?
- **Assertions**: What is validated?
- **Teardown**: What cleanup happens?

**Extract test data:**

- Input parameters (valid and invalid)
- Expected outputs and error messages
- Mock data and fixtures
- Environment variables
- Configuration files

### 3. Test Coverage Mapping (5-10 minutes)

**Understand what's covered:**

```bash
# Go coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Python coverage
pytest --cov=. --cov-report=html

# JavaScript coverage
npm test -- --coverage
```

**Identify:**

- Features with test coverage
- Features without test coverage (need manual tests)
- Edge cases being tested
- Error paths being validated
- Happy paths vs failure scenarios

### 4. Manual Test Conversion (20-40 minutes per test)

**For each automated test, create manual equivalent:**

**Example transformation:**

```go
// Automated test
func TestServiceMonitorCreation(t *testing.T) {
    sm := createServiceMonitor("test-sm", "default")
    err := k8sClient.Create(ctx, sm)
    require.NoError(t, err)

    fetched := &v1.ServiceMonitor{}
    err = k8sClient.Get(ctx, types.NamespacedName{
        Name: "test-sm",
        Namespace: "default",
    }, fetched)
    require.NoError(t, err)
    assert.Equal(t, "test-sm", fetched.Name)
}
```

**Becomes manual test case:**

```json
{
  "test_execution": {
    "steps": [
      {
        "step_number": 1,
        "title": "Create ServiceMonitor resource",
        "command": "oc apply -f - <<EOF\napiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: test-sm\n  namespace: default\nEOF",
        "expected_output": "servicemonitor.monitoring.coreos.com/test-sm created",
        "learning_note": "ServiceMonitors tell Prometheus which services to scrape metrics from",
        "why_this_step": "We're creating the monitoring configuration that the operator will reconcile"
      },
      {
        "step_number": 2,
        "title": "Verify ServiceMonitor was created",
        "command": "oc get servicemonitor test-sm -n default -o yaml",
        "expected_output": "name: test-sm\nnamespace: default\n...",
        "learning_note": "Kubernetes confirmation that resource is in etcd",
        "why_this_step": "Validates the resource was accepted by the API server"
      }
    ],
    "validation": {
      "success_criteria": [
        "ServiceMonitor resource exists in default namespace",
        "Resource name matches 'test-sm'",
        "No errors in operator logs"
      ]
    }
  }
}
```

### 5. Educational Enhancement (10-20 minutes per test)

**Add learning value:**

**From test code:**
```python
def test_invalid_cluster_name():
    """Cluster names must match RFC 1123 DNS label."""
    with pytest.raises(ValidationError):
        api.create_cluster(name="INVALID_NAME_123")
```

**Create learning content:**

```json
{
  "learning": {
    "objectives": [
      "Understand Kubernetes DNS label requirements",
      "Learn to validate resource names before creation"
    ],
    "concepts_covered": ["rfc_1123_dns_labels", "k8s_naming"],
    "common_beginner_mistakes": [
      {
        "mistake": "Using uppercase letters in resource names",
        "consequence": "API rejects the resource creation",
        "how_to_avoid": "Always use lowercase alphanumeric + hyphens for k8s names"
      }
    ]
  }
}
```

## Test Framework Patterns

### Go Testing (`*_test.go`)

**Common patterns to recognize:**

```go
// Table-driven tests
func TestValidator(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    bool
        wantErr bool
    }{
        {"valid", "test-123", true, false},
        {"invalid", "TEST", false, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Validate(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.Equal(t, tt.want, got)
            }
        })
    }
}
```

**Convert to:** Multiple test cases, one per table row, showing both valid and invalid inputs

**Integration test patterns:**

```go
// Uses testify suite
type IntegrationSuite struct {
    suite.Suite
    client client.Client
}

func (s *IntegrationSuite) TestE2E() {
    // Multi-step workflow test
}
```

**Convert to:** Complete end-to-end manual test with all steps

### Python pytest

**Common patterns:**

```python
# Fixtures for setup
@pytest.fixture
def cluster():
    c = create_test_cluster()
    yield c
    c.delete()

# Parametrize for multiple cases
@pytest.mark.parametrize("region,expected", [
    ("us-east-1", True),
    ("invalid", False),
])
def test_region(region, expected):
    assert validate_region(region) == expected
```

**Convert to:**
- Prerequisites section (from fixtures)
- Multiple test cases (from parametrize)
- Cleanup steps (from teardown)

### JavaScript/TypeScript Jest/Mocha

**Common patterns:**

```typescript
describe('ServiceMonitor Controller', () => {
  beforeEach(() => {
    // Setup
  });

  it('should create ServiceMonitor', async () => {
    const result = await controller.create(sm);
    expect(result.status).toBe('created');
  });

  afterEach(() => {
    // Cleanup
  });
});
```

**Convert to:** Structured test with setup, execution, validation, cleanup

## Output Structure

Generate test plan JSON with these key sections:

### Metadata (from test context)

```json
{
  "metadata": {
    "document_title": "Test Plan: [Component] - Derived from Automated Tests",
    "version": "1.0",
    "purpose": "Manual validation of [component] based on automated test coverage",
    "target_audience": "Engineers validating [component] functionality",
    "estimated_total_time": "[sum of all test times]"
  }
}
```

### Test Cases (from individual tests)

```json
{
  "testcases": {
    "test_[name]": {
      "metadata": {
        "id": "test_[name]",
        "title": "[Derived from test name/description]",
        "category": "[unit|integration|e2e]",
        "difficulty": "beginner",
        "source_test": "path/to/test_file.go:line_number"
      },
      "learning": {
        "objectives": ["What this test validates"],
        "concepts_covered": ["Extracted from test code"],
        "skills_gained": ["What you learn by running this test"]
      },
      "test_execution": {
        "objective": "Validate [what the automated test validates]",
        "steps": [
          {
            "step_number": 1,
            "title": "[From test setup/execution]",
            "command": "[Manual equivalent of test code]",
            "expected_output": "[From test assertions]"
          }
        ],
        "validation": {
          "success_criteria": ["From test assertions"],
          "how_to_verify": "[Manual verification steps]"
        }
      }
    }
  }
}
```

## Mapping Examples

### Example 1: Unit Test → Manual Test

**Unit test:**
```go
func TestParseClusteID(t *testing.T) {
    id, err := ParseClusterID("abc123")
    assert.NoError(t, err)
    assert.Equal(t, "abc123", id)
}
```

**Manual test case:**
```json
{
  "test_id_parsing": {
    "metadata": {
      "title": "Validate cluster ID parsing logic",
      "category": "validation",
      "source_test": "internal/cluster/cluster_test.go:45"
    },
    "test_execution": {
      "steps": [
        {
          "title": "Parse a valid cluster ID",
          "command": "echo 'abc123' | ./bin/parse-cluster-id",
          "expected_output": "Cluster ID: abc123",
          "why_this_step": "Tests that valid alphanumeric cluster IDs are accepted"
        }
      ]
    }
  }
}
```

### Example 2: Integration Test → Manual Test

**Integration test:**
```python
def test_operator_watches_servicemonitor(k8s_client):
    # Create ServiceMonitor
    sm = create_servicemonitor("test")
    k8s_client.create(sm)

    # Wait for operator to reconcile
    wait_for_prometheus_config_update()

    # Verify Prometheus config includes the ServiceMonitor
    config = get_prometheus_config()
    assert "test" in config["scrape_configs"]
```

**Manual test case:**
```json
{
  "test_operator_reconciliation": {
    "metadata": {
      "title": "Verify operator watches and reconciles ServiceMonitor resources",
      "category": "integration"
    },
    "test_execution": {
      "steps": [
        {
          "step_number": 1,
          "title": "Create ServiceMonitor resource",
          "command": "oc apply -f servicemonitor.yaml"
        },
        {
          "step_number": 2,
          "title": "Wait for operator reconciliation",
          "command": "sleep 30 && oc logs -n openshift-operators deployment/operator -f | grep 'Reconciled ServiceMonitor'",
          "why_this_step": "Operators use reconciliation loops to watch resources"
        },
        {
          "step_number": 3,
          "title": "Verify Prometheus configuration updated",
          "command": "oc get secret prometheus-config -n openshift-monitoring -o jsonpath='{.data.prometheus\\.yml}' | base64 -d | grep 'test'",
          "expected_output": "job_name: serviceMonitor/default/test"
        }
      ]
    }
  }
}
```

### Example 3: E2E Test → Manual Test Plan

**E2E test:**
```javascript
describe('Complete monitoring workflow', () => {
  it('deploys app, creates ServiceMonitor, validates metrics', async () => {
    // Deploy application
    await kubectl.apply('deployment.yaml');
    await waitForPods('app=myapp');

    // Create ServiceMonitor
    await kubectl.apply('servicemonitor.yaml');

    // Verify metrics are scraped
    const metrics = await prometheus.query('up{job="myapp"}');
    expect(metrics.value).toBe(1);
  });
});
```

**Manual test plan:**
```json
{
  "test_complete_monitoring_workflow": {
    "metadata": {
      "title": "End-to-end monitoring workflow validation",
      "category": "e2e",
      "estimated_time": "20 minutes",
      "hands_on_percentage": 90
    },
    "learning": {
      "objectives": [
        "Understand complete monitoring setup workflow",
        "Learn how Prometheus discovers and scrapes metrics",
        "Validate end-to-end metric collection"
      ]
    },
    "test_execution": {
      "objective": "Deploy an application and verify metrics are collected by Prometheus",
      "dependencies": [],
      "steps": [
        {
          "step_number": 1,
          "title": "Deploy sample application",
          "command": "oc apply -f deployment.yaml",
          "expected_output": "deployment.apps/myapp created"
        },
        {
          "step_number": 2,
          "title": "Wait for pods to be ready",
          "command": "oc wait --for=condition=Ready pod -l app=myapp --timeout=60s",
          "expected_output": "pod/myapp-xxxxx condition met"
        },
        {
          "step_number": 3,
          "title": "Create ServiceMonitor to scrape metrics",
          "command": "oc apply -f servicemonitor.yaml",
          "learning_note": "ServiceMonitor tells Prometheus where to find metrics endpoints"
        },
        {
          "step_number": 4,
          "title": "Query Prometheus for metrics",
          "command": "oc exec -n openshift-monitoring prometheus-0 -- promtool query instant http://localhost:9090 'up{job=\"myapp\"}'",
          "expected_output": "up{job=\"myapp\"} => 1"
        }
      ],
      "validation": {
        "success_criteria": [
          "Application pods are running",
          "ServiceMonitor exists and is valid",
          "Prometheus is scraping metrics (up=1)"
        ]
      }
    }
  }
}
```

## Quality Checklist

- [ ] All test files in scope have been analyzed
- [ ] Test categories identified (unit/integration/e2e)
- [ ] Automated assertions mapped to manual validation
- [ ] Test setup converted to prerequisites
- [ ] Test teardown converted to cleanup steps
- [ ] Source test files referenced for traceability
- [ ] Learning objectives added for each test
- [ ] Coverage gaps identified and documented
- [ ] JSON is valid and complete

## Example Usage

**User provides:**
```
Test files: internal/controller/*_test.go, e2e/tests/*.go
Framework: go test (testify)
Service: Observability Operator
```

**You should:**

1. Read all `*_test.go` files
2. Identify test categories and scenarios
3. For each test function:
   - Extract test name and description
   - Identify setup, execute, assert pattern
   - Convert to manual steps
   - Add educational context
4. Create test plan JSON with all tests
5. Group related tests into logical sequence
6. Add prerequisites from test fixtures
7. Save to `examples/observability-operator-from-tests.json`

## Tips

**Do:**
- Preserve test names for traceability
- Link to source test file and line number
- Extract actual error messages from test code
- Use test descriptions as learning notes
- Group related tests together
- Maintain test coverage mapping

**Don't:**
- Skip tests that seem too technical
- Lose information from test assertions
- Forget edge case tests
- Ignore test setup/teardown
- Create tests that can't be manually executed
