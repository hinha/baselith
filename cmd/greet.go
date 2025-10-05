package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	name string
)

var greetCmd = &cobra.Command{
	Use:   "greet",
	Short: "Greet someone",
	Long:  `Greet a person by name. This demonstrates a simple command with flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		if name == "" {
			name = "World"
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Hello, %s!\n", name)
	},
}

func init() {
	greetCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the person to greet")
}
