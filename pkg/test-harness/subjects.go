// Package testharness provides simple test subjects for the Go project indexing system.
// This package contains basic examples of each Phase 1.1 requirement.
package testharness

import (
	"context"
	"fmt"
	"time"
)

const (
	DefaultTimeout = 30 * time.Second
)
const MaxRetries = 3

// Global variables for testing variable extraction
var (
	GlobalCounter int64
	//internal variable for testing
	isDebugMode bool = false
)

type StringA string
type StringB = string

// User represents a simple user entity.
// Testing struct field extraction with types and tags.
type User struct {
	ID    int64  `json:"id" db:"user_id"`
	Name  string `json:"name" db:"full_name"`
	Email string `json:"email" db:"email"`
}

// UserService defines operations for user management.
// Testing interface method signature extraction.
type UserService interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, user *User) error
}

// Service implements user business logic.
type Service struct {
	userService UserService
}

// NewService creates a new Service instance.
// Testing standalone function extraction.
func NewService(userService UserService) *Service {
	return &Service{
		userService: userService,
	}
}

// ValidateEmail validates an email address format.
// Testing standalone function with parameters and return values.
func ValidateEmail(email string) bool {
	return len(email) > 0 && contains(email, "@")
}

// CreateUser creates a new user.
// Testing method extraction with receiver.
func (s *Service) CreateUser(ctx context.Context, name, email string) (*User, error) {
	if !ValidateEmail(email) {
		return nil, fmt.Errorf("invalid email: %s", email)
	}

	user := &User{
		Name:  name,
		Email: email,
	}

	return user, s.userService.Create(ctx, user)
}

// GetUser retrieves a user by ID.
// Testing method with different parameter types.
func (s *Service) GetUser(ctx context.Context, id int64) (*User, error) {
	return s.userService.GetByID(ctx, id)
}

// contains is a helper function for string operations.
// Testing private/unexported function extraction.
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
