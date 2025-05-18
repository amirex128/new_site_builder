package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

// IdentityService implements the IIdentityService interface
type IdentityService struct {
	jwtSecret    string
	jwtExpiresIn time.Duration
	claims       jwt.MapClaims
	userID       int64
	customerID   int64
}

// NewIdentityService creates a new instance of IdentityService
func NewIdentityService(jwtSecret string, jwtExpiresIn time.Duration) *IdentityService {
	return &IdentityService{
		jwtSecret:    jwtSecret,
		jwtExpiresIn: jwtExpiresIn,
		claims:       make(jwt.MapClaims),
	}
}

// AddClaim adds a claim to the token
func (s *IdentityService) AddClaim(name string, value string) common.IIdentityService {
	s.claims[name] = value
	return s
}

// AddRoles adds role claims to the token
func (s *IdentityService) AddRoles(roles []string) common.IIdentityService {
	s.claims["roles"] = roles
	return s
}

// TokenForUser creates a token for a user
func (s *IdentityService) TokenForUser(user domain.User) common.IIdentityService {
	s.userID = user.ID
	s.claims["email"] = user.Email
	s.claims["is_admin"] = user.IsAdmin
	s.claims["exp"] = time.Now().Add(s.jwtExpiresIn).Unix()
	return s
}

// TokenForCustomer creates a token for a customer
func (s *IdentityService) TokenForCustomer(customer domain.Customer) common.IIdentityService {
	s.customerID = customer.ID
	s.claims["customer_id"] = customer.ID
	s.claims["site_id"] = customer.SiteID
	s.claims["email"] = customer.Email
	s.claims["exp"] = time.Now().Add(s.jwtExpiresIn).Unix()
	return s
}

// Make generates and returns the JWT token string
func (s *IdentityService) Make() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, s.claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return ""
	}
	return tokenString
}

// VerifyPassword checks if the provided password matches the hashed password
func (s *IdentityService) VerifyPassword(password string, hashedPassword string, salt string) bool {
	// Hash the provided password with the salt
	hash := sha256.New()
	hash.Write([]byte(password + salt))
	calculatedHash := hex.EncodeToString(hash.Sum(nil))

	// Compare the calculated hash with the stored hash
	return calculatedHash == hashedPassword
}

// HashPassword creates a hashed password with a salt
func (s *IdentityService) HashPassword(password string) (string, string) {
	// Generate a random salt
	saltBytes := make([]byte, 16)
	_, err := rand.Read(saltBytes)
	if err != nil {
		// If random generation fails, use a timestamp-based salt
		saltBytes = []byte(time.Now().String())
	}
	salt := base64.StdEncoding.EncodeToString(saltBytes)

	// Hash the password with the salt
	hash := sha256.New()
	hash.Write([]byte(password + salt))
	hashedPassword := hex.EncodeToString(hash.Sum(nil))

	return hashedPassword, salt
}
