package auth

import (
	"context"
	"errors"
	"strconv"

	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
)

// AuthContextService implements the IAuthContextService interface
type AuthContextService struct {
	ctx      context.Context
	userRepo repository.IUserRepository
	roleRepo repository.IRoleRepository
}

// NewAuthContextService creates a new instance of AuthContextService
func NewAuthContextService(ctx context.Context, userRepo repository.IUserRepository, roleRepo repository.IRoleRepository) *AuthContextService {
	return &AuthContextService{
		ctx:      ctx,
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// GetCustomerID returns the current authenticated customer ID
func (s *AuthContextService) GetCustomerID() (int64, error) {
	// Get customer ID from context if available
	if s.ctx == nil {
		return 0, errors.New("context is nil")
	}

	customerID, ok := s.ctx.Value("customer_id").(int64)
	if !ok {
		// Try to get it as a string
		customerIDStr, ok := s.ctx.Value("customer_id").(string)
		if !ok {
			return 0, errors.New("customer ID not found in context")
		}

		// Convert string to int64
		var err error
		customerID, err = strconv.ParseInt(customerIDStr, 10, 64)
		if err != nil {
			return 0, errors.New("invalid customer ID format")
		}
	}

	return customerID, nil
}

// GetUserID returns the current authenticated user ID
func (s *AuthContextService) GetUserID() (int64, error) {
	// Get user ID from context if available
	if s.ctx == nil {
		return 0, errors.New("context is nil")
	}

	userID, ok := s.ctx.Value("user_id").(int64)
	if !ok {
		// Try to get it as a string
		userIDStr, ok := s.ctx.Value("user_id").(string)
		if !ok {
			return 0, errors.New("user ID not found in context")
		}

		// Convert string to int64
		var err error
		userID, err = strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return 0, errors.New("invalid user ID format")
		}
	}

	return userID, nil
}

// IsAdmin checks if the current user has admin privileges
func (s *AuthContextService) IsAdmin() (bool, error) {
	// Check if admin flag is set in context
	isAdmin, ok := s.ctx.Value("is_admin").(bool)
	if ok && isAdmin {
		return true, nil
	}

	// Get user ID
	userID, err := s.GetUserID()
	if err != nil {
		return false, err
	}

	// Get user
	_, err = s.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}

	// Check if user has admin role
	// Since GetAllByUserID is missing, we'll use a direct check for admin status from context
	// This should be replaced with proper role checking once the method is implemented
	return isAdmin, nil
}
