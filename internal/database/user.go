package db

import "golang.org/x/crypto/bcrypt"

// User represents a user in the system
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Category int    `json:"category"`
}

func NewUser(username string, password string, email string) *User {
	return &User{
		Username: username,
		Password: password,
		Email:    email,
		Category: -1,
	}
}

// Password hashing functions
func HashPassword(password string) (string, error) {
	// Use cost 12 for better security (adjust based on your performance requirements)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
