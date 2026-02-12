package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/example/scalable-learning-platform/backend/internal/config"
	"github.com/example/scalable-learning-platform/backend/internal/middleware"
)

// NewRouter wires all public API routes under /api/v1.
func NewRouter(cfg *config.Config) (http.Handler, error) {
	// Initialize services
	services, err := InitServices(cfg)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Auth endpoints
	r.Route("/auth", func(r chi.Router) {
		registerAuthRoutes(r, cfg, services)
	})

	// Student endpoints
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))
		r.Route("/student", func(r chi.Router) {
			registerStudentRoutes(r, services)
		})
		r.Route("/sessions", func(r chi.Router) {
			registerSessionRoutes(r, services)
		})
		r.Route("/exams", func(r chi.Router) {
			registerExamRoutes(r, services)
		})
	})

	// Admin / Instructor endpoints
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))
		r.Use(middleware.RequireAdminOrInstructor())
		r.Route("/admin", func(r chi.Router) {
			registerAdminRoutes(r, services)
		})
	})

	// Realtime
	r.Route("/realtime", func(r chi.Router) {
		registerRealtimeRoutes(r)
	})

	return r, nil
}

