package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/crdant/replicated-mcp-server/pkg/models"
)

// Test constants
const (
	testToken            = "test-token"
	testHTTPMethodGET    = "GET"
	testPathApplications = "/v1/applications"
)

// Helper function to create a test client with mock server
func setupTestClient(handler http.HandlerFunc) (*Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	client, _ := NewClient(ClientConfig{
		APIToken: testToken,
		BaseURL:  server.URL,
		Timeout:  30 * time.Second,
	})
	return client, server
}

func TestListApplications(t *testing.T) {
	tests := []struct {
		name           string
		opts           *ListOptions
		serverResponse string
		statusCode     int
		expectError    bool
		expectedCount  int
	}{
		{
			name: "successful list with default pagination",
			opts: nil,
			serverResponse: `{
				"applications": [
					{
						"id": "app-12345",
						"name": "Test Application",
						"slug": "test-app",
						"team_id": "team-67890",
						"team_name": "Test Team",
						"created_at": "2024-01-01T12:00:00Z",
						"updated_at": "2024-01-02T12:00:00Z",
						"description": "A test application for unit testing",
						"icon": "https://example.com/icon.png",
						"is_active": true
					}
				],
				"total_count": 1,
				"page": 1,
				"page_size": 20,
				"has_more": false
			}`,
			statusCode:    200,
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "successful list with custom pagination",
			opts: &ListOptions{Page: 2, PageSize: 10},
			serverResponse: `{
				"applications": [],
				"total_count": 5,
				"page": 2,
				"page_size": 10,
				"has_more": false
			}`,
			statusCode:    200,
			expectError:   false,
			expectedCount: 0,
		},
		{
			name:           "unauthorized error",
			opts:           nil,
			serverResponse: `{"message": "Unauthorized", "details": "Invalid API token"}`,
			statusCode:     401,
			expectError:    true,
			expectedCount:  0,
		},
		{
			name:           "server error",
			opts:           nil,
			serverResponse: `{"message": "Internal Server Error"}`,
			statusCode:     500,
			expectError:    true,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
				// Verify method and path
				if r.Method != testHTTPMethodGET {
					t.Errorf("Expected GET method, got %s", r.Method)
				}
				if r.URL.Path != testPathApplications {
					t.Errorf("Expected path %s, got %s", testPathApplications, r.URL.Path)
				}

				// Verify authentication
				if auth := r.Header.Get("Authorization"); auth != testToken {
					t.Errorf("Expected Authorization header '%s', got '%s'", testToken, auth)
				}

				// Check query parameters if opts provided
				if tt.opts != nil {
					if tt.opts.Page > 0 {
						if page := r.URL.Query().Get("page"); page == "" {
							t.Error("Expected page query parameter")
						}
					}
					if tt.opts.PageSize > 0 {
						if pageSize := r.URL.Query().Get("page_size"); pageSize == "" {
							t.Error("Expected page_size query parameter")
						}
					}
				}

				w.WriteHeader(tt.statusCode)
				fmt.Fprint(w, tt.serverResponse)
			})
			defer server.Close()

			ctx := context.Background()
			result, err := client.ListApplications(ctx, tt.opts)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("Expected result but got nil")
			}

			if len(result.Applications) != tt.expectedCount {
				t.Errorf("Expected %d applications, got %d", tt.expectedCount, len(result.Applications))
			}
		})
	}
}

func TestGetApplication(t *testing.T) {
	tests := []struct {
		name           string
		appID          string
		serverResponse string
		statusCode     int
		expectError    bool
		expectedID     string
	}{
		{
			name:  "successful get",
			appID: "app-12345",
			serverResponse: `{
				"id": "app-12345",
				"name": "Test Application",
				"slug": "test-app",
				"team_id": "team-67890",
				"team_name": "Test Team",
				"created_at": "2024-01-01T12:00:00Z",
				"updated_at": "2024-01-02T12:00:00Z",
				"description": "A test application for unit testing",
				"icon": "https://example.com/icon.png",
				"is_active": true
			}`,
			statusCode:  200,
			expectError: false,
			expectedID:  "app-12345",
		},
		{
			name:           "application not found",
			appID:          "app-nonexistent",
			serverResponse: `{"message": "Application not found"}`,
			statusCode:     404,
			expectError:    true,
			expectedID:     "",
		},
		{
			name:           "empty application ID",
			appID:          "",
			serverResponse: "",
			statusCode:     200,
			expectError:    true,
			expectedID:     "",
		},
		{
			name:           "unauthorized error",
			appID:          "app-12345",
			serverResponse: `{"message": "Unauthorized"}`,
			statusCode:     401,
			expectError:    true,
			expectedID:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
				// Verify method
				if r.Method != testHTTPMethodGET {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				// Verify path
				expectedPath := fmt.Sprintf("/v1/applications/%s", tt.appID)
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				// Verify authentication
				if auth := r.Header.Get("Authorization"); auth != testToken {
					t.Errorf("Expected Authorization header '%s', got '%s'", testToken, auth)
				}

				w.WriteHeader(tt.statusCode)
				fmt.Fprint(w, tt.serverResponse)
			})
			defer server.Close()

			ctx := context.Background()
			result, err := client.GetApplication(ctx, tt.appID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("Expected result but got nil")
			}

			if result.ID != tt.expectedID {
				t.Errorf("Expected ID %s, got %s", tt.expectedID, result.ID)
			}
		})
	}
}

