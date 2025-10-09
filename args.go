package baselith

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type PathFolder string

// Command line arguments
var (
	ConfigYaml bool
	ConfigPath string
	Folder     PathFolder

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
	rootCmd.Flags().StringVar((*string)(&Folder), "folder", "migrations", "Folder containing migrations")
	rootCmd.Flags().StringVar(&Driver, "driver", "postgres", "Database driver (postgres, mysql, etc.)")
	rootCmd.Flags().StringVar(&Host, "host", "localhost", "Database host")
	rootCmd.Flags().IntVar(&Port, "port", 5432, "Database port")
	rootCmd.Flags().StringVar(&Dbname, "dbname", "", "Database name")
	rootCmd.Flags().StringVar(&User, "user", "", "Database user")
	rootCmd.Flags().StringVar(&Password, "password", "", "Database password")
	rootCmd.Flags().StringVar(&Sub, "sub", "up", "Subcommand to execute: up, down, to, redo, history, status")
	rootCmd.Flags().StringVar(&ToID, "to", "", "Target migration ID for 'to' or 'down' subcommands")
}

func (f PathFolder) String() string {
	// Normalize the path to handle both Linux and Windows path separators
	normalizedPath := string(f)
	// Replace any forward slashes with the OS-specific separator
	if filepath.Separator != '/' {
		normalizedPath = strings.ReplaceAll(normalizedPath, "/", string(filepath.Separator))
	}
	return normalizedPath
}

// Path returns the properly formatted path for the current OS
func (f PathFolder) Path() string {
	return filepath.FromSlash(string(f))
}

// JoinPath joins the folder path with the given filename using the OS-specific separator
func (f PathFolder) JoinPath(filename string) string {
	return filepath.Join(string(f), filename)
}
