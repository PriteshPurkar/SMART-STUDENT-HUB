package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/example/scalable-learning-platform/backend/internal/models"
)

func registerSessionRoutes(r chi.Router, services *Services) {
	r.Get("/", handleListSessions(services))
	r.Get("/{id}", handleGetSession(services))
	r.Get("/{id}/status", handleGetSessionStatus(services))
	r.Get("/{id}/materials", handleGetSessionMaterials(services))
}

func handleListSessions(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessions, err := services.Sessions.GetAllSessions()
		if err != nil {
			http.Error(w, "failed to fetch sessions", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if sessions == nil {
			sessions = []models.Session{}
		}
		_ = json.NewEncoder(w).Encode(sessions)
	}
}

func handleGetSession(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid session id", http.StatusBadRequest)
			return
		}

		session, err := services.Sessions.GetSessionByID(id)
		if err != nil {
			http.Error(w, "session not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(session)
	}
}

func handleGetSessionStatus(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid session id", http.StatusBadRequest)
			return
		}

		session, err := services.Sessions.GetSessionByID(id)
		if err != nil {
			http.Error(w, "session not found", http.StatusNotFound)
			return
		}

		status := struct {
			Status interface{} `json:"status"`
		}{
			Status: session.Status,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(status)
	}
}

func handleGetSessionMaterials(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid session id", http.StatusBadRequest)
			return
		}

		materials, err := services.Materials.GetMaterialsBySessionID(id)
		if err != nil {
			http.Error(w, "failed to fetch materials", http.StatusInternalServerError)
			return
		}

		// Convert to interface{} slice for consistency with other handlers
		var asInterfaces []interface{}
		for _, m := range materials {
			asInterfaces = append(asInterfaces, m)
		}

		w.Header().Set("Content-Type", "application/json")
		if asInterfaces == nil {
			asInterfaces = []interface{}{}
		}
		_ = json.NewEncoder(w).Encode(asInterfaces)
	}
}

