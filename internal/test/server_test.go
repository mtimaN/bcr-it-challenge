package test

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	baseURL      = "https://localhost:8443/v1"
	testUsername = "testuser"
	testPassword = "testpassword"
	testEmail    = "test@example.com"
)

var (
	client *http.Client
	token  string
)

func TestMain(m *testing.M) {
	// Skip TLS verification for testing
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}

	// Give the server time to start if running tests concurrently
	time.Sleep(1 * time.Second)

	code := m.Run()
	// Clean up test user if needed
	cleanupTestUser()
	os.Exit(code)
}

func TestUserLifecycle(t *testing.T) {
	t.Run("Register new user", testRegister)
	t.Run("Login with created user", testLogin)
	t.Run("Get ads category", testGetAdsCategory)
	t.Run("Update user information", testUpdateUser)
	t.Run("Delete user", testDeleteUser)
}

func testRegister(t *testing.T) {
	payload := map[string]interface{}{
		"username": testUsername,
		"password": testPassword,
		"email":    testEmail,
	}

	resp, err := makeRequest("POST", "/register", payload, "")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, body)
	}
}

func testLogin(t *testing.T) {
	payload := map[string]interface{}{
		"username": testUsername,
		"password": testPassword,
	}

	resp, err := makeRequest("POST", "/login", payload, "")
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, body)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result["token"] == "" {
		t.Fatal("No token received in login response")
	}

	// Save token for later tests
	token = result["token"]
}

func testGetAdsCategory(t *testing.T) {
	if token == "" {
		t.Fatal("No auth token available")
	}

	resp, err := makeRequest("GET", "/get_ads", nil, token)
	if err != nil {
		t.Fatalf("Failed to get ads category: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, body)
	}

	var result map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result["category"] != 2 {
		t.Fatalf("Expected category 2, got %d", result["category"])
	}
}

func testUpdateUser(t *testing.T) {
	if token == "" {
		t.Fatal("No auth token available")
	}

	// Update email and category
	payload := map[string]interface{}{
		"password":     testPassword,
		"email":        "updated@example.com",
		"category":     2,
		"new_password": "newpassword",
	}

	resp, err := makeRequest("POST", "/update", payload, token)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, body)
	}

	// Login with new password to verify update
	loginPayload := map[string]interface{}{
		"username": testUsername,
		"password": "newpassword",
	}

	loginResp, err := makeRequest("POST", "/login", loginPayload, "")
	if err != nil {
		t.Fatalf("Failed to login with new password: %v", err)
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != http.StatusOK {
		t.Fatal("Failed to login with updated credentials")
	}

	// Update token
	var result map[string]string
	if err := json.NewDecoder(loginResp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	token = result["token"]

	// Check if category was updated
	catResp, err := makeRequest("GET", "/get_ads", nil, token)
	if err != nil {
		t.Fatalf("Failed to get updated category: %v", err)
	}
	defer catResp.Body.Close()

	var catResult map[string]int
	if err := json.NewDecoder(catResp.Body).Decode(&catResult); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if catResult["category"] != 2 {
		t.Fatalf("Expected updated category 2, got %d", catResult["category"])
	}
}

func testDeleteUser(t *testing.T) {
	if token == "" {
		t.Fatal("No auth token available")
	}

	resp, err := makeRequest("POST", "/delete", nil, token)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, body)
	}

	// Verify user is deleted by trying to login
	payload := map[string]interface{}{
		"username": testUsername,
		"password": "newpassword", // Using the updated password
	}

	loginResp, err := makeRequest("POST", "/login", payload, "")
	if err != nil {
		t.Fatalf("Error making login request: %v", err)
	}
	defer loginResp.Body.Close()

	// Should fail with unauthorized
	if loginResp.StatusCode >= 200 && loginResp.StatusCode < 300 {
		t.Fatalf("Expected error status for deleted user, got %d", loginResp.StatusCode)
	}
}

