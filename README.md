# RHOBS Test Plan Viewer

An interactive HTML generator for RHOBS test plans, optimized for **onboarding and training new developers**.

## Features

### 🎓 Learning-Focused Design
- **Progressive Learning Path**: Guides users through tests in optimal learning sequence
- **Difficulty Levels**: Beginner, Intermediate, Advanced filtering
- **Learning Objectives**: Clear goals for what each test teaches
- **Common Beginner Mistakes**: Proactive warnings to avoid pitfalls
- **Educational Notes**: Explains *why* each step matters

### 📊 Progress Tracking
- **Local Storage Persistence**: Progress saved automatically in browser
- **Step-by-Step Completion**: Mark individual steps and tests complete
- **Progress Bar**: Visual feedback on learning journey
- **Reset Capability**: Start fresh anytime

### 🔍 Interactive Features
- **Search**: Full-text search across all tests, steps, and concepts
- **Difficulty Filter**: Focus on your skill level
- **Expand/Collapse**: Manage information density
- **Copy-to-Clipboard**: One-click command copying
- **Syntax Highlighting**: Beautiful code blocks with highlight.js

### 🎨 Beautiful UI
- **Responsive Design**: Works on desktop and mobile
- **Tailwind CSS**: Modern, clean styling
- **Print-Friendly**: Generate printable guides
- **Color-Coded**: Visual cues for difficulty, status, and content type

## Installation

```bash
# Clone or navigate to the directory
cd testcase_poc

# Initialize Go modules (if needed)
go mod tidy

# Build the binary
go build -o testplan-viewer

# Or run directly
go run main.go
```

## Usage

### Basic Usage

```bash
# Generate HTML from JSON (uses defaults)
./testplan-viewer

# Specify input and output files
./testplan-viewer -i rhobs_test_plan_v2.json -o my_testplan.html

# Using go run
go run main.go --input rhobs_test_plan_v2.json --output testplan.html
```

### Flags

- `-i, --input`: Input JSON file path (default: `rhobs_test_plan_v2.json`)
- `-o, --output`: Output HTML file path (default: `testplan.html`)
- `-h, --help`: Show help message

### Example Output

```
📖 Reading test plan from: rhobs_test_plan_v2.json
🔍 Parsing test plan JSON...
✅ Validating test plan structure...
   ✓ Metadata valid
   ✓ Test cases valid
   ✓ Dependencies valid
   ✓ Learning path valid
🎨 Generating interactive HTML to: testplan.html

✨ Success! Generated test plan HTML with:
   • 3 test cases
   • 3 concepts explained
   • Learning path: Recommended sequence for new developers to learn the system
   • Target audience: New engineers onboarding to RHOBS/ROSA
   • Estimated time: 2-3 hours

🚀 Open testplan.html in your browser to get started!
```

## JSON Structure

The tool expects JSON with this structure (optimized for learning):

```json
{
  "metadata": {
    "document_title": "Learning Guide Title",
    "target_audience": "Who this is for",
    "estimated_total_time": "How long it takes",
    "learning_objectives": ["What", "You'll", "Learn"]
  },
  "learning_path": {
    "sequence": ["test_1", "test_2", "test_3"],
    "alternative_paths": {
      "experienced_path": ["test_1", "test_3"]
    }
  },
  "testcases": {
    "test_1": {
      "metadata": {
        "difficulty": "beginner|intermediate|advanced",
        "estimated_time": "15-20 minutes"
      },
      "learning": {
        "objectives": ["What this test teaches"],
        "common_beginner_mistakes": [...]
      },
      "test_execution": {
        "steps": [
          {
            "step_number": 1,
            "title": "Step title",
            "learning_note": "Why this matters",
            "command": "Shell command",
            "common_errors": [...]
          }
        ]
      }
    }
  }
}
```

See `rhobs_test_plan_v2.json` for a complete example.

## Architecture

```
testplan-viewer/
├── main.go                    # CLI entrypoint (Cobra)
├── models/
│   └── testplan.go           # Go structs matching JSON schema
├── generator/
│   ├── html.go               # HTML generation logic
│   └── templates/
│       └── testplan.html     # Alpine.js + Tailwind template
├── rhobs_test_plan_v2.json  # Example test plan (v2 - learning optimized)
└── testplan.html             # Generated output
```

### Technology Stack

**Backend (Go):**
- `encoding/json` - JSON parsing
- `html/template` - HTML generation
- `embed` - Embed templates in binary
- `spf13/cobra` - CLI framework

