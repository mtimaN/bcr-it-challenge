package test

import (
	"context"
	"internal/db"
	"testing"
)

func TestCassandraRepo(t *testing.T) {
	// Create test config using real Cassandra
	config := db.NewCassandraConfig("backend", "BPass0319", "cass_keyspace")

	// Create repo
	repo, err := db.NewCassandraRepo(config)
	if err != nil {
		t.Skipf("Skipping test: failed to connect to Cassandra: %v", err)
		return
	}
	defer repo.Close()

	ctx := context.Background()

	// Create test user
	testUser := &db.User{
		Credentials: &db.Credentials{
			Username: "testuser",
			Password: "password123",
		},
		Email:    "test@example.com",
		Category: 1,
	}

	// Clean up user if exists from previous tests
	_ = repo.DeleteUser(ctx, "testuser")

	// Test UsernameExists for non-existent user
	exists, err := repo.UsernameExists(ctx, "testuser")
	if err != nil {
		t.Errorf("UsernameExists failed: %v", err)
	}
	if exists {
		t.Error("User should not exist but does")
	}

	// Test GetUser for non-existent user
	_, err = repo.GetUser(ctx, "testuser")
	if err == nil {
		t.Error("Expected error for non-existent user but got none")
	}

	// Test AddUser
	if err := repo.AddUser(ctx, testUser); err != nil {
		t.Errorf("Failed to add user: %v", err)
	}

	// Test UsernameExists for existing user
	exists, err = repo.UsernameExists(ctx, "testuser")
	if err != nil || !exists {
		t.Error("User should exist after adding")
	}

	// Test GetUser for existing user
	user, err := repo.GetUser(ctx, "testuser")
	if err != nil {
		t.Errorf("GetUser failed: %v", err)
	}
	if user.Username != "testuser" || user.Email != "test@example.com" || user.Category != 1 {
		t.Error("Retrieved user doesn't match expected")
	}

	// Test UpdateUser
	testUser.Email = "updated@example.com"
	testUser.Category = 2
	if err := repo.UpdateUser(ctx, testUser); err != nil {
		t.Errorf("UpdateUser failed: %v", err)
	}

	// Verify update
	updatedUser, _ := repo.GetUser(ctx, "testuser")
	if updatedUser.Email != "updated@example.com" || updatedUser.Category != 2 {
		t.Error("User was not updated correctly")
	}

	// Test UpdateUser for non-existent user
	nonExistentUser := &db.User{
		Credentials: &db.Credentials{Username: "nonexistent", Password: "pass"},
		Email:       "none@example.com",
		Category:    1,
	}
	err = repo.UpdateUser(ctx, nonExistentUser)
	if err == nil {
		t.Errorf("Expected error when updating non-existent user")
	}

	// Test DeleteUser
	if err := repo.DeleteUser(ctx, "testuser"); err != nil {
		t.Errorf("DeleteUser failed: %v", err)
	}

	// Verify deleted
	exists, _ = repo.UsernameExists(ctx, "testuser")
	if exists {
		t.Error("User should not exist after deletion")
	}

	// Test DeleteUser for non-existent user
	err = repo.DeleteUser(ctx, "nonexistentuser")
	if err == nil {
		t.Error("Expected error when deleting non-existent user")
	}
}
