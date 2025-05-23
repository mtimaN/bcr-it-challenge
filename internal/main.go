package main

import (
	"encoding/json"
	"fmt"
	"internal/db"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func JSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// HTTP Handlers
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	payload, err := parseJSON(r)
	if err != nil {
		s.recordDBOperation("user_login", "error")
		JSONError(w, "Bad request: invalid json", http.StatusBadRequest)
		return
	}

	cred := payload.credentials()
	if cred == nil {
		s.recordDBOperation("user_login", "error")
		JSONError(w, "Bad request: invalid credentials", http.StatusBadRequest)
		return
	}

	token, err := s.login(r.Context(), cred)
	if err != nil {
		s.recordDBOperation("user_login", "error")
		log.Printf("login: %v", err)

		if strings.Contains("login: "+err.Error(), "unauthorized") {
			JSONError(w, err.Error(), http.StatusUnauthorized)
		} else {
			JSONError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	s.recordDBOperation("user_login", "success")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	payload, err := parseJSON(r)
	if err != nil {
		s.recordDBOperation("user_register", "error")
		JSONError(w, "Bad request: invalid json", http.StatusBadRequest)
		return
	}

	user := payload.user()
	if user == nil {
		s.recordDBOperation("user_register", "error")
		JSONError(w, "Bad request: invalid user data", http.StatusBadRequest)
		return
	}

	if err := s.register(r.Context(), user); err != nil {
		s.recordDBOperation("user_register", "error")
		log.Printf("register: %v", err)

		if strings.Contains(err.Error(), "validation:") {
			JSONError(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		} else {
			JSONError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	s.recordDBOperation("user_register", "success")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User added successfully"))
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, err := s.validateToken(r)
	if err != nil {
		s.recordDBOperation("user_update", "error")
		JSONError(w, "Unauthorized: invalid token", http.StatusUnauthorized)
		return
	}

	payload, err := parseJSON(r)
	if err != nil {
		s.recordDBOperation("user_update", "error")
		JSONError(w, "Bad request: invalid json", http.StatusBadRequest)
		return
	}

	password := payload.getString("password")
	if password == "" {
		s.recordDBOperation("user_update", "error")
		JSONError(w, "Unauthorized: password not provided", http.StatusUnauthorized)
		return
	}

	user, err := s.loginCheck(r.Context(), db.NewCredentials(username, password))
	if err != nil {
		s.recordDBOperation("user_update", "error")
		log.Printf("update: %v", err)

		if strings.Contains(err.Error(), "unauthorized") {
			JSONError(w, "update: "+err.Error(), http.StatusUnauthorized)
		} else if strings.Contains(err.Error(), "validation") {
			JSONError(w, "update: "+err.Error(), http.StatusBadRequest)
		} else {
			JSONError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	newPassword := payload.getString("new_password")
	if newPassword != "" && newPassword == password {
		s.recordDBOperation("user_update", "error")
		JSONError(w, "New password cannot be the same as the old one", http.StatusBadRequest)
		return
	}

	email := payload.getString("email")
	if email == "" {
		email = user.Email
	}

	updatedUser := db.NewUser(username, newPassword, email)
	updatedUser.Category = user.Category

	hashedPassword, err := db.HashPassword(updatedUser.Password)
	if err != nil {
		s.recordDBOperation("user_update", "error")
		JSONError(w, "internal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	updatedUser.Password = hashedPassword

	if err := s.userRepo.UpdateUser(r.Context(), updatedUser); err != nil {
		s.recordDBOperation("user_update", "error")
		log.Printf("Update user error: %v", err)

		if strings.Contains(err.Error(), "unauthorized") {
			JSONError(w, "update: "+err.Error(), http.StatusUnauthorized)
		} else if strings.Contains(err.Error(), "validation") {
			JSONError(w, "Bad request", http.StatusBadRequest)
		} else {
			JSONError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	s.recordDBOperation("user_update", "success")
	s.userCache.Add(r.Context(), updatedUser)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("User updated successfully"))
}

func (s *Server) handleGetAdsCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, err := s.validateToken(r)
	if err != nil {
		JSONError(w, "Unauthorized: invalid token", http.StatusUnauthorized)
		return
	}

	user, err := s.userCache.Get(r.Context(), username)
	if err != nil {
		user, err = s.userRepo.GetUser(r.Context(), username)
	}
	if err != nil {
		log.Printf("Get user error: %v", err)
		JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	s.userCache.Extend(r.Context(), username)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"category": user.Category})
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		JSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, err := s.validateToken(r)
	if err != nil {
		s.recordDBOperation("user_delete", "error")
		JSONError(w, "Unauthorized: invalid token", http.StatusUnauthorized)
		return
	}

	if err := s.userRepo.DeleteUser(r.Context(), username); err != nil {
		s.recordDBOperation("user_delete", "error")
		log.Printf("Delete user error: %v", err)
		JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	s.recordDBOperation("user_delete", "success")
	s.userCache.Delete(r.Context(), username)
	w.Write([]byte("User deleted successfully"))
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rStats, err := s.userCache.Stats(r.Context())
	if err != nil {
		log.Printf("Stats error: %v", err)
		JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rStats)
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Server initialization failed: %v", err)
	}

	middleware := []func(http.HandlerFunc) http.HandlerFunc{
		server.corsMiddleware,
		server.metricsMiddleware,
		server.rateLimitMiddleware,
	}

	routes := map[string]http.HandlerFunc{
		"/v1/login":    server.handleLogin,
		"/v1/register": server.handleRegister,
		"/v1/update":   server.handleUpdateUser,
		"/v1/delete":   server.handleDeleteUser,
		"/v1/stats":    server.handleStats,
		"/v1/get_ads":  server.handleGetAdsCategory,
	}

	for path, handler := range routes {
		for _, mware := range middleware {
			handler = mware(handler)
		}
		routes[path] = handler
	}

	mux := http.NewServeMux()
	for path, handler := range routes {
		mux.HandleFunc(path, handler)
	}

	certDir := getEnvOrDefault("TLS_CERT_DIR", "../certs")

	go func() {
		metricsMux := http.NewServeMux()
		metricsMux.Handle("/metrics", promhttp.Handler())

		fmt.Println("Metrics server starting on port :8080")
		log.Fatal(http.ListenAndServeTLS("0.0.0.0:8080", certDir+"/server.crt", certDir+"/server.key", metricsMux))
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			server.updateActiveUsersMetric()
		}
	}()

	port := ":8443"

	fmt.Printf("Server starting on port %s with rate limiting (100 req/min per IP)\n", port)
	log.Fatal(http.ListenAndServeTLS(port, certDir+"/server.crt", certDir+"/server.key", mux))
}
