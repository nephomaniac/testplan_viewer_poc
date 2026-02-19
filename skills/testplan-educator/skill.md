# Test Plan Educator Skill

Transform test plans into educational experiences. Write and review tests with a focus on user experience, teaching concepts while testing, and ensuring testers learn as they execute.

## ⚠️ Skill Self-Verification (REQUIRED - Run First)

**BEFORE executing this skill, ALWAYS verify you're using the latest skill definition.**

### Verification Steps

1. **Check if skill definition exists in repository:**
   ```bash
   ls -la skills/testplan-educator/skill.md
   ```

2. **Verify this is a git repository:**
   ```bash
   git rev-parse --is-inside-work-tree 2>/dev/null || echo "Not a git repo"
   ```

3. **Check if local skill has uncommitted changes:**
   ```bash
   git status skills/testplan-educator/skill.md
   ```

4. **Compare local vs committed version:**
   ```bash
   # Check if there are differences between working copy and HEAD
   git diff skills/testplan-educator/skill.md
   ```

5. **Check remote for updates (if applicable):**
   ```bash
   # Fetch latest from remote (don't merge)
   git fetch origin main 2>/dev/null

   # Compare local version with remote
   git diff HEAD origin/main -- skills/testplan-educator/skill.md
   ```

### Decision Tree

**If differences are found, PROMPT the user:**

```
⚠️ Skill Definition Verification

I've detected differences in the testplan-educator skill definition:

[Show summary of differences]

Options:
1. Continue with current version (may be outdated)
2. Read latest version from repository and use that
3. Show me the full diff to review
4. Cancel and let me update the skill first

What would you like to do?
```

**Response handling:**

- **Option 1 (Continue):** Proceed with current skill, but warn that results may differ from latest spec
- **Option 2 (Use latest):** Read `skills/testplan-educator/skill.md` from disk and use that definition
- **Option 3 (Show diff):** Display full diff, then ask again
- **Option 4 (Cancel):** Stop execution, advise user to:
  ```bash
  # Pull latest changes
  git pull origin main

  # Or if local changes exist
  git stash
  git pull origin main
  git stash pop
  ```

### Verification Output

After verification, display:

```
✅ Skill Verification Complete

Skill: testplan-educator
Version: [git commit hash of skills/testplan-educator/skill.md]
Last Modified: [file modification date]
Status: [Up to date | Using local changes | Using remote version]

Proceeding with skill execution...
```

### When to Skip Verification

Only skip this verification if:
- User explicitly says "skip verification"
- Already verified in the same conversation session
- Emergency/time-critical situation explicitly stated by user

**Note:** This self-verification ensures Claude always uses the most current skill definition, preventing drift between the skill prompt and the repository source of truth.

---

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

### Phase 7: Safety and Impact Education (CRITICAL - 20-30 minutes)

**REQUIRED:** Educate users about test safety, system impact, and state management to prevent accidents and build confidence.

#### Objective

Transform safety metadata into educational content that helps users:
- Understand what impact each test has on the system
- Know when backups are required
- Learn how to safely execute and recover from tests
- Build confidence by understanding risks before acting

#### For Every Test: Add Safety Education

**1. Explain System Impact Type**

Add learning content that explains the test's impact classification:

**For Read-Only Tests:**
```json
{
  "learning": {
    "safety_overview": {
      "impact_type": "read-only",
      "what_it_means": "This test only observes the system - it makes no changes whatsoever. Think of it like looking through a window: you can see everything, but you can't touch anything.",
      "why_its_safe": "Read-only tests can be run anywhere (including production) without risk of breaking things. They're perfect for learning because mistakes have no consequences.",
      "confidence_builder": "Feel free to run this test multiple times, experiment with variations, and explore the outputs. You literally cannot break anything with these commands.",
      "learning_value": "Read-only tests are great for understanding current state and building mental models of how systems work."
    }
  }
}
```

