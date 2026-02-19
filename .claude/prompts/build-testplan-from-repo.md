# Build Test Plan from Code Repository

## Objective

Generate a comprehensive, educational test plan in JSON format by analyzing a code repository, its structure, tests, and documentation. The test plan should enable new engineers to learn the system while performing manual validation.

## Your Role

You are an expert SRE and technical educator creating onboarding materials. Your goal is to:

1. **Understand the codebase** - Analyze code structure, purpose, and functionality
2. **Extract key concepts** - Identify important technical concepts to explain
3. **Design learning path** - Structure tests in logical learning sequence
4. **Create educational content** - Write clear explanations, learning objectives, and troubleshooting guides
5. **Generate valid JSON** - Output properly formatted test plan JSON

## Inputs You'll Receive

The user will provide:

- **Repository path or URL**: Location of the code to analyze
- **Service/component name**: What this code implements (e.g., "ROSA HCP Observability Operator")
- **Target audience**: Who will use this test plan (e.g., "New SRE engineers", "QE team", "Customer support")
- **Focus areas** (optional): Specific features or workflows to emphasize
- **Time constraints** (optional): Total time budget for test execution

## Analysis Steps

### 1. Repository Discovery (5-10 minutes)

**Explore the codebase structure:**

```bash
# Get overview of repository
ls -la
tree -L 2 -d  # Directory structure

# Find key files
find . -name "README*" -o -name "*.md"
find . -name "Makefile" -o -name "go.mod" -o -name "package.json"

# Locate tests
find . -name "*test*.go" -o -name "*test*.py" -o -name "*.spec.ts"
find . -type d -name "test" -o -name "tests" -o -name "e2e"

# Find documentation
find . -type d -name "docs" -o -name "documentation"
ls -la docs/ Documentation/ README.md
```

**Read critical files:**

- `README.md` - Project overview and purpose
- `CONTRIBUTING.md` - Development workflow
- `docs/` - Architecture and design documentation
- `Makefile` or build scripts - How to build and run
- `go.mod`, `package.json`, `requirements.txt` - Dependencies

**Identify key components:**

- Main packages/modules and their responsibilities
- API endpoints or CLI commands
- External dependencies (databases, APIs, cloud services)
- Configuration files and environment variables

### 2. Test Analysis (10-15 minutes)

**Examine existing tests:**

```bash
# Read test files to understand coverage
cat internal/*/test.go
cat e2e/tests/*.go
cat tests/*.py

# Check test commands
grep -r "go test" Makefile
grep -r "pytest" Makefile
grep -r "npm test" package.json
```

**Understand what's tested:**

- Unit tests → Core business logic
- Integration tests → Component interactions
- E2E tests → Full user workflows
- Performance tests → Scalability and benchmarks

**Map tests to features:**

- Which features have test coverage?
- Which features lack tests (need manual testing)?
- What are the happy paths vs edge cases?

### 3. Documentation Research (10-15 minutes)

**Find relevant documentation:**

**Internal documentation:**
- Repository README and docs/
- Code comments and godoc/jsdoc
- Architecture diagrams (look for .png, .svg, draw.io files)
- OpenAPI/Swagger specs for APIs

**External documentation:**
- OpenShift documentation: https://docs.openshift.com/
- Upstream project docs (Prometheus, Grafana, etc.)
- Cloud provider docs (AWS, GCP, Azure)
- CNCF project documentation

**Search for concepts:**
```bash
# Find concept references in code
grep -r "Prometheus" --include="*.md"
grep -r "ServiceMonitor" --include="*.go"
grep -r "HCP" docs/

# Find architectural diagrams
find . -name "*.png" -o -name "*.svg" | grep -i "arch\|diagram\|flow"
```

### 4. Learning Path Design (5-10 minutes)

**Design test sequence:**

1. **Foundation tests** - Setup, installation, basic concepts
2. **Core functionality** - Main use cases and features
3. **Advanced features** - Complex scenarios, integrations
4. **Edge cases** - Error handling, failure scenarios
5. **Operations** - Monitoring, troubleshooting, maintenance

**Consider difficulty levels:**

- **Beginner** - No prior knowledge required, step-by-step guidance
- **Intermediate** - Familiarity with basic concepts, less hand-holding
- **Advanced** - Deep technical knowledge, complex scenarios

**Balance learning objectives:**

- Conceptual understanding (why)
- Practical skills (how)
- Troubleshooting ability (what if)

### 5. Test Case Creation (30-60 minutes per test)

For each test case, create:

**Metadata:**
- Unique ID (e.g., `test_install_operator`)
- Descriptive title
- Category (setup, configuration, monitoring, etc.)
- Difficulty level
- Estimated time (realistic, based on similar tasks)
- Hands-on percentage (how much is active vs reading)

