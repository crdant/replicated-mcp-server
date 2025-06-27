package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// toolDefinition represents a complete tool definition with its handler function.
type toolDefinition struct {
	definition *mcp.Tool
	handler    server.ToolHandlerFunc
}

// defineTools returns all Phase 1 tools with their schemas and empty handler implementations.
// The actual business logic will be implemented in Step 7 (MCP Handlers).
//
// Tools are organized into four categories:
// - Application tools: list, get, search applications
// - Release tools: list, get, search releases  
// - Channel tools: list, get, search channels
// - Customer tools: list, get, search customers
//
// Note: ID parameters in tools accept both IDs and slugs (e.g., app_id accepts both 
// application IDs and application slugs). Handlers determine the parameter type at runtime.
//
// Each tool includes:
// - Proper JSON schema validation for arguments
// - Comprehensive documentation
// - Empty handler that returns placeholder responses
//
// Returns:
//   []toolDefinition: All tool definitions with handlers
func (s *Server) defineTools() []toolDefinition {
	return []toolDefinition{
		// Application Tools
		s.defineListApplicationsTool(),
		s.defineGetApplicationTool(),
		s.defineSearchApplicationsTool(),

		// Release Tools
		s.defineListReleasesTool(),
		s.defineGetReleaseTool(),
		s.defineSearchReleasesTool(),

		// Channel Tools
		s.defineListChannelsTool(),
		s.defineGetChannelTool(),
		s.defineSearchChannelsTool(),

		// Customer Tools
		s.defineListCustomersTool(),
		s.defineGetCustomerTool(),
		s.defineSearchCustomersTool(),
	}
}

// Application Tools

// defineListApplicationsTool creates the list_applications tool definition.
// Lists all applications accessible to the authenticated user.
func (s *Server) defineListApplicationsTool() toolDefinition {
	tool := mcp.NewTool("list_applications",
		mcp.WithDescription("List all applications in the Replicated Vendor Portal. Returns basic information about each application including ID, name, and status."),
		mcp.WithNumber("limit", 
			mcp.Description("Maximum number of applications to return (1-100)"),
			mcp.Min(1),
			mcp.Max(100),
		),
		mcp.WithNumber("offset", 
			mcp.Description("Number of applications to skip for pagination"),
			mcp.Min(0),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		s.logger.Info("list_applications tool called", "arguments", request.GetArguments())
		
		// TODO: Implement actual application listing in Step 7
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Applications listing will be implemented in Step 7 (MCP Handlers)"),
			},
		}, nil
	}

	return toolDefinition{definition: &tool, handler: handler}
}

// defineGetApplicationTool creates the get_application tool definition.
// Retrieves detailed information about a specific application.
func (s *Server) defineGetApplicationTool() toolDefinition {
	tool := mcp.NewTool("get_application",
		mcp.WithDescription("Get detailed information about a specific application by ID. Returns comprehensive application data including configuration and metadata."),
		mcp.WithString("app_id", 
			mcp.Required(),
			mcp.Description("The unique identifier of the application"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		s.logger.Info("get_application tool called", "arguments", request.GetArguments())
		
		// TODO: Implement actual application retrieval in Step 7
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Application details will be implemented in Step 7 (MCP Handlers)"),
			},
		}, nil
	}

	return toolDefinition{definition: &tool, handler: handler}
}

// defineSearchApplicationsTool creates the search_applications tool definition.
// Searches applications based on name or other criteria.
func (s *Server) defineSearchApplicationsTool() toolDefinition {
	tool := mcp.NewTool("search_applications",
		mcp.WithDescription("Search applications by name or other criteria. Returns matching applications with relevance scoring."),
		mcp.WithString("query", 
			mcp.Required(),
			mcp.Description("Search query string to match against application names and descriptions"),
		),
		mcp.WithNumber("limit", 
			mcp.Description("Maximum number of results to return (1-50)"),
			mcp.Min(1),
			mcp.Max(50),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		s.logger.Info("search_applications tool called", "arguments", request.GetArguments())
		
		// TODO: Implement actual application search in Step 7
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Application search will be implemented in Step 7 (MCP Handlers)"),
			},
		}, nil
	}

	return toolDefinition{definition: &tool, handler: handler}
}

// Release Tools

