package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
)

// ReleaseService provides methods for interacting with Replicated releases
type ReleaseService struct {
	client *Client
}

// NewReleaseService creates a new release service
func NewReleaseService(client *Client) *ReleaseService {
	return &ReleaseService{client: client}
}

// List retrieves a paginated list of releases for a specific application
func (s *ReleaseService) List(ctx context.Context, appID string, opts *ListOptions) (*PaginatedResponse[Release], error) {
	if appID == "" {
		return nil, fmt.Errorf("application ID is required")
	}

	path := fmt.Sprintf("/v3/app/%s/releases", url.PathEscape(appID))

	// Build query parameters
	params := url.Values{}
	if opts != nil {
		if opts.Limit > 0 {
			params.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Offset > 0 {
			params.Set("offset", strconv.Itoa(opts.Offset))
		}
	}

	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	resp, err := s.client.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}
	defer resp.Body.Close()

	if apiErr := s.client.ConvertHTTPError(resp); apiErr != nil {
		return nil, apiErr
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result PaginatedResponse[Release]
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Get retrieves a specific release by ID or sequence number
func (s *ReleaseService) Get(ctx context.Context, appID, releaseID string) (*Release, error) {
	if appID == "" {
		return nil, fmt.Errorf("application ID is required")
	}
	if releaseID == "" {
		return nil, fmt.Errorf("release ID is required")
	}

	path := fmt.Sprintf("/v3/app/%s/release/%s", url.PathEscape(appID), url.PathEscape(releaseID))

	resp, err := s.client.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get release: %w", err)
	}
	defer resp.Body.Close()

	if apiErr := s.client.ConvertHTTPError(resp); apiErr != nil {
		return nil, apiErr
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Release Release `json:"release"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result.Release, nil
}

// Search searches for releases within an application based on query criteria
func (s *ReleaseService) Search(
	ctx context.Context, appID string, opts *SearchOptions,
) (*PaginatedResponse[Release], error) {
	if appID == "" {
		return nil, fmt.Errorf("application ID is required")
	}
	if opts == nil || opts.Query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	path := fmt.Sprintf("/v3/app/%s/releases/search", url.PathEscape(appID))

	// Build query parameters
	params := url.Values{}
	params.Set("q", opts.Query)
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		params.Set("offset", strconv.Itoa(opts.Offset))
	}

	path += "?" + params.Encode()

	resp, err := s.client.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to search releases: %w", err)
	}
	defer resp.Body.Close()

	if apiErr := s.client.ConvertHTTPError(resp); apiErr != nil {
		return nil, apiErr
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result PaginatedResponse[Release]
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
