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

func createTestCredentials() *db.Credentials {
	cred := db.NewCredentials("test", "testStrongPw123")
	return cred
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
	client, err := db.NewCassandraRepo(config)
	if err != nil {
		t.Fatalf("failed to create Cassandra client: %v", err)
	}
	defer client.Close()

	cred := createTestCredentials()
	user := createTestUser()
	originalPassword := user.Password // Store original password for later use

	// Test update non-existent user
	if err := client.UpdateUser(ctx, user, user.Password); err == nil {
		t.Fatal("expected error on update non-existent user, got none")
	}

	exists, err := client.UsernameExists(ctx, user.Username)
	if err != nil {
		t.Fatalf("username exists check for non-existent failed: %v", err)
	}
	if exists {
		t.Error("nonexistent username should not exist")
	}

	// Test invalid get (non-existent user)
	if _, err := client.GetUser(ctx, cred); err == nil {
		t.Fatal("expected error on get non-existent user, got none")
	}

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
	got, err := client.GetUser(ctx, cred)
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

	wrongCred := createTestCredentials()
	wrongCred.Password = "WrongPassword123"

	// Test authentication with wrong password
	if _, err := client.GetUser(ctx, wrongCred); err == nil {
		t.Fatal("expected authentication to fail with wrong password")
	}

	wrongCred = createTestCredentials()
	wrongCred.Username = "wrongUsername"

	// Test authentication with non-existent user
	if _, err := client.GetUser(ctx, wrongCred); err == nil {
		t.Fatal("expected authentication to fail for non-existent user")
	}

	// Update user
	newPassword := "newPassword123"
	updatedUser := &db.User{
		Credentials: db.NewCredentials(user.Username, newPassword),
		Email:       "newEmail@test.com",
		Category:    1,
	}

	updatedCred := db.NewCredentials(updatedUser.Username, updatedUser.Password)

	// Update should work with correct old password
	if err := client.UpdateUser(ctx, updatedUser, originalPassword); err != nil {
		t.Fatalf("update user failed: %v", err)
	}

	// Verify update
	got, err = client.GetUser(ctx, updatedCred)
	if err != nil {
		t.Fatalf("get updated user failed: %v", err)
	}
	if got.Email != "newEmail@test.com" || got.Category != 1 {
		t.Fatalf("update verification failed: got email=%s category=%d", got.Email, got.Category)
	}

	// Test authentication with old password (should fail)
	if _, err := client.GetUser(ctx, cred); err == nil {
		t.Fatal("authentication with old password should fail after update")
	}

	// Test update with wrong old password (should fail)
	anotherUpdate := &db.User{
		Credentials: db.NewCredentials(user.Username, "anotherPassword"),
		Email:       "another@test.com",
		Category:    2,
	}
	if err := client.UpdateUser(ctx, anotherUpdate, "wrongOldPassword"); err == nil {
		t.Fatal("expected error when updating with wrong old password")
	}

	// Test username existence check
	exists, err = client.UsernameExists(ctx, user.Username)
	if err != nil {
		t.Fatalf("username exists check failed: %v", err)
	}
	if !exists {
		t.Error("username should exist")
	}

	// Delete user
	if err := client.DeleteUser(ctx, user.Username); err != nil {
		t.Fatalf("delete user failed: %v", err)
	}

	// Verify deletion
	if _, err := client.GetUser(ctx, cred); err == nil {
		t.Fatal("user still exists after deletion")
	}
}

func testCassandraInput(t *testing.T, ctx context.Context) {
	config := db.NewCassandraConfig("backend", "BPass0319", "cass_keyspace")
	client, err := db.NewCassandraRepo(config)
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
				Credentials: db.NewCredentials("", "password123"),
				Email:       "test@example.com",
				Category:    0,
			},
			shouldError: true,
		},
		{
			name: "short username",
			user: &db.User{
				Credentials: db.NewCredentials("ab", "password123"),
				Email:       "test@example.com",
				Category:    0,
			},
			shouldError: true,
		},
		{
			name: "long username",
			user: &db.User{
				Credentials: db.NewCredentials("this_username_should_be_way_too_long_for_validation", "password123"),
				Email:       "test@example.com",
				Category:    0,
			},
			shouldError: true,
		},
		{
			name: "invalid username characters",
			user: &db.User{
				Credentials: db.NewCredentials("user@name", "password123"),
				Email:       "test@example.com",
				Category:    0,
			},
			shouldError: true,
		},
		{
			name: "short password",
			user: &db.User{
				Credentials: db.NewCredentials("testuser", "short"),
				Email:       "test@example.com",
				Category:    0,
			},
			shouldError: true,
		},
		{
			name: "invalid email",
			user: &db.User{
				Credentials: db.NewCredentials("testuser", "password123"),
				Email:       "invalid-email",
				Category:    0,
			},
			shouldError: true,
		},
		{
			name: "negative category",
			user: &db.User{
				Credentials: db.NewCredentials("testuser", "password123"),
				Email:       "test@example.com",
				Category:    -1,
			},
			shouldError: true,
		},
		{
			name: "valid user",
			user: &db.User{
				Credentials: db.NewCredentials("validuser", "password123"),
				Email:       "valid@example.com",
				Category:    2,
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
				client.DeleteUser(ctx, tc.user.Username)
			}
		})
	}
}

func testRedisBasicCRUD(t *testing.T, ctx context.Context) {
	config := db.NewRedisConfig("RPass0319")
	client, err := db.NewRedisRepo(config)
	if err != nil {
		t.Fatalf("failed to create Redis client: %v", err)
	}
	defer client.Close()

	user := createTestUser()
	cred := createTestCredentials()

	// Invalid get (non-existent user)
	if _, err := client.Get(ctx, cred); err == nil {
		t.Fatal("expected error on get non-existent user, got none")
	}

	// Add user
	if err := client.Add(ctx, user); err != nil {
		t.Fatalf("add user failed: %v", err)
	}

	// Get and verify user
	got, err := client.Get(ctx, cred)
	if err != nil {
		t.Fatalf("get user failed: %v", err)
	}
	verifyUsersEqual(t, got, user)

	oldPassword := cred.Password
	cred.Password = "some_invalid_password232"

	if _, err := client.Get(ctx, cred); err == nil {
		t.Fatal("expected error on get invalid password, got none")
	}

	cred.Password = oldPassword

	err = client.Delete(ctx, user.Username)
	if err != nil {
		t.Fatalf("delete user failed: %v", err)
	}

	_, err = client.Get(ctx, cred)
	if err == nil {
		t.Fatalf("delete user failed: user still exists")
	}
}

func testRedisExpiration(t *testing.T, ctx context.Context) {
	config := db.NewRedisConfig("RPass0319")
	config.Expiration = 10 * time.Second
	client, err := db.NewRedisRepo(config)
	if err != nil {
		t.Fatalf("failed to create Redis client with expiration: %v", err)
	}
	defer client.Close()

	user := createTestUser()
	cred := createTestCredentials()

	if err := client.Add(ctx, user); err != nil {
		t.Fatalf("add user failed: %v", err)
	}

	// Verify user exists before expiration
	if _, err := client.Get(ctx, cred); err != nil {
		t.Fatalf("user should exist before expiration: %v", err)
	}

	// Wait for expiration + buffer
	time.Sleep(config.Expiration + 2*time.Second)

	// Verify user is expired
	if _, err := client.Get(ctx, cred); err == nil {
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
