package api

import (
	"fmt"
	"time"
)

// ClientConfig holds configuration for the API client
type ClientConfig struct {
	APIToken string
	BaseURL  string
	Timeout  time.Duration
}

// Validate ensures the configuration is valid
func (c ClientConfig) Validate() error {
	if c.APIToken == "" {
		return fmt.Errorf("API token is required")
	}
	if c.BaseURL == "" {
		return fmt.Errorf("base URL is required")
	}
	return nil
}

// Error represents an error response from the API
type Error struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
}

func (e Error) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("API error (status %d): %s - %s", e.StatusCode, e.Message, e.Details)
	}
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

// PaginatedResponse wraps paginated API responses
type PaginatedResponse[T any] struct {
	Data       []T  `json:"data"`
	TotalCount int  `json:"total_count"`
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	HasMore    bool `json:"has_more"`
}
