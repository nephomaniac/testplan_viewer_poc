.PHONY: build test clean run install dist help

# Build variables
BINARY_NAME=testplan-viewer
BUILD_DIR=build
DIST_DIR=dist
CMD_DIR=cmd/testplan-viewer

# Default target
all: build

## build: Build the application binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "✅ Built: $(BUILD_DIR)/$(BINARY_NAME)"

## test: Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

## clean: Remove build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR) $(DIST_DIR)
	rm -f testplan.html
	@echo "✅ Clean complete"

## run: Build and run with default example
run: build
	@echo "🚀 Running with example..."
	./$(BUILD_DIR)/$(BINARY_NAME) -i examples/camo/camo-testplan.json -o testplan.html

## install: Install binary to $GOPATH/bin
install:
	@echo "📦 Installing $(BINARY_NAME)..."
	go install ./$(CMD_DIR)
	@echo "✅ Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

## dist: Build release binaries for multiple platforms
dist:
	@echo "📦 Building distribution binaries..."
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)
	@echo "✅ Distribution binaries built in $(DIST_DIR)/"

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' Makefile | sed 's/## /  /'

.DEFAULT_GOAL := help
