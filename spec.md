# Replicated Vendor Portal MCP Server Specification

## Overview

This document specifies the requirements and design for a Machine Chat Protocol (MCP) server that interfaces with the Replicated Vendor Portal API. The server will enable AI agents to interact with Replicated Vendor Portal accounts through a standardized interface following the MCP protocol.

## Purpose and Goals

The primary goal is to allow AI agents to interact with Replicated Vendor Portal accounts, focusing initially on core functionality related to applications, releases, channels, and customers.

### Development Phases

1. **Phase 1:** Read-only access to core entities (applications, releases, channels, customers)
2. **Phase 2:** Write capabilities for core entities
3. **Phase 3:** Expand to additional entities and capabilities

## Technical Requirements

### Programming Language and Dependencies

- **Language:** Go
- **Key Libraries:**
  - `mark3labs/mcp-go` - For MCP protocol implementation
  - `cobra` - For CLI command handling
  - `ginkgo` - For testing

### Architecture

#### Protocol

- Implement the Machine Chat Protocol (MCP) using the `mark3labs/mcp-go` library
- Communication through stdio (standard input/output)
- MCP interactions only on stdout; all logs to stderr

#### Authentication

- Use API token-based authentication for the initial version
- Token provided via environment variable or command-line flag (with flag taking precedence)

#### Core Functionality

**Phase 1 (Read-Only):**

1. **Applications**
   - List applications
   - Get application details
   - Search/filter applications (as supported by API)

2. **Releases**
   - List releases for an application
   - Get release details
   - Search/filter releases (as supported by API)

3. **Channels**
   - List channels for an application
   - Get channel details
   - Search/filter channels (as supported by API)

4. **Customers**
   - List customers
   - Get customer details
   - Search/filter customers (as supported by API)

**Future Phases:**
- Write operations for core entities
- Additional entities as per the Replicated Vendor Portal API

### Configuration

The server should support the following configuration options via both environment variables and command-line flags (with flags taking precedence):

- API token
- Log level (fatal, error, info, debug, trace; default: fatal)
- Timeouts
- API endpoint (hidden option, not visible in help)

### Error Handling

- Provide detailed, actionable error information to agents
- Include sufficient context to allow agents to fix or work around errors
- Conform to MCP protocol error handling patterns

### Rate Limiting and Throttling

- Implement adaptive throttling based on API responses
- Back off when receiving rate limit responses (HTTP 429)
- Adjust request rates dynamically based on API behavior

### Logging

- Implement standard logging levels: fatal, error, info, debug, trace
- Default log level: fatal
- All logs directed to stderr
- MCP protocol interactions only on stdout

## Development Practices

### Testing

- Use Test-Driven Development (TDD) for unit tests
  - Write tests first, then implement functionality
- Create higher-level integration tests with ginkgo

### Versioning

- Follow Semantic Versioning (semver)
- Include version, build date, and git commit hash in binaries
- Handle Vendor Portal API version changes appropriately

### Documentation

- Provide man page
- Implement help text via cobra and `--help`
- Create comprehensive README
- Include code documentation (Go docs)

## CI/CD and Release Process

### CI Pipeline

- Run format checks, linting, and unit tests on every commit and PR
- Automate integration tests

### Release Pipeline

- Use goreleaser to build for macOS, Linux, and Windows
- Generate container images with signatures
- Include Software Bill of Materials (SBOM)
- Ensure SLSA (Supply chain Levels for Software Artifacts) compliance

## Future Considerations

- Caching mechanisms (not included in initial release)
- Additional security features beyond API token
- Expanded entity support
- Write operations for all entities
