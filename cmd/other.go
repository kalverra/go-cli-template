package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var otherCmd = &cobra.Command{
	Use:   "other",
	Short: "Other command",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Example value: %s\n", cfg.Nested.ExampleValue)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(otherCmd)

	otherCmd.Flags().String("example-value", "", "Value to set")
}