// defineListReleasesTool creates the list_releases tool definition.
// Lists releases for a specific application.
func (s *Server) defineListReleasesTool() toolDefinition {
	tool := mcp.NewTool("list_releases",
		mcp.WithDescription("List releases for a specific application. Returns release information including version, status, and deployment details."),
		mcp.WithString("app_id", 
			mcp.Required(),
			mcp.Description("The unique identifier of the application"),
		),
		mcp.WithNumber("limit", 
			mcp.Description("Maximum number of releases to return (1-100)"),
			mcp.Min(1),
			mcp.Max(100),
		),
		mcp.WithNumber("offset", 
			mcp.Description("Number of releases to skip for pagination"),
			mcp.Min(0),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		s.logger.Info("list_releases tool called", "arguments", request.GetArguments())
		
		// TODO: Implement actual release listing in Step 7
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Release listing will be implemented in Step 7 (MCP Handlers)"),
			},
		}, nil
	}

	return toolDefinition{definition: &tool, handler: handler}
}

// defineGetReleaseTool creates the get_release tool definition.
// Retrieves detailed information about a specific release.
func (s *Server) defineGetReleaseTool() toolDefinition {
	tool := mcp.NewTool("get_release",
		mcp.WithDescription("Get detailed information about a specific release by ID. Returns comprehensive release data including manifests and deployment configuration."),
		mcp.WithString("app_id", 
			mcp.Required(),
			mcp.Description("The unique identifier of the application"),
		),
		mcp.WithString("release_id", 
			mcp.Required(),
			mcp.Description("The unique identifier of the release"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		s.logger.Info("get_release tool called", "arguments", request.GetArguments())
		
		// TODO: Implement actual release retrieval in Step 7
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Release details will be implemented in Step 7 (MCP Handlers)"),
			},
		}, nil
	}

	return toolDefinition{definition: &tool, handler: handler}
}

// defineSearchReleasesTool creates the search_releases tool definition.
// Searches releases based on version or other criteria.
func (s *Server) defineSearchReleasesTool() toolDefinition {
	tool := mcp.NewTool("search_releases",
		mcp.WithDescription("Search releases by version or other criteria within a specific application. Returns matching releases with relevance scoring."),
		mcp.WithString("app_id", 
			mcp.Required(),
			mcp.Description("The unique identifier of the application"),
		),
		mcp.WithString("query", 
			mcp.Required(),
			mcp.Description("Search query string to match against release versions and descriptions"),
		),
		mcp.WithNumber("limit", 
			mcp.Description("Maximum number of results to return (1-50)"),
			mcp.Min(1),
			mcp.Max(50),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		s.logger.Info("search_releases tool called", "arguments", request.GetArguments())
		
		// TODO: Implement actual release search in Step 7
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Release search will be implemented in Step 7 (MCP Handlers)"),
			},
		}, nil
	}

	return toolDefinition{definition: &tool, handler: handler}
}

// Channel Tools

// defineListChannelsTool creates the list_channels tool definition.
// Lists all channels for a specific application.
func (s *Server) defineListChannelsTool() toolDefinition {
	tool := mcp.NewTool("list_channels",
		mcp.WithDescription("List channels for a specific application. Returns channel information including name, release assignments, and customer adoption."),
		mcp.WithString("app_id", 
			mcp.Required(),
			mcp.Description("The unique identifier of the application"),
		),
		mcp.WithNumber("limit", 
			mcp.Description("Maximum number of channels to return (1-100)"),
			mcp.Min(1),
			mcp.Max(100),
		),
		mcp.WithNumber("offset", 
			mcp.Description("Number of channels to skip for pagination"),
			mcp.Min(0),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		s.logger.Info("list_channels tool called", "arguments", request.GetArguments())
		
		// TODO: Implement actual channel listing in Step 7
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Channel listing will be implemented in Step 7 (MCP Handlers)"),
			},
		}, nil
	}

	return toolDefinition{definition: &tool, handler: handler}
}

// defineGetChannelTool creates the get_channel tool definition.
// Retrieves detailed information about a specific channel.
func (s *Server) defineGetChannelTool() toolDefinition {
	tool := mcp.NewTool("get_channel",
		mcp.WithDescription("Get detailed information about a specific channel by ID. Returns comprehensive channel data including release history and customer assignments."),
		mcp.WithString("app_id", 
			mcp.Required(),
			mcp.Description("The unique identifier of the application"),
		),
		mcp.WithString("channel_id", 
			mcp.Required(),
			mcp.Description("The unique identifier of the channel"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		s.logger.Info("get_channel tool called", "arguments", request.GetArguments())
		
		// TODO: Implement actual channel retrieval in Step 7
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Channel details will be implemented in Step 7 (MCP Handlers)"),
			},
		}, nil
	}

	return toolDefinition{definition: &tool, handler: handler}
}

