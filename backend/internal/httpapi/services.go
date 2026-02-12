package httpapi

import (
	"fmt"

	"github.com/example/scalable-learning-platform/backend/internal/config"
	"github.com/example/scalable-learning-platform/backend/internal/db"
)

// Services holds all database services
type Services struct {
	Users         *db.UserService
	Sessions      *db.SessionService
	Exams         *db.ExamService
	Notifications *db.NotificationService
	DB            *db.DB
}

// InitServices initializes all services with database connection
func InitServices(cfg *config.Config) (*Services, error) {
	database, err := db.New(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	services := &Services{
		Users:         db.NewUserService(database),
		Sessions:      db.NewSessionService(database),
		Exams:         db.NewExamService(database),
		Notifications: db.NewNotificationService(database),
		DB:            database,
	}

	return services, nil
}

// Close closes all connections
func (s *Services) Close() error {
	if s.DB != nil {
		return s.DB.Close()
	}
	return nil
}
