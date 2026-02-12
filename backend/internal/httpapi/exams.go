package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/example/scalable-learning-platform/backend/internal/middleware"
	"github.com/example/scalable-learning-platform/backend/internal/models"
)

func registerExamRoutes(r chi.Router) {
	r.Post("/{id}/submissions", handleSubmitExam)
	r.Get("/{id}/submissions/me", handleGetMySubmission)
}

type submitExamRequest struct {
	Answers string `json:"answers"` // simplified; could be structured
}

type submitExamResponse struct {
	SubmissionID int64  `json:"submission_id"`
	Message      string `json:"message"`
}

func handleSubmitExam(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	if user == nil || user.Role != models.RoleStudent {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var req submitExamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	resp := submitExamResponse{
		SubmissionID: time.Now().Unix(),
		Message:      "Submission received successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func handleGetMySubmission(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	if user == nil || user.Role != models.RoleStudent {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	submission := models.Submission{
		ID:         1,
		ExamID:     1,
		StudentID:  user.UserID,
		SubmittedAt: time.Now().Add(-1 * time.Hour),
		FileS3Key:  "",
		Score:      nil,
		Status:     models.SubmissionSubmitted,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(submission)
}

