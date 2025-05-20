package db

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
