package models

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestChannel_Validate(t *testing.T) {
	validTime := time.Now()
	laterTime := validTime.Add(time.Hour)

	tests := []struct {
		name        string
		channel     Channel
		wantErr     bool
		errContains []string
	}{
		{
			name: "valid channel",
			channel: Channel{
				ID:              "ch-123",
				ApplicationID:   "app-456",
				Name:            "Stable",
				Description:     "Stable release channel",
				ReleaseID:       "rel-789",
				ReleaseSequence: 1,
				CreatedAt:       validTime,
				UpdatedAt:       laterTime,
				IsDefault:       true,
				IsArchived:      false,
				ChannelSlug:     "stable",
			},
			wantErr: false,
		},
		{
			name: "minimal valid channel",
			channel: Channel{
				ID:            "ch-123",
				ApplicationID: "app-456",
				Name:          "Beta",
				CreatedAt:     validTime,
				UpdatedAt:     validTime,
				IsDefault:     false,
				IsArchived:    false,
				ChannelSlug:   "beta",
			},
			wantErr: false,
		},
		{
			name: "valid archived channel",
			channel: Channel{
				ID:            "ch-123",
				ApplicationID: "app-456",
				Name:          "Old Channel",
				CreatedAt:     validTime,
				UpdatedAt:     validTime,
				ArchivedAt:    &validTime,
				IsDefault:     false,
				IsArchived:    true,
				ChannelSlug:   "old-channel",
			},
			wantErr: false,
		},
		{
			name:        "missing ID",
			channel:     Channel{},
			wantErr:     true,
			errContains: []string{"channel ID is required"},
		},
		{
			name: "missing application ID",
			channel: Channel{
				ID: "ch-123",
			},
			wantErr:     true,
			errContains: []string{"application ID is required"},
		},
		{
			name: "missing name",
			channel: Channel{
				ID:            "ch-123",
				ApplicationID: "app-456",
			},
			wantErr:     true,
			errContains: []string{"channel name is required"},
		},
		{
			name: "name too long",
			channel: Channel{
				ID:            "ch-123",
				ApplicationID: "app-456",
				Name:          strings.Repeat("a", 101),
			},
			wantErr:     true,
			errContains: []string{"channel name must be 100 characters or less"},
		},
		{
			name: "missing channel slug",
			channel: Channel{
				ID:            "ch-123",
				ApplicationID: "app-456",
				Name:          "Test Channel",
			},
			wantErr:     true,
			errContains: []string{"channel slug is required"},
		},
		{
			name: "invalid channel slug with uppercase",
			channel: Channel{
				ID:            "ch-123",
				ApplicationID: "app-456",
				Name:          "Test Channel",
				ChannelSlug:   "Test-Channel",
			},
			wantErr:     true,
			errContains: []string{"channel slug must contain only lowercase letters, numbers, and hyphens"},
		},
		{
			name: "invalid channel slug with underscore",
			channel: Channel{
				ID:            "ch-123",
				ApplicationID: "app-456",
				Name:          "Test Channel",
				ChannelSlug:   "test_channel",
			},
			wantErr:     true,
			errContains: []string{"channel slug must contain only lowercase letters, numbers, and hyphens"},
		},
		{
			name: "invalid channel slug starting with hyphen",
			channel: Channel{
				ID:            "ch-123",
				ApplicationID: "app-456",
				Name:          "Test Channel",
				ChannelSlug:   "-test-channel",
			},
			wantErr:     true,
			errContains: []string{"channel slug must contain only lowercase letters, numbers, and hyphens"},
		},
		{
			name: "invalid channel slug ending with hyphen",
			channel: Channel{
				ID:            "ch-123",
				ApplicationID: "app-456",
				Name:          "Test Channel",
				ChannelSlug:   "test-channel-",
			},
			wantErr:     true,
			errContains: []string{"channel slug must contain only lowercase letters, numbers, and hyphens"},
		},
		{
			name: "missing timestamps",
			channel: Channel{
				ID:            "ch-123",
				ApplicationID: "app-456",
				Name:          "Test Channel",
				ChannelSlug:   "test-channel",
			},
			wantErr:     true,
			errContains: []string{"created_at timestamp is required", "updated_at timestamp is required"},
		},
		{
			name: "updated_at before created_at",
			channel: Channel{
				ID:            "ch-123",
				ApplicationID: "app-456",
				Name:          "Test Channel",
				ChannelSlug:   "test-channel",
				CreatedAt:     laterTime,
				UpdatedAt:     validTime,
			},
			wantErr:     true,
			errContains: []string{"updated_at must be equal to or after created_at"},
		},
		{
			name: "archived_at before created_at",
			channel: Channel{
				ID:            "ch-123",
				ApplicationID: "app-456",
				Name:          "Test Channel",
				ChannelSlug:   "test-channel",
				CreatedAt:     laterTime,
				UpdatedAt:     laterTime,
				ArchivedAt:    &validTime,
			},
			wantErr:     true,
			errContains: []string{"archived_at must be equal to or after created_at"},
		},
		{
			name: "archived_at set but is_archived false",
			channel: Channel{
				ID:            "ch-123",
				ApplicationID: "app-456",
				Name:          "Test Channel",
				ChannelSlug:   "test-channel",
				CreatedAt:     validTime,
				UpdatedAt:     validTime,
				ArchivedAt:    &validTime,
				IsArchived:    false,
			},
			wantErr:     true,
			errContains: []string{"is_archived must be true when archived_at is set"},
		},
		{
			name: "is_archived true but archived_at not set",
			channel: Channel{
				ID:            "ch-123",
				ApplicationID: "app-456",
				Name:          "Test Channel",
				ChannelSlug:   "test-channel",
				CreatedAt:     validTime,
				UpdatedAt:     validTime,
				IsArchived:    true,
			},
			wantErr:     true,
			errContains: []string{"archived_at is required when is_archived is true"},
		},
		{
			name: "release_id without sequence",
			channel: Channel{
				ID:            "ch-123",
				ApplicationID: "app-456",
				Name:          "Test Channel",
				ChannelSlug:   "test-channel",
				CreatedAt:     validTime,
				UpdatedAt:     validTime,
				ReleaseID:     "rel-789",
			},
			wantErr:     true,
			errContains: []string{"release_sequence must be positive when release_id is provided"},
		},
		{
			name: "release_sequence without id",
			channel: Channel{
				ID:              "ch-123",
				ApplicationID:   "app-456",
				Name:            "Test Channel",
				ChannelSlug:     "test-channel",
				CreatedAt:       validTime,
				UpdatedAt:       validTime,
				ReleaseSequence: 1,
			},
			wantErr:     true,
			errContains: []string{"release_id is required when release_sequence is provided"},
		},
		{
			name: "description too long",
			channel: Channel{
				ID:            "ch-123",
				ApplicationID: "app-456",
				Name:          "Test Channel",
				ChannelSlug:   "test-channel",
				CreatedAt:     validTime,
				UpdatedAt:     validTime,
				Description:   strings.Repeat("a", 501),
			},
			wantErr:     true,
			errContains: []string{"channel description must be 500 characters or less"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.channel.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Channel.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				for _, expectedErr := range tt.errContains {
					if !strings.Contains(err.Error(), expectedErr) {
						t.Errorf("Channel.Validate() error = %v, should contain %v", err, expectedErr)
					}
				}
			}
		})
	}
}

