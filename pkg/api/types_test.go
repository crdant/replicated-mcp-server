package api

import (
	"testing"
	"time"
)

func TestClientConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ClientConfig
		wantErr bool
	}{
		{
			name: "valid config",
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
			name: "both missing",
			config: ClientConfig{
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		apiError Error
		expected string
	}{
		{
			name: "error with details",
			apiError: Error{
				StatusCode: 400,
				Message:    "Bad Request",
				Details:    "Invalid parameters provided",
			},
			expected: "API error (status 400): Bad Request - Invalid parameters provided",
		},
		{
			name: "error without details",
			apiError: Error{
				StatusCode: 404,
				Message:    "Not Found",
			},
			expected: "API error (status 404): Not Found",
		},
		{
			name: "error with empty details",
			apiError: Error{
				StatusCode: 500,
				Message:    "Internal Server Error",
				Details:    "",
			},
			expected: "API error (status 500): Internal Server Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.apiError.Error()
			if result != tt.expected {
				t.Errorf("APIError.Error() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPaginatedResponse(t *testing.T) {
	// Test that PaginatedResponse works with different types
	t.Run("string slice", func(t *testing.T) {
		response := PaginatedResponse[string]{
			Data:       []string{"item1", "item2", "item3"},
			TotalCount: 100,
			Page:       1,
			PageSize:   3,
			HasMore:    true,
		}

		if len(response.Data) != 3 {
			t.Errorf("Expected 3 items, got %d", len(response.Data))
		}
		if response.Data[0] != "item1" {
			t.Errorf("Expected first item to be 'item1', got %v", response.Data[0])
		}
		if !response.HasMore {
			t.Error("Expected HasMore to be true")
		}
	})

	t.Run("int slice", func(t *testing.T) {
		response := PaginatedResponse[int]{
			Data:       []int{1, 2, 3},
			TotalCount: 50,
			Page:       2,
			PageSize:   3,
			HasMore:    false,
		}

		if len(response.Data) != 3 {
			t.Errorf("Expected 3 items, got %d", len(response.Data))
		}
		if response.Data[0] != 1 {
			t.Errorf("Expected first item to be 1, got %v", response.Data[0])
		}
		if response.HasMore {
			t.Error("Expected HasMore to be false")
		}
	})
}
