package mcp

import (
	"context"
	"testing"
	"time"

	"github.com/crdant/replicated-mcp-server/pkg/config"
	"github.com/crdant/replicated-mcp-server/pkg/logging"
)

// Test constants
const (
	listApplicationsToolName = "list_applications"
	applicationResourceURI   = "replicated://applications/{application}"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		logger      logging.Logger
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid configuration",
			config: &config.Config{
				APIToken: "test-token",
				LogLevel: "info",
				Timeout:  30 * time.Second,
				Endpoint: "https://api.replicated.com",
			},
			logger:      logging.NewLogger("info"),
			expectError: false,
		},
		{
			name:        "nil configuration",
			config:      nil,
			logger:      logging.NewLogger("info"),
			expectError: true,
			errorMsg:    "configuration is required",
		},
		{
			name: "nil logger",
			config: &config.Config{
				APIToken: "test-token",
				LogLevel: "info",
				Timeout:  30 * time.Second,
			},
			logger:      nil,
			expectError: true,
			errorMsg:    "logger is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := NewServer(tt.config, tt.logger)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
				if server != nil {
					t.Errorf("Expected nil server when error occurs, got %v", server)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if server == nil {
					t.Errorf("Expected server instance, got nil")
				}
				if server != nil {
					// Verify server fields are properly set
					if server.config != tt.config {
						t.Errorf("Expected config to be set correctly")
					}
					if server.logger != tt.logger {
						t.Errorf("Expected logger to be set correctly")
					}
					if server.mcpServer == nil {
						t.Errorf("Expected MCP server to be initialized")
					}
				}
			}
		})
	}
}

func TestServerStop(t *testing.T) {
	cfg := &config.Config{
		APIToken: "test-token",
		LogLevel: "info",
		Timeout:  30 * time.Second,
	}
	logger := logging.NewLogger("info")

	server, err := NewServer(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test graceful stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Stop(ctx)
	if err != nil {
		t.Errorf("Unexpected error during stop: %v", err)
	}
}

func TestServerToolRegistration(t *testing.T) {
	cfg := &config.Config{
		APIToken: "test-token",
		LogLevel: "debug", // Use debug level to see registration logs
		Timeout:  30 * time.Second,
	}
	logger := logging.NewLogger("debug")

	server, err := NewServer(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test that tools are registered - this happens during NewServer
	// We expect 12 tools to be registered (3 each for applications, releases, channels, customers)
	tools := server.defineTools()
	expectedToolCount := 12

	if len(tools) != expectedToolCount {
		t.Errorf("Expected %d tools to be defined, got %d", expectedToolCount, len(tools))
	}

	// Verify all expected tools are present
	expectedToolNames := []string{
		"list_applications", "get_application", "search_applications",
		"list_releases", "get_release", "search_releases",
		"list_channels", "get_channel", "search_channels",
		"list_customers", "get_customer", "search_customers",
	}

	foundTools := make(map[string]bool)
	for _, tool := range tools {
		foundTools[tool.definition.Name] = true
	}

	for _, expectedName := range expectedToolNames {
		if !foundTools[expectedName] {
			t.Errorf("Expected tool '%s' not found", expectedName)
		}
	}
}

func TestServerResourceRegistration(t *testing.T) {
	cfg := &config.Config{
		APIToken: "test-token",
		LogLevel: "debug",
		Timeout:  30 * time.Second,
	}
	logger := logging.NewLogger("debug")

	server, err := NewServer(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test that resources are registered
	resources := server.defineResources()
	expectedResourceCount := 4

	if len(resources) != expectedResourceCount {
		t.Errorf("Expected %d resources to be defined, got %d", expectedResourceCount, len(resources))
	}

	// Verify all expected resources are present
	expectedResourceURIs := []string{
		"replicated://applications/{application}",
		"replicated://applications/{application}/releases/{release}",
		"replicated://applications/{application}/channels/{channel}",
		"replicated://applications/{application}/customers/{customer}",
	}

	foundResources := make(map[string]bool)
	for _, resource := range resources {
		foundResources[resource.definition.URI] = true
	}

	for _, expectedURI := range expectedResourceURIs {
		if !foundResources[expectedURI] {
			t.Errorf("Expected resource with URI '%s' not found", expectedURI)
		}
	}
}

func TestServerToolDefinitions(t *testing.T) {
	cfg := &config.Config{
		APIToken: "test-token",
		LogLevel: "info",
		Timeout:  30 * time.Second,
	}
	logger := logging.NewLogger("info")

	server, err := NewServer(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	tools := server.defineTools()

	// Test a specific tool definition (list_applications)
	var listAppsTool *toolDefinition
	for _, tool := range tools {
		if tool.definition.Name == listApplicationsToolName {
			listAppsTool = &tool
			break
		}
	}

	if listAppsTool == nil {
		t.Fatal("list_applications tool not found")
	}

	// Verify tool properties
	if listAppsTool.definition.Name != listApplicationsToolName {
		t.Errorf("Expected tool name '%s', got '%s'", listApplicationsToolName, listAppsTool.definition.Name)
	}

	if listAppsTool.definition.Description == "" {
		t.Error("Expected tool to have a description")
	}

	if listAppsTool.handler == nil {
		t.Error("Expected tool to have a handler function")
	}
}

func TestServerResourceDefinitions(t *testing.T) {
	cfg := &config.Config{
		APIToken: "test-token",
		LogLevel: "info",
		Timeout:  30 * time.Second,
	}
	logger := logging.NewLogger("info")

	server, err := NewServer(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	resources := server.defineResources()

	// Test a specific resource definition (application)
	var appResource *resourceDefinition
	for _, resource := range resources {
		if resource.definition.URI == applicationResourceURI {
			appResource = &resource
			break
		}
	}

	if appResource == nil {
		t.Fatal("Application resource not found")
	}

	// Verify resource properties
	if appResource.definition.URI != applicationResourceURI {
		t.Errorf("Expected resource URI '%s', got '%s'", applicationResourceURI, appResource.definition.URI)
	}

	if appResource.definition.Name != "Application Data" {
		t.Errorf("Expected resource name 'Application Data', got '%s'", appResource.definition.Name)
	}

	if appResource.definition.Description == "" {
		t.Error("Expected resource to have a description")
	}

	if appResource.definition.MIMEType != "application/json" {
		t.Errorf("Expected MIME type 'application/json', got '%s'", appResource.definition.MIMEType)
	}

	if appResource.handler == nil {
		t.Error("Expected resource to have a handler function")
	}
}
