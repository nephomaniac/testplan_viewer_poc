# Skill Self-Verification System

All test plan skills include a built-in self-verification mechanism to ensure they're always using the latest skill definition from the repository.

## Why Self-Verification?

**Problem:** Skills can drift from their source of truth
- User might be working with an old version
- Skills might be updated in the repo but not locally
- Remote repo might have newer versions
- Uncommitted local changes might exist

**Solution:** Each skill verifies itself before execution
- Checks local vs committed version
- Checks local vs remote version
- Prompts user if differences found
- Offers options to resolve

## How It Works

### Automatic Verification

When you invoke a skill, Claude will:

1. **Check file existence**
   ```bash
   ls -la skills/[skill-name]/skill.md
   ```

2. **Verify git repository**
   ```bash
   git rev-parse --is-inside-work-tree
   ```

3. **Check for local changes**
   ```bash
   git status skills/[skill-name]/skill.md
   git diff skills/[skill-name]/skill.md
   ```

4. **Compare with remote** (if applicable)
   ```bash
   git fetch origin main
   git diff HEAD origin/main -- skills/[skill-name]/skill.md
   ```

### User Prompt

If differences are detected, you'll see:

```
⚠️ Skill Definition Verification

I've detected differences in the [skill-name] skill definition:

Changes detected:
- 15 lines added
- 3 lines removed
- Last commit: 2 hours ago
- Remote has newer version

Options:
1. Continue with current version (may be outdated)
2. Read latest version from repository and use that
3. Show me the full diff to review
4. Cancel and let me update the skill first

What would you like to do?
```

### Your Options

#### Option 1: Continue with current version
- Uses whatever skill definition Claude has in context
- May be outdated or different from repo
- Claude will warn results may differ from spec

**Choose when:**
- You know current version is correct
- Testing local changes
- Don't care about latest updates

#### Option 2: Read latest version from repository
- Claude reads `skills/[skill-name]/skill.md` from disk
- Uses that as the skill definition
- Ensures latest local version is used

**Choose when:**
- Want to use latest committed/pulled version
- Recently pulled from remote
- Trust local repo is current

#### Option 3: Show me the full diff
- Claude displays full git diff
- You can review what changed
- Then asked again which option to choose

**Choose when:**
- Want to understand what changed
- Reviewing skill updates
- Deciding if changes matter for your use case

#### Option 4: Cancel and update first
- Stops skill execution
- Claude provides update commands
- You can pull latest and re-run

**Choose when:**
- Want to ensure you have absolute latest
- Remote has updates you need
- Want to commit/stash local changes first

## Verification Output

After verification completes:

```
✅ Skill Verification Complete

Skill: testplan-generator
Version: abc123f (git commit hash)
Last Modified: 2026-02-18 14:30:00
Status: Up to date ✓

Proceeding with skill execution...
```

## When Verification is Skipped

Verification is only skipped if:

1. **Explicitly requested:** User says "skip verification"
2. **Already verified:** Verified earlier in same conversation
3. **Emergency:** User states time-critical situation

Otherwise, verification ALWAYS runs first.

## Example Scenarios

### Scenario 1: Everything up to date

```
User: "Use testplan-generator skill to create a test plan"

Claude:
✅ Skill Verification Complete

Skill: testplan-generator
Version: f21b0dd
Last Modified: 2026-02-18 15:00:00
Status: Up to date ✓

Proceeding with skill execution...
```

No prompt needed - skill executes immediately.

---

### Scenario 2: Local uncommitted changes

```
User: "Use testplan-educator skill"

Claude:
⚠️ Skill Definition Verification

I've detected uncommitted changes in the testplan-educator skill:

Changes:
M skills/testplan-educator/skill.md (modified, not staged)

You have local edits that aren't committed. These may be:
- Intentional improvements you're testing
- Accidental edits
- Work in progress

Options:
1. Continue with modified version
2. Use committed version (ignore local changes)
3. Show me the diff
4. Cancel and let me commit/discard changes

What would you like to do?
```

---

### Scenario 3: Remote has updates

