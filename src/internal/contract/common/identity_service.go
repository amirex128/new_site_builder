package common

import (
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

// IIdentityService provides methods for authentication and JWT token generation
type IIdentityService interface {
	// AddClaim adds a claim to the token
	AddClaim(name string, value string) IIdentityService

	// AddRoles adds role claims to the token
	AddRoles(roles []string) IIdentityService

	// TokenForUser creates a token for a user
	TokenForUser(user domain.User) IIdentityService

	// TokenForCustomer creates a token for a customer
	TokenForCustomer(customer domain.Customer) IIdentityService

	// Make generates and returns the JWT token string
	Make() string

	// VerifyPassword checks if the provided password matches the hashed password
	VerifyPassword(password string, hashedPassword string, salt string) bool

	// HashPassword creates a hashed password with a salt
	HashPassword(password string) (string, string)
}
