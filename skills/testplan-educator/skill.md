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

### Phase 8: Negative Testing Education (20-30 minutes)

**REQUIRED:** Add educational content to help testers understand and learn from negative tests.

#### Objective

Transform negative tests from simple error checks into learning experiences that teach:
- Why systems fail and how to handle failures gracefully
- How to identify invalid inputs and edge cases
- Error message interpretation and troubleshooting
- Defensive programming and input validation practices

#### For Every Negative Test: Add Educational Layers

**1. Explain the Value of Negative Testing**

Add conceptual overview explaining why we test failure scenarios:

```json
{
  "learning": {
    "negative_testing_overview": {
      "what_is_negative_testing": "Testing how a system handles invalid inputs, error conditions, and unexpected scenarios. Think of it like a crash test for cars - we intentionally try to break things in controlled ways to verify they fail safely.",

      "why_negative_tests_matter": {
        "error_handling": "Real users WILL provide invalid inputs - accidentally or intentionally. Negative tests verify your system handles these gracefully instead of crashing.",
        "security": "Many security vulnerabilities come from improper error handling. Negative tests help identify weaknesses before attackers do.",
        "user_experience": "Good error messages guide users to fix problems. Negative tests verify error messages are helpful, not cryptic.",
        "robustness": "Systems that handle errors well are more reliable and easier to troubleshoot in production."
      },

      "learning_value": {
        "for_testers": "Learn to think adversarially - what could go wrong? This mindset improves test coverage.",
        "for_developers": "Understanding common failure modes improves error handling code.",
        "for_operators": "Knowing how systems fail helps troubleshoot production issues faster."
      },

      "real_world_examples": [
        {
          "scenario": "User enters 'abc' instead of a number in replica count",
          "bad_system": "Silently fails or crashes with stack trace",
          "good_system": "Returns: 'Error: replicas must be a positive integer'",
          "what_we_test": "That the system responds like the 'good' example"
        },
        {
          "scenario": "Attempting to create resource in non-existent namespace",
          "bad_system": "Creates namespace automatically (unintended side effect)",
          "good_system": "Returns: 'Error: namespace \"xyz\" not found'",
          "what_we_test": "System doesn't auto-create namespaces, gives clear error"
        }
      ]
    }
  }
}
```

**2. Frame Failures as Learning Opportunities**

For each negative test, explain what we learn from the expected failure:

```json
{
  "child_tests": [{
    "title": "Attempt deployment with invalid image reference",
    "test_type": "negative",

    "learning_note": "This test intentionally uses an invalid image to verify Kubernetes catches the error BEFORE trying to run the pod. This protects against typos and configuration mistakes.",

    "what_we_expect_to_happen": {
      "failure_is_success": "⚠️ IMPORTANT: This test SHOULD fail! We're testing that the system correctly REJECTS invalid input.",
      "what_failure_means": "If we see an error message about image not found, that's actually SUCCESS - the system is working correctly.",
      "what_success_would_mean": "If the deployment succeeded with an invalid image, that would be a BUG in Kubernetes."
    },

    "error_message_education": {
      "expected_error": "Error: ErrImagePull or ImagePullBackOff",
      "what_each_part_means": {
        "ErrImagePull": "Kubernetes tried to download the image but couldn't find it",
        "ImagePullBackOff": "Kubernetes is waiting before retrying the pull (exponential backoff)"
      },
      "why_this_error_is_good": "It tells us exactly what went wrong (image not found) instead of failing silently or with a cryptic code.",
      "what_to_look_for": [
        "Error message mentions the specific image that failed",
        "Error clearly states the problem (not found vs. permission denied vs. network error)",
        "Pod status shows the error (oc describe pod shows events)"
      ]
    },

    "common_real_world_causes": [
      {
        "cause": "Typo in image name",
        "example": "nginz instead of nginx",
        "how_to_fix": "Double-check image name against registry",
        "prevention": "Use image tags from verified sources, not manual typing"
      },
      {
        "cause": "Image doesn't exist in registry",
        "example": "custom-app:v2.0 when only v1.0 exists",
        "how_to_fix": "Check available tags: docker search or registry UI",
        "prevention": "CI/CD should verify image exists before deploying"
      },
      {
        "cause": "Private image without pull secret",
        "example": "my-company/private-app:latest",
        "how_to_fix": "Create ImagePullSecret and reference in pod spec",
        "prevention": "Document which images require pull secrets"
      }
    ],

    "learning_objectives": [
      "Understand that 'expected failures' validate error handling",
      "Learn to interpret Kubernetes error messages",
      "Identify root causes from error symptoms",
      "Practice troubleshooting common deployment failures"
    ],

    "hands_on_exploration": {
      "try_variations": [
        "What happens with a valid image name but wrong tag? (ErrImagePull: tag not found)",
        "What if you spell nginx as 'nginz'? (Same error - image doesn't exist)",
        "Try an image from a private registry without credentials? (Error: pull access denied)"
      ],
      "questions_to_ponder": [
        "How quickly does Kubernetes detect the error? (Within seconds)",
        "Does the error message help you fix the problem? (Yes - tells you which image failed)",
        "What would happen in production if this got deployed? (Pods fail to start, service degraded)"
      ]
    }
  }]
}
```

**3. Teach Error Message Interpretation**

Add educational content about understanding and using error messages:

```json
{
  "learning": {
    "error_message_literacy": {
      "title": "Reading and Understanding Error Messages",

      "anatomy_of_error_message": {
        "typical_format": "Error from server (ErrorType): resource \"name\" reason",
        "example": "Error from server (Forbidden): secrets \"pd-secret\" is forbidden: User cannot create resource",
        "parts_explained": {
          "source": "'Error from server' = came from Kubernetes API server, not client-side",
          "error_type": "'(Forbidden)' = HTTP status category (401/403 auth issue)",
          "resource": "'secrets \"pd-secret\"' = what you tried to access",
          "reason": "'is forbidden: User cannot create' = why it failed"
        }
      },

      "common_error_types": [
        {
          "type": "NotFound",
          "what_it_means": "The resource doesn't exist",
          "common_causes": ["Typo in name", "Wrong namespace", "Resource not created yet"],
          "how_to_fix": "Verify resource exists: oc get <resource-type> <name> -n <namespace>"
        },
        {
          "type": "AlreadyExists",
          "what_it_means": "Resource with that name already exists",
          "common_causes": ["Trying to create duplicate", "Previous test didn't cleanup"],
          "how_to_fix": "Delete existing resource or use 'oc apply' instead of 'create'"
        },
        {
          "type": "Forbidden",
          "what_it_means": "You don't have permission to perform this operation",
          "common_causes": ["Insufficient RBAC permissions", "Wrong user/context"],
          "how_to_fix": "Check permissions: oc auth can-i <verb> <resource>"
        },
        {
          "type": "Invalid",
          "what_it_means": "Resource definition doesn't match required schema",
          "common_causes": ["Missing required field", "Wrong value type", "Invalid format"],
          "how_to_fix": "Validate YAML against schema: oc explain <resource>"
        }
      ],

      "debugging_strategy": {
        "step_1": "Read the error message completely - don't just skim",
        "step_2": "Identify the error type (Forbidden, NotFound, Invalid, etc.)",
        "step_3": "Look for specific details (resource name, field name, etc.)",
        "step_4": "Map error type to common causes (use table above)",
        "step_5": "Apply the recommended fix for that error type",
        "step_6": "If stuck, copy exact error message to search engine or ask for help"
      }
    }
  }
}
```

