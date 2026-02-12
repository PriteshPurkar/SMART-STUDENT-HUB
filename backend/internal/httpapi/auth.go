package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"

	"github.com/example/scalable-learning-platform/backend/internal/config"
	"github.com/example/scalable-learning-platform/backend/internal/models"
)

type authService interface {
	RegisterStudent(name, email, password string) (*models.User, error)
	Login(email, password string) (*models.User, error)
	GetUserByID(id int64) (*models.User, error)
}

// In this scaffold we use an in-memory auth service implementation.
var defaultAuthService authService = newInMemoryAuthService()

func registerAuthRoutes(r chi.Router, cfg *config.Config) {
	r.Post("/register", handleRegister)
	r.Post("/login", handleLogin(cfg))
	r.Get("/me", handleMe)
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	user, err := defaultAuthService.RegisterStudent(req.Name, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func handleLogin(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}
		user, err := defaultAuthService.Login(req.Email, req.Password)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		claims := &jwt.MapClaims{
			"user_id": user.ID,
			"role":    user.Role,
			"exp":     time.Now().Add(2 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			http.Error(w, "failed to sign token", http.StatusInternalServerError)
			return
		}

		resp := loginResponse{
			Token: signed,
			User:  user,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func handleMe(w http.ResponseWriter, r *http.Request) {
	// For simplicity in this scaffold, we just return 501 if not wired with middleware.
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = w.Write([]byte(`{"message":"not implemented in scaffold"}`))
}

