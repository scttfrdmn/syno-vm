package synology

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// WebAPIClient handles HTTP communication with Synology Web API
type WebAPIClient struct {
	baseURL    string
	httpClient *http.Client
	sessionID  string
	username   string
	password   string
}

// NewWebAPIClient creates a new Web API client
func NewWebAPIClient(host, username, password string) *WebAPIClient {
	baseURL := fmt.Sprintf("https://%s:5001", host)

	// Create HTTP client with SSL verification disabled (common for NAS devices)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	return &WebAPIClient{
		baseURL:    baseURL,
		httpClient: httpClient,
		username:   username,
		password:   password,
	}
}

// WebAPIResponse represents a standard Synology Web API response
type WebAPIResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   *WebAPIError           `json:"error,omitempty"`
}

// WebAPIError represents a Web API error
type WebAPIError struct {
	Code int `json:"code"`
}

// Login authenticates with the Synology Web API and obtains a session
func (w *WebAPIClient) Login() error {
	params := url.Values{}
	params.Set("api", "SYNO.API.Auth")
	params.Set("version", "3")
	params.Set("method", "login")
	params.Set("account", w.username)
	params.Set("passwd", w.password)
	params.Set("session", "VMM")
	params.Set("format", "cookie")

	resp, err := w.makeRequest("GET", "/webapi/auth.cgi", params, nil)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}

	var authResp WebAPIResponse
	if err := json.Unmarshal(resp, &authResp); err != nil {
		return fmt.Errorf("failed to parse login response: %w", err)
	}

	if !authResp.Success {
		code := 0
		if authResp.Error != nil {
			code = authResp.Error.Code
		}
		return fmt.Errorf("login failed with error code %d", code)
	}

	// Extract session ID from response data
	if authResp.Data != nil {
		if sid, ok := authResp.Data["sid"].(string); ok {
			w.sessionID = sid
		}
	}

	return nil
}

// Logout terminates the current session
func (w *WebAPIClient) Logout() error {
	if w.sessionID == "" {
		return nil // Already logged out
	}

	params := url.Values{}
	params.Set("api", "SYNO.API.Auth")
	params.Set("version", "3")
	params.Set("method", "logout")
	params.Set("session", "VMM")
	params.Set("_sid", w.sessionID)

	_, err := w.makeRequest("GET", "/webapi/auth.cgi", params, nil)
	w.sessionID = "" // Clear session regardless of result

	return err
}

// CallAPI makes an authenticated API call
func (w *WebAPIClient) CallAPI(api, method, version string, apiParams map[string]interface{}) (*WebAPIResponse, error) {
	if w.sessionID == "" {
		if err := w.Login(); err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
	}

	params := url.Values{}
	params.Set("api", api)
	params.Set("method", method)
	params.Set("version", version)
	params.Set("_sid", w.sessionID)

	// Add API-specific parameters
	for key, value := range apiParams {
		switch v := value.(type) {
		case string:
			params.Set(key, v)
		case int:
			params.Set(key, strconv.Itoa(v))
		case bool:
			params.Set(key, strconv.FormatBool(v))
		default:
			params.Set(key, fmt.Sprintf("%v", v))
		}
	}

	// Determine the endpoint based on the API
	endpoint := "/webapi/entry.cgi"
	if strings.Contains(api, "Virtualization") {
		// VMM APIs might use a different endpoint
		endpoint = "/webapi/entry.cgi"
	}

	resp, err := w.makeRequest("GET", endpoint, params, nil)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	var apiResp WebAPIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	// Handle session expiration
	if !apiResp.Success && apiResp.Error != nil && apiResp.Error.Code == 105 {
		// Session expired, try to re-login
		w.sessionID = ""
		if err := w.Login(); err != nil {
			return nil, fmt.Errorf("re-authentication failed: %w", err)
		}

		// Retry the API call with new session
		params.Set("_sid", w.sessionID)
		resp, err = w.makeRequest("GET", endpoint, params, nil)
		if err != nil {
			return nil, fmt.Errorf("API call retry failed: %w", err)
		}

		if err := json.Unmarshal(resp, &apiResp); err != nil {
			return nil, fmt.Errorf("failed to parse retry response: %w", err)
		}
	}

	return &apiResp, nil
}

// makeRequest performs an HTTP request
func (w *WebAPIClient) makeRequest(method, path string, params url.Values, body []byte) ([]byte, error) {
	var req *http.Request
	var err error

	fullURL := w.baseURL + path

	if method == "GET" && params != nil {
		fullURL += "?" + params.Encode()
		req, err = http.NewRequest(method, fullURL, nil)
	} else {
		var bodyReader io.Reader
		if body != nil {
			bodyReader = bytes.NewReader(body)
		} else if params != nil {
			bodyReader = strings.NewReader(params.Encode())
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		req, err = http.NewRequest(method, fullURL, bodyReader)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "syno-vm/0.1.0")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}