**4. Connect Negative Tests to Security**

Explain security implications of proper error handling:

```json
{
  "learning": {
    "security_implications": {
      "why_negative_tests_prevent_vulnerabilities": "Many security issues stem from improper input validation and error handling. Negative tests verify the system defends against malicious inputs.",

      "security_concepts": [
        {
          "concept": "Input Validation",
          "what_it_means": "Checking that user-provided data meets expected format and constraints before using it",
          "vulnerability_without_it": "SQL injection, command injection, XSS, buffer overflows",
          "how_negative_tests_help": "Verify the system rejects malformed, oversized, or malicious inputs"
        },
        {
          "concept": "Fail Securely",
          "what_it_means": "When errors occur, the system should default to denying access, not granting it",
          "vulnerability_without_it": "Authentication bypass, privilege escalation",
          "how_negative_tests_help": "Verify that errors don't accidentally grant permissions or expose data"
        },
        {
          "concept": "Information Disclosure",
          "what_it_means": "Error messages should be helpful but not reveal sensitive system details",
          "vulnerability_without_it": "Attackers learn about system internals (versions, paths, database structure)",
          "how_negative_tests_help": "Verify error messages are informative but don't leak sensitive details"
        }
      ],

      "examples_of_secure_vs_insecure_errors": [
        {
          "scenario": "Failed login attempt",
          "insecure": "Error: Password incorrect for user 'admin'",
          "why_insecure": "Confirms username 'admin' exists, attacker knows username is valid",
          "secure": "Error: Invalid username or password",
          "why_secure": "Doesn't reveal whether username or password was wrong"
        },
        {
          "scenario": "File access denied",
          "insecure": "Error: Permission denied: /etc/sensitive-config.yaml (requires admin role)",
          "why_insecure": "Reveals file path and what permission is needed",
          "secure": "Error: Access denied",
          "why_secure": "Confirms nothing about system structure"
        }
      ]
    }
  }
}
```

---

### Phase 9: Input Validation Education (20-30 minutes)

**REQUIRED:** Add educational content about input validation testing methodologies.

#### Objective

Teach testers to systematically validate all inputs using the comprehensive input validation categories.

#### Educational Content for Input Validation

**1. Explain the Input Validation Matrix**

```json
{
  "concepts": {
    "input_validation_testing": {
      "title": "Comprehensive Input Validation Testing",

      "description": "A systematic approach to testing all possible input scenarios - not just the 'happy path' where everything works, but also edge cases, invalid inputs, and malformed data.",

      "why_it_matters": {
        "quality": "Catches bugs before they reach production - invalid inputs are a leading cause of application crashes",
        "security": "Prevents injection attacks, buffer overflows, and other input-based vulnerabilities",
        "user_experience": "Ensures users get helpful error messages instead of cryptic failures",
        "robustness": "Makes applications resilient to unexpected or malicious inputs"
      },

      "five_categories_explained": {
        "positive": {
          "symbol": "✅",
          "what": "Valid inputs that should succeed",
          "example": "Creating a deployment with image='nginx:latest', replicas=3",
          "purpose": "Verify the system works correctly with valid, well-formed inputs",
          "coverage_target": "Test all documented valid input combinations"
        },

        "negative": {
          "symbol": "❌",
          "what": "Invalid inputs that should fail gracefully",
          "example": "Creating deployment with replicas=-5 (negative number invalid)",
          "purpose": "Verify the system rejects invalid inputs with helpful error messages",
          "coverage_target": "Test common invalid values for each input field"
        },

        "missing": {
          "symbol": "⚠️",
          "what": "Required inputs omitted entirely",
          "example": "Creating deployment without specifying 'image' field",
          "purpose": "Verify the system detects missing required fields and prompts for them",
          "coverage_target": "Test each required field by omitting it individually"
        },

        "corrupt": {
          "symbol": "💥",
          "what": "Malformed data that can't be parsed",
          "example": "YAML file with syntax errors, JSON with unclosed braces",
          "purpose": "Verify the system handles parsing errors gracefully",
          "coverage_target": "Test common formatting mistakes (bad YAML, invalid JSON, etc.)"
        },

        "boundary": {
          "symbol": "📊",
          "what": "Edge cases at minimum/maximum limits",
          "example": "Replicas=0 (minimum), replicas=10000 (maximum), namespace with 63 characters (DNS limit)",
          "purpose": "Verify the system handles limit cases correctly",
          "coverage_target": "Test min, max, and just-over-limit values"
        }
      },

      "how_to_apply_systematically": {
        "step_1": "List all inputs for the operation (image, replicas, namespace, labels, etc.)",
        "step_2": "For each input, identify: valid range, data type, format requirements, constraints",
        "step_3": "Create positive test with valid value",
        "step_4": "Create negative test with invalid value (wrong type, out of range, etc.)",
        "step_5": "Create missing test by omitting required input",
        "step_6": "Create corrupt test with malformed data",
        "step_7": "Create boundary tests for min/max values",
        "result": "5 tests per input field = comprehensive validation coverage"
      },

      "example_complete_matrix": {
        "input_field": "deployment.spec.replicas",
        "valid_range": "0-1000",
        "data_type": "integer",

        "tests": [
          {
            "category": "positive",
            "value": 3,
            "expected": "Success - creates deployment with 3 replicas"
          },
          {
            "category": "negative",
            "value": -5,
            "expected": "Error: replicas must be non-negative"
          },
          {
            "category": "missing",
            "value": null,
            "expected": "Uses default (1 replica) or error if required"
          },
          {
            "category": "corrupt",
            "value": "\"three\"",
            "expected": "Error: replicas must be an integer, got string"
          },
          {
            "category": "boundary",
            "values": [0, 1000, 1001],
            "expected": "0 and 1000 succeed, 1001 fails with 'exceeds maximum'"
          }
        ]
      }
    }
  }
}
```

**2. Add Learning Notes for Each Input Validation Category**

