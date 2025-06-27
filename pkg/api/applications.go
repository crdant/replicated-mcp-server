package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/crdant/replicated-mcp-server/pkg/models"
)

// ApplicationService provides methods for interacting with application APIs
type ApplicationService struct {
	client *Client
}

// NewApplicationService creates a new ApplicationService
func NewApplicationService(client *Client) *ApplicationService {
	return &ApplicationService{
		client: client,
	}
}

// ListOptions represents options for listing operations
type ListOptions struct {
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Limit    int `json:"limit,omitempty"`
}

// ApplicationList represents a paginated list of applications
type ApplicationList struct {
	Applications []models.Application `json:"applications"`
	Page         int                  `json:"page"`
	PageSize     int                  `json:"page_size"`
	TotalCount   int                  `json:"total_count"`
	HasMore      bool                 `json:"has_more"`
}

// ListApplications retrieves a paginated list of applications
func (s *ApplicationService) ListApplications(ctx context.Context, opts *ListOptions) (*ApplicationList, error) {
	path := "/vendor/v3/apps"
	
	// Build query parameters
	if opts != nil {
		params := url.Values{}
		if opts.Page > 0 {
			params.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.PageSize > 0 {
			params.Set("page_size", strconv.Itoa(opts.PageSize))
		}
		if opts.Limit > 0 {
			params.Set("limit", strconv.Itoa(opts.Limit))
		}
		if len(params) > 0 {
			path += "?" + params.Encode()
		}
	}

	s.client.logger.DebugContext(ctx, "Listing applications", "path", path)

	resp, err := s.client.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		apiErr := s.client.ConvertHTTPError(resp)
		return nil, fmt.Errorf("API error: %w", apiErr)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result ApplicationList
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	s.client.logger.DebugContext(ctx, "Successfully listed applications", 
		"count", len(result.Applications),
		"page", result.Page,
		"total", result.TotalCount)

	return &result, nil
}

// GetApplication retrieves a specific application by ID
func (s *ApplicationService) GetApplication(ctx context.Context, id string) (*models.Application, error) {
	if id == "" {
		return nil, fmt.Errorf("application ID is required")
	}

	path := fmt.Sprintf("/vendor/v3/app/%s", id)

	s.client.logger.DebugContext(ctx, "Getting application", "app_id", id)

	resp, err := s.client.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		apiErr := s.client.ConvertHTTPError(resp)
		return nil, fmt.Errorf("API error: %w", apiErr)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result models.Application
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	s.client.logger.DebugContext(ctx, "Successfully retrieved application", 
		"app_id", result.ID,
		"app_name", result.Name)

	return &result, nil
}

// SearchApplications searches for applications using a query string
func (s *ApplicationService) SearchApplications(ctx context.Context, query string, opts *ListOptions) (*ApplicationList, error) {
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("search query is required")
	}

	path := "/vendor/v3/apps"
	params := url.Values{}
	params.Set("search", query)

	// Add pagination options
	if opts != nil {
		if opts.Page > 0 {
			params.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.PageSize > 0 {
			params.Set("page_size", strconv.Itoa(opts.PageSize))
		}
		if opts.Limit > 0 {
			params.Set("limit", strconv.Itoa(opts.Limit))
		}
	}

	path += "?" + params.Encode()

	s.client.logger.DebugContext(ctx, "Searching applications", "query", query, "path", path)

	resp, err := s.client.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to search applications: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		apiErr := s.client.ConvertHTTPError(resp)
		return nil, fmt.Errorf("API error: %w", apiErr)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result ApplicationList
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	s.client.logger.DebugContext(ctx, "Successfully searched applications", 
		"query", query,
		"count", len(result.Applications),
		"total", result.TotalCount)

	return &result, nil
}