package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"internal/db"
	"log"
	"net/http"
	"os"
	"strings"
)

type Server struct {
	cassandra  *db.CassandraRepo
	redis      *db.RedisRepo
	jwtmanager *JWTManager
}

func NewServer() (*Server, error) {
	// Default credentials if environment variables not set
	username := getEnvOrDefault("CASS_USERNAME", "backend")
	password := getEnvOrDefault("CASS_PASSWORD", "BPass0319")
	keyspace := getEnvOrDefault("CASS_KEYSPACE", "cass_keyspace")

	cassandra, err := db.NewCassandraRepo(db.NewCassandraConfig(username, password, keyspace))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cassandra: %w", err)
	}

	redisPassword := getEnvOrDefault("REDIS_PASSWORD", "RPass0319")
	redis, err := db.NewRedisRepo(db.NewRedisConfig(redisPassword))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	return &Server{
		cassandra:  cassandra,
		redis:      redis,
		jwtmanager: NewJWTManager(),
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (s *Server) loginCheck(ctx context.Context, cred *db.Credentials) (*db.User, error) {
	// Try Redis first, fallback to Cassandra
	user, err := s.redis.Get(ctx, cred.Username)
	if err != nil && !strings.Contains(err.Error(), "incorrect password") {
		user, err = s.cassandra.GetUser(ctx, cred.Username)
	}
	if err != nil {
		return nil, err
	}
	if !db.CheckPasswordHash(cred.Password, user.Password) {
		return nil, errors.New("unauthorized: incorrect password")
	}

	s.redis.Add(ctx, user)
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
	ok, err := s.redis.Exists(ctx, user.Username)
	if err != nil || !ok {
		ok, err = s.cassandra.UsernameExists(ctx, user.Username)
	}
	if err != nil {
		return err
	}
	if ok {
		return errors.New("validation: username exists")
	}

	if err := s.cassandra.AddUser(ctx, user); err != nil {
		return err
	}

	s.redis.Add(ctx, user)
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
		Category:    p.getInt("category", -1),
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

// HTTP Handlers
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	payload, err := parseJSON(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	cred := payload.credentials()
	if cred == nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	token, err := s.login(r.Context(), cred)
	if err != nil {
		log.Printf("login: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	payload, err := parseJSON(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user := payload.user()
	if user == nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := s.register(r.Context(), user); err != nil {
		log.Printf("register: %v", err)
		if strings.Contains(err.Error(), "validation:") {
			http.Error(w, "Bad request", http.StatusBadRequest)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User added successfully"))
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, err := s.validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	payload, err := parseJSON(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	password := payload.getString("password")
	if password == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := s.loginCheck(r.Context(), db.NewCredentials(username, password))
	if err != nil {
		log.Printf("update: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
		} else if strings.Contains(err.Error(), "validation") {
			http.Error(w, "validation: "+err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}

	newPassword := payload.getString("new_password")
	if newPassword != "" && newPassword == password {
		http.Error(w, "New password cannot be the same as the old one", http.StatusBadRequest)
		return
	}

	email := payload.getString("email")
	if email == "" {
		email = user.Email
	}

	category := payload.getInt("category", user.Category)

	updatedUser := db.NewUser(username, newPassword, email)
	updatedUser.Category = category

	if err := s.cassandra.UpdateUser(r.Context(), updatedUser); err != nil {
		log.Printf("Update user error: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else if strings.Contains(err.Error(), "validation") {
			http.Error(w, "Bad request", http.StatusBadRequest)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	s.redis.Add(r.Context(), user)
	w.Write([]byte("User updated successfully"))
}

func (s *Server) handleGetAdsCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, err := s.validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := s.redis.Get(r.Context(), username)
	if err != nil {
		user, err = s.cassandra.GetUser(r.Context(), username)
	}
	if err != nil {
		log.Printf("Get user error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	s.redis.Extend(r.Context(), username)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"category": user.Category})
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, err := s.validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := s.cassandra.DeleteUser(r.Context(), username); err != nil {
		log.Printf("Delete user error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	s.redis.Delete(r.Context(), username)
	w.Write([]byte("User deleted successfully"))
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rStats, err := s.redis.Stats(r.Context())
	if err != nil {
		log.Printf("Stats error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rStats)
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Server initialization failed: %v", err)
	}

	routes := map[string]http.HandlerFunc{
		"/v1/login":    server.handleLogin,
		"/v1/register": server.handleRegister,
		"/v1/update":   server.handleUpdateUser,
		"/v1/delete":   server.handleDeleteUser,
		"/v1/stats":    server.handleStats,
		"/v1/get_ads":  server.handleGetAdsCategory,
	}

	mux := http.NewServeMux()
	for path, handler := range routes {
		mux.HandleFunc(path, handler)
	}

	port := ":8443"
	certDir := getEnvOrDefault("TLS_CERT_DIR", "../certs")

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServeTLS(port, certDir+"/server.crt", certDir+"/server.key", mux))
}
