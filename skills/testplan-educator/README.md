# Test Plan Educator Skill

**Focus:** Educational enhancement and user experience

**Purpose:** Transform tests into learning experiences that educate while validating

## What This Skill Does

The Test Plan Educator skill enhances test plans with educational content:

- 📚 **Concept explanations** - Plain-language definitions with analogies
- 🎓 **Learning objectives** - What testers will learn from each test
- ⚠️ **Common mistakes** - Pitfalls to avoid with prevention strategies
- 🔍 **Step context** - Why each step matters, not just what to do
- 📖 **Multi-modal resources** - Docs, videos, tutorials for different learning styles
- 🎯 **Progressive complexity** - Structured paths from beginner to advanced

## When to Use

Use this skill when you want to:

- Make tests beginner-friendly
- Add educational value to existing tests
- Explain complex concepts through testing
- Create onboarding materials
- Improve test documentation
- Help testers learn while validating

## Inputs Required

- **Test plan JSON** (required) - Test plan to enhance
- **Target audience** (required) - Who will use this
- **Learning objectives** (optional) - Specific learning goals
- **Documentation sources** (optional) - Where to find explanations

## Example Usage

### Example 1: Enhance for Beginners

```
Enhance this test plan for new engineers with no Kubernetes experience.

Test plan: examples/observability-operator-testplan.json
Target audience: New SRE engineers (0-6 months experience)
Learning objectives:
  - Understand Prometheus architecture
  - Learn ServiceMonitor pattern
  - Practice kubectl/oc commands
```

**Claude will:**
1. Add concept explanations (Prometheus, ServiceMonitor, etc.)
2. Include "learning notes" for every step
3. Explain kubectl/oc commands in detail
4. Add common beginner mistakes
5. Provide multiple learning resources
6. Output: `examples/observability-operator-educational.json`

### Example 2: Add Documentation Links

```
Add official documentation and learning resources to this test plan.

Test plan: examples/basic-testplan.json
Documentation sources:
  - OpenShift docs: https://docs.openshift.com/
  - Prometheus Operator: https://prometheus-operator.dev/
  - Internal wiki: https://wiki.company.com/monitoring
```

**Claude will:**
1. Research each documented concept
2. Find relevant documentation sections
3. Add links to every test step
4. Include video tutorials and guides
5. Reference internal wiki pages
6. Output: `examples/basic-testplan-documented.json`

### Example 3: Create Alternative Learning Paths

```
Add alternative learning paths for different skill levels.

Test plan: examples/comprehensive-testplan.json
Audiences:
  - Kubernetes beginners
  - Experienced K8s users (new to Prometheus)
  - Debugging-focused (troubleshooting emphasis)
```

**Claude will:**
1. Analyze test complexity
2. Create beginner path (fundamentals first)
3. Create experienced path (skip basics)
4. Create debugging path (troubleshooting focus)
5. Mark which concepts to skip/emphasize
6. Output: `examples/comprehensive-testplan-multipath.json`

## Output Structure

Enhanced test plan includes:

### Concept Explanations
```json
{
  "concepts": {
    "service_monitor": {
      "title": "ServiceMonitor Custom Resource",
      "description": "Plain-language explanation",
      "why_it_matters": "Practical relevance",
      "real_world_analogy": "Relatable comparison",
      "common_misconceptions": [...],
      "learn_more": [
        {"type": "video", "title": "...", "url": "..."},
        {"type": "docs", "title": "...", "url": "..."}
      ]
    }
  }
}
```

### Enhanced Steps
```json
{
  "steps": [
    {
      "learning_note": "What this teaches",
      "command_explanation": {
        "what_it_does": "...",
        "why_this_command": "...",
        "alternatives": [...]
      },
      "output_explanation": {
        "what_it_means": "...",
        "success_indicators": [...]
      },
      "why_this_step": "Purpose and context",
      "hands_on_exploration": {
        "suggested_experiments": [...],
        "questions_to_ponder": [...]
      }
    }
  ]
}
```

### Common Mistakes
```json
{
  "learning": {
    "common_beginner_mistakes": [
      {
        "mistake": "What people do wrong",
        "consequence": "What happens",
        "how_to_avoid": "Prevention strategy",
        "recovery": "How to fix it"
      }
    ]
  }
}
```

### Learning Paths
```json
{
  "learning_path": {
    "alternative_paths": {
      "experienced_with_kubernetes": {
        "sequence": ["test_2", "test_3", "test_1"],
        "skip_concepts": ["pods", "services"],
        "focus_on": ["prometheus", "monitoring"]
      }
    }
  }
}
```

## Integration with Other Skills

**Works well with:**
- **testplan-generator** - Enhance generated comprehensive test plans
- **testplan-executor** - Execute educational tests with learning context

**Typical workflow:**
1. **testplan-generator** → Generate comprehensive coverage
2. **testplan-educator** → Add educational content (this skill)
3. **testplan-executor** → Execute with learning feedback

## Best Practices

### Do:
- Write for your target audience (match their level)
- Use real-world analogies
- Explain "why" not just "what"
- Include visual aids where possible
- Provide multiple learning resources
- Celebrate success, normalize mistakes

### Don't:
- Use jargon without explaining it
- Assume prerequisite knowledge
- Skip the "why this matters"
- Forget common mistakes section
- Only provide one type of resource (text-only)
- Write in condescending or overly technical tone

## Success Criteria

Enhanced test plan achieves:
- ✅ Every concept explained in plain language
- ✅ Every step includes learning context
- ✅ Common mistakes identified
- ✅ Multiple learning modalities
- ✅ Appropriate for target audience
- ✅ Positive, encouraging tone

## Files Generated

After using this skill:
```
examples/
  [service]-educational.json       # Enhanced test plan
  [service]-concepts-guide.md      # Concept reference
  [service]-learning-resources.md  # Curated resource list
```

## Audience Profiles

### Beginner (0-6 months)
- Explain everything
- Provide extensive context
- Include lots of resources
- Anticipate many mistakes

### Intermediate (6-18 months)
- Focus on "why" over "what"
- Less hand-holding
- More complex scenarios
- Deeper explanations

### Advanced (18+ months)
- Architecture and design
- Optimization techniques
- Edge cases and gotchas
- Production considerations

## Next Steps

After educational enhancement:
1. Review with target audience member
2. Test with actual beginners
3. Iterate based on feedback
4. Generate HTML: `make run`
5. Use **testplan-executor** to validate flow
