package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"internal/db"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Server struct {
	userRepo    db.UserRepository
	userCache   db.UserCache
	jwtmanager  *JWTManager
	rateLimiter *RateLimiter
}

func NewServer() (*Server, error) {
	// Default credentials if environment variables not set
	cassUsername := getEnvOrDefault("CASS_USERNAME", "backend")
	cassPassword := getEnvOrDefault("CASS_PASSWORD", "BPass0319")
	cassKeyspace := getEnvOrDefault("CASS_KEYSPACE", "cass_keyspace")

	userRepo, err := db.NewCassandraRepo(db.NewCassandraConfig(cassUsername, cassPassword, cassKeyspace))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cassandra: %w", err)
	}

	redisPassword := getEnvOrDefault("REDIS_PASSWORD", "RPass0319")

	userCache, err := db.NewRedisRepo(db.NewRedisConfig(redisPassword))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	jwtSecret := getEnvOrDefault("JWT_SECRET", "some_secret")

	// Rate limiter: 100 requests per minute per IP
	rateLimiter := NewRateLimiter(100, time.Minute)

	return &Server{
		userRepo:    userRepo,
		userCache:   userCache,
		jwtmanager:  NewJWTManager(jwtSecret),
		rateLimiter: rateLimiter,
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (s *Server) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func (s *Server) rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := s.getClientIP(r)

		if !s.rateLimiter.Allow(clientIP) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}

func (s *Server) loginCheck(ctx context.Context, cred *db.Credentials) (*db.User, error) {
	// Try Redis first, fallback to Cassandra
	user, err := s.userCache.Get(ctx, cred.Username)
	if err != nil {
		user, err = s.userRepo.GetUser(ctx, cred.Username)
	}
	if err != nil {
		return nil, err
	}
	if !db.CheckPasswordHash(cred.Password, user.Password) {
		log.Println(cred.Password)
		h, err := db.HashPassword(cred.Password)
		if err != nil {
			panic(err)
		}
		hh, err := db.HashPassword(h)
		if err != nil {
			panic(err)
		}
		log.Printf("hash: %s, rehash: %s; stored: %s", h, hh, user.Password)
		return nil, errors.New("unauthorized: incorrect password")
	}

	s.userCache.Add(ctx, user)
	return user, nil
}

func (s *Server) login(ctx context.Context, cred *db.Credentials) (string, error) {
	_, err := s.loginCheck(ctx, cred)
	if err != nil {
		return "", err
	}

	jwt, err := s.jwtmanager.CreateToken(cred.Username)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func (s *Server) register(ctx context.Context, user *db.User) error {
	ok, err := s.userCache.Exists(ctx, user.Username)
	if err != nil || !ok {
		ok, err = s.userRepo.UsernameExists(ctx, user.Username)
	}
	if err != nil {
		return err
	}
	if ok {
		return errors.New("validation: username exists")
	}

	hashedPassword, err := db.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("internal: %w", err)
	}

	if user.Email == "" {
		return errors.New("validation: invalid email")
	}

	if err := db.ValidUser(user); err != nil {
		return fmt.Errorf("validation: %w", err)
	}

	user.Password = hashedPassword
	if err := s.userRepo.AddUser(ctx, user); err != nil {
		return err
	}

	s.userCache.Add(ctx, user)
	return nil
}

// Helper functions for handling requests
type Payload map[string]interface{}

func parseJSON(r *http.Request) (Payload, error) {
	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func (p Payload) getString(key string) string {
	if val, ok := p[key].(string); ok {
		return val
	}
	return ""
}

func (p Payload) getInt(key string, defaultVal int) int {
	if val, ok := p[key].(int); ok {
		return val
	}
	return defaultVal
}

func (p Payload) credentials() *db.Credentials {
	username := p.getString("username")
	password := p.getString("password")
	if username == "" || password == "" {
		return nil
	}
	return &db.Credentials{Username: username, Password: password}
}

func (p Payload) user() *db.User {
	cred := p.credentials()
	if cred == nil {
		return nil
	}

	return &db.User{
		Credentials: cred,
		Email:       p.getString("email"),
		Category:    2,
	}
}

// Token validation helper
func (s *Server) validateToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("missing or invalid authorization header")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := s.jwtmanager.ValidateToken(tokenStr)
	if err != nil {
		return "", err
	}

	return claims.Username, nil
}

func (s *Server) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // or https://localhost:5173
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
