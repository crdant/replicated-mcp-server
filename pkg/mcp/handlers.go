package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// Application Handlers

// handleListApplications handles the list_applications tool call
func (s *Server) handleListApplications(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Info("list_applications tool called", "arguments", request.GetArguments())

	// Extract arguments
	limit := 10 // default
	offset := 0 // default

	if args := request.GetArguments(); args != nil {
		if l, ok := args["limit"]; ok {
			if limitFloat, ok := l.(float64); ok {
				limit = int(limitFloat)
			}
		}
		if o, ok := args["offset"]; ok {
			if offsetFloat, ok := o.(float64); ok {
				offset = int(offsetFloat)
			}
		}
	}

	// Call API
	apps, err := s.apiClient.ListApplications(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(apps)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal applications: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonData)),
		},
	}, nil
}

// handleGetApplication handles the get_application tool call
func (s *Server) handleGetApplication(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Info("get_application tool called", "arguments", request.GetArguments())

	// Extract arguments
	args := request.GetArguments()
	if args == nil {
		return nil, fmt.Errorf("missing arguments")
	}

	appID, ok := args["app_id"].(string)
	if !ok {
		return nil, fmt.Errorf("app_id argument is required and must be a string")
	}

	// Call API
	app, err := s.apiClient.GetApplication(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(app)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal application: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonData)),
		},
	}, nil
}

// handleSearchApplications handles the search_applications tool call
func (s *Server) handleSearchApplications(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	s.logger.Info("search_applications tool called", "arguments", request.GetArguments())

	// Extract arguments
	args := request.GetArguments()
	if args == nil {
		return nil, fmt.Errorf("missing arguments")
	}

	query, ok := args["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query argument is required and must be a string")
	}

	limit := 10 // default
	if l, ok := args["limit"]; ok {
		if limitFloat, ok := l.(float64); ok {
			limit = int(limitFloat)
		}
	}

	// Call API
	apps, err := s.apiClient.SearchApplications(ctx, query, limit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to search applications: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(apps)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal applications: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonData)),
		},
	}, nil
}

// Release Handlers

// handleListReleases handles the list_releases tool call
func (s *Server) handleListReleases(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Info("list_releases tool called", "arguments", request.GetArguments())

	// Extract arguments
	args := request.GetArguments()
	if args == nil {
		return nil, fmt.Errorf("missing arguments")
	}

	appID, ok := args["app_id"].(string)
	if !ok {
		return nil, fmt.Errorf("app_id argument is required and must be a string")
	}

	limit := 10 // default
	offset := 0 // default

	if l, ok := args["limit"]; ok {
		if limitFloat, ok := l.(float64); ok {
			limit = int(limitFloat)
		}
	}
	if o, ok := args["offset"]; ok {
		if offsetFloat, ok := o.(float64); ok {
			offset = int(offsetFloat)
		}
	}

	// Call API
	releases, err := s.apiClient.ListReleases(ctx, appID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(releases)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal releases: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonData)),
		},
	}, nil
}

// handleGetRelease handles the get_release tool call
func (s *Server) handleGetRelease(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Info("get_release tool called", "arguments", request.GetArguments())

	// Extract arguments
	args := request.GetArguments()
	if args == nil {
		return nil, fmt.Errorf("missing arguments")
	}

	appID, ok := args["app_id"].(string)
	if !ok {
		return nil, fmt.Errorf("app_id argument is required and must be a string")
	}

	releaseID, ok := args["release_id"].(string)
	if !ok {
		return nil, fmt.Errorf("release_id argument is required and must be a string")
	}

	// Call API
	release, err := s.apiClient.GetRelease(ctx, appID, releaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get release: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(release)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal release: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonData)),
		},
	}, nil
}

// handleSearchReleases handles the search_releases tool call
func (s *Server) handleSearchReleases(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Info("search_releases tool called", "arguments", request.GetArguments())

	// Extract arguments
	args := request.GetArguments()
	if args == nil {
		return nil, fmt.Errorf("missing arguments")
	}

	appID, ok := args["app_id"].(string)
	if !ok {
		return nil, fmt.Errorf("app_id argument is required and must be a string")
	}

	query, ok := args["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query argument is required and must be a string")
	}

	limit := 10 // default
	if l, ok := args["limit"]; ok {
		if limitFloat, ok := l.(float64); ok {
			limit = int(limitFloat)
		}
	}

	// Call API
	releases, err := s.apiClient.SearchReleases(ctx, appID, query, limit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to search releases: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(releases)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal releases: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonData)),
		},
	}, nil
}

// Channel Handlers

