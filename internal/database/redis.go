package db

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConfig holds configuration for Redis connection
type RedisConfig struct {
	*redis.Options

	Expiration time.Duration
}

// RedisClient is a wrapper around the Redis client
type RedisClient struct {
	client *redis.Client
	config *RedisConfig
}

func NewRedisConfig(password string) *RedisConfig {
	return &RedisConfig{
		Options: &redis.Options{
			Addr:        "localhost:6379",
			Password:    password,
			DialTimeout: 5 * time.Second,
		},
		Expiration: 24 * time.Hour,
	}
}

func NewRedisClient(config *RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), config.DialTimeout)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("connect to redis: %w", err)
	}

	return &RedisClient{
		client: client,
		config: config,
	}, nil
}

func (c *RedisClient) Health(ctx context.Context) error {
	if c.client == nil {
		return errors.New("redis client is nil")
	}

	if _, err := c.client.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("redis health: %w", err)
	}

	return nil
}

func createKey(username string) string {
	h := sha256.New()
	h.Write([]byte(username))
	return hex.EncodeToString(h.Sum(nil))
}

func (c *RedisClient) GetUser(ctx context.Context, username string) (*User, error) {
	key := createKey(username)

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis get: %w", err)
	}

	var user User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, fmt.Errorf("redis get: %w", err)
	}

	return &user, nil
}

func (c *RedisClient) AddUser(ctx context.Context, user *User) error {
	key := createKey(user.Username)
	val, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("redis add: %w", err)
	}

	err = c.client.Set(ctx, key, val, c.config.Expiration).Err()
	if err != nil {
		return fmt.Errorf("redis add: %w", err)
	}
	
	return nil
}

func (c *RedisClient) Close() error {
	return c.client.Close()
}
