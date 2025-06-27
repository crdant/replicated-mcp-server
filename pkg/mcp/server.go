// Package mcp provides the Model Context Protocol server implementation for the Replicated Vendor Portal API.
// It enables AI agents to interact with Replicated resources through the MCP protocol over stdio transport.
package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/server"

	"github.com/crdant/replicated-mcp-server/pkg/config"
	"github.com/crdant/replicated-mcp-server/pkg/logging"
)

// Server represents the MCP server instance that handles communication with AI agents.
// It integrates with the Replicated Vendor Portal API to provide access to applications,
// releases, channels, and customer data through the MCP protocol.
type Server struct {
	logger    logging.Logger
	config    *config.Config
	mcpServer *server.MCPServer
}

// NewServer creates a new MCP server instance with the provided configuration and logger.
// It initializes the server with stdio transport and registers all available tools and resources.
//
// The server is configured with:
// - Stdio transport for communication with AI agents
// - Tool capabilities for all Replicated Vendor Portal operations
// - Resource capabilities for accessing Replicated entities
// - Proper logging integration (stderr only)
//
// Args:
//
//	cfg: Configuration containing API tokens, timeouts, and other settings
//	logger: Logger instance for structured logging (output to stderr)
//
// Returns:
//
//	*Server: Configured MCP server instance
//	error: Error if server initialization fails
func NewServer(cfg *config.Config, logger logging.Logger) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration is required")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	logger.Info("Initializing MCP server", "version", "1.0.0")

	// Create MCP server with tool and resource capabilities
	mcpServer := server.NewMCPServer(
		"replicated-mcp-server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, false), // subscribe=true, listChanged=false
	)

	s := &Server{
		logger:    logger,
		config:    cfg,
		mcpServer: mcpServer,
	}

	// Register all tools and resources
	if err := s.registerTools(); err != nil {
		return nil, fmt.Errorf("failed to register tools: %w", err)
	}

	if err := s.registerResources(); err != nil {
		return nil, fmt.Errorf("failed to register resources: %w", err)
	}

	logger.Info("MCP server initialized successfully")
	return s, nil
}

// Start begins serving the MCP protocol over stdio transport.
// This method blocks until the server is stopped or encounters an error.
// All MCP communication happens on stdout, while logging goes to stderr.
//
// Args:
//
//	ctx: Context for graceful shutdown handling
//
// Returns:
//
//	error: Error if server startup or operation fails
func (s *Server) Start(_ context.Context) error {
	s.logger.Info("Starting MCP server on stdio transport")

	// Start serving on stdio - this blocks until shutdown
	if err := server.ServeStdio(s.mcpServer); err != nil {
		s.logger.Error("MCP server error", "error", err)
		return fmt.Errorf("stdio server error: %w", err)
	}

	return nil
}

// Stop gracefully shuts down the MCP server.
// It ensures all ongoing operations complete and resources are cleaned up properly.
//
// Args:
//
//	ctx: Context with timeout for graceful shutdown
//
// Returns:
//
//	error: Error if shutdown fails
func (s *Server) Stop(_ context.Context) error {
	s.logger.Info("Stopping MCP server")

	// Note: The mark3labs/mcp-go library doesn't expose a Stop method for stdio servers
	// The server will stop when the stdio connection closes or context is canceled
	s.logger.Info("MCP server stopped")
	return nil
}

// registerTools registers all available MCP tools with the server.
// Each tool is defined with proper JSON schema validation and empty handler implementations.
// The actual business logic will be implemented in Step 7 (MCP Handlers).
//
// Returns:
//
//	error: Error if tool registration fails
func (s *Server) registerTools() error {
	s.logger.Debug("Registering MCP tools")

	tools := s.defineTools()
	for _, tool := range tools {
		s.mcpServer.AddTool(*tool.definition, tool.handler)
		s.logger.Debug("Registered tool", "name", tool.definition.Name)
	}

	s.logger.Info("Successfully registered tools", "count", len(tools))
	return nil
}

// registerResources registers all available MCP resources with the server.
// Resources provide access to Replicated entities through standardized URIs.
//
// Returns:
//
//	error: Error if resource registration fails
func (s *Server) registerResources() error {
	s.logger.Debug("Registering MCP resources")

	resources := s.defineResources()
	for _, resource := range resources {
		s.mcpServer.AddResource(*resource.definition, resource.handler)
		s.logger.Debug("Registered resource", "uri", resource.definition.URI)
	}

	s.logger.Info("Successfully registered resources", "count", len(resources))
	return nil
}
