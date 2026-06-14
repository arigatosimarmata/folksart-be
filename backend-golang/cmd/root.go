package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"react-example/backend-golang/config"
)

var rootCmd = &cobra.Command{
	Use:   "iam-api",
	Short: "IAM Governance API CLI",
	Long:  `CLI for managing IAM Governance API Server, Migrations, and Schedulers.`,
}

func Execute() {
	config.LoadConfig()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add global flags here if needed
}
