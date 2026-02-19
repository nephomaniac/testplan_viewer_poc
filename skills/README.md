# Test Plan Skills Suite

Three specialized Claude skills for comprehensive test plan generation, educational enhancement, and systematic execution.

## Overview

This suite provides end-to-end test plan management through three focused skills:

```
┌─────────────────────────────────────────────────────────────────┐
│                        Test Plan Lifecycle                       │
└─────────────────────────────────────────────────────────────────┘

    ┌─────────────────────┐
    │  1. GENERATION      │  → Comprehensive test coverage
    │  testplan-generator │     from artifacts
    └──────────┬──────────┘
               │
               ▼
    ┌─────────────────────┐
    │  2. EDUCATION       │  → Add learning content,
    │  testplan-educator  │     make tests educational
    └──────────┬──────────┘
               │
               ▼
    ┌─────────────────────┐
    │  3. EXECUTION       │  → Run tests, report results,
    │  testplan-executor  │     provide follow-up actions
    └─────────────────────┘
```

## The Three Skills

### 1. testplan-generator

**Focus:** Comprehensive test case generation

**What it does:**
- Analyzes code repositories, tests, and documentation
- Generates comprehensive test coverage
- Creates environment variants (AWS, GCP, Azure)
- Covers input permutations and edge cases
- Identifies coverage gaps

**Use when:**
- Starting a new test plan from scratch
- Expanding coverage to new environments
- Systematically covering all scenarios
- Need gap analysis

**Input:**
- Code repository OR documentation OR existing tests
- Environment types to cover
- Coverage goals

**Output:**
- Comprehensive test plan JSON
- Coverage matrix
- Gap analysis
- Recommendations

[→ Full documentation](testplan-generator/)

---

### 2. testplan-educator

**Focus:** Educational enhancement and UX

**What it does:**
- Adds concept explanations
- Writes learning objectives
- Documents common mistakes
- Provides multi-modal resources
- Creates progressive learning paths

**Use when:**
- Making tests beginner-friendly
- Creating onboarding materials
- Adding documentation to tests
- Improving test educational value
- Teaching through testing

**Input:**
- Test plan JSON
- Target audience
- Learning objectives

**Output:**
- Educationally enhanced test plan JSON
- Concept reference guide
- Learning resource list

[→ Full documentation](testplan-educator/)

---

### 3. testplan-executor

**Focus:** Execution, reporting, and follow-up

**What it does:**
- Executes tests systematically
- Validates results
- Generates comprehensive reports
- Identifies failure patterns
- Provides remediation steps

**Use when:**
- Executing test plans
- Validating test accuracy
- Generating stakeholder reports
- Identifying root causes
- Creating action plans

**Input:**
- Test plan JSON
- Environment context
- Execution scope and mode

**Output:**
- Executive summary
- Detailed technical report
- Follow-up action plan
- Gap analysis
- Test results JSON

[→ Full documentation](testplan-executor/)

## Usage Patterns

### Pattern 1: End-to-End New Test Plan

**Goal:** Create and validate a complete test plan for a new service

```
Step 1: Generate comprehensive coverage
  → Use: testplan-generator
  → Input: Code repository + documentation
  → Output: examples/service-testplan.json

Step 2: Add educational content
  → Use: testplan-educator
  → Input: examples/service-testplan.json
  → Output: examples/service-testplan-educational.json

Step 3: Execute and validate
  → Use: testplan-executor
  → Input: examples/service-testplan-educational.json
  → Output: test-results/YYYYMMDD_HHMMSS/ (reports)

Step 4: Iterate based on findings
  → Update test plan based on execution results
  → Re-execute
```

**Time estimate:** 4-6 hours total

---

### Pattern 2: Quick Coverage Generation

**Goal:** Rapidly generate test coverage without education overhead

```
Step 1: Generate from existing tests
  → Use: testplan-generator
  → Input: E2E test files
  → Output: examples/coverage-testplan.json

Step 2: Execute to validate
  → Use: testplan-executor
  → Input: examples/coverage-testplan.json
  → Output: Gap analysis and recommendations
```

**Time estimate:** 1-2 hours

---

### Pattern 3: Educational Enhancement Only

**Goal:** Make existing tests more beginner-friendly

```
Step 1: Enhance existing test plan
  → Use: testplan-educator
  → Input: examples/existing-testplan.json
  → Output: examples/existing-testplan-educational.json

Step 2: Generate HTML and review
  → make run
  → Review with target audience

Step 3: Iterate based on feedback
```

**Time estimate:** 2-3 hours

---

### Pattern 4: Execution and Remediation

**Goal:** Run tests and fix issues

```
Step 1: Execute test plan
  → Use: testplan-executor
  → Input: examples/testplan.json
  → Output: Reports showing failures

Step 2: Follow remediation steps
  → Use: Follow-up action plan
  → Fix critical blockers

Step 3: Re-execute failed tests
  → Use: testplan-executor --resume-from test_3
  → Validate fixes work

Step 4: Update test plan
  → Incorporate learnings
  → Use testplan-educator to add troubleshooting notes
```

**Time estimate:** 2-4 hours (varies by issues)

---

### Pattern 5: Continuous Validation

**Goal:** Regular test plan execution for CI/CD

```
Scheduled execution (nightly, weekly, etc.):
  → Use: testplan-executor (automated mode)
  → Input: examples/regression-testplan.json
  → Output: Automated reports
  → Alert on failures

On failure:
  → Review detailed report
  → Follow remediation plan
  → Update test plan if needed
```

**Time estimate:** 1-3 hours per run

## Skill Selection Guide

