# HTML Viewer - Test Notes and Status Tracking Update

**Date**: 2026-02-19
**Purpose**: Add test notes, status tracking, and export/import functionality to HTML viewer

---

## Features to Add

### 1. Test Status Tracking
- Status dropdown for each test: Not Started, In Progress, Passed, Failed, Blocked, Skipped
- Visual indicators (colors, icons) for each status
- Status shown in test card header

### 2. Test Notes
- Textarea for each test to write notes
- Auto-save on blur
- Markdown support (optional)
- Timestamp of last update

### 3. Export/Import/Save
- **Export**: Download test data as JSON file (notes, statuses, completion, timestamps)
- **Import**: Load test data from JSON file
- **Save**: Auto-save to localStorage (already exists for completion)
- **Clear Data**: Reset all notes and statuses

---

## Data Model Changes

### JavaScript State (add to testPlanApp function)

```javascript
testStatuses: {},      // {testId: 'passed'|'failed'|'in-progress'|'blocked'|'skipped'|'not-started'}
testNotes: {},         // {testId: 'user notes text'}
testTimestamps: {},    // {testId: 'ISO timestamp of last update'}
```

### localStorage Structure

```json
{
  "completedTests": ["test_id_1", "test_id_2"],
  "completedSteps": {"test_id_1-1": true},
  "testStatuses": {"test_id_1": "passed", "test_id_2": "in-progress"},
  "testNotes": {
    "test_id_1": "Test passed but observed 2s latency spike...",
    "test_id_2": "Blocked waiting for JIRA-123 fix"
  },
  "testTimestamps": {
    "test_id_1": "2026-02-19T14:30:00Z",
    "test_id_2": "2026-02-19T15:45:00Z"
  },
  "lastUpdated": "2026-02-19T15:45:00Z",
  "testPlanTitle": "RHOBS Next Test Plan"
}
```

---

## UI Changes

### 1. Export/Import Buttons (add after Reset Progress button)

```html
<button @click="exportTestData()"
        class="px-4 py-2 bg-green-100 text-green-700 rounded-lg hover:bg-green-200 transition text-sm">
    📥 Export Test Data
</button>

<label class="px-4 py-2 bg-purple-100 text-purple-700 rounded-lg hover:bg-purple-200 transition text-sm cursor-pointer">
    📤 Import Test Data
    <input type="file"
           @change="importTestData($event)"
           accept=".json"
           class="hidden">
</label>
```

### 2. Test Card Header - Add Status Badge

```html
<!-- Add after difficulty badge, before completion checkbox -->
<div class="flex items-center gap-2">
    <span class="text-xs font-medium text-gray-600">Status:</span>
    <span x-text="getStatusLabel('{{$test.Metadata.ID}}')"
          :class="getStatusBadgeClass('{{$test.Metadata.ID}}')"
          class="px-3 py-1 rounded-full text-xs font-semibold">
    </span>
</div>
```

### 3. Test Details - Add Notes and Status Section

```html
<!-- Add after Safety section, before Learning Objectives -->
<div class="p-6 bg-gray-50 border-b border-gray-200">
    <h4 class="font-bold text-gray-900 mb-4 flex items-center">
        📝 Test Notes & Status
    </h4>

    <!-- Status Selector -->
    <div class="mb-4">
        <label class="block text-sm font-medium text-gray-700 mb-2">
            Test Status
        </label>
        <select :value="getTestStatus('{{$test.Metadata.ID}}')"
                @change="setTestStatus('{{$test.Metadata.ID}}', $event.target.value)"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500">
            <option value="not-started">⚪ Not Started</option>
            <option value="in-progress">🔵 In Progress</option>
            <option value="passed">✅ Passed</option>
            <option value="failed">❌ Failed</option>
            <option value="blocked">🚫 Blocked</option>
            <option value="skipped">⏭️ Skipped</option>
        </select>
    </div>

    <!-- Notes Textarea -->
    <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">
            Test Notes
            <span class="text-xs text-gray-500 font-normal ml-2"
                  x-show="getTestTimestamp('{{$test.Metadata.ID}}')"
                  x-text="'(Last updated: ' + formatTimestamp(getTestTimestamp('{{$test.Metadata.ID}}')) + ')'">
            </span>
        </label>
        <textarea :value="getTestNote('{{$test.Metadata.ID}}')"
                  @blur="setTestNote('{{$test.Metadata.ID}}', $event.target.value)"
                  rows="6"
                  placeholder="Record observations, issues encountered, environment details, or any relevant notes..."
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 font-mono text-sm">
        </textarea>
        <p class="text-xs text-gray-500 mt-1">
            💡 Notes auto-save when you click outside the text area
        </p>
    </div>
</div>
```

