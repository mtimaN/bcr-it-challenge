package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080"

type TestUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Category int    `json:"category"`
}

type TestPayload struct {
	User        TestUser `json:"user"`
	OldPassword string   `json:"old_password,omitempty"`
}

func makeRequest(t *testing.T, endpoint string, payload TestPayload) *http.Response {
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post(baseURL+endpoint, "application/json", bytes.NewReader(body))
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
		resp, err := http.Post(baseURL+"/v1/add", "application/json", bytes.NewReader(body))
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
		user := TestUser{
			Username: "", // required field is empty
			Password: "testpass",
			Email:    "email@example.com",
			Category: 1,
		}
		resp := makeRequest(t, "/v1/add", TestPayload{User: user})
		body := parseBody(t, resp)
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			t.Fatalf("expected validation error, got %d - %s", resp.StatusCode, body)
		}
	})

	// Invalid email format
	t.Run("InvalidEmail", func(t *testing.T) {
		user := TestUser{
			Username: "bademail",
			Password: "pass1234",
			Email:    "not-an-email",
			Category: 1,
		}
		resp := makeRequest(t, "/v1/add", TestPayload{User: user})
		body := parseBody(t, resp)
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			t.Fatalf("expected email validation failure, got %d - %s", resp.StatusCode, body)
		}
	})
}

func TestAPI(t *testing.T) {
	user := TestUser{
		Username: "apitest",
		Password: "testpass",
		Email:    "apitest@example.com",
		Category: 1,
	}

	// --- Add User
	t.Run("AddUser", func(t *testing.T) {
		resp := makeRequest(t, "/v1/add", TestPayload{User: user})
		body := parseBody(t, resp)
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("add user failed: %d - %s", resp.StatusCode, body)
		}
	})

	// --- Get User
	t.Run("GetUser", func(t *testing.T) {
		resp := makeRequest(t, "/v1/get", TestPayload{User: user})
		body := parseBody(t, resp)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("get user failed: %d - %s", resp.StatusCode, body)
		}
	})

	// --- Update User
	t.Run("UpdateUser", func(t *testing.T) {
		updated := user
		updated.Email = "updated@example.com"
		updated.Password = "newpass"
		updated.Category = 2

		resp := makeRequest(t, "/v1/update", TestPayload{
			User:        updated,
			OldPassword: user.Password,
		})
		body := parseBody(t, resp)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("update user failed: %d - %s", resp.StatusCode, body)
		}

		// use updated credentials for remaining tests
		user = updated
	})

	// --- Delete User
	t.Run("DeleteUser", func(t *testing.T) {
		resp := makeRequest(t, "/v1/delete", TestPayload{User: user})
		body := parseBody(t, resp)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("delete user failed: %d - %s", resp.StatusCode, body)
		}
	})

	// --- Get Deleted User
	t.Run("GetDeletedUser", func(t *testing.T) {
		resp := makeRequest(t, "/v1/get", TestPayload{User: user})
		body := parseBody(t, resp)
		if resp.StatusCode == http.StatusOK {
			t.Fatalf("expected failure for deleted user, got 200 - %s", body)
		}
	})
}
