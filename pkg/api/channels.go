package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

// ChannelService provides methods for interacting with release channels
type ChannelService struct {
	client *Client
}

// NewChannelService creates a new channel service
func NewChannelService(client *Client) *ChannelService {
	return &ChannelService{client: client}
}

// List retrieves a paginated list of channels for a specific application
func (s *ChannelService) List(ctx context.Context, appID string, opts *ListOptions) (*PaginatedResponse[Channel], error) {
	if appID == "" {
		return nil, fmt.Errorf("application ID is required")
	}

	path := fmt.Sprintf("/v3/app/%s/channels", url.PathEscape(appID))

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
		return nil, fmt.Errorf("failed to list channels: %w", err)
	}
	defer resp.Body.Close()

	if apiErr := s.client.ConvertHTTPError(resp); apiErr != nil {
		return nil, apiErr
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result PaginatedResponse[Channel]
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Get retrieves a specific channel by ID
func (s *ChannelService) Get(ctx context.Context, appID, channelID string) (*Channel, error) {
	if appID == "" {
		return nil, fmt.Errorf("application ID is required")
	}
	if channelID == "" {
		return nil, fmt.Errorf("channel ID is required")
	}

	path := fmt.Sprintf("/v3/app/%s/channel/%s", url.PathEscape(appID), url.PathEscape(channelID))

	resp, err := s.client.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %w", err)
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
		Channel Channel `json:"channel"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result.Channel, nil
}

// Search searches for channels within an application based on query criteria
func (s *ChannelService) Search(ctx context.Context, appID string, opts *SearchOptions) (*PaginatedResponse[Channel], error) {
	if appID == "" {
		return nil, fmt.Errorf("application ID is required")
	}
	if opts == nil || strings.TrimSpace(opts.Query) == "" {
		return nil, fmt.Errorf("search query is required")
	}

	// Since there's no dedicated search endpoint for channels, use list and filter client-side
	listOpts := &ListOptions{
		Limit:  100, // Get more results for better search coverage
		Offset: opts.Offset,
	}

	channels, err := s.List(ctx, appID, listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list channels for search: %w", err)
	}

	// Filter channels client-side based on query
	var filteredChannels []Channel
	queryLower := strings.ToLower(strings.TrimSpace(opts.Query))

	for _, channel := range channels.Data {
		// Search in name, slug, and description (case-insensitive)
		if strings.Contains(strings.ToLower(channel.Name), queryLower) ||
			strings.Contains(strings.ToLower(channel.Slug), queryLower) ||
			strings.Contains(strings.ToLower(channel.Description), queryLower) {
			filteredChannels = append(filteredChannels, channel)
		}
	}

	// Apply limit if specified
	if opts.Limit > 0 && len(filteredChannels) > opts.Limit {
		filteredChannels = filteredChannels[:opts.Limit]
	}

	result := &PaginatedResponse[Channel]{
		Data:       filteredChannels,
		TotalCount: len(filteredChannels),
		Page:       1, // Since we're doing client-side filtering
		PageSize:   len(filteredChannels),
		HasMore:    false, // We're not implementing true pagination for search
	}

	return result, nil
}