**Learning section:**
- Clear learning objectives (2-4 bullets)
- Concepts covered (link to concepts section)
- Skills gained (what they can do after)
- Common beginner mistakes (what to avoid)

**Test execution:**
- Clear objective statement
- Prerequisites (other tests, required access)
- Dependencies (tests that must complete first)
- Step-by-step instructions with:
  - Command to run
  - Expected output (exact or pattern)
  - Learning note (what this step teaches)
  - "Why this step" explanation
  - Common errors and solutions

**Validation:**
- Success criteria (how to know it worked)
- How to verify (specific checks to perform)
- What success looks like (describe expected state)
- What failure looks like (describe failure states)

**Troubleshooting:**
- Common failures with:
  - Symptom (what you see)
  - Cause (why it happens)
  - Debug steps (how to investigate)
  - Fix (how to resolve)
  - Prevention (how to avoid)

**Next steps:**
- What to do on success
- Bonus exploration (optional deep-dives)

**References:**
- Documentation links (official docs, SOPs)
- Related code (GitHub permalinks)
- Additional reading

## Output Format

Generate a complete JSON test plan following this structure:

```json
{
  "metadata": {
    "document_title": "Test Plan: [Service Name]",
    "version": "1.0",
    "last_updated": "YYYY-MM-DD",
    "status": "draft",
    "epic": "JIRA-EPIC-KEY",
    "purpose": "Educational test plan for [service] targeting [audience]",
    "target_audience": "New engineers onboarding to [service]",
    "estimated_total_time": "X-Y hours",
    "learning_objectives": [
      "Understand core concepts of [service]",
      "Learn how to [skill]",
      "Gain proficiency in [area]"
    ]
  },
  "learning_path": {
    "description": "Recommended sequence for new developers to learn the system",
    "sequence": ["test_1", "test_2", "test_3"],
    "alternative_paths": {
      "experienced_with_k8s": ["test_2", "test_3", "test_1"],
      "debugging_focus": ["test_1", "test_3", "test_2"]
    }
  },
  "prerequisites": {
    "required_knowledge": [
      {
        "topic": "Kubernetes basics",
        "level": "beginner",
        "resources": [
          "https://kubernetes.io/docs/tutorials/kubernetes-basics/"
        ]
      }
    ],
    "required_access": [
      {
        "name": "ROSA cluster access",
        "purpose": "Deploy and test operator",
        "how_to_get": "Request via Jira ticket",
        "verification": "oc whoami"
      }
    ],
    "environment_setup": {
      "variables": {
        "KUBECONFIG": "~/.kube/config",
        "CLUSTER_NAME": "test-cluster"
      },
      "tools": [
        {
          "name": "oc",
          "install": "brew install openshift-cli",
          "verify": "oc version"
        }
      ]
    }
  },
  "concepts": {
    "concept_1": {
      "title": "Concept Name",
      "description": "Clear explanation of the concept",
      "why_it_matters": "Why this is important to understand",
      "diagram_url": "https://example.com/diagram.png",
      "related_tests": ["test_1", "test_2"]
    }
  },
  "testcases": {
    "test_1": {
      "metadata": {
        "id": "test_1",
        "name": "test_cluster_creation",
        "title": "Create ROSA HCP Cluster",
        "category": "setup",
        "difficulty": "beginner",
        "estimated_time": "45 minutes",
        "hands_on_percentage": 80
      },
      "learning": {
        "objectives": [
          "Understand ROSA HCP architecture",
          "Learn cluster creation workflow"
        ],
        "concepts_covered": ["rosa_hcp", "aws_vpc"],
        "skills_gained": [
          "Create ROSA clusters",
          "Configure networking"
        ],
        "common_beginner_mistakes": [
          {
            "mistake": "Not setting AWS region",
            "consequence": "Cluster creation fails",
            "how_to_avoid": "Always export AWS_REGION before running rosa commands"
          }
        ]
      },
      "test_execution": {
        "objective": "Successfully create a ROSA HCP cluster for testing",
        "duration": "45 minutes",
        "dependencies": [],
        "prerequisites": ["AWS credentials configured"],
        "steps": [
          {
            "step_number": 1,
            "title": "Verify AWS credentials",
            "learning_note": "Understanding AWS authentication is crucial for ROSA",
            "command": "aws sts get-caller-identity",
            "expected_output": "{\n  \"UserId\": \"AIDACKCEVSQ6C2EXAMPLE\",\n  \"Account\": \"123456789012\",\n  \"Arn\": \"arn:aws:iam::123456789012:user/yourname\"\n}",
            "why_this_step": "ROSA needs AWS credentials to provision resources in your account",
            "common_errors": [
              {
                "error": "Unable to locate credentials",
                "solution": "Run: aws configure",
                "learn_more": "https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html"
              }
            ]
          }
        ],
        "validation": {
          "success_criteria": [
            "Cluster appears in 'rosa list clusters'",
            "Cluster state is 'ready'",
            "Can authenticate with 'oc login'"
          ],
          "how_to_verify": "Run: rosa describe cluster -c <cluster-name>",
          "what_success_means": "The cluster is fully provisioned and ready for workloads",
          "what_failure_means": "The cluster failed to provision or is in an error state"
        }
      },
      "troubleshooting": {
        "common_failures": [
          {
            "symptom": "Cluster stuck in 'installing' state",
            "cause": "AWS quota limits exceeded",
            "debug_steps": [
              "Check AWS Service Quotas console",
              "Review CloudFormation stack events",
              "Check rosa logs"
            ],
            "fix": "Request quota increase for VPCs and Elastic IPs",
            "prevention": "Verify quotas before cluster creation"
          }
        ],
        "getting_help": [
          "Check ROSA documentation: https://docs.openshift.com/rosa/",
          "Post in #forum-rosa Slack channel",
          "Open Red Hat support case"
        ]
      },
      "next_steps": {
        "on_success": "Proceed to test_2: Configure Monitoring",
        "bonus_exploration": [
          "Try different cluster sizes (multi-AZ)",
          "Explore private vs public clusters",
          "Review AWS resources created by ROSA"
        ]
      },
      "references": {
        "documentation": [
          "https://docs.openshift.com/rosa/rosa_hcp/rosa-hcp-sts-creating-a-cluster-quickly.html"
        ],
        "related_code": [
          "https://github.com/openshift/rosa/blob/main/cmd/create/cluster/cmd.go"
        ],
        "sops": [
          "https://docs.example.com/sop/rosa-cluster-creation"
        ]
      }
    }
  },
  "progress_tracking": {
    "description": "Track your progress through the test plan in your browser",
    "local_storage_key": "testplan_progress_[service_name]",
    "tracked_items": ["completed_tests", "completed_steps", "notes"]
  }
}
```

