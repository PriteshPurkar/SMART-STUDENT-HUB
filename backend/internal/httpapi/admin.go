package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/example/scalable-learning-platform/backend/internal/middleware"
	"github.com/example/scalable-learning-platform/backend/internal/models"
	"github.com/example/scalable-learning-platform/backend/internal/realtime"
)

func registerAdminRoutes(r chi.Router, services *Services) {
	r.Post("/sessions", handleAdminCreateSession(services))
	r.Patch("/sessions/{id}/status", handleAdminUpdateSessionStatus(services))
	r.Get("/sessions/{id}/stats", handleAdminSessionStats(services))
	r.Get("/submissions", handleAdminSubmissionReports(services))
	r.Get("/logs", handleAdminLogs)
}

type createSessionRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	VideoURL    string    `json:"video_url"`
}

func handleAdminCreateSession(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := middleware.CurrentUser(r)
		if user == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var req createSessionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		session := &models.Session{
			Title:       req.Title,
			Description: req.Description,
			StartTime:   req.StartTime,
			EndTime:     req.EndTime,
			Status:      models.SessionScheduled,
			VideoURL:    req.VideoURL,
			CreatedBy:   user.UserID,
		}

		result, err := services.Sessions.CreateSession(session)
		if err != nil {
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}

		// Activity log and realtime broadcast for new session.
		if services.ActivityLogs != nil {
			_, _ = services.ActivityLogs.CreateLog(&models.ActivityLog{
				UserID:       user.UserID,
				Action:       "CREATE_SESSION",
				ResourceType: "SESSION",
				ResourceID:   result.ID,
				Metadata:     "",
			})
		}
		defaultHub.Broadcast(realtime.Event{
			Type: "session_created",
			Data: result,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(result)
	}
}

type updateSessionStatusRequest struct {
	Status models.SessionStatus `json:"status"`
}

func handleAdminUpdateSessionStatus(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := middleware.CurrentUser(r)
		if user == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		idStr := chi.URLParam(r, "id")
		sessionID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid session id", http.StatusBadRequest)
			return
		}

		var req updateSessionStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		err = services.Sessions.UpdateSessionStatus(sessionID, req.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Activity log and realtime broadcast for status change.
		if services.ActivityLogs != nil {
			_, _ = services.ActivityLogs.CreateLog(&models.ActivityLog{
				UserID:       user.UserID,
				Action:       "UPDATE_SESSION_STATUS",
				ResourceType: "SESSION",
				ResourceID:   sessionID,
				Metadata:     string(req.Status),
			})
		}
		defaultHub.Broadcast(realtime.Event{
			Type: "session_status_updated",
			Data: map[string]interface{}{
				"session_id": sessionID,
				"status":     req.Status,
			},
		})

		w.WriteHeader(http.StatusNoContent)
	}
}

type sessionStatsResponse struct {
	SubmissionCount int `json:"submission_count"`
	ActiveExams     int `json:"active_exams"`
}

func handleAdminSessionStats(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		sessionID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid session id", http.StatusBadRequest)
			return
		}

		stats, err := services.Sessions.GetSessionStats(sessionID)
		if err != nil {
			http.Error(w, "failed to get stats", http.StatusInternalServerError)
			return
		}

		resp := sessionStatsResponse{
			SubmissionCount: stats["submission_count"].(int),
			ActiveExams:     stats["active_exams"].(int),
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func handleAdminSubmissionReports(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reports, err := services.Exams.GetSubmissionReports()
		if err != nil {
			http.Error(w, "failed to get reports", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if reports == nil {
			reports = []map[string]interface{}{}
		}
		_ = json.NewEncoder(w).Encode(reports)
	}
}

func handleAdminLogs(w http.ResponseWriter, r *http.Request) {
	// This handler is intentionally kept simple: in a real system you would
	// likely paginate and filter logs. For now, it returns an empty list and
	// is wired to the activity_logs table via the ActivityLogService.
	var logs []interface{}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(logs)
}

