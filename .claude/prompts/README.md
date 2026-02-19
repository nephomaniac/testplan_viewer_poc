# Test Plan Generation and Execution Prompts

This directory contains prompts and templates for Claude to automatically generate educational test plans from existing codebases, tests, and documentation, as well as execute those test plans.

## Workflow Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                    Test Plan Generation                         │
│                                                                 │
│  Code Repo + Tests + Docs  ──>  JSON Test Plan  ──>  HTML     │
│                                                                 │
│  Claude Prompts:                                               │
│  1. build-testplan-from-repo.md                               │
│  2. build-testplan-from-tests.md                              │
│  3. extract-documentation.md                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ (Manual execution or automation)
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Test Plan Execution                          │
│                                                                 │
│  JSON Test Plan  ──>  Execute Tests  ──>  Update Results      │
│                                                                 │
│  Claude Prompts:                                               │
│  4. execute-testplan.md                                        │
│  5. update-testplan-results.md                                 │
└─────────────────────────────────────────────────────────────────┘
```

## Prompt Files

### 1. build-testplan-from-repo.md
**Purpose**: Generate a complete test plan from a code repository

**Inputs**:
- Repository path or URL
- Target service/component name
- Optional: Specific features to focus on

**Outputs**:
- Complete JSON test plan with:
  - Metadata (title, audience, objectives)
  - Prerequisites (knowledge, access, setup)
  - Test cases with step-by-step instructions
  - Learning objectives and common mistakes
  - Troubleshooting guides

**Example Usage**:
```bash
# In your project directory
cat .claude/prompts/build-testplan-from-repo.md
# Then provide the repository context to Claude
```

### 2. build-testplan-from-tests.md
**Purpose**: Generate test plan from existing unit/integration/e2e tests

**Inputs**:
- Path to test files (pytest, go test, jest, etc.)
- Test framework being used
- Optional: Coverage reports

**Outputs**:
- Test plan JSON based on actual test implementations
- Maps test cases to code coverage
- Explains what each test validates

**Example Usage**:
```bash
# Analyze Go tests
cat .claude/prompts/build-testplan-from-tests.md
# Provide test file paths like: internal/*/test.go, e2e/tests/
```

### 3. extract-documentation.md
**Purpose**: Find and incorporate relevant documentation into test cases

**Inputs**:
- Repository with docs/ directory
- Links to external documentation (OpenShift, upstream)
- Concept names to research

**Outputs**:
- Enhanced test plan with:
  - Concept explanations
  - Links to official documentation
  - Architecture diagrams
  - Best practices

**Example Usage**:
```bash
# Extract docs for ROSA operators
cat .claude/prompts/extract-documentation.md
# Provide service name: "ROSA HCP Monitoring"
```

### 4. execute-testplan.md
**Purpose**: Execute a test plan and validate results

**Inputs**:
- Test plan JSON file
- Environment context (cluster, credentials)
- Execution scope (all tests, specific test, resume from step)

**Outputs**:
- Test execution results (pass/fail/skip)
- Actual command outputs
- Error messages and logs
- Validation of success criteria

**Example Usage**:
```bash
# Execute a test plan
cat .claude/prompts/execute-testplan.md
# Provide: examples/rhobs_test_plan_v2.json
```

### 5. update-testplan-results.md
**Purpose**: Update test plan JSON with execution results

**Inputs**:
- Original test plan JSON
- Execution results (from execute-testplan.md)
- Test status (passed, failed, needs_review, skipped)

**Outputs**:
- Updated JSON test plan with:
  - Execution timestamps
  - Pass/fail status per test case
  - Actual outputs captured
  - Notes on failures
  - Reviewer comments

**Example Usage**:
```bash
# Update results after execution
cat .claude/prompts/update-testplan-results.md
# Provide execution logs and original JSON
```

## Templates

### testplan-template.json
Complete JSON schema template with all required and optional fields.

### testcase-template.json
Template for a single test case with comprehensive structure.

## Getting Started

### Generate a Test Plan from Scratch

1. **Clone or navigate to your target repository**
   ```bash
   cd /path/to/your/repository
   ```

2. **Provide Claude with the build prompt**
   ```bash
   cat .claude/prompts/build-testplan-from-repo.md
   ```

3. **Give Claude context about your service**
   - Service name: "My ROSA Operator"
   - Target audience: "New SRE engineers"
   - Focus areas: "Installation, configuration, monitoring"

4. **Review and refine the generated JSON**
   - Claude will create a JSON file in `examples/`
   - Validate with: `./build/testplan-viewer -i examples/your-plan.json`

5. **Generate HTML for manual testing**
   ```bash
   make run
   ```

### Execute an Existing Test Plan

1. **Prepare your environment**
   - Ensure cluster access
   - Set required environment variables
   - Verify prerequisites from test plan

2. **Provide Claude with the execution prompt**
   ```bash
   cat .claude/prompts/execute-testplan.md
   ```

3. **Specify the test plan**
   ```bash
   # Provide path to JSON test plan
   examples/rhobs_test_plan_v2.json
   ```

4. **Claude will execute each test and report results**
   - Run commands from test steps
   - Validate outputs against expected results
   - Mark tests as pass/fail/needs_review

5. **Update the JSON with results**
   ```bash
   cat .claude/prompts/update-testplan-results.md
   ```

## Use Cases

### Use Case 1: Onboarding New Engineer
**Scenario**: New engineer needs to learn ROSA HCP monitoring

**Steps**:
1. Use `build-testplan-from-repo.md` on observability-operator repo
2. Claude analyzes code, tests, and documentation
3. Generates educational test plan with learning objectives
4. Engineer uses HTML to step through tests manually
5. Engineer learns by doing, with context at each step

### Use Case 2: Validating New Feature
**Scenario**: New feature needs manual validation before release

**Steps**:
1. Use `build-testplan-from-tests.md` on feature branch tests
2. Claude generates test plan from unit + e2e tests
3. QE engineer reviews HTML and executes manually
4. Use `update-testplan-results.md` to record outcomes
5. Share results with team for review

### Use Case 3: Automated Test Execution
**Scenario**: Run nightly validation of service health

**Steps**:
1. Existing test plan JSON in examples/
2. Use `execute-testplan.md` in CI/CD pipeline
3. Claude runs all tests automatically
4. Results updated with `update-testplan-results.md`
5. Failed tests trigger alerts

### Use Case 4: Documentation Generation
**Scenario**: Need to document manual testing procedures

**Steps**:
1. Use `build-testplan-from-repo.md` + `extract-documentation.md`
2. Claude creates comprehensive test plan with docs
3. Generate HTML with diagrams and explanations
4. Publish as internal documentation
5. Update as code evolves

## Advanced Features

### Incremental Test Plan Building

Start with basic structure, then enhance:

```bash
# Step 1: Basic structure from code
cat .claude/prompts/build-testplan-from-repo.md

