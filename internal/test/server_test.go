package test

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// Configuration constants
const (
	defaultBaseURL     = "https://localhost:8443"
	defaultCertDir     = "../../certs"
	requestTimeout     = 30 * time.Second
	contentTypeJSON    = "application/json"
	testUserPrefix     = "apitest"
	testEmailDomain    = "example.com"
	testPasswordStrong = "TestPassword123!"
)

// API endpoints
const (
	endpointRegister = "/v1/register"
	endpointLogin    = "/v1/login"
	endpointUpdate   = "/v1/update"
	endpointDelete   = "/v1/delete"
	endpointStats    = "/v1/stats"
)

// Test payload structures
type TestPayload struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email,omitempty"`
	OldPassword string `json:"old_password,omitempty"`
	Category    int    `json:"category,omitempty"`
}

type APIResponse struct {
	Message string                 `json:"message,omitempty"`
	Error   string                 `json:"error,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

type StatsResponse struct {
	TotalUsers     int                    `json:"total_users"`
	ActiveSessions int                    `json:"active_sessions"`
	DatabaseStats  map[string]interface{} `json:"database_stats,omitempty"`
}

type TestReporter interface {
	Helper()
	Errorf(string, ...any)
	Fatalf(string, ...any)
	Fatal(...any)
	Logf(string, ...any)
}

// Test client management
type APITestClient struct {
	baseURL    string
	httpClient *http.Client
	t          TestReporter
}

func newAPITestClient(t TestReporter) *APITestClient {
	t.Helper()

	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	certDir := os.Getenv("TLS_CERT_DIR")
	if certDir == "" {
		certDir = defaultCertDir
	}

	client := &APITestClient{
		baseURL:    baseURL,
		httpClient: createSecureHTTPClient(t, certDir),
		t:          t,
	}

	return client
}

func createSecureHTTPClient(t TestReporter, certDir string) *http.Client {
	t.Helper()

	certPath := fmt.Sprintf("%s/server.crt", certDir)
	cert, err := os.ReadFile(certPath)
	if err != nil {
		t.Fatalf("failed to read certificate from %s: %v", certPath, err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(cert) {
		t.Fatal("failed to append certificate to CA pool")
	}

	return &http.Client{
		Timeout: requestTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:    caCertPool,
				MinVersion: tls.VersionTLS12, // Enforce minimum TLS version
			},
		},
	}
}

// Test data factory functions
func createTestPayload(suffix string) TestPayload {
	return TestPayload{
		Username: fmt.Sprintf("%s_%s", testUserPrefix, suffix),
		Password: testPasswordStrong,
		Email:    fmt.Sprintf("%s_%s@%s", testUserPrefix, suffix, testEmailDomain),
		Category: 0,
	}
}

func createInvalidPayloads() map[string]TestPayload {
	return map[string]TestPayload{
		"empty_username": {
			Username: "",
			Password: testPasswordStrong,
			Email:    "test@example.com",
		},
		"short_username": {
			Username: "ab",
			Password: testPasswordStrong,
			Email:    "test@example.com",
		},
		"weak_password": {
			Username: "testuser",
			Password: "123",
			Email:    "test@example.com",
		},
		"invalid_email": {
			Username: "testuser",
			Password: testPasswordStrong,
			Email:    "not-an-email",
		},
		"missing_email": {
			Username: "testuser",
			Password: testPasswordStrong,
			Email:    "",
		},
	}
}

// HTTP request helpers
func (c *APITestClient) makeRequest(method, endpoint string, payload interface{}) (*http.Response, error) {
	c.t.Helper()

	url := c.baseURL + endpoint

	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", contentTypeJSON)
	}
	req.Header.Set("User-Agent", "API-Test-Client/1.0")

	return c.httpClient.Do(req)
}

func (c *APITestClient) makeJSONRequest(method, endpoint string, payload interface{}) (*http.Response, error) {
	return c.makeRequest(method, endpoint, payload)
}

func (c *APITestClient) post(endpoint string, payload interface{}) (*http.Response, error) {
	return c.makeJSONRequest("POST", endpoint, payload)
}

func (c *APITestClient) get(endpoint string) (*http.Response, error) {
	return c.makeRequest("GET", endpoint, nil)
}

func (c *APITestClient) parseResponseBody(resp *http.Response) ([]byte, error) {
	c.t.Helper()

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil
}

func (c *APITestClient) parseJSONResponse(resp *http.Response, target interface{}) error {
	c.t.Helper()

	body, err := c.parseResponseBody(resp)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON response: %w, body: %s", err, string(body))
	}

	return nil
}

func (c *APITestClient) expectStatusCode(resp *http.Response, expected int) {
	c.t.Helper()

	if resp.StatusCode != expected {
		body, _ := c.parseResponseBody(resp)
		c.t.Fatalf("expected status %d, got %d. Response body: %s",
			expected, resp.StatusCode, string(body))
	}
}

func (c *APITestClient) expectStatusCodeIn(resp *http.Response, expected ...int) {
	c.t.Helper()

	for _, code := range expected {
		if resp.StatusCode == code {
			return
		}
	}

	body, _ := c.parseResponseBody(resp)
	c.t.Fatalf("expected status in %v, got %d. Response body: %s",
		expected, resp.StatusCode, string(body))
}

// Cleanup helper
func (c *APITestClient) cleanupUser(username string) {
	c.t.Helper()

	payload := TestPayload{
		Username: username,
		Password: testPasswordStrong, // Assuming default password
	}

	resp, err := c.post(endpointDelete, payload)
	if err != nil {
		c.t.Logf("cleanup failed for user %s: %v", username, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		c.t.Logf("cleanup returned unexpected status for user %s: %d", username, resp.StatusCode)
	}
}

// Test functions
func TestAPI_InvalidRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping API integration test in short mode")
	}

	client := newAPITestClient(t)

	t.Run("malformed_json", func(t *testing.T) {
		// Test malformed JSON syntax
		malformedJSON := `{"username": "test", "email": "bad@example.com", "password": "1234",}` // extra comma

		req, err := http.NewRequest("POST", client.baseURL+endpointRegister, strings.NewReader(malformedJSON))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", contentTypeJSON)

		resp, err := client.httpClient.Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			body, _ := client.parseResponseBody(resp)
			t.Fatalf("expected error for malformed JSON, got %d: %s", resp.StatusCode, string(body))
		}
	})

	t.Run("validation_errors", func(t *testing.T) {
		invalidPayloads := createInvalidPayloads()

		for name, payload := range invalidPayloads {
			t.Run(name, func(t *testing.T) {
				resp, err := client.post(endpointRegister, payload)
				if err != nil {
					t.Fatalf("request failed: %v", err)
				}
				defer resp.Body.Close()

				// Should return 4xx error for validation failures
				if resp.StatusCode < 400 || resp.StatusCode >= 500 {
					body, _ := client.parseResponseBody(resp)
					t.Errorf("expected 4xx validation error, got %d: %s", resp.StatusCode, string(body))
				}
			})
		}
	})

	t.Run("unsupported_content_type", func(t *testing.T) {
		req, err := http.NewRequest("POST", client.baseURL+endpointRegister, strings.NewReader("plain text"))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "text/plain")

		resp, err := client.httpClient.Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			t.Fatalf("expected error for unsupported content type, got %d", resp.StatusCode)
		}
	})

	t.Run("nonexistent_endpoint", func(t *testing.T) {
		resp, err := client.get("/v1/nonexistent")
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404 for nonexistent endpoint, got %d", resp.StatusCode)
		}
	})
}

func TestAPI_Stats(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping API integration test in short mode")
	}

	client := newAPITestClient(t)

	resp, err := client.get(endpointStats)
	if err != nil {
		t.Fatalf("stats request failed: %v", err)
	}
	defer resp.Body.Close()

	client.expectStatusCode(resp, http.StatusOK)

	var stats StatsResponse
	if err := client.parseJSONResponse(resp, &stats); err != nil {
		// Fallback to generic map if specific structure fails
		var genericStats map[string]interface{}
		body, _ := client.parseResponseBody(resp)
		if err := json.Unmarshal(body, &genericStats); err != nil {
			t.Fatalf("failed to parse stats response: %v, body: %s", err, string(body))
		}
		t.Logf("Stats (generic): %+v", genericStats)
		return
	}

	t.Logf("Stats: TotalUsers=%d, ActiveSessions=%d, DatabaseStats=%+v",
		stats.TotalUsers, stats.ActiveSessions, stats.DatabaseStats)

	// Basic validation
	if stats.TotalUsers < 0 {
		t.Errorf("total users should not be negative: %d", stats.TotalUsers)
	}
	if stats.ActiveSessions < 0 {
		t.Errorf("active sessions should not be negative: %d", stats.ActiveSessions)
	}
}

func TestAPI_UserLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping API integration test in short mode")
	}

	client := newAPITestClient(t)
	payload := createTestPayload("lifecycle")

	// Ensure cleanup
	t.Cleanup(func() {
		client.cleanupUser(payload.Username)
	})

	t.Run("register_user", func(t *testing.T) {
		resp, err := client.post(endpointRegister, payload)
		if err != nil {
			t.Fatalf("register request failed: %v", err)
		}
		defer resp.Body.Close()

		client.expectStatusCode(resp, http.StatusCreated)

		var response APIResponse
		if err := client.parseJSONResponse(resp, &response); err != nil {
			t.Logf("failed to parse register response as structured JSON: %v", err)
		} else {
			t.Logf("Register response: %s", response.Message)
		}
	})

	t.Run("duplicate_registration", func(t *testing.T) {
		// Try to register the same user again
		resp, err := client.post(endpointRegister, payload)
		if err != nil {
			t.Fatalf("duplicate register request failed: %v", err)
		}
		defer resp.Body.Close()

		// Should return conflict or bad request
		client.expectStatusCodeIn(resp, http.StatusConflict, http.StatusBadRequest)
	})

	t.Run("login_user", func(t *testing.T) {
		loginPayload := TestPayload{
			Username: payload.Username,
			Password: payload.Password,
		}

		resp, err := client.post(endpointLogin, loginPayload)
		if err != nil {
			t.Fatalf("login request failed: %v", err)
		}
		defer resp.Body.Close()

		client.expectStatusCode(resp, http.StatusOK)

		var response APIResponse
		if err := client.parseJSONResponse(resp, &response); err != nil {
			t.Logf("failed to parse login response as structured JSON: %v", err)
		}
	})

	t.Run("login_wrong_password", func(t *testing.T) {
		wrongLoginPayload := TestPayload{
			Username: payload.Username,
			Password: "WrongPassword123!",
		}

		resp, err := client.post(endpointLogin, wrongLoginPayload)
		if err != nil {
			t.Fatalf("wrong password login request failed: %v", err)
		}
		defer resp.Body.Close()

		// Should return unauthorized
		client.expectStatusCode(resp, http.StatusUnauthorized)
	})

	t.Run("update_user", func(t *testing.T) {
		updatedPayload := TestPayload{
			Username:    payload.Username,
			Password:    "NewPassword123!",
			Email:       fmt.Sprintf("updated_%s@%s", payload.Username, testEmailDomain),
			Category:    2,
			OldPassword: payload.Password,
		}

		resp, err := client.post(endpointUpdate, updatedPayload)
		if err != nil {
			t.Fatalf("update request failed: %v", err)
		}
		defer resp.Body.Close()

		client.expectStatusCode(resp, http.StatusOK)

		// Update our test payload for subsequent tests
		payload.Password = updatedPayload.Password
		payload.Email = updatedPayload.Email
		payload.Category = updatedPayload.Category

		// Verify we can login with new password
		loginPayload := TestPayload{
			Username: payload.Username,
			Password: payload.Password,
		}

		loginResp, err := client.post(endpointLogin, loginPayload)
		if err != nil {
			t.Fatalf("login with new password failed: %v", err)
		}
		defer loginResp.Body.Close()

		client.expectStatusCode(loginResp, http.StatusOK)
	})

	t.Run("update_wrong_old_password", func(t *testing.T) {
		wrongUpdatePayload := TestPayload{
			Username:    payload.Username,
			Password:    "AnotherPassword123!",
			Email:       "another@example.com",
			Category:    3,
			OldPassword: "WrongOldPassword123!",
		}

		resp, err := client.post(endpointUpdate, wrongUpdatePayload)
		if err != nil {
			t.Fatalf("update with wrong old password request failed: %v", err)
		}
		defer resp.Body.Close()

		// Should return unauthorized or bad request
		client.expectStatusCodeIn(resp, http.StatusUnauthorized, http.StatusBadRequest)
	})

	t.Run("delete_user", func(t *testing.T) {
		deletePayload := TestPayload{
			Username: payload.Username,
			Password: payload.Password,
		}

		resp, err := client.post(endpointDelete, deletePayload)
		if err != nil {
			t.Fatalf("delete request failed: %v", err)
		}
		defer resp.Body.Close()

		client.expectStatusCode(resp, http.StatusOK)
	})

	t.Run("login_deleted_user", func(t *testing.T) {
		loginPayload := TestPayload{
			Username: payload.Username,
			Password: payload.Password,
		}

		resp, err := client.post(endpointLogin, loginPayload)
		if err != nil {
			t.Fatalf("login deleted user request failed: %v", err)
		}
		defer resp.Body.Close()

		// Should return unauthorized or not found
		client.expectStatusCodeIn(resp, http.StatusUnauthorized, http.StatusNotFound)
	})
}

func TestAPI_ConcurrentUsers(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping API concurrent test in short mode")
	}

	client := newAPITestClient(t)
	const numUsers = 5

	// Create multiple users concurrently
	t.Run("concurrent_registration", func(t *testing.T) {
		results := make(chan error, numUsers)

		for i := 0; i < numUsers; i++ {
			go func(i int) {
				payload := createTestPayload(fmt.Sprintf("concurrent_%d", i))

				// Cleanup
				defer client.cleanupUser(payload.Username)

				resp, err := client.post(endpointRegister, payload)
				if err != nil {
					results <- fmt.Errorf("user %d registration failed: %w", i, err)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusCreated {
					body, _ := client.parseResponseBody(resp)
					results <- fmt.Errorf("user %d got status %d: %s", i, resp.StatusCode, string(body))
					return
				}

				results <- nil
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < numUsers; i++ {
			if err := <-results; err != nil {
				t.Errorf("concurrent registration error: %v", err)
			}
		}
	})
}

// Benchmark test
func BenchmarkAPI_RegisterUser(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping API benchmark in short mode")
	}

	client := &APITestClient{
		baseURL:    defaultBaseURL,
		httpClient: createSecureHTTPClient(b, defaultCertDir),
		t:          b,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		payload := createTestPayload(fmt.Sprintf("bench_%d", i))

		resp, err := client.post(endpointRegister, payload)
		if err != nil {
			b.Fatalf("register request failed: %v", err)
		}
		resp.Body.Close()

		// Cleanup
		client.cleanupUser(payload.Username)
	}
}
