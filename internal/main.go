package main

import (
	"db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	got, err := s.cassandra.GetUser(r.Context(), user.Username)

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

	_, err = s.cassandra.GetUser(r.Context(), user.Username)
	if err == nil {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	if err := s.cassandra.AddUser(r.Context(), user); err != nil {
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

	w.Write([]byte("User deleted successfully"))
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Server initialization failed: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/get/user", server.handleGetUser)
	mux.HandleFunc("/add/user", server.handleAddUser)
	mux.HandleFunc("/update/user", server.handleUpdateUser)
	mux.HandleFunc("/delete/user", server.handleDeleteUser)

	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
