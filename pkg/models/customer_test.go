package models

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestCustomer_Validate(t *testing.T) {
	validTime := time.Now()
	laterTime := validTime.Add(time.Hour)
	futureTime := validTime.Add(24 * time.Hour)

	tests := []struct {
		name        string
		customer    Customer
		wantErr     bool
		errContains []string
	}{
		{
			name: "valid customer",
			customer: Customer{
				ID:                "cust-123",
				ApplicationID:     "app-456",
				Name:              "Test Customer",
				Email:             "test@example.com",
				ChannelID:         "ch-789",
				ChannelName:       "Stable",
				CreatedAt:         validTime,
				UpdatedAt:         laterTime,
				ExpiresAt:         &futureTime,
				Type:              CustomerTypePaid,
				IsArchived:        false,
				IsGitOpsSupported: true,
				LicenseID:         "lic-abc",
				LicenseType:       LicenseTypePaid,
			},
			wantErr: false,
		},
		{
			name: "minimal valid customer",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          "Test Customer",
				ChannelID:     "ch-789",
				CreatedAt:     validTime,
				UpdatedAt:     validTime,
				Type:          CustomerTypeTrial,
				LicenseID:     "lic-abc",
				LicenseType:   LicenseTypeTrial,
			},
			wantErr: false,
		},
		{
			name: "valid archived customer",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          "Test Customer",
				ChannelID:     "ch-789",
				CreatedAt:     validTime,
				UpdatedAt:     validTime,
				ArchivedAt:    &validTime,
				Type:          CustomerTypeTrial,
				IsArchived:    true,
				LicenseID:     "lic-abc",
				LicenseType:   LicenseTypeTrial,
			},
			wantErr: false,
		},
		{
			name:        "missing ID",
			customer:    Customer{},
			wantErr:     true,
			errContains: []string{"customer ID is required"},
		},
		{
			name: "missing application ID",
			customer: Customer{
				ID: "cust-123",
			},
			wantErr:     true,
			errContains: []string{"application ID is required"},
		},
		{
			name: "missing name",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
			},
			wantErr:     true,
			errContains: []string{"customer name is required"},
		},
		{
			name: "name too long",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          strings.Repeat("a", 256),
			},
			wantErr:     true,
			errContains: []string{"customer name must be 255 characters or less"},
		},
		{
			name: "invalid email",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          "Test Customer",
				Email:         "invalid-email",
			},
			wantErr:     true,
			errContains: []string{"customer email must be a valid email address"},
		},
		{
			name: "missing channel ID",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          "Test Customer",
			},
			wantErr:     true,
			errContains: []string{"channel ID is required"},
		},
		{
			name: "missing type",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          "Test Customer",
				ChannelID:     "ch-789",
			},
			wantErr:     true,
			errContains: []string{"customer type is required"},
		},
		{
			name: "invalid type",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          "Test Customer",
				ChannelID:     "ch-789",
				Type:          "invalid",
			},
			wantErr:     true,
			errContains: []string{"invalid customer type 'invalid'"},
		},
		{
			name: "missing license ID",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          "Test Customer",
				ChannelID:     "ch-789",
				Type:          CustomerTypePaid,
			},
			wantErr:     true,
			errContains: []string{"license ID is required"},
		},
		{
			name: "missing license type",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          "Test Customer",
				ChannelID:     "ch-789",
				Type:          CustomerTypePaid,
				LicenseID:     "lic-abc",
			},
			wantErr:     true,
			errContains: []string{"license type is required"},
		},
		{
			name: "invalid license type",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          "Test Customer",
				ChannelID:     "ch-789",
				Type:          CustomerTypePaid,
				LicenseID:     "lic-abc",
				LicenseType:   "invalid",
			},
			wantErr:     true,
			errContains: []string{"invalid license type 'invalid'"},
		},
		{
			name: "missing timestamps",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          "Test Customer",
				ChannelID:     "ch-789",
				Type:          CustomerTypePaid,
				LicenseID:     "lic-abc",
				LicenseType:   LicenseTypePaid,
			},
			wantErr:     true,
			errContains: []string{"created_at timestamp is required", "updated_at timestamp is required"},
		},
		{
			name: "updated_at before created_at",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          "Test Customer",
				ChannelID:     "ch-789",
				CreatedAt:     laterTime,
				UpdatedAt:     validTime,
				Type:          CustomerTypePaid,
				LicenseID:     "lic-abc",
				LicenseType:   LicenseTypePaid,
			},
			wantErr:     true,
			errContains: []string{"updated_at must be equal to or after created_at"},
		},
		{
			name: "archived_at before created_at",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          "Test Customer",
				ChannelID:     "ch-789",
				CreatedAt:     laterTime,
				UpdatedAt:     laterTime,
				ArchivedAt:    &validTime,
				Type:          CustomerTypePaid,
				LicenseID:     "lic-abc",
				LicenseType:   LicenseTypePaid,
			},
			wantErr:     true,
			errContains: []string{"archived_at must be equal to or after created_at"},
		},
		{
			name: "archived_at set but is_archived false",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          "Test Customer",
				ChannelID:     "ch-789",
				CreatedAt:     validTime,
				UpdatedAt:     validTime,
				ArchivedAt:    &validTime,
				IsArchived:    false,
				Type:          CustomerTypePaid,
				LicenseID:     "lic-abc",
				LicenseType:   LicenseTypePaid,
			},
			wantErr:     true,
			errContains: []string{"is_archived must be true when archived_at is set"},
		},
		{
			name: "is_archived true but archived_at not set",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          "Test Customer",
				ChannelID:     "ch-789",
				CreatedAt:     validTime,
				UpdatedAt:     validTime,
				IsArchived:    true,
				Type:          CustomerTypePaid,
				LicenseID:     "lic-abc",
				LicenseType:   LicenseTypePaid,
			},
			wantErr:     true,
			errContains: []string{"archived_at is required when is_archived is true"},
		},
		{
			name: "expires_at before created_at",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          "Test Customer",
				ChannelID:     "ch-789",
				CreatedAt:     laterTime,
				UpdatedAt:     laterTime,
				ExpiresAt:     &validTime,
				Type:          CustomerTypeTrial,
				LicenseID:     "lic-abc",
				LicenseType:   LicenseTypeTrial,
			},
			wantErr:     true,
			errContains: []string{"expires_at must be equal to or after created_at"},
		},
		{
			name: "entitlements validation",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          "Test Customer",
				ChannelID:     "ch-789",
				CreatedAt:     validTime,
				UpdatedAt:     validTime,
				Type:          CustomerTypePaid,
				LicenseID:     "lic-abc",
				LicenseType:   LicenseTypePaid,
				Entitlements: map[string]string{
					"":                       "value",
					"key":                    strings.Repeat("a", 501),
					strings.Repeat("k", 101): "value",
				},
			},
			wantErr:     true,
			errContains: []string{"entitlement keys cannot be empty", "entitlement values must be 500 characters or less", "entitlement keys must be 100 characters or less"},
		},
		{
			name: "custom fields validation",
			customer: Customer{
				ID:            "cust-123",
				ApplicationID: "app-456",
				Name:          "Test Customer",
				ChannelID:     "ch-789",
				CreatedAt:     validTime,
				UpdatedAt:     validTime,
				Type:          CustomerTypePaid,
				LicenseID:     "lic-abc",
				LicenseType:   LicenseTypePaid,
				CustomFields: map[string]string{
					"":                       "value",
					"field":                  strings.Repeat("a", 501),
					strings.Repeat("f", 101): "value",
				},
			},
			wantErr:     true,
			errContains: []string{"custom field keys cannot be empty", "custom field values must be 500 characters or less", "custom field keys must be 100 characters or less"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.customer.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Customer.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				for _, expectedErr := range tt.errContains {
					if !strings.Contains(err.Error(), expectedErr) {
						t.Errorf("Customer.Validate() error = %v, should contain %v", err, expectedErr)
					}
				}
			}
		})
	}
}

