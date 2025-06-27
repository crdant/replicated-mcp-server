// Package api provides a client for the Replicated Vendor Portal API
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/crdant/replicated-mcp-server/pkg/config"
	"github.com/crdant/replicated-mcp-server/pkg/logging"
)

const (
	defaultBaseURL     = "https://api.replicated.com"
	userAgent          = "replicated-mcp-server"
	httpErrorThreshold = 400
)

// Client is the Replicated API client
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiToken   string
	logger     logging.Logger
}

// NewClient creates a new Replicated API client
func NewClient(cfg *config.Config, logger logging.Logger) (*Client, error) {
	if cfg.APIToken == "" {
		return nil, fmt.Errorf("API token is required")
	}

	baseURL := defaultBaseURL
	if cfg.Endpoint != "" {
		baseURL = cfg.Endpoint
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		baseURL:  baseURL,
		apiToken: cfg.APIToken,
		logger:   logger,
	}, nil
}

// doRequest performs an HTTP request and handles common response processing
//nolint:unparam // method parameter designed for extensibility, currently only GET is used
func (c *Client) doRequest(ctx context.Context, method, path string, params url.Values) ([]byte, error) {
	u := fmt.Sprintf("%s%s", c.baseURL, path)
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	c.logger.Debug("Making API request", "method", method, "url", u)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= httpErrorThreshold {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, errorResp.Error.Message)
		}
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// ListApplications retrieves a list of applications
func (c *Client) ListApplications(ctx context.Context, limit, offset int) ([]Application, error) {
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))
	}

	body, err := c.doRequest(ctx, "GET", "/v3/apps", params)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}

	var response ApplicationListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Applications, nil
}

// GetApplication retrieves a single application by ID or slug
func (c *Client) GetApplication(ctx context.Context, appIdentifier string) (*Application, error) {
	path := fmt.Sprintf("/v3/app/%s", url.PathEscape(appIdentifier))

	body, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	var response ApplicationResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response.Application, nil
}

// SearchApplications searches for applications by name
func (c *Client) SearchApplications(ctx context.Context, query string, limit, offset int) ([]Application, error) {
	params := url.Values{}
	params.Set("q", query)
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))
	}

	body, err := c.doRequest(ctx, "GET", "/v3/apps/search", params)
	if err != nil {
		return nil, fmt.Errorf("failed to search applications: %w", err)
	}

	var response ApplicationListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Applications, nil
}

// ListReleases retrieves releases for an application
func (c *Client) ListReleases(ctx context.Context, appIdentifier string, limit, offset int) ([]Release, error) {
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))
	}

	path := fmt.Sprintf("/v3/app/%s/releases", url.PathEscape(appIdentifier))

	body, err := c.doRequest(ctx, "GET", path, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	var response ReleaseListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Releases, nil
}

// GetRelease retrieves a single release
func (c *Client) GetRelease(ctx context.Context, appIdentifier, releaseIdentifier string) (*Release, error) {
	path := fmt.Sprintf("/v3/app/%s/release/%s",
		url.PathEscape(appIdentifier),
		url.PathEscape(releaseIdentifier))

	body, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get release: %w", err)
	}

	var response ReleaseResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response.Release, nil
}

// SearchReleases searches for releases by version
func (c *Client) SearchReleases(ctx context.Context, appIdentifier, query string, limit, offset int) ([]Release, error) {
	params := url.Values{}
	params.Set("q", query)
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))
	}

	path := fmt.Sprintf("/v3/app/%s/releases/search", url.PathEscape(appIdentifier))

	body, err := c.doRequest(ctx, "GET", path, params)
	if err != nil {
		return nil, fmt.Errorf("failed to search releases: %w", err)
	}

	var response ReleaseListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Releases, nil
}

// ListChannels retrieves channels for an application
func (c *Client) ListChannels(ctx context.Context, appIdentifier string, limit, offset int) ([]Channel, error) {
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))
	}

	path := fmt.Sprintf("/v3/app/%s/channels", url.PathEscape(appIdentifier))

	body, err := c.doRequest(ctx, "GET", path, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list channels: %w", err)
	}

	var response ChannelListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Channels, nil
}

// GetChannel retrieves a single channel
func (c *Client) GetChannel(ctx context.Context, appIdentifier, channelIdentifier string) (*Channel, error) {
	path := fmt.Sprintf("/v3/app/%s/channel/%s",
		url.PathEscape(appIdentifier),
		url.PathEscape(channelIdentifier))

	body, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}

	var response ChannelResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response.Channel, nil
}

// SearchChannels searches for channels by name
func (c *Client) SearchChannels(ctx context.Context, appIdentifier, query string, limit, offset int) ([]Channel, error) {
	params := url.Values{}
	params.Set("q", query)
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))
	}

	path := fmt.Sprintf("/v3/app/%s/channels/search", url.PathEscape(appIdentifier))

	body, err := c.doRequest(ctx, "GET", path, params)
	if err != nil {
		return nil, fmt.Errorf("failed to search channels: %w", err)
	}

	var response ChannelListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Channels, nil
}

// ListCustomers retrieves customers for an application
func (c *Client) ListCustomers(ctx context.Context, appIdentifier string, limit, offset int) ([]Customer, error) {
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))
	}

	path := fmt.Sprintf("/v3/app/%s/customers", url.PathEscape(appIdentifier))

	body, err := c.doRequest(ctx, "GET", path, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}

	var response CustomerListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Customers, nil
}

// GetCustomer retrieves a single customer
func (c *Client) GetCustomer(ctx context.Context, appIdentifier, customerIdentifier string) (*Customer, error) {
	path := fmt.Sprintf("/v3/app/%s/customer/%s",
		url.PathEscape(appIdentifier),
		url.PathEscape(customerIdentifier))

	body, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	var response CustomerResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response.Customer, nil
}

// SearchCustomers searches for customers by name or email
func (c *Client) SearchCustomers(ctx context.Context, appIdentifier, query string, limit, offset int) ([]Customer, error) {
	params := url.Values{}
	params.Set("q", query)
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))
	}

	path := fmt.Sprintf("/v3/app/%s/customers/search", url.PathEscape(appIdentifier))

	body, err := c.doRequest(ctx, "GET", path, params)
	if err != nil {
		return nil, fmt.Errorf("failed to search customers: %w", err)
	}

	var response CustomerListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Customers, nil
}