```json
{
  "child_tests": [{
    "title": "Create deployment with valid inputs",
    "input_validation": "positive",

    "learning_note": "Positive tests establish the baseline - what 'correct' looks like. They verify that valid inputs are accepted and processed correctly. Think of this as the reference implementation.",

    "what_makes_this_positive": "All inputs match their expected format, type, and constraints. This represents real-world usage with properly configured values.",

    "learning_objectives": [
      "Understand what constitutes 'valid' input for each field",
      "Learn the expected data types and formats",
      "See how successful operations should behave"
    ]
  }, {
    "title": "Attempt deployment with negative replicas",
    "input_validation": "negative",

    "learning_note": "Negative tests verify the system validates inputs BEFORE trying to use them. Replicas=-5 makes no sense (can't have negative pods), so the system should reject it immediately with a clear error.",

    "why_we_test_invalid_inputs": {
      "prevent_crashes": "Invalid inputs shouldn't crash the system - they should be caught during validation",
      "helpful_errors": "Error messages should explain what's wrong and how to fix it",
      "early_detection": "Catching errors early (at API level) is better than failing later (at runtime)"
    },

    "what_good_validation_looks_like": {
      "immediate_rejection": "API returns error before creating any resources",
      "specific_message": "Error clearly states the problem: 'replicas must be >= 0, got -5'",
      "suggested_fix": "Some systems suggest valid range: 'replicas must be 0-1000'"
    }
  }, {
    "title": "Attempt deployment without required image field",
    "input_validation": "missing",

    "learning_note": "Missing required field tests verify the system enforces completeness. The 'image' field is required - without it, Kubernetes doesn't know what container to run.",

    "why_testing_missing_fields_matters": {
      "schema_enforcement": "Required fields are part of the API contract - system should enforce them",
      "user_guidance": "Error messages should tell users what's missing, not just 'invalid request'",
      "prevent_partial_creation": "Better to fail fast than create incomplete/broken resources"
    },

    "common_missing_field_scenarios": [
      "User forgot to specify the field",
      "Template/script has a variable that's empty",
      "Copy-paste from example but didn't fill in placeholder"
    ]
  }, {
    "title": "Attempt deployment with malformed YAML",
    "input_validation": "corrupt",

    "learning_note": "Corrupt data tests verify the system handles parsing failures gracefully. Malformed YAML/JSON can't be processed, so the system should detect this immediately and explain what's wrong.",

    "types_of_data_corruption": [
      {
        "type": "Syntax errors",
        "example": "Unclosed quotes, missing colons, wrong indentation in YAML",
        "expected_error": "Error parsing YAML: line 5, column 10"
      },
      {
        "type": "Type mismatches",
        "example": "String where integer expected: replicas: \"three\"",
        "expected_error": "Error: cannot unmarshal string into Go value of type int32"
      },
      {
        "type": "Invalid characters",
        "example": "Namespace with spaces or special characters",
        "expected_error": "Error: namespace must match DNS-1123 label format"
      }
    ],

    "why_graceful_parsing_matters": "Parser errors are common developer mistakes. Good error messages that point to the exact line/character save debugging time."
  }, {
    "title": "Create deployment with 0 replicas (minimum boundary)",
    "input_validation": "boundary",

    "learning_note": "Boundary tests explore edge cases at limits. Replicas=0 is the minimum valid value - it's allowed (creates deployment but no pods), but it's an edge case worth testing.",

    "what_are_boundaries": "The minimum, maximum, and just-beyond limits for a value. These are where bugs often hide because they're not frequently tested.",

    "common_boundaries_to_test": [
      {
        "field": "replicas",
        "min": 0,
        "max": "Cluster-dependent (typically 1000-5000)",
        "edge_cases": "0 (special case - scales down completely), 1 (single instance), max+1 (should fail)"
      },
      {
        "field": "namespace",
        "min": "1 character",
        "max": "63 characters (DNS label limit)",
        "edge_cases": "Empty string (invalid), 63 chars (max valid), 64 chars (invalid)"
      },
      {
        "field": "port",
        "min": 1,
        "max": 65535,
        "edge_cases": "0 (invalid), 80 (common), 65535 (max), 65536 (invalid)"
      }
    ],

    "why_boundary_testing_catches_bugs": {
      "off_by_one": "Common bugs: >= vs >, <= vs <",
      "overflow": "Max values can cause integer overflow if not handled",
      "special_handling": "Min/max values sometimes have special code paths that need testing"
    }
  }]
}
```

---

### Phase 10: RBAC Testing Education (20-30 minutes)

**REQUIRED:** Add educational content about role-based access control and permission testing.

#### Objective

Teach testers about Kubernetes RBAC, permission levels, and security boundaries through hands-on testing.

#### Educational Content for RBAC Testing

**1. Explain RBAC Fundamentals**

```json
{
  "concepts": {
    "kubernetes_rbac": {
      "title": "Role-Based Access Control (RBAC) in Kubernetes",

      "description": "RBAC is Kubernetes's security model for controlling who can do what. Instead of giving everyone full access, RBAC lets you grant specific permissions to specific users for specific resources.",

      "real_world_analogy": "Like access levels in a building: lobby access (view), office access (edit), manager office (namespace-admin), building security office (cluster-admin). Each level can access more areas and do more things.",

      "why_rbac_matters": {
        "security": "Prevents unauthorized access and accidental damage - users can only do what they need to do",
        "compliance": "Many regulations require least-privilege access (users have minimal permissions needed for their job)",
        "safety": "Limits blast radius of mistakes - a junior engineer can't accidentally delete production namespaces",
        "multi_tenancy": "Allows multiple teams to share a cluster safely - each team can only see/modify their own resources"
      },

      "four_permission_levels": {
        "view": {
          "icon": "👁️",
          "permissions": "Read-only access to most resources",
          "can_do": ["oc get pods", "oc describe deployment", "oc logs pod-name"],
          "cannot_do": ["oc create/delete/edit anything"],
          "use_case": "Developers reviewing logs, SREs investigating issues, auditors checking compliance",
          "risk_level": "Very low - cannot make changes"
        },

        "edit": {
          "icon": "✏️",
          "permissions": "Create, update, delete namespace-scoped resources",
          "can_do": ["oc create deployment", "oc delete pod", "oc apply -f config.yaml", "oc scale deployment"],
          "cannot_do": ["Create namespaces", "Modify RBAC roles", "Access other namespaces", "Create CRDs"],
          "use_case": "Application developers managing their app's resources",
          "risk_level": "Medium - can modify resources in allowed namespaces"
        },

        "namespace_admin": {
          "icon": "🔐",
          "permissions": "Full control within a namespace including RBAC",
          "can_do": ["Everything 'edit' can do", "Create Roles/RoleBindings in namespace", "Grant permissions to others"],
          "cannot_do": ["Access other namespaces", "Create cluster-scoped resources", "Modify cluster RBAC"],
          "use_case": "Team leads managing team access to their namespace",
          "risk_level": "Medium-high - can grant permissions within namespace"
        },

        "cluster_admin": {
          "icon": "👑",
          "permissions": "Full access to everything in the cluster",
          "can_do": ["EVERYTHING - no restrictions"],
          "cannot_do": ["Nothing - this is God mode"],
          "use_case": "Platform team, cluster operators, emergency break-glass access",
          "risk_level": "Critical - can break entire cluster"
        }
      },

      "how_rbac_works": {
        "components": [
          {
            "name": "User/ServiceAccount",
            "what": "Who is trying to perform the action",
            "example": "User 'alice' or ServiceAccount 'ci-deploy'"
          },
          {
            "name": "Role/ClusterRole",
            "what": "What permissions are defined",
            "example": "Role 'pod-reader' with 'get, list' permissions on 'pods'"
          },
          {
            "name": "RoleBinding/ClusterRoleBinding",
            "what": "Which users have which roles",
            "example": "RoleBinding giving user 'alice' the 'pod-reader' role"
          },
          {
            "name": "Resource",
            "what": "What is being accessed",
            "example": "Pod 'nginx-abc123'"
          },
          {
            "name": "Verb",
            "what": "What action is being attempted",
            "example": "get, list, create, update, delete, watch"
          }
        ],

        "permission_check_flow": "User attempts action → API server checks: Does user have a RoleBinding? → Does that Role allow this verb on this resource? → YES: Allow / NO: Forbidden",

        "visual_example": {
          "scenario": "Alice tries: oc get pods -n dev",
          "check_1": "Does Alice have any RoleBindings in namespace 'dev'? → Yes, 'dev-viewer' RoleBinding",
          "check_2": "What Role does 'dev-viewer' grant? → Role 'view'",
          "check_3": "Does 'view' Role allow 'get' verb on 'pods' resource? → Yes",
          "result": "✅ ALLOWED - Alice can get pods in dev namespace"
        }
      }
    }
  }
}
```

