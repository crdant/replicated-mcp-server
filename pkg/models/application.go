// Package models provides data structures for Replicated Vendor Portal entities.
// These models are designed for Phase 1 read-only operations and include
// comprehensive validation and JSON marshaling support.
package models

import (
	"fmt"
	"strings"
	"time"
)

// Validation constants
const (
	MaxNameLength        = 255
	MaxDescriptionLength = 1000
)

// Application represents a Replicated application in the Vendor Portal
type Application struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	TeamID      string    `json:"team_id"`
	TeamName    string    `json:"team_name,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Description string    `json:"description,omitempty"`
	Icon        string    `json:"icon,omitempty"`
	IsActive    bool      `json:"is_active"`
}

// Validate ensures the Application struct contains valid data
func (a *Application) Validate() error {
	var errors []string

	// Validate ID
	if a.ID == "" {
		errors = append(errors, "application ID is required")
	}

	// Validate Name
	if a.Name == "" {
		errors = append(errors, "application name is required")
	} else if len(a.Name) > MaxNameLength {
		errors = append(errors, "application name must be 255 characters or less")
	}

	// Validate Slug
	if a.Slug == "" {
		errors = append(errors, "application slug is required")
	} else if !isValidSlug(a.Slug) {
		errors = append(errors, "application slug must contain only lowercase letters, numbers, and hyphens")
	}

	// Validate TeamID
	if a.TeamID == "" {
		errors = append(errors, "team ID is required")
	}

	// Validate timestamps
	if a.CreatedAt.IsZero() {
		errors = append(errors, "created_at timestamp is required")
	}
	if a.UpdatedAt.IsZero() {
		errors = append(errors, "updated_at timestamp is required")
	}
	if !a.CreatedAt.IsZero() && !a.UpdatedAt.IsZero() && a.UpdatedAt.Before(a.CreatedAt) {
		errors = append(errors, "updated_at must be equal to or after created_at")
	}

	// Validate optional fields
	if a.Description != "" && len(a.Description) > MaxDescriptionLength {
		errors = append(errors, "application description must be 1000 characters or less")
	}

	if len(errors) > 0 {
		return fmt.Errorf("application validation errors:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

// isValidSlug checks if a slug contains only valid characters
func isValidSlug(slug string) bool {
	if slug == "" {
		return false
	}

	for _, r := range slug {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
			return false
		}
	}

	// Slug cannot start or end with hyphen
	return !strings.HasPrefix(slug, "-") && !strings.HasSuffix(slug, "-")
}

// String returns a string representation of the Application
func (a *Application) String() string {
	return fmt.Sprintf("Application{ID: %s, Name: %s, Slug: %s, TeamID: %s, IsActive: %t}",
		a.ID, a.Name, a.Slug, a.TeamID, a.IsActive)
}
