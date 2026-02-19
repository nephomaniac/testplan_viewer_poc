# Master Prompt: AI-Assisted Test Plan Generation and Execution

## Overview

This directory contains a comprehensive system for using Claude to generate, enhance, execute, and track educational test plans for software projects. The system transforms code repositories, tests, and documentation into interactive HTML test plans that educate users while they validate functionality.

## The Complete Workflow

```
                            ┌──────────────────────────────────┐
                            │   INPUT: Code Repository         │
                            │   - Source code                  │
                            │   - Existing tests               │
                            │   - Documentation                │
                            └──────────┬───────────────────────┘
                                       │
                    ┌──────────────────┴────────────────────┐
                    │                                        │
         ┌──────────▼────────────┐              ┌──────────▼────────────┐
         │  Build from Repo      │              │  Build from Tests     │
         │  Prompt               │              │  Prompt               │
         └──────────┬────────────┘              └──────────┬────────────┘
                    │                                        │
                    └──────────────────┬────────────────────┘
                                       │
                            ┌──────────▼───────────────────────┐
                            │  Extract Documentation           │
                            │  Prompt                          │
                            └──────────┬───────────────────────┘
                                       │
                            ┌──────────▼───────────────────────┐
                            │  OUTPUT: Test Plan JSON          │
                            │  - Metadata                      │
                            │  - Concepts                      │
                            │  - Prerequisites                 │
                            │  - Test Cases                    │
                            └──────────┬───────────────────────┘
                                       │
                     ┌─────────────────┴─────────────────┐
                     │  Generate HTML with Viewer        │
                     └─────────────────┬─────────────────┘
                                       │
                            ┌──────────▼───────────────────────┐
                            │  Interactive HTML Test Plan      │
                            │  - Manual test execution         │
                            │  - Educational content           │
                            │  - Progress tracking             │
                            └──────────┬───────────────────────┘
                                       │
                     ┌─────────────────┴─────────────────┐
                     │                                    │
          ┌──────────▼────────────┐          ┌──────────▼────────────┐
          │  Execute Test Plan    │          │  Manual Execution     │
          │  (Automated)          │          │  by User              │
          └──────────┬────────────┘          └──────────┬────────────┘
                     │                                    │
                     └──────────────────┬─────────────────┘
                                        │
                             ┌──────────▼───────────────────────┐
                             │  Update with Results             │
                             │  Prompt                          │
                             └──────────┬───────────────────────┘
                                        │
                             ┌──────────▼───────────────────────┐
                             │  OUTPUT:                         │
                             │  - Updated JSON with results     │
                             │  - Execution report              │
                             │  - Updated HTML                  │
                             └──────────────────────────────────┘
```

## Prompts and Their Purposes

### 1. build-testplan-from-repo.md

**When to use:** You have a code repository and want to create a comprehensive test plan

**What it does:**
- Analyzes repository structure
- Examines code and identifies key components
- Designs learning path
- Creates test cases covering main functionality
- Generates complete test plan JSON

**Typical interaction:**
```
User: "Create a test plan for the observability-operator repository targeting new SRE engineers"

Claude: [Uses build-testplan-from-repo.md]
1. Explores repository structure
2. Reads README and documentation
3. Analyzes code organization
4. Creates 3-5 core test cases
5. Outputs: examples/observability-operator-testplan.json
```

### 2. build-testplan-from-tests.md

**When to use:** You have existing automated tests and want manual test plan

**What it does:**
- Analyzes test files (unit, integration, E2E)
- Extracts test scenarios and assertions
- Converts automated tests to manual procedures
- Maps test coverage to manual test cases
- Generates test plan JSON from test code

**Typical interaction:**
```
User: "Create a test plan based on the E2E tests in e2e/tests/"

Claude: [Uses build-testplan-from-tests.md]
1. Reads all test files in e2e/tests/
2. Identifies test scenarios
3. Converts each test to manual steps
4. Maps assertions to validation criteria
5. Outputs: examples/component-from-tests.json
```

