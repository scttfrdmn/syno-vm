package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool

	// Version info set by main
	appVersion = "0.1.0"
	buildCommit = "unknown"
	buildDate = "unknown"
	builtBy = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "syno-vm",
	Short: "A CLI tool for managing Synology Virtual Machine Manager",
	Long: `syno-vm is a command-line tool for managing virtual machines on Synology NAS
devices with Virtual Machine Manager (VMM). It provides an easy way to create,
start, stop, and manage VMs through Synology's VMM API.

This tool is adapted from qnap-vm for Synology DSM 7.x+ systems.`,
	Version: appVersion,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

// SetVersionInfo sets the version information from build-time variables
func SetVersionInfo(version, commit, date, builder string) {
	if version != "" && version != "dev" {
		appVersion = version
	}
	if commit != "" && commit != "none" {
		buildCommit = commit
	}
	if date != "" && date != "unknown" {
		buildDate = date
	}
	if builder != "" && builder != "unknown" {
		builtBy = builder
	}
	rootCmd.Version = appVersion
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.syno-vm/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")) // nolint:errcheck // CLI setup
}

// initConfig reads in config file and ENV variables.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Search config in home directory with name ".syno-vm" (without extension).
		configDir := home + "/.syno-vm"
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")

		// Create config directory if it doesn't exist
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			_ = os.MkdirAll(configDir, 0755) // Best effort directory creation
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// Set default values
	viper.SetDefault("port", 22)
	viper.SetDefault("timeout", 30)
}