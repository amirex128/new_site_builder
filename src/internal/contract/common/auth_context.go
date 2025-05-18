package common

// IAuthContextService provides methods to access authentication and authorization context
type IAuthContextService interface {
	// GetCustomerID returns the current authenticated customer ID
	GetCustomerID() (int64, error)

	// GetUserID returns the current authenticated user ID
	GetUserID() (int64, error)

	// IsAdmin checks if the current user has admin privileges
	IsAdmin() (bool, error)
}