### 3. extract-documentation.md

**When to use:** You want to enhance an existing test plan with better documentation

**What it does:**
- Searches for relevant documentation (repo, OpenShift, upstream)
- Extracts key concepts and creates explanations
- Finds architecture diagrams
- Adds learning resources
- Enhances test cases with context

**Typical interaction:**
```
User: "Enhance this test plan with documentation for Prometheus and ServiceMonitor concepts"

Claude: [Uses extract-documentation.md]
1. Searches Prometheus Operator docs
2. Finds OpenShift monitoring guides
3. Locates architecture diagrams
4. Creates concept explanations
5. Adds references to test steps
6. Outputs: examples/enhanced-testplan.json
```

### 4. execute-testplan.md

**When to use:** You want Claude to execute a test plan automatically

**What it does:**
- Loads test plan JSON
- Verifies prerequisites
- Executes each test step
- Captures command outputs
- Validates results
- Determines pass/fail status
- Logs execution details

**Typical interaction:**
```
User: "Execute the RHOBS test plan on my current cluster"

Claude: [Uses execute-testplan.md]
1. Verifies cluster access
2. Checks prerequisites
3. Runs each test step sequentially
4. Captures all outputs
5. Validates against success criteria
6. Outputs: test-results/YYYYMMDD_HHMMSS/results.json
```

### 5. update-testplan-results.md

**When to use:** After test execution, to merge results into test plan

**What it does:**
- Loads original test plan
- Loads execution results
- Merges results into test plan structure
- Adds execution metadata
- Generates summary report
- Creates updated HTML

**Typical interaction:**
```
User: "Update the test plan with the execution results from today's run"

Claude: [Uses update-testplan-results.md]
1. Loads examples/testplan.json
2. Loads test-results/20260218_103000/results.json
3. Merges results into test plan
4. Generates markdown report
5. Outputs:
   - examples/testplan-results.json
   - test-results/20260218_103000/REPORT.md
   - test-results/20260218_103000/testplan-results.html
```

## Usage Patterns

### Pattern 1: Fresh Test Plan from New Repository

**Goal:** Create test plan for a repository you've never tested before

**Steps:**
1. Clone or navigate to repository
2. Use `build-testplan-from-repo.md`:
   - Provide repository path
   - Specify target audience
   - Define focus areas
3. Use `extract-documentation.md`:
   - Research concepts
   - Add official documentation links
4. Generate HTML and review
5. Iterate based on feedback

**Example:**
```
User: "I need a test plan for the new monitoring operator"

Claude: "I'll create a comprehensive test plan. Let me start by exploring the repository..."

[Claude follows build-testplan-from-repo.md]
[Then enhances with extract-documentation.md]

Claude: "I've created examples/monitoring-operator-testplan.json with 5 test cases covering installation, configuration, and validation. Generating HTML now..."
```

### Pattern 2: Leverage Existing Tests

**Goal:** Create manual test plan based on existing automated tests

**Steps:**
1. Identify test files
2. Use `build-testplan-from-tests.md`:
   - Provide test file paths
   - Specify test framework
3. Use `extract-documentation.md`:
   - Add concept explanations
   - Link to official docs
4. Review and refine

**Example:**
```
User: "We have comprehensive unit tests in internal/*/test.go. Can you create a manual test plan from them?"

Claude: "I'll analyze your test files and create manual test procedures..."

[Claude follows build-testplan-from-tests.md]

Claude: "I've created examples/component-from-tests.json with 12 test cases derived from your unit tests. Each automated assertion is now a manual validation step."
```

### Pattern 3: Automated Test Execution

**Goal:** Run test plan automatically and capture results