**2. Add Learning Notes for RBAC Test Levels**

```json
{
  "child_tests": [{
    "title": "View user lists pods in namespace",
    "test_type": "rbac",
    "rbac_level": "view",

    "learning_note": "This test uses the 'view' permission level - read-only access. View users can observe resources but not change anything. This is the safest permission level for troubleshooting and auditing.",

    "rbac_education": {
      "permission_level_explained": {
        "level": "view",
        "what_you_can_do": "Get, list, describe, and watch most resources in the namespace",
        "what_you_cannot_do": "Create, update, delete, or modify any resources",
        "real_world_use": "Junior engineers investigating issues, auditors reviewing configurations, SREs reading logs during incidents"
      },

      "why_this_level_appropriate": "Listing pods doesn't change anything - it's purely observational. This is why 'view' permission is sufficient and why granting more (like 'edit') would violate least-privilege principle.",

      "how_to_verify_permissions": {
        "command": "oc auth can-i list pods -n namespace-name --as=view-user",
        "expected": "yes",
        "interpretation": "The '--as' flag simulates being another user. If it returns 'yes', that user has permission."
      }
    },

    "learning_objectives": [
      "Understand view-level permissions and their limitations",
      "Practice checking permissions with 'oc auth can-i'",
      "Learn the least-privilege principle (grant minimum required access)"
    ]
  }, {
    "title": "View user attempts to create deployment (should fail)",
    "test_type": "negative",
    "rbac_level": "view",

    "learning_note": "⚠️ This test SHOULD fail! We're verifying RBAC works by confirming view users CANNOT create resources. The 'Forbidden' error means security is working correctly.",

    "rbac_education": {
      "testing_permission_boundaries": "We don't just test what users CAN do - we test what they CANNOT do. This verifies RBAC actually enforces restrictions.",

      "expected_failure": {
        "error": "Error from server (Forbidden): deployments.apps is forbidden: User \"view-user\" cannot create resource \"deployments\"",
        "why_this_is_good": "RBAC blocked an unauthorized operation - this protects the cluster",
        "parts_explained": {
          "Forbidden": "HTTP 403 - permission denied",
          "User_cannot_create": "Explicitly states what's not allowed",
          "resource_type": "Tells you what resource type was blocked"
        }
      },

      "what_this_teaches": {
        "security_boundaries": "RBAC creates clear security boundaries - view users have read-only access, period",
        "fail_secure": "When in doubt, Kubernetes denies access. This is 'fail secure' design.",
        "troubleshooting": "If you get Forbidden errors, check permissions with 'oc auth can-i'"
      }
    }
  }, {
    "title": "Create deployment with edit permissions",
    "test_type": "install",
    "rbac_level": "edit",

    "learning_note": "This test requires 'edit' permission level - ability to create and modify namespace-scoped resources. Edit is the standard level for developers managing their applications.",

    "rbac_education": {
      "permission_level_explained": {
        "level": "edit",
        "what_you_can_do": "Create, update, delete deployments, pods, services, configmaps, secrets (in assigned namespaces)",
        "what_you_cannot_do": "Create namespaces, modify RBAC roles, create cluster-scoped resources like CRDs",
        "real_world_use": "Application developers deploying their apps, platform users managing their workloads"
      },

      "why_edit_not_admin": "Edit is sufficient for deployment operations. Granting admin would allow modifying RBAC, which developers don't need. Least privilege = use minimum required permissions.",

      "namespace_scoped_security": "Edit permissions are scoped to specific namespaces. You can have 'edit' in 'team-a-dev' namespace but no access to 'team-b-dev' namespace. This enables multi-tenancy."
    }
  }, {
    "title": "Create CustomResourceDefinition (requires cluster-admin)",
    "test_type": "rbac",
    "rbac_level": "cluster-admin",

    "learning_note": "⚠️ CLUSTER-ADMIN REQUIRED: CRDs are cluster-scoped resources that extend Kubernetes's API. Only cluster admins can install them because they affect the entire cluster, not just one namespace.",

    "rbac_education": {
      "permission_level_explained": {
        "level": "cluster-admin",
        "what_you_can_do": "Everything - no restrictions",
        "what_you_cannot_do": "Nothing - this is full cluster control",
        "real_world_use": "Platform team installing operators, cluster maintenance, emergency access"
      },

      "why_cluster_scoped_resources_need_cluster_admin": {
        "impact": "CRDs, ClusterRoles, Nodes, PersistentVolumes affect entire cluster",
        "risk": "Mistakes with cluster-scoped resources can break all namespaces",
        "security": "Cluster-admin can grant themselves access to anything - must be restricted"
      },

      "when_cluster_admin_is_appropriate": [
        "Installing cluster-wide operators",
        "Creating new resource types (CRDs)",
        "Configuring cluster authentication/authorization",
        "Emergency incident response (break-glass access)"
      ],

      "when_cluster_admin_is_inappropriate": [
        "Day-to-day application deployments (use 'edit')",
        "Routine troubleshooting (use 'view')",
        "Automated CI/CD pipelines (use ServiceAccount with specific permissions)",
        "Junior engineers learning (start with 'view', progress to 'edit')"
      ],

      "best_practices": {
        "time_limited": "Grant cluster-admin temporarily for specific tasks, then remove",
        "audit_logged": "Cluster-admin actions should be logged and reviewed",
        "approval_required": "Changes requiring cluster-admin should go through review",
        "just_in_time": "Use temporary elevation systems rather than permanent cluster-admin"
      }
    }
  }]
}
```

---

### Phase 11: Load Testing Education (25-35 minutes)

**REQUIRED:** Add educational content about load testing, performance testing, and resource constraints.

#### Objective

Teach testers about performance testing, resource management, and how systems behave under stress.

#### Educational Content for Load Testing

**1. Explain Load Testing Concepts**

