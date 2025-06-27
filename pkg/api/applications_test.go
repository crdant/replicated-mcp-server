package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Test constants
const (
	testMethodGET = "GET"
	testPathApps  = "/vendor/v3/apps"
)

func TestApplicationService_ListApplications(t *testing.T) {
	tests := []struct {
		name                  string
		opts                  *ListApplicationsOptions
		mockResponse          string
		mockStatus            int
		expectError           bool
		expectedCount         int
		expectExcludeChannels bool
	}{
		{
			name: "successful list with default options",
			opts: nil,
			mockResponse: `{
				"applications": [
					{
						"id": "app-1",
						"name": "Test App 1",
						"slug": "test-app-1",
						"team_id": "team-1",
						"team_name": "Test Team",
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z",
						"description": "Test application 1",
						"is_active": true
					},
					{
						"id": "app-2",
						"name": "Test App 2",
						"slug": "test-app-2",
						"team_id": "team-1",
						"team_name": "Test Team",
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z",
						"description": "Test application 2",
						"is_active": true
					}
				]
			}`,
			mockStatus:            http.StatusOK,
			expectError:           false,
			expectedCount:         2,
			expectExcludeChannels: false,
		},
		{
			name: "successful list with exclude channels",
			opts: &ListApplicationsOptions{ExcludeChannels: true},
			mockResponse: `{
				"applications": [
					{
						"id": "app-1",
						"name": "Test App 1",
						"slug": "test-app-1",
						"team_id": "team-1",
						"team_name": "Test Team",
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z",
						"description": "Test application 1",
						"is_active": true
					}
				]
			}`,
			mockStatus:            http.StatusOK,
			expectError:           false,
			expectedCount:         1,
			expectExcludeChannels: true,
		},
		{
			name:          "empty list",
			opts:          nil,
			mockResponse:  `{"applications": []}`,
			mockStatus:    http.StatusOK,
			expectError:   false,
			expectedCount: 0,
		},
		{
			name:         "unauthorized error",
			opts:         nil,
			mockResponse: `{"message": "Unauthorized", "details": "Invalid API token"}`,
			mockStatus:   http.StatusUnauthorized,
			expectError:  true,
		},
		{
			name:         "internal server error",
			opts:         nil,
			mockResponse: `{"message": "Internal Server Error"}`,
			mockStatus:   http.StatusInternalServerError,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method and path
				if r.Method != testMethodGET {
					t.Errorf("Expected GET request, got %s", r.Method)
				}
				if r.URL.Path != testPathApps {
					t.Errorf("Expected path /vendor/v3/apps, got %s", r.URL.Path)
				}

				// Verify authentication header
				if auth := r.Header.Get("Authorization"); auth == "" {
					t.Error("Expected Authorization header")
				}

				// Check excludeChannels parameter
				query := r.URL.Query()
				excludeChannels := query.Get("excludeChannels")
				if tt.expectExcludeChannels {
					if excludeChannels != "true" {
						t.Errorf("Expected excludeChannels=true, got excludeChannels=%s", excludeChannels)
					}
				} else {
					if excludeChannels != "" {
						t.Errorf("Expected no excludeChannels parameter, got excludeChannels=%s", excludeChannels)
					}
				}

				w.WriteHeader(tt.mockStatus)
				fmt.Fprint(w, tt.mockResponse)
			}))
			defer server.Close()

			client, err := NewClient(ClientConfig{
				APIToken: "test-token",
				BaseURL:  server.URL,
				Timeout:  30 * time.Second,
			})
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			appService := NewApplicationService(client)
			ctx := context.Background()

			result, err := appService.ListApplications(ctx, tt.opts)

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

			// Validate individual applications
			for _, app := range result.Applications {
				if err := app.Validate(); err != nil {
					t.Errorf("Application validation failed: %v", err)
				}
			}
		})
	}
}

func TestApplicationService_GetApplication(t *testing.T) {
	tests := []struct {
		name         string
		appID        string
		mockResponse string
		mockStatus   int
		expectError  bool
		expectedID   string
		expectedName string
	}{
		{
			name:  "successful get",
			appID: "app-1",
			mockResponse: `{
				"id": "app-1",
				"name": "Test App 1",
				"slug": "test-app-1",
				"team_id": "team-1",
				"team_name": "Test Team",
				"created_at": "2023-01-01T00:00:00Z",
				"updated_at": "2023-01-01T00:00:00Z",
				"description": "Test application 1",
				"is_active": true
			}`,
			mockStatus:   http.StatusOK,
			expectError:  false,
			expectedID:   "app-1",
			expectedName: "Test App 1",
		},
		{
			name:         "empty app ID",
			appID:        "",
			mockResponse: "",
			mockStatus:   0,
			expectError:  true,
		},
		{
			name:         "not found",
			appID:        "nonexistent-app",
			mockResponse: `{"message": "Not Found", "details": "Application not found"}`,
			mockStatus:   http.StatusNotFound,
			expectError:  true,
		},
		{
			name:         "unauthorized",
			appID:        "app-1",
			mockResponse: `{"message": "Unauthorized"}`,
			mockStatus:   http.StatusUnauthorized,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != testMethodGET {
					t.Errorf("Expected GET request, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/vendor/v3/app/%s", tt.appID)
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				if auth := r.Header.Get("Authorization"); auth == "" {
					t.Error("Expected Authorization header")
				}

				w.WriteHeader(tt.mockStatus)
				fmt.Fprint(w, tt.mockResponse)
			}))
			defer server.Close()

			client, err := NewClient(ClientConfig{
				APIToken: "test-token",
				BaseURL:  server.URL,
				Timeout:  30 * time.Second,
			})
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			appService := NewApplicationService(client)
			ctx := context.Background()

			result, err := appService.GetApplication(ctx, tt.appID)

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

			if result.Name != tt.expectedName {
				t.Errorf("Expected name %s, got %s", tt.expectedName, result.Name)
			}

			if err := result.Validate(); err != nil {
				t.Errorf("Application validation failed: %v", err)
			}
		})
	}
}

