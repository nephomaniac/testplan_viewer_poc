# testplan-educator Skill Update Summary

**Date**: 2026-02-19
**Purpose**: Add comprehensive educational content for new test types (Phases 8-12)

---

## Overview

Updated `skills/testplan-educator/skill.md` to add educational content for all comprehensive test types introduced in the testplan-generator skill update (negative testing, input validation, RBAC, load testing, network testing).

---

## Changes Summary

**File**: `skills/testplan-educator/skill.md`
- **Before**: 1,007 lines, 7 phases
- **After**: ~2,100 lines, 12 phases
- **Added**: ~1,100 lines, 5 new phases

---

## New Phases Added

### Phase 8: Negative Testing Education (20-30 minutes)

Teaches testers about negative testing methodology and error handling:

**Educational Content:**
- What is negative testing (crash tests for software)
- Why negative tests matter (error handling, security, user experience)
- Error message literacy (anatomy of error messages, common error types)
- Security implications of error handling
- Real-world examples of negative test scenarios

**Learning Framework:**
```json
{
  "negative_testing_overview": {
    "what_is_negative_testing": "Testing how a system handles invalid inputs...",
    "why_negative_tests_matter": {
      "error_handling": "Real users WILL provide invalid inputs",
      "security": "Many vulnerabilities from improper error handling",
      "user_experience": "Good error messages guide users to fix problems"
    }
  },
  "error_message_literacy": {
    "anatomy_of_error_message": "Error from server (ErrorType): resource...",
    "common_error_types": ["NotFound", "AlreadyExists", "Forbidden", "Invalid"]
  }
}
```

**Test Categories Covered:**
1. Invalid Configuration - Malformed YAML, wrong data types
2. Permission Denials - Insufficient RBAC permissions
3. Resource Conflicts - AlreadyExists, namespace conflicts
4. Missing Dependencies - Prerequisites not met

---

### Phase 9: Input Validation Education (20-30 minutes)

Teaches the 5-category input validation matrix:

**Five Categories Explained:**

1. **Positive (✅)** - Valid inputs that should succeed
   - Purpose: Verify system works correctly with valid inputs
   - Example: Valid image name, correct replica count

2. **Negative (❌)** - Invalid inputs that should fail gracefully
   - Purpose: Verify system rejects invalid inputs with helpful errors
   - Example: Non-existent image, negative replica count

3. **Missing (⚠️)** - Required inputs omitted entirely
   - Purpose: Verify system handles missing required fields
   - Example: No image specified, no namespace

4. **Corrupt (💥)** - Malformed data that can't be parsed
   - Purpose: Verify parser errors are handled gracefully
   - Example: Invalid YAML syntax, wrong data type

5. **Boundary (📊)** - Edge cases at minimum/maximum limits
   - Purpose: Test limits and edge cases
   - Example: 0 replicas, 10,000 replicas, 63-character namespace

**Educational Content:**
- Input validation matrix with examples
- Systematic testing approach
- Common validation mistakes
- Real-world input validation failures
- Hands-on exploration suggestions

---

### Phase 10: RBAC Testing Education (20-30 minutes)

Teaches Kubernetes RBAC (Role-Based Access Control):

**Four Permission Levels:**

1. **View (👁️)** - Read-only access
   - Permissions: Get, list, watch resources
   - Risk level: Very low - cannot make changes
   - Example operations: oc get pods, oc describe deployment

2. **Edit (✏️)** - Create/update/delete namespace resources
   - Permissions: Create, update, delete namespace-scoped resources
   - Risk level: Medium - can modify resources in allowed namespaces
   - Example operations: oc create deployment, oc delete pod

3. **Namespace-Admin (🔐)** - Full control within namespace
   - Permissions: Full control including RBAC within namespace
   - Risk level: Medium-high - can grant permissions to others
   - Example operations: Create Roles, RoleBindings

4. **Cluster-Admin (👑)** - Full cluster access
   - Permissions: Unlimited - can do anything in cluster
   - Risk level: Critical - can break entire cluster
   - Example operations: Create CRDs, ClusterRoles, modify nodes

**Educational Content:**
- RBAC fundamentals (who can do what to which resources)
- Security best practices (principle of least privilege)
- Common RBAC mistakes
- Real-world RBAC incidents
- Testing permission boundaries

