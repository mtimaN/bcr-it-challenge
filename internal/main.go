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
)

type Server struct {
	cassandra *db.CassandraClient
	redis     *db.RedisClient
}

type Payload struct {
	User        db.User `json:"user"`
	OldPassword string  `json:"old_password"`
}

func NewServer() (*Server, error) {
	cConfig := db.NewCassandraConfig("backend", "BPass0319", "cass_keyspace")
	cassandra, err := db.NewCassandraClient(cConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cassandra: %w", err)
	}

	rConfig := db.NewRedisConfig("RPass0319")
	redis, err := db.NewRedisClient(rConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	return &Server{
		cassandra: cassandra,
		redis:     redis,
	}, nil
}

func (s *Server) getUser(ctx context.Context, user *db.User) (*db.User, error) {
	got, err := s.redis.Get(ctx, user)
	if err != nil {
		got, err = s.cassandra.GetUser(ctx, user.Username)
	}

	if err != nil {
		return nil, fmt.Errorf("server get: %w", err)
	}

	s.redis.Add(ctx, got)
	return got, nil
}

func (s *Server) addUser(ctx context.Context, user *db.User) error {
	_, err := s.getUser(ctx, user)
	if err == nil {
		return errors.New("server post: user already exists")
	}

	err = s.cassandra.AddUser(ctx, user)
	if err != nil {
		return fmt.Errorf("server post: %w", err)
	}

	s.redis.Add(ctx, user)
	return nil
}

func parseJSONPayload(w http.ResponseWriter, r *http.Request) (*Payload, error) {
	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return nil, err
	}
	return &payload, nil
}

func payloadToUser(payload *Payload) *db.User {
	user := db.NewUser(payload.User.Username, payload.User.Password, payload.User.Email)
	user.Category = payload.User.Category

	return user
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	payload, err := parseJSONPayload(w, r)
	if err != nil {
		return
	}

	user := payloadToUser(payload)
	got, err := s.getUser(r.Context(), user)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.Password != got.Password {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (s *Server) handleAddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	payload, err := parseJSONPayload(w, r)
	if err != nil {
		return
	}

	user := payloadToUser(payload)
	if err := s.addUser(r.Context(), user); err != nil {
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

	payload, err := parseJSONPayload(w, r)
	if err != nil {
		return
	}

	user := payloadToUser(payload)

	if err := s.cassandra.UpdateUser(r.Context(), user, payload.OldPassword); err != nil {
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

	payload, err := parseJSONPayload(w, r)
	if err != nil {
		return
	}

	user := payloadToUser(payload)

	if err := s.cassandra.DeleteUser(r.Context(), user); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	s.redis.Delete(r.Context(), user.Username)
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
	mux.HandleFunc("/v1/get", server.handleGetUser)
	mux.HandleFunc("/v1/add", server.handleAddUser)
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
