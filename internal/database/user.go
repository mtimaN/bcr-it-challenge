package db

import (
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Category int    `json:"category"`
}

func UserID(username string) (gocql.UUID, error) {
	namespace := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	id := uuid.NewSHA1(namespace, []byte(username))

	cqlID, err := gocql.UUIDFromBytes(id[:])
	if err != nil {
		return gocql.UUID{}, err
	}

	return cqlID, nil
}

func NewUser(username string, password string, email string) *User {
	return &User{
		Username: username,
		Password: password,
		Email:    email,
		Category: -1,
	}
}
