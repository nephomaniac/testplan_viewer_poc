# Test Plan Skills for Claude Code

Specialized Claude Code skills for generating, enhancing, and executing test plans from code repositories, documentation, and other artifacts.

---

## Quick Start

### Installation

These skills are designed for **Claude Code CLI** (the official Claude desktop/terminal application).

**Prerequisites**:
- Claude Code CLI installed ([download here](https://claude.ai/download))
- This repository cloned locally
- Git repository to test (or documentation/existing tests)

**Setup**:
```bash
# 1. Clone this repository
git clone https://github.com/nephomaniac/testplan_tools_poc
cd testplan_tools_poc

# 2. Skills are automatically available when you run Claude Code from this directory
cd /path/to/testplan_tools_poc
claude-code

# OR add skills to your global Claude Code skills directory (optional)
mkdir -p ~/.claude/skills
ln -s $(pwd)/skills/* ~/.claude/skills/

# 3. Verify skills are available
# In Claude Code, type:
#   What skills are available?
# You should see: testplan-generator, testplan-educator, testplan-executor
```

### Updating Skills

```bash
# Pull latest changes from repository
cd /path/to/testplan_tools_poc
git pull origin main

# Skills are automatically updated when you restart Claude Code
# or when you navigate to this directory
```

---

## Available Skills

| Skill | Purpose | When to Use | Duration |
|-------|---------|-------------|----------|
| **testplan-generator** | Analyze artifacts and generate comprehensive test plan JSON | Starting a new test plan from code/docs/tests | 1-2 hours |
| **testplan-educator** | Enhance test plan with educational content, hyperlinks, learning notes | After initial generation, or to improve existing plan | 30-60 min |
| **testplan-executor** | Execute tests, collect results, generate reports | Running tests against a cluster, or reviewing execution approach | 1-4 hours |

---

## How to Use the Skills

### Basic Workflow

Skills are invoked conversationally in Claude Code:

```
# In Claude Code chat:

"Please use the testplan-generator skill to create a test plan for
the operator at ~/sandbox/my-operator"
```

Claude will:
1. Recognize the skill name
2. Load the skill definition from `skills/testplan-generator/skill.md`
3. Follow the methodology defined in the skill
4. Generate comprehensive test plan

### Providing Source Artifacts

**The quality of your test plan depends on the artifacts you provide.** Follow this checklist:

#### ✅ Minimum Required (Choose ONE)

**Option 1: Git Repository** (Most Common)
```
Source: Git repository at ~/sandbox/my-operator
Purpose: Generate test plan for operator installation and operations

Please analyze:
  - Code in controllers/, api/, pkg/
  - Existing tests in test/e2e/, test/integration/
  - Documentation in docs/, README.md
  - Configuration examples in config/, deploy/
```

**Option 2: Documentation**
```
Source: Documentation at ~/sandbox/product-docs/
Purpose: Generate user acceptance tests for product features

Please analyze:
  - User guides in guides/*.md
  - API documentation in api-reference/
  - Installation instructions in install/
  - Troubleshooting guides in troubleshooting/
```

**Option 3: Existing Test Suite**
```
Source: JUnit test results at test-results/*.xml
Purpose: Convert existing tests to educational test plan format

Please analyze:
  - Test structure and naming conventions
  - Setup and teardown patterns
  - Test data and fixtures
  - Assertions and validation logic
```

#### ✅ Highly Recommended Supporting Artifacts

**Provide these to create robust, comprehensive test plans**:

**Code & Configuration**:
- ✅ Related repositories (dependencies, integrations): `~/sandbox/related-repo`
- ✅ Configuration examples: `config/samples/*.yaml`
- ✅ Helm charts / deployment manifests: `deploy/`, `charts/`
- ✅ API specifications: `api/openapi.yaml`, `proto/*.proto`

**Documentation**:
- ✅ Architecture documentation: `docs/architecture.md`, design docs
- ✅ User guides: `docs/user-guide.md`, getting started guides
- ✅ Runbooks/SOPs: `docs/runbooks/`, operational procedures
- ✅ Troubleshooting guides: Known issues, FAQs, debugging guides
- ✅ Release notes: What's new, breaking changes, deprecations

**Tests & Quality**:
- ✅ Existing tests: `test/e2e/`, `test/integration/`, `test/unit/`
- ✅ Test fixtures: `testdata/`, sample configurations, mock data
- ✅ CI/CD configurations: `.github/workflows/`, `Jenkinsfile`
- ✅ Test coverage reports: What's tested vs not tested

**Operational Data**:
- ✅ Bug reports: JIRA queries, GitHub issues (provide links or exports)
- ✅ Incident postmortems: What went wrong and how was it fixed
- ✅ Monitoring dashboards: Grafana URLs, Prometheus queries
- ✅ Known limitations: Performance constraints, supported configurations

#### ✅ Context Information

**Always provide this context**:

```
Target system: "ROSA HCP Monitoring Operator"

Environment types to cover:
  - AWS, GCP (cloud providers)
  - Single-AZ, Multi-AZ (availability zones)
  - HCP, Classic (cluster types)

User personas:
  - SRE engineers (primary)
  - QE team (testing)
  - New engineers (onboarding)

Coverage goals:
  - Installation and upgrades (focus)
  - Configuration and customization
  - Monitoring and observability
  - Troubleshooting common issues

Difficulty level:
  - 30% beginner (basic operations)
  - 50% intermediate (standard workflows)
  - 20% advanced (edge cases, debugging)
```

#### 📋 Example: Comprehensive Artifact Specification

**Best practice for creating a robust test plan**:

```
Please create a comprehensive test plan for the Configure AlertManager Operator (CAMO).

PRIMARY ARTIFACT:
  Git repository: ~/sandbox/camo/configure-alertmanager-operator

SUPPORTING ARTIFACTS:
  Code:
    - Controllers: controllers/*.go
    - API definitions: api/v1alpha1/*.go
    - E2E tests: test/e2e/*.go
    - Unit tests: *_test.go

  Documentation:
    - README.md (overview and quick start)
    - docs/architecture.md (system design)
    - docs/LOCAL_TESTING.md (development guide)
    - Inline godoc comments

  Configuration:
    - CRD examples: config/samples/*.yaml
    - RBAC definitions: config/rbac/*.yaml
    - Deployment manifests: deploy/*.yaml

  Related Repositories:
    - Alertmanager: (for integration understanding)
    - Prometheus Operator: (for ServiceMonitor patterns)

  Bug Reports:
    - JIRA query: project=SREP AND component=CAMO
    - Common issues: Webhook failures, RBAC errors, config validation

CONTEXT:
  Target System: Configure AlertManager Operator (CAMO)
  Purpose: Automate Alertmanager configuration for OpenShift Dedicated clusters

  Environment Types:
    - ROSA HCP clusters
    - OpenShift Dedicated clusters
    - Stage environment (for testing)

  User Personas:
    - SRE engineers (install, configure, troubleshoot)
    - QE team (E2E testing)
    - New engineers onboarding to CAMO

  Coverage Goals:
    - Installation and RBAC setup
    - Secret creation and validation
    - Webhook configuration
    - Alertmanager config updates
    - Troubleshooting common failures

  Difficulty Distribution:
    - 40% beginner (verify installation, check logs)
    - 40% intermediate (create secrets, configure webhooks)
    - 20% advanced (debug validation, handle edge cases)

  Safety Requirements:
    - Read-only tests for production validation
    - Modifying tests must have cleanup
    - All tests must be idempotent where possible

OUTPUT:
  Directory: examples/camo/
  Files needed:
    - camo-testplan.json (test plan)
    - camo-testplan.html (interactive viewer)
    - README.md (example documentation)
    - EXECUTION-REVIEW.md (safety analysis)

WORKFLOW:
  1. Use testplan-generator skill to analyze and generate tests
  2. Use testplan-educator skill to add educational content
  3. Use testplan-viewer to generate HTML
  4. Use testplan-executor to review execution approach (don't run yet)
  5. Document improvements to testplan_tools_poc repository

Please proceed with comprehensive analysis and test generation.
```

---

## Skill Usage Patterns

### Pattern 1: New Test Plan from Git Repository

**When**: You have a codebase and want comprehensive test coverage

**Prompt Template**:
```
Create a comprehensive test plan for <SYSTEM_NAME> at <PATH>.

Primary artifact: Git repository at <PATH>
Supporting artifacts:
  - Documentation: <PATHS>
  - Existing tests: <PATHS>
  - Configuration examples: <PATHS>
  - Related repositories: <PATHS>

Context:
  - Target audience: <PERSONAS>
  - Environment types: <ENVIRONMENTS>
  - Coverage focus: <FOCUS_AREAS>
  - Difficulty: <DISTRIBUTION>

Output: examples/<SYSTEM_NAME>/

Workflow:
  1. testplan-generator - Analyze and generate
  2. testplan-educator - Add educational content
  3. testplan-viewer - Generate HTML
  4. testplan-executor - Review approach (don't run)
  5. Document improvements
```

**Example**:
```
Create a comprehensive test plan for route-monitor-operator at ~/sandbox/route-monitor-operator.

Primary artifact: Git repository at ~/sandbox/route-monitor-operator
Supporting artifacts:
  - E2E tests: test/e2e/*.go
  - API definitions: api/v1alpha1/*.go
  - Documentation: README.md, docs/*.md
  - Related: RHOBS Synthetics API integration

Context:
  - Target audience: SRE engineers, QE team
  - Environment types: ROSA HCP, OpenShift Dedicated
  - Coverage focus: Installation, RouteMonitor CR creation, metrics validation
  - Difficulty: 30% beginner, 50% intermediate, 20% advanced

Output: examples/route-monitor-operator/

Please proceed with comprehensive test plan generation.
```

### Pattern 2: Enhance Existing Test Plan

**When**: You have a test plan and want to improve educational content

**Prompt Template**:
```
Enhance the test plan at examples/<SYSTEM_NAME>/<SYSTEM_NAME>-testplan.json
with educational content.

Use testplan-educator skill to:
  1. Add hyperlinks for technical terms (first occurrence)
  2. Enhance learning_note for every step
  3. Add conceptual_overview for complex tests
  4. Expand common_errors with troubleshooting
  5. Add "why this matters" context

Regenerate HTML when done.
```

**Example**:
```
Enhance examples/camo/camo-testplan.json with educational content.

Focus:
  - Hyperlink these terms: Secret, ValidatingWebhook, ConfigMap, ServiceMonitor
  - Add conceptual_overview for test_configure_pagerduty_integration
  - Expand common_errors for test_validate_alertmanager_config
  - Add real-world context for when to use each test

Use testplan-educator skill, then regenerate HTML.
```

### Pattern 3: Execute Tests (or Review Execution)

**When**: You want to run tests or understand execution approach

**Prompt Template for Review** (Safe):
```
Review the execution approach for examples/<SYSTEM_NAME>/<SYSTEM_NAME>-testplan.json

Use testplan-executor skill to:
  1. Analyze dependencies and execution order
  2. Classify safety (read-only vs modifies-state)
  3. Identify cleanup requirements
  4. Create pre-execution checklist
  5. Document execution flow
  6. Identify gaps and improvements

DO NOT execute tests, only review the approach.
Output: examples/<SYSTEM_NAME>/EXECUTION-REVIEW.md
```

**Prompt Template for Execution** (Requires Cluster):
```
Execute tests from examples/<SYSTEM_NAME>/<SYSTEM_NAME>-testplan.json

Environment:
  - Cluster: <CLUSTER_TYPE> (verified NON-PRODUCTION)
  - Access: KUBECONFIG=<PATH>
  - Credentials: <ENV_VARS or credentials file>

Mode: <read-only | full>

If read-only:
  - Only run tests where system_impact.type == "read-only"
  - Generate validation report
  - Safe to run in any environment

If full:
  - Run all tests in dependency order
  - Track cleanup requirements
  - Prompt before each modifying test
  - Run cleanup tests at end
  - Generate comprehensive report

Please confirm cluster is non-production before starting.
```

**Example (Review)**:
```
Review the execution approach for examples/camo/camo-testplan.json

Use testplan-executor skill to analyze:
  - Test dependencies
  - Safety classifications
  - Cleanup coverage
  - Execution order
  - Pre-flight checklist

Output: examples/camo/EXECUTION-REVIEW.md

Do NOT execute tests.
```

**Example (Read-Only Execution)**:
```
Execute read-only tests from examples/camo/camo-testplan.json

Environment:
  - Cluster: ROSA HCP test cluster (non-production)
  - Access: KUBECONFIG=/path/to/test-cluster.kubeconfig
  - Namespace: openshift-configure-alertmanager-operator

Mode: read-only

Please:
  1. Verify cluster access
  2. Run only tests where system_impact.type == "read-only"
  3. Generate validation report showing installation health
  4. Highlight any issues found

This is safe to run on the test cluster.
```

### Pattern 4: Update Test Plan

**When**: Code changed, bugs fixed, features added

**Prompt Template**:
```
Update examples/<SYSTEM_NAME>/<SYSTEM_NAME>-testplan.json

Changes:
  <Describe what changed>

Updates needed:
  1. <Add new tests, update commands, fix errors, etc.>
  2. <Update learning notes, references, etc.>
  3. <Regenerate HTML>

Please:
  - Maintain safety classifications
  - Update learning path if needed
  - Test JSON parses correctly
  - Document in CHANGELOG
```

**Example**:
```
Update examples/camo/camo-testplan.json

Changes:
  - CAMO v2.0 now uses different secret format
  - New CLI flag: --validate-only
  - Bug fixed: webhook timeout is now configurable

Updates needed:
  1. Update test_create_pagerduty_secret to use new secret format
  2. Add test_validate_configuration_only (uses --validate-only flag)
  3. Update test_configure_webhook to show timeout configuration
  4. Update expected outputs where format changed

Please update the test plan, regenerate HTML, and document changes in CHANGELOG.md.
```

### Pattern 5: Implement Repository Improvements

**When**: After creating test plan, you identified improvements to the tooling

**Prompt Template**:
```
Implement improvements to testplan_tools_poc repository.

Based on: examples/<SYSTEM_NAME>/IMPROVEMENTS.md

Priority: <Critical | High | Medium>

Improvements to implement:
  1. <Improvement description>
  2. <Improvement description>

Please:
  - Update Go code, skills, or templates
  - Test with existing test plans
  - Regenerate all example HTMLs to verify
  - Update documentation
  - Mark improvements as completed in REPOSITORY-IMPROVEMENTS.md
```

**Example**:
```
Implement critical improvements from examples/rhobs-next/IMPROVEMENTS.md

Priority: Critical (blocks test plan generation)

Improvements:
  1. Add schema validation to testplan-generator skill
  2. Auto-generate cleanup tests for modifies-state tests
  3. Validate alternative_paths include dependencies

Please:
  - Update skills/testplan-generator/skill.md
  - Update internal/parser/validator.go
  - Add schema generation: go install jsonschema, generate schema
  - Test with camo and rhobs-next examples
  - Mark as ✅ Complete in REPOSITORY-IMPROVEMENTS.md
```

---

## Skill Self-Verification

**All skills auto-verify themselves against the repository** to ensure they're using the latest definitions.

When you invoke a skill, it will:
1. Check if local skill definition has changes
2. Compare with committed version
3. Prompt if differences found

**Prompt example**:
```
⚠️ Skill Definition Verification

I've detected differences in the testplan-generator skill definition:
  - local has uncommitted changes in Phase 6 (safety analysis)
  - remote has updates to schema validation

Options:
1. Continue with current version (may be outdated)
2. Read latest version from repository and use that
3. Show me the full diff to review
4. Cancel and let me update the skill first

What would you like to do?
```

**Best practice**: Choose option 2 (use latest from repository)

**To skip verification** (if you're developing/testing skills):
```
"skip verification" - Use the testplan-generator skill to...
```

---

## Advanced Usage

### Customizing Skills with Parameters

Skills accept parameters to control behavior:

```
# Generate only test cases, reuse existing concepts
Use testplan-generator skill with --tests-only to generate
additional tests for examples/my-operator/my-operator-testplan.json

# Add only hyperlinks, don't modify other content
Use testplan-educator skill with --hyperlinks-only to enhance
examples/my-operator/my-operator-testplan.json

# Dry-run execution review without accessing cluster
Use testplan-executor skill in --dry-run mode to review
examples/my-operator/my-operator-testplan.json execution approach
```

### Chaining Skills

```
Please create a complete test plan for ~/sandbox/my-operator:

1. testplan-generator: Analyze repo and generate test plan JSON
2. testplan-educator: Enhance with educational content and hyperlinks
3. testplan-viewer: Generate interactive HTML
4. testplan-executor --dry-run: Review execution approach
5. Document improvements to repository

Output: examples/my-operator/
```

### Providing Credentials Securely

**Option 1: Environment file** (recommended)
```bash
# Create credentials file (add to .gitignore)
cat > test-credentials.env <<EOF
export KUBECONFIG=/path/to/test-cluster.kubeconfig
export RHOBS_API_URL=https://api-stage.example.com
export ACCESS_TOKEN=<token>
EOF

# Source before using Claude Code
source test-credentials.env
```

**Option 2: Secure credential storage**
```bash
# Store in secure location
mkdir -p ~/.config/testplan-tools/
chmod 700 ~/.config/testplan-tools/
echo "export ACCESS_TOKEN=..." > ~/.config/testplan-tools/credentials
chmod 600 ~/.config/testplan-tools/credentials

# Reference in prompt
```

```
Execute tests with credentials from ~/.config/testplan-tools/credentials

Please:
  - Source the credentials file
  - Verify access to cluster
  - Run read-only tests
  - Do NOT log credentials in output
```

---

## Troubleshooting

### Skill Not Found

**Symptom**: "Unknown skill: testplan-generator"

**Solution**:
```bash
# Verify skills directory exists
ls -la skills/

# Verify you're in the repository directory
pwd
# Should be: /path/to/testplan_tools_poc

# Restart Claude Code from repository directory
cd /path/to/testplan_tools_poc
claude-code
```

### Skill Definition Outdated

**Symptom**: Generated test plan doesn't match expected structure

**Solution**:
```bash
# Pull latest changes
git pull origin main

# Let skill auto-verify (choose option 2: use latest)
# OR manually update
cd skills/testplan-generator/
git diff skill.md  # Review changes
git checkout skill.md  # Revert to repository version
```

### Generated JSON Doesn't Parse

**Symptom**: `Error: failed to parse test plan`

**Solution**:
```
The generated test plan at examples/my-operator/my-operator-testplan.json
failed to parse.

Please:
  1. Read internal/models/testplan.go to see expected structure
  2. Validate JSON structure with jq
  3. Fix structure mismatches (prerequisites, references, etc.)
  4. Regenerate HTML to verify
  5. Document issue in IMPROVEMENTS.md as a schema validation gap
```

### Skills Keep Prompting for Verification

**Symptom**: Every skill use asks to verify against repository

**Cause**: Local changes to skill.md files

**Solution**:
```bash
# Check status
git status skills/

# If you want to keep local changes
git stash
git pull
git stash pop

# If you want repository version
git checkout -- skills/
git pull
```

---

## Examples

See `examples/` directory for complete examples:

- **camo/** - Configure AlertManager Operator
- **rhobs-next/** - RHOBS Next Synthetic Monitoring (comprehensive example)
- **rhobs/** - Original RHOBS examples (v1 format)

Each example includes:
- Test plan JSON
- Interactive HTML viewer
- README with creation notes
- EXECUTION-REVIEW for safety analysis
- IMPROVEMENTS documenting findings

---

## Contributing

Found a bug or have an improvement?

1. Document in `REPOSITORY-IMPROVEMENTS.md`
2. Prioritize: Critical → High → Medium
3. Ask Claude to implement
4. Test with existing examples
5. Submit PR if contributing back

---

## Questions?

- Read skill.md files in each skill directory for detailed methodology
- Check examples/README.md for usage patterns
- Review existing examples in examples/
- See docs/ for detailed documentation

**Quick Help**:
```
In Claude Code, ask:
"How do I use the testplan-generator skill?"
"Show me an example of creating a test plan"
"What artifacts do I need to provide?"
```

Happy test planning! 🧪
