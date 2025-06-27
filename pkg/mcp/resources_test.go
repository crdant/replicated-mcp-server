package mcp

import (
	"context"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crdant/replicated-mcp-server/pkg/config"
	"github.com/crdant/replicated-mcp-server/pkg/logging"
)

func TestResourceHandlers(t *testing.T) {
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

	tests := []struct {
		resourceURI string
		testURI     string
	}{
		{
			resourceURI: "replicated://applications/{application}",
			testURI:     "replicated://applications/test-app-123",
		},
		{
			resourceURI: "replicated://applications/{application}/releases/{release}",
			testURI:     "replicated://applications/test-app-123/releases/test-release-456",
		},
		{
			resourceURI: "replicated://applications/{application}/channels/{channel}",
			testURI:     "replicated://applications/test-app-123/channels/test-channel-789",
		},
		{
			resourceURI: "replicated://applications/{application}/customers/{customer}",
			testURI:     "replicated://applications/test-app-123/customers/test-customer-101",
		},
	}

	for _, tt := range tests {
		t.Run(tt.resourceURI, func(t *testing.T) {
			// Find the resource
			var resource *resourceDefinition
			for _, resourceDef := range resources {
				if resourceDef.definition.URI == tt.resourceURI {
					resource = &resourceDef
					break
				}
			}

			if resource == nil {
				t.Fatalf("Resource '%s' not found", tt.resourceURI)
			}

			// Create a mock request
			request := createMockReadResourceRequest(tt.testURI)

			// Call the handler
			ctx := context.Background()
			contents, err := resource.handler(ctx, request)

			// Verify the response
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if contents == nil {
				t.Error("Expected contents slice, got nil")
				return
			}

			// For now, we expect empty content since actual implementation is in Step 7
			// The test mainly verifies the handler executes without error and returns proper type
		})
	}
}

func TestResourceDefinitions(t *testing.T) {
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

	tests := []struct {
		uri         string
		name        string
		description string
		mimeType    string
	}{
		{
			uri:         "replicated://applications/{application}",
			name:        "Application Data",
			description: "Access to detailed application information",
			mimeType:    "application/json",
		},
		{
			uri:         "replicated://applications/{application}/releases/{release}",
			name:        "Release Data",
			description: "Access to detailed release information",
			mimeType:    "application/json",
		},
		{
			uri:         "replicated://applications/{application}/channels/{channel}",
			name:        "Channel Data",
			description: "Access to detailed channel information",
			mimeType:    "application/json",
		},
		{
			uri:         "replicated://applications/{application}/customers/{customer}",
			name:        "Customer Data",
			description: "Access to detailed customer information",
			mimeType:    "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.uri, func(t *testing.T) {
			// Find the resource
			var resource *resourceDefinition
			for _, resourceDef := range resources {
				if resourceDef.definition.URI == tt.uri {
					resource = &resourceDef
					break
				}
			}

			if resource == nil {
				t.Fatalf("Resource '%s' not found", tt.uri)
			}

			// Verify resource properties
			if resource.definition.URI != tt.uri {
				t.Errorf("Expected URI '%s', got '%s'", tt.uri, resource.definition.URI)
			}

			if resource.definition.Name != tt.name {
				t.Errorf("Expected name '%s', got '%s'", tt.name, resource.definition.Name)
			}

			if !contains(resource.definition.Description, tt.description) {
				t.Errorf("Expected description to contain '%s', got '%s'", tt.description, resource.definition.Description)
			}

			if resource.definition.MIMEType != tt.mimeType {
				t.Errorf("Expected MIME type '%s', got '%s'", tt.mimeType, resource.definition.MIMEType)
			}

			if resource.handler == nil {
				t.Error("Expected resource to have a handler function")
			}
		})
	}
}

func TestResourceURIPatterns(t *testing.T) {
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

	// Test that all resources follow the expected URI pattern
	expectedPatterns := []struct {
		pattern     string
		description string
	}{
		{
			pattern:     "replicated://applications/{application}",
			description: "Application resources should follow replicated://applications/{application} pattern",
		},
		{
			pattern:     "replicated://applications/{application}/releases/{release}",
			description: "Release resources should follow replicated://applications/{application}/releases/{release} pattern",
		},
		{
			pattern:     "replicated://applications/{application}/channels/{channel}",
			description: "Channel resources should follow replicated://applications/{application}/channels/{channel} pattern",
		},
		{
			pattern:     "replicated://applications/{application}/customers/{customer}",
			description: "Customer resources should follow replicated://applications/{application}/customers/{customer} pattern",
		},
	}

	foundPatterns := make(map[string]bool)
	for _, resource := range resources {
		foundPatterns[resource.definition.URI] = true
	}

	for _, expected := range expectedPatterns {
		t.Run(expected.pattern, func(t *testing.T) {
			if !foundPatterns[expected.pattern] {
				t.Errorf("%s - pattern '%s' not found", expected.description, expected.pattern)
			}
		})
	}

	// Verify we don't have unexpected patterns
	if len(foundPatterns) != len(expectedPatterns) {
		t.Errorf("Expected exactly %d resource patterns, found %d", len(expectedPatterns), len(foundPatterns))
	}
}

func TestResourceHandlerErrorHandling(t *testing.T) {
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

	// Test with empty URI
	emptyRequest := createMockReadResourceRequest("")

	for _, resource := range resources {
		t.Run(resource.definition.URI+"_empty_uri", func(t *testing.T) {
			ctx := context.Background()
			contents, err := resource.handler(ctx, emptyRequest)

			// The handler should still work with empty URI (it's just logged)
			// The actual URI validation would happen in the MCP library
			if err != nil {
				t.Errorf("Unexpected error with empty URI: %v", err)
			}

			if contents == nil {
				t.Error("Expected contents slice even with empty URI")
			}
		})
	}
}

// Helper function to create a mock ReadResourceRequest
func createMockReadResourceRequest(uri string) mcp.ReadResourceRequest {
	return mcp.ReadResourceRequest{
		Request: mcp.Request{
			Method: "resources/read",
		},
		Params: mcp.ReadResourceParams{
			URI: uri,
		},
	}
}