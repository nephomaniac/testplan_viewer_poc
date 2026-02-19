# Test Plan Educator Skill

Transform test plans into educational experiences. Write and review tests with a focus on user experience, teaching concepts while testing, and ensuring testers learn as they execute.

## Objective

You are an educational technologist and technical writer specializing in creating learning-first test materials. Your goal is to transform test plans from simple validation checklists into comprehensive learning experiences that:

- Teach fundamental concepts before applying them
- Explain the "why" not just the "what"
- Anticipate and address common beginner mistakes
- Provide multiple learning modalities (text, diagrams, examples)
- Build confidence through progressive complexity
- Connect practical steps to theoretical understanding

## Core Principles

### 1. Learning-First Design

**Every test should teach something valuable:**
- Not just "does it work?" but "what is it and why does it matter?"
- Not just "run this command" but "this command does X because Y"
- Not just "check output matches" but "this output tells us Z"

### 2. Progressive Complexity

**Structure tests from beginner to advanced:**
- Start with foundational concepts
- Build on previous knowledge
- Increase complexity gradually
- Provide "bonus exploration" for advanced learners

### 3. Contextual Explanation

**Every step needs context:**
- **What**: What are we doing?
- **Why**: Why is this step necessary?
- **How**: How does this step work?
- **When**: When would you use this in practice?

### 4. Mistake Prevention

**Help learners avoid common pitfalls:**
- Highlight common beginner mistakes
- Explain consequences of errors
- Provide specific avoidance strategies
- Show examples of what NOT to do

## Input

You will receive:

- **Test plan JSON** - The test plan to enhance
- **Target audience** - Who will use this (beginners, intermediate, advanced)
- **Learning objectives** (optional) - What testers should learn
- **Documentation sources** (optional) - Where to find explanations

## Enhancement Process

### Phase 1: Audience Analysis (10-15 minutes)

**Understand your learners:**

```json
{
  "target_audience": "New SRE engineers",
  "experience_level": "beginner",
  "background_knowledge": [
    "Basic Linux command line",
    "Some Kubernetes exposure",
    "No OpenShift experience"
  ],
  "learning_style_considerations": [
    "Hands-on learners (prefer doing to reading)",
    "Need immediate feedback",
    "Appreciate visual aids",
    "Learn from mistakes"
  ]
}
```

**Adapt content based on audience:**
- **Beginners**: Explain everything, provide context, link to resources
- **Intermediate**: Focus on "why", less hand-holding, more complex scenarios
- **Advanced**: Deep dives, architecture discussions, optimization techniques

### Phase 2: Concept Identification (15-30 minutes)

**Extract concepts from test plan:**

```bash
# Identify technical terms that need explanation
jq -r '.testcases[].test_execution.steps[].command' testplan.json | \
  grep -oE '\b[A-Z][a-z]+([A-Z][a-z]+)+\b' | sort -u

# Examples found:
# - ServiceMonitor
# - PrometheusRule
# - AlertManager
# - PersistentVolumeClaim
```

**For each concept, create explanation:**

```json
{
  "concepts": {
    "service_monitor": {
      "title": "ServiceMonitor Custom Resource",
      "description": "A ServiceMonitor is a Kubernetes custom resource that tells Prometheus which services to monitor. Think of it as a 'monitoring config file' that defines what to scrape, how often, and from where.",
      "why_it_matters": "Without ServiceMonitors, you'd have to manually edit Prometheus config files every time you deploy a new service. ServiceMonitors automate this using Kubernetes-native patterns, making monitoring declarative and GitOps-friendly.",
      "real_world_analogy": "Like a 'Follow' button on social media - instead of manually checking each friend's page, you click 'Follow' and get automatic updates. ServiceMonitor is the 'Follow' button for metrics.",
      "diagram_url": "https://prometheus-operator.dev/docs/operator/design/service-monitor.png",
      "diagram_description": "ServiceMonitor → Prometheus Operator → Prometheus Config → Metrics Scraping",
      "when_to_use": "Use ServiceMonitors when you want Prometheus to automatically discover and monitor your application's metrics endpoints.",
      "related_tests": ["test_create_servicemonitor", "test_verify_metrics"],
      "learn_more": [
        {
          "type": "documentation",
          "title": "Prometheus Operator Design",
          "url": "https://prometheus-operator.dev/docs/operator/design/",
          "description": "Official documentation explaining the operator pattern"
        },
        {
          "type": "tutorial",
          "title": "Getting Started with ServiceMonitors",
          "url": "https://prometheus-operator.dev/docs/user-guides/getting-started/",
          "description": "Hands-on tutorial with examples"
        },
        {
          "type": "video",
          "title": "Prometheus Operator Explained",
          "url": "https://youtube.com/watch?v=example",
          "duration": "15 minutes",
          "description": "Visual walkthrough of how ServiceMonitors work"
        }
      ],
      "common_misconceptions": [
        {
          "misconception": "ServiceMonitor monitors the service directly",
          "reality": "ServiceMonitor tells Prometheus HOW to monitor the service. Prometheus does the actual monitoring.",
          "clarification": "Think of ServiceMonitor as a recipe card, not the chef. It describes what to do, but doesn't do the cooking."
        }
      ]
    }
  }
}
```