```json
{
  "concepts": {
    "load_testing": {
      "title": "Load Testing and Performance Validation",

      "description": "Load testing verifies how systems behave under stress - high CPU usage, memory pressure, or many concurrent requests. It answers the question: 'Does this still work when resources are constrained or demand is high?'",

      "real_world_analogy": "Like stress-testing a bridge by driving heavy trucks across it. You want to know: At what point does it start to sag? Can it handle rush-hour traffic? What happens if load exceeds design capacity?",

      "why_load_testing_matters": {
        "capacity_planning": "Know how much load your system can handle before adding more resources",
        "performance_regression": "Catch performance degradations before deploying to production",
        "resilience": "Verify system degrades gracefully under load rather than crashing",
        "resource_limits": "Validate that resource limits prevent one app from starving others"
      },

      "types_of_load_tests": {
        "resource_limit_tests": {
          "what": "Test behavior when hitting resource limits (CPU, memory)",
          "example": "Deploy pod with memory limit of 50Mi when app needs 100Mi",
          "what_we_learn": "Does app crash? Get OOMKilled? Slow down gracefully?",
          "impact_type": "modifies-state"
        },

        "load_generation_tests": {
          "what": "Deploy stress-generating pods to consume resources",
          "example": "Run CPU stress pods to max out cluster CPUs",
          "what_we_learn": "How does operator/app behave when cluster is under pressure?",
          "impact_type": "load-generation",
          "safety": "Test environments only - never in production"
        },

        "performance_tests": {
          "what": "Measure response times, throughput, latency under load",
          "example": "Send 1000 requests/second to API, measure average response time",
          "what_we_learn": "Does performance meet SLOs? Where are bottlenecks?",
          "impact_type": "read-only (if just measuring)"
        }
      }
    }
  }
}
```

**2. Add Educational Content for Resource Limit Tests**

```json
{
  "child_tests": [{
    "title": "Deploy pod with insufficient memory limit",
    "test_type": "load",
    "impact_type": "modifies-state",
    "resource_requirements": {
      "memory": "50Mi"
    },

    "learning_note": "This test intentionally sets memory limit below what the application needs. This tests how Kubernetes handles out-of-memory (OOM) conditions - does it kill the pod? Restart it? Notify operators?",

    "load_testing_education": {
      "what_are_resource_limits": {
        "definition": "Maximum amount of CPU/memory a pod can use",
        "why_they_exist": "Prevent one app from consuming all cluster resources and starving others (the 'noisy neighbor' problem)",
        "two_types": {
          "requests": "Minimum guaranteed resources - used for scheduling",
          "limits": "Maximum allowed resources - enforced at runtime"
        }
      },

      "what_happens_when_limit_exceeded": {
        "memory_limit": "Pod gets OOMKilled (Out Of Memory Killed) and may restart",
        "cpu_limit": "Pod gets throttled (slowed down), not killed",
        "why_different": "Memory can't be 'throttled' - you either have enough RAM or you don't. CPU can be time-sliced."
      },

      "how_kubernetes_handles_oom": {
        "step_1": "Pod exceeds memory limit",
        "step_2": "Linux kernel's OOM killer terminates the process",
        "step_3": "Kubernetes sees pod crashed with exit code 137 (SIGKILL)",
        "step_4": "Kubelet restarts pod (if RestartPolicy allows)",
        "step_5": "Pod enters CrashLoopBackOff if it keeps OOM-ing"
      },

      "what_this_test_teaches": [
        "How to set appropriate memory limits (not too low)",
        "What OOMKilled looks like (oc describe pod shows 'OOMKilled' in last state)",
        "How to diagnose memory issues (check metrics, logs before crash)",
        "Importance of monitoring memory usage in production"
      ],

      "real_world_application": "In production, pods occasionally hit memory limits due to traffic spikes or memory leaks. Understanding OOM behavior helps operators diagnose and fix these issues quickly."
    },

    "hands_on_exploration": {
      "observe_the_crash": [
        "Watch pod status: oc get pods -w (should see Running → OOMKilled → CrashLoopBackOff)",
        "Check events: oc describe pod (look for 'OOMKilled' in events)",
        "See exit code: oc get pod -o jsonpath='{.status.containerStatuses[0].lastState.terminated.exitCode}' (should be 137)"
      ],

      "fix_the_problem": [
        "Increase memory limit: edit pod spec to memory: 200Mi",
        "Verify it works: pod should stay Running",
        "Find right-sized limit: monitor actual usage with 'oc adm top pods'"
      ]
    }
  }]
}
```

**3. Add Educational Content for Load Generation Tests**

```json
{
  "child_tests": [{
    "title": "Deploy CPU stress pods to test operator under load",
    "test_type": "load",
    "impact_type": "load-generation",
    "resource_requirements": {
      "cpu": "2",
      "memory": "512Mi",
      "load_generation": true
    },

    "learning_note": "This test deploys 'stress' containers that intentionally consume CPU. We use this to test how the operator behaves when cluster CPUs are maxed out - does it still reconcile? Slow down? Time out?",

    "load_generation_education": {
      "what_is_load_generation": {
        "definition": "Deliberately creating artificial load (CPU, memory, network, disk I/O) to stress-test a system",
        "tools_used": ["stress-ng", "stress", "ab (Apache Bench)", "hey", "k6"],
        "why_artificial_load": "Provides controlled, repeatable stress scenarios without waiting for real production load"
      },

      "how_stress_containers_work": {
        "cpu_stress": "Runs CPU-intensive calculations in tight loops to max out CPU cores",
        "memory_stress": "Allocates large memory blocks to consume RAM",
        "disk_stress": "Performs intensive read/write operations",
        "timeout": "Stress runs for specified duration (e.g., 300s) then stops automatically"
      },

      "what_we_test": {
        "operator_resilience": "Does the operator continue reconciling when CPU is scarce?",
        "api_responsiveness": "Does API server still respond when under CPU pressure?",
        "scheduling": "Can Kubernetes still schedule new pods when cluster is busy?",
        "graceful_degradation": "Does system slow down gracefully or crash hard?"
      },

      "safety_critical": {
        "why_test_only": "Load generation can degrade or crash production systems",
        "blast_radius": "CPU stress affects ALL workloads on the same node",
        "cleanup_essential": "Must delete stress pods when done - they won't stop on their own (unless timeout set)",
        "monitoring": "Watch cluster metrics during test to see impact"
      },

      "interpreting_results": {
        "good_outcome": "Operator slows down but continues working, recovers when stress ends",
        "concerning": "Operator stops reconciling, times out, or gets stuck",
        "bad_outcome": "Operator crashes, cluster becomes unstable",
        "action_items": "If concerning/bad: optimize operator code, increase resource requests, add retries"
      }
    },

    "hands_on_exploration": {
      "monitor_the_load": [
        "Watch node CPU: oc adm top nodes (should show high CPU usage)",
        "Watch pod CPU: oc adm top pods (stress pods should show high usage)",
        "Check if operator still reconciles: oc get customresource -w (should still update)"
      ],

      "measure_impact": [
        "Time an API call before stress: time oc get pods",
        "Time same call during stress: time oc get pods (should be slower)",
        "Compare response times to quantify degradation"
      ]
    }
  }]
}
```

---

### Phase 12: Network Testing Education (25-35 minutes)

**REQUIRED:** Add educational content about network testing, impairment, and resilience validation.

#### Objective

Teach testers about network behavior, failure modes, and how distributed systems handle network issues.

#### Educational Content for Network Testing

**1. Explain Network Testing Concepts**

