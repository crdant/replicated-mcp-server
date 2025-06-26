// Package config provides configuration management for the Replicated MCP Server.
// It supports loading configuration from environment variables and CLI flags,
// with comprehensive validation and helpful error messages.
package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Config represents the application configuration
type Config struct {
	APIToken string
	LogLevel string
	Timeout  time.Duration
	Endpoint string
}

// Validation constants
const (
	DefaultLogLevel = "fatal"
	DefaultTimeout  = 30 * time.Second
	MinTimeout      = 1 * time.Second
	MaxTimeout      = 300 * time.Second
)

// ValidLogLevels contains all supported log level names
var ValidLogLevels = []string{"fatal", "error", "info", "debug", "trace"}

// Load creates a new Config by loading from environment variables and CLI flags
// CLI flags take precedence over environment variables
func Load(cmd *cobra.Command) (*Config, error) {
	config := &Config{}

	// Load from environment variables first
	if err := config.loadFromEnv(); err != nil {
		return nil, fmt.Errorf("failed to load configuration from environment: %w", err)
	}

	// Override with CLI flags if provided
	if err := config.loadFromFlags(cmd.Flags()); err != nil {
		return nil, fmt.Errorf("failed to load configuration from flags: %w", err)
	}

	// Validate the final configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// loadFromEnv loads configuration from environment variables
func (c *Config) loadFromEnv() error {
	// API Token (required)
	if token := os.Getenv("REPLICATED_API_TOKEN"); token != "" {
		c.APIToken = token
	}

	// Log Level (optional, has default)
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		c.LogLevel = level
	} else {
		c.LogLevel = DefaultLogLevel
	}

	// Timeout (optional, has default)
	if timeoutStr := os.Getenv("TIMEOUT"); timeoutStr != "" {
		timeout, err := strconv.Atoi(timeoutStr)
		if err != nil {
			return fmt.Errorf("invalid TIMEOUT environment variable '%s': must be a number of seconds", timeoutStr)
		}
		c.Timeout = time.Duration(timeout) * time.Second
	} else {
		c.Timeout = DefaultTimeout
	}

	// Endpoint (optional)
	if endpoint := os.Getenv("ENDPOINT"); endpoint != "" {
		c.Endpoint = endpoint
	}

	return nil
}

// loadFromFlags loads configuration from CLI flags, overriding environment variables
func (c *Config) loadFromFlags(flags *pflag.FlagSet) error {
	// API Token
	if flags.Changed("api-token") {
		token, err := flags.GetString("api-token")
		if err != nil {
			return fmt.Errorf("failed to get api-token flag: %w", err)
		}
		c.APIToken = token
	}

	// Log Level
	if flags.Changed("log-level") {
		level, err := flags.GetString("log-level")
		if err != nil {
			return fmt.Errorf("failed to get log-level flag: %w", err)
		}
		c.LogLevel = level
	}

	// Timeout
	if flags.Changed("timeout") {
		timeoutSeconds, err := flags.GetInt("timeout")
		if err != nil {
			return fmt.Errorf("failed to get timeout flag: %w", err)
		}
		c.Timeout = time.Duration(timeoutSeconds) * time.Second
	}

	// Endpoint
	if flags.Changed("endpoint") {
		endpoint, err := flags.GetString("endpoint")
		if err != nil {
			return fmt.Errorf("failed to get endpoint flag: %w", err)
		}
		c.Endpoint = endpoint
	}

	return nil
}

// Validate ensures the configuration is valid
func (c *Config) Validate() error {
	var errors []string

	// Validate API Token
	if c.APIToken == "" {
		errors = append(errors, "API token is required. Set REPLICATED_API_TOKEN environment variable "+
			"or use --api-token flag")
	}

	// Validate Log Level
	if !isValidLogLevel(c.LogLevel) {
		errors = append(errors, fmt.Sprintf("invalid log level '%s'. Valid levels are: %s",
			c.LogLevel, strings.Join(ValidLogLevels, ", ")))
	}

	// Validate Timeout
	if c.Timeout < MinTimeout || c.Timeout > MaxTimeout {
		errors = append(errors, fmt.Sprintf("timeout must be between %v and %v seconds, got %v",
			MinTimeout.Seconds(), MaxTimeout.Seconds(), c.Timeout.Seconds()))
	}

	// Validate Endpoint (if provided)
	if c.Endpoint != "" {
		if u, err := url.Parse(c.Endpoint); err != nil {
			errors = append(errors, fmt.Sprintf("invalid endpoint URL '%s': %v", c.Endpoint, err))
		} else if u.Scheme == "" || u.Host == "" {
			errors = append(errors, fmt.Sprintf("invalid endpoint URL '%s': must include scheme and host "+
				"(e.g., https://api.example.com)", c.Endpoint))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration validation errors:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

// isValidLogLevel checks if the provided log level is valid
func isValidLogLevel(level string) bool {
	level = strings.ToLower(level)
	for _, valid := range ValidLogLevels {
		if level == valid {
			return true
		}
	}
	return false
}

// String returns a string representation of the configuration (without sensitive data)
func (c *Config) String() string {
	endpoint := c.Endpoint
	if endpoint == "" {
		endpoint = "(default)"
	}

	token := "(not set)"
	if c.APIToken != "" {
		token = "(set)"
	}

	return fmt.Sprintf("Config{APIToken: %s, LogLevel: %s, Timeout: %v, Endpoint: %s}",
		token, c.LogLevel, c.Timeout, endpoint)
}
