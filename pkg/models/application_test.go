package models

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestApplication_Validate(t *testing.T) {
	validTime := time.Now()
	laterTime := validTime.Add(time.Hour)

	tests := []struct {
		name        string
		app         Application
		wantErr     bool
		errContains []string
	}{
		{
			name: "valid application",
			app: Application{
				ID:        "app-123",
				Name:      "Test Application",
				Slug:      "test-app",
				TeamID:    "team-456",
				TeamName:  "Test Team",
				CreatedAt: validTime,
				UpdatedAt: laterTime,
				IsActive:  true,
			},
			wantErr: false,
		},
		{
			name: "minimal valid application",
			app: Application{
				ID:        "app-123",
				Name:      "Test App",
				Slug:      "test-app",
				TeamID:    "team-456",
				CreatedAt: validTime,
				UpdatedAt: validTime,
				IsActive:  false,
			},
			wantErr: false,
		},
		{
			name:        "missing ID",
			app:         Application{},
			wantErr:     true,
			errContains: []string{"application ID is required"},
		},
		{
			name: "missing name",
			app: Application{
				ID:     "app-123",
				Slug:   "test-app",
				TeamID: "team-456",
			},
			wantErr:     true,
			errContains: []string{"application name is required"},
		},
		{
			name: "name too long",
			app: Application{
				ID:     "app-123",
				Name:   strings.Repeat("a", 256),
				Slug:   "test-app",
				TeamID: "team-456",
			},
			wantErr:     true,
			errContains: []string{"application name must be 255 characters or less"},
		},
		{
			name: "missing slug",
			app: Application{
				ID:     "app-123",
				Name:   "Test App",
				TeamID: "team-456",
			},
			wantErr:     true,
			errContains: []string{"application slug is required"},
		},
		{
			name: "invalid slug with uppercase",
			app: Application{
				ID:     "app-123",
				Name:   "Test App",
				Slug:   "Test-App",
				TeamID: "team-456",
			},
			wantErr:     true,
			errContains: []string{"application slug must contain only lowercase letters, numbers, and hyphens"},
		},
		{
			name: "invalid slug with underscore",
			app: Application{
				ID:     "app-123",
				Name:   "Test App",
				Slug:   "test_app",
				TeamID: "team-456",
			},
			wantErr:     true,
			errContains: []string{"application slug must contain only lowercase letters, numbers, and hyphens"},
		},
		{
			name: "invalid slug starting with hyphen",
			app: Application{
				ID:     "app-123",
				Name:   "Test App",
				Slug:   "-test-app",
				TeamID: "team-456",
			},
			wantErr:     true,
			errContains: []string{"application slug must contain only lowercase letters, numbers, and hyphens"},
		},
		{
			name: "invalid slug ending with hyphen",
			app: Application{
				ID:     "app-123",
				Name:   "Test App",
				Slug:   "test-app-",
				TeamID: "team-456",
			},
			wantErr:     true,
			errContains: []string{"application slug must contain only lowercase letters, numbers, and hyphens"},
		},
		{
			name: "missing team ID",
			app: Application{
				ID:   "app-123",
				Name: "Test App",
				Slug: "test-app",
			},
			wantErr:     true,
			errContains: []string{"team ID is required"},
		},
		{
			name: "missing created_at",
			app: Application{
				ID:        "app-123",
				Name:      "Test App",
				Slug:      "test-app",
				TeamID:    "team-456",
				UpdatedAt: validTime,
			},
			wantErr:     true,
			errContains: []string{"created_at timestamp is required"},
		},
		{
			name: "missing updated_at",
			app: Application{
				ID:        "app-123",
				Name:      "Test App",
				Slug:      "test-app",
				TeamID:    "team-456",
				CreatedAt: validTime,
			},
			wantErr:     true,
			errContains: []string{"updated_at timestamp is required"},
		},
		{
			name: "updated_at before created_at",
			app: Application{
				ID:        "app-123",
				Name:      "Test App",
				Slug:      "test-app",
				TeamID:    "team-456",
				CreatedAt: laterTime,
				UpdatedAt: validTime,
			},
			wantErr:     true,
			errContains: []string{"updated_at must be equal to or after created_at"},
		},
		{
			name: "description too long",
			app: Application{
				ID:          "app-123",
				Name:        "Test App",
				Slug:        "test-app",
				TeamID:      "team-456",
				CreatedAt:   validTime,
				UpdatedAt:   validTime,
				Description: strings.Repeat("a", 1001),
			},
			wantErr:     true,
			errContains: []string{"application description must be 1000 characters or less"},
		},
		{
			name: "multiple validation errors",
			app: Application{
				Name:        strings.Repeat("a", 256),
				Slug:        "Test_App",
				Description: strings.Repeat("a", 1001),
			},
			wantErr:     true,
			errContains: []string{"application ID is required", "application name must be 255 characters or less", "application slug must contain only lowercase letters, numbers, and hyphens", "team ID is required"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.app.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Application.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				for _, expectedErr := range tt.errContains {
					if !strings.Contains(err.Error(), expectedErr) {
						t.Errorf("Application.Validate() error = %v, should contain %v", err, expectedErr)
					}
				}
			}
		})
	}
}