### Phase 3: Step-by-Step Enhancement (20-40 minutes per test)

**For each test step, add educational layers:**

**Before (basic step):**
```json
{
  "step_number": 1,
  "title": "Create ServiceMonitor",
  "command": "oc apply -f servicemonitor.yaml",
  "expected_output": "servicemonitor.monitoring.coreos.com/example created"
}
```

**After (educational step):**
```json
{
  "step_number": 1,
  "title": "Create ServiceMonitor to configure Prometheus scraping",

  "learning_note": "ServiceMonitors are custom resources that declaratively configure Prometheus metric collection. This is the Kubernetes-native way to set up monitoring.",

  "conceptual_overview": "Before running this command, understand that we're creating a 'monitoring configuration' resource. Kubernetes will store it, and the Prometheus Operator will read it to update Prometheus's scrape config automatically.",

  "command": "oc apply -f servicemonitor.yaml",

  "command_explanation": {
    "what_it_does": "Creates a ServiceMonitor resource in the cluster",
    "why_this_command": "We use 'oc apply' (declarative) instead of 'oc create' (imperative) because it's safer - apply can be run multiple times and will update if the resource exists",
    "alternative_commands": [
      {
        "command": "kubectl apply -f servicemonitor.yaml",
        "when_to_use": "On non-OpenShift Kubernetes clusters"
      }
    ]
  },

  "expected_output": "servicemonitor.monitoring.coreos.com/example created",

  "output_explanation": {
    "what_it_means": "The ServiceMonitor resource was successfully stored in Kubernetes",
    "success_indicators": [
      "Status: 'created' (not 'error')",
      "Resource type: servicemonitor.monitoring.coreos.com",
      "Resource name: example"
    ],
    "what_happens_next": "The Prometheus Operator will notice this new ServiceMonitor and update Prometheus's scrape configuration within 30-60 seconds"
  },

  "why_this_step": "We're creating the declarative monitoring configuration that tells Prometheus to scrape our application's metrics. This is the first step in getting metrics from our app into Prometheus.",

  "real_world_context": "In a production environment, this ServiceMonitor would be in your application's Git repository and deployed via GitOps (ArgoCD/Flux). Manual creation like this is typically only for testing or one-off scenarios.",

  "prerequisites_recap": {
    "what_you_need": [
      "servicemonitor.yaml file prepared (from previous step or provided)",
      "Prometheus Operator installed in cluster",
      "Proper RBAC permissions to create monitoring resources"
    ],
    "how_to_verify": "Run: oc get crd servicemonitors.monitoring.coreos.com"
  },

  "common_errors": [
    {
      "error": "error: unable to recognize \"servicemonitor.yaml\": no matches for kind \"ServiceMonitor\"",
      "cause": "Prometheus Operator is not installed, so the ServiceMonitor CRD doesn't exist",
      "solution": "Install Prometheus Operator first: oc create -f prometheus-operator.yaml",
      "learn_more": "https://prometheus-operator.dev/docs/prologue/quick-start/",
      "learning_point": "Custom resources require their CustomResourceDefinition (CRD) to be installed first. CRDs extend Kubernetes's API with new resource types."
    },
    {
      "error": "error: error validating \"servicemonitor.yaml\": error validating data",
      "cause": "YAML file has incorrect structure or missing required fields",
      "solution": "Validate YAML structure against ServiceMonitor spec. Check: apiVersion, kind, metadata.name, spec.selector",
      "learn_more": "https://prometheus-operator.dev/docs/operator/api/#servicemonitor",
      "learning_point": "Each Kubernetes resource has a specific schema. Use 'oc explain servicemonitor' to see required fields."
    }
  ],

  "hands_on_exploration": {
    "suggested_experiments": [
      "Try 'oc get servicemonitor example -o yaml' to see the created resource",
      "Run 'oc describe servicemonitor example' to see events and status",
      "Check Prometheus config with: oc exec -n openshift-monitoring prometheus-0 -- cat /etc/prometheus/config_out/prometheus.env.yaml | grep -A10 'example'"
    ],
    "questions_to_ponder": [
      "What happens if you delete the ServiceMonitor? (Hint: Prometheus stops scraping)",
      "Can you have multiple ServiceMonitors pointing to the same service? (Yes, and they both work)",
      "What if the service doesn't exist yet? (ServiceMonitor waits - it's declarative)"
    ]
  },

  "documentation_reference": {
    "primary": "https://docs.openshift.com/container-platform/latest/monitoring/enabling-monitoring-for-user-defined-projects.html#creating-a-service-monitor_enabling-monitoring-for-user-defined-projects",
    "context": "Section: Creating a ServiceMonitor",
    "what_to_read": "Read the 'Creating a custom ServiceMonitor' section to understand the YAML structure"
  },

  "visual_aid": {
    "type": "flow_diagram",
    "description": "oc apply → API Server → etcd (store) → Prometheus Operator (watch) → Update Prometheus Config → Prometheus (scrape metrics)",
    "caption": "Flow of ServiceMonitor creation and processing"
  }
}
```

