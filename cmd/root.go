// Package cmd contains command execution logic.
package cmd

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"

	"github.com/kalverra/go-cli-template/internal/config"
)

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "go-cli-template",
	Short: "Template CLI application in Go",
	Long:  `Template CLI application in Go`,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		var err error
		cfg, err = config.Load(config.WithFlags(cmd.Flags()))
		return err
	},
}

func init() {
	rootCmd.PersistentFlags().
		StringVarP(&cfg.LogLevel, "log-level", "l", config.DefaultLogLevel, "Log level (env: LOG_LEVEL)")
}

// Execute runs the root command.
func Execute() {
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}
