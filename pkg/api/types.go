package api

import "time"

// Application represents a Replicated application
type Application struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ApplicationListResponse represents the response from listing applications
type ApplicationListResponse struct {
	Applications []Application `json:"apps"`
}

// ApplicationResponse represents the response from getting a single application
type ApplicationResponse struct {
	Application Application `json:"app"`
}

// Release represents a release in the Replicated Vendor Portal
type Release struct {
	Sequence     int       `json:"sequence"`
	Version      string    `json:"version"`
	ReleaseNotes string    `json:"release_notes"`
	Required     bool      `json:"required"`
	CreatedAt    time.Time `json:"created_at"`
	EditedAt     time.Time `json:"edited_at"`
}

// ReleaseListResponse represents the response from listing releases
type ReleaseListResponse struct {
	Releases []Release `json:"releases"`
}

// ReleaseResponse represents the response from getting a single release
type ReleaseResponse struct {
	Release Release `json:"release"`
}

// Channel represents a release channel
type Channel struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ChannelListResponse represents the response from listing channels
type ChannelListResponse struct {
	Channels []Channel `json:"channels"`
}

// ChannelResponse represents the response from getting a single channel
type ChannelResponse struct {
	Channel Channel `json:"channel"`
}

// Customer represents a customer in the Replicated Vendor Portal
type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CustomerListResponse represents the response from listing customers
type CustomerListResponse struct {
	Customers []Customer `json:"customers"`
}

// CustomerResponse represents the response from getting a single customer
type CustomerResponse struct {
	Customer Customer `json:"customer"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
}

