package main

import (
	"fmt"
	"os"

	"github.com/hinha/baselith"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "baselith",
	Short: "Database migration management tool",
	Long:  `A cross-platform CLI tool for managing database migrations, making schema changes easier without manual intervention.`,
	Run:   baselith.Run,
}

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of Baselith",
		Long:  `All software has versions. This is Baselith's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "Baselith v1.0.0")
		},
	})
	baselith.ReadFlags(rootCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
