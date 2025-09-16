package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestConfigOperations(t *testing.T) {
	// Create temporary config directory for testing
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".syno-vm")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp config dir: %v", err)
	}

	// Set up viper to use temp directory
	viper.Reset()
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	tests := []struct {
		name     string
		key      string
		value    string
		wantErr  bool
	}{
		{
			name:    "set host",
			key:     "host",
			value:   "test-host.local",
			wantErr: false,
		},
		{
			name:    "set username",
			key:     "username",
			value:   "admin",
			wantErr: false,
		},
		{
			name:    "set port",
			key:     "port",
			value:   "22",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(tt.key, tt.value)

			// Verify the value was set
			if viper.GetString(tt.key) != tt.value {
				t.Errorf("Expected %s to be %s, got %s", tt.key, tt.value, viper.GetString(tt.key))
			}
		})
	}
}

func TestExpandTildePath(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		shouldExpand   bool
	}{
		{
			name:         "tilde path",
			path:         "~/.ssh/id_rsa",
			shouldExpand: true, // We'll check if it's expanded, not the exact path
		},
		{
			name:         "absolute path",
			path:         "/home/user/.ssh/id_rsa",
			shouldExpand: false, // Should remain unchanged
		},
		{
			name:         "relative path",
			path:         ".ssh/id_rsa",
			shouldExpand: false, // Should remain unchanged
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := expandTildePath(tt.path)

			if err != nil {
				t.Errorf("expandTildePath() error = %v", err)
				return
			}

			if tt.shouldExpand {
				// For tilde paths, result should be different and not start with ~
				if result == tt.path || len(result) > 0 && result[0] == '~' {
					t.Errorf("expandTildePath() failed to expand tilde: got %s", result)
				}
			} else {
				// For non-tilde paths, result should be unchanged
				if result != tt.path {
					t.Errorf("expandTildePath() changed non-tilde path: got %s, want %s", result, tt.path)
				}
			}
		})
	}
}

// Helper function for expanding tilde paths (would be added to config.go)
func expandTildePath(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	if len(path) > 1 && path[1] != '/' {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path, err
	}

	return filepath.Join(home, path[2:]), nil
}