package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

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

// ListApplicationsOptions represents options for listing applications
type ListApplicationsOptions struct {
	ExcludeChannels bool `json:"exclude_channels,omitempty"`
}

// ApplicationList represents a list of applications
type ApplicationList struct {
	Applications []models.Application `json:"applications"`
}

// ListApplications retrieves all applications accessible to the authenticated team
func (s *ApplicationService) ListApplications(ctx context.Context, opts *ListApplicationsOptions) (*ApplicationList, error) {
	path := "/vendor/v3/apps"
	
	// Build query parameters
	if opts != nil && opts.ExcludeChannels {
		params := url.Values{}
		params.Set("excludeChannels", "true")
		path += "?" + params.Encode()
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
		"count", len(result.Applications))

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

