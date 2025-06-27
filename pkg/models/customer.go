package models

import (
	"fmt"
	"strings"
	"time"
)

// Customer validation constants
const (
	MaxCustomerNameLength = 255
	EmailParts            = 2
)

// Customer represents a Replicated customer with license and installation details
type Customer struct {
	ID                string            `json:"id"`
	ApplicationID     string            `json:"application_id"`
	Name              string            `json:"name"`
	Email             string            `json:"email,omitempty"`
	ChannelID         string            `json:"channel_id"`
	ChannelName       string            `json:"channel_name,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	ArchivedAt        *time.Time        `json:"archived_at,omitempty"`
	ExpiresAt         *time.Time        `json:"expires_at,omitempty"`
	Type              string            `json:"type"`
	IsArchived        bool              `json:"is_archived"`
	IsGitOpsSupported bool              `json:"is_gitops_supported"`
	LicenseID         string            `json:"license_id"`
	LicenseType       string            `json:"license_type"`
	Entitlements      map[string]string `json:"entitlements,omitempty"`
	CustomFields      map[string]string `json:"custom_fields,omitempty"`
}

// Customer type constants
const (
	CustomerTypeTrial       = "trial"
	CustomerTypePaid        = "paid"
	CustomerTypeCommunity   = "community"
	CustomerTypeDevelopment = "development"
)

// License type constants
const (
	LicenseTypeTrial       = "trial"
	LicenseTypePaid        = "paid"
	LicenseTypeCommunity   = "community"
	LicenseTypeDevelopment = "development"
	LicenseTypeEmbedded    = "embedded"
)

var validCustomerTypes = []string{
	CustomerTypeTrial,
	CustomerTypePaid,
	CustomerTypeCommunity,
	CustomerTypeDevelopment,
}

var validLicenseTypes = []string{
	LicenseTypeTrial,
	LicenseTypePaid,
	LicenseTypeCommunity,
	LicenseTypeDevelopment,
	LicenseTypeEmbedded,
}

// Validate ensures the Customer struct contains valid data
func (c *Customer) Validate() error {
	var errors []string

	errors = append(errors, c.validateBasicFields()...)
	errors = append(errors, c.validateTimestamps()...)
	errors = append(errors, c.validateKeyValueMaps()...)

	if len(errors) > 0 {
		return fmt.Errorf("customer validation errors:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

// validateBasicFields validates basic customer fields
func (c *Customer) validateBasicFields() []string {
	var errors []string

	if c.ID == "" {
		errors = append(errors, "customer ID is required")
	}
	if c.ApplicationID == "" {
		errors = append(errors, "application ID is required")
	}
	if c.Name == "" {
		errors = append(errors, "customer name is required")
	} else if len(c.Name) > MaxCustomerNameLength {
		errors = append(errors, "customer name must be 255 characters or less")
	}
	if c.Email != "" && !isValidEmail(c.Email) {
		errors = append(errors, "customer email must be a valid email address")
	}
	if c.ChannelID == "" {
		errors = append(errors, "channel ID is required")
	}
	if c.Type == "" {
		errors = append(errors, "customer type is required")
	} else if !isValidCustomerType(c.Type) {
		errors = append(errors, fmt.Sprintf("invalid customer type '%s'. Valid types are: %s",
			c.Type, strings.Join(validCustomerTypes, ", ")))
	}
	if c.LicenseID == "" {
		errors = append(errors, "license ID is required")
	}
	if c.LicenseType == "" {
		errors = append(errors, "license type is required")
	} else if !isValidLicenseType(c.LicenseType) {
		errors = append(errors, fmt.Sprintf("invalid license type '%s'. Valid types are: %s",
			c.LicenseType, strings.Join(validLicenseTypes, ", ")))
	}

	return errors
}

// validateTimestamps validates customer timestamp fields
func (c *Customer) validateTimestamps() []string {
	var errors []string

	if c.CreatedAt.IsZero() {
		errors = append(errors, "created_at timestamp is required")
	}
	if c.UpdatedAt.IsZero() {
		errors = append(errors, "updated_at timestamp is required")
	}
	if !c.CreatedAt.IsZero() && !c.UpdatedAt.IsZero() && c.UpdatedAt.Before(c.CreatedAt) {
		errors = append(errors, "updated_at must be equal to or after created_at")
	}
	if c.ArchivedAt != nil {
		if c.ArchivedAt.Before(c.CreatedAt) {
			errors = append(errors, "archived_at must be equal to or after created_at")
		}
		if !c.IsArchived {
			errors = append(errors, "is_archived must be true when archived_at is set")
		}
	}
	if c.IsArchived && c.ArchivedAt == nil {
		errors = append(errors, "archived_at is required when is_archived is true")
	}
	if c.ExpiresAt != nil && c.ExpiresAt.Before(c.CreatedAt) {
		errors = append(errors, "expires_at must be equal to or after created_at")
	}

	return errors
}

// validateKeyValueMaps validates entitlements and custom fields
func (c *Customer) validateKeyValueMaps() []string {
	var errors []string

	errors = append(errors, validateKeyValueMap(c.Entitlements, "entitlement")...)
	errors = append(errors, validateKeyValueMap(c.CustomFields, "custom field")...)

	return errors
}

// isValidCustomerType checks if the provided customer type is valid
func isValidCustomerType(customerType string) bool {
	for _, valid := range validCustomerTypes {
		if customerType == valid {
			return true
		}
	}
	return false
}

// isValidLicenseType checks if the provided license type is valid
func isValidLicenseType(licenseType string) bool {
	for _, valid := range validLicenseTypes {
		if licenseType == valid {
			return true
		}
	}
	return false
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	// Basic email validation - contains @ and has characters before and after
	parts := strings.Split(email, "@")
	if len(parts) != EmailParts {
		return false
	}
	if parts[0] == "" || parts[1] == "" {
		return false
	}
	// Domain part should contain at least one dot
	return strings.Contains(parts[1], ".")
}

// IsActive returns true if the customer is not archived
func (c *Customer) IsActive() bool {
	return !c.IsArchived
}

// IsExpired returns true if the customer license has expired
func (c *Customer) IsExpired() bool {
	return c.ExpiresAt != nil && c.ExpiresAt.Before(time.Now())
}

// IsTrialCustomer returns true if the customer is a trial customer
func (c *Customer) IsTrialCustomer() bool {
	return c.Type == CustomerTypeTrial || c.LicenseType == LicenseTypeTrial
}

// String returns a string representation of the Customer
func (c *Customer) String() string {
	return fmt.Sprintf("Customer{ID: %s, ApplicationID: %s, Name: %s, Type: %s, LicenseType: %s, IsArchived: %t}",
		c.ID, c.ApplicationID, c.Name, c.Type, c.LicenseType, c.IsArchived)
}
