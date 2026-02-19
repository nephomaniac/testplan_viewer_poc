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

Edit your test plan JSON file (e.g., `examples/camo/camo-testplan.json`):

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
- ✅ You can execute tests confidently
- ✅ You understand the system under test
- ✅ You can troubleshoot common issues
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
localStorage.getItem('testplan_progress')

# Reset if corrupted
localStorage.removeItem('testplan_progress')
```

## Next Steps After Completion

### 🎓 You've Learned
- System architecture and components
- Testing patterns and workflows
- Troubleshooting techniques
- Best practices

### 🚀 What's Next
1. **Apply to real scenarios** - Use knowledge in production
2. **Explore advanced topics** - Dive deeper into complex areas
3. **Contribute back** - Add your learnings to test plan
4. **Create new test plans** - Document other systems

### 📚 Further Reading
- See project-specific documentation
- Explore operator/system source code
- Review related test plans in examples/

---

**Time Investment:** Varies by test plan
**Outcome:** Understanding of system under test
**Next:** Apply knowledge to real scenarios

🎉 **Happy Learning!**
