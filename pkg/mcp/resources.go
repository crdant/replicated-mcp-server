package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// resourceDefinition represents a complete resource definition with its handler function.
type resourceDefinition struct {
	definition *mcp.Resource
	handler    server.ResourceHandlerFunc
}

// defineResources returns all MCP resource definitions for Replicated entities.
// Resources provide standardized access to Replicated data through URI-based addressing.
//
// Resource URI patterns:
// - Applications: replicated://applications/{app_id}
// - Releases: replicated://applications/{app_id}/releases/{release_id}
// - Channels: replicated://applications/{app_id}/channels/{channel_id}
// - Customers: replicated://applications/{app_id}/customers/{customer_id}
//
// Each resource includes:
// - Standardized URI scheme for consistent addressing
// - MIME type specification for content format
// - Comprehensive metadata and descriptions
// - Empty handler that returns placeholder responses
//
// Returns:
//   []resourceDefinition: All resource definitions with handlers
func (s *Server) defineResources() []resourceDefinition {
	return []resourceDefinition{
		s.defineApplicationResource(),
		s.defineReleaseResource(),
		s.defineChannelResource(),
		s.defineCustomerResource(),
	}
}

// defineApplicationResource creates the application resource definition.
// Provides access to application data through the replicated://applications/{app_id} URI pattern.
func (s *Server) defineApplicationResource() resourceDefinition {
	resource := mcp.NewResource(
		"replicated://applications/{app_id}",
		"Application Data",
		mcp.WithResourceDescription("Access to detailed application information including configuration, status, and metadata from the Replicated Vendor Portal"),
		mcp.WithMIMEType("application/json"),
	)

	handler := func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		s.logger.Info("Application resource accessed", "uri", request.Params.URI)
		
		// TODO: Implement actual application resource retrieval in Step 7
		// Return empty slice for now - actual implementation will be in Step 7
		// The MCP library expects this signature for resource handlers
		return []mcp.ResourceContents{}, nil
	}

	return resourceDefinition{definition: &resource, handler: handler}
}

// defineReleaseResource creates the release resource definition.
// Provides access to release data through the replicated://applications/{app_id}/releases/{release_id} URI pattern.
func (s *Server) defineReleaseResource() resourceDefinition {
	resource := mcp.NewResource(
		"replicated://applications/{app_id}/releases/{release_id}",
		"Release Data",
		mcp.WithResourceDescription("Access to detailed release information including version, manifests, deployment configuration, and changelog from the Replicated Vendor Portal"),
		mcp.WithMIMEType("application/json"),
	)

	handler := func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		s.logger.Info("Release resource accessed", "uri", request.Params.URI)
		
		// TODO: Implement actual release resource retrieval in Step 7
		// Return empty slice for now - actual implementation will be in Step 7
		return []mcp.ResourceContents{}, nil
	}

	return resourceDefinition{definition: &resource, handler: handler}
}

// defineChannelResource creates the channel resource definition.
// Provides access to channel data through the replicated://applications/{app_id}/channels/{channel_id} URI pattern.
func (s *Server) defineChannelResource() resourceDefinition {
	resource := mcp.NewResource(
		"replicated://applications/{app_id}/channels/{channel_id}",
		"Channel Data",
		mcp.WithResourceDescription("Access to detailed channel information including release assignments, customer adoption, and deployment policies from the Replicated Vendor Portal"),
		mcp.WithMIMEType("application/json"),
	)

	handler := func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		s.logger.Info("Channel resource accessed", "uri", request.Params.URI)
		
		// TODO: Implement actual channel resource retrieval in Step 7
		// TODO: Implement actual channel resource retrieval in Step 7
		// Return empty slice for now - actual implementation will be in Step 7
		return []mcp.ResourceContents{}, nil
	}

	return resourceDefinition{definition: &resource, handler: handler}
}

// defineCustomerResource creates the customer resource definition.
// Provides access to customer data through the replicated://applications/{app_id}/customers/{customer_id} URI pattern.
func (s *Server) defineCustomerResource() resourceDefinition {
	resource := mcp.NewResource(
		"replicated://applications/{app_id}/customers/{customer_id}",
		"Customer Data",
		mcp.WithResourceDescription("Access to detailed customer information including license details, deployment status, and usage analytics from the Replicated Vendor Portal"),
		mcp.WithMIMEType("application/json"),
	)

	handler := func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		s.logger.Info("Customer resource accessed", "uri", request.Params.URI)
		
		// TODO: Implement actual customer resource retrieval in Step 7
		// TODO: Implement actual customer resource retrieval in Step 7
		// Return empty slice for now - actual implementation will be in Step 7
		return []mcp.ResourceContents{}, nil
	}

	return resourceDefinition{definition: &resource, handler: handler}
}