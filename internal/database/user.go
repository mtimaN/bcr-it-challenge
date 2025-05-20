package db

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User represents a user in the system
type User struct {
	*Credentials
	Email    string `json:"email"`
	Category int    `json:"category"`
}

func NewCredentials(username, password string) *Credentials {
	return &Credentials{
		Username: username,
		Password: password,
	}
}

func NewUser(username, password, email string) *User {
	return &User{
		Credentials: NewCredentials(username, password),
		Email:       email,
		Category:    0,
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

func ValidCredentials(username, password string) error {
	// Username validation
	if username == "" {
		return errors.New("username is required")
	}
	if len(username) < MinUsernameLength || len(username) > MaxUsernameLength {
		return fmt.Errorf("username must be %d-%d characters", MinUsernameLength, MaxUsernameLength)
	}

	// Allow alphanumeric, underscore, and hyphen
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validUsername.MatchString(username) {
		return errors.New("username contains invalid characters")
	}

	// Password validation
	if password == "" {
		return errors.New("password is required")
	}
	if len(password) < MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", MinPasswordLength)
	}
	if len(password) > MaxPasswordLength {
		return errors.New("password too long")
	}

	return nil
}

func ValidUser(user *User) error {
	// Email validation
	if user.Email == "" {
		return errors.New("email is required")
	}
	if len(user.Email) > MaxEmailLength {
		return errors.New("email too long")
	}
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return errors.New("invalid email format")
	}

	if err := ValidCredentials(user.Username, user.Password); err != nil {
		return errors.New("invalid credentials")
	}

	// Check UTF-8 validity
	if !utf8.ValidString(user.Username) || !utf8.ValidString(user.Email) {
		return errors.New("invalid character encoding")
	}

	if user.Category < 0 || user.Category > 3 {
		return errors.New("invalid category")
	}

	return nil
}
