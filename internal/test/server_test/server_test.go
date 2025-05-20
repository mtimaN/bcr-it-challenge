package server_test

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
)

const baseURL = "https://localhost:8443"

type TestPayload struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email,omitempty"`
	OldPassword string `json:"old_password,omitempty"`
	Category    int    `json:"category,omitempty"`
}

var secureClient = func() *http.Client {
	certDir := os.Getenv("tls_cert_dir")
	if certDir == "" {
		certDir = "../../../certs"
	}

	cert, err := os.ReadFile(certDir + "/server.crt")
	if err != nil {
		panic(err)
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(cert); !ok {
		panic(err)
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
}()

func makeRequest(t *testing.T, endpoint string, payload TestPayload) *http.Response {
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	resp, err := secureClient.Post(baseURL+endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}

	return resp
}

func parseBody(t *testing.T, resp *http.Response) string {
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body failed: %v", err)
	}
	return string(data)
}

func TestInvalidRequests(t *testing.T) {
	// Invalid JSON (malformed syntax)
	t.Run("InvalidJSON", func(t *testing.T) {
		body := []byte(`{"user": { "username": "test", "email": "bad@example.com", "password": "1234", }}`) // extra comma
		resp, err := secureClient.Post(baseURL+"/v1/register", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("POST request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			t.Fatalf("expected failure on malformed JSON, got %d", resp.StatusCode)
		}
	})

	// Missing required fields
	t.Run("MissingFields", func(t *testing.T) {
		resp := makeRequest(t, "/v1/register", TestPayload{
			Username: "",
			Password: "testpass123",
			Email:    "email@example.com",
		})
		body := parseBody(t, resp)
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			t.Fatalf("expected validation error, got %d - %s", resp.StatusCode, body)
		}
	})

	// Invalid email format
	t.Run("InvalidEmail", func(t *testing.T) {
		resp := makeRequest(t, "/v1/register", TestPayload{
			Username: "bademail",
			Password: "pass1234",
			Email:    "not-an-email",
		})
		body := parseBody(t, resp)
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			t.Fatalf("expected email validation failure, got %d - %s", resp.StatusCode, body)
		}
	})
}

func TestStats(t *testing.T) {
	resp, err := secureClient.Get(baseURL + "/v1/stats")
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	body := parseBody(t, resp)

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		t.Fatalf("unmarshal failed: %v", body)
	}

	t.Logf("Stats: %v", result)
}

func TestAPI(t *testing.T) {
	payload := TestPayload{
		Username: "apitest",
		Password: "testpass123",
		Email:    "apitest@example.com",
	}

	// --- Add User
	t.Run("AddUser", func(t *testing.T) {
		resp := makeRequest(t, "/v1/register", payload)
		body := parseBody(t, resp)
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("add user failed: %d - %s", resp.StatusCode, body)
		}
	})

	// --- Get User
	t.Run("GetUser", func(t *testing.T) {
		resp := makeRequest(t, "/v1/login", payload)
		body := parseBody(t, resp)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("get user failed: %d - %s", resp.StatusCode, body)
		}
	})

	// --- Update User
	t.Run("UpdateUser", func(t *testing.T) {
		updated := TestPayload{
			Username:    payload.Username,
			Password:    "newpass123",
			Email:       "updated@example.com",
			Category:    2,
			OldPassword: payload.Password,
		}

		resp := makeRequest(t, "/v1/update", updated)
		body := parseBody(t, resp)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("update user failed: %d - %s", resp.StatusCode, body)
		}

		payload = updated
	})

	// --- Delete User
	t.Run("DeleteUser", func(t *testing.T) {
		resp := makeRequest(t, "/v1/delete", payload)
		body := parseBody(t, resp)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("delete user failed: %d - %s", resp.StatusCode, body)
		}
	})

	// --- Get Deleted User
	t.Run("GetDeletedUser", func(t *testing.T) {
		resp := makeRequest(t, "/v1/login", payload)
		body := parseBody(t, resp)
		if resp.StatusCode == http.StatusOK {
			t.Fatalf("expected failure for deleted user, got 200 - %s", body)
		}
	})
}
