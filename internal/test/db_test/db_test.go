package db_test

import (
	"context"
	"testing"
	"time"

	"db"
)

func createTestUser() *db.User {
	user := db.NewUser("test", "testStrongPw123", "test232@gmail.com")
	user.Category = 0
	return user
}

func verifyUsersEqual(t *testing.T, got, want *db.User) {
	t.Helper()
	if got.Username != want.Username ||
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
	originalPassword := user.Password // Store original password for later use

	// Add user
	if err := client.AddUser(ctx, user); err != nil {
		t.Fatalf("add user failed: %v", err)
	}

	// Test duplicate user creation (should fail)
	duplicateUser := createTestUser() // Same username
	if err := client.AddUser(ctx, duplicateUser); err == nil {
		t.Fatal("expected error when adding duplicate user, got none")
	}

	// Get and verify user (password should be hashed now)
	got, err := client.GetUser(ctx, user.Username)
	if err != nil {
		t.Fatalf("get user failed: %v", err)
	}

	// Verify user data (except password which is now hashed)
	verifyUsersEqual(t, got, user)

	// Verify password is hashed (should not match original)
	if got.Password == originalPassword {
		t.Error("password should be hashed, but matches original")
	}
	if len(got.Password) < 20 { // bcrypt hashes are typically 60+ chars
		t.Error("password doesn't appear to be properly hashed")
	}

	// Test authentication with new method
	authUser, err := client.AuthenticateUser(ctx, user.Username, originalPassword)
	if err != nil {
		t.Fatalf("authentication failed: %v", err)
	}
	if authUser.Password != "" {
		t.Error("authenticated user should not return password hash")
	}

	// Test authentication with wrong password
	if _, err := client.AuthenticateUser(ctx, user.Username, "wrongpassword"); err == nil {
		t.Fatal("expected authentication to fail with wrong password")
	}

	// Test authentication with non-existent user
	if _, err := client.AuthenticateUser(ctx, "nonexistent", "password"); err == nil {
		t.Fatal("expected authentication to fail for non-existent user")
	}

	// Update user
	newPassword := "newPassword123"
	updatedUser := &db.User{
		Username: user.Username,
		Password: newPassword,
		Email:    "newEmail@test.com",
		Category: 1,
	}

	// Update should work with correct old password
	if err := client.UpdateUser(ctx, updatedUser, originalPassword); err != nil {
		t.Fatalf("update user failed: %v", err)
	}

	// Verify update
	got, err = client.GetUser(ctx, user.Username)
	if err != nil {
		t.Fatalf("get updated user failed: %v", err)
	}
	if got.Email != "newEmail@test.com" || got.Category != 1 {
		t.Fatalf("update verification failed: got email=%s category=%d", got.Email, got.Category)
	}

	// Test authentication with new password
	if _, err := client.AuthenticateUser(ctx, user.Username, newPassword); err != nil {
		t.Fatalf("authentication with new password failed: %v", err)
	}

	// Test authentication with old password (should fail)
	if _, err := client.AuthenticateUser(ctx, user.Username, originalPassword); err == nil {
		t.Fatal("authentication with old password should fail after update")
	}

	// Test update with wrong old password (should fail)
	anotherUpdate := &db.User{
		Username: user.Username,
		Password: "anotherPassword",
		Email:    "another@test.com",
		Category: 2,
	}
	if err := client.UpdateUser(ctx, anotherUpdate, "wrongOldPassword"); err == nil {
		t.Fatal("expected error when updating with wrong old password")
	}

	// Test update non-existent user
	nonExistentUser := &db.User{
		Username: "ghost",
		Password: "password123",
		Email:    "ghost@test.com",
		Category: 2,
	}
	if err := client.UpdateUser(ctx, nonExistentUser, "password123"); err == nil {
		t.Fatal("expected error on update non-existent user, got none")
	}

	// Test username existence check
	exists, err := client.UsernameExists(ctx, user.Username)
	if err != nil {
		t.Fatalf("username exists check failed: %v", err)
	}
	if !exists {
		t.Error("username should exist")
	}

	exists, err = client.UsernameExists(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("username exists check for non-existent failed: %v", err)
	}
	if exists {
		t.Error("nonexistent username should not exist")
	}

	// Test invalid get (non-existent user)
	if _, err := client.GetUser(ctx, "nonexistent"); err == nil {
		t.Fatal("expected error on get non-existent user, got none")
	}

	// Delete user with correct password
	deleteUser := &db.User{
		Username: user.Username,
		Password: newPassword, // Use current password
		Email:    got.Email,
		Category: got.Category,
	}
	if err := client.DeleteUser(ctx, deleteUser); err != nil {
		t.Fatalf("delete user failed: %v", err)
	}

	// Verify deletion
	if _, err := client.GetUser(ctx, user.Username); err == nil {
		t.Fatal("user still exists after deletion")
	}

	// Test delete with wrong password (on a new user)
	testUser2 := &db.User{
		Username: "testuser2",
		Password: "password123",
		Email:    "test2@example.com",
		Category: 0,
	}
	if err := client.AddUser(ctx, testUser2); err != nil {
		t.Fatalf("failed to add test user 2: %v", err)
	}

	wrongPasswordUser := &db.User{
		Username: "testuser2",
		Password: "wrongpassword",
		Email:    "test2@example.com",
		Category: 0,
	}
	if err := client.DeleteUser(ctx, wrongPasswordUser); err == nil {
		t.Fatal("expected error when deleting with wrong password")
	}

	// Clean up test user 2
	if err := client.DeleteUser(ctx, testUser2); err != nil {
		t.Fatalf("failed to clean up test user 2: %v", err)
	}
}