**Steps:**
1. Have existing test plan JSON
2. Prepare environment (cluster, credentials)
3. Use `execute-testplan.md`:
   - Specify test plan file
   - Choose execution mode (automated/interactive)
   - Set scope (all tests or specific tests)
4. Use `update-testplan-results.md`:
   - Merge results into test plan
   - Generate report

**Example:**
```
User: "Execute examples/rhobs_test_plan_v2.json on my cluster and generate a report"

Claude: "I'll execute the test plan and capture results..."

[Claude follows execute-testplan.md]

Claude: "Execution complete. 2 of 3 tests passed. Now updating the test plan with results..."

[Claude follows update-testplan-results.md]

Claude: "Results have been merged. See test-results/20260218_103000/REPORT.md for details."
```

### Pattern 4: Incremental Enhancement

**Goal:** Start with basic test plan and progressively enhance it

**Steps:**
1. Create minimal test plan (build-testplan-from-repo.md)
2. Add test coverage (build-testplan-from-tests.md)
3. Enhance with docs (extract-documentation.md)
4. Execute and capture results (execute-testplan.md)
5. Update with results (update-testplan-results.md)
6. Iterate based on findings

**Example:**
```
Day 1: Basic structure from code
Day 2: Add tests from test suite
Day 3: Research and add documentation
Day 4: Execute and identify gaps
Day 5: Refine based on execution results
```

## Key Concepts

### Test Plan JSON Schema

All prompts work with the same JSON schema defined in `templates/testplan-template.json`:

**Core sections:**
- `metadata` - Document info, audience, objectives
- `learning_path` - Recommended test sequence
- `prerequisites` - Required knowledge, access, environment
- `concepts` - Technical concepts explained
- `testcases` - Individual test cases with steps
- `progress_tracking` - Browser-based progress config

**Test case structure:**
- `metadata` - Test info (ID, title, difficulty, time)
- `learning` - Educational content (objectives, skills, mistakes)
- `test_execution` - Steps, validation, dependencies
- `troubleshooting` - Common failures and solutions
- `next_steps` - What to do after completion
- `references` - Documentation links

### Educational Philosophy

These prompts are designed to create **learning-first** test plans:

**Principles:**
1. **Context before commands** - Explain why before showing how
2. **Progressive complexity** - Start simple, build to advanced
3. **Learn by doing** - Hands-on execution with guidance
4. **Mistakes as learning** - Highlight common pitfalls
5. **Multiple resources** - Docs, videos, tutorials, code

**Not just test execution:**
- Traditional: "Run this command, check output matches"
- Educational: "This command validates X. You're learning Y. Common mistake is Z."

### Status Values

**Test status:**
- `passed` - All steps successful, criteria met
- `failed` - Critical failure or criteria not met
- `needs_review` - Ambiguous result, human judgment needed
- `skipped` - Not executed due to blocker
- `blocked` - Environmental issue prevented execution

**Step status:**
- `passed` - Command succeeded, output matched
- `failed` - Command failed or output incorrect
- `needs_review` - Succeeded but output format differs
- `skipped` - Not executed

## Best Practices

### For Test Plan Generation

1. **Be specific about audience** - "New engineers" vs "Experienced SREs" vs "QE team"
2. **Start with real use cases** - What do users actually do with this service?
3. **Balance breadth and depth** - 3-5 comprehensive tests > 20 shallow tests
4. **Include failure scenarios** - Not just happy path
5. **Link to official docs** - Don't recreate documentation, reference it
6. **Validate early** - Generate HTML frequently to check formatting

### For Test Execution

1. **Verify prerequisites** - Don't skip environmental checks
2. **Capture everything** - All outputs, timestamps, environment details
3. **Be patient** - Long-running tests need monitoring
4. **Document deviations** - Note anything different from expected
5. **Take screenshots** - Visual validation requires images
6. **Keep artifacts** - Logs, configs, screenshots for debugging

### For Result Updates

