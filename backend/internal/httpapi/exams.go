package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/example/scalable-learning-platform/backend/internal/middleware"
	"github.com/example/scalable-learning-platform/backend/internal/models"
)

func registerExamRoutes(r chi.Router, services *Services) {
	r.Post("/{id}/submissions", handleSubmitExam(services))
	r.Get("/{id}/submissions/me", handleGetMySubmission(services))
}

type submitExamRequest struct {
	Answers string `json:"answers"` // simplified; could be structured
	FileS3Key string `json:"file_s3_key,omitempty"`
}

type submitExamResponse struct {
	SubmissionID int64  `json:"submission_id"`
	Message      string `json:"message"`
}

func handleSubmitExam(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := middleware.CurrentUser(r)
		if user == nil || user.Role != models.RoleStudent {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		idStr := chi.URLParam(r, "id")
		examID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid exam id", http.StatusBadRequest)
			return
		}

		var req submitExamRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		submission := &models.Submission{
			ExamID:    examID,
			StudentID: user.UserID,
			FileS3Key: req.FileS3Key,
			Status:    models.SubmissionSubmitted,
		}

		result, err := services.Exams.CreateSubmission(submission)
		if err != nil {
			http.Error(w, "failed to create submission", http.StatusInternalServerError)
			return
		}

		resp := submitExamResponse{
			SubmissionID: result.ID,
			Message:      "Submission received successfully",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func handleGetMySubmission(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := middleware.CurrentUser(r)
		if user == nil || user.Role != models.RoleStudent {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		idStr := chi.URLParam(r, "id")
		examID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid exam id", http.StatusBadRequest)
			return
		}

		submissions, err := services.Exams.GetSubmissionsByStudentID(user.UserID)
		if err != nil {
			http.Error(w, "failed to fetch submissions", http.StatusInternalServerError)
			return
		}

		// Find the submission for this exam
		for _, sub := range submissions {
			if sub.ExamID == examID {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(sub)
				return
			}
		}

		http.Error(w, "submission not found", http.StatusNotFound)
	}
}

