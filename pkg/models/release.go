package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Release validation constants
const (
	MaxNotesLength = 10000
)

// Release represents a Replicated application release
type Release struct {
	ID            string            `json:"id"`
	ApplicationID string            `json:"application_id"`
	Version       string            `json:"version"`
	Sequence      int64             `json:"sequence"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	ReleasedAt    *time.Time        `json:"released_at,omitempty"`
	Notes         string            `json:"notes,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	IsRequired    bool              `json:"is_required"`
	IsPrerelease  bool              `json:"is_prerelease"`
	Status        string            `json:"status"`
	Config        string            `json:"config,omitempty"`
}

// Release status constants
const (
	ReleaseStatusDraft      = "draft"
	ReleaseStatusReleased   = "released"
	ReleaseStatusArchived   = "archived"
	ReleaseStatusSuperseded = "superseded"
)

var validReleaseStatuses = []string{
	ReleaseStatusDraft,
	ReleaseStatusReleased,
	ReleaseStatusArchived,
	ReleaseStatusSuperseded,
}

// semVerRegex matches semantic version format (X.Y.Z with optional pre-release and build metadata)
var semVerRegex = regexp.MustCompile(
	`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)` +
		`(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?` +
		`(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`,
)

// Validate ensures the Release struct contains valid data
func (r *Release) Validate() error {
	var errors []string

	errors = append(errors, r.validateBasicFields()...)
	errors = append(errors, r.validateTimestamps()...)
	errors = append(errors, r.validateOptionalFields()...)

	if len(errors) > 0 {
		return fmt.Errorf("release validation errors:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

// validateBasicFields validates basic release fields
func (r *Release) validateBasicFields() []string {
	var errors []string

	if r.ID == "" {
		errors = append(errors, "release ID is required")
	}
	if r.ApplicationID == "" {
		errors = append(errors, "application ID is required")
	}
	if r.Version == "" {
		errors = append(errors, "release version is required")
	} else if !isValidSemanticVersion(r.Version) {
		errors = append(errors, "release version must follow semantic versioning format (e.g., 1.0.0)")
	}
	if r.Sequence < 0 {
		errors = append(errors, "release sequence must be non-negative")
	}
	if r.Status == "" {
		errors = append(errors, "release status is required")
	} else if !isValidReleaseStatus(r.Status) {
		errors = append(errors, fmt.Sprintf("invalid release status '%s'. Valid statuses are: %s",
			r.Status, strings.Join(validReleaseStatuses, ", ")))
	}

	return errors
}

// validateTimestamps validates release timestamp fields
func (r *Release) validateTimestamps() []string {
	var errors []string

	if r.CreatedAt.IsZero() {
		errors = append(errors, "created_at timestamp is required")
	}
	if r.UpdatedAt.IsZero() {
		errors = append(errors, "updated_at timestamp is required")
	}
	if !r.CreatedAt.IsZero() && !r.UpdatedAt.IsZero() && r.UpdatedAt.Before(r.CreatedAt) {
		errors = append(errors, "updated_at must be equal to or after created_at")
	}
	if r.ReleasedAt != nil && r.ReleasedAt.Before(r.CreatedAt) {
		errors = append(errors, "released_at must be equal to or after created_at")
	}
	if r.Status == ReleaseStatusReleased && r.ReleasedAt == nil {
		errors = append(errors, "released_at is required when status is 'released'")
	}

	return errors
}

// validateOptionalFields validates optional release fields
func (r *Release) validateOptionalFields() []string {
	var errors []string

	if r.Notes != "" && len(r.Notes) > MaxNotesLength {
		errors = append(errors, "release notes must be 10000 characters or less")
	}

	errors = append(errors, validateKeyValueMap(r.Metadata, "metadata")...)

	return errors
}

// isValidSemanticVersion checks if the version follows semantic versioning
func isValidSemanticVersion(version string) bool {
	return semVerRegex.MatchString(version)
}

// isValidReleaseStatus checks if the provided status is valid
func isValidReleaseStatus(status string) bool {
	for _, valid := range validReleaseStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

// IsReleased returns true if the release has been released
func (r *Release) IsReleased() bool {
	return r.Status == ReleaseStatusReleased && r.ReleasedAt != nil
}

// String returns a string representation of the Release
func (r *Release) String() string {
	return fmt.Sprintf("Release{ID: %s, ApplicationID: %s, Version: %s, Sequence: %d, Status: %s}",
		r.ID, r.ApplicationID, r.Version, r.Sequence, r.Status)
}
