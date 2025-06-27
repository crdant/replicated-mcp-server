# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Core Commands
- **Build**: `go build -o replicated-mcp-server ./cmd/server`
- **Test**: `go test -v -race -coverprofile=coverage.out ./...`
- **Test with coverage**: `go tool cover -func=coverage.out`
- **Lint**: Uses golangci-lint via GitHub Actions (no local Makefile)
- **Run server**: `REPLICATED_API_TOKEN="your-token" ./replicated-mcp-server`
- **Dependencies**: `go mod download && go mod verify`

### Testing Individual Packages
- Single package: `go test -v ./pkg/config`
- Specific test: `go test -v ./pkg/config -run TestLoad`
- All packages: `go test ./...`

## Architecture Overview

### MCP Server Implementation
This is a Machine Chat Protocol (MCP) server that interfaces with the Replicated Vendor Portal API. Key architectural components:

- **Transport**: Uses stdio (standard input/output) for MCP communication
- **Protocol**: All MCP interactions on stdout, logging to stderr only
- **Authentication**: API token-based (via environment variable or CLI flag)

### Package Structure
- `cmd/server/`: Main application entry point with Cobra CLI
- `pkg/mcp/`: MCP server implementation, tools, and resources
- `pkg/api/`: HTTP client for Replicated Vendor Portal API
- `pkg/models/`: Data structures for Replicated entities (Application, Release, Channel, Customer)
- `pkg/config/`: Configuration management (env vars + CLI flags)
- `pkg/logging/`: Structured logging (stderr only)

### Key Integration Points
- **MCP Library**: Uses `mark3labs/mcp-go` for protocol implementation
- **CLI Framework**: Cobra for command-line interface
- **HTTP Client**: Custom client in `pkg/api/` with authentication headers
- **Configuration**: Environment variables with CLI flag precedence

### Development Phases
- **Phase 1** (Current): Read-only access to core entities
- **Phase 2**: Write capabilities
- **Phase 3**: Extended entities and features

### Configuration
All configuration supports both environment variables and CLI flags (flags take precedence):
- `REPLICATED_API_TOKEN` / `--api-token`: Required API token
- `LOG_LEVEL` / `--log-level`: fatal, error, info, debug, trace (default: fatal)
- `TIMEOUT` / `--timeout`: API timeout in seconds (default: 30)

### Testing Requirements
- Coverage threshold: 70%
- TDD approach: write tests first
- Integration tests in CI pipeline
- No tests for main.go (focused on package testing)

#### Red-Green-Refactor TDD Process
**Always follow this strict cycle when implementing new functionality:**

1. **RED**: Write a failing test first
   - Run `go test ./pkg/[package]` to confirm the test fails
   - Test should fail for the right reason (missing function/method, not compilation error)
   - Write the minimal test that expresses the behavior you want

2. **GREEN**: Write minimal code to make the test pass
   - Implement only enough code to make the failing test pass
   - Don't worry about code quality yet - focus on making it work
   - Run `go test ./pkg/[package]` to confirm the test now passes
   - All existing tests must still pass

3. **REFACTOR**: Improve the code while keeping tests green
   - Clean up implementation without changing behavior
   - Extract functions, improve naming, reduce duplication
   - Run tests frequently during refactoring: `go test ./pkg/[package]`
   - Ensure coverage stays above 70%

**TDD Guidelines:**
- Never write production code without a failing test first
- Keep test cycles short (minutes, not hours)
- Each test should focus on one specific behavior
- Use table-driven tests for multiple scenarios
- Mock external dependencies (HTTP calls, file I/O)
- Test error cases as thoroughly as success cases

### Pull Request Guidelines

#### Title Requirements
- **Format**: Start with verb ending in 's' (Adds, Fixes, Updates, Implements)
- **Length**: Maximum 40 characters
- **Tense**: Present tense only, no past/future tense
- **Voice**: No first person (avoid "I", "We", "My")

#### Body Structure
```markdown
TL;DR
-----
[1-2 sentences explaining the change and impact]

Details
--------
[Context explaining WHY this change was needed]
[Implementation details only if essential for understanding]
```

#### Content Guidelines
- **Professional yet conversational**: Technical authority with approachable tone
- **Impact-focused**: Explain practical benefits and deployment implications
- **Active voice**: "Updates the client" not "The client was updated"
- **Technical precision**: Use exact version numbers and specific configuration names
- **Problem-solving narrative**: Tell the story of what problem this solves
- **Avoid**: "this PR", "this change", "this commit" - use direct descriptions instead

#### Examples
- ✅ **Good**: "Fixes authentication timeout in API client"
- ❌ **Bad**: "Fix auth bug" (wrong verb form, too vague)
- ✅ **Good**: "Implements rate limiting for Replicated API calls"
- ❌ **Bad**: "I added rate limiting" (first person, wrong tense)

### CI/CD
- **Lint**: golangci-lint with 5m timeout
- **Test**: Go 1.21, 1.22, 1.23 matrix
- **Build**: Multi-platform via GitHub Actions
- **Integration**: End-to-end configuration and logging tests