func TestCustomer_JSONMarshaling(t *testing.T) {
	validTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	futureTime := validTime.Add(24 * time.Hour)

	customer := Customer{
		ID:                "cust-123",
		ApplicationID:     "app-456",
		Name:              "Test Customer",
		Email:             "test@example.com",
		ChannelID:         "ch-789",
		ChannelName:       "Stable",
		CreatedAt:         validTime,
		UpdatedAt:         validTime,
		ExpiresAt:         &futureTime,
		Type:              CustomerTypePaid,
		IsArchived:        false,
		IsGitOpsSupported: true,
		LicenseID:         "lic-abc",
		LicenseType:       LicenseTypePaid,
		Entitlements: map[string]string{
			"feature_a": "enabled",
			"max_users": "100",
		},
		CustomFields: map[string]string{
			"department": "engineering",
			"region":     "us-west",
		},
	}

	// Test marshaling
	jsonData, err := json.Marshal(customer)
	if err != nil {
		t.Fatalf("Failed to marshal Customer: %v", err)
	}

	// Test unmarshaling
	var unmarshaledCustomer Customer
	err = json.Unmarshal(jsonData, &unmarshaledCustomer)
	if err != nil {
		t.Fatalf("Failed to unmarshal Customer: %v", err)
	}

	// Verify fields
	if unmarshaledCustomer.ID != customer.ID {
		t.Errorf("ID mismatch: got %v, want %v", unmarshaledCustomer.ID, customer.ID)
	}
	if unmarshaledCustomer.Name != customer.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaledCustomer.Name, customer.Name)
	}
	if unmarshaledCustomer.Type != customer.Type {
		t.Errorf("Type mismatch: got %v, want %v", unmarshaledCustomer.Type, customer.Type)
	}
	if unmarshaledCustomer.LicenseType != customer.LicenseType {
		t.Errorf("LicenseType mismatch: got %v, want %v", unmarshaledCustomer.LicenseType, customer.LicenseType)
	}

	// Verify entitlements and custom fields
	if len(unmarshaledCustomer.Entitlements) != len(customer.Entitlements) {
		t.Errorf("Entitlements length mismatch: got %v, want %v", len(unmarshaledCustomer.Entitlements), len(customer.Entitlements))
	}
	if len(unmarshaledCustomer.CustomFields) != len(customer.CustomFields) {
		t.Errorf("CustomFields length mismatch: got %v, want %v", len(unmarshaledCustomer.CustomFields), len(customer.CustomFields))
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"valid simple email", "test@example.com", true},
		{"valid email with subdomain", "user@mail.example.com", true},
		{"valid email with numbers", "user123@example.org", true},
		{"valid email with hyphens", "test-user@example-domain.com", true},
		{"invalid no @ symbol", "testexample.com", false},
		{"invalid multiple @ symbols", "test@@example.com", false},
		{"invalid no domain", "test@", false},
		{"invalid no user", "@example.com", false},
		{"invalid no dot in domain", "test@example", false},
		{"empty email", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidEmail(tt.email); got != tt.want {
				t.Errorf("isValidEmail(%v) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestIsValidCustomerType(t *testing.T) {
	tests := []struct {
		name  string
		cType string
		want  bool
	}{
		{"valid trial", CustomerTypeTrial, true},
		{"valid paid", CustomerTypePaid, true},
		{"valid community", CustomerTypeCommunity, true},
		{"valid development", CustomerTypeDevelopment, true},
		{"invalid type", "invalid", false},
		{"empty type", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidCustomerType(tt.cType); got != tt.want {
				t.Errorf("isValidCustomerType(%v) = %v, want %v", tt.cType, got, tt.want)
			}
		})
	}
}

func TestIsValidLicenseType(t *testing.T) {
	tests := []struct {
		name  string
		lType string
		want  bool
	}{
		{"valid trial", LicenseTypeTrial, true},
		{"valid paid", LicenseTypePaid, true},
		{"valid community", LicenseTypeCommunity, true},
		{"valid development", LicenseTypeDevelopment, true},
		{"valid embedded", LicenseTypeEmbedded, true},
		{"invalid type", "invalid", false},
		{"empty type", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidLicenseType(tt.lType); got != tt.want {
				t.Errorf("isValidLicenseType(%v) = %v, want %v", tt.lType, got, tt.want)
			}
		})
	}
}