func TestApplication_JSONMarshaling(t *testing.T) {
	validTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	app := Application{
		ID:          "app-123",
		Name:        "Test Application",
		Slug:        "test-app",
		TeamID:      "team-456",
		TeamName:    "Test Team",
		CreatedAt:   validTime,
		UpdatedAt:   validTime,
		Description: "A test application",
		Icon:        "https://example.com/icon.png",
		IsActive:    true,
	}

	// Test marshaling
	jsonData, err := json.Marshal(app)
	if err != nil {
		t.Fatalf("Failed to marshal Application: %v", err)
	}

	// Test unmarshaling
	var unmarshaledApp Application
	err = json.Unmarshal(jsonData, &unmarshaledApp)
	if err != nil {
		t.Fatalf("Failed to unmarshal Application: %v", err)
	}

	// Verify fields
	if unmarshaledApp.ID != app.ID {
		t.Errorf("ID mismatch: got %v, want %v", unmarshaledApp.ID, app.ID)
	}
	if unmarshaledApp.Name != app.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaledApp.Name, app.Name)
	}
	if unmarshaledApp.Slug != app.Slug {
		t.Errorf("Slug mismatch: got %v, want %v", unmarshaledApp.Slug, app.Slug)
	}
	if unmarshaledApp.TeamID != app.TeamID {
		t.Errorf("TeamID mismatch: got %v, want %v", unmarshaledApp.TeamID, app.TeamID)
	}
	if unmarshaledApp.IsActive != app.IsActive {
		t.Errorf("IsActive mismatch: got %v, want %v", unmarshaledApp.IsActive, app.IsActive)
	}

	// Verify JSON contains expected fields
	expectedFields := []string{
		`"id":"app-123"`,
		`"name":"Test Application"`,
		`"slug":"test-app"`,
		`"team_id":"team-456"`,
		`"is_active":true`,
	}

	jsonString := string(jsonData)
	for _, field := range expectedFields {
		if !strings.Contains(jsonString, field) {
			t.Errorf("JSON should contain %v, got %v", field, jsonString)
		}
	}
}

func TestIsValidSlug(t *testing.T) {
	tests := []struct {
		name string
		slug string
		want bool
	}{
		{"valid simple slug", "test", true},
		{"valid slug with numbers", "test123", true},
		{"valid slug with hyphens", "test-app-123", true},
		{"empty slug", "", false},
		{"slug with uppercase", "Test", false},
		{"slug with underscore", "test_app", false},
		{"slug with spaces", "test app", false},
		{"slug starting with hyphen", "-test", false},
		{"slug ending with hyphen", "test-", false},
		{"slug with special characters", "test@app", false},
		{"valid complex slug", "my-app-v2", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidSlug(tt.slug); got != tt.want {
				t.Errorf("isValidSlug(%v) = %v, want %v", tt.slug, got, tt.want)
			}
		})
	}
}

func TestApplication_String(t *testing.T) {
	app := Application{
		ID:       "app-123",
		Name:     "Test App",
		Slug:     "test-app",
		TeamID:   "team-456",
		IsActive: true,
	}

	str := app.String()
	expected := "Application{ID: app-123, Name: Test App, Slug: test-app, TeamID: team-456, IsActive: true}"

	if str != expected {
		t.Errorf("Application.String() = %v, want %v", str, expected)
	}
}
