# Test Plan Examples

This directory contains example test plans demonstrating the complete workflow from code analysis to interactive HTML test plans.

## 🚀 Quick Start

**New to these skills?**
1. Read [skills/README.md](../skills/README.md) - Installation and skill usage
2. Review an example below - See what's possible
3. Follow "Creating New Test Plans" section - Step-by-step guide

**Want to create a test plan right now?**
```
# In Claude Code, paste this template and fill in <YOUR_INFO>:

Create a comprehensive test plan for <SYSTEM_NAME>.

Primary artifact: Git repository at <PATH>
Supporting artifacts:
  - Documentation: <PATHS> (README.md, docs/*.md)
  - Existing tests: <PATHS> (test/e2e/, test/integration/)
  - Examples: <PATHS> (config/samples/, examples/)

Context:
  - Target audience: <WHO_USES_THIS> (SRE, QE, developers, etc.)
  - Environment types: <ENVIRONMENTS> (AWS, GCP, HCP, Classic, etc.)
  - Coverage focus: <WHAT_TO_TEST> (installation, config, monitoring, etc.)

Output: examples/<SYSTEM_NAME>/

Please analyze artifacts and create comprehensive, educational test plan.
```

**Want to run tests against a cluster?**
→ See "Running Tests Against a Test Cluster" section below

---

## Directory Structure

Each test plan is organized in its own subdirectory:

```
examples/
├── README.md                    # This file - usage guide
├── camo/                        # Configure AlertManager Operator (CAMO)
│   ├── camo-testplan.json      # Test plan JSON
│   ├── camo-testplan.html      # Interactive HTML viewer
│   ├── README.md               # Example-specific documentation
│   ├── EXECUTION-REVIEW.md     # Execution safety analysis
│   └── SUMMARY.md              # Creation summary and lessons learned
├── rhobs-next/                  # RHOBS Next Synthetic Monitoring
│   ├── rhobs-next-testplan.json
│   ├── rhobs-next-testplan.html
│   ├── README.md
│   ├── EXECUTION-REVIEW.md
│   ├── IMPROVEMENTS.md          # Improvements identified during creation
│   ├── ENHANCEMENT_SUMMARY.md   # Educational enhancement details
│   └── STRUCTURE_FIXES.md       # JSON structure alignment notes
└── rhobs/                       # Original RHOBS examples (v1 format)
    ├── rhobs_test_plan_v2.json  # Learning-optimized test plan
    ├── rhobs_test_plan.json     # Original test plan structure
    └── RHOBS_Manual_Test_Plan.md # Markdown version
```

---

## Creating New Test Plans with Claude

### Overview

Use the Claude skills defined in `skills/` to create comprehensive test plans from various source artifacts. The typical workflow is:

1. **testplan-generator** - Analyze artifacts and generate test plan JSON
2. **testplan-educator** - Enhance with educational content and hyperlinks
3. **testplan-viewer** - Generate interactive HTML
4. **testplan-executor** - Review execution approach (or execute tests)

**📚 First Time Using Skills?** → See [skills/README.md](../skills/README.md) for installation and usage guide

### Understanding Source Artifacts

**The quality and completeness of your test plan directly depends on the artifacts you provide.** Think of Claude as a test architect who needs to understand your system - the more context and examples you provide, the better the test plan.

#### How Claude Guides You

When you start creating a test plan, Claude will:

**1. Identify what you've provided**:
```
Claude: "I see you've provided:
  ✓ Git repository at ~/sandbox/my-operator
  ✓ README.md with basic usage

  To create a comprehensive, educational test plan, I recommend also providing:
  ⚠ Existing E2E tests (if available) - helps understand test patterns
  ⚠ Architecture documentation - needed for conceptual overviews
  ⚠ Known issues or bug reports - creates better troubleshooting content

  Should I proceed with what's available, or would you like to provide additional artifacts?"
```

**2. Ask clarifying questions**:
```
Claude: "I found both installation instructions and upgrade procedures.

  Which should I prioritize for this test plan?
  A) Installation (for new users)
  B) Upgrades (for existing users)
  C) Both with equal weight

  This helps me allocate test coverage appropriately."
```

