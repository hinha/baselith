package cmd

import (
	"bytes"
	"testing"
)

func TestRootCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("rootCmd.Execute() error = %v", err)
	}

	output := buf.String()
	expected := "Welcome to Baselith! Use --help for more information.\n"
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}
}

func TestVersionCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"version"})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("versionCmd.Execute() error = %v", err)
	}

	output := buf.String()
	expected := "Baselith v1.0.0\n"
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}
}

func TestGreetCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "greet without name",
			args:     []string{"greet"},
			expected: "Hello, World!\n",
		},
		{
			name:     "greet with name",
			args:     []string{"greet", "--name", "Alice"},
			expected: "Hello, Alice!\n",
		},
		{
			name:     "greet with short flag",
			args:     []string{"greet", "-n", "Bob"},
			expected: "Hello, Bob!\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(tt.args)

			err := rootCmd.Execute()
			if err != nil {
				t.Errorf("Execute() error = %v", err)
			}

			output := buf.String()
			if output != tt.expected {
				t.Errorf("Expected output %q, got %q", tt.expected, output)
			}
		})
	}
}
