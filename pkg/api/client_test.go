package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

// Test constants
const (
	testUserAgent = "replicated-mcp-server"
	testPath      = "/test"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  ClientConfig
		wantErr bool
	}{
		{
			name: "valid configuration",
			config: ClientConfig{
				APIToken: "valid-token",
				BaseURL:  "https://api.replicated.com",
				Timeout:  30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "missing API token",
			config: ClientConfig{
				BaseURL: "https://api.replicated.com",
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "missing base URL",
			config: ClientConfig{
				APIToken: "valid-token",
				Timeout:  30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "zero timeout defaults to 30s",
			config: ClientConfig{
				APIToken: "valid-token",
				BaseURL:  "https://api.replicated.com",
				Timeout:  0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if client == nil {
					t.Error("NewClient() returned nil client")
					return
				}
				if client.config.APIToken != tt.config.APIToken {
					t.Errorf("NewClient() APIToken = %v, want %v", client.config.APIToken, tt.config.APIToken)
				}
				if client.config.BaseURL != tt.config.BaseURL {
					t.Errorf("NewClient() BaseURL = %v, want %v", client.config.BaseURL, tt.config.BaseURL)
				}
				// Check default timeout
				if tt.config.Timeout == 0 && client.config.Timeout != 30*time.Second {
					t.Errorf("NewClient() default timeout = %v, want %v", client.config.Timeout, 30*time.Second)
				}
			}
		})
	}
}

func TestClient_Authentication(t *testing.T) {
	tests := []struct {
		name      string
		apiToken  string
		wantError bool
	}{
		{
			name:      "valid token",
			apiToken:  "valid-token-12345",
			wantError: false,
		},
		{
			name:      "empty token should be caught in NewClient",
			apiToken:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ClientConfig{
				APIToken: tt.apiToken,
				BaseURL:  "https://api.replicated.com",
				Timeout:  30 * time.Second,
			}

			client, err := NewClient(config)
			if (err != nil) != tt.wantError {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				// Test that the client has the expected headers
				headers := client.GetAuthHeaders()
				expectedAuth := tt.apiToken
				if headers.Get("Authorization") != expectedAuth {
					t.Errorf("Authorization header = %v, want %v", headers.Get("Authorization"), expectedAuth)
				}
				if headers.Get("User-Agent") != testUserAgent {
					t.Errorf("User-Agent header = %v, want %v", headers.Get("User-Agent"), testUserAgent)
				}
			}
		})
	}
}

func TestClient_HTTPMethods(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify authentication headers
		if auth := r.Header.Get("Authorization"); auth != "test-token" {
			t.Errorf("Expected Authorization header 'test-token', got '%s'", auth)
		}
		if ua := r.Header.Get("User-Agent"); ua != testUserAgent {
			t.Errorf("Expected User-Agent '%s', got '%s'", testUserAgent, ua)
		}

		// Respond based on method and path
		switch {
		case r.Method == "GET" && r.URL.Path == testPath:
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message": "success"}`)
		case r.Method == "POST" && r.URL.Path == testPath:
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"message": "created"}`)
		case r.Method == "PUT" && r.URL.Path == testPath:
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message": "updated"}`)
		case r.Method == "DELETE" && r.URL.Path == testPath:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
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

	ctx := context.Background()

	t.Run("GET request", func(t *testing.T) {
		resp, err := client.Get(ctx, testPath)
		if err != nil {
			t.Fatalf("GET request failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST request", func(t *testing.T) {
		body := strings.NewReader(`{"key": "value"}`)
		resp, err := client.Post(ctx, testPath, "application/json", body)
		if err != nil {
			t.Fatalf("POST request failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT request", func(t *testing.T) {
		body := strings.NewReader(`{"key": "updated"}`)
		resp, err := client.Put(ctx, testPath, "application/json", body)
		if err != nil {
			t.Fatalf("PUT request failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE request", func(t *testing.T) {
		resp, err := client.Delete(ctx, testPath)
		if err != nil {
			t.Fatalf("DELETE request failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", resp.StatusCode)
		}
	})
}

func TestClient_ErrorHandling(t *testing.T) {
	// Create a test server that returns various error responses
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/400":
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"message": "Bad Request", "details": "Invalid parameters"}`)
		case "/401":
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"message": "Unauthorized"}`)
		case "/404":
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"message": "Not Found"}`)
		case "/500":
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message": "Internal Server Error"}`)
		default:
			w.WriteHeader(http.StatusOK)
		}
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

	ctx := context.Background()

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "400 Bad Request",
			path:           "/400",
			expectedStatus: 400,
			expectError:    true,
		},
		{
			name:           "401 Unauthorized",
			path:           "/401",
			expectedStatus: 401,
			expectError:    true,
		},
		{
			name:           "404 Not Found",
			path:           "/404",
			expectedStatus: 404,
			expectError:    true,
		},
		{
			name:           "500 Internal Server Error",
			path:           "/500",
			expectedStatus: 500,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Get(ctx, tt.path)
			// We should still get a response, but it will have error status
			if err != nil {
				t.Fatalf("Unexpected request error: %v", err)
			}
			if resp == nil {
				t.Fatal("Expected response but got nil")
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			// Test that we can convert HTTP errors to API errors
			if resp.StatusCode >= 400 {
				apiErr := client.ConvertHTTPError(resp)
				if apiErr == nil {
					t.Error("Expected ConvertHTTPError to return an Error")
					return
				}
				if apiErr.StatusCode != tt.expectedStatus {
					t.Errorf("Expected Error status %d, got %d", tt.expectedStatus, apiErr.StatusCode)
				}
			}
		})
	}
}

func TestClient_ContextCancellation(t *testing.T) {
	// Create a test server with a slow response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
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

	// Create a context that cancels quickly
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	resp, err := client.Get(ctx, "/slow")
	if err == nil {
		t.Error("Expected context cancellation error")
	}
	if resp != nil {
		resp.Body.Close()
	}
}

func TestClient_Logging(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"message": "success"}`)
	}))
	defer server.Close()

	// Create a client with debug logging enabled
	client, err := NewClientWithLogger(ClientConfig{
		APIToken: "test-token",
		BaseURL:  server.URL,
		Timeout:  30 * time.Second,
	}, slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// This should trigger debug logging
	resp, err := client.Get(ctx, "/test")
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Test that the client has a logger
	if client.logger == nil {
		t.Error("Expected client to have a logger")
	}
}
