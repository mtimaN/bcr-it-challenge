package db

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/gocql/gocql"
)

// CassandraConfig holds the configuration for Cassandra connection
type CassandraConfig struct {
	username string
	password string

	keyspace       string
	hosts          []string
	timeout        time.Duration
	connectTimeout time.Duration
}

// CassandraClient is a wrapper around the Cassandra connection
type CassandraClient struct {
	session *gocql.Session
	config  *CassandraConfig
}

const (
	MinPasswordLength = 8
	MaxPasswordLength = 128
	MinUsernameLength = 3
	MaxUsernameLength = 20
	MaxEmailLength    = 254 // RFC 5321 limit
	BcryptCost        = 12  // increase for better security
)

func NewCassandraConfig(username string, password string, keyspace string) *CassandraConfig {
	hosts := []string{"localhost"}

	return &CassandraConfig{
		username: username,
		password: password,

		keyspace:       keyspace,
		hosts:          hosts,
		timeout:        5 * time.Second,
		connectTimeout: 10 * time.Second,
	}
}

func NewCassandraClient(config *CassandraConfig) (*CassandraClient, error) {
	cluster := gocql.NewCluster(config.hosts...)
	cluster.Keyspace = config.keyspace
	cluster.Consistency = gocql.Quorum

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.username,
		Password: config.password,
	}

	cluster.Timeout = config.timeout
	cluster.ConnectTimeout = config.connectTimeout

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("cassandra session: %w", err)
	}

	return &CassandraClient{
		session: session,
		config:  config,
	}, nil
}

func (c *CassandraClient) Health(ctx context.Context) error {
	if c.session == nil {
		return errors.New("cassandra session is nil")
	}

	var test string
	if err := c.session.Query("SELECT release_version FROM system.local").WithContext(ctx).Scan(&test); err != nil {
		return fmt.Errorf("cassandra health: %w", err)
	}

	return nil
}

func validUser(user *User) error {
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

	// Username validation
	if user.Username == "" {
		return errors.New("username is required")
	}
	if len(user.Username) < MinUsernameLength || len(user.Username) > MaxUsernameLength {
		return fmt.Errorf("username must be %d-%d characters", MinUsernameLength, MaxUsernameLength)
	}

	// Allow alphanumeric, underscore, and hyphen
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validUsername.MatchString(user.Username) {
		return errors.New("username contains invalid characters")
	}

	// Password validation
	if user.Password == "" {
		return errors.New("password is required")
	}
	if len(user.Password) < MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", MinPasswordLength)
	}
	if len(user.Password) > MaxPasswordLength {
		return errors.New("password too long")
	}

	// Check UTF-8 validity
	if !utf8.ValidString(user.Username) || !utf8.ValidString(user.Email) {
		return errors.New("invalid character encoding")
	}

	return nil
}

func (c *CassandraClient) AddUser(ctx context.Context, user *User) error {
	if err := validUser(user); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if user already exists
	existing, _ := c.GetUser(ctx, user.Username)
	if existing != nil {
		return errors.New("username already exists")
	}

	// Hash password before storing
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return errors.New("password processing failed")
	}

	if err := c.session.Query(
		"INSERT INTO users (username, password, email, category) VALUES (?, ?, ?, ?)",
		user.Username, hashedPassword, user.Email, user.Category).WithContext(ctx).Exec(); err != nil {
		return errors.New("user creation failed")
	}

	return nil
}

var ErrAuthenticationFailed error = errors.New("authentication failed")

func (c *CassandraClient) UpdateUser(ctx context.Context, user *User, oldPassword string) error {
	if err := validUser(user); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if user.Password == oldPassword {
		return errors.New("validation failed: new passsword cannot be old password")
	}

	// Get existing user
	old, err := c.GetUser(ctx, user.Username)
	if err != nil {
		return ErrAuthenticationFailed
	}

	// Verify old password using constant-time comparison
	if !CheckPasswordHash(oldPassword, old.Password) {
		return ErrAuthenticationFailed
	}

	// Hash new password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return errors.New("password processing failed")
	}

	if err := c.session.Query(
		"UPDATE users SET password = ?, email = ?, category = ? WHERE username = ?",
		hashedPassword, user.Email, user.Category, user.Username).WithContext(ctx).Exec(); err != nil {
		return errors.New("update failed")
	}

	return nil
}

func (c *CassandraClient) DeleteUser(ctx context.Context, user *User) error {
	// Get existing user
	got, err := c.GetUser(ctx, user.Username)
	if err != nil {
		return ErrAuthenticationFailed
	}

	// Verify password using constant-time comparison
	if !CheckPasswordHash(user.Password, got.Password) {
		return ErrAuthenticationFailed
	}

	if err := c.session.Query(
		"DELETE FROM users WHERE username = ?",
		user.Username).WithContext(ctx).Exec(); err != nil {
		return errors.New("deletion failed")
	}

	return nil
}

func (c *CassandraClient) GetUser(ctx context.Context, username string) (*User, error) {
	// Basic input validation for username
	if username == "" {
		return nil, errors.New("username required")
	}
	if len(username) > MaxUsernameLength {
		return nil, errors.New("username too long")
	}

	user := &User{}
	if err := c.session.Query(
		"SELECT username, password, email, category FROM users WHERE username = ? LIMIT 1",
		username).WithContext(ctx).Scan(&user.Username, &user.Password, &user.Email, &user.Category); err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("database error")
	}

	return user, nil
}

// New method for authentication (recommended)
func (c *CassandraClient) AuthenticateUser(ctx context.Context, username, password string) (*User, error) {
	user, err := c.GetUser(ctx, username)
	if err != nil {
		return nil, ErrAuthenticationFailed
	}

	if !CheckPasswordHash(password, user.Password) {
		return nil, ErrAuthenticationFailed
	}

	// Don't return the password hash
	user.Password = ""
	return user, nil
}

// Method to check if username exists (for registration)
func (c *CassandraClient) UsernameExists(ctx context.Context, username string) (bool, error) {
	user, err := c.GetUser(ctx, username)
	if err != nil {
		if err.Error() == "user not found" {
			return false, nil
		}
		return false, err
	}
	return user != nil, nil
}

// attempts to retrieve stats from system.stats
func (c *CassandraClient) Stats(ctx context.Context) (map[string]interface{}, error) {
	var tableName string
	var sstables int
	var readLatency float64
	var writeLatency float64

	err := c.session.Query(`SELECT table_name, sstables_count, read_latency_avg,
	write_latency_avg FROM system.stats WHERE table_name = ?`, "users").
		WithContext(ctx).
		Scan(&tableName, &sstables, &readLatency, &writeLatency)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"table_name":    tableName,
		"sstables":      sstables,
		"read_latency":  readLatency,
		"write_latency": writeLatency,
	}, nil
}

func (c *CassandraClient) Close() {
	if c.session != nil {
		c.session.Close()
	}
}
