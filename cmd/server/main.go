package main

import (
	"fmt"
	"os"

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
	Run: func(cmd *cobra.Command, args []string) {
		// Main server logic will go here
		fmt.Println("Replicated MCP Server starting...")
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

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}