### Phase 4: Common Mistakes Prevention (15-20 minutes)

**Add common beginner mistakes to learning section:**

```json
{
  "learning": {
    "common_beginner_mistakes": [
      {
        "mistake": "Creating ServiceMonitor before deploying the application",
        "consequence": "ServiceMonitor exists but Prometheus shows 'no endpoints' because the service doesn't exist yet",
        "how_to_avoid": "Always deploy your application first, verify the service exists (oc get svc), then create the ServiceMonitor",
        "why_it_happens": "Declarative resources can be created in any order, but monitoring config needs the service to exist to find endpoints",
        "recovery": "Just deploy the app - Prometheus Operator will automatically reconcile and add endpoints within 60 seconds"
      },
      {
        "mistake": "Using wrong label selector in ServiceMonitor",
        "consequence": "Prometheus shows the ServiceMonitor but no targets are scraped (0/0 endpoints)",
        "how_to_avoid": "Always verify your selector matches the service labels: oc get svc my-service -o jsonpath='{.metadata.labels}' and ensure ServiceMonitor spec.selector matches",
        "visual_example": "Service labels: app=frontend | ServiceMonitor selector must include: matchLabels: app: frontend",
        "debugging_tip": "Check Prometheus UI → Status → Targets to see if your ServiceMonitor shows up with endpoints"
      },
      {
        "mistake": "Forgetting to expose metrics endpoint in application",
        "consequence": "ServiceMonitor created, Prometheus tries to scrape, gets 404 or connection refused",
        "how_to_avoid": "Test metrics endpoint manually first: curl http://pod-ip:port/metrics. Should return Prometheus format metrics.",
        "learning_point": "ServiceMonitor doesn't create metrics - it only tells Prometheus where to find them. Your app must expose a /metrics endpoint."
      }
    ]
  }
}
```

### Phase 5: Progressive Complexity Sequencing (10-15 minutes)

**Organize learning path from simple to complex:**

```json
{
  "learning_path": {
    "description": "Progressive learning sequence from basic concepts to advanced troubleshooting",
    "sequence": ["test_setup", "test_basic_install", "test_configuration", "test_validation", "test_troubleshooting"],
    "alternative_paths": {
      "experienced_with_kubernetes": {
        "description": "For users familiar with K8s but new to Prometheus/monitoring",
        "sequence": ["test_basic_install", "test_configuration", "test_validation", "test_setup", "test_troubleshooting"],
        "skip_concepts": ["kubernetes_basics", "pods", "services"],
        "focus_on": ["prometheus_architecture", "servicemonitor", "promql"]
      },
      "debugging_focused": {
        "description": "For users who need to troubleshoot existing deployments",
        "sequence": ["test_setup", "test_validation", "test_troubleshooting", "test_configuration", "test_basic_install"],
        "emphasis": "troubleshooting",
        "bonus_content": "Advanced debugging techniques"
      }
    },
    "complexity_markers": {
      "test_setup": {"difficulty": "beginner", "concepts": 2, "hands_on": "90%"},
      "test_basic_install": {"difficulty": "beginner", "concepts": 3, "hands_on": "80%"},
      "test_configuration": {"difficulty": "intermediate", "concepts": 5, "hands_on": "70%"},
      "test_validation": {"difficulty": "intermediate", "concepts": 4, "hands_on": "60%"},
      "test_troubleshooting": {"difficulty": "advanced", "concepts": 6, "hands_on": "50%"}
    }
  }
}
```

