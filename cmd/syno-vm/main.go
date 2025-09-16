package main

import (
	"fmt"
	"os"

	"github.com/scttfrdmn/syno-vm/internal/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	// Set version info for the CLI
	cmd.SetVersionInfo(version, commit, date, builtBy)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}