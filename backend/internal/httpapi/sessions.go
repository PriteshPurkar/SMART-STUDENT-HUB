package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/example/scalable-learning-platform/backend/internal/models"
)

func registerSessionRoutes(r chi.Router) {
	r.Get("/", handleListSessions)
	r.Get("/{id}", handleGetSession)
	r.Get("/{id}/status", handleGetSessionStatus)
	r.Get("/{id}/materials", handleGetSessionMaterials)
}

func handleListSessions(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	sessions := []models.Session{
		{
			ID:          1,
			Title:       "Live Exam Prep",
			Description: "Final preparation for exams",
			StartTime:   now.Add(30 * time.Minute),
			EndTime:     now.Add(90 * time.Minute),
			Status:      models.SessionScheduled,
			VideoURL:    "https://video.example.com/session/1",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sessions)
}

func handleGetSession(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	now := time.Now()
	s := models.Session{
		ID:          id,
		Title:       "Live Exam Prep",
		Description: "Details for session",
		StartTime:   now.Add(30 * time.Minute),
		EndTime:     now.Add(90 * time.Minute),
		Status:      models.SessionScheduled,
		VideoURL:    "https://video.example.com/session/1",
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s)
}

func handleGetSessionStatus(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status models.SessionStatus `json:"status"`
	}{
		Status: models.SessionScheduled,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

func handleGetSessionMaterials(w http.ResponseWriter, r *http.Request) {
	materials := []models.StudyMaterial{
		{
			ID:        1,
			Title:     "Exam Syllabus",
			Type:      models.MaterialPDF,
			S3Key:     "materials/exam-syllabus.pdf",
			URL:       "https://cdn.example.com/materials/exam-syllabus.pdf",
			UploadedBy: 100,
			CreatedAt: time.Now(),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(materials)
}

