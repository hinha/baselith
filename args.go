package baselith

import "github.com/spf13/cobra"

// Command line arguments
var (
	ConfigYaml bool
	ConfigPath string

	// Database connection flags
	Driver   string
	Host     string
	Port     int
	Dbname   string
	User     string
	Password string
	Schema   string
	Sub      string
	ToID     string
)

func ReadFlags(rootCmd *cobra.Command) {
	// YAML configuration flag
	rootCmd.Flags().BoolVar(&ConfigYaml, "yaml", false, "output YAML")
	rootCmd.Flags().StringVar(&ConfigPath, "config", "", "path to config file")

	// Direct database connection flags
	rootCmd.Flags().StringVar(&Driver, "driver", "postgres", "Database driver (postgres, mysql, etc.)")
	rootCmd.Flags().StringVar(&Host, "host", "localhost", "Database host")
	rootCmd.Flags().IntVar(&Port, "port", 5432, "Database port")
	rootCmd.Flags().StringVar(&Dbname, "dbname", "", "Database name")
	rootCmd.Flags().StringVar(&User, "user", "", "Database user")
	rootCmd.Flags().StringVar(&Password, "password", "", "Database password")
	rootCmd.Flags().StringVar(&Schema, "schema", "public", "Database schema (for PostgreSQL)")
	rootCmd.Flags().StringVar(&Sub, "s", "up", "Subcommand to execute: up, down, to, redo, history, status")
	rootCmd.Flags().StringVar(&ToID, "to", "", "Target migration ID for 'to' or 'down' subcommands")
}
