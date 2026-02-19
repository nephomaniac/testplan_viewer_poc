# Test Plan Generator Skill

Generate comprehensive, well-structured test plans from any combination of artifacts including code repositories, existing tests, documentation, SOPs, and guides. Focus on complete coverage across different environments and configurations.

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
