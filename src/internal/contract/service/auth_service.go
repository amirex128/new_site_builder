package service

// IAuthService provides methods to access authentication and authorization context
type IAuthService interface {
	GetRoles() ([]string, error)
	GetSiteIDs() ([]int64, error)
	GetUserID() (int64, error)
	GetCustomerID() (int64, error)
	GetUserType() (string, error)
	GetEmail() (string, error)
	IsAdmin() (bool, error)
}