func testCassandraInput(t *testing.T, ctx context.Context) {
	config := db.NewCassandraConfig("backend", "BPass0319", "cass_keyspace")
	client, err := db.NewCassandraClient(config)
	if err != nil {
		t.Fatalf("failed to create Cassandra client: %v", err)
	}
	defer client.Close()

	// Test cases for invalid users
	testCases := []struct {
		name        string
		user        *db.User
		shouldError bool
	}{
		{
			name: "empty username",
			user: &db.User{
				Username: "",
				Password: "password123",
				Email:    "test@example.com",
				Category: 0,
			},
			shouldError: true,
		},
		{
			name: "short username",
			user: &db.User{
				Username: "ab",
				Password: "password123",
				Email:    "test@example.com",
				Category: 0,
			},
			shouldError: true,
		},
		{
			name: "long username",
			user: &db.User{
				Username: "this_username_should_be_way_too_long_for_validation",
				Password: "password123",
				Email:    "test@example.com",
				Category: 0,
			},
			shouldError: true,
		},
		{
			name: "invalid username characters",
			user: &db.User{
				Username: "user@name",
				Password: "password123",
				Email:    "test@example.com",
				Category: 0,
			},
			shouldError: true,
		},
		{
			name: "short password",
			user: &db.User{
				Username: "testuser",
				Password: "short",
				Email:    "test@example.com",
				Category: 0,
			},
			shouldError: true,
		},
		{
			name: "invalid email",
			user: &db.User{
				Username: "testuser",
				Password: "password123",
				Email:    "invalid-email",
				Category: 0,
			},
			shouldError: true,
		},
		{
			name: "negative category",
			user: &db.User{
				Username: "testuser",
				Password: "password123",
				Email:    "test@example.com",
				Category: -1,
			},
			shouldError: false,
		},
		{
			name: "valid user",
			user: &db.User{
				Username: "validuser",
				Password: "password123",
				Email:    "valid@example.com",
				Category: 5,
			},
			shouldError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := client.AddUser(ctx, tc.user)
			if tc.shouldError && err == nil {
				t.Errorf("expected error for %s, got none", tc.name)
			}
			if !tc.shouldError && err != nil {
				t.Errorf("unexpected error for %s: %v", tc.name, err)
			}

			// Clean up if user was created
			if !tc.shouldError && err == nil {
				client.DeleteUser(ctx, tc.user)
			}
		})
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
	if _, err := client.Get(ctx, user); err == nil {
		t.Fatal("expected error on get non-existent user, got none")
	}

	// Add user
	if err := client.Add(ctx, user); err != nil {
		t.Fatalf("add user failed: %v", err)
	}

	// Get and verify user
	got, err := client.Get(ctx, user)
	if err != nil {
		t.Fatalf("get user failed: %v", err)
	}
	verifyUsersEqual(t, got, user)

	oldPassword := user.Password
	user.Password = "some_invalid_password232"

	if _, err := client.Get(ctx, user); err == nil {
		t.Fatal("expected error on get invalid password, got none")
	}

	user.Password = oldPassword

	err = client.Delete(ctx, user.Username)
	if err != nil {
		t.Fatalf("delete user failed: %v", err)
	}

	_, err = client.Get(ctx, user)
	if err == nil {
		t.Fatalf("delete user failed: user still exists")
	}
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

	if err := client.Add(ctx, user); err != nil {
		t.Fatalf("add user failed: %v", err)
	}

	// Verify user exists before expiration
	if _, err := client.Get(ctx, user); err != nil {
		t.Fatalf("user should exist before expiration: %v", err)
	}

	// Wait for expiration + buffer
	time.Sleep(config.Expiration + 2*time.Second)

	// Verify user is expired
	if _, err := client.Get(ctx, user); err == nil {
		t.Fatal("user still exists after expiration")
	}
}

func TestCassandraCRUD(t *testing.T) {
	ctx := context.Background()
	testCassandraAddGetUpdateDelete(t, ctx)
}

func TestCassandraInput(t *testing.T) {
	ctx := context.Background()
	testCassandraInput(t, ctx)
}

func TestRedisCRUD(t *testing.T) {
	ctx := context.Background()
	testRedisBasicCRUD(t, ctx)
}

func TestRedisExpiration(t *testing.T) {
	ctx := context.Background()
	testRedisExpiration(t, ctx)
}
