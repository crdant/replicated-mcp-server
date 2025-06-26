package config

import (
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		flags       map[string]interface{}
		want        *Config
		wantErr     bool
		errContains string
	}{
		{
			name: "load from environment variables only",
			envVars: map[string]string{
				"REPLICATED_API_TOKEN": "test-token",
				"LOG_LEVEL":            "info",
				"TIMEOUT":              "60",
				"ENDPOINT":             "https://api.example.com",
			},
			want: &Config{
				APIToken: "test-token",
				LogLevel: "info",
				Timeout:  60 * time.Second,
				Endpoint: "https://api.example.com",
			},
			wantErr: false,
		},
		{
			name: "minimal configuration with defaults",
			envVars: map[string]string{
				"REPLICATED_API_TOKEN": "test-token",
			},
			want: &Config{
				APIToken: "test-token",
				LogLevel: DefaultLogLevel,
				Timeout:  DefaultTimeout,
				Endpoint: "",
			},
			wantErr: false,
		},
		{
			name: "CLI flags override environment variables",
			envVars: map[string]string{
				"REPLICATED_API_TOKEN": "env-token",
				"LOG_LEVEL":            "error",
				"TIMEOUT":              "30",
			},
			flags: map[string]interface{}{
				"api-token": "flag-token",
				"log-level": "debug",
				"timeout":   120,
			},
			want: &Config{
				APIToken: "flag-token",
				LogLevel: "debug",
				Timeout:  120 * time.Second,
				Endpoint: "",
			},
			wantErr: false,
		},
		{
			name:        "missing API token",
			envVars:     map[string]string{},
			wantErr:     true,
			errContains: "API token is required",
		},
		{
			name: "invalid log level",
			envVars: map[string]string{
				"REPLICATED_API_TOKEN": "test-token",
				"LOG_LEVEL":            "invalid",
			},
			wantErr:     true,
			errContains: "invalid log level",
		},
		{
			name: "invalid timeout",
			envVars: map[string]string{
				"REPLICATED_API_TOKEN": "test-token",
				"TIMEOUT":              "abc",
			},
			wantErr:     true,
			errContains: "invalid TIMEOUT environment variable",
		},
		{
			name: "timeout out of range",
			envVars: map[string]string{
				"REPLICATED_API_TOKEN": "test-token",
				"TIMEOUT":              "500",
			},
			wantErr:     true,
			errContains: "timeout must be between",
		},
		{
			name: "invalid endpoint URL",
			envVars: map[string]string{
				"REPLICATED_API_TOKEN": "test-token",
				"ENDPOINT":             "not-a-url",
			},
			wantErr:     true,
			errContains: "invalid endpoint URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			clearTestEnv()

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer clearTestEnv()

			// Create test command with flags
			cmd := createTestCommand()
			
			// Build command line args from flags
			var args []string
			if tt.flags != nil {
				for flag, value := range tt.flags {
					switch v := value.(type) {
					case string:
						args = append(args, "--"+flag, v)
					case int:
						args = append(args, "--"+flag, strconv.Itoa(v))
					}
				}
			}
			
			// Parse the flags to simulate actual command execution
			if len(args) > 0 {
				cmd.SetArgs(args)
				cmd.ParseFlags(args)
			}

			got, err := Load(cmd)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Load() expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Load() error = %v, expected to contain %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Load() unexpected error = %v", err)
				return
			}

			if got.APIToken != tt.want.APIToken {
				t.Errorf("Load() APIToken = %v, want %v", got.APIToken, tt.want.APIToken)
			}
			if got.LogLevel != tt.want.LogLevel {
				t.Errorf("Load() LogLevel = %v, want %v", got.LogLevel, tt.want.LogLevel)
			}
			if got.Timeout != tt.want.Timeout {
				t.Errorf("Load() Timeout = %v, want %v", got.Timeout, tt.want.Timeout)
			}
			if got.Endpoint != tt.want.Endpoint {
				t.Errorf("Load() Endpoint = %v, want %v", got.Endpoint, tt.want.Endpoint)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		wantErr     bool
		errContains string
	}{
		{
			name: "valid configuration",
			config: &Config{
				APIToken: "test-token",
				LogLevel: "info",
				Timeout:  30 * time.Second,
				Endpoint: "https://api.example.com",
			},
			wantErr: false,
		},
		{
			name: "valid minimal configuration",
			config: &Config{
				APIToken: "test-token",
				LogLevel: "fatal",
				Timeout:  1 * time.Second,
				Endpoint: "",
			},
			wantErr: false,
		},
		{
			name: "missing API token",
			config: &Config{
				APIToken: "",
				LogLevel: "info",
				Timeout:  30 * time.Second,
			},
			wantErr:     true,
			errContains: "API token is required",
		},
		{
			name: "invalid log level",
			config: &Config{
				APIToken: "test-token",
				LogLevel: "INVALID",
				Timeout:  30 * time.Second,
			},
			wantErr:     true,
			errContains: "invalid log level",
		},
		{
			name: "timeout too short",
			config: &Config{
				APIToken: "test-token",
				LogLevel: "info",
				Timeout:  500 * time.Millisecond,
			},
			wantErr:     true,
			errContains: "timeout must be between",
		},
		{
			name: "timeout too long",
			config: &Config{
				APIToken: "test-token",
				LogLevel: "info",
				Timeout:  400 * time.Second,
			},
			wantErr:     true,
			errContains: "timeout must be between",
		},
		{
			name: "invalid endpoint",
			config: &Config{
				APIToken: "test-token",
				LogLevel: "info",
				Timeout:  30 * time.Second,
				Endpoint: "not-a-valid-url",
			},
			wantErr:     true,
			errContains: "invalid endpoint URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Validate() error = %v, expected to contain %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Validate() unexpected error = %v", err)
			}
		})
	}
}

func TestIsValidLogLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  bool
	}{
		{"valid fatal", "fatal", true},
		{"valid error", "error", true},
		{"valid info", "info", true},
		{"valid debug", "debug", true},
		{"valid trace", "trace", true},
		{"valid uppercase", "INFO", true},
		{"valid mixed case", "Debug", true},
		{"invalid level", "invalid", false},
		{"empty string", "", false},
		{"number", "123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidLogLevel(tt.level); got != tt.want {
				t.Errorf("isValidLogLevel(%q) = %v, want %v", tt.level, got, tt.want)
			}
		})
	}
}

func TestConfig_String(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   []string // substrings that should be present
	}{
		{
			name: "config with all fields",
			config: &Config{
				APIToken: "secret-token",
				LogLevel: "info",
				Timeout:  30 * time.Second,
				Endpoint: "https://api.example.com",
			},
			want: []string{"(set)", "info", "30s", "https://api.example.com"},
		},
		{
			name: "config without token",
			config: &Config{
				APIToken: "",
				LogLevel: "fatal",
				Timeout:  15 * time.Second,
				Endpoint: "",
			},
			want: []string{"(not set)", "fatal", "15s", "(default)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.String()
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf("Config.String() = %v, expected to contain %v", got, want)
				}
			}
			// Ensure we don't leak the actual token
			if tt.config.APIToken != "" && strings.Contains(got, tt.config.APIToken) {
				t.Errorf("Config.String() = %v, should not contain actual token", got)
			}
		})
	}
}

// Helper functions for testing

func clearTestEnv() {
	os.Unsetenv("REPLICATED_API_TOKEN")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("TIMEOUT")
	os.Unsetenv("ENDPOINT")
}

func createTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {},
	}

	// Add the same flags as the real application
	cmd.PersistentFlags().String("api-token", "", "Replicated Vendor Portal API token")
	cmd.PersistentFlags().String("log-level", "fatal", "Log level (fatal, error, info, debug, trace)")
	cmd.PersistentFlags().Int("timeout", 30, "API request timeout in seconds")
	cmd.PersistentFlags().String("endpoint", "", "API endpoint (hidden)")

	return cmd
}