package httpapi

import (
	"errors"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/example/scalable-learning-platform/backend/internal/models"
)

type inMemoryAuthService struct {
	mu    sync.Mutex
	users map[string]*models.User // key by email
	lastID int64
}

func newInMemoryAuthService() *inMemoryAuthService {
	return &inMemoryAuthService{
		users: make(map[string]*models.User),
	}
}

func (s *inMemoryAuthService) RegisterStudent(name, email, password string) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[email]; exists {
		return nil, errors.New("email already registered")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	s.lastID++
	u := &models.User{
		ID:        s.lastID,
		Name:      name,
		Email:     email,
		Password:  string(hash),
		Role:      models.RoleStudent,
		CreatedAt: time.Now(),
	}
	s.users[email] = u
	return u, nil
}

func (s *inMemoryAuthService) Login(email, password string) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[email]
	if !ok {
		return nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return u, nil
}

func (s *inMemoryAuthService) GetUserByID(id int64) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range s.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("not found")
}

