package test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"internal/db"
)

// Test configuration constants
const (
	testTimeout       = 30 * time.Second
	defaultExpiration = 10 * time.Second
	expirationBuffer  = 2 * time.Second
)

// Test data factory functions
func createTestUser(suffix ...string) *db.User {
	id := ""
	if len(suffix) > 0 {
		id = suffix[0]
	}
	user := db.NewUser(
		fmt.Sprintf("test%s", id),
		"testStrongPw123",
		fmt.Sprintf("test%s@gmail.com", id),
	)
	user.Category = 0
	return user
}

func createTestCredentials(suffix ...string) *db.Credentials {
	id := ""
	if len(suffix) > 0 {
		id = suffix[0]
	}
	return db.NewCredentials(fmt.Sprintf("test%s", id), "testStrongPw123")
}

// Test helper functions
func verifyUsersEqual(t *testing.T, got, want *db.User) {
	t.Helper()

	if got.Username != want.Username {
		t.Errorf("username mismatch: got %q, want %q", got.Username, want.Username)
	}
	if got.Password != want.Password {
		t.Errorf("password mismatch: got %q, want %q", got.Password, want.Password)
	}
	if got.Email != want.Email {
		t.Errorf("email mismatch: got %q, want %q", got.Email, want.Email)
	}
	if got.Category != want.Category {
		t.Errorf("category mismatch: got %d, want %d", got.Category, want.Category)
	}
}

func setupCassandraClient(t *testing.T) *db.CassandraRepo {
	t.Helper()

	config := db.NewCassandraConfig("backend", "BPass0319", "cass_keyspace")
	client, err := db.NewCassandraRepo(config)
	if err != nil {
		t.Fatalf("failed to create Cassandra client: %v", err)
	}

	t.Cleanup(func() {
		client.Close()
	})

	return client
}

func setupRedisClient(t *testing.T, expiration ...time.Duration) *db.RedisRepo {
	t.Helper()

	config := db.NewRedisConfig("RPass0319")
	if len(expiration) > 0 {
		config.Expiration = expiration[0]
	}

	client, err := db.NewRedisRepo(config)
	if err != nil {
		t.Fatalf("failed to create Redis client: %v", err)
	}

	t.Cleanup(func() {
		client.Close()
	})

	return client
}

func createContextWithTimeout(t *testing.T) context.Context {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	t.Cleanup(cancel)

	return ctx
}

