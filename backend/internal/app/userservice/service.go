package userservice

import (
	"context"
	"fmt"
	"strings"

	"github.com/axelfrache/savesync/internal/domain"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Service handles user business logic
type Service struct {
	userRepo domain.UserRepository
	logger   *zap.Logger
}

// New creates a new user service
func New(userRepo domain.UserRepository, logger *zap.Logger) *Service {
	return &Service{
		userRepo: userRepo,
		logger:   logger,
	}
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, email, password string) (*domain.User, error) {
	// Validate email
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if !strings.Contains(email, "@") {
		return nil, fmt.Errorf("invalid email format")
	}

	// Validate password
	if len(password) < 8 {
		return nil, fmt.Errorf("password must be at least 8 characters")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && err != domain.ErrNotFound {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("user registered", zap.String("email", email), zap.Int64("id", user.ID))

	// Clear password hash before returning
	user.PasswordHash = ""
	return user, nil
}

// Authenticate verifies email and password
func (s *Service) Authenticate(ctx context.Context, email, password string) (*domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err == domain.ErrNotFound {
		return nil, fmt.Errorf("invalid email or password")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.logger.Warn("failed login attempt", zap.String("email", email))
		return nil, fmt.Errorf("invalid email or password")
	}

	s.logger.Info("user authenticated", zap.String("email", email), zap.Int64("id", user.ID))

	// Clear password hash before returning
	user.PasswordHash = ""
	return user, nil
}

// GetByID retrieves a user by ID
func (s *Service) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Clear password hash
	user.PasswordHash = ""
	return user, nil
}

// GetAll retrieves all users
func (s *Service) GetAll(ctx context.Context) ([]*domain.User, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Clear password hashes
	for _, user := range users {
		user.PasswordHash = ""
	}

	return users, nil
}

// Update updates a user
func (s *Service) Update(ctx context.Context, user *domain.User) error {
	return s.userRepo.Update(ctx, user)
}

// SetAdminStatus updates the admin status of a user
func (s *Service) SetAdminStatus(ctx context.Context, id int64, isAdmin bool) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	user.IsAdmin = isAdmin
	return s.userRepo.Update(ctx, user)
}

// Delete removes a user
func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.userRepo.Delete(ctx, id)
}
