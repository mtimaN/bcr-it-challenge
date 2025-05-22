package test

import (
	"context"
	"internal/db"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

func TestRedisRepo(t *testing.T) {
	// Setup miniredis for testing
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	// Create test config using miniredis
	config := &db.RedisConfig{
		Addr:         "localhost:6379",
		Password:     "RPass0319",
		DB:           0,
		PoolSize:     1,
		MinIdleConns: 1,
		DialTimeout:  1 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		MaxRetries:   1,
		Expiration:   1 * time.Hour,
		KeyPrefix:    "test:user:",
	}

	// Create repo
	repo, err := db.NewRedisRepo(config)
	if err != nil {
		t.Fatalf("Failed to create Redis repo: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()

	// Test Health
	if err := repo.Health(ctx); err != nil {
		t.Errorf("Health check failed: %v", err)
	}

	// Create test user
	testUser := &db.User{
		Credentials: &db.Credentials{
			Username: "testuser",
			Password: "password123",
		},
		Email:    "test@example.com",
		Category: 0,
	}

	// Test Add
	if err := repo.Add(ctx, testUser); err != nil {
		t.Errorf("Failed to add user: %v", err)
	}

	// Test Exists
	exists, err := repo.Exists(ctx, "testuser")
	if err != nil {
		t.Errorf("Failed to check if user exists: %v", err)
	}
	if !exists {
		t.Error("User should exist but doesn't")
	}

	// Test Get
	retrievedUser, err := repo.Get(ctx, "testuser")
	if err != nil {
		t.Errorf("Failed to get user: %v", err)
	}
	if retrievedUser.Username != testUser.Username || retrievedUser.Email != testUser.Email {
		t.Error("Retrieved user doesn't match original user")
	}

	// Test Extend
	if err := repo.Extend(ctx, "testuser"); err != nil {
		t.Errorf("Failed to extend user expiration: %v", err)
	}

	// Test Stats
	stats, err := repo.Stats(ctx)
	if err != nil {
		t.Errorf("Failed to get stats: %v", err)
	}
	if stats == nil {
		t.Error("Stats should not be nil")
	}

	// Test Delete
	if err := repo.Delete(ctx, "testuser"); err != nil {
		t.Errorf("Failed to delete user: %v", err)
	}

	// Verify Delete
	exists, err = repo.Exists(ctx, "testuser")
	if err != nil {
		t.Errorf("Failed to check if user exists: %v", err)
	}
	if exists {
		t.Error("User should not exist after deletion")
	}
}
