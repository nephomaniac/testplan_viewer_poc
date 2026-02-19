# Test Plan Generator Skill

**Focus:** Comprehensive test case generation from artifacts

**Purpose:** Generate complete test coverage across different environments and configurations

## What This Skill Does

The Test Plan Generator skill analyzes provided artifacts (code, tests, documentation, SOPs) and generates comprehensive test plans that cover:

- ✅ **All functional paths** - Happy path, edge cases, error handling
- ✅ **Multiple environments** - AWS, GCP, Azure, different cluster types
- ✅ **Input permutations** - Valid, invalid, boundary conditions
- ✅ **Integration scenarios** - Dependencies and external services
- ✅ **Gap analysis** - What's covered and what's missing

## When to Use

Use this skill when you need to:

- Create a new test plan from scratch
- Generate test coverage for a new feature or service
- Expand test coverage across multiple environments
- Systematically cover all configuration variants
- Identify gaps in existing test coverage

## Inputs Required

### Primary Focus (Choose ONE)
- Code repository
- Documentation
- Existing test suite

### Supporting Artifacts (Optional)
- Git repositories
- Test cases in any format
- Product documentation
- SOPs and guides
- Architecture diagrams
- API specifications

### Context
- Target system name
- Environment types to cover
- User personas
- Coverage goals

## Example Usage

### Example 1: Generate from Code Repository

```
I need a comprehensive test plan for the observability-operator.

Focus artifact: /path/to/observability-operator
Supporting artifacts:
  - E2E tests in e2e/tests/
  - OpenShift monitoring docs
  - Troubleshooting SOPs

Environment types:
  - AWS ROSA (single-AZ, multi-AZ)
  - GCP regional clusters
  - Azure ARO

Coverage goals: Installation, configuration, monitoring, troubleshooting
```

**Claude will:**
1. Analyze the repository structure
2. Examine existing E2E tests
3. Review documentation
4. Create coverage matrix for all environments
5. Generate 40-60 test cases covering all scenarios
6. Document gaps and recommendations
7. Output: `examples/observability-operator-comprehensive.json`

### Example 2: Expand Coverage Across Environments

```
I have a basic test plan but need variants for different cloud providers.

Existing test plan: examples/basic-testplan.json
Add coverage for: AWS, GCP, Azure, bare metal
Focus on: Network configuration differences per cloud
```

**Claude will:**
1. Load existing test plan
2. For each test, create cloud-specific variants
3. Add environment-specific setup and validation
4. Document cloud-specific considerations
5. Output: `examples/multi-cloud-testplan.json`

### Example 3: Generate from Documentation

```
Generate test plan from product documentation.

Focus artifact: docs/user-guide.md
Supporting: API specification in docs/api.yaml
Coverage: All documented features and API endpoints
Environments: Development, staging, production
```

**Claude will:**
1. Parse documentation for features
2. Extract API endpoints from spec
3. Create tests for each documented workflow
4. Add environment-specific considerations
5. Output: `examples/docs-based-testplan.json`

## Output Structure

Generated test plan includes:

### Coverage Summary
```json
{
  "coverage_summary": {
    "total_tests": 45,
    "environment_coverage": {
      "aws": 20,
      "gcp": 15,
      "azure": 10
    },
    "functional_coverage": {
      "installation": "complete",
      "configuration": "complete",
      "monitoring": "partial"
    },
    "gaps_identified": [
      "No IPv6 testing",
      "Missing ARM architecture tests"
    ]
  }
}
```

### Environment Variants
Each test includes variants for different environments:
```json
{
  "test_install_operator": {
    "variants": {
      "aws_single_az": {...},
      "aws_multi_az": {...},
      "gcp_regional": {...}
    }
  }
}
```

### Gap Analysis
```json
{
  "gaps_identified": [
    "No tests for GCP private clusters",
    "Missing upgrade path testing"
  ],
  "recommendations": [
    "Add GCP private cluster variant tests",
    "Create upgrade test suite"
  ]
}
```

## Integration with Other Skills

**Works well with:**
- **testplan-educator** - Enhance generated tests with educational content
- **testplan-executor** - Execute the generated comprehensive test plan

**Typical workflow:**
1. **testplan-generator** → Generate comprehensive coverage
2. **testplan-educator** → Add learning objectives and documentation
3. **testplan-executor** → Execute and validate

## Best Practices

### Do:
- Provide as many supporting artifacts as possible
- Specify all environment types you need coverage for
- Define coverage goals clearly (what's most important?)
- Review the coverage summary for gaps
- Iterate if gaps are unacceptable

### Don't:
- Expect perfect coverage on first pass (iterate!)
- Skip environment variants to save time
- Ignore the gap analysis
- Generate more tests than you can reasonably execute

## Success Criteria

A successful generation includes:
- ✅ 70%+ coverage of documented features
- ✅ 3+ environment variants for critical tests
- ✅ Negative test cases included
- ✅ Integration tests for dependencies
- ✅ Documented gaps with recommendations
- ✅ Realistic time estimates

## Files Generated

After using this skill, you'll have:
```
examples/
  [service]-testplan.json          # Complete test plan
  [service]-coverage-report.md     # Coverage analysis
  [service]-gap-analysis.md        # Identified gaps
```

## Next Steps

After generating test plan:
1. Review coverage summary
2. Address high-priority gaps
3. Use **testplan-educator** to enhance with learning content
4. Generate HTML: `make run`
5. Use **testplan-executor** to validate tests
