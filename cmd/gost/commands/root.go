package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "github.com/JsCodeDevlopment/gost",
	Short: "Gost CLI - Scaffolding for Go applications with NestJS-like architecture",
	Long: `Gost CLI is a powerful tool to quickly scaffold new projects, modules, and CRUDs
following the Gost Framework standards of modularity and clean architecture.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
}
