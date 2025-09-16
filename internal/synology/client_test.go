package synology

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name          string
		host          string
		username      string
		expectedError bool
	}{
		{
			name:          "valid configuration",
			host:          "test-host",
			username:      "admin",
			expectedError: false,
		},
		{
			name:          "missing host",
			host:          "",
			username:      "admin",
			expectedError: true,
		},
		{
			name:          "missing username",
			host:          "test-host",
			username:      "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear viper settings
			viper.Reset()

			// Set test values
			if tt.host != "" {
				viper.Set("host", tt.host)
			}
			if tt.username != "" {
				viper.Set("username", tt.username)
			}
			viper.Set("port", 22)
			viper.Set("timeout", 30)

			client, err := NewClient()

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if client == nil {
					t.Errorf("expected client but got nil")
				} else {
					if client.host != tt.host {
						t.Errorf("expected host %s, got %s", tt.host, client.host)
					}
					if client.username != tt.username {
						t.Errorf("expected username %s, got %s", tt.username, client.username)
					}
				}
			}
		})
	}
}

func TestBuildAPICommand(t *testing.T) {
	tests := []struct {
		name     string
		api      string
		method   string
		version  string
		params   map[string]string
		expected string
	}{
		{
			name:     "basic command",
			api:      "SYNO.Virtualization.API.Guest.Action",
			method:   "poweron",
			version:  "1",
			params:   nil,
			expected: "synowebapi --exec api=SYNO.Virtualization.API.Guest.Action method=poweron version=1",
		},
		{
			name:    "command with parameters",
			api:     "SYNO.Virtualization.API.Guest.Action",
			method:  "poweron",
			version: "1",
			params: map[string]string{
				"runner":     "admin",
				"guest_name": "test-vm",
			},
			expected: "synowebapi --exec api=SYNO.Virtualization.API.Guest.Action method=poweron version=1 runner=admin guest_name=test-vm",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildAPICommand(tt.api, tt.method, tt.version, tt.params)

			// Since map iteration order is not guaranteed, we check if all expected parts are present
			if !strings.Contains(result, tt.api) {
				t.Errorf("result doesn't contain API: %s", result)
			}
			if !strings.Contains(result, "method="+tt.method) {
				t.Errorf("result doesn't contain method: %s", result)
			}
			if !strings.Contains(result, "version="+tt.version) {
				t.Errorf("result doesn't contain version: %s", result)
			}

			for key, value := range tt.params {
				param := key + "=" + value
				if !strings.Contains(result, param) {
					t.Errorf("result doesn't contain parameter %s: %s", param, result)
				}
			}
		})
	}
}


// Helper function to build API command (this would be extracted to the main client code)
func buildAPICommand(api, method, version string, params map[string]string) string {
	cmd := "synowebapi --exec api=" + api + " method=" + method + " version=" + version

	for key, value := range params {
		cmd += " " + key + "=" + value
	}

	return cmd
}

func TestVMConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    VMConfig
		wantError bool
	}{
		{
			name: "valid config",
			config: VMConfig{
				Name:   "test-vm",
				CPU:    2,
				Memory: 2048,
			},
			wantError: false,
		},
		{
			name: "missing name",
			config: VMConfig{
				CPU:    2,
				Memory: 2048,
			},
			wantError: true,
		},
		{
			name: "invalid CPU",
			config: VMConfig{
				Name:   "test-vm",
				CPU:    0,
				Memory: 2048,
			},
			wantError: true,
		},
		{
			name: "invalid memory",
			config: VMConfig{
				Name:   "test-vm",
				CPU:    2,
				Memory: 0,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("VMConfig.Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}