---

## JavaScript Methods

### Test Status Methods

```javascript
getTestStatus(testId) {
    return this.testStatuses[testId] || 'not-started';
},

setTestStatus(testId, status) {
    this.testStatuses[testId] = status;
    this.testTimestamps[testId] = new Date().toISOString();
    this.saveProgress();
    this.calculateProgress(); // Update progress to reflect status
},

getStatusLabel(testId) {
    const status = this.getTestStatus(testId);
    const labels = {
        'not-started': '⚪ Not Started',
        'in-progress': '🔵 In Progress',
        'passed': '✅ Passed',
        'failed': '❌ Failed',
        'blocked': '🚫 Blocked',
        'skipped': '⏭️ Skipped'
    };
    return labels[status] || labels['not-started'];
},

getStatusBadgeClass(testId) {
    const status = this.getTestStatus(testId);
    const classes = {
        'not-started': 'bg-gray-100 text-gray-700',
        'in-progress': 'bg-blue-100 text-blue-700',
        'passed': 'bg-green-100 text-green-700',
        'failed': 'bg-red-100 text-red-700',
        'blocked': 'bg-orange-100 text-orange-700',
        'skipped': 'bg-purple-100 text-purple-700'
    };
    return classes[status] || classes['not-started'];
},
```

### Test Notes Methods

```javascript
getTestNote(testId) {
    return this.testNotes[testId] || '';
},

setTestNote(testId, note) {
    this.testNotes[testId] = note;
    this.testTimestamps[testId] = new Date().toISOString();
    this.saveProgress();
},

getTestTimestamp(testId) {
    return this.testTimestamps[testId] || null;
},

formatTimestamp(timestamp) {
    if (!timestamp) return '';
    const date = new Date(timestamp);
    return date.toLocaleString();
},
```

### Export/Import Methods

```javascript
exportTestData() {
    const testData = {
        completedTests: this.completedTests,
        completedSteps: this.completedSteps,
        testStatuses: this.testStatuses,
        testNotes: this.testNotes,
        testTimestamps: this.testTimestamps,
        lastUpdated: new Date().toISOString(),
        testPlanTitle: '{{.Metadata.DocumentTitle}}',
        exportVersion: '1.0.0'
    };

    // Create downloadable JSON file
    const dataStr = JSON.stringify(testData, null, 2);
    const dataBlob = new Blob([dataStr], { type: 'application/json' });
    const url = URL.createObjectURL(dataBlob);

    const link = document.createElement('a');
    link.href = url;
    link.download = `testplan-data-${new Date().toISOString().split('T')[0]}.json`;
    link.click();

    URL.revokeObjectURL(url);

    alert('✅ Test data exported successfully!');
},

importTestData(event) {
    const file = event.target.files[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
        try {
            const imported = JSON.parse(e.target.result);

            // Validate import data
            if (!imported.exportVersion) {
                if (!confirm('This file may be from an older version. Import anyway?')) {
                    return;
                }
            }

            // Confirm import
            const message = `Import test data from "${imported.testPlanTitle}"?\n\n` +
                          `This will replace your current progress:\n` +
                          `- ${imported.completedTests?.length || 0} completed tests\n` +
                          `- ${Object.keys(imported.testNotes || {}).length} test notes\n` +
                          `- Last updated: ${new Date(imported.lastUpdated).toLocaleString()}\n\n` +
                          `Current progress will be overwritten!`;

            if (!confirm(message)) {
                // Reset file input
                event.target.value = '';
                return;
            }

            // Import data
            this.completedTests = imported.completedTests || [];
            this.completedSteps = imported.completedSteps || {};
            this.testStatuses = imported.testStatuses || {};
            this.testNotes = imported.testNotes || {};
            this.testTimestamps = imported.testTimestamps || {};

            this.saveProgress();
            this.calculateProgress();

            alert('✅ Test data imported successfully!');

            // Reset file input
            event.target.value = '';
        } catch (err) {
            alert('❌ Error importing file: ' + err.message);
            event.target.value = '';
        }
    };
    reader.readAsText(file);
},

downloadTestData() {
    // Alias for exportTestData
    this.exportTestData();
},
```