**3. Suggest missing artifacts during generation**:
```
Claude: "While analyzing the controller code, I noticed references to webhook configuration,
but I couldn't find webhook examples in config/samples/.

Should I:
  A) Skip webhook tests for now
  B) Generate webhook tests based on code analysis (less accurate)
  C) Wait while you provide webhook configuration examples (recommended)
```

**4. Adapt to what's available**:
```
Claude: "I don't see E2E tests in this repository, so I'll base test generation on:
  - Controller logic and reconciliation loops
  - API definitions (CRDs)
  - Example configurations
  - README instructions

  The test plan will be comprehensive but may benefit from later enhancement
  if E2E tests become available."
```

### Common Source Artifacts

Claude can build test plans from any combination of:

#### Primary Focus (Choose One)
- **Git Repository** - Operator code, application source, infrastructure as code
- **Documentation** - User guides, architecture docs, API specifications
- **Existing Tests** - Unit tests, E2E tests, integration tests (any format)
- **Product Documentation** - Official docs, runbooks, SOPs

#### Supporting Artifacts (Optional but Recommended)
- Related code repositories
- OpenAPI/Swagger specifications
- Architecture diagrams (provide URLs)
- Configuration examples (YAML, JSON, etc.)
- Bug reports and known issues (JIRA, GitHub Issues)
- Deployment manifests (Helm charts, operators, etc.)
- Monitoring dashboards (Grafana, Prometheus)
- Incident postmortems

### Step-by-Step: Creating a Test Plan

#### 1. Prepare Your Artifacts

**For a Git Repository** (most common):
```bash
# Clone or update the repository
cd ~/sandbox
git clone https://github.com/org/my-operator
cd my-operator
git pull origin main

# Optionally gather related repos
git clone https://github.com/org/related-component
```

**Document the sources**:
Create a brief summary of what you want tested:
- Primary artifact: `/Users/maclark/sandbox/my-operator` (git repo)
- Target: "Comprehensive E2E test plan for my-operator installation and operations"
- Focus areas: "Installation, configuration, monitoring, troubleshooting"
- Audience: "SRE engineers and QE team"

#### 2. Ask Claude to Generate the Test Plan

**💡 TIP**: Be specific about what you provide. Claude will guide you if more artifacts are needed.

**Example prompt**:
```
Please create a comprehensive test plan for the my-operator at ~/sandbox/my-operator.

Primary artifact: Git repository at ~/sandbox/my-operator
Supporting artifacts:
  - Documentation at docs/
  - E2E tests at test/e2e/
  - README.md

Target audience: SRE engineers new to this operator
Focus: Installation, configuration, and troubleshooting
Output: examples/my-operator/

Please:
1. Analyze the codebase, tests, and documentation
2. Generate comprehensive test cases with safety classifications
3. Add educational content with hyperlinks for technical terms
4. Generate the HTML viewer
5. Review the execution approach (but don't run tests yet)
6. Document any improvements to the testplan_tools_poc project
```

Claude will follow the workflow defined in the skills automatically.

#### 3. Review the Generated Test Plan

```bash
# View the HTML
open examples/my-operator/my-operator-testplan.html

# Review execution safety
cat examples/my-operator/EXECUTION-REVIEW.md

# Check for improvements
cat examples/my-operator/IMPROVEMENTS.md
```

#### 4. Iterate and Improve

Ask Claude to:
- Add more tests for specific scenarios: "Add tests for upgrade paths"
- Enhance educational content: "Add more troubleshooting guidance for common errors"
- Fix identified issues: "Implement the critical improvements from IMPROVEMENTS.md"
- Update the testplan_tools_poc project: "Add the missing cleanup test validation to the parser"

---

## Running Tests Against a Test Cluster

### Prerequisites

#### 1. Set Up Cluster Access

