package db_test

import (
	"context"
	"testing"
	"time"

	"db"
)

func createTestUser() *db.User {
	user := db.NewUser("test", "testPw", "test232")
	user.Category = 0
	return user
}

func verifyUsersEqual(t *testing.T, got, want *db.User) {
	t.Helper()
	if got.Username != want.Username ||
		got.Password != want.Password ||
		got.Email != want.Email ||
		got.Category != want.Category {
		t.Fatalf("user mismatch: got %+v, want %+v", got, want)
	}
}

func testCassandraAddGetUpdateDelete(t *testing.T, ctx context.Context) {
	config := db.NewCassandraConfig("backend", "BPass0319", "cass_keyspace")
	client, err := db.NewCassandraClient(config)
	if err != nil {
		t.Fatalf("failed to create Cassandra client: %v", err)
	}
	defer client.Close()

	user := createTestUser()

	// Add user
	if err := client.AddUser(ctx, user); err != nil {
		t.Fatalf("add user failed: %v", err)
	}

	// Get and verify
	got, err := client.GetUser(ctx, user.Username)
	if err != nil {
		t.Fatalf("get user failed: %v", err)
	}
	verifyUsersEqual(t, got, user)

	// Update user
	user.Email = "newEmail"
	user.Category = 1
	if err := client.UpdateUser(ctx, user); err != nil {
		t.Fatalf("update user failed: %v", err)
	}

	got, err = client.GetUser(ctx, user.Username)
	if err != nil {
		t.Fatalf("get updated user failed: %v", err)
	}
	if got.Email != "newEmail" || got.Category != 1 {
		t.Fatalf("update verification failed: got %+v", got)
	}

	// Invalid get (non-existent user)
	if _, err := client.GetUser(ctx, "nonexistent"); err == nil {
		t.Fatal("expected error on get non-existent user, got none")
	}

	// Invalid update (non-existent user)
	nonExistentUser := db.NewUser("ghost", "password", "email")
	nonExistentUser.Category = 2
	if err := client.UpdateUser(ctx, nonExistentUser); err == nil {
		t.Fatal("expected error on update non-existent user, got none")
	}

	// Delete user
	if err := client.DeleteUser(ctx, user.Username); err != nil {
		t.Fatalf("delete user failed: %v", err)
	}

	// Verify deletion
	if _, err := client.GetUser(ctx, user.Username); err == nil {
		t.Fatal("user still exists after deletion")
	}
}

func testRedisBasicCRUD(t *testing.T, ctx context.Context) {
	config := db.NewRedisConfig("RPass0319")
	client, err := db.NewRedisClient(config)
	if err != nil {
		t.Fatalf("failed to create Redis client: %v", err)
	}
	defer client.Close()

	user := createTestUser()

	// Invalid get (non-existent user)
	if _, err := client.GetUser(ctx, "test123"); err == nil {
		t.Fatal("expected error on get non-existent user, got none")
	}

	// Add user
	if err := client.AddUser(ctx, user); err != nil {
		t.Fatalf("add user failed: %v", err)
	}

	// Get and verify user
	got, err := client.GetUser(ctx, user.Username)
	if err != nil {
		t.Fatalf("get user failed: %v", err)
	}
	verifyUsersEqual(t, got, user)
}

func testRedisExpiration(t *testing.T, ctx context.Context) {
	config := db.NewRedisConfig("RPass0319")
	config.Expiration = 10 * time.Second
	client, err := db.NewRedisClient(config)
	if err != nil {
		t.Fatalf("failed to create Redis client with expiration: %v", err)
	}
	defer client.Close()

	user := createTestUser()

	if err := client.AddUser(ctx, user); err != nil {
		t.Fatalf("add user failed: %v", err)
	}

	// Verify user exists before expiration
	if _, err := client.GetUser(ctx, user.Username); err != nil {
		t.Fatalf("user should exist before expiration: %v", err)
	}

	// Wait for expiration + buffer
	time.Sleep(config.Expiration + 2*time.Second)

	// Verify user is expired
	if _, err := client.GetUser(ctx, user.Username); err == nil {
		t.Fatal("user still exists after expiration")
	}
}

func TestCassandraCRUD(t *testing.T) {
	ctx := context.Background()
	testCassandraAddGetUpdateDelete(t, ctx)
}

func TestRedisCRUD(t *testing.T) {
	ctx := context.Background()
	testRedisBasicCRUD(t, ctx)
}

func TestRedisExpiration(t *testing.T) {
	ctx := context.Background()
	testRedisExpiration(t, ctx)
}
