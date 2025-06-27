package models

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestRelease_Validate(t *testing.T) {
	validTime := time.Now()
	laterTime := validTime.Add(time.Hour)

	tests := []struct {
		name        string
		release     Release
		wantErr     bool
		errContains []string
	}{
		{
			name: "valid release",
			release: Release{
				ID:            "rel-123",
				ApplicationID: "app-456",
				Version:       "1.0.0",
				Sequence:      1,
				CreatedAt:     validTime,
				UpdatedAt:     laterTime,
				ReleasedAt:    &laterTime,
				Notes:         "Initial release",
				Status:        ReleaseStatusReleased,
				IsRequired:    false,
				IsPrerelease:  false,
			},
			wantErr: false,
		},
		{
			name: "valid prerelease version",
			release: Release{
				ID:            "rel-123",
				ApplicationID: "app-456",
				Version:       "1.0.0-beta.1",
				Sequence:      1,
				CreatedAt:     validTime,
				UpdatedAt:     validTime,
				Status:        ReleaseStatusDraft,
				IsPrerelease:  true,
			},
			wantErr: false,
		},
		{
			name: "valid version with build metadata",
			release: Release{
				ID:            "rel-123",
				ApplicationID: "app-456",
				Version:       "1.0.0+20130313144700",
				Sequence:      1,
				CreatedAt:     validTime,
				UpdatedAt:     validTime,
				Status:        ReleaseStatusDraft,
			},
			wantErr: false,
		},
		{
			name:        "missing ID",
			release:     Release{},
			wantErr:     true,
			errContains: []string{"release ID is required"},
		},
		{
			name: "missing application ID",
			release: Release{
				ID: "rel-123",
			},
			wantErr:     true,
			errContains: []string{"application ID is required"},
		},
		{
			name: "missing version",
			release: Release{
				ID:            "rel-123",
				ApplicationID: "app-456",
			},
			wantErr:     true,
			errContains: []string{"release version is required"},
		},
		{
			name: "invalid semantic version",
			release: Release{
				ID:            "rel-123",
				ApplicationID: "app-456",
				Version:       "1.0",
			},
			wantErr:     true,
			errContains: []string{"release version must follow semantic versioning format"},
		},
		{
			name: "invalid version format",
			release: Release{
				ID:            "rel-123",
				ApplicationID: "app-456",
				Version:       "v1.0.0",
			},
			wantErr:     true,
			errContains: []string{"release version must follow semantic versioning format"},
		},
		{
			name: "negative sequence",
			release: Release{
				ID:            "rel-123",
				ApplicationID: "app-456",
				Version:       "1.0.0",
				Sequence:      -1,
			},
			wantErr:     true,
			errContains: []string{"release sequence must be non-negative"},
		},
		{
			name: "missing status",
			release: Release{
				ID:            "rel-123",
				ApplicationID: "app-456",
				Version:       "1.0.0",
				Sequence:      1,
			},
			wantErr:     true,
			errContains: []string{"release status is required"},
		},
		{
			name: "invalid status",
			release: Release{
				ID:            "rel-123",
				ApplicationID: "app-456",
				Version:       "1.0.0",
				Sequence:      1,
				Status:        "invalid",
			},
			wantErr:     true,
			errContains: []string{"invalid release status 'invalid'"},
		},
		{
			name: "missing timestamps",
			release: Release{
				ID:            "rel-123",
				ApplicationID: "app-456",
				Version:       "1.0.0",
				Sequence:      1,
				Status:        ReleaseStatusDraft,
			},
			wantErr:     true,
			errContains: []string{"created_at timestamp is required", "updated_at timestamp is required"},
		},
		{
			name: "updated_at before created_at",
			release: Release{
				ID:            "rel-123",
				ApplicationID: "app-456",
				Version:       "1.0.0",
				Sequence:      1,
				CreatedAt:     laterTime,
				UpdatedAt:     validTime,
				Status:        ReleaseStatusDraft,
			},
			wantErr:     true,
			errContains: []string{"updated_at must be equal to or after created_at"},
		},
		{
			name: "released_at before created_at",
			release: Release{
				ID:            "rel-123",
				ApplicationID: "app-456",
				Version:       "1.0.0",
				Sequence:      1,
				CreatedAt:     laterTime,
				UpdatedAt:     laterTime,
				ReleasedAt:    &validTime,
				Status:        ReleaseStatusReleased,
			},
			wantErr:     true,
			errContains: []string{"released_at must be equal to or after created_at"},
		},
		{
			name: "released status without released_at",
			release: Release{
				ID:            "rel-123",
				ApplicationID: "app-456",
				Version:       "1.0.0",
				Sequence:      1,
				CreatedAt:     validTime,
				UpdatedAt:     validTime,
				Status:        ReleaseStatusReleased,
			},
			wantErr:     true,
			errContains: []string{"released_at is required when status is 'released'"},
		},
		{
			name: "notes too long",
			release: Release{
				ID:            "rel-123",
				ApplicationID: "app-456",
				Version:       "1.0.0",
				Sequence:      1,
				CreatedAt:     validTime,
				UpdatedAt:     validTime,
				Status:        ReleaseStatusDraft,
				Notes:         strings.Repeat("a", 10001),
			},
			wantErr:     true,
			errContains: []string{"release notes must be 10000 characters or less"},
		},
		{
			name: "metadata validation",
			release: Release{
				ID:            "rel-123",
				ApplicationID: "app-456",
				Version:       "1.0.0",
				Sequence:      1,
				CreatedAt:     validTime,
				UpdatedAt:     validTime,
				Status:        ReleaseStatusDraft,
				Metadata: map[string]string{
					"":                       "value",
					"key":                    strings.Repeat("a", 501),
					strings.Repeat("k", 101): "value",
				},
			},
			wantErr:     true,
			errContains: []string{"metadata keys cannot be empty", "metadata values must be 500 characters or less", "metadata keys must be 100 characters or less"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.release.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Release.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				for _, expectedErr := range tt.errContains {
					if !strings.Contains(err.Error(), expectedErr) {
						t.Errorf("Release.Validate() error = %v, should contain %v", err, expectedErr)
					}
				}
			}
		})
	}
}

