package service

import (
	"errors"
	"fmt"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthContextService implements the IAuthContextService interface
type AuthContextService struct {
	ctx      *gin.Context
	identity service.IIdentityService
}

// NewAuthContextService creates a new instance of AuthContextService
func NewAuthContextService(c *gin.Context, identity service.IIdentityService) *AuthContextService {
	return &AuthContextService{
		ctx:      c,
		identity: identity,
	}
}

// GetRoles returns the roles of the current user
func (s *AuthContextService) GetRoles() ([]string, error) {
	rolesStr, err := s.identity.GetClaim(s.ctx, "roles")
	if err != nil {
		return nil, errors.New("authorization access exception: roles not found")
	}

	if rolesStr == "" {
		return []string{}, nil
	}

	return strings.Split(rolesStr, ","), nil
}

// GetSiteIDs returns the site IDs the user has access to
func (s *AuthContextService) GetSiteIDs() ([]int64, error) {
	siteIDsStr, err := s.identity.GetClaim(s.ctx, "site_id")
	if err != nil {
		return nil, errors.New("authorization access exception: site IDs not found")
	}

	if siteIDsStr == "" {
		return []int64{}, nil
	}

	return parseSiteIDs(siteIDsStr)
}

// Helper function to parse site IDs from a comma-separated string
func parseSiteIDs(siteIDsStr string) ([]int64, error) {
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
	userIDStr, err := s.identity.GetClaim(s.ctx, "user_id")
	if err != nil {
		return 0, errors.New("authorization access exception: user ID not found")
	}

	return parseID(userIDStr, "user ID")
}

// Helper function to parse an ID from a string
func parseID(idStr string, idType string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format: %s", idType, err)
	}
	return id, nil
}

// GetCustomerID returns the current authenticated customer ID
func (s *AuthContextService) GetCustomerID() (int64, error) {
	customerIDStr, err := s.identity.GetClaim(s.ctx, "customer_id")
	if err != nil {
		return 0, errors.New("authorization access exception: customer ID not found")
	}

	return parseID(customerIDStr, "customer ID")
}

// GetUserType returns the type of the current user
func (s *AuthContextService) GetUserType() (string, error) {
	userType, err := s.identity.GetClaim(s.ctx, "type")
	if err != nil {
		return "", errors.New("authorization access exception: user type not found")
	}

	return validateUserType(userType)
}

// Helper function to validate user type
func validateUserType(userType string) (string, error) {
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
	email, err := s.identity.GetClaim(s.ctx, "email")
	if err != nil {
		return "", errors.New("authorization access exception: email not found")
	}

	return email, nil
}

// IsAdmin checks if the current user has admin privileges
func (s *AuthContextService) IsAdmin() (bool, error) {
	isAdminStr, err := s.identity.GetClaim(s.ctx, "is_admin")
	if err != nil {
		return false, nil // Not an admin if claim not found
	}

	return isAdminStr == "true", nil
}
