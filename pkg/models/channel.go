package models

import (
	"fmt"
	"strings"
	"time"
)

// Channel validation constants
const (
	MaxChannelNameLength        = 100
	MaxChannelDescriptionLength = 500
)

// Channel represents a Replicated release channel
type Channel struct {
	ID              string     `json:"id"`
	ApplicationID   string     `json:"application_id"`
	Name            string     `json:"name"`
	Description     string     `json:"description,omitempty"`
	ReleaseID       string     `json:"release_id,omitempty"`
	ReleaseSequence int64      `json:"release_sequence,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	ArchivedAt      *time.Time `json:"archived_at,omitempty"`
	IsDefault       bool       `json:"is_default"`
	IsArchived      bool       `json:"is_archived"`
	ChannelSlug     string     `json:"channel_slug"`
}

// Validate ensures the Channel struct contains valid data
func (c *Channel) Validate() error {
	var errors []string

	// Validate ID
	if c.ID == "" {
		errors = append(errors, "channel ID is required")
	}

	// Validate ApplicationID
	if c.ApplicationID == "" {
		errors = append(errors, "application ID is required")
	}

	// Validate Name
	if c.Name == "" {
		errors = append(errors, "channel name is required")
	} else if len(c.Name) > MaxChannelNameLength {
		errors = append(errors, "channel name must be 100 characters or less")
	}

	// Validate ChannelSlug
	if c.ChannelSlug == "" {
		errors = append(errors, "channel slug is required")
	} else if !isValidChannelSlug(c.ChannelSlug) {
		errors = append(errors, "channel slug must contain only lowercase letters, numbers, and hyphens")
	}

	// Validate timestamps
	if c.CreatedAt.IsZero() {
		errors = append(errors, "created_at timestamp is required")
	}
	if c.UpdatedAt.IsZero() {
		errors = append(errors, "updated_at timestamp is required")
	}
	if !c.CreatedAt.IsZero() && !c.UpdatedAt.IsZero() && c.UpdatedAt.Before(c.CreatedAt) {
		errors = append(errors, "updated_at must be equal to or after created_at")
	}

	// Validate ArchivedAt if provided
	if c.ArchivedAt != nil {
		if c.ArchivedAt.Before(c.CreatedAt) {
			errors = append(errors, "archived_at must be equal to or after created_at")
		}
		// If archived_at is set, is_archived should be true
		if !c.IsArchived {
			errors = append(errors, "is_archived must be true when archived_at is set")
		}
	}

	// Validate archived state consistency
	if c.IsArchived && c.ArchivedAt == nil {
		errors = append(errors, "archived_at is required when is_archived is true")
	}

	// Validate release relationship
	if c.ReleaseID != "" && c.ReleaseSequence <= 0 {
		errors = append(errors, "release_sequence must be positive when release_id is provided")
	}
	if c.ReleaseID == "" && c.ReleaseSequence > 0 {
		errors = append(errors, "release_id is required when release_sequence is provided")
	}

	// Validate optional fields
	if c.Description != "" && len(c.Description) > MaxChannelDescriptionLength {
		errors = append(errors, "channel description must be 500 characters or less")
	}

	if len(errors) > 0 {
		return fmt.Errorf("channel validation errors:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

// isValidChannelSlug checks if a channel slug contains only valid characters
func isValidChannelSlug(slug string) bool {
	if slug == "" {
		return false
	}

	for _, r := range slug {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
			return false
		}
	}

	// Channel slug cannot start or end with hyphen
	return !strings.HasPrefix(slug, "-") && !strings.HasSuffix(slug, "-")
}

// HasRelease returns true if the channel has an assigned release
func (c *Channel) HasRelease() bool {
	return c.ReleaseID != "" && c.ReleaseSequence > 0
}

// IsActive returns true if the channel is not archived
func (c *Channel) IsActive() bool {
	return !c.IsArchived
}

// String returns a string representation of the Channel
func (c *Channel) String() string {
	return fmt.Sprintf("Channel{ID: %s, ApplicationID: %s, Name: %s, Slug: %s, IsDefault: %t, IsArchived: %t}",
		c.ID, c.ApplicationID, c.Name, c.ChannelSlug, c.IsDefault, c.IsArchived)
}