**For OpenShift clusters**:
```bash
# Set KUBECONFIG to point to your test cluster
export KUBECONFIG=/path/to/test-cluster-kubeconfig

# Verify access
oc whoami
oc get nodes

# Confirm it's NOT a production cluster
oc get infrastructure cluster -o jsonpath='{.status.infrastructureName}'
# Should NOT contain "prod" or "production"
```

**For standard Kubernetes clusters**:
```bash
# Set KUBECONFIG
export KUBECONFIG=/path/to/kubeconfig

# Verify access
kubectl cluster-info
kubectl get nodes

# Confirm cluster identity
kubectl config current-context
```

#### 2. Set Required Environment Variables

Different test plans may require different credentials. Common variables:

```bash
# Basic cluster access
export KUBECONFIG=/path/to/kubeconfig

# API endpoints (for RHOBS, other APIs)
export RHOBS_API_URL="https://rhobs-api.example.com"
export ACCESS_TOKEN="your-oidc-token"

# Cloud provider credentials (if needed for cluster operations)
export AWS_ACCESS_KEY_ID="..."
export AWS_SECRET_ACCESS_KEY="..."

# OCM credentials (for ROSA/OSD clusters)
export OCM_TOKEN="$(ocm token)"

# Other operator-specific variables (check test plan prerequisites)
export OPERATOR_NAMESPACE="openshift-operators"
```

**Store in a credentials file** (gitignored):
```bash
# Create credentials file
cat > test-credentials.env <<EOF
export KUBECONFIG=/path/to/test-cluster.kubeconfig
export RHOBS_API_URL=https://rhobs-api-stage.example.com
export ACCESS_TOKEN=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
EOF

# Source before running tests
source test-credentials.env
```

#### 3. Verify Prerequisites

Each test plan has a prerequisites section. Verify all requirements:

```bash
# Check tools are installed
oc version
kubectl version
jq --version
curl --version

# Verify cluster access
oc auth can-i create deployment
oc get all -n openshift-operators

# Check operator is deployed (if applicable)
oc get deployment my-operator -n openshift-operators
```

### Execution Modes

#### Read-Only Mode (Safe for Any Cluster)

Read-only tests query existing resources without making changes. They provide **runtime validation** of an existing installation.

**Use cases**:
- Verify operator is installed and healthy
- Validate configuration is correct
- Check monitoring is working
- Troubleshoot issues in production (safe to run)
- Smoke test after deployment

**How to run read-only tests only**:

**Option 1: Use testplan-executor skill**
```
Please run only the read-only tests from examples/my-operator/my-operator-testplan.json against my test cluster.

Cluster: export KUBECONFIG=/path/to/test-cluster.kubeconfig
Credentials: (provide any API tokens needed)

Only run tests where system_impact.type == "read-only"
Generate a validation report showing the health of the installation.
```

**Option 2: Filter manually** (future tool will automate this)
```bash
# Extract read-only test IDs
jq -r '.testcases | to_entries | .[] | select(.value.metadata.system_impact.type == "read-only") | .key' \
  examples/my-operator/my-operator-testplan.json

# Run specific tests
# (Execution automation tool coming soon)
```

**Read-only validation benefits**:
- ✅ Safe to run in production
- ✅ No cleanup required
- ✅ Can run repeatedly without side effects
- ✅ Idempotent
- ✅ Provides health check
- ✅ Helps troubleshoot issues

#### Full Test Mode (Modifies State - Test Clusters Only)

Full test mode includes modifying tests that create, update, or delete resources.

**⚠️ CRITICAL**: Only run in test/dev environments, NEVER in production!

**Safety checklist before running**:
```bash
# 1. Confirm cluster is NOT production
CLUSTER_NAME=$(oc get infrastructure cluster -o jsonpath='{.status.infrastructureName}')
if [[ "$CLUSTER_NAME" =~ prod|production ]]; then
  echo "❌ PRODUCTION CLUSTER DETECTED - ABORT"
  exit 1
fi

# 2. Review execution plan
cat examples/my-operator/EXECUTION-REVIEW.md

# 3. Identify tests that need cleanup
jq -r '.testcases | to_entries | .[] | select(.value.metadata.system_impact.type == "modifies-state") | "\(.key) → cleanup: \(.value.metadata.state_management.cleanup_test)"' \
  examples/my-operator/my-operator-testplan.json

# 4. Confirm all cleanup tests exist
# (Check output above - all should have cleanup test)

# 5. Ask Claude to run tests with cleanup tracking
```