# Step 2: Enhance with test coverage
cat .claude/prompts/build-testplan-from-tests.md

# Step 3: Add documentation and context
cat .claude/prompts/extract-documentation.md
```

### Partial Execution

Execute specific tests or resume from failure:

```bash
# Execute only test_1 and test_2
cat .claude/prompts/execute-testplan.md
# Specify: --tests test_1,test_2

# Resume from failed step
cat .claude/prompts/execute-testplan.md
# Specify: --resume-from test_2:step_5
```

### Result Comparison

Compare execution results over time:

```bash
# Keep versioned results
examples/results/
  rhobs_test_plan_v2_results_2026-02-01.json
  rhobs_test_plan_v2_results_2026-02-15.json
  rhobs_test_plan_v2_results_2026-02-18.json
```

## Best Practices

### For Test Plan Generation

1. **Be Specific**: Provide clear context about the service and audience
2. **Include Examples**: Share existing test files for reference
3. **Iterate**: Start simple, then enhance with more detail
4. **Validate Early**: Generate HTML frequently to check formatting
5. **Link Documentation**: Always provide URLs to official docs

### For Test Execution

1. **Prepare Environment**: Verify all prerequisites before execution
2. **Capture Output**: Save all command outputs for debugging
3. **Note Deviations**: Record any differences from expected results
4. **Time Execution**: Track how long tests actually take
5. **Update Estimates**: Refine time estimates based on actual execution

### For Result Updates

1. **Preserve History**: Keep previous results for comparison
2. **Add Context**: Include environment details (cluster version, etc.)
3. **Flag Blockers**: Mark tests that cannot run due to blockers
4. **Include Screenshots**: Attach images for visual validation steps
5. **Review Before Merge**: Have second set of eyes on results

## Troubleshooting

### Claude Can't Find Tests
- Ensure test files are in the working directory
- Provide explicit paths: `find . -name "*test*.go"`
- Check test file naming conventions

### Generated JSON Invalid
- Validate with: `./build/testplan-viewer -i file.json`
- Check required fields in schema
- Ensure proper JSON escaping

### Execution Fails
- Verify environment prerequisites
- Check cluster connectivity
- Validate credentials and permissions
- Review test dependencies order

### Documentation Not Found
- Provide direct URLs to docs
- Check repo for README, docs/ directory
- Search upstream project documentation

## Contributing

To add new prompts:

1. Create markdown file in `.claude/prompts/`
2. Follow existing prompt structure
3. Include examples and expected outputs
4. Update this README with usage instructions
5. Test the prompt with Claude before committing

## Related Documentation

- [Test Plan Schema](../../examples/README.md) - JSON schema reference
- [Project README](../../README.md) - Main project documentation
- [Quick Start](../../docs/QUICKSTART.md) - Getting started guide