---

### Phase 11: Load Testing Education (25-35 minutes)

Teaches load testing and performance validation:

**Three Load Test Types:**

1. **Resource Limit Tests**
   - What: Test behavior when hitting resource limits (CPU, memory)
   - What we learn: Does app crash? Get OOMKilled? Slow down gracefully?
   - Example: Deploy pod with 128Mi memory limit, try to allocate 256Mi

2. **Load Generation Tests**
   - What: Deploy stress-generating pods to consume resources
   - Safety: Test environments only - never in production
   - Tools: stress-ng, sysbench, Apache Bench
   - Example: CPU stress pods, memory stress pods

3. **Performance Tests**
   - What: Measure and validate performance metrics
   - Metrics: Response time, throughput, error rate
   - Baselines: Establish performance baselines for comparison

**Educational Content:**
- Load testing concepts (stress testing vs load testing)
- Resource limits in Kubernetes (requests vs limits)
- What happens when limits are exceeded (throttling, OOMKilled)
- Interpreting results (good/concerning/bad outcomes)
- Hands-on exploration (monitor load, measure impact)

---

### Phase 12: Network Testing Education (25-35 minutes)

Teaches network testing and resilience validation:

**Four Network Impairment Types:**

1. **Network Policy Tests**
   - What: Kubernetes firewall rules restricting traffic
   - Example: Deny-all policy blocks all ingress/egress
   - Tests: Does app handle connection failures gracefully?
   - Safety: High risk - can break connectivity

2. **Latency Injection**
   - What: Artificially add delay to network packets
   - Example: Add 100ms delay using tc qdisc
   - Tests: Does app handle slow responses? Timeout appropriately?
   - Safety: Critical risk - affects entire node

3. **Packet Loss**
   - What: Randomly drop percentage of packets
   - Example: Drop 10% of packets
   - Tests: Does TCP retransmit? Does app retry failed requests?
   - Safety: Critical risk - degrades all traffic

4. **Security Group Modification (AWS)**
   - What: AWS security group rules block traffic at infrastructure level
   - Example: Remove ingress rule for port 443
   - Tests: Does load balancer health check fail? Does traffic reroute?
   - Safety: Critical risk - can cause outage

**Educational Content:**
- Network testing concepts
- Fallacies of distributed computing (network is reliable, latency is zero, etc.)
- How resilient apps handle network issues (retries, timeouts, circuit breakers)
- Network impairment tools (tc qdisc, NetworkPolicy, AWS security groups)
- Real-world network incidents
- Safety critical guidance for each impairment type

---

## Updated Sections

### Quality Checklist

Added new checkboxes for comprehensive test types:

**Comprehensive Test Type Coverage:**
- [ ] Positive/validation tests include learning objectives
- [ ] Negative tests explain expected failures and error messages
- [ ] Input validation tests cover all 5 categories (positive, negative, missing, corrupt, boundary)
- [ ] RBAC tests explain permission levels and security implications
- [ ] Load tests explain resource constraints and performance impact
- [ ] Network tests explain impairment types and failure modes
- [ ] All test types include safety guidance and risk assessment

**Safety and Risk Management:**
- [ ] All high-risk tests clearly labeled with risk level
- [ ] Production safety guidelines provided
- [ ] Cleanup verification steps included
- [ ] Blast radius documented for impairment tests
- [ ] Incident coordination guidance provided

### Success Criteria

Added comprehensive test type education criteria:

**Comprehensive Test Type Education:**
- ✅ Negative tests explain error handling and expected failures
- ✅ Input validation tests explain all 5 categories with examples
- ✅ RBAC tests explain Kubernetes permission model
- ✅ Load tests explain resource constraints and performance impact
- ✅ Network tests explain failure modes and resilience patterns
- ✅ All high-risk tests include safety guidance and cleanup verification
- ✅ Real-world incident examples provided where relevant

---

## Educational Content Pattern

Each new phase follows a consistent educational framework:

### 1. Conceptual Overview
- What is this test type?
- Why does it matter?
- Real-world analogies

### 2. Educational Deep Dive
- How does it work technically?
- What are the different categories/types?
- Common misconceptions