**Frontend (Embedded in HTML):**
- **Alpine.js 3.x** - Lightweight reactivity (18KB)
- **Tailwind CSS** - Utility-first styling (CDN)
- **Highlight.js** - Syntax highlighting for code blocks
- **LocalStorage API** - Progress persistence

## Key Design Decisions

### Why Go + Alpine.js Instead of Pure SPA?

1. **Single Binary Distribution**: Go compiles to a single executable - easy to share
2. **Self-Contained Output**: Generated HTML has no external dependencies (CDN only)
3. **No Build Step for Users**: Recipients just open HTML in browser
4. **Offline Capable**: Works without internet (except initial CDN load)
5. **Simple**: No Node.js, npm, webpack, or complex tooling required

### Why This JSON Structure?

The v2 JSON structure prioritizes **learning outcomes** over just test execution:

**v1 (Original)** - Test-focused:
- Organized by test type (metrics, logs, synthetics)
- Steps are procedural
- Minimal educational context

**v2 (Learning-Optimized)** - Education-focused:
- Organized by learning path
- Difficulty levels for progressive learning
- Learning objectives for each test
- Common beginner mistakes highlighted
- "Why this matters" explanations
- Related concepts cross-referenced

This makes the JSON **consumable as training material**, not just test documentation.

## Use Cases

### 1. New Engineer Onboarding
```bash
# Generate HTML
./testplan-viewer -i rhobs_test_plan_v2.json -o onboarding.html

# Send onboarding.html to new team member
# They work through tests, track progress in browser
```

### 2. Workshop/Training Sessions
```bash
# Generate for specific difficulty level
# (Filter in browser UI after generation)
./testplan-viewer
# Instructor shares testplan.html
# Students follow along, track individual progress
```

### 3. CI/CD Integration
```bash
# Auto-generate docs on test plan updates
git clone repo
./testplan-viewer -i test_plan.json -o docs/learning-guide.html
git add docs/learning-guide.html
git commit -m "Update learning guide"
```

### 4. Offline Learning
```bash
# Generate HTML
./testplan-viewer

# Copy testplan.html to USB drive or shared network
# Works offline after initial CDN caching
```

## Customization

### Adding Custom Concepts

Edit `rhobs_test_plan_v2.json`:

```json
"concepts": {
  "my_concept": {
    "title": "Concept Name",
    "description": "Detailed explanation",
    "why_it_matters": "Why engineers need to understand this",
    "related_tests": ["test_1", "test_2"]
  }
}
```

### Modifying Learning Paths

```json
"learning_path": {
  "sequence": ["test_1", "test_2", "test_3"],
  "alternative_paths": {
    "fast_track": ["test_1", "test_3"],
    "deep_dive": ["test_1", "test_2", "test_3", "bonus_test"]
  }
}
```

### Styling Customization

Edit `generator/templates/testplan.html`:
- Tailwind classes can be modified directly
- Custom CSS in `<style>` section
- Color scheme via Tailwind config

## Troubleshooting

### Command Not Found: testplan-viewer

```bash
# Build first
go build -o testplan-viewer

# Or run directly
go run main.go
```

### Template Parse Error

```bash
# Ensure you're in the correct directory
cd testcase_poc

# Verify templates exist
ls -la generator/templates/
```

### JSON Validation Errors

```bash
# Check JSON syntax
cat rhobs_test_plan_v2.json | jq .

# Validate against schema
./testplan-viewer -i rhobs_test_plan_v2.json
# Look for specific error messages
```

## Future Enhancements

Potential features for future versions:

- [ ] Export to PDF
- [ ] Quiz mode for testing knowledge
- [ ] Multi-language support
- [ ] Video embedding for step demonstrations
- [ ] Integration with Jira for ticket creation
- [ ] Time tracking and analytics
- [ ] Collaborative mode (team progress)
- [ ] LLM-assisted troubleshooting chatbot
- [ ] Auto-generate test plan from code analysis

## Contributing

To improve the test plan or viewer:

1. **Edit JSON**: Update `rhobs_test_plan_v2.json` with new tests/concepts
2. **Modify Template**: Edit `generator/templates/testplan.html` for UI changes
3. **Extend Models**: Update `models/testplan.go` for new JSON fields
4. **Test**: Run `./testplan-viewer` and open in browser
5. **Iterate**: Repeat until perfect

## License

Internal Red Hat tooling - see team documentation for usage guidelines.

## Contact

- **Slack**: #forum-rhobs-core
- **Team**: SREP Observability
- **Epic**: SREP-3109
