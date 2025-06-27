package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/crdant/replicated-mcp-server/pkg/models"
)

// ListOptions represents pagination and filtering options for API requests
type ListOptions struct {
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}

// ApplicationList represents a paginated list of applications
type ApplicationList struct {
	Applications []models.Application `json:"applications"`
	TotalCount   int                  `json:"total_count"`
	Page         int                  `json:"page"`
	PageSize     int                  `json:"page_size"`
	HasMore      bool                 `json:"has_more"`
}

// parseJSONResponse is a helper function that handles common HTTP response processing
func (c *Client) parseJSONResponse(ctx context.Context, resp *http.Response, target interface{}, operation string) error {
	defer resp.Body.Close()

	// Check for HTTP errors
	if apiErr := c.ConvertHTTPError(resp); apiErr != nil {
		return apiErr
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body for %s: %w", operation, err)
	}

	// Parse JSON response
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to parse JSON response for %s: %w", operation, err)
	}

	return nil
}

// ListApplications retrieves a paginated list of applications from the Replicated API
func (c *Client) ListApplications(ctx context.Context, opts *ListOptions) (*ApplicationList, error) {
	// Build query parameters
	params := url.Values{}
	if opts != nil {
		if opts.Page > 0 {
			params.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.PageSize > 0 {
			params.Set("page_size", strconv.Itoa(opts.PageSize))
		}
	}

	// Build the URL path
	path := "/v1/applications"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	// Make the HTTP request
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}

	// Parse JSON response
	var appList ApplicationList
	if err := c.parseJSONResponse(ctx, resp, &appList, "list applications"); err != nil {
		return nil, err
	}

	// Validate applications
	for i, app := range appList.Applications {
		if err := app.Validate(); err != nil {
			c.logger.WarnContext(ctx, "Application validation failed",
				"index", i,
				"app_id", app.ID,
				"error", err,
			)
		}
	}

	c.logger.InfoContext(ctx, "Successfully listed applications",
		"count", len(appList.Applications),
		"total_count", appList.TotalCount,
		"page", appList.Page,
	)

	return &appList, nil
}

// GetApplication retrieves a single application by ID from the Replicated API
func (c *Client) GetApplication(ctx context.Context, id string) (*models.Application, error) {
	// Validate input
	if id == "" {
		return nil, fmt.Errorf("application ID is required")
	}

	// Build the URL path
	path := fmt.Sprintf("/v1/applications/%s", id)

	// Make the HTTP request
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get application %s: %w", id, err)
	}

	// Parse JSON response
	var app models.Application
	if err := c.parseJSONResponse(ctx, resp, &app, "get application"); err != nil {
		return nil, err
	}

	// Validate application
	if err := app.Validate(); err != nil {
		c.logger.WarnContext(ctx, "Application validation failed",
			"app_id", app.ID,
			"error", err,
		)
	}

	c.logger.InfoContext(ctx, "Successfully retrieved application",
		"app_id", app.ID,
		"app_name", app.Name,
	)

	return &app, nil
}

// SearchApplications searches for applications using a query string with optional pagination
func (c *Client) SearchApplications(ctx context.Context, query string, opts *ListOptions) (*ApplicationList, error) {
	// Validate input
	if query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	// Build query parameters
	params := url.Values{}
	params.Set("q", query)

	if opts != nil {
		if opts.Page > 0 {
			params.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.PageSize > 0 {
			params.Set("page_size", strconv.Itoa(opts.PageSize))
		}
	}

	// Build the URL path
	path := "/v1/applications/search?" + params.Encode()

	// Make the HTTP request
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to search applications: %w", err)
	}

	// Parse JSON response
	var appList ApplicationList
	if err := c.parseJSONResponse(ctx, resp, &appList, "search applications"); err != nil {
		return nil, err
	}

	// Validate applications
	for i, app := range appList.Applications {
		if err := app.Validate(); err != nil {
			c.logger.WarnContext(ctx, "Application validation failed",
				"index", i,
				"app_id", app.ID,
				"error", err,
			)
		}
	}

	c.logger.InfoContext(ctx, "Successfully searched applications",
		"query", query,
		"count", len(appList.Applications),
		"total_count", appList.TotalCount,
		"page", appList.Page,
	)

	return &appList, nil
}
