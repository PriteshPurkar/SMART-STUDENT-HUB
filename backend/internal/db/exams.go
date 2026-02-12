package db

import (
	"fmt"

	"github.com/example/scalable-learning-platform/backend/internal/models"
)

// ExamService provides database operations for exams
type ExamService struct {
	db *DB
}

// NewExamService creates a new exam service
func NewExamService(db *DB) *ExamService {
	return &ExamService{db: db}
}

// CreateExam creates a new exam in the database
func (es *ExamService) CreateExam(exam *models.Exam) (*models.Exam, error) {
	var examID int64
	err := es.db.QueryRow(
		`INSERT INTO exams (session_id, title, open_time, close_time, max_score, type)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id`,
		exam.SessionID, exam.Title, exam.OpenTime, exam.CloseTime, exam.MaxScore, string(exam.Type),
	).Scan(&examID)

	if err != nil {
		return nil, fmt.Errorf("failed to create exam: %w", err)
	}

	exam.ID = examID
	return exam, nil
}

// GetExamByID retrieves an exam by ID
func (es *ExamService) GetExamByID(id int64) (*models.Exam, error) {
	exam := &models.Exam{}
	err := es.db.QueryRow(
		`SELECT id, session_id, title, open_time, close_time, max_score, type
		 FROM exams WHERE id = $1`,
		id,
	).Scan(&exam.ID, &exam.SessionID, &exam.Title, &exam.OpenTime, &exam.CloseTime, &exam.MaxScore, &exam.Type)

	if err != nil {
		return nil, fmt.Errorf("failed to get exam: %w", err)
	}

	return exam, nil
}

// GetExamsBySessionID retrieves all exams for a session
func (es *ExamService) GetExamsBySessionID(sessionID int64) ([]models.Exam, error) {
	rows, err := es.db.Query(
		`SELECT id, session_id, title, open_time, close_time, max_score, type
		 FROM exams WHERE session_id = $1 ORDER BY open_time DESC`,
		sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query exams: %w", err)
	}
	defer rows.Close()

	var exams []models.Exam
	for rows.Next() {
		var exam models.Exam
		if err := rows.Scan(&exam.ID, &exam.SessionID, &exam.Title, &exam.OpenTime, &exam.CloseTime, &exam.MaxScore, &exam.Type); err != nil {
			return nil, fmt.Errorf("failed to scan exam: %w", err)
		}
		exams = append(exams, exam)
	}

	return exams, rows.Err()
}

// GetAllExams retrieves all exams
func (es *ExamService) GetAllExams() ([]models.Exam, error) {
	rows, err := es.db.Query(
		`SELECT id, session_id, title, open_time, close_time, max_score, type
		 FROM exams ORDER BY open_time DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query exams: %w", err)
	}
	defer rows.Close()

	var exams []models.Exam
	for rows.Next() {
		var exam models.Exam
		if err := rows.Scan(&exam.ID, &exam.SessionID, &exam.Title, &exam.OpenTime, &exam.CloseTime, &exam.MaxScore, &exam.Type); err != nil {
			return nil, fmt.Errorf("failed to scan exam: %w", err)
		}
		exams = append(exams, exam)
	}

	return exams, rows.Err()
}

// CreateSubmission creates a new submission
func (es *ExamService) CreateSubmission(submission *models.Submission) (*models.Submission, error) {
	var submissionID int64
	err := es.db.QueryRow(
		`INSERT INTO submissions (exam_id, student_id, file_s3_key, status)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, submitted_at`,
		submission.ExamID, submission.StudentID, submission.FileS3Key, string(submission.Status),
	).Scan(&submissionID, &submission.SubmittedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create submission: %w", err)
	}

	submission.ID = submissionID
	return submission, nil
}

// GetSubmissionsByExamID retrieves all submissions for an exam
func (es *ExamService) GetSubmissionsByExamID(examID int64) ([]models.Submission, error) {
	rows, err := es.db.Query(
		`SELECT id, exam_id, student_id, submitted_at, file_s3_key, score, status
		 FROM submissions WHERE exam_id = $1 ORDER BY submitted_at DESC`,
		examID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query submissions: %w", err)
	}
	defer rows.Close()

	var submissions []models.Submission
	for rows.Next() {
		var submission models.Submission
		if err := rows.Scan(&submission.ID, &submission.ExamID, &submission.StudentID, &submission.SubmittedAt,
			&submission.FileS3Key, &submission.Score, &submission.Status); err != nil {
			return nil, fmt.Errorf("failed to scan submission: %w", err)
		}
		submissions = append(submissions, submission)
	}

	return submissions, rows.Err()
}

// GetSubmissionsByStudentID retrieves all submissions for a student
func (es *ExamService) GetSubmissionsByStudentID(studentID int64) ([]models.Submission, error) {
	rows, err := es.db.Query(
		`SELECT id, exam_id, student_id, submitted_at, file_s3_key, score, status
		 FROM submissions WHERE student_id = $1 ORDER BY submitted_at DESC`,
		studentID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query submissions: %w", err)
	}
	defer rows.Close()

	var submissions []models.Submission
	for rows.Next() {
		var submission models.Submission
		if err := rows.Scan(&submission.ID, &submission.ExamID, &submission.StudentID, &submission.SubmittedAt,
			&submission.FileS3Key, &submission.Score, &submission.Status); err != nil {
			return nil, fmt.Errorf("failed to scan submission: %w", err)
		}
		submissions = append(submissions, submission)
	}

	return submissions, rows.Err()
}

// UpdateSubmissionScore updates a submission's score and status
func (es *ExamService) UpdateSubmissionScore(submissionID int64, score int, status models.SubmissionStatus) error {
	result, err := es.db.Exec(
		`UPDATE submissions SET score = $1, status = $2 WHERE id = $3`,
		score, string(status), submissionID,
	)
	if err != nil {
		return fmt.Errorf("failed to update submission: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("submission not found")
	}

	return nil
}

// GetSubmissionReports gets submission statistics
func (es *ExamService) GetSubmissionReports() ([]map[string]interface{}, error) {
	rows, err := es.db.Query(
		`SELECT e.id as exam_id, e.title, COUNT(s.id) as total_submissions,
		 COUNT(CASE WHEN s.status = 'GRADED' THEN 1 END) as graded_count,
		 COUNT(CASE WHEN s.status = 'SUBMITTED' THEN 1 END) as pending_count
		 FROM exams e
		 LEFT JOIN submissions s ON e.id = s.exam_id
		 GROUP BY e.id, e.title
		 ORDER BY e.id DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query submission reports: %w", err)
	}
	defer rows.Close()

	var reports []map[string]interface{}
	for rows.Next() {
		var examID int64
		var title string
		var totalSubmissions, gradedCount, pendingCount int

		if err := rows.Scan(&examID, &title, &totalSubmissions, &gradedCount, &pendingCount); err != nil {
			return nil, fmt.Errorf("failed to scan report: %w", err)
		}

		report := map[string]interface{}{
			"exam_id":             examID,
			"title":               title,
			"total_submissions":   totalSubmissions,
			"graded_count":        gradedCount,
			"pending_count":       pendingCount,
		}
		reports = append(reports, report)
	}

	return reports, rows.Err()
}
