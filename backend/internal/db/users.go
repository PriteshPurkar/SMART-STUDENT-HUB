package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/example/scalable-learning-platform/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// UserService provides database operations for users
type UserService struct {
	db *DB
}

// NewUserService creates a new user service
func NewUserService(db *DB) *UserService {
	return &UserService{db: db}
}

// CreateUser creates a new user in the database
func (us *UserService) CreateUser(name, email, password string, role models.Role) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	var userID int64
	err = us.db.QueryRow(
		`INSERT INTO users (name, email, password_hash, role)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		name, email, string(hashedPassword), string(role),
	).Scan(&userID)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &models.User{
		ID:        userID,
		Name:      name,
		Email:     email,
		Password:  "", // Don't return password
		Role:      role,
		CreatedAt: time.Now(),
	}, nil
}

// GetUserByEmail retrieves a user by email
func (us *UserService) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := us.db.QueryRow(
		`SELECT id, name, email, password_hash, role, created_at
		 FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (us *UserService) GetUserByID(id int64) (*models.User, error) {
	user := &models.User{}
	err := us.db.QueryRow(
		`SELECT id, name, email, password_hash, role, created_at
		 FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// VerifyPassword verifies a user's password
func (us *UserService) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GetAllUsers retrieves all users (for admin purposes)
func (us *UserService) GetAllUsers() ([]models.User, error) {
	rows, err := us.db.Query(
		`SELECT id, name, email, password_hash, role, created_at
		 FROM users ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		user.Password = "" // Don't return password
		users = append(users, user)
	}

	return users, rows.Err()
}