**Running full tests with Claude**:
```
Please execute the full test plan at examples/my-operator/my-operator-testplan.json.

Environment:
  - Cluster: TEST cluster (verified non-production)
  - KUBECONFIG: /path/to/test-cluster.kubeconfig
  - Credentials: (provided separately)

Requirements:
  1. Run tests in dependency order
  2. Track which cleanup tests are needed
  3. Prompt before each modifying test
  4. Collect evidence on failures
  5. Run all required cleanup tests at the end
  6. Generate comprehensive test report

Please confirm the cluster is non-production before starting.
```

#### Dry-Run Mode (Simulate Execution)

Review what would happen without actually running commands:

```
Please perform a dry-run of examples/my-operator/my-operator-testplan.json.

Show:
  - Execution order
  - Which tests would run
  - Safety confirmations that would be required
  - Cleanup operations that would be needed
  - Estimated duration

Do NOT execute any commands against the cluster.
```

---

## Runtime Validation Use Case

Read-only tests serve as **operational validation tools** for installed operators/applications.

### Example: Validating CAMO Installation

After deploying Configure AlertManager Operator (CAMO), run read-only tests to verify:

1. **Operator Health**
   - Deployment exists and is ready
   - Pods are running
   - RBAC permissions are correct
   - No error logs in recent output

2. **Configuration Validation**
   - ServiceMonitors are created
   - Prometheus is scraping metrics
   - AlertManager config is valid
   - Integration secrets exist

3. **Metrics Pipeline**
   - Metrics are being collected
   - Queries return expected results
   - No gaps in time series

**When to run**:
- ✅ After initial installation (smoke test)
- ✅ After upgrades (regression check)
- ✅ Periodically (health monitoring)
- ✅ During troubleshooting (diagnostic)
- ✅ In CI/CD pipelines (integration test)
- ✅ Before planned maintenance (baseline)

**How it helps troubleshooting**:

**Scenario**: Alerts aren't firing in AlertManager

**Run validation tests**:
```
Please run the read-only validation tests from examples/camo/camo-testplan.json.

Focus on tests that validate:
  - AlertManager configuration
  - ServiceMonitor creation
  - Prometheus scraping
  - Metric collection

Generate a report highlighting what's working vs broken.
```

**Claude will identify**:
- ✅ CAMO operator is running
- ✅ ServiceMonitor exists
- ❌ Prometheus NOT scraping (no metrics)
- Root cause: ServiceMonitor selector doesn't match service labels

**Immediate fix provided**: Update ServiceMonitor selector to match service

---

## Updating Existing Test Plans

As code changes, documentation evolves, and bugs are fixed, test plans need updates.

### Adding New Tests

**Scenario**: New feature added to operator

```
The my-operator now supports automatic scaling. Please add tests to examples/my-operator/my-operator-testplan.json:

New tests needed:
  1. test_enable_autoscaling - Enable autoscaling on deployment
  2. test_verify_autoscaling - Verify HPA is created and working
  3. test_cleanup_autoscaling - Remove autoscaling configuration

Source: New code in controllers/autoscaler_controller.go
Documentation: docs/autoscaling.md

Please:
  - Generate the 3 new tests with safety classifications
  - Add them to the learning path
  - Update the HTML
  - Review execution approach
```

### Updating Test Content

**Scenario**: Command syntax changed in new operator version

```
The my-operator CLI changed in v2.0:
  - Old: myop create --type=foo
  - New: myop create foo

Please update all tests in examples/my-operator/my-operator-testplan.json to use the new CLI syntax.

Also update:
  - Expected outputs
  - Common errors (old syntax is now an error)
  - Learning notes to mention the v2.0 change
```

### Fixing Issues

**Scenario**: Test plan has gaps identified in IMPROVEMENTS.md