func TestInvalidRequests(t *testing.T) {
	t.Run("Register with existing username", testRegisterDuplicate)
	t.Run("Login with invalid credentials", testInvalidLogin)
	t.Run("Access protected endpoint without token", testUnauthorizedAccess)
}

func testRegisterDuplicate(t *testing.T) {
	// First, register a user
	payload := map[string]interface{}{
		"username": "dupluser",
		"password": "testpassword",
		"email":    "dupl@example.com",
		"category": 1,
	}

	resp, err := makeRequest("POST", "/register", payload, "")
	if err != nil {
		t.Fatalf("Failed to register user for duplicate test: %v", err)
	}
	resp.Body.Close()

	// Try to register again with same username
	resp, err = makeRequest("POST", "/register", payload, "")
	if err != nil {
		t.Fatalf("Failed to make duplicate registration request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for duplicate registration, got %d", resp.StatusCode)
	}

	// Clean up
	loginResp, err := makeRequest("POST", "/login", payload, "")
	if err == nil {
		var result map[string]string
		json.NewDecoder(loginResp.Body).Decode(&result)
		loginResp.Body.Close()

		if result["token"] != "" {
			deleteResp, _ := makeRequest("POST", "/delete", nil, result["token"])
			if deleteResp != nil {
				deleteResp.Body.Close()
			}
		}
	}
}

func testInvalidLogin(t *testing.T) {
	payload := map[string]interface{}{
		"username": "nonexistentuser",
		"password": "wrongpassword",
	}

	resp, err := makeRequest("POST", "/login", payload, "")
	if err != nil {
		t.Fatalf("Failed to make invalid login request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		t.Fatalf("Expected error status for invalid login, got %d", resp.StatusCode)
	}
}

func testUnauthorizedAccess(t *testing.T) {
	resp, err := makeRequest("GET", "/get_ads", nil, "")
	if err != nil {
		t.Fatalf("Failed to make unauthorized request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected status 401 for unauthorized access, got %d", resp.StatusCode)
	}
}

func TestServerStats(t *testing.T) {
	// No authentication needed for stats endpoint in this implementation
	resp, err := makeRequest("GET", "/stats", nil, "")
	if err != nil {
		t.Fatalf("Failed to get server stats: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 for stats, got %d", resp.StatusCode)
	}

	// Just check if we get valid JSON, not testing specific values
	var stats map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		t.Fatalf("Failed to decode stats response: %v", err)
	}

	printStats(stats)
}

// Helper functions

func printStats(stats map[string]interface{}) {
	fmt.Println("-- Redis Stats --")
	for k, v := range stats {
		if val, ok := v.(float64); ok {
			fmt.Printf("%s: %d\n", k, uint32(val))
		} else {
			fmt.Printf("%s: (invalid type)\n", k)
		}
	}
}

func makeRequest(method, endpoint string, payload interface{}, authToken string) (*http.Response, error) {
	var reqBody io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, baseURL+endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return client.Do(req)
}

func cleanupTestUser() {
	// Try to login and delete the test user if it exists
	payload := map[string]interface{}{
		"username": testUsername,
		"password": testPassword,
	}

	resp, err := makeRequest("POST", "/login", payload, "")
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			resp.Body.Close()
		}
		return
	}

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	if err != nil || result["token"] == "" {
		return
	}

	// Try with new password if the user was updated
	if result["token"] == "" {
		payload["password"] = "newpassword"
		resp, err = makeRequest("POST", "/login", payload, "")
		if err != nil || resp.StatusCode != http.StatusOK {
			if resp != nil {
				resp.Body.Close()
			}
			return
		}

		err = json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()
		if err != nil || result["token"] == "" {
			return
		}
	}

	// Delete the user
	deleteResp, _ := makeRequest("POST", "/delete", nil, result["token"])
	if deleteResp != nil {
		deleteResp.Body.Close()
	}
}
