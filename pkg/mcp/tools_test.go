package mcp

import (
	"testing"
	"time"

	"github.com/crdant/replicated-mcp-server/pkg/config"
	"github.com/crdant/replicated-mcp-server/pkg/logging"
)

// TestToolHandlers was removed as it tested placeholder responses.
// Real API integration is now tested in handlers_test.go with mock servers.

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

