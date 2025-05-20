package main

import (
	"context"
	"db"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Server struct {
	cassandra *db.CassandraRepo
	redis     *db.RedisRepo
}

func NewServer() (*Server, error) {
	username := os.Getenv("cass_username")
	password := os.Getenv("cass_password")
	keyspace := os.Getenv("cass_keyspace")

	if username == "" || password == "" || keyspace == "" {
		username = "backend"
		password = "BPass0319"
		keyspace = "cass_keyspace"
	}

	cConfig := db.NewCassandraConfig(username, password, keyspace)
	cassandra, err := db.NewCassandraRepo(cConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cassandra: %w", err)
	}

	password = os.Getenv("redis_password")
	if password == "" {
		password = "RPass0319"
	}

	rConfig := db.NewRedisConfig(password)
	redis, err := db.NewRedisRepo(rConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	return &Server{
		cassandra: cassandra,
		redis:     redis,
	}, nil
}

func (s *Server) login(ctx context.Context, cred *db.Credentials) (*db.User, error) {
	got, err := s.redis.Get(ctx, cred)
	if err != nil && !strings.Contains(err.Error(), "incorrect password") {
		got, err = s.cassandra.GetUser(ctx, cred)
	}

	if err != nil {
		return nil, fmt.Errorf("server get: %w", err)
	}

	s.redis.Add(ctx, got)
	return got, nil
}

func (s *Server) register(ctx context.Context, user *db.User) error {
	ok, err := s.cassandra.UsernameExists(ctx, user.Username)
	if err != nil {
		return err
	}
	if ok {
		return errors.New("username exists")
	}

	err = s.cassandra.AddUser(ctx, user)
	if err != nil {
		return fmt.Errorf("server post: %w", err)
	}

	s.redis.Add(ctx, user)
	return nil
}

type Payload map[string]interface{}

func parseJSON(w http.ResponseWriter, r *http.Request) (Payload, error) {
	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return nil, err
	}
	return payload, nil
}

func (p Payload) credentials() *db.Credentials {
	username, ok := p["username"].(string)
	if !ok {
		return nil
	}

	password, ok := p["password"].(string)
	if !ok {
		return nil
	}

	return &db.Credentials{
		Username: username,
		Password: password,
	}
}

func (p Payload) user() *db.User {
	email, ok := p["email"].(string)
	if !ok {
		email = ""
	}

	cred := p.credentials()
	if cred == nil {
		return nil
	}

	category, ok := p["category"].(int)
	if !ok {
		category = -1
	}

	return &db.User{
		Credentials: cred,
		Email:       email,
		Category:    category,
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	payload, err := parseJSON(w, r)
	if err != nil {
		return
	}

	cred := payload.credentials()
	if cred == nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err = s.login(r.Context(), cred)
	if err != nil {
		http.Error(w, "Invalid user data", http.StatusNotFound)
		return
	}

	w.Write([]byte("token"))
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	payload, err := parseJSON(w, r)
	if err != nil {
		return
	}

	user := payload.user()
	if user == nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := s.register(r.Context(), user); err != nil {
		http.Error(w, "Failed to add user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User added successfully"))
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	payload, err := parseJSON(w, r)
	if err != nil {
		return
	}

	password, ok := payload["old_password"].(string)
	if !ok {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user := payload.user()
	if user == nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if user.Email == "" && user.Category == -1 && password == user.Password {
		http.Error(w, "New password cannot be the same as the old one", http.StatusBadRequest)
		return
	}

	if err := s.cassandra.UpdateUser(r.Context(), user, password); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	s.redis.Add(r.Context(), user)
	w.Write([]byte("User updated successfully"))
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	payload, err := parseJSON(w, r)
	if err != nil {
		return
	}

	username, ok := payload["username"].(string)
	if !ok {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := s.cassandra.DeleteUser(r.Context(), username); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	s.redis.Delete(r.Context(), username)
	w.Write([]byte("User deleted successfully"))
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	rStats, err := s.redis.Stats(r.Context())
	if err != nil {
		http.Error(w, "Failed to get cache stats", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rStats)
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Server initialization failed: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/login", server.handleLogin)
	mux.HandleFunc("/v1/register", server.handleRegister)
	mux.HandleFunc("/v1/update", server.handleUpdateUser)
	mux.HandleFunc("/v1/delete", server.handleDeleteUser)
	mux.HandleFunc("/v1/stats", server.handleStats)

	port := ":8443"
	fmt.Printf("Server starting on port %s\n", port)

	certDir := os.Getenv("tls_cert_dir")
	if certDir == "" {
		certDir = "../certs"
	}

	if err := http.ListenAndServeTLS(port, certDir+"/server.crt", certDir+"/server.key", mux); err != nil {
		log.Fatalf("server listen: %v", err)
	}
}