### Updated Save/Load Progress

```javascript
saveProgress() {
    const progress = {
        completedTests: this.completedTests,
        completedSteps: this.completedSteps,
        testStatuses: this.testStatuses,        // NEW
        testNotes: this.testNotes,              // NEW
        testTimestamps: this.testTimestamps,    // NEW
        lastUpdated: new Date().toISOString()
    };
    localStorage.setItem('testplan_progress', JSON.stringify(progress));
},

loadProgress() {
    const saved = localStorage.getItem('testplan_progress');
    if (saved) {
        try {
            const progress = JSON.parse(saved);
            this.completedTests = progress.completedTests || [];
            this.completedSteps = progress.completedSteps || {};
            this.testStatuses = progress.testStatuses || {};      // NEW
            this.testNotes = progress.testNotes || {};            // NEW
            this.testTimestamps = progress.testTimestamps || {};  // NEW
        } catch (e) {
            console.error('Error loading progress:', e);
        }
    }
},

resetProgress() {
    if (confirm('Are you sure you want to reset ALL data (progress, notes, statuses)? This cannot be undone.')) {
        this.completedTests = [];
        this.completedSteps = {};
        this.testStatuses = {};      // NEW
        this.testNotes = {};          // NEW
        this.testTimestamps = {};     // NEW
        this.saveProgress();
        this.calculateProgress();
    }
},
```

### Enhanced Progress Calculation

```javascript
calculateProgress() {
    const totalTests = {{len .TestCases}};

    // Count tests by status
    const passedCount = Object.values(this.testStatuses).filter(s => s === 'passed').length;
    const failedCount = Object.values(this.testStatuses).filter(s => s === 'failed').length;
    const inProgressCount = Object.values(this.testStatuses).filter(s => s === 'in-progress').length;
    const blockedCount = Object.values(this.testStatuses).filter(s => s === 'blocked').length;
    const skippedCount = Object.values(this.testStatuses).filter(s => s === 'skipped').length;

    // Progress based on passed tests
    this.progressPercentage = totalTests > 0 ? (passedCount / totalTests) * 100 : 0;

    // Optional: Display detailed stats
    this.testStats = {
        total: totalTests,
        passed: passedCount,
        failed: failedCount,
        inProgress: inProgressCount,
        blocked: blockedCount,
        skipped: skippedCount,
        notStarted: totalTests - (passedCount + failedCount + inProgressCount + blockedCount + skippedCount)
    };
},
```

---

## Enhanced Progress Bar (Optional)

```html
<!-- Replace simple progress bar with detailed stats -->
<div class="bg-white border-b sticky top-0 z-40 shadow-sm no-print">
    <div class="container mx-auto px-6 py-4">
        <div class="flex items-center justify-between mb-2">
            <span class="text-sm font-medium text-gray-700">Test Execution Progress</span>
            <span class="text-sm text-gray-600">
                <span x-show="testStats" x-text="testStats.passed + ' / ' + testStats.total + ' passed'"></span>
            </span>
        </div>

        <!-- Multi-colored progress bar -->
        <div class="w-full bg-gray-200 rounded-full h-3 flex overflow-hidden">
            <div x-show="testStats && testStats.passed > 0"
                 :style="'width: ' + (testStats.passed / testStats.total * 100) + '%'"
                 class="bg-green-600 h-3"
                 title="Passed">
            </div>
            <div x-show="testStats && testStats.failed > 0"
                 :style="'width: ' + (testStats.failed / testStats.total * 100) + '%'"
                 class="bg-red-600 h-3"
                 title="Failed">
            </div>
            <div x-show="testStats && testStats.inProgress > 0"
                 :style="'width: ' + (testStats.inProgress / testStats.total * 100) + '%'"
                 class="bg-blue-600 h-3"
                 title="In Progress">
            </div>
            <div x-show="testStats && testStats.blocked > 0"
                 :style="'width: ' + (testStats.blocked / testStats.total * 100) + '%'"
                 class="bg-orange-600 h-3"
                 title="Blocked">
            </div>
            <div x-show="testStats && testStats.skipped > 0"
                 :style="'width: ' + (testStats.skipped / testStats.total * 100) + '%'"
                 class="bg-purple-600 h-3"
                 title="Skipped">
            </div>
        </div>

        <!-- Legend -->
        <div class="flex gap-4 mt-2 text-xs">
            <span x-show="testStats && testStats.passed > 0" class="flex items-center gap-1">
                <span class="w-3 h-3 bg-green-600 rounded"></span>
                <span x-text="testStats.passed + ' Passed'"></span>
            </span>
            <span x-show="testStats && testStats.failed > 0" class="flex items-center gap-1">
                <span class="w-3 h-3 bg-red-600 rounded"></span>
                <span x-text="testStats.failed + ' Failed'"></span>
            </span>
            <span x-show="testStats && testStats.inProgress > 0" class="flex items-center gap-1">
                <span class="w-3 h-3 bg-blue-600 rounded"></span>
                <span x-text="testStats.inProgress + ' In Progress'"></span>
            </span>
            <span x-show="testStats && testStats.blocked > 0" class="flex items-center gap-1">
                <span class="w-3 h-3 bg-orange-600 rounded"></span>
                <span x-text="testStats.blocked + ' Blocked'"></span>
            </span>
            <span x-show="testStats && testStats.skipped > 0" class="flex items-center gap-1">
                <span class="w-3 h-3 bg-purple-600 rounded"></span>
                <span x-text="testStats.skipped + ' Skipped'"></span>
            </span>
            <span x-show="testStats && testStats.notStarted > 0" class="flex items-center gap-1 text-gray-500">
                <span class="w-3 h-3 bg-gray-300 rounded"></span>
                <span x-text="testStats.notStarted + ' Not Started'"></span>
            </span>
        </div>
    </div>
</div>
```

