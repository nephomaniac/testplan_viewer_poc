# RHOBS-Next Test Plan Structure Fixes

## Summary

Fixed all structure mismatches in `rhobs-next-testplan.json` to match the Go data model at `/Users/maclark/sandbox/testplan_tools_poc/internal/models/testplan.go`.

## Changes Made

### 1. ✅ Metadata Structure
- **Changed**: Renamed `last_updated` field to `date`
- **Before**: `"last_updated": "2026-02-18"`
- **After**: `"date": "2026-02-18"`
- **Reason**: Go model expects `date` field, not `last_updated`

### 2. ✅ References Structure in All Test Cases
- **Changed**: Transformed references from array of objects to categorized URL arrays
- **Before**:
  ```json
  "references": [
    {
      "title": "Route Monitor Operator GitHub",
      "url": "https://github.com/openshift/route-monitor-operator",
      "type": "documentation"
    }
  ]
  ```
- **After**:
  ```json
  "references": {
    "documentation": [
      "https://github.com/openshift/route-monitor-operator"
    ],
    "related_code": [],
    "sops": []
  }
  ```
- **Mapping Logic**:
  - `type: "documentation"`, `"tutorial"`, `"guide"` → `documentation[]`
  - `type: "code"`, `"repository"`, `"related_code"` → `related_code[]`
  - `type: "sop"`, `"runbook"`, `"procedure"` → `sops[]`
  - Default → `documentation[]`

### 3. ✅ Removed Conceptual Overview Field
- **Removed**: `conceptual_overview` field from all test cases
- **Reason**: Not present in Go model's TestCase struct
- **Note**: This field contained valuable learning content but didn't match the expected structure

### 4. ✅ Added Progress Tracking
- **Added**: `progress_tracking` top-level section
- **Structure**:
  ```json
  "progress_tracking": {
    "description": "Track your progress through the RHOBS Next test plan",
    "local_storage_key": "rhobs_next_testplan_progress",
    "tracked_items": [
      "test_verify_route_monitor_deployment",
      "test_create_routemonitor_cr",
      "test_cleanup_routemonitor_cr",
      "test_verify_probe_metrics",
      "test_create_probe_via_api"
    ]
  }
  ```
- **Reason**: Required by Go model's TestPlan struct

## Test Cases Updated

All 5 test cases had their references transformed:
1. `test_verify_route_monitor_deployment` - 4 documentation refs
2. `test_create_routemonitor_cr` - 4 documentation refs
3. `test_cleanup_routemonitor_cr` - 4 documentation refs
4. `test_verify_probe_metrics` - 5 documentation refs
5. `test_create_probe_via_api` - 6 documentation refs

## Validation Results

✅ All required top-level keys exist:
- metadata
- learning_path
- prerequisites
- concepts
- testcases
- progress_tracking

✅ All test cases have correct structure:
- metadata
- learning
- test_execution
- troubleshooting
- next_steps
- references (with documentation, related_code, sops arrays)

✅ JSON syntax is valid and parseable with jq

## Scripts Created

1. `fix_json_structure.py` - Main transformation script
2. `restore_references.py` - Script to restore references from individual test case files
3. `validate_structure.sh` - Comprehensive validation script

## Next Steps

The JSON file is now ready to parse correctly with the Go parser at `/Users/maclark/sandbox/testplan_tools_poc/internal/parser/parser.go`.
