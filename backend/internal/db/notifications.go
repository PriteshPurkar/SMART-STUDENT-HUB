package db

import (
	"fmt"

	"github.com/example/scalable-learning-platform/backend/internal/models"
)

// NotificationService provides database operations for notifications
type NotificationService struct {
	db *DB
}

// NewNotificationService creates a new notification service
func NewNotificationService(db *DB) *NotificationService {
	return &NotificationService{db: db}
}

// CreateNotification creates a new notification
func (ns *NotificationService) CreateNotification(notification *models.Notification) (*models.Notification, error) {
	var notifID int64
	err := ns.db.QueryRow(
		`INSERT INTO notifications (user_id, type, message, is_read)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at`,
		notification.UserID, notification.Type, notification.Message, notification.IsRead,
	).Scan(&notifID, &notification.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	notification.ID = notifID
	return notification, nil
}

// GetNotificationsByUserID retrieves notifications for a user
func (ns *NotificationService) GetNotificationsByUserID(userID int64) ([]models.Notification, error) {
	rows, err := ns.db.Query(
		`SELECT id, user_id, type, message, is_read, created_at
		 FROM notifications WHERE user_id = $1
		 ORDER BY created_at DESC LIMIT 50`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notif models.Notification
		if err := rows.Scan(&notif.ID, &notif.UserID, &notif.Type, &notif.Message, &notif.IsRead, &notif.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notif)
	}

	return notifications, rows.Err()
}

// MarkAsRead marks a notification as read
func (ns *NotificationService) MarkAsRead(notifID int64) error {
	result, err := ns.db.Exec(
		`UPDATE notifications SET is_read = TRUE WHERE id = $1`,
		notifID,
	)
	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (ns *NotificationService) MarkAllAsRead(userID int64) error {
	_, err := ns.db.Exec(
		`UPDATE notifications SET is_read = TRUE WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update notifications: %w", err)
	}

	return nil
}

// GetUnreadCount gets unread notification count for a user
func (ns *NotificationService) GetUnreadCount(userID int64) (int, error) {
	var count int
	err := ns.db.QueryRow(
		`SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE`,
		userID,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}