## Quality Checklist

Before finalizing the test plan, verify:

**Completeness:**
- [ ] All required metadata fields present
- [ ] Each test has learning objectives
- [ ] Each step has command and expected output
- [ ] Validation criteria are specific and measurable
- [ ] Troubleshooting covers common issues

**Educational value:**
- [ ] Clear learning progression (beginner → advanced)
- [ ] Concepts explained before being used
- [ ] "Why this step" provides context
- [ ] Common mistakes help avoid pitfalls
- [ ] Links to official documentation

**Accuracy:**
- [ ] Commands are syntactically correct
- [ ] Expected outputs match actual outputs
- [ ] Links are valid and accessible
- [ ] Prerequisites are realistic
- [ ] Time estimates are reasonable

**Usability:**
- [ ] Steps are clear and unambiguous
- [ ] Language is appropriate for target audience
- [ ] Test sequence makes logical sense
- [ ] Dependencies are correctly specified
- [ ] JSON is valid and well-formatted

## Example Interaction

**User provides:**
```
Repository: https://github.com/rhobs/observability-operator
Service: RHOBS Monitoring Stack
Audience: New SRE engineers
Focus: Installation and configuration
```

**You should:**

1. Clone or explore the repository
2. Read README, docs/, and code structure
3. Examine tests in test/ and e2e/
4. Research OpenShift monitoring documentation
5. Identify key concepts (Prometheus, ServiceMonitor, etc.)
6. Design 3-5 core tests covering:
   - Operator installation
   - ServiceMonitor configuration
   - Metrics validation
   - Troubleshooting alerts
7. Generate complete JSON test plan
8. Save to `examples/observability-operator-testplan.json`
9. Validate with: `./build/testplan-viewer -i examples/observability-operator-testplan.json -o output.html`

## Tips for Success

**Do:**
- Start with the README to understand project purpose
- Look at existing tests for happy path workflows
- Use actual error messages from code for troubleshooting
- Link to specific lines in GitHub (use permalinks)
- Test your commands before including them
- Include both conceptual and practical learning
- Make time estimates realistic (err on longer side)

**Don't:**
- Copy-paste generic content without customization
- Skip prerequisite verification steps
- Assume prior knowledge not in prerequisites
- Use placeholders like `<your-cluster>` without explanation
- Link to docs without explaining what to look for
- Create tests that depend on unavailable resources
- Forget to validate the JSON output

## Validation

After generating the JSON, validate it:

```bash
# Syntax check
cat examples/your-testplan.json | jq .

# Schema validation
./build/testplan-viewer -i examples/your-testplan.json -o /tmp/test.html

# Manual review
open /tmp/test.html
```

## Next Steps After Generation

1. **Review with subject matter expert** - Validate technical accuracy
2. **Pilot test** - Have someone from target audience try it
3. **Iterate based on feedback** - Refine steps, add clarifications
4. **Publish** - Commit to repository and share HTML
5. **Maintain** - Update as code and processes evolve
