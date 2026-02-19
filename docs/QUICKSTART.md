# Quick Start Guide

## Generate Interactive HTML in 30 Seconds

```bash
# 1. Build the tool (one time)
go build -o testplan-viewer

# 2. Generate HTML
./testplan-viewer

# 3. Open in browser
open testplan.html
```

That's it! 🎉

## What You'll See

### 📊 Progress Dashboard
- Visual progress bar tracking completion
- Filter by difficulty (Beginner/Intermediate/Advanced)
- Search all content

### 📚 Interactive Learning Path
- **Test 1: Create Your First HCP Cluster** (Beginner, 15-20 min)
  - Learn ROSA HCP creation
  - Understand provision shards
  - Practice OCM CLI

- **Test 2: Trace Logs from Application to Loki** (Intermediate, 30-40 min)
  - Master LogQL queries
  - Debug log forwarding
  - Explore Grafana

- **Test 3: Watch Synthetic Monitoring Detect Your Cluster** (Advanced, 45-60 min)
  - Understand distributed monitoring
  - Trace probe lifecycle
  - Multi-cluster debugging

### ✨ Interactive Features

**For Each Test:**
- ✅ Learning objectives (what you'll learn)
- ⚠️ Common beginner mistakes (what to avoid)
- 💡 Learning notes (why each step matters)
- 📋 Copy-to-clipboard commands
- 🚨 Expected errors with solutions
- 📊 Progress checkboxes

## Example Workflow

### Day 1: Beginner Test (20 minutes)

1. Open `testplan.html`
2. Filter: "Beginner" difficulty
3. Expand "Test 1: Create Your First HCP Cluster"
4. Read learning objectives
5. Follow steps 1-4:
   - Copy commands with one click
   - Read "learning notes" to understand why
   - Check off completed steps
6. Mark test complete when done
7. Progress bar shows 33% ✅

### Day 2: Intermediate Test (40 minutes)

1. Reopen `testplan.html`
2. Progress automatically restored (33%)
3. Expand "Test 2: Trace Logs..."
4. Work through log pipeline
5. Learn LogQL basics
6. Mark complete → 66% ✅

### Day 3: Advanced Test (60 minutes)

1. Complete "Test 3: Synthetic Monitoring"
2. Progress: 100% 🎉
3. Explore "Bonus Exploration" sections
4. You're now onboarded! 🎓

## Features Demo

### Search
```
Search: "logs"
→ Highlights all tests/steps related to logging
→ Jump directly to relevant content
```

### Difficulty Filter
```
Filter: "Beginner"
→ Shows only Test 1
→ Perfect for Day 1
```

### Copy Commands
```
Hover over any command block
→ "Copy" button appears
→ One click to clipboard
→ Paste into terminal
```

### Progress Tracking
```
Check a step → Saved to browser
Close browser → Progress persists
Reopen next day → Pick up where you left off
```

## Customization

### Change Test Plan Content

Edit `rhobs_test_plan_v2.json`:

```json
{
  "testcases": {
    "my_new_test": {
      "metadata": {
        "difficulty": "intermediate",
        "estimated_time": "30 minutes"
      },
      "learning": {
        "objectives": ["Learn X", "Understand Y"]
      }
    }
  }
}
```

Then regenerate:
```bash
./testplan-viewer
```

### Add New Learning Path

```json
{
  "learning_path": {
    "sequence": ["test_1", "test_2", "my_new_test", "test_3"],
    "alternative_paths": {
      "fast_track": ["test_1", "test_3"]
    }
  }
}
```

## Tips for Learners

### 📌 Best Practices
1. **Follow the sequence** - Tests build on each other
2. **Read learning notes** - Understand why, not just what
3. **Don't skip errors** - Learn from expected failures
4. **Take breaks** - Spread over 2-3 days
5. **Explore bonuses** - Deepen understanding

### ⚠️ Common Pitfalls
- **Skipping pre-requisites** → Set up access first
- **Rushing through** → Read learning notes
- **Ignoring errors** → Expected errors teach debugging
- **Not tracking progress** → Check boxes to stay motivated

### 🎯 Success Signals
- ✅ You can create HCP clusters confidently
- ✅ You understand the observability data flows
- ✅ You can query Loki and RHOBS
- ✅ You can debug synthetic monitoring issues
- ✅ You know where to get help

## Troubleshooting

### HTML Not Opening?
```bash
# Try explicit browser
firefox testplan.html
# or
google-chrome testplan.html
```

### Progress Not Saving?
- Check browser allows localStorage
- Clear browser cache and retry
- Try incognito mode (warning: won't persist)

### Commands Failing?
- Verify all prerequisites are met
- Check environment variables are set
- Read the "Common Errors" section for that step
- Each step includes troubleshooting guidance

### Lost Progress?
```bash
# Check localStorage in browser console
localStorage.getItem('rhobs_testplan_progress')

# Reset if corrupted
localStorage.removeItem('rhobs_testplan_progress')
```

## Next Steps After Completion

### 🎓 You've Learned
- ROSA HCP cluster lifecycle
- Observability data flows (metrics, logs, probes)
- LogQL and PromQL basics
- Multi-cluster debugging
- RHOBS architecture

### 🚀 What's Next
1. **Shadow experienced SRE** - Watch them debug real issues
2. **Take on-call shift** - With mentorship
3. **Explore advanced topics** - Custom metrics, alerting rules
4. **Contribute back** - Add your learnings to test plan

### 📚 Further Reading
- RHOBS Architecture Docs
- Route Monitor Operator Code
- OpenShift Logging Documentation
- Synthetic Monitoring Design Proposal

---

**Time Investment:** 2-3 hours hands-on + breaks
**Outcome:** Full understanding of RHOBS observability stack
**Next:** Apply knowledge to real production issues

🎉 **Happy Learning!**
