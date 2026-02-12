package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"

	"github.com/example/scalable-learning-platform/backend/internal/config"
	"github.com/example/scalable-learning-platform/backend/internal/middleware"
	"github.com/example/scalable-learning-platform/backend/internal/models"
)

func registerAuthRoutes(r chi.Router, cfg *config.Config, services *Services) {
	r.Post("/register", handleRegister(services))
	r.Post("/login", handleLogin(cfg, services))

	// Authenticated profile endpoint
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))
		r.Get("/me", handleMe(services))
	})
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"`
}

func handleRegister(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req registerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		role := models.RoleStudent
		if req.Role != "" && req.Role == "INSTRUCTOR" {
			role = models.RoleInstructor
		}

		user, err := services.Users.CreateUser(req.Name, req.Email, req.Password, role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(user)
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func handleLogin(cfg *config.Config, services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		user, err := services.Users.GetUserByEmail(req.Email)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		if !services.Users.VerifyPassword(user.Password, req.Password) {
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

		// Don't return password in response
		user.Password = ""

		resp := loginResponse{
			Token: signed,
			User:  user,
		}

		// Best-effort activity log; failures should not break login.
		if services.ActivityLogs != nil {
			_, _ = services.ActivityLogs.CreateLog(&models.ActivityLog{
				UserID:       user.ID,
				Action:       "LOGIN",
				ResourceType: "USER",
				ResourceID:   user.ID,
				Metadata:     "",
			})
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func handleMe(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := middleware.CurrentUser(r)
		if user == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		dbUser, err := services.Users.GetUserByID(user.UserID)
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		dbUser.Password = ""

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(dbUser)
	}
}

