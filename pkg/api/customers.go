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

// CustomerService provides methods for interacting with customers
type CustomerService struct {
	client *Client
}

// NewCustomerService creates a new customer service
func NewCustomerService(client *Client) *CustomerService {
	return &CustomerService{client: client}
}

// List retrieves a paginated list of customers for a specific application
func (s *CustomerService) List(ctx context.Context, appID string, opts *ListOptions) (*PaginatedResponse[Customer], error) {
	if appID == "" {
		return nil, fmt.Errorf("application ID is required")
	}

	path := fmt.Sprintf("/v3/app/%s/customers", url.PathEscape(appID))

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
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}
	defer resp.Body.Close()

	if apiErr := s.client.ConvertHTTPError(resp); apiErr != nil {
		return nil, apiErr
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result PaginatedResponse[Customer]
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Get retrieves a specific customer by ID
func (s *CustomerService) Get(ctx context.Context, appID, customerID string) (*Customer, error) {
	if appID == "" {
		return nil, fmt.Errorf("application ID is required")
	}
	if customerID == "" {
		return nil, fmt.Errorf("customer ID is required")
	}

	path := fmt.Sprintf("/v3/app/%s/customer/%s", url.PathEscape(appID), url.PathEscape(customerID))

	resp, err := s.client.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
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
		Customer Customer `json:"customer"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result.Customer, nil
}

// Search searches for customers within an application based on query criteria
func (s *CustomerService) Search(ctx context.Context, appID string, opts *SearchOptions) (*PaginatedResponse[Customer], error) {
	if appID == "" {
		return nil, fmt.Errorf("application ID is required")
	}
	if opts == nil || strings.TrimSpace(opts.Query) == "" {
		return nil, fmt.Errorf("search query is required")
	}

	// Since there's no dedicated search endpoint for customers, use list and filter client-side
	listOpts := &ListOptions{
		Limit:  100, // Get more results for better search coverage
		Offset: opts.Offset,
	}

	customers, err := s.List(ctx, appID, listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers for search: %w", err)
	}

	// Filter customers client-side based on query
	var filteredCustomers []Customer
	queryLower := strings.ToLower(strings.TrimSpace(opts.Query))

	for _, customer := range customers.Data {
		// Search in name, email, type, and status (case-insensitive)
		if strings.Contains(strings.ToLower(customer.Name), queryLower) ||
			strings.Contains(strings.ToLower(customer.Email), queryLower) ||
			strings.Contains(strings.ToLower(customer.Type), queryLower) ||
			strings.Contains(strings.ToLower(customer.Status), queryLower) ||
			strings.Contains(strings.ToLower(customer.LicenseID), queryLower) ||
			strings.Contains(strings.ToLower(customer.ChannelName), queryLower) {
			filteredCustomers = append(filteredCustomers, customer)
		}
	}

	// Apply limit if specified
	if opts.Limit > 0 && len(filteredCustomers) > opts.Limit {
		filteredCustomers = filteredCustomers[:opts.Limit]
	}

	result := &PaginatedResponse[Customer]{
		Data:       filteredCustomers,
		TotalCount: len(filteredCustomers),
		Page:       1, // Since we're doing client-side filtering
		PageSize:   len(filteredCustomers),
		HasMore:    false, // We're not implementing true pagination for search
	}

	return result, nil
}
