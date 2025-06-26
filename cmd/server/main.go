package main

import (
	"fmt"
	"os"

	"github.com/crdant/replicated-mcp-server/pkg/config"
	"github.com/crdant/replicated-mcp-server/pkg/logging"
	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	buildDate = "unknown"
	commit    = "none"
)

var rootCmd = &cobra.Command{
	Use:   "replicated-mcp-server",
	Short: "MCP server for Replicated Vendor Portal API",
	Long: `A Machine Chat Protocol (MCP) server that interfaces with the Replicated Vendor Portal API, 
enabling AI agents to interact with Replicated Vendor Portal accounts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServer(cmd, args)
	},
	Version: fmt.Sprintf("%s (Built: %s, Commit: %s)", version, buildDate, commit),
}

func init() {
	// Define flags and configuration settings
	rootCmd.PersistentFlags().String("api-token", "", "Replicated Vendor Portal API token")
	rootCmd.PersistentFlags().String("log-level", "fatal", "Log level (fatal, error, info, debug, trace)")
	rootCmd.PersistentFlags().Int("timeout", 30, "API request timeout in seconds")
	rootCmd.PersistentFlags().String("endpoint", "", "API endpoint (hidden)")
	rootCmd.PersistentFlags().MarkHidden("endpoint")
}

func runServer(cmd *cobra.Command, args []string) error {
	// Load configuration from environment variables and CLI flags
	cfg, err := config.Load(cmd)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize structured logger
	logger := logging.NewLogger(cfg.LogLevel)

	// Log startup information
	logger.Info("Replicated MCP Server starting",
		"version", version,
		"build_date", buildDate,
		"commit", commit,
		"config", cfg.String())

	// TODO: Initialize MCP server (Step 4)
	// TODO: Initialize API client (Step 3)
	// TODO: Start server

	logger.Info("Server initialization complete - MCP server would start here")
	
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		// Use fmt.Fprintf to ensure error goes to stderr
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}