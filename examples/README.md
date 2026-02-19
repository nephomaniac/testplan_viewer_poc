# Test Plan Examples

This directory contains example test plan JSON files for the RHOBS Test Plan Viewer.

## Available Examples

### rhobs_test_plan_v2.json
- **Purpose**: Learning-optimized test plan for RHOBS onboarding
- **Target Audience**: New engineers joining the RHOBS team
- **Test Count**: 3 comprehensive tests
- **Features**: Beginner-friendly with common mistakes highlighted, detailed learning objectives

### rhobs_test_plan.json
- **Purpose**: Original test plan structure (v1 format)
- **Target Audience**: Testing and comparison purposes
- **Test Count**: Multiple test scenarios

### RHOBS_Manual_Test_Plan.md
- **Purpose**: Markdown version of the manual test plan
- **Format**: Human-readable markdown documentation

## Adding New Test Plan Examples

1. **Create your JSON file** following the test plan schema (see below)
2. **Place it in this directory**: `examples/your_testplan.json`
3. **Test it** with the viewer:
   ```bash
   ./build/testplan-viewer -i examples/your_testplan.json -o output.html
   ```
4. **Document it** by adding an entry to this README

## Minimal Valid Test Plan

Here's the minimum structure required for a valid test plan:

```json
{
  "metadata": {
    "document_title": "My Test Plan",
    "version": "1.0",
    "last_updated": "2026-02-18",
    "status": "draft",
    "epic": "SREP-1234",
    "purpose": "Test plan purpose",
    "target_audience": "Engineers",
    "estimated_total_time": "1 hour",
    "learning_objectives": [
      "Learn objective 1"
    ]
  },
  "learning_path": {
    "description": "Basic learning path",
    "sequence": ["test_1"]
  },
  "testcases": {
    "test_1": {
      "metadata": {
        "id": "test_1",
        "name": "test_1",
        "title": "First Test",
        "category": "setup",
        "difficulty": "beginner",
        "estimated_time": "30 minutes",
        "hands_on_percentage": 70
      },
      "learning": {
        "objectives": ["Learn something useful"],
        "concepts_covered": ["Basic concept"],
        "skills_gained": ["Basic skill"]
      },
      "test_execution": {
        "objective": "Test objective description",
        "duration": "30 minutes",
        "dependencies": [],
        "prerequisites": [],
        "steps": [
          {
            "step_number": 1,
            "title": "First step",
            "command": "echo 'Hello World'",
            "expected_output": "Hello World"
          }
        ],
        "validation": {
          "success_criteria": ["Step completed successfully"],
          "what_success_looks_like": "You see the expected output",
          "what_failure_looks_like": "Error messages appear"
        }
      }
    }
  }
}
```

## JSON Schema Reference

### Required Top-Level Fields
- `metadata` - Document metadata
- `learning_path` - Learning path configuration
- `testcases` - Map of test case IDs to test case objects

### Metadata Fields
- `document_title` (required) - Title of the test plan
- `version` - Version string
- `last_updated` - ISO date string
- `status` - Status (draft, review, approved, etc.)
- `epic` - Associated JIRA epic
- `purpose` - Purpose of this test plan
- `target_audience` - Who should use this
- `estimated_total_time` - Total time estimate
- `learning_objectives` - Array of learning objectives

### Learning Path Fields
- `description` - Description of the learning path
- `sequence` - Array of test IDs in recommended order
- `alternative_paths` - Optional alternative paths for different audiences

### Test Case Structure
Each test case must include:
- `metadata` - Test metadata (id, title, difficulty, time, etc.)
- `learning` - Learning objectives and concepts
- `test_execution` - Execution details (steps, validation, etc.)

Optional sections:
- `troubleshooting` - Common issues and solutions
- `next_steps` - What to do after completing the test
- `references` - Related documentation

## Validation

All test plans are automatically validated when you generate HTML. The validator checks:

1. **Metadata validation**
   - `document_title` is required
   - At least one test case must exist

2. **Learning path validation**
   - All test IDs in `sequence` must exist in `testcases`
   - All test IDs in `alternative_paths` must exist

3. **Dependency validation**
   - All test dependencies must reference existing tests
   - No circular dependencies

To test validation:

```bash
# This should succeed
./build/testplan-viewer -i examples/rhobs_test_plan_v2.json -o test.html

# This will fail with validation errors
./build/testplan-viewer -i invalid.json -o test.html
```

## Tips for Creating Great Test Plans

1. **Start simple** - Begin with one or two tests, then expand
2. **Focus on learning** - Include learning objectives and "why this step" explanations
3. **Add context** - Include troubleshooting and common beginner mistakes
4. **Test thoroughly** - Run through your test plan before sharing
5. **Include examples** - Show expected output for every command
6. **Link resources** - Add references to related documentation

## Need Help?

- See the full schema in `docs/README.md`
- Look at `rhobs_test_plan_v2.json` for a complete example
- Check `docs/QUICKSTART.md` for getting started