func TestRelease_JSONMarshaling(t *testing.T) {
	validTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	releasedTime := validTime.Add(time.Hour)

	release := Release{
		ID:            "rel-123",
		ApplicationID: "app-456",
		Version:       "1.0.0",
		Sequence:      1,
		CreatedAt:     validTime,
		UpdatedAt:     releasedTime,
		ReleasedAt:    &releasedTime,
		Notes:         "Initial release",
		Status:        ReleaseStatusReleased,
		IsRequired:    true,
		IsPrerelease:  false,
		Metadata: map[string]string{
			"git_commit": "abc123",
			"builder":    "github-actions",
		},
	}

	// Test marshaling
	jsonData, err := json.Marshal(release)
	if err != nil {
		t.Fatalf("Failed to marshal Release: %v", err)
	}

	// Test unmarshaling
	var unmarshaledRelease Release
	err = json.Unmarshal(jsonData, &unmarshaledRelease)
	if err != nil {
		t.Fatalf("Failed to unmarshal Release: %v", err)
	}

	// Verify fields
	if unmarshaledRelease.ID != release.ID {
		t.Errorf("ID mismatch: got %v, want %v", unmarshaledRelease.ID, release.ID)
	}
	if unmarshaledRelease.Version != release.Version {
		t.Errorf("Version mismatch: got %v, want %v", unmarshaledRelease.Version, release.Version)
	}
	if unmarshaledRelease.Sequence != release.Sequence {
		t.Errorf("Sequence mismatch: got %v, want %v", unmarshaledRelease.Sequence, release.Sequence)
	}
	if unmarshaledRelease.Status != release.Status {
		t.Errorf("Status mismatch: got %v, want %v", unmarshaledRelease.Status, release.Status)
	}
	if unmarshaledRelease.IsRequired != release.IsRequired {
		t.Errorf("IsRequired mismatch: got %v, want %v", unmarshaledRelease.IsRequired, release.IsRequired)
	}

	// Verify metadata
	if len(unmarshaledRelease.Metadata) != len(release.Metadata) {
		t.Errorf("Metadata length mismatch: got %v, want %v", len(unmarshaledRelease.Metadata), len(release.Metadata))
	}
}

func TestIsValidSemanticVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{"valid basic version", "1.0.0", true},
		{"valid version with pre-release", "1.0.0-alpha", true},
		{"valid version with pre-release and number", "1.0.0-alpha.1", true},
		{"valid version with build metadata", "1.0.0+20130313144700", true},
		{"valid complex version", "1.0.0-beta.2+exp.sha.5114f85", true},
		{"invalid missing patch", "1.0", false},
		{"invalid with v prefix", "v1.0.0", false},
		{"invalid negative number", "-1.0.0", false},
		{"invalid leading zeros", "01.0.0", false},
		{"invalid letters in version", "1.a.0", false},
		{"empty version", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidSemanticVersion(tt.version); got != tt.want {
				t.Errorf("isValidSemanticVersion(%v) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

func TestIsValidReleaseStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{"valid draft", ReleaseStatusDraft, true},
		{"valid released", ReleaseStatusReleased, true},
		{"valid archived", ReleaseStatusArchived, true},
		{"valid superseded", ReleaseStatusSuperseded, true},
		{"invalid status", "invalid", false},
		{"empty status", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidReleaseStatus(tt.status); got != tt.want {
				t.Errorf("isValidReleaseStatus(%v) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestRelease_IsReleased(t *testing.T) {
	validTime := time.Now()

	tests := []struct {
		name    string
		release Release
		want    bool
	}{
		{
			name: "released with timestamp",
			release: Release{
				Status:     ReleaseStatusReleased,
				ReleasedAt: &validTime,
			},
			want: true,
		},
		{
			name: "released without timestamp",
			release: Release{
				Status: ReleaseStatusReleased,
			},
			want: false,
		},
		{
			name: "draft status",
			release: Release{
				Status: ReleaseStatusDraft,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.release.IsReleased(); got != tt.want {
				t.Errorf("Release.IsReleased() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelease_String(t *testing.T) {
	release := Release{
		ID:            "rel-123",
		ApplicationID: "app-456",
		Version:       "1.0.0",
		Sequence:      1,
		Status:        ReleaseStatusReleased,
	}

	str := release.String()
	expected := "Release{ID: rel-123, ApplicationID: app-456, Version: 1.0.0, Sequence: 1, Status: released}"

	if str != expected {
		t.Errorf("Release.String() = %v, want %v", str, expected)
	}
}
