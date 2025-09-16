package main

import (
	"fmt"
	"os"

	"github.com/scttfrdmn/syno-vm/internal/synology"
	"github.com/spf13/viper"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <host> <username>\n", os.Args[0])
		os.Exit(1)
	}

	host := os.Args[1]
	username := os.Args[2]

	// Set up configuration
	viper.Set("host", host)
	viper.Set("username", username)
	viper.Set("port", 22)
	viper.Set("timeout", 30)

	fmt.Printf("Testing connection to %s@%s...\n", username, host)

	client, err := synology.NewClient()
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		os.Exit(1)
	}

	// Test basic SSH connection
	output, err := client.ExecuteCommand("echo 'SSH connection successful'")
	if err != nil {
		fmt.Printf("SSH connection failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SSH test result: %s", output)

	// Check if synowebapi is available
	output, err = client.ExecuteCommand("which synowebapi 2>/dev/null || echo 'synowebapi not found'")
	if err != nil {
		fmt.Printf("Failed to check synowebapi: %v\n", err)
	} else {
		fmt.Printf("synowebapi check: %s", output)
	}

	// Check for VMM-related processes
	output, err = client.ExecuteCommand("ps aux | grep -i vmm | head -3 || echo 'No VMM processes found'")
	if err != nil {
		fmt.Printf("Failed to check VMM processes: %v\n", err)
	} else {
		fmt.Printf("VMM processes:\n%s", output)
	}

	// Check installed packages
	output, err = client.ExecuteCommand("ls /var/packages/ | grep -i virtual || echo 'No virtual packages found'")
	if err != nil {
		fmt.Printf("Failed to check packages: %v\n", err)
	} else {
		fmt.Printf("Virtual packages:\n%s", output)
	}

	// Check system info
	output, err = client.ExecuteCommand("uname -a")
	if err != nil {
		fmt.Printf("Failed to get system info: %v\n", err)
	} else {
		fmt.Printf("System info: %s", output)
	}

	fmt.Println("Connection test completed successfully!")
}