package main

import (
	"context"
	"db"
	"errors"
	"fmt"
)

func testCassandra(ctx context.Context) error {
	config := db.NewCassandraConfig(
		"backend", "BPass0508", "cass_keyspace",
	)
	client, err := db.NewCassandraClient(config)
	if err != nil {
		return err
	}
	defer client.Close()

	// Create user
	user := db.NewUser("radu", "fmmdepraf", "radu232")
	user.Category = 0

	if err := client.AddUser(ctx, user); err != nil {
		return err
	}

	// Get user
	user_copy, err := client.GetUser(ctx, user.Username)
	if err != nil {
		return err
	}

	if user_copy.Username != user.Username ||
		user_copy.Password != user.Password ||
		user_copy.Email != user.Email ||
		user_copy.Category != user.Category {
		return errors.New("invalid user data")
	}

	fmt.Printf("User data: %v\n", user_copy)

	// Update user
	user.Email = "newEmail"
	user.Category = 1
	if err := client.UpdateUser(ctx, user); err != nil {
		return err
	}

	user, err = client.GetUser(ctx, user.Username)
	if err != nil {
		return err
	}

	if user.Email != "newEmail" || user.Category != 1 {
		return errors.New("update user failed")
	}

	// Invalid user get
	_, err = client.GetUser(ctx, "fmmdepraf")
	if err == nil {
		return errors.New("invalid user get - should fail with non-existent username")
	}

	// Invalid user update
	nonExistentUser := db.NewUser("nonexistent", "password", "email")
	nonExistentUser.Category = 2
	if err := client.UpdateUser(ctx, nonExistentUser); err == nil {
		return errors.New("invalid user update - should fail with non-existent username")
	}

	// Delete user
	if err := client.DeleteUser(ctx, user.Username); err != nil {
		return err
	}

	_, err = client.GetUser(ctx, "radu")
	if err == nil {
		return errors.New("user deletion failed - user still exists")
	}

	// Pass
	return nil
}

func main() {
	ctx := context.Background()
	if err := testCassandra(ctx); err != nil {
		fmt.Println(fmt.Errorf("cassandra: %w", err))
	} else {
		fmt.Println("Cassandra tests passed!")
	}
}
