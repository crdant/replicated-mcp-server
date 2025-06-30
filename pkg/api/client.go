// Package api provides HTTP client functionality for communicating with the Replicated Vendor Portal API.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

// Constants for HTTP client configuration
const (
	DefaultTimeout     = 30 * time.Second
	DefaultUserAgent   = "replicated-mcp-server"
	HTTPErrorThreshold = 400
)

// Client provides HTTP client functionality for the Replicated API
type Client struct {
	config     ClientConfig
	httpClient *http.Client
	logger     *slog.Logger
	
	// Services
	Applications *ApplicationService
}

// NewClient creates a new API client with the given configuration
func NewClient(config ClientConfig) (*Client, error) {
	// Use a no-op logger by default
	return NewClientWithLogger(config, slog.New(slog.NewTextHandler(io.Discard, nil)))
}

// NewClientWithLogger creates a new API client with the given configuration and logger
func NewClientWithLogger(config ClientConfig, logger *slog.Logger) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Set default timeout if not specified
	if config.Timeout == 0 {
		config.Timeout = DefaultTimeout
	}

	client := &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		logger: logger,
	}
	
	// Initialize services
	client.Applications = NewApplicationService(client)

	return client, nil
}

// GetAuthHeaders returns the authentication headers for API requests
func (c *Client) GetAuthHeaders() http.Header {
	headers := make(http.Header)
	headers.Set("Authorization", c.config.APIToken)
	headers.Set("User-Agent", DefaultUserAgent)
	return headers
}

// makeRequest creates and executes an HTTP request with proper authentication
func (c *Client) makeRequest(
	ctx context.Context, method, path, contentType string, body io.Reader,
) (*http.Response, error) {
	// Build full URL
	baseURL, err := url.Parse(c.config.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	fullURL, err := baseURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	// Log the request
	c.logger.DebugContext(ctx, "Making API request",
		"method", method,
		"url", fullURL.String(),
		"content_type", contentType,
	)

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, fullURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers
	headers := c.GetAuthHeaders()
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Set content type if provided
	if contentType != "" && body != nil {
		req.Header.Set("Content-Type", contentType)
	}

	// Execute request
	start := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(start)

	if err != nil {
		c.logger.ErrorContext(ctx, "API request failed",
			"method", method,
			"url", fullURL.String(),
			"duration", duration,
			"error", err,
		)
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Log the response
	c.logger.DebugContext(ctx, "API request completed",
		"method", method,
		"url", fullURL.String(),
		"status", resp.StatusCode,
		"duration", duration,
	)

	return resp, nil
}

// Get performs a GET request to the specified path
func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.makeRequest(ctx, "GET", path, "", nil)
}

// Post performs a POST request to the specified path
func (c *Client) Post(ctx context.Context, path, contentType string, body io.Reader) (*http.Response, error) {
	return c.makeRequest(ctx, "POST", path, contentType, body)
}

// Put performs a PUT request to the specified path
func (c *Client) Put(ctx context.Context, path, contentType string, body io.Reader) (*http.Response, error) {
	return c.makeRequest(ctx, "PUT", path, contentType, body)
}

// Delete performs a DELETE request to the specified path
func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	return c.makeRequest(ctx, "DELETE", path, "", nil)
}

// ConvertHTTPError converts an HTTP error response to an Error
func (c *Client) ConvertHTTPError(resp *http.Response) *Error {
	if resp.StatusCode < HTTPErrorThreshold {
		return nil
	}

	apiError := &Error{
		StatusCode: resp.StatusCode,
		Message:    http.StatusText(resp.StatusCode),
	}

	// Try to parse JSON error response
	if resp.Body != nil {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			var errorResponse struct {
				Message string `json:"message"`
				Details string `json:"details"`
			}
			if json.Unmarshal(body, &errorResponse) == nil {
				if errorResponse.Message != "" {
					apiError.Message = errorResponse.Message
				}
				apiError.Details = errorResponse.Details
			}
		}
	}

	return apiError
}
