package db

import (
	"context"
	"errors"
	"fmt"
	"time"

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

// CassandraRepo is a wrapper around the Cassandra connection
type CassandraRepo struct {
	session *gocql.Session
	config  *CassandraConfig
}

func NewCassandraConfig(username, password, keyspace string) *CassandraConfig {
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

func NewCassandraRepo(config *CassandraConfig) (*CassandraRepo, error) {
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

	return &CassandraRepo{
		session: session,
		config:  config,
	}, nil
}

func (c *CassandraRepo) Health(ctx context.Context) error {
	if c.session == nil {
		return errors.New("cassandra session is nil")
	}

	var test string
	if err := c.session.Query("SELECT release_version FROM system.local").WithContext(ctx).Scan(&test); err != nil {
		return fmt.Errorf("cassandra health: %w", err)
	}

	return nil
}

func (c *CassandraRepo) GetUser(ctx context.Context, cred *Credentials) (*User, error) {
	user := &User{}
	user.Credentials = &Credentials{}

	if err := c.session.Query(
		"SELECT username, password, email, category FROM users WHERE username = ? LIMIT 1",
		cred.Username).WithContext(ctx).Scan(&user.Username, &user.Password, &user.Email, &user.Category); err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("database error")
	}

	user.Password = cred.Password
	return user, nil
}

func (c *CassandraRepo) AddUser(ctx context.Context, user *User) error {
	if err := ValidUser(user); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if user already exists
	exists, err := c.UsernameExists(ctx, user.Username)
	if err != nil {
		return err
	}
	if exists {
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

func (c *CassandraRepo) UpdateUser(ctx context.Context, user *User, oldPassword string) error {
	if err := ValidUser(user); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Validate password
	old, err := c.GetUser(ctx, NewCredentials(user.Username, oldPassword))
	if err != nil {
		return ErrAuthenticationFailed
	}

	email := user.Email
	if email == "" {
		email = old.Email
	}

	category := user.Category
	if category < 0 || category > 3 {
		category = old.Category
	}

	// Hash new password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return errors.New("password processing failed")
	}

	if err := c.session.Query(
		"UPDATE users SET password = ?, email = ?, category = ? WHERE username = ?",
		hashedPassword, email, category, user.Username).WithContext(ctx).Exec(); err != nil {
		return errors.New("update failed")
	}

	return nil
}

func (c *CassandraRepo) DeleteUser(ctx context.Context, username string) error {
	ok, err := c.UsernameExists(ctx, username)
	if err != nil {
		return fmt.Errorf("deletion: %w", err)
	}
	if !ok {
		return errors.New("username does not exist")
	}

	if err := c.session.Query(
		"DELETE FROM users WHERE username = ?",
		username).WithContext(ctx).Exec(); err != nil {
		return errors.New("deletion failed")
	}

	return nil
}

// Method to check if username exists (for registration)
func (c *CassandraRepo) UsernameExists(ctx context.Context, username string) (bool, error) {
	if len(username) < MinUsernameLength {
		return false, errors.New("username required")
	}
	if len(username) > MaxUsernameLength {
		return false, errors.New("username too long")
	}

	var dummy string
	if err := c.session.Query(
		"SELECT username FROM users WHERE username = ? LIMIT 1",
		username).WithContext(ctx).Scan(&dummy); err != nil {
		if err == gocql.ErrNotFound {
			return false, nil
		}
		return false, errors.New("database error")
	}

	return true, nil
}

// Attempts to retrieve stats from system.stats table
func (c *CassandraRepo) Stats(ctx context.Context) (map[string]interface{}, error) {
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

func (c *CassandraRepo) Close() {
	if c.session != nil {
		c.session.Close()
	}
}
