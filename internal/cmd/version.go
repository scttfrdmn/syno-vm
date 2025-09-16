package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Show detailed version information including build metadata.`,
	RunE:  runVersion,
}

var (
	showBuildInfo bool
)

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolVarP(&showBuildInfo, "build", "b", false, "Show build information")
}

func runVersion(cmd *cobra.Command, args []string) error {
	fmt.Printf("syno-vm version %s\n", appVersion)

	if showBuildInfo {
		fmt.Printf("\nBuild Information:\n")
		fmt.Printf("  Version:    %s\n", appVersion)
		fmt.Printf("  Commit:     %s\n", buildCommit)
		fmt.Printf("  Build Date: %s\n", buildDate)
		fmt.Printf("  Built By:   %s\n", builtBy)
		fmt.Printf("  Go Version: %s\n", runtime.Version())
		fmt.Printf("  Platform:   %s/%s\n", runtime.GOOS, runtime.GOARCH)
	}

	return nil
}