```
Please implement the critical improvements from examples/my-operator/IMPROVEMENTS.md:

Priority 1: Add missing cleanup test for test_create_custom_resource
Priority 2: Fix alternative_paths to include dependencies

Update the test plan JSON and regenerate HTML.
```

### Enhancing Educational Content

```
Please enhance the educational content in examples/my-operator/my-operator-testplan.json:

Improvements:
  1. Add hyperlinks for these technical terms: CustomResourceDefinition, Webhook, Finalizer
  2. Add conceptual_overview for advanced tests
  3. Expand common_errors for test_install_operator (add timeout scenario)
  4. Add troubleshooting section for RBAC permission errors

Regenerate HTML when done.
```

---

## Improving the testplan_tools_poc Project

Each time you create a test plan, you'll discover improvements to the tooling itself. Claude can help implement these.

### Process

1. **Identify Improvements** - Documented in `IMPROVEMENTS.md` for each example
2. **Prioritize** - Critical → High → Medium
3. **Implement** - Ask Claude to update the repo
4. **Test** - Verify with existing test plans
5. **Document** - Update repo README

### Example Improvement Workflow

**Scenario**: After creating 3 test plans, you notice the testplan-generator always creates invalid prerequisites structure.

**Step 1: Document the Issue**
```
Issue: testplan-generator creates prerequisites.required_access as array of strings,
but the Go model expects array of RequiredAccess objects.

Impact: All generated test plans fail to parse until manually fixed.
Occurred in: camo, rhobs-next examples

Solution: Update testplan-generator skill to follow Go model structure.
```

**Step 2: Ask Claude to Fix**
```
Please fix the testplan-generator skill to generate prerequisites that match the Go model in internal/models/testplan.go.

The skill currently generates:
  "required_access": ["string1", "string2"]

But the model expects:
  "required_access": [
    {"name": "...", "purpose": "...", "how_to_get": "...", "verification": "..."}
  ]

Please:
  1. Update skills/testplan-generator/skill.md Phase 3 to use correct structure
  2. Add schema validation step to detect this early
  3. Test with a simple test plan to verify it works
  4. Update the skill documentation

Also check and fix:
  - required_knowledge structure
  - environment_setup.tools structure
```

**Step 3: Verify the Fix**
```
Please regenerate the camo test plan prerequisites using the updated skill and verify:
  1. JSON structure matches Go model
  2. Parser accepts it without errors
  3. HTML renders prerequisites correctly
```

**Step 4: Document in Repo**
```
Please update REPOSITORY-IMPROVEMENTS.md to mark improvement #1 as ✅ Complete and add implementation notes.
```

### Common Improvements to Track

As you create more test plans, ask Claude to implement:

**Schema & Validation**:
- ✅ Add JSON schema validation to testplan-generator
- ✅ Validate alternative_paths include dependencies
- ✅ Check cleanup tests exist for modifying tests

**Safety & Execution**:
- ✅ Add production environment detection
- ✅ Implement cleanup tracking system
- ✅ Auto-generate cleanup tests

**Educational Features**:
- ✅ Render markdown hyperlinks in HTML
- ✅ Support conceptual_overview field in Go model
- ✅ Add inline glossary for technical terms

**Usability**:
- ✅ Create test filtering tool (read-only vs modifying)
- ✅ Add dry-run mode to testplan-executor
- ✅ Generate test execution reports

**Learning from Each Test Plan**:

After creating each new test plan, ask:
```
What improvements to testplan_tools_poc did we identify while creating this test plan?

Please:
  1. Review the IMPROVEMENTS.md file
  2. Identify 3 highest-priority improvements
  3. Implement them if straightforward
  4. Document if complex (for later implementation)

Focus on changes that will benefit ALL future test plans.
```

---

## Future Tool: Test Filtering for Validation

**Coming Soon**: A CLI tool to filter and run subsets of tests.

