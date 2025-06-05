package service

import (
	"github.com/amirex128/new_site_builder/internal/domain/enums"
)

// IAuthService provides methods to access authentication and authorization context
type IAuthService interface {
	GetRoles() ([]string, error)
	GetSiteIDs() ([]int64, error)
	GetUserID() (*int64, error)
	GetCustomerID() (*int64, error)
	GetUserType() (*enums.UserTypeEnum, error)
	GetEmail() (string, error)
	IsAdmin() (bool, error)
	GetUserOrCustomerID() (*int64, *int64, *enums.UserTypeEnum, error)
}