**For Modifies-State Tests:**
```json
{
  "learning": {
    "safety_overview": {
      "impact_type": "modifies-state",
      "what_it_means": "This test creates or changes resources in the cluster. Think of it like rearranging furniture: you can always move things back, but it takes effort.",
      "what_changes": "Specifically, this test creates a ServiceMonitor resource. This resource will persist after the test completes.",
      "why_cleanup_matters": "If you don't clean up test resources, they accumulate and can: 1) waste cluster resources, 2) interfere with future tests, 3) create confusion about what's 'real' vs 'test'.",
      "how_to_be_safe": {
        "before_running": "Make sure you understand what resources will be created (see 'Affected Resources' below)",
        "while_running": "Pay attention to the output to confirm resources are created in the right namespace",
        "after_completing": "Run the cleanup test (test_cleanup_servicemonitor) to restore the system to its original state"
      },
      "backup_restore_pattern": {
        "explanation": "Because this test modifies state, we follow a 3-step safety pattern:",
        "steps": [
          {
            "step": "1. Verify state",
            "description": "Check current system state before making changes (like taking a mental snapshot)"
          },
          {
            "step": "2. Make changes",
            "description": "Execute the test and create/modify resources"
          },
          {
            "step": "3. Cleanup",
            "description": "Remove test resources and verify system returned to original state (run test_cleanup_servicemonitor)"
          }
        ]
      },
      "common_safety_mistakes": [
        {
          "mistake": "Running modifying tests in production without approval",
          "consequence": "Accidentally creates test resources in production clusters",
          "how_to_avoid": "Always verify your kubeconfig context before running: oc config current-context. If it says 'production', stop and switch to a test cluster.",
          "check_command": "oc config current-context | grep -i prod && echo 'WARNING: This is a production cluster!' || echo 'OK: Non-production cluster'"
        },
        {
          "mistake": "Forgetting to run cleanup tests",
          "consequence": "Test resources accumulate, wasting resources and creating confusion",
          "how_to_avoid": "After completing this test, immediately run the cleanup test: test_cleanup_servicemonitor. Set a reminder or add to your checklist.",
          "pro_tip": "Some teams use GitOps to auto-delete resources in test namespaces older than 24 hours"
        }
      ]
    },
    "affected_resources_explained": {
      "what_we_create": [
        {
          "resource": "ServiceMonitor 'test-monitor'",
          "namespace": "test-namespace",
          "why_it_exists": "Tells Prometheus to scrape metrics from our test service",
          "how_to_verify": "oc get servicemonitor test-monitor -n test-namespace",
          "how_to_remove": "oc delete servicemonitor test-monitor -n test-namespace (or run test_cleanup_servicemonitor)"
        }
      ]
    }
  }
}
```

**For Destructive Tests:**
```json
{
  "learning": {
    "safety_overview": {
      "impact_type": "destructive",
      "what_it_means": "⚠️ CRITICAL: This test DELETES resources permanently. Think of it like tearing down a building: once gone, it's gone forever unless you have blueprints (backups).",
      "what_gets_deleted": "This test deletes a PersistentVolumeClaim, which means all data in the volume is permanently lost.",
      "why_extremely_risky": "Destructive tests can cause data loss, service outages, and irreversible changes. They should ONLY be run in isolated test environments with proper backups.",
      "required_safety_steps": {
        "mandatory_backup": {
          "requirement": "You MUST run the backup test first: test_backup_pvc_data",
          "why": "Without a backup, there's no way to recover if something goes wrong",
          "what_it_backs_up": "The backup test exports all data from the PVC to a tar archive in your local filesystem",
          "verify_backup": "Check that backup/pvc-data-TIMESTAMP.tar exists and is not zero bytes before proceeding"
        },
        "confirmation_required": {
          "why": "Destructive tests require explicit confirmation to prevent accidents",
          "what_to_confirm": "Verify: 1) You're in the RIGHT cluster, 2) Backup completed successfully, 3) You understand this is permanent, 4) You have approval if needed"
        }
      },
      "backup_restore_pattern": {
        "explanation": "Destructive tests require a 4-step safety pattern:",
        "steps": [
          {
            "step": "1. Backup (MANDATORY)",
            "test_id": "test_backup_pvc_data",
            "description": "Export all data that will be deleted",
            "verification": "Verify backup file exists and contains data"
          },
          {
            "step": "2. Delete (DESTRUCTIVE)",
            "test_id": "test_delete_pvc",
            "description": "Permanently delete the PVC and its data",
            "cannot_undo": true
          },
          {
            "step": "3. Restore",
            "test_id": "test_restore_pvc_data",
            "description": "Recreate PVC and restore data from backup"
          },
          {
            "step": "4. Verify",
            "test_id": "test_verify_pvc_restored",
            "description": "Confirm data was restored correctly"
          }
        ]
      },
      "when_to_run_destructive_tests": {
        "appropriate": [
          "In isolated test clusters that can be rebuilt",
          "During disaster recovery drills (with backups)",
          "Testing backup/restore procedures",
          "Validating delete operations work as expected"
        ],
        "never": [
          "In production clusters (even with backups)",
          "Without completing the backup test first",
          "If you're unsure what will be deleted",
          "If you don't have approval from cluster owner"
        ]
      },
      "common_safety_mistakes": [
        {
          "mistake": "Skipping the backup test to save time",
          "consequence": "Permanent data loss with no recovery option",
          "how_to_avoid": "ALWAYS run backup test first. It takes 2 minutes and could save hours of recovery work.",
          "real_story": "An engineer once skipped backup before testing PVC deletion on what they thought was a test cluster. It was production. Don't be that engineer."
        },
        {
          "mistake": "Not verifying the backup completed successfully",
          "consequence": "Backup exists but is empty or corrupted, no way to restore",
          "how_to_avoid": "After backup test, verify: 1) Backup file exists, 2) File size is reasonable (not 0 bytes), 3) Can list contents: tar -tzf backup.tar | head",
          "pro_tip": "Some teams require two different backup methods for destructive tests"
        }
      ],
      "emergency_contacts": {
        "if_something_goes_wrong": [
          "STOP immediately - don't make it worse",
          "Check if backup exists and is valid",
          "Run restore test if backup is good",
          "Contact cluster admin if restore fails",
          "Document what happened for incident review"
        ]
      }
    }
  }
}
```