```
User: "Use testplan-executor skill"

Claude:
⚠️ Skill Definition Verification

Remote repository has a newer version:

Remote updates:
+ Added new failure pattern detection
+ Improved gap analysis reporting
+ Fixed bug in result validation

Last remote commit: 1 hour ago
Your local: 2 days ago

Options:
1. Continue with local version (older)
2. Pull updates and use latest (recommended)
3. Show full diff
4. Cancel - I'll pull manually

What would you like to do?
```

---

### Scenario 4: Complex situation (local changes + remote updates)

```
User: "Use testplan-generator skill"

Claude:
⚠️ Skill Definition Verification

Complex situation detected:
- You have uncommitted local changes
- Remote has newer version
- Changes may conflict

Recommendation: Resolve git state before proceeding

Suggested steps:
1. Review your local changes: git diff skills/testplan-generator/skill.md
2. Decide: keep, commit, or discard
3. Pull remote updates: git pull origin main
4. Resolve any conflicts
5. Re-run skill

Options:
1. Continue with local modified version (ignoring remote)
2. Show me both diffs (local changes + remote changes)
3. Cancel - I'll fix git state first

What would you like to do?
```

## Best Practices

### For Users

**Do:**
- Let verification run (it's quick)
- Pull latest before starting
- Commit/stash local changes
- Review diffs when prompted

**Don't:**
- Skip verification habitually
- Ignore update prompts
- Work with outdated skills
- Make skill changes without committing

### For Skill Developers

**When updating skills:**

1. Make changes in local file
2. Test with Claude
3. Commit changes with clear message
4. Push to remote
5. Document what changed

**Version tracking:**

- Each commit hash becomes a version
- File modification time shows recency
- Git log shows change history

## Troubleshooting

### "Git command not found"

If not in a git repository:
```bash
git init
git add skills/
git commit -m "Initialize skills"
```

### "Remote not configured"

If remote isn't set up:
```bash
git remote add origin <repository-url>
git fetch origin
```

### "Permission denied"

If can't fetch from remote:
```bash
# Check remote URL
git remote -v

# Update credentials if needed
git config credential.helper store
```

### "Conflicts detected"

If local and remote conflict:
```bash
# Option 1: Keep local, ignore remote
git stash
git pull origin main
git stash pop
# Resolve conflicts

# Option 2: Discard local, use remote
git checkout -- skills/
git pull origin main
```

## Architecture

### Verification Flow

```
Skill Invocation
     ↓
Check Local File Exists
     ↓
Git Repo Check
     ↓
Local vs Committed Diff
     ↓
Local vs Remote Diff
     ↓
   ┌─── No Differences ────┐
   │                       │
   ✓                       ↓
Proceed            Differences Detected
                          ↓
                   Prompt User (4 options)
                          ↓
              ┌────────────┼────────────┐
              ↓            ↓            ↓
          Option 1     Option 2     Option 3
          Continue     Use Latest   Show Diff
              ↓            ↓            ↓
          Proceed      Read File    Display → Prompt Again
                           ↓
                       Proceed
```

### File References

Each skill's verification section references:
- **Local file:** `skills/[skill-name]/skill.md`
- **Git HEAD:** Committed version
- **Remote:** `origin/main:skills/[skill-name]/skill.md`

## Impact on Skill Execution

### Performance
- Verification adds ~2-5 seconds
- Mostly git command execution
- Negligible compared to skill execution time

### Reliability
- Ensures latest features available
- Prevents using deprecated approaches
- Maintains consistency across team

### User Experience
- Transparent process
- Clear prompts when action needed
- Quick when everything is current

## Future Enhancements

Potential improvements:

1. **Version tags:** Semantic versioning for skills
2. **Change log:** Automatic change summaries
3. **Auto-update:** Option to auto-pull latest
4. **Diff visualization:** Better diff display
5. **Notification:** Alert when remote updates available

## Summary

The self-verification system ensures:

✅ Skills are always current
✅ Users are informed of changes
✅ No silent drift from source of truth
✅ Transparent and user-controlled
✅ Minimal overhead

**Result:** Reliable, up-to-date skills that maintain their effectiveness over time.