func TestApplicationService_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"applications": []}`)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		APIToken: "test-token",
		BaseURL:  server.URL,
		Timeout:  30 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	appService := NewApplicationService(client)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err = appService.ListApplications(ctx, nil)
	if err == nil {
		t.Error("Expected context cancellation error")
	}
}

func TestApplicationService_SearchApplications(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		opts          *ListApplicationsOptions
		mockResponse  string
		mockStatus    int
		expectError   bool
		expectedCount int
	}{
		{
			name:  "successful search by name",
			query: "Test App 1",
			opts:  nil,
			mockResponse: `{
				"applications": [
					{
						"id": "app-1",
						"name": "Test App 1",
						"slug": "test-app-1",
						"team_id": "team-1",
						"team_name": "Test Team",
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z",
						"description": "Test application 1",
						"is_active": true
					},
					{
						"id": "app-2",
						"name": "Different App",
						"slug": "different-app",
						"team_id": "team-1",
						"team_name": "Test Team",
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z",
						"description": "Different application",
						"is_active": true
					}
				]
			}`,
			mockStatus:    http.StatusOK,
			expectError:   false,
			expectedCount: 1,
		},
		{
			name:  "successful search by slug",
			query: "different",
			opts:  nil,
			mockResponse: `{
				"applications": [
					{
						"id": "app-1",
						"name": "Test App 1",
						"slug": "test-app-1",
						"team_id": "team-1",
						"team_name": "Test Team",
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z",
						"description": "Test application 1",
						"is_active": true
					},
					{
						"id": "app-2",
						"name": "Some App",
						"slug": "different-app",
						"team_id": "team-1",
						"team_name": "Test Team",
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z",
						"description": "Some application",
						"is_active": true
					}
				]
			}`,
			mockStatus:    http.StatusOK,
			expectError:   false,
			expectedCount: 1,
		},
		{
			name:  "successful search by description",
			query: "special feature",
			opts:  nil,
			mockResponse: `{
				"applications": [
					{
						"id": "app-1",
						"name": "Test App 1",
						"slug": "test-app-1",
						"team_id": "team-1",
						"team_name": "Test Team",
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z",
						"description": "App with special feature",
						"is_active": true
					},
					{
						"id": "app-2",
						"name": "Different App",
						"slug": "different-app",
						"team_id": "team-1",
						"team_name": "Test Team",
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z",
						"description": "Standard application",
						"is_active": true
					}
				]
			}`,
			mockStatus:    http.StatusOK,
			expectError:   false,
			expectedCount: 1,
		},
		{
			name:  "case insensitive search",
			query: "TEST",
			opts:  nil,
			mockResponse: `{
				"applications": [
					{
						"id": "app-1",
						"name": "Test App 1",
						"slug": "test-app-1",
						"team_id": "team-1",
						"team_name": "Test Team",
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z",
						"description": "Test application 1",
						"is_active": true
					}
				]
			}`,
			mockStatus:    http.StatusOK,
			expectError:   false,
			expectedCount: 1,
		},
		{
			name:  "no results found",
			query: "nonexistent",
			opts:  nil,
			mockResponse: `{
				"applications": [
					{
						"id": "app-1",
						"name": "Test App 1",
						"slug": "test-app-1",
						"team_id": "team-1",
						"team_name": "Test Team",
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z",
						"description": "Test application 1",
						"is_active": true
					}
				]
			}`,
			mockStatus:    http.StatusOK,
			expectError:   false,
			expectedCount: 0,
		},
		{
			name:        "empty query",
			query:       "",
			opts:        nil,
			expectError: true,
		},
		{
			name:        "whitespace only query",
			query:       "   ",
			opts:        nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != testMethodGET {
					t.Errorf("Expected GET request, got %s", r.Method)
				}
				if r.URL.Path != testPathApps {
					t.Errorf("Expected path /vendor/v3/apps, got %s", r.URL.Path)
				}

				w.WriteHeader(tt.mockStatus)
				fmt.Fprint(w, tt.mockResponse)
			}))
			defer server.Close()

			client, err := NewClient(ClientConfig{
				APIToken: "test-token",
				BaseURL:  server.URL,
				Timeout:  30 * time.Second,
			})
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			appService := NewApplicationService(client)
			ctx := context.Background()

			result, err := appService.SearchApplications(ctx, tt.query, tt.opts)

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

			// Validate that results actually match the query
			if len(result.Applications) > 0 {
				queryLower := strings.ToLower(tt.query)
				for _, app := range result.Applications {
					nameMatch := strings.Contains(strings.ToLower(app.Name), queryLower)
					slugMatch := strings.Contains(strings.ToLower(app.Slug), queryLower)
					descMatch := strings.Contains(strings.ToLower(app.Description), queryLower)

					if !nameMatch && !slugMatch && !descMatch {
						t.Errorf("Application %s does not match query %s", app.Name, tt.query)
					}
				}
			}
		})
	}
}
