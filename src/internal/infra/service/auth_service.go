package service

import (
	"errors"
	"fmt"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
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
		return nil, errors.New("خطا در بررسی دسترسی کاربر")
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
		return nil, errors.New("خطا در بررسی دسترسی کاربر")
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
func (s *AuthContextService) GetUserID() (*int64, error) {
	userIDStr, err := s.identity.GetClaim(s.ctx, "user_id")
	if err != nil {
		return nil, errors.New("خطا در بررسی دسترسی کاربر")
	}

	id, err := parseID(userIDStr, "user ID")
	return &id, err
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
func (s *AuthContextService) GetCustomerID() (*int64, error) {
	customerIDStr, err := s.identity.GetClaim(s.ctx, "customer_id")
	if err != nil {
		return nil, errors.New("خطا در بررسی دسترسی کاربر")
	}

	id, err := parseID(customerIDStr, "customer ID")
	return &id, err
}

// GetUserType returns the type of the current user
func (s *AuthContextService) GetUserType() (*enums.UserTypeEnum, error) {
	userType, err := s.identity.GetClaim(s.ctx, "type")
	if err != nil {
		return nil, errors.New("خطا در بررسی دسترسی کاربر")
	}

	s2, err := validateUserType(userType)
	if err != nil {
		return nil, err
	}
	if s2 == "user" {
		return (*enums.UserTypeEnum)(&s2), nil
	}
	if s2 == "customer" {
		return (*enums.UserTypeEnum)(&s2), nil
	}
	return nil, errors.New("خطا در بررسی دسترسی کاربر")
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
		return "", errors.New("خطا در بررسی دسترسی کاربر")
	}

	return email, nil
}

// IsAdmin checks if the current user has admin privileges
func (s *AuthContextService) IsAdmin() (bool, error) {
	isAdminStr, err := s.identity.GetClaim(s.ctx, "is_admin")
	if err != nil {
		return false, errors.New("خطا در بررسی دسترسی کاربر")
	}
	if isAdminStr == "true" {
		return true, nil
	}
	if isAdminStr == "false" {
		return false, errors.New("خطا در بررسی دسترسی کاربر")
	}
	return false, errors.New("خطا در بررسی دسترسی کاربر")
}

func (s *AuthContextService) GetUserOrCustomerID() (*int64, *int64, *enums.UserTypeEnum, error) {
	var customerID, userID *int64
	var err error

	userType, err := s.GetUserType()
	if err != nil {
		return nil, nil, nil, errors.New("خطا در احراز هویت کاربر")
	}

	if *userType == enums.CustomerTypeValue {
		customerID, err = s.GetCustomerID()
		if err != nil {
			return nil, nil, nil, errors.New("خطا در احراز هویت کاربر")
		}
		return userID, customerID, userType, nil
	}
	if *userType == enums.UserTypeValue {
		userID, err = s.GetUserID()
		if err != nil {
			return nil, nil, nil, errors.New("خطا در احراز هویت کاربر")
		}
		return userID, customerID, userType, nil
	}

	return nil, nil, nil, errors.New("خطا در احراز هویت کاربر")
}