### 3. Learning Notes for Each Test
- What you're testing
- What good apps do vs bad apps
- Interpreting results
- Security/safety implications

### 4. Real-World Examples
- Real incidents and outages
- Lessons learned
- Prevention strategies

### 5. Hands-On Exploration
- Commands to run
- What to observe
- How to verify
- Cleanup procedures

### 6. Safety Guidance
- Risk level assessment
- Blast radius documentation
- Production safety guidelines
- Cleanup verification

---

## Benefits

### For Learners
- **Progressive learning** from basic concepts to advanced techniques
- **Real-world context** connecting theory to practice
- **Safety-first approach** with clear risk assessment
- **Hands-on exploration** encouraging experimentation
- **Error literacy** understanding what failures mean

### For Test Executors
- **Clear understanding** of what each test validates
- **Safety awareness** knowing risks before execution
- **Troubleshooting guidance** when tests fail
- **Cleanup procedures** for safe test execution

### For System Quality
- **Comprehensive coverage** across all test types
- **Security awareness** understanding permission boundaries
- **Resilience validation** testing failure modes
- **Performance baselines** establishing performance expectations

---

## Files Modified

1. **skills/testplan-educator/skill.md**
   - Added 1,100 lines
   - Added Phases 8-12
   - Updated Quality Checklist
   - Updated Success Criteria

---

## Next Steps

### Immediate
1. ✅ Update testplan-educator skill (complete)
2. ⏳ Update testplan-executor skill (Phase 3 - pending)

### Short Term
1. ⏳ Migrate CAMO example to new format
2. ⏳ Migrate RHOBS-Next example to new format
3. ⏳ Generate comprehensive tests for both examples

### Documentation
1. ⏳ Update examples/README.md with filtering examples
2. ⏳ Update skills/README.md with new test type descriptions
3. ⏳ Create comprehensive test generation guide

### Skills Review (NEW)
1. ⏳ Review test-case-generator skill from MCPmarket
2. ⏳ Review test-driven-development skill from superpowers repo
3. ⏳ Determine if test case reviewer skill would be beneficial

---

## Complete Phase Sequence

The testplan-educator skill now follows this comprehensive workflow:

1. **Artifact Analysis** - Explore primary artifact
2. **Prerequisites** - Explain prerequisites with verification
3. **Difficulty Assessment** - Assess and explain difficulty levels
4. **Dependency Education** - Explain test dependencies
5. **Test Case Enhancement** - Add learning notes to tests
6. **Conceptual Foundation** - Add conceptual overviews
7. **Learning Path Creation** - Create progressive learning paths
8. **Negative Testing Education** ✨ NEW - Error handling education
9. **Input Validation Education** ✨ NEW - 5-category validation matrix
10. **RBAC Testing Education** ✨ NEW - Permission model education
11. **Load Testing Education** ✨ NEW - Performance and resource education
12. **Network Testing Education** ✨ NEW - Resilience and failure modes
13. **Quality Checklist** - Final validation

---

## Integration with Existing Workflow

The new phases integrate seamlessly:

**Execution Flow:**
1. Phases 1-7 enhance base tests with educational content
2. **Phase 8-9** add negative testing and input validation education
3. **Phase 10** adds RBAC security education
4. **Phase 11-12** add load and network resilience education
5. Quality checklist validates completeness

**Backward Compatibility:**
- Skill can still enhance tests without Phases 8-12 if needed
- Each phase is independent and can be skipped if not applicable
- Legacy test plans continue to work

---

## Verification

To verify the skill update:

1. **Test with existing example:**
   ```
   Use testplan-educator skill to enhance examples/rhobs-next/rhobs-next-testplan.json
   Verify it adds educational content for all test types
   Check all phases are present
   ```

2. **Check educational content:**
   - Should have conceptual overviews for all test types
   - Should explain error messages and failure modes
   - Should have hands-on exploration suggestions
   - Should have safety guidance for high-risk tests

3. **Validate format:**
   - All tests have learning notes
   - All concepts explained in plain language
   - All phases follow consistent pattern
   - Quality checklist validates all test types

---

**Status**: testplan-educator skill updated with comprehensive educational content ✅
**Next**: Review skills from MCPmarket and superpowers repo, then update testplan-executor skill