// defineSearchChannelsTool creates the search_channels tool definition.
// Searches channels based on name or other criteria.
func (s *Server) defineSearchChannelsTool() toolDefinition {
	tool := mcp.NewTool("search_channels",
		mcp.WithDescription("Search channels by name or other criteria within a specific application. Returns matching channels with relevance scoring."),
		mcp.WithString("app_id", 
			mcp.Required(),
			mcp.Description("The unique identifier of the application"),
		),
		mcp.WithString("query", 
			mcp.Required(),
			mcp.Description("Search query string to match against channel names and descriptions"),
		),
		mcp.WithNumber("limit", 
			mcp.Description("Maximum number of results to return (1-50)"),
			mcp.Min(1),
			mcp.Max(50),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		s.logger.Info("search_channels tool called", "arguments", request.GetArguments())
		
		// TODO: Implement actual channel search in Step 7
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Channel search will be implemented in Step 7 (MCP Handlers)"),
			},
		}, nil
	}

	return toolDefinition{definition: &tool, handler: handler}
}

// Customer Tools

// defineListCustomersTool creates the list_customers tool definition.
// Lists all customers for a specific application.
func (s *Server) defineListCustomersTool() toolDefinition {
	tool := mcp.NewTool("list_customers",
		mcp.WithDescription("List customers for a specific application. Returns customer information including name, status, and channel assignments."),
		mcp.WithString("app_id", 
			mcp.Required(),
			mcp.Description("The unique identifier of the application"),
		),
		mcp.WithNumber("limit", 
			mcp.Description("Maximum number of customers to return (1-100)"),
			mcp.Min(1),
			mcp.Max(100),
		),
		mcp.WithNumber("offset", 
			mcp.Description("Number of customers to skip for pagination"),
			mcp.Min(0),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		s.logger.Info("list_customers tool called", "arguments", request.GetArguments())
		
		// TODO: Implement actual customer listing in Step 7
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Customer listing will be implemented in Step 7 (MCP Handlers)"),
			},
		}, nil
	}

	return toolDefinition{definition: &tool, handler: handler}
}

// defineGetCustomerTool creates the get_customer tool definition.
// Retrieves detailed information about a specific customer.
func (s *Server) defineGetCustomerTool() toolDefinition {
	tool := mcp.NewTool("get_customer",
		mcp.WithDescription("Get detailed information about a specific customer by ID. Returns comprehensive customer data including license details and deployment status."),
		mcp.WithString("app_id", 
			mcp.Required(),
			mcp.Description("The unique identifier of the application"),
		),
		mcp.WithString("customer_id", 
			mcp.Required(),
			mcp.Description("The unique identifier of the customer"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		s.logger.Info("get_customer tool called", "arguments", request.GetArguments())
		
		// TODO: Implement actual customer retrieval in Step 7
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Customer details will be implemented in Step 7 (MCP Handlers)"),
			},
		}, nil
	}

	return toolDefinition{definition: &tool, handler: handler}
}

// defineSearchCustomersTool creates the search_customers tool definition.
// Searches customers based on name or other criteria.
func (s *Server) defineSearchCustomersTool() toolDefinition {
	tool := mcp.NewTool("search_customers",
		mcp.WithDescription("Search customers by name or other criteria within a specific application. Returns matching customers with relevance scoring."),
		mcp.WithString("app_id", 
			mcp.Required(),
			mcp.Description("The unique identifier of the application"),
		),
		mcp.WithString("query", 
			mcp.Required(),
			mcp.Description("Search query string to match against customer names and metadata"),
		),
		mcp.WithNumber("limit", 
			mcp.Description("Maximum number of results to return (1-50)"),
			mcp.Min(1),
			mcp.Max(50),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		s.logger.Info("search_customers tool called", "arguments", request.GetArguments())
		
		// TODO: Implement actual customer search in Step 7
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Customer search will be implemented in Step 7 (MCP Handlers)"),
			},
		}, nil
	}

	return toolDefinition{definition: &tool, handler: handler}
}