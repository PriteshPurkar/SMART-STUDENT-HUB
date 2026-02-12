package db

import (
	"fmt"
	"time"

	"github.com/example/scalable-learning-platform/backend/internal/models"
)

// SessionService provides database operations for sessions
type SessionService struct {
	db *DB
}

// NewSessionService creates a new session service
func NewSessionService(db *DB) *SessionService {
	return &SessionService{db: db}
}

// CreateSession creates a new session in the database
func (ss *SessionService) CreateSession(session *models.Session) (*models.Session, error) {
	var sessionID int64
	err := ss.db.QueryRow(
		`INSERT INTO sessions (title, description, start_time, end_time, status, video_url, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id`,
		session.Title, session.Description, session.StartTime, session.EndTime,
		string(session.Status), session.VideoURL, session.CreatedBy,
	).Scan(&sessionID)

	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	session.ID = sessionID
	session.CreatedAt = time.Now()
	return session, nil
}

// GetSessionByID retrieves a session by ID
func (ss *SessionService) GetSessionByID(id int64) (*models.Session, error) {
	session := &models.Session{}
	err := ss.db.QueryRow(
		`SELECT id, title, description, start_time, end_time, status, video_url, created_by, created_at
		 FROM sessions WHERE id = $1`,
		id,
	).Scan(&session.ID, &session.Title, &session.Description, &session.StartTime,
		&session.EndTime, &session.Status, &session.VideoURL, &session.CreatedBy, &session.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

// GetAllSessions retrieves all sessions
func (ss *SessionService) GetAllSessions() ([]models.Session, error) {
	rows, err := ss.db.Query(
		`SELECT id, title, description, start_time, end_time, status, video_url, created_by, created_at
		 FROM sessions ORDER BY start_time DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var session models.Session
		if err := rows.Scan(&session.ID, &session.Title, &session.Description, &session.StartTime,
			&session.EndTime, &session.Status, &session.VideoURL, &session.CreatedBy, &session.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

// GetUpcomingSessions retrieves sessions starting after now
func (ss *SessionService) GetUpcomingSessions(limit int) ([]models.Session, error) {
	rows, err := ss.db.Query(
		`SELECT id, title, description, start_time, end_time, status, video_url, created_by, created_at
		 FROM sessions WHERE start_time > NOW()
		 ORDER BY start_time ASC LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query upcoming sessions: %w", err)
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var session models.Session
		if err := rows.Scan(&session.ID, &session.Title, &session.Description, &session.StartTime,
			&session.EndTime, &session.Status, &session.VideoURL, &session.CreatedBy, &session.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

// GetPastSessions retrieves completed sessions
func (ss *SessionService) GetPastSessions(limit int) ([]models.Session, error) {
	rows, err := ss.db.Query(
		`SELECT id, title, description, start_time, end_time, status, video_url, created_by, created_at
		 FROM sessions WHERE end_time < NOW()
		 ORDER BY start_time DESC LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query past sessions: %w", err)
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var session models.Session
		if err := rows.Scan(&session.ID, &session.Title, &session.Description, &session.StartTime,
			&session.EndTime, &session.Status, &session.VideoURL, &session.CreatedBy, &session.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

// UpdateSessionStatus updates a session's status
func (ss *SessionService) UpdateSessionStatus(sessionID int64, status models.SessionStatus) error {
	result, err := ss.db.Exec(
		`UPDATE sessions SET status = $1, updated_at = NOW() WHERE id = $2`,
		string(status), sessionID,
	)
	if err != nil {
		return fmt.Errorf("failed to update session status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// GetSessionStats returns statistics for a session
func (ss *SessionService) GetSessionStats(sessionID int64) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get submission count
	var submissionCount int
	err := ss.db.QueryRow(
		`SELECT COUNT(*) FROM submissions WHERE exam_id IN (
			SELECT id FROM exams WHERE session_id = $1
		)`,
		sessionID,
	).Scan(&submissionCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get submission count: %w", err)
	}

	stats["submission_count"] = submissionCount

	// Get active exams count
	var activeExams int
	err = ss.db.QueryRow(
		`SELECT COUNT(*) FROM exams WHERE session_id = $1 AND open_time <= NOW() AND close_time > NOW()`,
		sessionID,
	).Scan(&activeExams)
	if err != nil {
		return nil, fmt.Errorf("failed to get active exams: %w", err)
	}

	stats["active_exams"] = activeExams

	return stats, nil
}