```json
{
  "concepts": {
    "network_testing": {
      "title": "Network Testing and Resilience Validation",

      "description": "Network testing verifies how systems behave when network conditions degrade - latency spikes, packets drop, connections fail. It answers: 'Does this still work when the network is slow or unreliable?'",

      "why_network_testing_matters": {
        "distributed_systems": "Microservices depend on network - if network fails, services fail",
        "real_world_conditions": "Networks are never perfect - latency varies, packets get lost, connections time out",
        "resilience": "Production systems must handle network issues gracefully (retries, timeouts, circuit breakers)",
        "debugging": "Understanding network failure modes helps diagnose production incidents"
      },

      "types_of_network_impairment": {
        "network_policies": {
          "what": "Kubernetes rules that restrict traffic (like firewall rules)",
          "example": "Deny all ingress traffic to a pod",
          "tests": "Does app handle connection refused? Timeout? Retry?",
          "safety": "High risk - can break connectivity"
        },

        "latency_injection": {
          "what": "Artificially add delay to network packets",
          "example": "Add 100ms delay to all outgoing packets",
          "tests": "Does app handle slow responses? Timeout appropriately?",
          "safety": "Critical risk - affects entire node"
        },

        "packet_loss": {
          "what": "Randomly drop percentage of packets",
          "example": "Drop 10% of packets",
          "tests": "Does TCP retransmit? Does app retry failed requests?",
          "safety": "Critical risk - degrades all traffic"
        },

        "security_group_modification": {
          "what": "AWS security group rules block traffic at infrastructure level",
          "example": "Remove ingress rule for port 443",
          "tests": "Does load balancer health check fail? Does traffic reroute?",
          "safety": "Critical risk - can cause outage"
        }
      },

      "fallacies_of_distributed_computing": {
        "explanation": "Common false assumptions developers make about networks",
        "fallacies": [
          {
            "fallacy": "The network is reliable",
            "reality": "Networks fail, packets get lost, connections drop",
            "how_we_test": "Packet loss tests, connection failure tests"
          },
          {
            "fallacy": "Latency is zero",
            "reality": "Network calls take time - varies from microseconds to seconds",
            "how_we_test": "Latency injection tests"
          },
          {
            "fallacy": "Bandwidth is infinite",
            "reality": "Network bandwidth is limited and shared",
            "how_we_test": "Bandwidth limit tests, concurrent request tests"
          }
        ]
      }
    }
  }
}
```

**2. Add Educational Content for Network Policy Tests**

```json
{
  "child_tests": [{
    "title": "Apply deny-all network policy",
    "test_type": "network",
    "impact_type": "impairment",
    "network_requirements": {
      "impairment_needed": true,
      "impairment_type": "network-policy",
      "network_policies": ["deny-all-policy.yaml"]
    },

    "learning_note": "This test applies a NetworkPolicy that denies ALL traffic to pods in the namespace. This simulates complete network isolation - like unplugging the network cable - to test how apps handle connectivity loss.",

    "network_testing_education": {
      "what_are_network_policies": {
        "definition": "Kubernetes firewall rules that control traffic between pods",
        "default_behavior": "By default, all pods can talk to all pods (allow all)",
        "when_applied": "NetworkPolicies restrict traffic - anything not explicitly allowed is denied",
        "scope": "Apply to pods matching a selector in a specific namespace"
      },

      "deny_all_policy_explained": {
        "what_it_does": "Blocks all ingress (incoming) and egress (outgoing) traffic",
        "effect": "Pod becomes isolated - can't receive requests or make outbound calls",
        "use_cases": [
          "Testing: Verify app handles connection failures gracefully",
          "Security: Isolate compromised pods",
          "Debugging: Eliminate network as variable in troubleshooting"
        ]
      },

      "what_happens_to_apps": {
        "incoming_requests": "Connections time out - clients can't reach the pod",
        "outgoing_requests": "Pod can't reach databases, APIs, other services",
        "dns_resolution": "May fail if DNS server is outside pod (depends on policy)",
        "symptoms": "Connection refused, connection timeout, DNS resolution failures"
      },

      "how_resilient_apps_handle_this": {
        "good": [
          "Detect connection failure quickly (with timeout)",
          "Retry with exponential backoff",
          "Return user-friendly error (Service temporarily unavailable)",
          "Degrade gracefully (serve cached data if possible)"
        ],
        "bad": [
          "Hang indefinitely waiting for response",
          "Crash with unhandled connection error",
          "No retry logic - fail permanently on first error",
          "Expose stack traces to users"
        ]
      },

      "safety_critical": {
        "risk_level": "High - breaks ALL network connectivity for affected pods",
        "impact": "Services become unavailable while policy is in place",
        "duration": "Keep test short - verify behavior then immediately cleanup",
        "monitoring": "Watch for cascading failures - other services may call this one"
      }
    },

    "hands_on_exploration": {
      "observe_the_isolation": [
        "Try to curl the service: oc exec test-pod -- curl http://service (should timeout)",
        "Check pod logs: look for connection errors",
        "Watch metrics: connection errors should spike"
      ],

      "verify_cleanup": [
        "Delete NetworkPolicy: oc delete networkpolicy deny-all",
        "Retry curl: should work now",
        "Confirm: oc get networkpolicy (should show 'No resources found')"
      ]
    }
  }]
}
```

**3. Add Educational Content for Latency Injection Tests**

