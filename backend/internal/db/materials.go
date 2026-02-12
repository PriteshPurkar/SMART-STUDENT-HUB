package db

import (
	"fmt"

	"github.com/example/scalable-learning-platform/backend/internal/models"
)

// MaterialService provides database operations for study materials.
type MaterialService struct {
	db *DB
}

// NewMaterialService creates a new material service.
func NewMaterialService(db *DB) *MaterialService {
	return &MaterialService{db: db}
}

// GetMaterialsBySessionID retrieves all study materials attached to a session.
func (ms *MaterialService) GetMaterialsBySessionID(sessionID int64) ([]models.StudyMaterial, error) {
	rows, err := ms.db.Query(
		`SELECT id, session_id, title, type, s3_key, url, uploaded_by, created_at
		 FROM study_materials
		 WHERE session_id = $1
		 ORDER BY created_at DESC`,
		sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query materials: %w", err)
	}
	defer rows.Close()

	var materials []models.StudyMaterial
	for rows.Next() {
		var m models.StudyMaterial
		if err := rows.Scan(
			&m.ID,
			&m.SessionID,
			&m.Title,
			&m.Type,
			&m.S3Key,
			&m.URL,
			&m.UploadedBy,
			&m.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan material: %w", err)
		}
		materials = append(materials, m)
	}

	return materials, rows.Err()
}

