package main

import (
	"encoding/json"
	"fmt"
	"internal/db"
	"log"
	"net/http"
	"strings"
)

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
			http.Error(w, "Unauthorized"+err.Error(), http.StatusUnauthorized)
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
		http.Error(w, "Bad request1", http.StatusBadRequest)
		return
	}

	user := payload.user()
	if user == nil {
		http.Error(w, "Bad request2", http.StatusBadRequest)
		return
	}

	if err := s.register(r.Context(), user); err != nil {
		log.Printf("register: %v", err)
		if strings.Contains(err.Error(), "validation:") {
			http.Error(w, "Bad request"+err.Error(), http.StatusBadRequest)
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
			http.Error(w, "update: "+err.Error(), http.StatusUnauthorized)
		} else if strings.Contains(err.Error(), "validation") {
			http.Error(w, "update: "+err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
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

	hashedPassword, err := db.HashPassword(updatedUser.Password)
	if err != nil {
		http.Error(w, "internal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	updatedUser.Password = hashedPassword

	if err := s.userRepo.UpdateUser(r.Context(), updatedUser); err != nil {
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

	s.userCache.Add(r.Context(), updatedUser)
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

	user, err := s.userCache.Get(r.Context(), username)
	if err != nil {
		user, err = s.userRepo.GetUser(r.Context(), username)
	}
	if err != nil {
		log.Printf("Get user error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	s.userCache.Extend(r.Context(), username)

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

	if err := s.userRepo.DeleteUser(r.Context(), username); err != nil {
		log.Printf("Delete user error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	s.userCache.Delete(r.Context(), username)
	w.Write([]byte("User deleted successfully"))
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rStats, err := s.userCache.Stats(r.Context())
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
		"/v1/login":    server.corsMiddleware(server.rateLimitMiddleware(server.handleLogin)),
		"/v1/register": server.corsMiddleware(server.rateLimitMiddleware(server.handleRegister)),
		"/v1/update":   server.corsMiddleware(server.rateLimitMiddleware(server.handleUpdateUser)),
		"/v1/delete":   server.corsMiddleware(server.rateLimitMiddleware(server.handleDeleteUser)),
		"/v1/stats":    server.corsMiddleware(server.rateLimitMiddleware(server.handleStats)),
		"/v1/get_ads":  server.corsMiddleware(server.rateLimitMiddleware(server.handleGetAdsCategory)),
	}

	mux := http.NewServeMux()
	for path, handler := range routes {
		mux.HandleFunc(path, handler)
	}

	port := ":8443"
	certDir := getEnvOrDefault("TLS_CERT_DIR", "../certs")

	fmt.Printf("Server starting on port %s with rate limiting (100 req/min per IP)\n", port)
	log.Fatal(http.ListenAndServeTLS(port, certDir+"/server.crt", certDir+"/server.key", mux))
}