**Choose testplan-generator when you need:**
- ✅ Comprehensive test coverage across environments
- ✅ Systematic coverage of all scenarios
- ✅ Gap identification
- ✅ Environment variant generation
- ✅ Input permutation testing

**Choose testplan-educator when you need:**
- ✅ Beginner-friendly tests
- ✅ Learning while testing
- ✅ Concept explanations
- ✅ Common mistake prevention
- ✅ Multi-modal learning resources

**Choose testplan-executor when you need:**
- ✅ Automated test execution
- ✅ Result validation
- ✅ Stakeholder reports
- ✅ Failure pattern analysis
- ✅ Remediation action plans

**Use all three when:**
- ✅ Creating comprehensive educational test plans
- ✅ Onboarding new team members
- ✅ Building production validation suites
- ✅ Establishing test-driven processes

## Integration with Test Plan Viewer

All skills work with the test plan viewer:

```bash
# After any skill generates/enhances a test plan:
make run

# Or specify the test plan:
./build/testplan-viewer -i examples/your-testplan.json -o output.html

# Open in browser:
open output.html
```

The HTML viewer provides:
- Interactive test execution
- Progress tracking
- Copy-to-clipboard commands
- Collapsible sections
- Search and filter
- Browser-based persistence

## File Organization

```
skills/
├── README.md                    # This file
├── testplan-generator/
│   ├── skill.md                # Generator skill prompt
│   └── README.md               # Generator documentation
├── testplan-educator/
│   ├── skill.md                # Educator skill prompt
│   └── README.md               # Educator documentation
└── testplan-executor/
    ├── skill.md                # Executor skill prompt
    └── README.md               # Executor documentation

examples/
├── [service]-testplan.json              # Generated test plan
├── [service]-testplan-educational.json  # With education
└── test-results/
    └── YYYYMMDD_HHMMSS/
        ├── REPORT_EXECUTIVE.md
        ├── REPORT_DETAILED.md
        └── REPORT_FOLLOWUP.md
```

## How to Use a Skill

### Method 1: Direct Invocation

```
# Tell Claude which skill to use:
"Use the testplan-generator skill to create a test plan for my operator"
```

Claude will follow the skill's instructions and generate the output.

### Method 2: Reference Skill File

```
# Share the skill file:
cat skills/testplan-generator/skill.md

# Then provide context:
"Generate test plan for: /path/to/my-service"
```

### Method 3: Use Skill Documentation

```
# Share the README:
cat skills/testplan-generator/README.md

# Then request:
"Generate following the example in the README"
```

## Best Practices

### For testplan-generator

**Do:**
- Provide as many artifacts as possible
- Specify all environments to cover
- Define clear coverage goals
- Review gap analysis

**Don't:**
- Skip environment variants
- Ignore the gap analysis
- Expect 100% coverage on first pass

### For testplan-educator

**Do:**
- Specify target audience clearly
- Provide documentation sources
- Review with actual beginners
- Iterate based on feedback

**Don't:**
- Use jargon without explaining
- Assume prerequisite knowledge
- Forget common mistakes
- Only provide one learning modality

### For testplan-executor

**Do:**
- Verify prerequisites first
- Capture all outputs
- Follow remediation plans
- Re-execute after fixes

**Don't:**
- Skip prerequisite checks
- Continue blindly after failures
- Ignore gap analysis
- Forget to preserve artifacts

## Success Metrics

### For testplan-generator
- ✅ 70%+ feature coverage
- ✅ 3+ environment variants per test
- ✅ Negative tests included
- ✅ Gaps documented with recommendations

### For testplan-educator
- ✅ Every concept explained
- ✅ Every step has learning context
- ✅ Common mistakes documented
- ✅ Multiple learning resources
- ✅ Passes "beginner test"

### For testplan-executor
- ✅ 100% output capture
- ✅ Clear pass/fail determination
- ✅ Root cause analysis on failures
- ✅ Prioritized action plan
- ✅ Multiple report formats

## Troubleshooting

### "Which skill should I use?"

See [Skill Selection Guide](#skill-selection-guide) above.

### "Can I use multiple skills together?"

Yes! See [Usage Patterns](#usage-patterns) for common workflows.

### "How long does each skill take?"

- **testplan-generator**: 1-3 hours (depending on complexity)
- **testplan-educator**: 1-2 hours (per test plan)
- **testplan-executor**: Variable (depends on test plan size)

### "What if the output isn't what I need?"

1. Review the skill's README for examples
2. Provide more specific context
3. Iterate with Claude to refine
4. Use AskUserQuestion for clarification

## Next Steps

1. **Read skill documentation**
   - [testplan-generator/README.md](testplan-generator/README.md)
   - [testplan-educator/README.md](testplan-educator/README.md)
   - [testplan-executor/README.md](testplan-executor/README.md)

2. **Try a simple example**
   - Start with testplan-generator
   - Use a small repository
   - Review the output

3. **Iterate and improve**
   - Use findings from executor
   - Enhance with educator
   - Build your test library

4. **Share and collaborate**
   - Commit test plans to git
   - Share HTML with team
   - Track in version control

## Contributing

To add a new skill to this suite:

1. Create directory: `skills/new-skill/`
2. Write `skill.md` with comprehensive instructions
3. Write `README.md` with examples and documentation
4. Update this README with skill information
5. Test the skill with Claude
6. Commit to repository

## Related Documentation

- [Master Prompt](.claude/prompts/MASTER-PROMPT.md) - Complete prompt system
- [Prompts README](.claude/prompts/README.md) - Individual prompts
- [Examples README](examples/README.md) - Test plan examples
- [Project README](README.md) - Main project documentation