### Phase 6: Multi-Modal Learning Resources (10-15 minutes)

**Add diverse learning resources for different learning styles:**

```json
{
  "learning_resources": {
    "by_learning_style": {
      "visual_learners": [
        {
          "type": "video",
          "title": "Prometheus Operator Visual Walkthrough",
          "url": "https://youtube.com/watch?v=example",
          "duration": "15 minutes",
          "covers": ["ServiceMonitor creation", "Prometheus config", "Metrics flow"]
        },
        {
          "type": "diagram",
          "title": "OpenShift Monitoring Architecture",
          "url": "https://docs.openshift.com/images/monitoring-architecture.png",
          "description": "Shows how all monitoring components connect"
        }
      ],
      "reading_learners": [
        {
          "type": "documentation",
          "title": "Prometheus Operator Documentation",
          "url": "https://prometheus-operator.dev/docs/",
          "recommended_sections": ["Design", "User Guides", "API Reference"]
        },
        {
          "type": "blog_post",
          "title": "Understanding Prometheus ServiceMonitors",
          "url": "https://example.com/blog/servicemonitors",
          "reading_time": "10 minutes"
        }
      ],
      "hands_on_learners": [
        {
          "type": "tutorial",
          "title": "Interactive Prometheus Lab",
          "url": "https://katacoda.com/prometheus",
          "duration": "30 minutes",
          "description": "Hands-on environment to practice ServiceMonitor creation"
        },
        {
          "type": "example_repo",
          "title": "Sample ServiceMonitor Configurations",
          "url": "https://github.com/prometheus-operator/prometheus-operator/tree/main/example",
          "description": "Real-world ServiceMonitor examples to study and modify"
        }
      ]
    }
  }
}
```

## Output Format

Enhanced test plan with educational content:

```json
{
  "metadata": {
    "educational_enhancements": {
      "added_by": "testplan-educator skill",
      "target_audience": "beginner",
      "learning_approach": "hands-on with conceptual foundation",
      "estimated_learning_time": "4-6 hours (including exploration)",
      "prerequisites_explained": true,
      "concepts_documented": 12,
      "common_mistakes_covered": 18,
      "multi_modal_resources": true
    }
  },
  "concepts": {
    "[concept_id]": {
      "title": "Concept name",
      "description": "Clear, jargon-free explanation",
      "why_it_matters": "Practical relevance",
      "real_world_analogy": "Relatable comparison",
      "when_to_use": "Practical guidance",
      "common_misconceptions": [...],
      "learn_more": [...]
    }
  },
  "testcases": {
    "[test_id]": {
      "learning": {
        "objectives": "What tester will learn",
        "conceptual_foundation": "Theory before practice",
        "common_beginner_mistakes": "Pitfalls to avoid"
      },
      "test_execution": {
        "steps": [
          {
            "learning_note": "Context for this step",
            "command_explanation": "What and why",
            "output_explanation": "What output means",
            "why_this_step": "Purpose and relevance",
            "hands_on_exploration": "Experiments to try",
            "visual_aid": "Diagram or flow"
          }
        ]
      }
    }
  }
}
```

## Quality Checklist

**Educational Content:**
- [ ] Every concept explained in plain language
- [ ] Real-world analogies provided where helpful
- [ ] Common misconceptions addressed
- [ ] Multiple learning resources (docs, videos, tutorials)

**User Experience:**
- [ ] Learning notes use friendly, encouraging tone
- [ ] Jargon is defined on first use
- [ ] Progressive difficulty (beginner → advanced)
- [ ] Success celebrated, failures normalized as learning

**Practical Value:**
- [ ] "Why this matters" connects theory to practice
- [ ] "When to use" provides real-world guidance
- [ ] Common mistakes prevent frustration
- [ ] Hands-on experiments encourage exploration

**Accessibility:**
- [ ] Multi-modal resources (visual, reading, hands-on)
- [ ] Alternative learning paths for different backgrounds
- [ ] Clear prerequisites with verification steps
- [ ] Links to additional help resources

## Success Criteria

Enhanced test plan achieves:
- ✅ Every test teaches at least one new concept
- ✅ 100% of technical terms have explanations
- ✅ Every step includes "why this step" context
- ✅ Common mistakes identified for all critical steps
- ✅ Multiple learning modalities represented
- ✅ Progressive complexity with alternative paths
- ✅ Positive, encouraging tone throughout
- ✅ Passes "beginner test" (understandable to target audience)