```json
{
  "child_tests": [{
    "title": "Inject 100ms network latency",
    "test_type": "network",
    "impact_type": "impairment",
    "network_requirements": {
      "impairment_needed": true,
      "impairment_type": "latency",
      "impairment_config": "100ms delay on eth0"
    },

    "learning_note": "This test adds artificial delay to every network packet, simulating a slow network. It helps verify timeout handling, user experience degradation, and retry logic.",

    "latency_injection_education": {
      "what_is_latency": {
        "definition": "Time delay between sending a packet and receiving a response",
        "normal_latency": "1-10ms within same datacenter, 50-200ms cross-continent",
        "why_latency_matters": "Slow responses frustrate users, cause timeouts, reduce throughput",
        "cumulative_effect": "Multiple microservice calls multiply latency (1 request → 10 services → 1000ms total)"
      },

      "how_latency_injection_works": {
        "tool": "tc (traffic control) - Linux kernel packet scheduler",
        "mechanism": "Adds delay to network interface egress queue",
        "scope": "Affects ALL traffic from the node (not just one pod)",
        "command_example": "tc qdisc add dev eth0 root netem delay 100ms",
        "cleanup": "tc qdisc del dev eth0 root"
      },

      "what_100ms_latency_means": {
        "comparison": "About the blink of an eye - barely noticeable for one request",
        "real_impact": "User clicks button → 100ms to backend → 100ms response → 200ms round trip",
        "perception": "Humans notice delays over 100ms as 'sluggish'",
        "timeout_risk": "Default timeouts (30s) won't be hit, but multi-hop calls add up"
      },

      "testing_objectives": {
        "verify_timeouts": "Do requests complete or timeout?",
        "check_retries": "Does app retry on timeout? How many times?",
        "measure_ux_impact": "Is UI still responsive? Do users see loading indicators?",
        "validate_circuit_breakers": "Does circuit breaker open to prevent cascading failures?"
      },

      "what_good_apps_do": {
        "reasonable_timeouts": [
          "Set timeouts appropriate for operation (5s for API call, 30s for file upload)",
          "Don't use infinite timeouts - always have a limit"
        ],
        "retry_with_backoff": [
          "Retry failed requests with exponential backoff",
          "Don't retry immediately - wait longer each attempt"
        ],
        "graceful_degradation": [
          "Show loading indicators to users",
          "Serve stale/cached data if real-time unavailable",
          "Provide partial results instead of complete failure"
        ],
        "circuit_breakers": [
          "Stop sending requests to failing service",
          "Fail fast instead of waiting for timeout",
          "Periodically check if service recovered"
        ]
      },

      "what_bad_apps_do": {
        "no_timeouts": "Wait forever for response - user interface freezes",
        "aggressive_retries": "Retry immediately without backoff - overwhelm already slow service",
        "no_feedback": "User has no idea app is working - thinks it's frozen",
        "cascade_failures": "Slow service causes timeouts upstream, timeouts cascade across system"
      },

      "safety_critical": {
        "risk_level": "Critical - affects ALL traffic on the node, not just test pods",
        "blast_radius": "Every pod on the node experiences latency",
        "monitoring": "Watch for timeout alerts across entire cluster",
        "duration": "Keep test very short (1-2 minutes max)",
        "verification": "Test latency injection on test pod first before applying to node",
        "cleanup": "MUST remove tc qdisc immediately after test - verify with 'tc qdisc show'"
      }
    },

    "hands_on_exploration": {
      "measure_latency": [
        "Before injection: time oc get pods (measure baseline)",
        "Apply latency: tc qdisc add dev eth0 root netem delay 100ms",
        "After injection: time oc get pods (should be ~200ms slower - 100ms each way)",
        "Verify: ping 8.8.8.8 (should show 100ms added to each ping)"
      ],

      "test_application_behavior": [
        "Make API call: curl http://service/api (measure response time)",
        "Check timeout: curl --max-time 1 http://service/api (should timeout if total > 1s)",
        "Watch retries: Check app logs for retry attempts"
      ],

      "cleanup_verification": [
        "Remove: tc qdisc del dev eth0 root",
        "Verify: tc qdisc show dev eth0 (should show only 'pfifo_fast' default)",
        "Confirm: ping 8.8.8.8 (latency should return to normal)"
      ]
    }
  }]
}
```

**4. Add Educational Content for Packet Loss Tests**

```json
{
  "child_tests": [{
    "title": "Inject 10% packet loss",
    "test_type": "network",
    "impact_type": "impairment",
    "network_requirements": {
      "impairment_needed": true,
      "impairment_type": "packet-loss",
      "impairment_config": "10% loss on eth0"
    },

    "learning_note": "This test randomly drops 10% of network packets, simulating unreliable network conditions. TCP will automatically retransmit lost packets, but apps must handle increased latency and potential timeouts.",

    "packet_loss_education": {
      "what_is_packet_loss": {
        "definition": "When network packets fail to reach their destination",
        "causes": "Congested routers, faulty hardware, wireless interference, buffer overflows",
        "normal_rates": "0.1-0.5% on good networks, 5%+ on wifi, 10%+ indicates serious problem",
        "tcp_vs_udp": {
          "tcp": "Automatically retransmits lost packets - increases latency but eventually succeeds",
          "udp": "No retransmission - lost packets are gone forever (used for video streaming where old frames don't matter)"
        }
      },

      "how_packet_loss_injection_works": {
        "tool": "tc (traffic control) with netem (network emulation)",
        "mechanism": "Randomly drops percentage of outgoing packets",
        "scope": "Affects ALL traffic from the node",
        "command_example": "tc qdisc add dev eth0 root netem loss 10%",
        "randomness": "Uses random drop - some connections lucky, others very unlucky",
        "cleanup": "tc qdisc del dev eth0 root"
      },

      "what_10_percent_loss_means": {
        "math": "Out of 100 packets, 10 randomly disappear",
        "tcp_impact": "TCP retransmits after timeout (~200ms-1s), so each lost packet adds delay",
        "worst_case": "Could lose same packet multiple times - exponential backoff means longer delays",
        "throughput_impact": "Effective bandwidth reduced - retransmissions consume capacity",
        "user_experience": "Page loads feel slow, videos buffer, file downloads crawl"
      },

      "testing_objectives": {
        "verify_retransmissions": "Does TCP successfully retransmit and complete requests?",
        "measure_latency_increase": "How much does packet loss slow down operations?",
        "check_timeout_handling": "Do apps timeout appropriately on severe loss?",
        "validate_error_messages": "Do users get clear feedback about network issues?"
      },

      "what_good_apps_do": {
        "rely_on_tcp": "Let TCP handle retransmissions - don't reinvent the wheel",
        "set_realistic_timeouts": "Account for retransmission delays (10s not 1s for network calls)",
        "show_progress": "Display loading indicators - users tolerate delays if they know app is working",
        "handle_eventual_timeout": "Gracefully handle when retries exhaust timeout",
        "log_network_errors": "Log connection issues to help diagnose network problems"
      },

      "what_bad_apps_do": {
        "short_timeouts": "Timeout before TCP can retransmit - fail on first packet loss",
        "no_retry_logic": "Give up immediately instead of letting TCP recover",
        "assume_reliable_network": "No error handling for network failures",
        "poor_error_messages": "Generic 'Error' instead of 'Network connection unstable'"
      },

      "interpreting_results": {
        "good_outcome": "Requests complete successfully but slower (300-500ms instead of 50ms)",
        "acceptable": "Some timeouts but app recovers with retry",
        "concerning": "Frequent timeouts, poor user experience, cascading failures",
        "bad_outcome": "Complete failure to communicate, no recovery",
        "action_items": "If concerning/bad: increase timeouts, add retries, improve error handling"
      },

      "safety_critical": {
        "risk_level": "Critical - degrades ALL traffic on the node",
        "blast_radius": "Every connection from node affected - API servers, etcd, monitoring",
        "cascading_risk": "Could trigger cascading timeouts across cluster",
        "monitoring": "Watch for widespread connection errors and timeout alerts",
        "duration": "Keep test very short (30-60 seconds max)",
        "incident_risk": "10% loss could trigger paging alerts - coordinate with on-call"
      }
    },

    "hands_on_exploration": {
      "observe_packet_loss": [
        "Baseline ping: ping -c 100 8.8.8.8 (note 0% loss)",
        "Apply loss: tc qdisc add dev eth0 root netem loss 10%",
        "Test ping: ping -c 100 8.8.8.8 (should show ~10% packet loss)",
        "Watch retransmissions: netstat -s | grep retransmit (count increasing)"
      ],

      "measure_impact": [
        "File download before: time curl -O http://example.com/largefile (baseline)",
        "File download during: time curl -O http://example.com/largefile (should be slower)",
        "Calculate slowdown: Compare times to quantify impact"
      ],

      "cleanup_verification": [
        "Remove: tc qdisc del dev eth0 root",
        "Verify: ping -c 20 8.8.8.8 (should show 0% loss)",
        "Confirm: tc qdisc show dev eth0 (should show default qdisc only)"
      ]
    }
  }]
}
```