func TestCustomer_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		customer Customer
		want     bool
	}{
		{
			name: "active customer",
			customer: Customer{
				IsArchived: false,
			},
			want: true,
		},
		{
			name: "archived customer",
			customer: Customer{
				IsArchived: true,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.customer.IsActive(); got != tt.want {
				t.Errorf("Customer.IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomer_IsExpired(t *testing.T) {
	pastTime := time.Now().Add(-time.Hour)
	futureTime := time.Now().Add(time.Hour)

	tests := []struct {
		name     string
		customer Customer
		want     bool
	}{
		{
			name: "expired customer",
			customer: Customer{
				ExpiresAt: &pastTime,
			},
			want: true,
		},
		{
			name: "non-expired customer",
			customer: Customer{
				ExpiresAt: &futureTime,
			},
			want: false,
		},
		{
			name: "customer without expiration",
			customer: Customer{
				ExpiresAt: nil,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.customer.IsExpired(); got != tt.want {
				t.Errorf("Customer.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomer_IsTrialCustomer(t *testing.T) {
	tests := []struct {
		name     string
		customer Customer
		want     bool
	}{
		{
			name: "trial customer type",
			customer: Customer{
				Type:        CustomerTypeTrial,
				LicenseType: LicenseTypePaid,
			},
			want: true,
		},
		{
			name: "trial license type",
			customer: Customer{
				Type:        CustomerTypePaid,
				LicenseType: LicenseTypeTrial,
			},
			want: true,
		},
		{
			name: "paid customer",
			customer: Customer{
				Type:        CustomerTypePaid,
				LicenseType: LicenseTypePaid,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.customer.IsTrialCustomer(); got != tt.want {
				t.Errorf("Customer.IsTrialCustomer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomer_String(t *testing.T) {
	customer := Customer{
		ID:            "cust-123",
		ApplicationID: "app-456",
		Name:          "Test Customer",
		Type:          CustomerTypePaid,
		LicenseType:   LicenseTypePaid,
		IsArchived:    false,
	}

	str := customer.String()
	expected := "Customer{ID: cust-123, ApplicationID: app-456, Name: Test Customer, Type: paid, LicenseType: paid, IsArchived: false}"

	if str != expected {
		t.Errorf("Customer.String() = %v, want %v", str, expected)
	}
}
