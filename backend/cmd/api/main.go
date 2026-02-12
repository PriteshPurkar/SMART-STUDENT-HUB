package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/example/scalable-learning-platform/backend/internal/config"
	"github.com/example/scalable-learning-platform/backend/internal/httpapi"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	apiRouter, err := httpapi.NewRouter(cfg)
	if err != nil {
		log.Fatalf("init router: %v", err)
	}

	r.Mount("/api/v1", apiRouter)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("API listening on :%s", cfg.Port)
	log.Printf("Environment: %s", cfg.Environment)
	log.Printf("Database: %s", maskDatabaseURL(cfg.DatabaseURL))

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("server error: %v", err)
		os.Exit(1)
	}
}

// maskDatabaseURL masks sensitive information in database URL for logging
func maskDatabaseURL(url string) string {
	if len(url) == 0 {
		return "(not configured)"
	}
	// Just show that it's configured without showing credentials
	return "postgres://***:***@***:***/**?"
}
}

