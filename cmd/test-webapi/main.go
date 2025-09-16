package main

import (
	"fmt"
	"os"

	"github.com/scttfrdmn/syno-vm/internal/synology"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: %s <host> <username> <password>\n", os.Args[0])
		os.Exit(1)
	}

	host := os.Args[1]
	username := os.Args[2]
	password := os.Args[3]

	fmt.Printf("Testing Web API connection to %s@%s...\n", username, host)

	// Create Web API client directly
	client := synology.NewWebAPIClient(host, username, password)

	// Test login
	err := client.Login()
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Login successful!")

	// Test API call to list VMs
	resp, err := client.CallAPI("SYNO.Virtualization.API.Guest", "list", "1", map[string]interface{}{})
	if err != nil {
		fmt.Printf("API call failed: %v\n", err)
		// Don't exit, just show the error
	} else {
		fmt.Printf("✅ API call successful! Response: %+v\n", resp)
	}

	// Logout
	err = client.Logout()
	if err != nil {
		fmt.Printf("Logout failed: %v\n", err)
	} else {
		fmt.Println("✅ Logout successful!")
	}
}