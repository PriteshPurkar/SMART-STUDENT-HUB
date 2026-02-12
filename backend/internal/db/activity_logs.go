package db

import (
	"fmt"

	"github.com/example/scalable-learning-platform/backend/internal/models"
)

// ActivityLogService provides database operations for activity logs.
type ActivityLogService struct {
	db *DB
}

// NewActivityLogService creates a new activity log service.
func NewActivityLogService(db *DB) *ActivityLogService {
	return &ActivityLogService{db: db}
}

// CreateLog inserts a new activity log entry.
func (als *ActivityLogService) CreateLog(log *models.ActivityLog) (*models.ActivityLog, error) {
	err := als.db.QueryRow(
		`INSERT INTO activity_logs (user_id, action, resource_type, resource_id, metadata)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at`,
		log.UserID,
		log.Action,
		log.ResourceType,
		log.ResourceID,
		log.Metadata,
	).Scan(&log.ID, &log.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create activity log: %w", err)
	}
	return log, nil
}