// Cassandra tests
func TestCassandra_UserLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Cassandra integration test in short mode")
	}

	ctx := createContextWithTimeout(t)
	client := setupCassandraClient(t)

	user := createTestUser("_lifecycle")
	cred := createTestCredentials("_lifecycle")
	originalPassword := user.Password

	// Cleanup any existing test data
	t.Cleanup(func() {
		client.DeleteUser(ctx, user.Username)
	})

	t.Run("operations_on_nonexistent_user", func(t *testing.T) {
		// Update non-existent user should fail
		err := client.UpdateUser(ctx, user, user.Password)
		if err == nil {
			t.Error("expected error when updating non-existent user")
		}

		// Username should not exist
		exists, err := client.UsernameExists(ctx, user.Username)
		if err != nil {
			t.Fatalf("username exists check failed: %v", err)
		}
		if exists {
			t.Error("non-existent username should not exist")
		}

		// Get non-existent user should fail
		_, err = client.GetUser(ctx, cred)
		if err == nil {
			t.Error("expected error when getting non-existent user")
		}
	})

	t.Run("create_user", func(t *testing.T) {
		err := client.AddUser(ctx, user)
		if err != nil {
			t.Fatalf("failed to add user: %v", err)
		}

		// Verify username now exists
		exists, err := client.UsernameExists(ctx, user.Username)
		if err != nil {
			t.Fatalf("username exists check failed: %v", err)
		}
		if !exists {
			t.Error("username should exist after creation")
		}
	})

	t.Run("prevent_duplicate_user", func(t *testing.T) {
		duplicateUser := createTestUser("_lifecycle")
		err := client.AddUser(ctx, duplicateUser)
		if err == nil {
			t.Error("expected error when adding duplicate user")
		}
	})

	t.Run("retrieve_and_verify_user", func(t *testing.T) {
		got, err := client.GetUser(ctx, cred)
		if err != nil {
			t.Fatalf("failed to get user: %v", err)
		}

		verifyUsersEqual(t, got, user)

		// Verify password is properly hashed
		if got.Password != originalPassword {
			t.Error("returned password does not match")
		}
	})

	t.Run("authentication_failures", func(t *testing.T) {
		// Wrong password
		wrongCred := createTestCredentials("_lifecycle")
		wrongCred.Password = "WrongPassword123"
		_, err := client.GetUser(ctx, wrongCred)
		if err == nil {
			t.Error("authentication should fail with wrong password")
		}

		// Non-existent user
		wrongCred = createTestCredentials("_nonexistent")
		_, err = client.GetUser(ctx, wrongCred)
		if err == nil {
			t.Error("authentication should fail for non-existent user")
		}
	})

	t.Run("update_user", func(t *testing.T) {
		newPassword := "newPassword123"
		updatedUser := &db.User{
			Credentials: db.NewCredentials(user.Username, newPassword),
			Email:       "newEmail@test.com",
			Category:    1,
		}

		// Update with correct old password
		err := client.UpdateUser(ctx, updatedUser, originalPassword)
		if err != nil {
			t.Fatalf("failed to update user: %v", err)
		}

		// Verify update
		updatedCred := db.NewCredentials(updatedUser.Username, newPassword)
		got, err := client.GetUser(ctx, updatedCred)
		if err != nil {
			t.Fatalf("failed to get updated user: %v", err)
		}

		if got.Email != "newEmail@test.com" {
			t.Errorf("email not updated: got %q, want %q", got.Email, "newEmail@test.com")
		}
		if got.Category != 1 {
			t.Errorf("category not updated: got %d, want %d", got.Category, 1)
		}

		// Old password should no longer work
		_, err = client.GetUser(ctx, cred)
		if err == nil {
			t.Error("old password should not work after update")
		}

		// Update with wrong old password should fail
		anotherUpdate := &db.User{
			Credentials: db.NewCredentials(user.Username, "anotherPassword"),
			Email:       "another@test.com",
			Category:    2,
		}
		err = client.UpdateUser(ctx, anotherUpdate, "wrongOldPassword")
		if err == nil {
			t.Error("update should fail with wrong old password")
		}
	})

	t.Run("delete_user", func(t *testing.T) {
		err := client.DeleteUser(ctx, user.Username)
		if err != nil {
			t.Fatalf("failed to delete user: %v", err)
		}

		// Verify deletion
		_, err = client.GetUser(ctx, cred)
		if err == nil {
			t.Error("user should not exist after deletion")
		}

		// Username should not exist
		exists, err := client.UsernameExists(ctx, user.Username)
		if err != nil {
			t.Fatalf("username exists check failed: %v", err)
		}
		if exists {
			t.Error("username should not exist after deletion")
		}
	})
}

func TestCassandra_InputValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Cassandra integration test in short mode")
	}

	ctx := createContextWithTimeout(t)
	client := setupCassandraClient(t)

	testCases := []struct {
		name        string
		user        *db.User
		shouldError bool
		description string
	}{
		{
			name:        "empty_username",
			user:        &db.User{Credentials: db.NewCredentials("", "password123"), Email: "test@example.com", Category: 0},
			shouldError: true,
			description: "empty username should be rejected",
		},
		{
			name:        "short_username",
			user:        &db.User{Credentials: db.NewCredentials("ab", "password123"), Email: "test@example.com", Category: 0},
			shouldError: true,
			description: "username too short should be rejected",
		},
		{
			name:        "long_username",
			user:        &db.User{Credentials: db.NewCredentials(strings.Repeat("a", 100), "password123"), Email: "test@example.com", Category: 0},
			shouldError: true,
			description: "username too long should be rejected",
		},
		{
			name:        "invalid_username_chars",
			user:        &db.User{Credentials: db.NewCredentials("user@name", "password123"), Email: "test@example.com", Category: 0},
			shouldError: true,
			description: "username with invalid characters should be rejected",
		},
		{
			name:        "short_password",
			user:        &db.User{Credentials: db.NewCredentials("testuser", "short"), Email: "test@example.com", Category: 0},
			shouldError: true,
			description: "password too short should be rejected",
		},
		{
			name:        "invalid_email",
			user:        &db.User{Credentials: db.NewCredentials("testuser", "password123"), Email: "invalid-email", Category: 0},
			shouldError: true,
			description: "invalid email format should be rejected",
		},
		{
			name:        "negative_category",
			user:        &db.User{Credentials: db.NewCredentials("testuser1", "password123"), Email: "test@example.com", Category: -1},
			shouldError: false,
			description: "negative category should be allowed",
		},
		{
			name:        "valid_user",
			user:        &db.User{Credentials: db.NewCredentials("validuser", "password123"), Email: "valid@example.com", Category: 2},
			shouldError: false,
			description: "valid user should be accepted",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Ensure cleanup regardless of test outcome
			defer func() {
				if tc.user != nil && tc.user.Username != "" {
					client.DeleteUser(ctx, tc.user.Username)
				}
			}()

			err := client.AddUser(ctx, tc.user)

			if tc.shouldError && err == nil {
				t.Errorf("%s: expected error but got none", tc.description)
			}
			if !tc.shouldError && err != nil {
				t.Errorf("%s: unexpected error: %v", tc.description, err)
			}
		})
	}
}

