package auth

import (
	"errors"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/gin-gonic/gin"
)

// AuthContextService implements the IAuthContextService interface
type AuthContextService struct {
	logger        sflogger.Logger
	contextGetter func() (*gin.Context, bool)
}

// NewAuthContextService creates a new instance of AuthContextService
func NewAuthContextService(logger sflogger.Logger, contextGetter func() (*gin.Context, bool)) *AuthContextService {
	return &AuthContextService{
		logger:        logger,
		contextGetter: contextGetter,
	}
}

// GetCustomerID returns the current authenticated customer ID
func (s *AuthContextService) GetCustomerID() (int64, error) {
	ctx, exists := s.contextGetter()
	if !exists {
		return 0, errors.New("context not found")
	}

	// Get the customer ID from the context
	// The actual implementation would depend on your authentication middleware
	// This is just a placeholder
	customerID, exists := ctx.Get("customer_id")
	if !exists {
		return 0, errors.New("customer ID not found in context")
	}

	// Convert to int64
	id, ok := customerID.(int64)
	if !ok {
		return 0, errors.New("customer ID is not of type int64")
	}

	return id, nil
}

// GetUserID returns the current authenticated user ID
func (s *AuthContextService) GetUserID() (int64, error) {
	ctx, exists := s.contextGetter()
	if !exists {
		return 0, errors.New("context not found")
	}

	// Get the user ID from the context
	// The actual implementation would depend on your authentication middleware
	// This is just a placeholder
	userID, exists := ctx.Get("user_id")
	if !exists {
		return 0, errors.New("user ID not found in context")
	}

	// Convert to int64
	id, ok := userID.(int64)
	if !ok {
		return 0, errors.New("user ID is not of type int64")
	}

	return id, nil
}

// IsAdmin checks if the current user has admin privileges
func (s *AuthContextService) IsAdmin() (bool, error) {
	ctx, exists := s.contextGetter()
	if !exists {
		return false, errors.New("context not found")
	}

	// Check if the user is an admin
	// The actual implementation would depend on your authentication middleware
	// This is just a placeholder
	isAdmin, exists := ctx.Get("is_admin")
	if !exists {
		return false, errors.New("admin status not found in context")
	}

	// Convert to bool
	admin, ok := isAdmin.(bool)
	if !ok {
		return false, errors.New("admin status is not of type bool")
	}

	return admin, nil
}