**5. Add Educational Content for Security Group Modification Tests (AWS)**

```json
{
  "child_tests": [{
    "title": "Modify security group to block traffic",
    "test_type": "network",
    "impact_type": "impairment",
    "network_requirements": {
      "impairment_needed": true,
      "impairment_type": "security-group",
      "security_groups": ["sg-abc123"]
    },

    "learning_note": "This test modifies AWS security group rules to block traffic at the infrastructure level, simulating firewall restrictions or misconfigurations. This is more severe than NetworkPolicies as it affects the entire node.",

    "security_group_education": {
      "what_are_security_groups": {
        "definition": "AWS virtual firewall that controls traffic to/from EC2 instances",
        "scope": "Applied at instance level - affects everything on that node",
        "stateful": "If you allow inbound, response is automatically allowed outbound",
        "default_behavior": "Deny all inbound by default, allow all outbound",
        "kubernetes_context": "Each node has security groups controlling access to that machine"
      },

      "how_this_differs_from_network_policies": {
        "network_policies": {
          "level": "Kubernetes/pod level",
          "scope": "Individual pods or namespaces",
          "granularity": "Fine-grained - specific pods, ports, protocols"
        },
        "security_groups": {
          "level": "Infrastructure/EC2 level",
          "scope": "Entire node",
          "granularity": "Coarse - whole machine, all pods on node affected"
        },
        "when_each_matters": "Security groups catch misconfiguration/attacks at infrastructure level; NetworkPolicies provide defense-in-depth within cluster"
      },

      "common_security_group_modifications": {
        "block_specific_port": {
          "example": "Remove rule allowing inbound 443 (HTTPS)",
          "effect": "Load balancer can't reach service, health checks fail",
          "symptom": "503 Service Unavailable errors"
        },
        "block_ssh": {
          "example": "Remove rule allowing inbound 22 (SSH)",
          "effect": "Can't SSH to node for debugging",
          "symptom": "Connection refused when trying to access node"
        },
        "block_api_server": {
          "example": "Remove rule allowing inbound 6443 (Kubernetes API)",
          "effect": "Control plane can't reach node, node goes NotReady",
          "symptom": "Node marked NotReady, pods evicted"
        }
      },

      "testing_objectives": {
        "validate_health_checks": "Do health checks detect blocked traffic and mark node unhealthy?",
        "test_load_balancer": "Does LB stop routing traffic to affected node?",
        "verify_monitoring": "Do alerts fire for unreachable nodes?",
        "check_recovery": "When SG fixed, does node rejoin cluster automatically?"
      },

      "what_good_infrastructure_does": {
        "health_checks": [
          "Load balancer health checks detect unreachable nodes",
          "Stop routing traffic to failed nodes",
          "Automatically resume traffic when node healthy"
        ],
        "monitoring": [
          "Alert when nodes become unreachable",
          "Track security group changes in audit logs",
          "Monitor connection refused errors"
        ],
        "redundancy": [
          "Multiple nodes so one failure doesn't break service",
          "Pods spread across availability zones",
          "Auto-scaling to replace unhealthy nodes"
        ]
      },

      "real_world_incidents": {
        "accidental_sg_change": {
          "scenario": "Engineer removes wrong security group rule during troubleshooting",
          "impact": "All nodes in that security group become unreachable",
          "resolution": "Revert SG change, wait for health checks to restore traffic",
          "prevention": "Review SG changes carefully, use infrastructure-as-code, require approvals"
        },
        "automation_gone_wrong": {
          "scenario": "Automated script modifies SG with incorrect rules",
          "impact": "Production cluster nodes isolated",
          "resolution": "Emergency rollback of automation, manual SG fix",
          "prevention": "Test automation in staging, implement dry-run mode, require human approval for production"
        }
      },

      "safety_critical": {
        "risk_level": "Critical - can cause complete outage",
        "blast_radius": "Affects all pods on nodes in security group",
        "production_risk": "NEVER run in production - staging/test only",
        "backup_plan": "Have console access ready to revert SG changes",
        "communication": "Coordinate with team - this test will trigger alerts",
        "rollback": "Know exact SG rules to restore BEFORE starting test",
        "verification": "Test on single non-critical node first"
      }
    },

    "hands_on_exploration": {
      "prepare_backup": [
        "List current rules: aws ec2 describe-security-groups --group-ids sg-abc123 > sg-backup.json",
        "Document exact rules to restore",
        "Have AWS console open in browser as backup access method"
      ],

      "apply_impairment": [
        "Remove ingress rule: aws ec2 revoke-security-group-ingress --group-id sg-abc123 --protocol tcp --port 443 --cidr 0.0.0.0/0",
        "Verify change: aws ec2 describe-security-groups --group-ids sg-abc123 | grep 443 (should not appear)",
        "Test connectivity: curl https://service (should fail with connection timeout)"
      ],

      "observe_impact": [
        "Check node status: oc get nodes (may show NotReady)",
        "Check load balancer: Check AWS LB health checks (should mark unhealthy)",
        "Check monitoring: Verify alerts fire for unreachable node"
      ],

      "cleanup_and_verify": [
        "Restore rule: aws ec2 authorize-security-group-ingress --group-id sg-abc123 --protocol tcp --port 443 --cidr 0.0.0.0/0",
        "Verify restoration: aws ec2 describe-security-groups --group-ids sg-abc123 | grep 443",
        "Test connectivity: curl https://service (should succeed)",
        "Confirm health: oc get nodes (should show Ready)"
      ]
    }
  }]
}
```

---

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

**Comprehensive Test Type Coverage:**
- [ ] Positive/validation tests include learning objectives
- [ ] Negative tests explain expected failures and error messages
- [ ] Input validation tests cover all 5 categories (positive, negative, missing, corrupt, boundary)
- [ ] RBAC tests explain permission levels and security implications
- [ ] Load tests explain resource constraints and performance impact
- [ ] Network tests explain impairment types and failure modes
- [ ] All test types include safety guidance and risk assessment

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
- [ ] Real-world incident examples provided where relevant
- [ ] Cleanup procedures clearly documented

**Accessibility:**
- [ ] Multi-modal resources (visual, reading, hands-on)
- [ ] Alternative learning paths for different backgrounds
- [ ] Clear prerequisites with verification steps
- [ ] Links to additional help resources

**Safety and Risk Management:**
- [ ] All high-risk tests clearly labeled with risk level
- [ ] Production safety guidelines provided
- [ ] Cleanup verification steps included
- [ ] Blast radius documented for impairment tests
- [ ] Incident coordination guidance provided

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

**Comprehensive Test Type Education:**
- ✅ Negative tests explain error handling and expected failures
- ✅ Input validation tests explain all 5 categories with examples
- ✅ RBAC tests explain Kubernetes permission model
- ✅ Load tests explain resource constraints and performance impact
- ✅ Network tests explain failure modes and resilience patterns
- ✅ All high-risk tests include safety guidance and cleanup verification
- ✅ Real-world incident examples provided where relevant
