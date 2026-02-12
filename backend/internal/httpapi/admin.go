package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/example/scalable-learning-platform/backend/internal/models"
)

func registerAdminRoutes(r chi.Router) {
	r.Post("/sessions", handleAdminCreateSession)
	r.Patch("/sessions/{id}/status", handleAdminUpdateSessionStatus)
	r.Get("/sessions/{id}/stats", handleAdminSessionStats)
	r.Get("/submissions", handleAdminSubmissionReports)
	r.Get("/logs", handleAdminLogs)
}

type createSessionRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	VideoURL    string    `json:"video_url"`
}

func handleAdminCreateSession(w http.ResponseWriter, r *http.Request) {
	var req createSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// In a real implementation, persist in DB and return created ID.
	s := models.Session{
		ID:          time.Now().Unix(),
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Status:      models.SessionScheduled,
		VideoURL:    req.VideoURL,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(s)
}

type updateSessionStatusRequest struct {
	Status models.SessionStatus `json:"status"`
}

func handleAdminUpdateSessionStatus(w http.ResponseWriter, r *http.Request) {
	var req updateSessionStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	// Here we would update DB and publish realtime event.
	w.WriteHeader(http.StatusNoContent)
}

type sessionStatsResponse struct {
	ActiveStudents int `json:"active_students"`
}

func handleAdminSessionStats(w http.ResponseWriter, r *http.Request) {
	resp := sessionStatsResponse{
		ActiveStudents: 123, // placeholder, would come from analytics or realtime hub
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func handleAdminSubmissionReports(w http.ResponseWriter, r *http.Request) {
	// Placeholder empty list; real implementation would query DB with filters.
	var submissions []models.Submission
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(submissions)
}

func handleAdminLogs(w http.ResponseWriter, r *http.Request) {
	// Placeholder empty list; real implementation would query logs table.
	var logs []models.ActivityLog
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(logs)
}

