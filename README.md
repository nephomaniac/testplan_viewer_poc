# RHOBS Test Plan Viewer

A CLI tool that generates interactive, educational HTML documentation from RHOBS test plan JSON files.

## Quick Start

```bash
# Build the application
make build

# Generate HTML from an example test plan
make run

# Or use the binary directly
./build/testplan-viewer -i examples/rhobs_test_plan_v2.json -o testplan.html

# Open testplan.html in your browser
```

## Features

- **Interactive Learning**: Progress tracking, collapsible sections, and copy-to-clipboard commands
- **Educational Focus**: Learning objectives, common mistakes, and troubleshooting guides
- **Filtering & Search**: Filter by difficulty, search across tests and steps
- **Browser-Based**: No server required, progress saved in localStorage
- **Single Binary**: Distributes as a single executable with embedded templates

## Project Structure

```
testplan_tools_poc/
├── cmd/testplan-viewer/        # CLI entry point
├── internal/                    # Internal packages
│   ├── models/                 # Data structures
│   ├── parser/                 # JSON parsing and validation
│   └── generator/              # HTML generation
├── examples/                    # Sample test plan JSON files
├── docs/                        # Detailed documentation
├── build/                       # Build output
└── dist/                        # Distribution binaries
```

## Documentation

- [Detailed README](docs/README.md) - Complete documentation
- [Quick Start Guide](docs/QUICKSTART.md) - Get started quickly
- [Implementation Summary](docs/IMPLEMENTATION_SUMMARY.md) - Technical details
- [Examples README](examples/README.md) - How to create test plans

## Building

```bash
# Build for current platform
make build

# Run tests
make test

# Build for multiple platforms
make dist

# Install to $GOPATH/bin
make install

# Clean build artifacts
make clean
```

## Usage

```bash
# Basic usage
testplan-viewer -i <input.json> -o <output.html>

# With default paths
testplan-viewer  # Uses examples/rhobs_test_plan_v2.json

# Examples
testplan-viewer -i examples/rhobs_test_plan.json -o my_testplan.html
```

## Creating Test Plans

See [examples/README.md](examples/README.md) for:
- JSON schema reference
- Minimal valid example
- Validation requirements
- Tips for creating great test plans

## Development

```bash
# Install dependencies
go mod download

# Build
make build

# Run with example
make run

# Run tests
make test
```

## Architecture

- **Parser**: Reads and validates test plan JSON files
- **Generator**: Transforms test plans into interactive HTML
- **Models**: Pure data structures for test plan schema
- **CLI**: Command-line interface using Cobra

Each component is independently testable and follows the single responsibility principle.

## Requirements

- Go 1.22 or later
- No runtime dependencies (generates standalone HTML)

## License

Internal Red Hat tool for RHOBS team use.

## Contributing

For questions or improvements, contact the RHOBS SRE team.
