package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "baselith",
	Short: "Baselith is a cross-platform CLI application",
	Long:  `A cross-platform command-line application built with Cobra that works on Linux, macOS, and Windows.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), "Welcome to Baselith! Use --help for more information.")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(greetCmd)
}