// Redis tests
func TestRedis_UserLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Redis integration test in short mode")
	}

	ctx := createContextWithTimeout(t)
	client := setupRedisClient(t)

	user := createTestUser("_redis")
	cred := createTestCredentials("_redis")

	// Cleanup any existing test data
	t.Cleanup(func() {
		client.Delete(ctx, user.Username)
	})

	t.Run("get_nonexistent_user", func(t *testing.T) {
		_, err := client.Get(ctx, cred)
		if err == nil {
			t.Error("expected error when getting non-existent user")
		}
	})

	t.Run("add_and_get_user", func(t *testing.T) {
		err := client.Add(ctx, user)
		if err != nil {
			t.Fatalf("failed to add user: %v", err)
		}

		got, err := client.Get(ctx, cred)
		if err != nil {
			t.Fatalf("failed to get user: %v", err)
		}

		verifyUsersEqual(t, got, user)
	})

	t.Run("authentication_with_wrong_password", func(t *testing.T) {
		wrongCred := createTestCredentials("_redis")
		wrongCred.Password = "some_invalid_password232"

		_, err := client.Get(ctx, wrongCred)
		if err == nil {
			t.Error("authentication should fail with wrong password")
		}
	})

	t.Run("delete_user", func(t *testing.T) {
		err := client.Delete(ctx, user.Username)
		if err != nil {
			t.Fatalf("failed to delete user: %v", err)
		}

		_, err = client.Get(ctx, cred)
		if err == nil {
			t.Error("user should not exist after deletion")
		}
	})
}

func TestRedis_Expiration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Redis expiration test in short mode")
	}

	ctx := createContextWithTimeout(t)
	client := setupRedisClient(t, defaultExpiration)

	user := createTestUser("_expiry")
	cred := createTestCredentials("_expiry")

	t.Cleanup(func() {
		client.Delete(ctx, user.Username)
	})

	// Add user
	err := client.Add(ctx, user)
	if err != nil {
		t.Fatalf("failed to add user: %v", err)
	}

	// Verify user exists before expiration
	_, err = client.Get(ctx, cred)
	if err != nil {
		t.Fatalf("user should exist before expiration: %v", err)
	}

	// Wait for expiration
	time.Sleep(defaultExpiration + expirationBuffer)

	// Verify user is expired
	_, err = client.Get(ctx, cred)
	if err == nil {
		t.Error("user should be expired and not accessible")
	}
}

// Benchmark tests
func BenchmarkCassandra_AddUser(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark in short mode")
	}

	config := db.NewCassandraConfig("backend", "BPass0319", "cass_keyspace")
	client, err := db.NewCassandraRepo(config)
	if err != nil {
		b.Fatalf("failed to create Cassandra client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := createTestUser(fmt.Sprintf("_bench_%d", i))
		err := client.AddUser(ctx, user)
		if err != nil {
			b.Fatalf("failed to add user: %v", err)
		}

		// Cleanup
		client.DeleteUser(ctx, user.Username)
	}
}

func BenchmarkRedis_AddUser(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark in short mode")
	}

	config := db.NewRedisConfig("RPass0319")
	client, err := db.NewRedisRepo(config)
	if err != nil {
		b.Fatalf("failed to create Redis client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := createTestUser(fmt.Sprintf("_bench_%d", i))
		err := client.Add(ctx, user)
		if err != nil {
			b.Fatalf("failed to add user: %v", err)
		}

		// Cleanup
		client.Delete(ctx, user.Username)
	}
}
