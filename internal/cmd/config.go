package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage syno-vm configuration",
	Long:  `Configure connection settings for your Synology NAS.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	Long:  `Set configuration values for syno-vm.`,
	RunE:  runConfigSet,
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get configuration values",
	Long:  `Get configuration values from syno-vm config.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runConfigGet,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration",
	Long:  `List all configuration values.`,
	RunE:  runConfigList,
}

var (
	host     string
	username string
	password string
	port     int
	keyfile  string
	timeout  int
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)

	// Set command flags
	configSetCmd.Flags().StringVar(&host, "host", "", "Synology NAS hostname or IP address")
	configSetCmd.Flags().StringVar(&username, "username", "", "Username for authentication")
	configSetCmd.Flags().StringVar(&password, "password", "", "Password for Web API authentication")
	configSetCmd.Flags().IntVar(&port, "port", 22, "SSH port")
	configSetCmd.Flags().StringVar(&keyfile, "keyfile", "", "SSH private key file path")
	configSetCmd.Flags().IntVar(&timeout, "timeout", 30, "Connection timeout in seconds")
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	configChanged := false

	if cmd.Flags().Changed("host") {
		viper.Set("host", host)
		configChanged = true
		fmt.Printf("Set host: %s\n", host)
	}

	if cmd.Flags().Changed("username") {
		viper.Set("username", username)
		configChanged = true
		fmt.Printf("Set username: %s\n", username)
	}

	if cmd.Flags().Changed("password") {
		viper.Set("password", password)
		configChanged = true
		fmt.Printf("Set password: [hidden]\n")
	}

	if cmd.Flags().Changed("port") {
		viper.Set("port", port)
		configChanged = true
		fmt.Printf("Set port: %d\n", port)
	}

	if cmd.Flags().Changed("keyfile") {
		// Expand tilde to home directory
		if keyfile[:2] == "~/" {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}
			keyfile = filepath.Join(home, keyfile[2:])
		}
		viper.Set("keyfile", keyfile)
		configChanged = true
		fmt.Printf("Set keyfile: %s\n", keyfile)
	}

	if cmd.Flags().Changed("timeout") {
		viper.Set("timeout", timeout)
		configChanged = true
		fmt.Printf("Set timeout: %d\n", timeout)
	}

	if !configChanged {
		return fmt.Errorf("no configuration values provided")
	}

	// Write config to file
	return viper.WriteConfig()
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return runConfigList(cmd, args)
	}

	key := args[0]
	value := viper.Get(key)
	if value == nil {
		return fmt.Errorf("configuration key '%s' not found", key)
	}

	fmt.Printf("%s: %v\n", key, value)
	return nil
}

func runConfigList(cmd *cobra.Command, args []string) error {
	fmt.Println("Current configuration:")

	keys := []string{"host", "username", "password", "port", "keyfile", "timeout"}
	for _, key := range keys {
		value := viper.Get(key)
		if value != nil {
			if key == "password" {
				fmt.Printf("  %s: [hidden]\n", key)
			} else {
				fmt.Printf("  %s: %v\n", key, value)
			}
		}
	}

	return nil
}