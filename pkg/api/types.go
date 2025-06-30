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

// ListOptions provides common options for list operations
type ListOptions struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// SearchOptions provides common options for search operations
type SearchOptions struct {
	Query  string `json:"query"`
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

// Application represents a Replicated application
type Application struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Release represents a release in the Replicated Vendor Portal
type Release struct {
	ID           string    `json:"id"`
	Sequence     int       `json:"sequence"`
	Version      string    `json:"version"`
	ReleaseNotes string    `json:"release_notes,omitempty"`
	Required     bool      `json:"required"`
	Status       string    `json:"status,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Channel represents a release channel
type Channel struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Customer represents a customer in the Replicated Vendor Portal
type Customer struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email,omitempty"`
	Type        string    `json:"type,omitempty"`
	Status      string    `json:"status,omitempty"`
	LicenseID   string    `json:"license_id,omitempty"`
	ChannelID   string    `json:"channel_id,omitempty"`
	ChannelName string    `json:"channel_name,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