---

## Implementation Summary

### Files to Modify
1. `internal/generator/templates/testplan.html`

### Sections to Update
1. **Line ~940**: Add new data properties (testStatuses, testNotes, testTimestamps, testStats)
2. **Line ~216**: Add Export/Import buttons
3. **Line ~305-450**: Add status badge and notes section to test cards
4. **Line ~930-1150**: Add new JavaScript methods
5. **Line ~94-105**: Optionally enhance progress bar with detailed stats

### Testing Checklist
- [ ] Status dropdown changes test status
- [ ] Status badge displays correctly in test card header
- [ ] Notes textarea saves on blur
- [ ] Notes persist in localStorage
- [ ] Export downloads JSON file
- [ ] Import loads data from JSON file
- [ ] Import confirms before overwriting
- [ ] Reset clears all data (notes, statuses, progress)
- [ ] Timestamp updates when notes/status change
- [ ] Progress bar reflects status distribution

---

## Export File Format Example

```json
{
  "completedTests": ["test_verify_deployment", "test_create_resource"],
  "completedSteps": {
    "test_verify_deployment-1": true,
    "test_verify_deployment-2": true
  },
  "testStatuses": {
    "test_verify_deployment": "passed",
    "test_create_resource": "in-progress",
    "test_cleanup": "blocked"
  },
  "testNotes": {
    "test_verify_deployment": "All pods running. Observed 2s startup delay on replica 3.",
    "test_create_resource": "ResourceQuota issue - need to increase namespace limits.",
    "test_cleanup": "Blocked - waiting for JIRA-456 approval before deleting resources."
  },
  "testTimestamps": {
    "test_verify_deployment": "2026-02-19T14:30:15Z",
    "test_create_resource": "2026-02-19T15:22:30Z",
    "test_cleanup": "2026-02-19T15:45:00Z"
  },
  "lastUpdated": "2026-02-19T15:45:00Z",
  "testPlanTitle": "RHOBS Next and Synthetic Monitoring E2E Test Plan",
  "exportVersion": "1.0.0"
}
```

---

## Benefits

### For Test Executors
- **Track Progress**: Know exactly which tests passed/failed/blocked
- **Document Issues**: Record errors, observations, environment details
- **Share Results**: Export and share with team
- **Resume Later**: Import previous session and continue testing

### For Teams
- **Standardized Reporting**: Consistent format for test results
- **Knowledge Sharing**: Notes capture tribal knowledge
- **Audit Trail**: Timestamps show when tests were executed
- **Collaboration**: Share test data files via git, email, Slack

### For Test Plan Authors
- **Feedback Loop**: See which tests are problematic (failed/blocked)
- **Improve Tests**: Use notes to identify unclear steps
- **Validate Coverage**: See which tests are skipped frequently

---

**Implementation Status**: Documented, ready for HTML template updates
**Next Steps**: Update internal/generator/templates/testplan.html with these changes
