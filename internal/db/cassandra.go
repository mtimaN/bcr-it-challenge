package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

var (
	// Error definitions
	ErrSessionNotInitialized = errors.New("session: cassandra not initialized")
	ErrUserNotFound          = errors.New("not found: user not found")
	ErrDatabaseError         = errors.New("internal: database error")
	ErrUsernameExists        = errors.New("validation: username already exists")
	ErrUsernameNotExists     = errors.New("validation: username does not exist")
	ErrInvalidEmail          = errors.New("validation: invalid email")
	ErrPasswordProcessing    = errors.New("validation: password processing failed")
	ErrUserCreationFailed    = errors.New("validation: user creation failed")
	ErrUpdateFailed          = errors.New("update failed")
	ErrDeletionFailed        = errors.New("internal: deletion failed")
)

// CassandraConfig holds the configuration for Cassandra connection
type CassandraConfig struct {
	Username       string
	Password       string
	Keyspace       string
	Hosts          []string
	Timeout        time.Duration
	ConnectTimeout time.Duration
}

// NewCassandraConfig creates a new configuration with defaults
func NewCassandraConfig(username, password, keyspace string) *CassandraConfig {
	return &CassandraConfig{
		Username:       username,
		Password:       password,
		Keyspace:       keyspace,
		Hosts:          []string{"localhost"},
		Timeout:        5 * time.Second,
		ConnectTimeout: 10 * time.Second,
	}
}

// UserRepository defines the interface for database operations
type UserRepository interface {
	Health(ctx context.Context) error
	GetUser(ctx context.Context, username string) (*User, error)
	AddUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, username string) error
	UsernameExists(ctx context.Context, username string) (bool, error)
	Stats(ctx context.Context) (map[string]interface{}, error)
	Close()
}

// CassandraRepo is a wrapper around the Cassandra connection
type CassandraRepo struct {
	session *gocql.Session
	config  *CassandraConfig
}

var _ UserRepository = (*CassandraRepo)(nil)

// NewCassandraRepo creates a new Cassandra UserRepository
func NewCassandraRepo(config *CassandraConfig) (UserRepository, error) {
	cluster := gocql.NewCluster(config.Hosts...)
	cluster.Keyspace = config.Keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.Username,
		Password: config.Password,
	}
	cluster.Timeout = config.Timeout
	cluster.ConnectTimeout = config.ConnectTimeout

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("cassandra session: %w", err)
	}

	return &CassandraRepo{
		session: session,
		config:  config,
	}, nil
}

// ensureSession checks if the session is initialized
func (c *CassandraRepo) ensureSession() error {
	if c.session == nil {
		return ErrSessionNotInitialized
	}
	return nil
}

// Health checks if the connection to Cassandra is working
func (c *CassandraRepo) Health(ctx context.Context) error {
	if err := c.ensureSession(); err != nil {
		return err
	}

	var version string
	if err := c.session.Query("SELECT release_version FROM system.local").
		WithContext(ctx).Scan(&version); err != nil {
		return fmt.Errorf("health: %w", err)
	}

	return nil
}

// GetUser retrieves a user by username
func (c *CassandraRepo) GetUser(ctx context.Context, username string) (*User, error) {
	if err := c.ensureSession(); err != nil {
		return nil, err
	}

	user := &User{Credentials: &Credentials{}}

	err := c.session.Query(
		"SELECT username, password, email, category FROM users WHERE username = ? LIMIT 1",
		username).WithContext(ctx).Scan(
		&user.Username, &user.Password, &user.Email, &user.Category)

	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, ErrDatabaseError
	}

	return user, nil
}

// AddUser adds a new user to the database
func (c *CassandraRepo) AddUser(ctx context.Context, user *User) error {
	if err := c.ensureSession(); err != nil {
		return err
	}

	if user.Email == "" {
		return ErrInvalidEmail
	}

	if err := ValidUser(user); err != nil {
		return fmt.Errorf("validation: %w", err)
	}

	if err := c.session.Query(
		"INSERT INTO users (username, password, email, category) VALUES (?, ?, ?, ?)",
		user.Username, user.Password, user.Email, user.Category).
		WithContext(ctx).Exec(); err != nil {
		return ErrUserCreationFailed
	}

	return nil
}

// UpdateUser updates an existing user
func (c *CassandraRepo) UpdateUser(ctx context.Context, user *User) error {
	if err := c.ensureSession(); err != nil {
		return err
	}

	ok, err := c.UsernameExists(ctx, user.Username)
	if err != nil {
		return err
	}
	if !ok {
		return ErrUserNotFound
	}

	if err := c.session.Query(
		"UPDATE users SET password = ?, email = ?, category = ? WHERE username = ?",
		user.Password, user.Email, user.Category, user.Username).
		WithContext(ctx).Exec(); err != nil {
		return ErrUpdateFailed
	}

	return nil
}

// DeleteUser deletes a user by username
func (c *CassandraRepo) DeleteUser(ctx context.Context, username string) error {
	if err := c.ensureSession(); err != nil {
		return err
	}

	exists, err := c.UsernameExists(ctx, username)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUsernameNotExists
	}

	if err := c.session.Query(
		"DELETE FROM users WHERE username = ?", username).
		WithContext(ctx).Exec(); err != nil {
		return ErrDeletionFailed
	}

	return nil
}

// UsernameExists checks if a username exists
func (c *CassandraRepo) UsernameExists(ctx context.Context, username string) (bool, error) {
	if err := c.ensureSession(); err != nil {
		return false, err
	}

	if len(username) < MinUsernameLength {
		return false, errors.New("validation: username required")
	}
	if len(username) > MaxUsernameLength {
		return false, errors.New("validation: username too long")
	}

	var dummy string
	err := c.session.Query(
		"SELECT username FROM users WHERE username = ? LIMIT 1", username).
		WithContext(ctx).Scan(&dummy)

	if err != nil {
		if err == gocql.ErrNotFound {
			return false, nil
		}
		return false, ErrDatabaseError
	}

	return true, nil
}

// Stats retrieves statistics from the system.stats table
func (c *CassandraRepo) Stats(ctx context.Context) (map[string]interface{}, error) {
	if err := c.ensureSession(); err != nil {
		return nil, err
	}

	var tableName string
	var sstables int
	var readLatency, writeLatency float64

	err := c.session.Query(`
		SELECT table_name, sstables_count, read_latency_avg, write_latency_avg
		FROM system.stats WHERE table_name = ?`, "users").
		WithContext(ctx).Scan(&tableName, &sstables, &readLatency, &writeLatency)

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

// Close closes the Cassandra session
func (c *CassandraRepo) Close() {
	if c.session != nil {
		c.session.Close()
	}
}
