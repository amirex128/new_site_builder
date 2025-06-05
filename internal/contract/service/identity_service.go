package service

import (
	domain2 "github.com/amirex128/new_site_builder/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// IIdentityService provides methods for authentication and JWT token generation
type IIdentityService interface {
	// AddClaim adds a claim to the token
	AddClaim(name string, value string) IIdentityService

	// AddRoles adds role claims to the token
	AddRoles(roles []string) IIdentityService

	// TokenForUser creates a token for a user
	TokenForUser(user domain2.User) IIdentityService

	// TokenForCustomer creates a token for a customer
	TokenForCustomer(customer domain2.Customer) IIdentityService

	// Make generates and returns the JWT token string
	Make() string

	// VerifyPassword checks if the provided password matches the hashed password
	VerifyPassword(password string, hashedPassword string, salt string) bool

	// HashPassword creates a hashed password with a salt
	HashPassword(password string) (string, string)

	// GetToken GetTokenFromContext extracts the JWT token from the request context
	GetToken(c *gin.Context) (*jwt.Token, error)
	GetTokenString(c *gin.Context) (string, error)

	// VerifyToken validates a token string and returns the parsed token
	VerifyToken(tokenString string) (*jwt.Token, error)
	VerifyTokenContext(ctx *gin.Context) (*jwt.Token, error)

	// GetClaim extracts a specific claim from a JWT token
	GetClaim(ctx *gin.Context, claimName string) (string, error)
}
