package main

import (
	"os"

	"github.com/hinha/baselith/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
