# Makefile for Replicated MCP Server

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
BINARY_NAME=replicated-mcp-server
BINARY_PATH=./cmd/server

# Build the application
.PHONY: build
build:
	$(GOBUILD) -v -o $(BINARY_NAME) $(BINARY_PATH)

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) verify

# Run tests with coverage
.PHONY: test
test:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
.PHONY: test-coverage
test-coverage: test
	$(GOCMD) tool cover -func=coverage.out
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code using gofmt
.PHONY: format
format:
	@echo "Formatting Go code..."
	$(GOFMT) -s -w .
	@echo "✅ Code formatted"

# Check if code is formatted
.PHONY: format-check
format-check:
	@echo "Checking code formatting..."
	@if [ -n "$$($(GOFMT) -l .)" ]; then \
		echo "❌ Code is not formatted. Run 'make format' to fix."; \
		$(GOFMT) -l .; \
		exit 1; \
	else \
		echo "✅ Code is properly formatted"; \
	fi

# Run linting (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
		echo "✅ Linting completed"; \
	else \
		echo "❌ golangci-lint not found. Install it from https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

# Install golangci-lint (Linux/macOS) - matches CI version
.PHONY: install-linter
install-linter:
	@echo "Installing golangci-lint (latest version to match CI)..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest
	@echo "✅ golangci-lint installed (latest version)"

# Fix linting issues that can be auto-fixed
.PHONY: lint-fix
lint-fix:
	@echo "Auto-fixing linting issues..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --fix --timeout=5m; \
		echo "✅ Auto-fixable linting issues resolved"; \
	else \
		echo "❌ golangci-lint not found. Install it with 'make install-linter'"; \
		exit 1; \
	fi

# Tidy go modules
.PHONY: tidy
tidy:
	$(GOMOD) tidy
	@echo "✅ Go modules tidied"

# Run all checks (format, lint, test)
.PHONY: check
check: format-check lint test
	@echo "✅ All checks passed"

# Run CI-like checks locally
.PHONY: ci
ci: deps tidy format-check lint test
	@echo "✅ CI checks completed successfully"

# Run the application with help
.PHONY: run-help
run-help: build
	./$(BINARY_NAME) --help

# Run the application with version
.PHONY: run-version
run-version: build
	./$(BINARY_NAME) --version

# Development setup - install tools and dependencies
.PHONY: setup
setup: install-linter deps
	@echo "✅ Development environment setup complete"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download and verify dependencies"
	@echo "  test          - Run tests with coverage"
	@echo "  test-coverage - Run tests and generate coverage report"
	@echo "  format        - Format code using gofmt"
	@echo "  format-check  - Check if code is properly formatted"
	@echo "  lint          - Run golangci-lint"
	@echo "  lint-fix      - Auto-fix linting issues"
	@echo "  install-linter- Install golangci-lint"
	@echo "  tidy          - Tidy go modules"
	@echo "  check         - Run all checks (format, lint, test)"
	@echo "  ci            - Run CI-like checks locally"
	@echo "  run-help      - Build and run application with --help"
	@echo "  run-version   - Build and run application with --version"
	@echo "  setup         - Set up development environment"
	@echo "  help          - Show this help message"

# Default target
.DEFAULT_GOAL := help