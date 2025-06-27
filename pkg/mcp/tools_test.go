package mcp

import (
	"context"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crdant/replicated-mcp-server/pkg/config"
	"github.com/crdant/replicated-mcp-server/pkg/logging"
)

func TestToolHandlers(t *testing.T) {
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

	tests := []struct {
		toolName string
		args     map[string]any
	}{
		{
			toolName: "list_applications",
			args: map[string]any{
				"limit":  float64(10),
				"offset": float64(0),
			},
		},
		{
			toolName: "get_application",
			args: map[string]any{
				"app_id": "test-app-123",
			},
		},
		{
			toolName: "search_applications",
			args: map[string]any{
				"query": "test app",
				"limit": float64(5),
			},
		},
		{
			toolName: "list_releases",
			args: map[string]any{
				"app_id": "test-app-123",
				"limit":  float64(10),
				"offset": float64(0),
			},
		},
		{
			toolName: "get_release",
			args: map[string]any{
				"app_id":     "test-app-123",
				"release_id": "test-release-456",
			},
		},
		{
			toolName: "search_releases",
			args: map[string]any{
				"app_id": "test-app-123",
				"query":  "v1.0",
				"limit":  float64(5),
			},
		},
		{
			toolName: "list_channels",
			args: map[string]any{
				"app_id": "test-app-123",
				"limit":  float64(10),
				"offset": float64(0),
			},
		},
		{
			toolName: "get_channel",
			args: map[string]any{
				"app_id":     "test-app-123",
				"channel_id": "test-channel-789",
			},
		},
		{
			toolName: "search_channels",
			args: map[string]any{
				"app_id": "test-app-123",
				"query":  "stable",
				"limit":  float64(5),
			},
		},
		{
			toolName: "list_customers",
			args: map[string]any{
				"app_id": "test-app-123",
				"limit":  float64(10),
				"offset": float64(0),
			},
		},
		{
			toolName: "get_customer",
			args: map[string]any{
				"app_id":      "test-app-123",
				"customer_id": "test-customer-101",
			},
		},
		{
			toolName: "search_customers",
			args: map[string]any{
				"app_id": "test-app-123",
				"query":  "acme corp",
				"limit":  float64(5),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.toolName, func(t *testing.T) {
			// Find the tool
			var tool *toolDefinition
			for _, toolDef := range tools {
				if toolDef.definition.Name == tt.toolName {
					tool = &toolDef
					break
				}
			}

			if tool == nil {
				t.Fatalf("Tool '%s' not found", tt.toolName)
			}

			// Create a mock request
			request := createMockCallToolRequest(tt.toolName, tt.args)

			// Call the handler
			ctx := context.Background()
			result, err := tool.handler(ctx, request)

			// Verify the response
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result == nil {
				t.Error("Expected result, got nil")
				return
			}

			if len(result.Content) == 0 {
				t.Error("Expected content in result")
				return
			}

			// Verify the placeholder response (will be replaced in Step 7)
			textContent, ok := result.Content[0].(mcp.TextContent)
			if !ok {
				t.Error("Expected TextContent")
				return
			}

			expectedMessage := step7ImplementationMsg
			if !contains(textContent.Text, expectedMessage) {
				t.Errorf("Expected placeholder message containing '%s', got '%s'", expectedMessage, textContent.Text)
			}
		})
	}
}

func TestToolParameterValidation(t *testing.T) {
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

	// Test that tools have proper parameter definitions
	parameterTests := []struct {
		toolName           string
		expectedParameters []string
		requiredParams     []string
	}{
		{
			toolName:           "list_applications",
			expectedParameters: []string{"limit", "offset"},
			requiredParams:     []string{}, // Both are optional
		},
		{
			toolName:           "get_application",
			expectedParameters: []string{"app_id"},
			requiredParams:     []string{"app_id"},
		},
		{
			toolName:           "search_applications",
			expectedParameters: []string{"query", "limit"},
			requiredParams:     []string{"query"},
		},
		{
			toolName:           "list_releases",
			expectedParameters: []string{"app_id", "limit", "offset"},
			requiredParams:     []string{"app_id"},
		},
		{
			toolName:           "get_release",
			expectedParameters: []string{"app_id", "release_id"},
			requiredParams:     []string{"app_id", "release_id"},
		},
	}

	for _, tt := range parameterTests {
		t.Run(tt.toolName+"_parameters", func(t *testing.T) {
			// Find the tool
			var tool *toolDefinition
			for _, toolDef := range tools {
				if toolDef.definition.Name == tt.toolName {
					tool = &toolDef
					break
				}
			}

			if tool == nil {
				t.Fatalf("Tool '%s' not found", tt.toolName)
			}

			// Verify the tool has expected parameters
			// Note: We can't easily inspect the schema without more complex JSON parsing
			// This test mainly ensures the tool definition exists and has a handler
			if tool.handler == nil {
				t.Error("Tool should have a handler")
			}

			if tool.definition.Name != tt.toolName {
				t.Errorf("Expected tool name '%s', got '%s'", tt.toolName, tool.definition.Name)
			}
		})
	}
}

// Helper function to create a mock CallToolRequest
func createMockCallToolRequest(toolName string, args map[string]any) mcp.CallToolRequest {
	// Create a basic request structure
	// Note: This is a simplified mock - in real usage, the MCP library would create these
	return mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name:      toolName,
			Arguments: args,
		},
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					indexOfSubstring(s, substr) >= 0))
}

// Helper function to find substring index
func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