func TestSearchApplications(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		opts           *ListOptions
		serverResponse string
		statusCode     int
		expectError    bool
		expectedCount  int
	}{
		{
			name:  "successful search with query",
			query: "test",
			opts:  nil,
			serverResponse: `{
				"applications": [
					{
						"id": "app-12345",
						"name": "Test Application",
						"slug": "test-app",
						"team_id": "team-67890",
						"team_name": "Test Team",
						"created_at": "2024-01-01T12:00:00Z",
						"updated_at": "2024-01-02T12:00:00Z",
						"description": "A test application for unit testing",
						"icon": "https://example.com/icon.png",
						"is_active": true
					}
				],
				"total_count": 1,
				"page": 1,
				"page_size": 20,
				"has_more": false
			}`,
			statusCode:    200,
			expectError:   false,
			expectedCount: 1,
		},
		{
			name:  "search with no results",
			query: "nonexistent",
			opts:  nil,
			serverResponse: `{
				"applications": [],
				"total_count": 0,
				"page": 1,
				"page_size": 20,
				"has_more": false
			}`,
			statusCode:    200,
			expectError:   false,
			expectedCount: 0,
		},
		{
			name:  "search with pagination",
			query: "app",
			opts:  &ListOptions{Page: 1, PageSize: 5},
			serverResponse: `{
				"applications": [],
				"total_count": 0,
				"page": 1,
				"page_size": 5,
				"has_more": false
			}`,
			statusCode:    200,
			expectError:   false,
			expectedCount: 0,
		},
		{
			name:           "empty query",
			query:          "",
			opts:           nil,
			serverResponse: "",
			statusCode:     200,
			expectError:    true,
			expectedCount:  0,
		},
		{
			name:           "server error",
			query:          "test",
			opts:           nil,
			serverResponse: `{"message": "Internal Server Error"}`,
			statusCode:     500,
			expectError:    true,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
				// Verify method and path
				if r.Method != testHTTPMethodGET {
					t.Errorf("Expected GET method, got %s", r.Method)
				}
				if r.URL.Path != "/v1/applications/search" {
					t.Errorf("Expected path /v1/applications/search, got %s", r.URL.Path)
				}

				// Verify authentication
				if auth := r.Header.Get("Authorization"); auth != testToken {
					t.Errorf("Expected Authorization header '%s', got '%s'", testToken, auth)
				}

				// Check query parameters
				if tt.query != "" {
					if q := r.URL.Query().Get("q"); q != tt.query {
						t.Errorf("Expected query parameter q=%s, got %s", tt.query, q)
					}
				}

				if tt.opts != nil {
					if tt.opts.Page > 0 {
						if page := r.URL.Query().Get("page"); page == "" {
							t.Error("Expected page query parameter")
						}
					}
					if tt.opts.PageSize > 0 {
						if pageSize := r.URL.Query().Get("page_size"); pageSize == "" {
							t.Error("Expected page_size query parameter")
						}
					}
				}

				w.WriteHeader(tt.statusCode)
				fmt.Fprint(w, tt.serverResponse)
			})
			defer server.Close()

			ctx := context.Background()
			result, err := client.SearchApplications(ctx, tt.query, tt.opts)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("Expected result but got nil")
			}

			if len(result.Applications) != tt.expectedCount {
				t.Errorf("Expected %d applications, got %d", tt.expectedCount, len(result.Applications))
			}
		})
	}
}

func TestApplicationResponseParsing(t *testing.T) {
	t.Run("parse valid application JSON", func(t *testing.T) {
		jsonData := `{
			"id": "app-12345",
			"name": "Test Application",
			"slug": "test-app",
			"team_id": "team-67890",
			"team_name": "Test Team",
			"created_at": "2024-01-01T12:00:00Z",
			"updated_at": "2024-01-02T12:00:00Z",
			"description": "A test application for unit testing",
			"icon": "https://example.com/icon.png",
			"is_active": true
		}`

		var app models.Application
		err := json.Unmarshal([]byte(jsonData), &app)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		if app.ID != "app-12345" {
			t.Errorf("Expected ID 'app-12345', got %s", app.ID)
		}
		if app.Name != "Test Application" {
			t.Errorf("Expected Name 'Test Application', got %s", app.Name)
		}
		if !app.IsActive {
			t.Error("Expected IsActive to be true")
		}
	})

	t.Run("parse application list JSON", func(t *testing.T) {
		jsonData := `{
			"applications": [
				{
					"id": "app-12345",
					"name": "Test Application",
					"slug": "test-app",
					"team_id": "team-67890",
					"created_at": "2024-01-01T12:00:00Z",
					"updated_at": "2024-01-02T12:00:00Z",
					"is_active": true
				}
			],
			"total_count": 1,
			"page": 1,
			"page_size": 20,
			"has_more": false
		}`

		var appList ApplicationList
		err := json.Unmarshal([]byte(jsonData), &appList)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		if len(appList.Applications) != 1 {
			t.Errorf("Expected 1 application, got %d", len(appList.Applications))
		}
		if appList.TotalCount != 1 {
			t.Errorf("Expected TotalCount 1, got %d", appList.TotalCount)
		}
		if appList.Page != 1 {
			t.Errorf("Expected Page 1, got %d", appList.Page)
		}
	})
}
