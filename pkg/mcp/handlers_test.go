package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crdant/replicated-mcp-server/pkg/api"
	"github.com/crdant/replicated-mcp-server/pkg/config"
	"github.com/crdant/replicated-mcp-server/pkg/logging"
)

const (
	testAppID = "app1"
)

func TestHandleListApplications(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/apps" {
			t.Errorf("Expected path /v3/apps, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"apps": [{"id": "app1", "name": "Test App", "slug": "test-app"}]}`))
	}))
	defer server.Close()

	// Create API client
	cfg := &config.Config{
		APIToken: "test-token",
		Timeout:  30 * time.Second,
		Endpoint: server.URL,
	}
	logger := logging.NewLogger("error")
	apiClient, err := api.NewClient(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	// Create MCP server with API client
	mcpServer, err := NewServerWithClient(cfg, logger, apiClient)
	if err != nil {
		t.Fatalf("Failed to create MCP server: %v", err)
	}

	// Create request
	request := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name: "list_applications",
			Arguments: map[string]interface{}{
				"limit":  10,
				"offset": 0,
			},
		},
	}

	// Call handler
	result, err := mcpServer.handleListApplications(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	// Verify result
	if len(result.Content) != 1 {
		t.Fatalf("Expected 1 content item, got %d", len(result.Content))
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}

	// Parse JSON response
	var apps []api.Application
	if err := json.Unmarshal([]byte(textContent.Text), &apps); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(apps) != 1 {
		t.Fatalf("Expected 1 app, got %d", len(apps))
	}
	if apps[0].ID != testAppID {
		t.Errorf("Expected app ID 'app1', got %s", apps[0].ID)
	}
}

func TestHandleGetApplication(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/app/test-app" {
			t.Errorf("Expected path /v3/app/test-app, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"app": {"id": "app1", "name": "Test App", "slug": "test-app"}}`))
	}))
	defer server.Close()

	// Create API client
	cfg := &config.Config{
		APIToken: "test-token",
		Timeout:  30 * time.Second,
		Endpoint: server.URL,
	}
	logger := logging.NewLogger("error")
	apiClient, err := api.NewClient(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	// Create MCP server with API client
	mcpServer, err := NewServerWithClient(cfg, logger, apiClient)
	if err != nil {
		t.Fatalf("Failed to create MCP server: %v", err)
	}

	// Create request
	request := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name: "get_application",
			Arguments: map[string]interface{}{
				"app_id": "test-app",
			},
		},
	}

	// Call handler
	result, err := mcpServer.handleGetApplication(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	// Verify result
	if len(result.Content) != 1 {
		t.Fatalf("Expected 1 content item, got %d", len(result.Content))
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}

	// Parse JSON response
	var app api.Application
	if err := json.Unmarshal([]byte(textContent.Text), &app); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if app.ID != testAppID {
		t.Errorf("Expected app ID 'app1', got %s", app.ID)
	}
}

func TestHandleSearchApplications(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/apps/search" {
			t.Errorf("Expected path /v3/apps/search, got %s", r.URL.Path)
		}
		query := r.URL.Query().Get("q")
		if query != "test" {
			t.Errorf("Expected query 'test', got %s", query)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"apps": [{"id": "app1", "name": "Test App", "slug": "test-app"}]}`))
	}))
	defer server.Close()

	// Create API client
	cfg := &config.Config{
		APIToken: "test-token",
		Timeout:  30 * time.Second,
		Endpoint: server.URL,
	}
	logger := logging.NewLogger("error")
	apiClient, err := api.NewClient(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	// Create MCP server with API client
	mcpServer, err := NewServerWithClient(cfg, logger, apiClient)
	if err != nil {
		t.Fatalf("Failed to create MCP server: %v", err)
	}

	// Create request
	request := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name: "search_applications",
			Arguments: map[string]interface{}{
				"query": "test",
				"limit": 10,
			},
		},
	}

	// Call handler
	result, err := mcpServer.handleSearchApplications(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	// Verify result
	if len(result.Content) != 1 {
		t.Fatalf("Expected 1 content item, got %d", len(result.Content))
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}

	// Parse JSON response
	var apps []api.Application
	if err := json.Unmarshal([]byte(textContent.Text), &apps); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(apps) != 1 {
		t.Fatalf("Expected 1 app, got %d", len(apps))
	}
	if apps[0].ID != testAppID {
		t.Errorf("Expected app ID 'app1', got %s", apps[0].ID)
	}
}