1. **Preserve history** - Don't overwrite previous results
2. **Add context** - Environment, versions, conditions
3. **Be specific** - "ServiceMonitor label mismatch" > "config error"
4. **Include next steps** - What should happen next?
5. **Track trends** - Compare with previous executions

## Advanced Usage

### Chaining Prompts

Prompts can be used in sequence for comprehensive workflow:

```bash
# 1. Generate from repo
cat .claude/prompts/build-testplan-from-repo.md
# → outputs: examples/service-testplan.json

# 2. Enhance with tests
cat .claude/prompts/build-testplan-from-tests.md
# → outputs: examples/service-testplan-enhanced.json

# 3. Add documentation
cat .claude/prompts/extract-documentation.md
# → outputs: examples/service-testplan-final.json

# 4. Execute tests
cat .claude/prompts/execute-testplan.md
# → outputs: test-results/YYYYMMDD_HHMMSS/results.json

# 5. Update with results
cat .claude/prompts/update-testplan-results.md
# → outputs: examples/service-testplan-results.json + REPORT.md
```

### Partial Execution

Execute specific tests only:

```
User: "Execute only test_2 and test_3 from the plan"

Claude: [Uses execute-testplan.md with scope: ["test_2", "test_3"]]
```

### Resume from Failure

Continue execution after fixing an issue:

```
User: "Resume execution from test_2, step 4"

Claude: [Uses execute-testplan.md with resume_from: "test_2:step_4"]
```

### CI/CD Integration

Automate test execution in pipeline:

```yaml
# .github/workflows/test-execution.yml
- name: Execute test plan
  run: |
    claude execute-testplan \
      --plan examples/testplan.json \
      --mode automated \
      --output test-results/

- name: Update results
  run: |
    claude update-testplan-results \
      --plan examples/testplan.json \
      --results test-results/results.json \
      --output examples/testplan-results.json
```

## Troubleshooting

### "Claude can't find the repository"

**Solution:** Provide explicit path
```
User: "Create test plan for repository at /Users/me/projects/my-operator"
```

### "Test plan JSON is invalid"

**Solution:** Validate with viewer
```bash
./build/testplan-viewer -i examples/testplan.json -o /tmp/test.html
```

### "Execution fails with authentication error"

**Solution:** Verify prerequisites before execution
```bash
oc whoami
aws sts get-caller-identity
```

### "Results merge failed"

**Solution:** Ensure test IDs match
```bash
# Check original test IDs
jq '.testcases | keys' examples/testplan.json

# Check results test IDs
jq '.test_results | keys' test-results/results.json
```

## Files and Directories

```
.claude/
├── prompts/
│   ├── README.md                          # Prompts overview
│   ├── MASTER-PROMPT.md                   # This file
│   ├── build-testplan-from-repo.md       # Generate from code
│   ├── build-testplan-from-tests.md      # Generate from tests
│   ├── extract-documentation.md          # Add documentation
│   ├── execute-testplan.md               # Execute tests
│   └── update-testplan-results.md        # Merge results
└── templates/
    ├── testplan-template.json            # Full test plan template
    └── testcase-template.json            # Single test case template

examples/
├── README.md                              # Examples documentation
├── rhobs_test_plan_v2.json               # Example test plan
└── testplan-results.json                 # Example with results

test-results/
└── YYYYMMDD_HHMMSS/
    ├── results.json                      # Execution results
    ├── execution.log                     # Execution log
    ├── REPORT.md                         # Summary report
    └── testplan-results.html             # HTML with results
```

## Next Steps

1. **Read the prompts** - Familiarize yourself with each prompt's purpose
2. **Review templates** - Understand the JSON schema
3. **Try an example** - Generate a test plan for a small project
4. **Execute it** - Run through the test plan manually or automatically
5. **Iterate** - Refine based on what you learn

## Questions?

See detailed documentation in each prompt file:
- Specific usage examples
- Input/output formats
- Quality checklists
- Tips for success
