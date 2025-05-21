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

	jwtManager := NewJWTManager()

	return &Server{
		cassandra:  cassandra,
		redis:      redis,
		jwtmanager: jwtManager,
	}, nil
}

func (s *Server) login(ctx context.Context, cred *db.Credentials) (string, error) {
	got, err := s.redis.Get(ctx, cred)
	if err != nil && !strings.Contains(err.Error(), "incorrect password") {
		got, err = s.cassandra.GetUser(ctx, cred)
	}

	if err != nil {
		return "", fmt.Errorf("server get: %w", err)
	}

	jwt, err := s.jwtmanager.CreateToken(got.Username)

	if err != nil {
		return "", fmt.Errorf("server jwt: %w", err)
	}

	s.redis.Add(ctx, got)
	return jwt, nil
}

func (s *Server) register(ctx context.Context, user *db.User) error {
	ok, err := s.cassandra.UsernameExists(ctx, user.Username)
	if err != nil {
		return err
	}
	if ok {
		return errors.New("validation: username exists")
	}

	err = s.cassandra.AddUser(ctx, user)
	if err != nil {
		return fmt.Errorf("server post: %w", err)
	}

	s.redis.Add(ctx, user)
	return nil
}

type Payload map[string]interface{}

func parseJSON(r *http.Request) (Payload, error) {
	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("json: %w", err)
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

func (p Payload) userData(skipCred bool) *db.User {
	email, ok := p["email"].(string)
	if !ok {
		email = ""
	}

	cred := &db.Credentials{}
	if !skipCred {
		cred = p.credentials()
		if cred == nil {
			return nil
		}
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

func (p Payload) user() *db.User {
	return p.userData(false)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "login: invalid method", http.StatusMethodNotAllowed)
		return
	}

	payload, err := parseJSON(r)
	if err != nil {
		http.Error(w, "login: "+err.Error(), http.StatusBadRequest)
		return
	}

	cred := payload.credentials()
	if cred == nil {
		http.Error(w, "login: invalid request", http.StatusBadRequest)
		return
	}

	token, err := s.login(r.Context(), cred)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, "login: "+err.Error(), http.StatusUnauthorized)
		} else {
			fmt.Println(err)
			http.Error(w, "Could not log in", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"token": token}
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "register: invalid method", http.StatusMethodNotAllowed)
		return
	}

	payload, err := parseJSON(r)
	if err != nil {
		http.Error(w, "register: "+err.Error(), http.StatusBadRequest)
		return
	}

	user := payload.user()
	if user == nil {
		http.Error(w, "register: invalid request", http.StatusBadRequest)
		return
	}

	if err := s.register(r.Context(), user); err != nil {
		if strings.Contains(err.Error(), "validation:") {
			http.Error(w, "register: "+err.Error(), http.StatusBadRequest)
		} else {
			fmt.Println(err)
			http.Error(w, "Failed to add user", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User added successfully"))
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "update: invalid method", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := s.jwtmanager.ValidateToken(tokenStr)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	payload, err := parseJSON(r)
	if err != nil {
		http.Error(w, "update: "+err.Error(), http.StatusBadRequest)
		return
	}

	password, ok := payload["password"].(string)
	if !ok {
		http.Error(w, "update: invalid password", http.StatusUnauthorized)
		return
	}

	user, err := s.cassandra.GetUser(r.Context(), db.NewCredentials(claims.Username, password))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	if !db.CheckPasswordHash(password, user.Password) {
		http.Error(w, "unauthorized: incorrect password", http.StatusUnauthorized)
		return
	}

	userData := payload.userData(true)
	userData.Credentials = db.NewCredentials(user.Username, password)

	newPassword, ok := payload["new_password"].(string)

	if ok {
		userData.Password = newPassword
		if newPassword == password {
			http.Error(w, "New password cannot be the same as the old one", http.StatusBadRequest)
			return
		}
	}

	if userData.Email == "" {
		userData.Email = user.Email
	}
	if userData.Category == -1 {
		userData.Category = user.Category
	}

	if err := s.cassandra.UpdateUser(r.Context(), db.NewUser(claims.Username, newPassword, user.Email)); err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, "update: "+err.Error(), http.StatusUnauthorized)
		} else if strings.Contains(err.Error(), "validation") {
			http.Error(w, "update: "+err.Error(), http.StatusBadRequest)
		} else {
			fmt.Println(err)
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
		}
		return
	}

	s.redis.Add(r.Context(), user)
	w.Write([]byte("User updated successfully"))
}

func (s *Server) handleGetAdsCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := s.jwtmanager.ValidateToken(tokenStr)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	log.Printf("Authenticated user: %s", claims.Username)

	user, err := s.cassandra.GetUser(r.Context(), &db.Credentials{Username: claims.Username, Password: ""})

	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	category := user.Category
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]int{"category": category}
	json.NewEncoder(w).Encode(response)
	w.Write([]byte("User category retrieved successfully"))
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "delete: invalid method", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := s.jwtmanager.ValidateToken(tokenStr)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	if err := s.cassandra.DeleteUser(r.Context(), claims.Username); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	s.redis.Delete(r.Context(), claims.Username)
	w.Write([]byte("User deleted successfully"))
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	rStats, err := s.redis.Stats(r.Context())
	if err != nil {
		fmt.Println(err)
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
	mux.HandleFunc("/v1/get_ads", server.handleGetAdsCategory)

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
