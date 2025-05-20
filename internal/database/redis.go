package db

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConfig holds configuration for Redis connection
type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	MaxRetries   int
	Expiration   time.Duration
	KeyPrefix    string // Add namespace for keys
}

// RedisRepo is a wrapper around the Redis client
type RedisRepo struct {
	client *redis.Client
	config *RedisConfig
}

// NewRedisConfig creates a new Redis configuration with secure defaults
func NewRedisConfig(password string) *RedisConfig {
	return &RedisConfig{
		Addr:         "localhost:6379",
		Password:     password,
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 2,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		MaxRetries:   3,
		Expiration:   24 * time.Hour,
		KeyPrefix:    "cache:user:", // Namespace for cache keys
	}
}

// NewRedisRepo creates a new Redis client with improved configuration
func NewRedisRepo(config *RedisConfig) (*RedisRepo, error) {
	if config == nil {
		return nil, errors.New("redis config cannot be nil")
	}

	client := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		MaxRetries:   config.MaxRetries,
	})

	ctx, cancel := context.WithTimeout(context.Background(), config.DialTimeout)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("connect to redis: %w", err)
	}

	return &RedisRepo{
		client: client,
		config: config,
	}, nil
}

// Health checks the Redis connection health
func (c *RedisRepo) Health(ctx context.Context) error {
	if c.client == nil {
		return errors.New("redis client is nil")
	}

	if _, err := c.client.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("redis health check: %w", err)
	}

	return nil
}

// createKey creates a secure, namespaced key from username
func (c *RedisRepo) createKey(username string) (string, error) {
	// Input validation
	if username == "" {
		return "", errors.New("username cannot be empty")
	}

	// Normalize username (trim spaces, convert to lowercase)
	username = strings.ToLower(strings.TrimSpace(username))

	// Additional validation - ensure username doesn't contain control characters
	if strings.ContainsAny(username, "\r\n\t\000") {
		return "", errors.New("username contains invalid characters")
	}

	// Create hash
	h := sha256.New()
	h.Write([]byte(username))
	hash := hex.EncodeToString(h.Sum(nil))

	// Return namespaced key
	return c.config.KeyPrefix + hash, nil
}

// Get retrieves a user from cache
func (c *RedisRepo) Get(ctx context.Context, cred *Credentials) (*User, error) {
	if c.client == nil {
		return nil, errors.New("redis client is not initialized")
	}

	key, err := c.createKey(cred.Username)
	if err != nil {
		return nil, fmt.Errorf("create cache key: %w", err)
	}

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("user not found in cache")
		}
		return nil, fmt.Errorf("get from redis: %w", err)
	}

	var user User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		// If we can't unmarshal, delete the corrupted entry
		_ = c.client.Del(ctx, key)
		return nil, fmt.Errorf("unmarshal cached data: %w", err)
	}

	if !CheckPasswordHash(cred.Password, user.Password) {
		return nil, errors.New("cache: incorrect password")
	}

	user.Password = cred.Password
	return &user, nil
}

// Add stores a user in cache
func (c *RedisRepo) Add(ctx context.Context, user *User) error {
	if c.client == nil {
		return errors.New("redis client is not initialized")
	}
	if user == nil {
		return errors.New("user cannot be nil")
	}

	key, err := c.createKey(user.Username)
	if err != nil {
		return fmt.Errorf("create cache key: %w", err)
	}

	rUser := &User{
		Credentials: &Credentials{
			Username: user.Username,
			Password: user.Password,
		},
		Email:    user.Email,
		Category: user.Category,
	}

	rUser.Password, err = HashPassword(rUser.Password)
	if err != nil {
		return fmt.Errorf("cache password hash: %w", err)
	}

	val, err := json.Marshal(rUser)
	if err != nil {
		return fmt.Errorf("marshal user data: %w", err)
	}

	if err := c.client.Set(ctx, key, val, c.config.Expiration).Err(); err != nil {
		return fmt.Errorf("set in redis: %w", err)
	}

	return nil
}

// Delete removes a user from cache
func (c *RedisRepo) Delete(ctx context.Context, username string) error {
	if c.client == nil {
		return errors.New("redis client is not initialized")
	}

	key, err := c.createKey(username)
	if err != nil {
		return fmt.Errorf("create cache key: %w", err)
	}

	deleted, err := c.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("delete from redis: %w", err)
	}

	if deleted == 0 {
		return fmt.Errorf("user not found in cache")
	}

	return nil
}

// Exists checks if a user exists in cache without retrieving the full data
func (c *RedisRepo) Exists(ctx context.Context, username string) (bool, error) {
	if c.client == nil {
		return false, errors.New("redis client is not initialized")
	}

	key, err := c.createKey(username)
	if err != nil {
		return false, fmt.Errorf("create cache key: %w", err)
	}

	exists, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("check existence in redis: %w", err)
	}

	return exists > 0, nil
}

// Close gracefully closes the Redis connection
func (c *RedisRepo) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// Stats returns basic Redis statistics
func (c *RedisRepo) Stats(ctx context.Context) (map[string]interface{}, error) {
	if c.client == nil {
		return nil, errors.New("redis client is not initialized")
	}

	poolStats := c.client.PoolStats()

	return map[string]interface{}{
		"total_connections": poolStats.TotalConns,
		"idle_connections":  poolStats.IdleConns,
		"stale_connections": poolStats.StaleConns,
		"hits":              poolStats.Hits,
		"misses":            poolStats.Misses,
		"timeouts":          poolStats.Timeouts,
	}, nil
}