func TestChannel_JSONMarshaling(t *testing.T) {
	validTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	channel := Channel{
		ID:              "ch-123",
		ApplicationID:   "app-456",
		Name:            "Stable",
		Description:     "Stable release channel",
		ReleaseID:       "rel-789",
		ReleaseSequence: 1,
		CreatedAt:       validTime,
		UpdatedAt:       validTime,
		IsDefault:       true,
		IsArchived:      false,
		ChannelSlug:     "stable",
	}

	// Test marshaling
	jsonData, err := json.Marshal(channel)
	if err != nil {
		t.Fatalf("Failed to marshal Channel: %v", err)
	}

	// Test unmarshaling
	var unmarshaledChannel Channel
	err = json.Unmarshal(jsonData, &unmarshaledChannel)
	if err != nil {
		t.Fatalf("Failed to unmarshal Channel: %v", err)
	}

	// Verify fields
	if unmarshaledChannel.ID != channel.ID {
		t.Errorf("ID mismatch: got %v, want %v", unmarshaledChannel.ID, channel.ID)
	}
	if unmarshaledChannel.Name != channel.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaledChannel.Name, channel.Name)
	}
	if unmarshaledChannel.ChannelSlug != channel.ChannelSlug {
		t.Errorf("ChannelSlug mismatch: got %v, want %v", unmarshaledChannel.ChannelSlug, channel.ChannelSlug)
	}
	if unmarshaledChannel.IsDefault != channel.IsDefault {
		t.Errorf("IsDefault mismatch: got %v, want %v", unmarshaledChannel.IsDefault, channel.IsDefault)
	}
	if unmarshaledChannel.ReleaseSequence != channel.ReleaseSequence {
		t.Errorf("ReleaseSequence mismatch: got %v, want %v", unmarshaledChannel.ReleaseSequence, channel.ReleaseSequence)
	}
}

func TestIsValidChannelSlug(t *testing.T) {
	validSlugs := []string{"stable", "beta123", "release-candidate", "release-candidate-v2"}
	invalidSlugs := []string{"", "Stable", "release_candidate", "release candidate", "-stable", "stable-", "stable@channel"}

	testSlugValidation(t, "isValidChannelSlug", isValidChannelSlug, validSlugs, invalidSlugs)
}

func TestChannel_HasRelease(t *testing.T) {
	tests := []struct {
		name    string
		channel Channel
		want    bool
	}{
		{
			name: "channel with release",
			channel: Channel{
				ReleaseID:       "rel-123",
				ReleaseSequence: 1,
			},
			want: true,
		},
		{
			name: "channel without release ID",
			channel: Channel{
				ReleaseSequence: 1,
			},
			want: false,
		},
		{
			name: "channel without release sequence",
			channel: Channel{
				ReleaseID: "rel-123",
			},
			want: false,
		},
		{
			name:    "channel without release",
			channel: Channel{},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.channel.HasRelease(); got != tt.want {
				t.Errorf("Channel.HasRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannel_IsActive(t *testing.T) {
	activeChannel := Channel{IsArchived: false}
	archivedChannel := Channel{IsArchived: true}

	testIsActiveValidation(t, activeChannel.IsActive, archivedChannel.IsActive)
}

func TestChannel_String(t *testing.T) {
	channel := Channel{
		ID:            "ch-123",
		ApplicationID: "app-456",
		Name:          "Stable",
		ChannelSlug:   "stable",
		IsDefault:     true,
		IsArchived:    false,
	}

	str := channel.String()
	expected := "Channel{ID: ch-123, ApplicationID: app-456, Name: Stable, Slug: stable, IsDefault: true, IsArchived: false}"

	if str != expected {
		t.Errorf("Channel.String() = %v, want %v", str, expected)
	}
}
