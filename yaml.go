package baselith

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// DBConfigYAML represents the expected structure of the YAML config file.
type DBConfigYAML struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Dbname   string `yaml:"dbname"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Schema   string `yaml:"schema"`
}

// ReadConfigYAML reads and parses the YAML config file if ConfigYaml is true and ConfigPath is set.
// Returns DBConfigYAML and error.
func ReadConfigYAML() (*DBConfigYAML, error) {
	if !ConfigYaml || ConfigPath == "" {
		return nil, fmt.Errorf("YAML config not enabled or path not set")
	}
	file, err := os.Open(ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open YAML file: %w", err)
	}
	defer file.Close()

	var cfg DBConfigYAML
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return &cfg, nil
}