// handleListChannels handles the list_channels tool call
func (s *Server) handleListChannels(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Info("list_channels tool called", "arguments", request.GetArguments())

	// Extract arguments
	args := request.GetArguments()
	if args == nil {
		return nil, fmt.Errorf("missing arguments")
	}

	appID, ok := args["app_id"].(string)
	if !ok {
		return nil, fmt.Errorf("app_id argument is required and must be a string")
	}

	limit := 10 // default
	offset := 0 // default

	if l, ok := args["limit"]; ok {
		if limitFloat, ok := l.(float64); ok {
			limit = int(limitFloat)
		}
	}
	if o, ok := args["offset"]; ok {
		if offsetFloat, ok := o.(float64); ok {
			offset = int(offsetFloat)
		}
	}

	// Call API
	channels, err := s.apiClient.ListChannels(ctx, appID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list channels: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(channels)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal channels: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonData)),
		},
	}, nil
}

// handleGetChannel handles the get_channel tool call
func (s *Server) handleGetChannel(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Info("get_channel tool called", "arguments", request.GetArguments())

	// Extract arguments
	args := request.GetArguments()
	if args == nil {
		return nil, fmt.Errorf("missing arguments")
	}

	appID, ok := args["app_id"].(string)
	if !ok {
		return nil, fmt.Errorf("app_id argument is required and must be a string")
	}

	channelID, ok := args["channel_id"].(string)
	if !ok {
		return nil, fmt.Errorf("channel_id argument is required and must be a string")
	}

	// Call API
	channel, err := s.apiClient.GetChannel(ctx, appID, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(channel)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal channel: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonData)),
		},
	}, nil
}

// handleSearchChannels handles the search_channels tool call
func (s *Server) handleSearchChannels(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Info("search_channels tool called", "arguments", request.GetArguments())

	// Extract arguments
	args := request.GetArguments()
	if args == nil {
		return nil, fmt.Errorf("missing arguments")
	}

	appID, ok := args["app_id"].(string)
	if !ok {
		return nil, fmt.Errorf("app_id argument is required and must be a string")
	}

	query, ok := args["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query argument is required and must be a string")
	}

	limit := 10 // default
	if l, ok := args["limit"]; ok {
		if limitFloat, ok := l.(float64); ok {
			limit = int(limitFloat)
		}
	}

	// Call API
	channels, err := s.apiClient.SearchChannels(ctx, appID, query, limit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to search channels: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(channels)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal channels: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonData)),
		},
	}, nil
}

// Customer Handlers

// handleListCustomers handles the list_customers tool call
func (s *Server) handleListCustomers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Info("list_customers tool called", "arguments", request.GetArguments())

	// Extract arguments
	args := request.GetArguments()
	if args == nil {
		return nil, fmt.Errorf("missing arguments")
	}

	appID, ok := args["app_id"].(string)
	if !ok {
		return nil, fmt.Errorf("app_id argument is required and must be a string")
	}

	limit := 10 // default
	offset := 0 // default

	if l, ok := args["limit"]; ok {
		if limitFloat, ok := l.(float64); ok {
			limit = int(limitFloat)
		}
	}
	if o, ok := args["offset"]; ok {
		if offsetFloat, ok := o.(float64); ok {
			offset = int(offsetFloat)
		}
	}

	// Call API
	customers, err := s.apiClient.ListCustomers(ctx, appID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(customers)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal customers: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonData)),
		},
	}, nil
}

// handleGetCustomer handles the get_customer tool call
func (s *Server) handleGetCustomer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Info("get_customer tool called", "arguments", request.GetArguments())

	// Extract arguments
	args := request.GetArguments()
	if args == nil {
		return nil, fmt.Errorf("missing arguments")
	}

	appID, ok := args["app_id"].(string)
	if !ok {
		return nil, fmt.Errorf("app_id argument is required and must be a string")
	}

	customerID, ok := args["customer_id"].(string)
	if !ok {
		return nil, fmt.Errorf("customer_id argument is required and must be a string")
	}

	// Call API
	customer, err := s.apiClient.GetCustomer(ctx, appID, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(customer)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal customer: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonData)),
		},
	}, nil
}

// handleSearchCustomers handles the search_customers tool call
func (s *Server) handleSearchCustomers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Info("search_customers tool called", "arguments", request.GetArguments())

	// Extract arguments
	args := request.GetArguments()
	if args == nil {
		return nil, fmt.Errorf("missing arguments")
	}

	appID, ok := args["app_id"].(string)
	if !ok {
		return nil, fmt.Errorf("app_id argument is required and must be a string")
	}

	query, ok := args["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query argument is required and must be a string")
	}

	limit := 10 // default
	if l, ok := args["limit"]; ok {
		if limitFloat, ok := l.(float64); ok {
			limit = int(limitFloat)
		}
	}

	// Call API
	customers, err := s.apiClient.SearchCustomers(ctx, appID, query, limit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to search customers: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(customers)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal customers: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonData)),
		},
	}, nil
}
