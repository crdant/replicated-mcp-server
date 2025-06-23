# Replicated Vendor Portal MCP Server

A Machine Chat Protocol (MCP) server that interfaces with the Replicated Vendor Portal API, enabling AI agents to interact with Replicated Vendor Portal accounts.

## Overview

This project implements an MCP server that allows AI agents to access and manipulate Replicated Vendor Portal resources through a standardized interface following the MCP protocol.

### Features

- Read-only access to core Replicated entities (applications, releases, channels, customers)
- Simple configuration via environment variables or command-line flags
- MCP-compliant interface for AI agent interaction

## Installation

```bash
go install github.com/crdant/replicated-mcp-server/cmd/server@latest
```

Or download pre-built binaries from the [releases page](https://github.com/crdant/replicated-mcp-server/releases).

## Usage

```bash
REPLICATED_API_TOKEN="your-api-token" replicated-mcp-server
```

Or with command-line flags:

```bash
replicated-mcp-server --api-token="your-api-token" --log-level=info
```

## Configuration

| Flag | Environment Variable | Description | Default |
|------|---------------------|-------------|---------|
| `--api-token` | `REPLICATED_API_TOKEN` | Replicated Vendor Portal API token | *(required)* |
| `--log-level` | `LOG_LEVEL` | Log level (fatal, error, info, debug, trace) | `fatal` |
| `--timeout` | `TIMEOUT` | API request timeout in seconds | `30` |

## Development

This project uses standard Go development practices.

```bash
# Clone the repository
git clone https://github.com/crdant/replicated-mcp-server.git
cd replicated-mcp-server

# Install dependencies
go mod download

# Run tests
go test ./...

# Build the binary
go build -o replicated-mcp-server ./cmd/server
```

## License

[MIT](LICENSE)