package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestApplicationService_ListApplications(t *testing.T) {
	tests := []struct {
		name           string
		opts           *ListOptions
		mockResponse   string
		mockStatus     int
		expectError    bool
		expectedCount  int
		expectedPage   int
		expectedTotal  int
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
				],
				"page": 1,
				"page_size": 20,
				"total_count": 2,
				"has_more": false
			}`,
			mockStatus:    http.StatusOK,
			expectError:   false,
			expectedCount: 2,
			expectedPage:  1,
			expectedTotal: 2,
		},
		{
			name: "successful list with pagination options",
			opts: &ListOptions{Page: 2, PageSize: 10},
			mockResponse: `{
				"applications": [
					{
						"id": "app-3",
						"name": "Test App 3",
						"slug": "test-app-3",
						"team_id": "team-1",
						"team_name": "Test Team",
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z",
						"description": "Test application 3",
						"is_active": true
					}
				],
				"page": 2,
				"page_size": 10,
				"total_count": 21,
				"has_more": true
			}`,
			mockStatus:    http.StatusOK,
			expectError:   false,
			expectedCount: 1,
			expectedPage:  2,
			expectedTotal: 21,
		},
		{
			name:         "empty list",
			opts:         nil,
			mockResponse: `{"applications": [], "page": 1, "page_size": 20, "total_count": 0, "has_more": false}`,
			mockStatus:   http.StatusOK,
			expectError:  false,
			expectedCount: 0,
			expectedPage: 1,
			expectedTotal: 0,
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
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}
				if r.URL.Path != "/vendor/v3/apps" {
					t.Errorf("Expected path /vendor/v3/apps, got %s", r.URL.Path)
				}

				// Verify authentication header
				if auth := r.Header.Get("Authorization"); auth == "" {
					t.Error("Expected Authorization header")
				}

				// Check query parameters if options provided
				if tt.opts != nil {
					query := r.URL.Query()
					if tt.opts.Page > 0 {
						if page := query.Get("page"); page != fmt.Sprintf("%d", tt.opts.Page) {
							t.Errorf("Expected page=%d, got page=%s", tt.opts.Page, page)
						}
					}
					if tt.opts.PageSize > 0 {
						if pageSize := query.Get("page_size"); pageSize != fmt.Sprintf("%d", tt.opts.PageSize) {
							t.Errorf("Expected page_size=%d, got page_size=%s", tt.opts.PageSize, pageSize)
						}
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

			if result.Page != tt.expectedPage {
				t.Errorf("Expected page %d, got %d", tt.expectedPage, result.Page)
			}

			if result.TotalCount != tt.expectedTotal {
				t.Errorf("Expected total count %d, got %d", tt.expectedTotal, result.TotalCount)
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
				if r.Method != "GET" {
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

func TestApplicationService_SearchApplications(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		opts           *ListOptions
		mockResponse   string
		mockStatus     int
		expectError    bool
		expectedCount  int
		expectedQuery  string
	}{
		{
			name:  "successful search",
			query: "test",
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
				],
				"page": 1,
				"page_size": 20,
				"total_count": 1,
				"has_more": false
			}`,
			mockStatus:    http.StatusOK,
			expectError:   false,
			expectedCount: 1,
			expectedQuery: "test",
		},
		{
			name:          "empty query",
			query:         "",
			opts:          nil,
			mockResponse:  "",
			mockStatus:    0,
			expectError:   true,
		},
		{
			name:  "no results",
			query: "nonexistent",
			opts:  nil,
			mockResponse: `{
				"applications": [],
				"page": 1,
				"page_size": 20,
				"total_count": 0,
				"has_more": false
			}`,
			mockStatus:    http.StatusOK,
			expectError:   false,
			expectedCount: 0,
			expectedQuery: "nonexistent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}
				if r.URL.Path != "/vendor/v3/apps" {
					t.Errorf("Expected path /vendor/v3/apps, got %s", r.URL.Path)
				}

				// Verify query parameter
				query := r.URL.Query()
				if searchQuery := query.Get("search"); searchQuery != tt.expectedQuery {
					t.Errorf("Expected search=%s, got search=%s", tt.expectedQuery, searchQuery)
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
		})
	}
}

func TestApplicationService_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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