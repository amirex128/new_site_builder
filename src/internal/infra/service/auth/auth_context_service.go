package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
)

// AuthContextService implements the IAuthContextService interface
type AuthContextService struct {
	ctx      context.Context
	request  *http.Request
	userRepo repository.IUserRepository
	roleRepo repository.IRoleRepository
}

// NewAuthContextService creates a new instance of AuthContextService
func NewAuthContextService(ctx context.Context, request *http.Request, userRepo repository.IUserRepository, roleRepo repository.IRoleRepository) *AuthContextService {
	return &AuthContextService{
		ctx:      ctx,
		request:  request,
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// GetRoles returns the roles of the current user
func (s *AuthContextService) GetRoles() ([]string, error) {
	roles := s.getHeaderValue("X-Roles")
	if roles == "" {
		if s.hasFreeRoute() {
			return []string{}, nil
		}
		return nil, errors.New("authorization access exception: roles not found")
	}

	return strings.Split(roles, ","), nil
}

// GetSiteIDs returns the site IDs the user has access to
func (s *AuthContextService) GetSiteIDs() ([]int64, error) {
	siteIDsStr := s.getHeaderValue("X-Site-Ids")
	if siteIDsStr == "" {
		if s.hasFreeRoute() {
			return []int64{}, nil
		}
		return nil, errors.New("authorization access exception: site IDs not found")
	}

	siteIDsStrArray := strings.Split(siteIDsStr, ",")
	siteIDs := make([]int64, 0, len(siteIDsStrArray))

	for _, idStr := range siteIDsStrArray {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid site ID format: %s", err)
		}
		siteIDs = append(siteIDs, id)
	}

	return siteIDs, nil
}

// GetUserID returns the current authenticated user ID
func (s *AuthContextService) GetUserID() (int64, error) {
	userIDStr := s.getHeaderValue("X-User-Id")
	if userIDStr == "" {
		if s.hasFreeRoute() {
			return 0, nil
		}
		return 0, errors.New("authorization access exception: user ID not found")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID format: %s", err)
	}

	return userID, nil
}

// GetCustomerID returns the current authenticated customer ID
func (s *AuthContextService) GetCustomerID() (int64, error) {
	customerIDStr := s.getHeaderValue("X-Customer-Id")
	if customerIDStr == "" {
		if s.hasFreeRoute() {
			return 0, nil
		}
		return 0, errors.New("authorization access exception: customer ID not found")
	}

	customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid customer ID format: %s", err)
	}

	return customerID, nil
}

// GetUserType returns the type of the current user
func (s *AuthContextService) GetUserType() (string, error) {
	userType := s.getHeaderValue("X-Type")
	if userType == "" {
		if s.hasFreeRoute() {
			return "guest", nil
		}
		return "", errors.New("authorization access exception: user type not found")
	}

	if userType == "user" {
		return "user", nil
	}

	if userType == "customer" {
		return "customer", nil
	}

	return "", errors.New("authorization access exception: invalid user type")
}

// GetEmail returns the email of the current user
func (s *AuthContextService) GetEmail() (string, error) {
	email := s.getHeaderValue("X-Email")
	if email == "" {
		if s.hasFreeRoute() {
			return "", nil
		}
		return "", errors.New("authorization access exception: email not found")
	}

	return email, nil
}

// IsAdmin checks if the current user has admin privileges
func (s *AuthContextService) IsAdmin() (bool, error) {
	// Check if admin flag is set in context or header
	isAdminStr := s.getHeaderValue("X-Is-Admin")
	if isAdminStr == "true" {
		return true, nil
	}

	// Get user ID
	userID, err := s.GetUserID()
	if err != nil {
		return false, err
	}

	// If user ID is 0, the user is not authenticated or is in a free route
	if userID == 0 {
		return false, nil
	}

	// Get user from repository
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}

	// Check if user has admin flag
	return user.IsAdmin, nil
}

// getHeaderValue gets a value from the header, with JWT token fallback for local development
func (s *AuthContextService) getHeaderValue(headerName string) string {
	// If we have a direct context value, use it
	if s.ctx != nil {
		if value, ok := s.ctx.Value(strings.ToLower(strings.TrimPrefix(headerName, "X-"))).(string); ok && value != "" {
			return value
		}
	}

	// If we don't have a request, we can't get header values
	if s.request == nil {
		return ""
	}

	// Try to get from header directly
	if headerValue := s.request.Header.Get(headerName); headerValue != "" {
		return headerValue
	}

	return ""
}

// hasFreeRoute checks if the current request path is a free route
func (s *AuthContextService) hasFreeRoute() bool {
	if s.request == nil {
		return false
	}

	currentPath := s.request.URL.Path
	pattern := `^/[^/]+/free/`
	matched, _ := regexp.MatchString(pattern, currentPath)
	return matched
}
