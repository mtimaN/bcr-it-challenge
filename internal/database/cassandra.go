package db

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/gocql/gocql"
)

// CassandraConfig holds the configuration for Cassandra connection
type CassandraConfig struct {
	username string
	password string

	keyspace string
	hosts    []string

	timeout        time.Duration
	connectTimeout time.Duration
}

// CassandraClient is a wrapper around the Cassandra connection
type CassandraClient struct {
	session *gocql.Session
	config  *CassandraConfig
}

func NewCassandraConfig(username string, password string, keyspace string) *CassandraConfig {
	hosts := []string{"localhost"}

	return &CassandraConfig{
		username: username,
		password: password,

		keyspace: keyspace,
		hosts:    hosts,

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
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return err
	}

	if len(user.Username) == 0 || len(user.Password) == 0 {
		return errors.New("required fields are empty")
	}

	return nil
}

func (c *CassandraClient) AddUser(ctx context.Context, user *User) error {
	if err := validUser(user); err != nil {
		return fmt.Errorf("user creation: %w", err)
	}

	if err := c.session.Query(
		"INSERT INTO users (username, password, email, category) VALUES (?, ?, ?, ?)",
		user.Username, user.Password, user.Email, user.Category).WithContext(ctx).Exec(); err != nil {
		return fmt.Errorf("user creation: %w", err)
	}

	return nil
}

func (c *CassandraClient) UpdateUser(ctx context.Context, user *User, oldPassword string) error {
	old, err := c.GetUser(ctx, user.Username)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	if old.Password != oldPassword {
		return errors.New("update user: old password does not match")
	}

	if err := c.session.Query(
		"UPDATE users SET password = ?, email = ?, category = ? WHERE username = ?",
		user.Password, user.Email, user.Category, user.Username).WithContext(ctx).Exec(); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (c *CassandraClient) DeleteUser(ctx context.Context, user *User) error {
	got, err := c.GetUser(ctx, user.Username)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if got.Password != user.Password {
		return errors.New("delete user: old password does not match")
	}

	if err := c.session.Query(
		"DELETE FROM users WHERE username = ?",
		user.Username).WithContext(ctx).Exec(); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}

func (c *CassandraClient) GetUser(ctx context.Context, username string) (*User, error) {
	user := &User{}

	if err := c.session.Query(
		"SELECT username, password, email, category FROM users WHERE username = ? LIMIT 1",
		username).WithContext(ctx).Scan(&user.Username, &user.Password, &user.Email, &user.Category); err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

func (c *CassandraClient) Close() {
	c.session.Close()
}
