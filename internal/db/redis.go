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

// UserCache defines the interface for cache operations
type UserCache interface {
	Get(ctx context.Context, username string) (*User, error)
	Add(ctx context.Context, user *User) error
	Delete(ctx context.Context, username string) error
	Exists(ctx context.Context, username string) (bool, error)
	Extend(ctx context.Context, username string) error
	Health(ctx context.Context) error
	Close() error
	Stats(ctx context.Context) (map[string]interface{}, error)
}

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
	KeyPrefix    string
}

// RedisRepo implements the UserCache interface using Redis
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
		KeyPrefix:    "cache:user:",
	}
}

// NewRedisRepo creates a new Redis client with improved configuration
func NewRedisRepo(config *RedisConfig) (UserCache, error) {
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

	return &RedisRepo{client: client, config: config}, nil
}

// createKey creates a secure, namespaced key from username
func (r *RedisRepo) createKey(username string) (string, error) {
	if username == "" {
		return "", errors.New("validation: username cannot be empty")
	}

	username = strings.ToLower(strings.TrimSpace(username))
	if strings.ContainsAny(username, "\r\n\t\000") {
		return "", errors.New("validation: username contains invalid characters")
	}

	h := sha256.New()
	h.Write([]byte(username))
	return r.config.KeyPrefix + hex.EncodeToString(h.Sum(nil)), nil
}

// Health checks the Redis connection health
func (r *RedisRepo) Health(ctx context.Context) error {
	if r.client == nil {
		return errors.New("internal: redis client is nil")
	}
	_, err := r.client.Ping(ctx).Result()
	return err
}

// Get retrieves a user from cache
func (r *RedisRepo) Get(ctx context.Context, username string) (*User, error) {
	if r.client == nil {
		return nil, errors.New("internal: redis client is not initialized")
	}

	key, err := r.createKey(username)
	if err != nil {
		return nil, fmt.Errorf("internal: %w", err)
	}

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("not found: user not found in cache")
		}
		return nil, err
	}

	var user User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		_ = r.client.Del(ctx, key)
		return nil, fmt.Errorf("internal: %w", err)
	}

	return &user, nil
}

// Extend extends the expiration time of a cached user
func (r *RedisRepo) Extend(ctx context.Context, username string) error {
	if r.client == nil {
		return errors.New("internal: redis client is not initialized")
	}

	key, err := r.createKey(username)
	if err != nil {
		return fmt.Errorf("internal: %w", err)
	}

	success, err := r.client.Expire(ctx, key, r.config.Expiration).Result()
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("internal: redis: extend failed")
	}
	return nil
}

// Add stores a user in cache
func (r *RedisRepo) Add(ctx context.Context, user *User) error {
	if r.client == nil {
		return errors.New("internal: redis client is not initialized")
	}
	if user == nil {
		return errors.New("validation: user cannot be nil")
	}

	key, err := r.createKey(user.Username)
	if err != nil {
		return fmt.Errorf("internal: %w", err)
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
		return fmt.Errorf("internal: %w", err)
	}

	val, err := json.Marshal(rUser)
	if err != nil {
		return fmt.Errorf("internal: %w", err)
	}

	return r.client.Set(ctx, key, val, r.config.Expiration).Err()
}

// Delete removes a user from cache
func (r *RedisRepo) Delete(ctx context.Context, username string) error {
	if r.client == nil {
		return errors.New("internal: redis client is not initialized")
	}

	key, err := r.createKey(username)
	if err != nil {
		return fmt.Errorf("internal: %w", err)
	}

	deleted, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("internal: %w", err)
	}

	if deleted == 0 {
		return fmt.Errorf("not found: user not found in cache")
	}

	return nil
}

// Exists checks if a user exists in cache
func (r *RedisRepo) Exists(ctx context.Context, username string) (bool, error) {
	if r.client == nil {
		return false, errors.New("internal: redis client is not initialized")
	}

	key, err := r.createKey(username)
	if err != nil {
		return false, fmt.Errorf("internal: create cache key: %w", err)
	}

	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("internal: %w", err)
	}

	return exists > 0, nil
}

// Close gracefully closes the Redis connection
func (r *RedisRepo) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// Stats returns basic Redis statistics
func (r *RedisRepo) Stats(ctx context.Context) (map[string]interface{}, error) {
	if r.client == nil {
		return nil, errors.New("internal: redis client is not initialized")
	}

	poolStats := r.client.PoolStats()
	return map[string]interface{}{
		"total_connections": poolStats.TotalConns,
		"idle_connections":  poolStats.IdleConns,
		"stale_connections": poolStats.StaleConns,
		"hits":              poolStats.Hits,
		"misses":            poolStats.Misses,
		"timeouts":          poolStats.Timeouts,
	}, nil
}