**Planned Usage**:
```bash
# Run only read-only tests (validation mode)
./build/testplan-runner --filter read-only \
  --testplan examples/my-operator/my-operator-testplan.json

# Run only tests for specific component
./build/testplan-runner --filter category=installation \
  --testplan examples/my-operator/my-operator-testplan.json

# Run troubleshooting diagnostics
./build/testplan-runner --filter difficulty=beginner,read-only \
  --testplan examples/my-operator/my-operator-testplan.json \
  --report diagnostics.html

# Run full test suite with cleanup tracking
./build/testplan-runner --mode full \
  --testplan examples/my-operator/my-operator-testplan.json \
  --cleanup-tracking
```

**Until this tool exists**, ask Claude:
```
Please extract and run only the read-only tests from examples/my-operator/my-operator-testplan.json.

Filter criteria: system_impact.type == "read-only"
```

---

## Best Practices

### When Creating Test Plans

1. **Start with Read-Only Tests** - Build confidence before modifying state
2. **Classify Safety Early** - Every test should have system_impact classification
3. **Always Create Cleanup Tests** - For every modifying test, create matching cleanup
4. **Use Real Clusters for Examples** - Test against actual installations when possible
5. **Document as You Go** - Capture IMPROVEMENTS.md during creation
6. **Iterate** - First pass gets 70% right, iterate to 100%

### When Running Tests

1. **Verify Cluster Identity** - ALWAYS confirm non-production before modifying tests
2. **Read EXECUTION-REVIEW.md** - Understand safety before running
3. **Start with Read-Only** - Validate environment before modifying
4. **Track Cleanup** - Note which cleanup tests are needed
5. **Collect Evidence** - Save logs and outputs on failures
6. **Generate Reports** - Document what passed/failed

### When Updating Test Plans

1. **Test After Updates** - Regenerate HTML and verify parsing
2. **Maintain Safety Classifications** - Don't remove safety attributes
3. **Update Learning Path** - Add new tests to sequence
4. **Version Test Plans** - Consider test-plan-v2.json for major changes
5. **Document Breaking Changes** - Note CLI syntax changes, deprecated features

### When Improving the Repository

1. **Fix Critical First** - Blocking issues before nice-to-haves
2. **Test with All Examples** - Verify changes don't break existing test plans
3. **Update Skills Documentation** - Keep skill.md files current
4. **Schema-First** - Define data structures in Go model, generate JSON schema
5. **Backward Compatible** - Don't break existing test plans if possible

---

## Example Workflow Summary

**Creating a new test plan for `my-operator`**:

```bash
# 1. Prepare repository
cd ~/sandbox
git clone https://github.com/org/my-operator
cd my-operator
git pull

# 2. Ask Claude (using Claude Code CLI or Web)
```

**Prompt**:
```
Create comprehensive test plan for ~/sandbox/my-operator

Sources: Git repo, docs/, test/e2e/, README.md
Output: examples/my-operator/
Audience: SRE engineers
Focus: Installation, config, monitoring, troubleshooting

Workflow:
  1. Analyze artifacts (testplan-generator)
  2. Enhance education (testplan-educator)
  3. Generate HTML (testplan-viewer)
  4. Review execution (testplan-executor - don't run yet)
  5. Document improvements
```

```bash
# 3. Review output
open examples/my-operator/my-operator-testplan.html
cat examples/my-operator/EXECUTION-REVIEW.md

# 4. Run read-only validation
export KUBECONFIG=/path/to/test-cluster.kubeconfig
```

**Prompt**:
```
Run read-only tests from examples/my-operator/my-operator-testplan.json
Generate validation report for the installation
```

```bash
# 5. Implement improvements
```

**Prompt**:
```
Implement critical improvements from examples/my-operator/IMPROVEMENTS.md
Focus on schema validation and cleanup test generation
```

```bash
# 6. Commit and document
git add examples/my-operator/
git commit -m "Add my-operator test plan"
```

---

## Questions?

For help:
- Read example READMEs in each subdirectory
- Review skills documentation in `skills/`
- Check `docs/` for detailed guides
- See `REPOSITORY-IMPROVEMENTS.md` for roadmap

Happy testing! 🧪
