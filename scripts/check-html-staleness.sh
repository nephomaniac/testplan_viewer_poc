#!/usr/bin/env bash
#
# check-html-staleness.sh - Check if HTML files need regeneration
#
# Usage: ./scripts/check-html-staleness.sh [example-dir]
#   - No args: Check all examples
#   - With arg: Check specific example (e.g., "camo" or "rhobs-next")
#
# Exit codes:
#   0 - All HTML files are up-to-date
#   1 - One or more HTML files need regeneration
#   2 - Error (missing files, etc.)

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track if any HTML is stale
STALE_COUNT=0
TOTAL_COUNT=0

check_example() {
    local example_dir="$1"
    local example_name=$(basename "$example_dir")
    local json_file="${example_dir}/${example_name}-testplan.json"
    local html_file="${example_dir}/${example_name}-testplan.html"

    echo "Checking ${example_name}..."

    # Check if JSON exists
    if [ ! -f "$json_file" ]; then
        echo -e "${RED}  ✗ JSON not found: $json_file${NC}"
        return 2
    fi

    # Check if HTML exists
    if [ ! -f "$html_file" ]; then
        echo -e "${YELLOW}  ⚠ HTML not found: $html_file${NC}"
        echo "  → Run: make html-${example_name}"
        STALE_COUNT=$((STALE_COUNT + 1))
        return 1
    fi

    # Compare timestamps
    json_mtime=$(stat -f "%m" "$json_file" 2>/dev/null || stat -c "%Y" "$json_file" 2>/dev/null)
    html_mtime=$(stat -f "%m" "$html_file" 2>/dev/null || stat -c "%Y" "$html_file" 2>/dev/null)

    if [ "$json_mtime" -gt "$html_mtime" ]; then
        echo -e "${YELLOW}  ⚠ HTML is stale (JSON modified after HTML)${NC}"
        echo "  JSON: $(date -r "$json_mtime" '+%Y-%m-%d %H:%M:%S' 2>/dev/null || date -d "@$json_mtime" '+%Y-%m-%d %H:%M:%S')"
        echo "  HTML: $(date -r "$html_mtime" '+%Y-%m-%d %H:%M:%S' 2>/dev/null || date -d "@$html_mtime" '+%Y-%m-%d %H:%M:%S')"
        echo "  → Run: make html-${example_name}"
        STALE_COUNT=$((STALE_COUNT + 1))
        return 1
    else
        echo -e "${GREEN}  ✓ HTML is up-to-date${NC}"
        return 0
    fi
}

# Main execution
if [ $# -eq 0 ]; then
    # Check all examples
    echo "Checking all test plan examples..."
    echo ""

    for example_dir in examples/*/; do
        if [ -d "$example_dir" ]; then
            check_example "$example_dir" || true
            TOTAL_COUNT=$((TOTAL_COUNT + 1))
            echo ""
        fi
    done
else
    # Check specific example
    example_name="$1"
    example_dir="examples/${example_name}"

    if [ ! -d "$example_dir" ]; then
        echo -e "${RED}Error: Example directory not found: $example_dir${NC}"
        exit 2
    fi

    check_example "$example_dir" || true
    TOTAL_COUNT=1
fi

# Summary
echo "─────────────────────────────────────"
if [ $STALE_COUNT -eq 0 ]; then
    echo -e "${GREEN}✓ All HTML files are up-to-date${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠ ${STALE_COUNT}/${TOTAL_COUNT} HTML file(s) need regeneration${NC}"
    echo ""
    echo "To regenerate all:"
    echo "  make html-all"
    exit 1
fi