**2. Add Safety Badges and Visual Warnings**

Enhance the test metadata with visual safety cues:

```json
{
  "metadata": {
    "safety_badges": {
      "read-only": {
        "icon": "👁️",
        "color": "green",
        "message": "Safe - Read-Only",
        "confidence_level": "Run freely, experiment safely"
      },
      "modifies-state": {
        "icon": "✏️",
        "color": "yellow",
        "message": "Caution - Modifies System",
        "confidence_level": "Run carefully, cleanup required"
      },
      "destructive": {
        "icon": "🚨",
        "color": "red",
        "message": "DANGER - Destructive",
        "confidence_level": "Backup required, expert supervision recommended"
      }
    }
  }
}
```

**3. Create Safety-Focused Learning Objectives**

Add specific learning objectives about safe test execution:

```json
{
  "learning": {
    "objectives": [
      "Understand the difference between read-only and state-modifying tests",
      "Learn the backup-modify-cleanup-verify pattern for safe testing",
      "Practice verifying cluster context before running modifying tests",
      "Build confidence identifying risky operations in test commands",
      "Master the cleanup process to restore system state"
    ],
    "safety_skills_gained": [
      "Identify whether a test command modifies state by analyzing the operations",
      "Verify you're in the correct cluster before running tests",
      "Execute backup/cleanup patterns to ensure safe test execution",
      "Recognize when a test requires expert review or approval"
    ]
  }
}
```

**4. Add Cleanup Test Education**

For every cleanup test, explain its purpose and importance:

```json
{
  "test_cleanup_servicemonitor": {
    "learning": {
      "cleanup_education": {
        "what_this_test_does": "Removes all resources created by test_create_servicemonitor, returning the system to its original state",
        "why_cleanup_matters": {
          "resource_waste": "Unused resources consume cluster capacity (CPU, memory, storage)",
          "test_interference": "Leftover test resources can cause future tests to fail or produce unexpected results",
          "confusion": "Makes it hard to distinguish between 'real' production resources and abandoned test artifacts",
          "cost": "In cloud environments, unused resources cost money"
        },
        "when_to_run_cleanup": [
          "Immediately after completing the main test (preferred)",
          "At the end of your testing session",
          "Before running the same test again (to ensure clean slate)",
          "If a test fails midway and leaves partial resources"
        ],
        "how_to_verify_cleanup_worked": [
          {
            "check": "Resource is gone",
            "command": "oc get servicemonitor test-monitor -n test-namespace",
            "expected": "Error from server (NotFound): servicemonitors.monitoring.coreos.com \"test-monitor\" not found",
            "meaning": "Success - resource was deleted"
          },
          {
            "check": "Namespace is clean",
            "command": "oc get all -n test-namespace",
            "expected": "No resources found in test-namespace namespace",
            "meaning": "All test resources removed"
          }
        ],
        "common_cleanup_mistakes": [
          {
            "mistake": "Assuming cleanup happened automatically",
            "reality": "Resources persist until explicitly deleted",
            "how_to_avoid": "Always verify cleanup worked by checking for the resource after deletion"
          }
        ]
      }
    }
  }
}
```

**5. Add Pre-Flight Safety Checklist**

For modifying and destructive tests, add a pre-execution safety checklist:

```json
{
  "test_execution": {
    "safety_checklist": {
      "title": "⚠️ Pre-Flight Safety Check",
      "description": "Complete this checklist BEFORE running this test",
      "required_checks": [
        {
          "check": "Verify cluster context",
          "command": "oc config current-context",
          "requirement": "Must NOT be a production cluster",
          "how_to_verify": "Context name should contain 'test', 'dev', or 'sandbox'",
          "if_wrong": "Switch context with: oc login <test-cluster-url>"
        },
        {
          "check": "Confirm backup completed (destructive tests only)",
          "command": "ls -lh backup/",
          "requirement": "Backup file must exist and have non-zero size",
          "how_to_verify": "See a recent .tar or .yaml file with size > 0",
          "if_missing": "Run backup test first: test_backup_pvc_data"
        },
        {
          "check": "Verify test namespace exists",
          "command": "oc get namespace test-namespace",
          "requirement": "Namespace must exist",
          "how_to_verify": "Command succeeds without error",
          "if_missing": "Create namespace: oc create namespace test-namespace"
        },
        {
          "check": "Check RBAC permissions",
          "command": "oc auth can-i create servicemonitor -n test-namespace",
          "requirement": "Must return 'yes'",
          "how_to_verify": "Output is exactly: yes",
          "if_no": "Request permissions from cluster admin"
        }
      ],
      "all_checks_passed": "✅ All safety checks passed - you may proceed with the test",
      "if_any_failed": "❌ Do NOT proceed until all checks pass. Fix the failures first."
    }
  }
}
```

#### Safety Education Quality Checklist

Before finalizing enhanced test plan:

- [ ] Every test has safety_overview explaining its impact type
- [ ] Modifying tests explain what resources they create/change
- [ ] Destructive tests have clear warnings and backup requirements
- [ ] Cleanup tests explain why they matter and how to verify success
- [ ] Pre-flight safety checklists added for risky tests
- [ ] Common safety mistakes documented with real consequences
- [ ] Backup/restore pattern explained step-by-step
- [ ] Visual safety badges/icons used for quick recognition
- [ ] "When to run" and "when NOT to run" guidance provided
- [ ] Emergency recovery steps included for destructive tests

#### Example: Complete Safety Education for a Modifying Test

```json
{
  "test_create_servicemonitor": {
    "metadata": {
      "title": "Create ServiceMonitor for metrics collection",
      "system_impact": {
        "type": "modifies-state",
        "description": "Creates a ServiceMonitor custom resource in test-namespace",
        "affected_resources": ["ServiceMonitor: test-monitor in test-namespace"],
        "reversible": true,
        "persistence": "permanent",
        "risk_level": "medium"
      },
      "state_management": {
        "requires_backup": false,
        "requires_cleanup": true,
        "cleanup_test": "test_cleanup_servicemonitor"
      },
      "safety": {
        "can_run_in_production": false,
        "requires_confirmation": true,
        "warning_message": "⚠️ This test creates a ServiceMonitor resource that will persist until cleanup. Run test_cleanup_servicemonitor when done.",
        "safe_to_retry": true,
        "idempotent": true
      }
    },
    "learning": {
      "safety_overview": {
        "impact_type": "modifies-state",
        "what_it_means": "This test creates a resource that persists after completion",
        "affected_resources_explained": "We create ServiceMonitor 'test-monitor' in test-namespace to demonstrate Prometheus configuration",
        "cleanup_required": "Yes - run test_cleanup_servicemonitor after completing this test",
        "why_cleanup_matters": "Prevents resource waste and test interference"
      },
      "safety_skills_gained": [
        "Identify state-modifying commands (oc create, apply, patch)",
        "Verify cluster context before modifying resources",
        "Execute cleanup tests to restore original state"
      ]
    },
    "test_execution": {
      "safety_checklist": {
        "required_checks": [
          {
            "check": "Verify not in production",
            "command": "oc config current-context | grep -v prod",
            "requirement": "Must NOT contain 'prod'"
          }
        ]
      },
      "steps": [
        {
          "step_number": 1,
          "title": "Verify current cluster context (safety check)",
          "learning_note": "ALWAYS verify you're in the right cluster before modifying resources. This prevents accidental changes to production.",
          "command": "oc config current-context",
          "expected_output": "Contains 'test', 'dev', or 'sandbox' - NOT 'prod'",
          "why_this_step": "Safety first! Verifying context prevents accidentally creating test resources in production clusters.",
          "common_errors": [
            {
              "error": "Context shows 'production-cluster'",
              "solution": "STOP. Switch to test cluster: oc login <test-cluster-url>",
              "learning_point": "Production clusters should never have test resources. Always double-check context."
            }
          ]
        },
        {
          "step_number": 2,
          "title": "Create ServiceMonitor resource",
          "learning_note": "This command creates a new Kubernetes resource that will persist. Remember to run cleanup later!",
          "command": "oc apply -f servicemonitor.yaml",
          "expected_output": "servicemonitor.monitoring.coreos.com/test-monitor created"
        }
      ]
    },
    "next_steps": {
      "on_success": "✅ Test complete! Next: Run test_cleanup_servicemonitor to remove test resources and restore the cluster.",
      "cleanup_reminder": {
        "message": "Don't forget cleanup!",
        "test_to_run": "test_cleanup_servicemonitor",
        "why": "Removes the ServiceMonitor we created, preventing resource waste"
      }
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
