package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/example/scalable-learning-platform/backend/internal/middleware"
	"github.com/example/scalable-learning-platform/backend/internal/models"
)

func registerStudentRoutes(r chi.Router) {
	r.Get("/dashboard", handleStudentDashboard)
}

type dashboardResponse struct {
	UpcomingSessions []models.Session `json:"upcoming_sessions"`
	PastSessions     []models.Session `json:"past_sessions"`
	Notifications    []models.Notification `json:"notifications"`
}

func handleStudentDashboard(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	now := time.Now()
	// For the scaffold, return dummy data. In a real system this would query Postgres.
	resp := dashboardResponse{
		UpcomingSessions: []models.Session{
			{
				ID:          1,
				Title:       "Live Exam Prep",
				Description: "Final preparation for exams",
				StartTime:   now.Add(30 * time.Minute),
				EndTime:     now.Add(90 * time.Minute),
				Status:      models.SessionScheduled,
				VideoURL:    "https://video.example.com/session/1",
				CreatedBy:   100,
			},
		},
		PastSessions: []models.Session{
			{
				ID:          2,
				Title:       "Recorded Lecture: System Design",
				Description: "Scalable platform design overview",
				StartTime:   now.Add(-48 * time.Hour),
				EndTime:     now.Add(-47 * time.Hour),
				Status:      models.SessionCompleted,
				VideoURL:    "https://video.example.com/session/2",
				CreatedBy:   101,
			},
		},
		Notifications: []models.Notification{
			{
				ID:        1,
				UserID:    user.UserID,
				Type:      "SESSION_STARTING_SOON",
				Message:   "Your exam prep session starts in 30 minutes.",
				IsRead:    false,
				CreatedAt: now,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

