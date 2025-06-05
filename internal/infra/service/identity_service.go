package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/amirex128/new_site_builder/internal/contract/service"
	domain2 "github.com/amirex128/new_site_builder/internal/domain"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// IdentityService implements the IIdentityService interface
type IdentityService struct {
	jwtSecret    string
	jwtIssuer    string
	jwtAudience  string
	jwtExpiresIn time.Duration
	claims       []jwt.Claims
}

// NewIdentityService creates a new instance of IdentityService
func NewIdentityService(jwtSecret string, issuer string, audience string, jwtExpiresIn time.Duration) *IdentityService {
	return &IdentityService{
		jwtSecret:    jwtSecret,
		jwtIssuer:    issuer,
		jwtAudience:  audience,
		jwtExpiresIn: jwtExpiresIn,
		claims:       make([]jwt.Claims, 0),
	}
}

// AddClaim adds a claim to the token
func (s *IdentityService) AddClaim(name string, value string) service.IIdentityService {
	s.claims = append(s.claims, jwt.RegisteredClaims{
		Subject: name,
		ID:      value,
	})
	return s
}

// AddRoles adds role claims to the token
func (s *IdentityService) AddRoles(roles []string) service.IIdentityService {
	if roles == nil || len(roles) == 0 {
		return s
	}

	// Join roles into a comma-separated string like in .NET implementation
	s.claims = append(s.claims, jwt.RegisteredClaims{
		Subject: "roles",
		ID:      joinRoles(roles),
	})
	return s
}

// joinRoles joins roles into a comma-separated string
func joinRoles(roles []string) string {
	if len(roles) == 0 {
		return ""
	}

	result := roles[0]
	for i := 1; i < len(roles); i++ {
		result += "," + roles[i]
	}
	return result
}

// TokenForUser creates a token for a user
func (s *IdentityService) TokenForUser(user domain2.User) service.IIdentityService {
	s.claims = append(s.claims,
		jwt.RegisteredClaims{
			Subject: "user_id",
			ID:      strconv.FormatInt(user.ID, 10),
		},
		jwt.RegisteredClaims{
			Subject: "email",
			ID:      user.Email,
		},
		jwt.RegisteredClaims{
			Subject: "type",
			ID:      "user",
		},
	)

	if user.IsAdmin {
		s.claims = append(s.claims, jwt.RegisteredClaims{
			Subject: "is_admin",
			ID:      "true",
		})
	}

	return s
}

// TokenForCustomer creates a token for a customer
func (s *IdentityService) TokenForCustomer(customer domain2.Customer) service.IIdentityService {
	s.claims = append(s.claims,
		jwt.RegisteredClaims{
			Subject: "customer_id",
			ID:      strconv.FormatInt(customer.ID, 10),
		},
		jwt.RegisteredClaims{
			Subject: "email",
			ID:      customer.Email,
		},
		jwt.RegisteredClaims{
			Subject: "type",
			ID:      "customer",
		},
		jwt.RegisteredClaims{
			Subject: "site_id",
			ID:      strconv.FormatInt(customer.SiteID, 10),
		},
	)
	return s
}

// Make generates and returns the JWT token string
func (s *IdentityService) Make() string {
	// Create a map of claims from our slice
	mapClaims := jwt.MapClaims{
		"iss": s.jwtIssuer,
		"aud": s.jwtAudience,
		"exp": time.Now().Add(s.jwtExpiresIn).Unix(),
	}

	// Add all claims to the map
	for _, claim := range s.claims {
		if rc, ok := claim.(jwt.RegisteredClaims); ok {
			mapClaims[rc.Subject] = rc.ID
		}
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)

	// Sign the token
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		fmt.Printf("Error signing token: %v\n", err)
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

// GetToken extracts the JWT token from the request context
func (s *IdentityService) GetToken(c *gin.Context) (*jwt.Token, error) {
	// Get the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header not found")
	}

	// Check if it's a Bearer token
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("invalid authorization header format")
	}

	// Extract the token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse and verify the token
	return s.VerifyToken(tokenString)
}

// GetTokenString extracts the JWT token from the request context
func (s *IdentityService) GetTokenString(c *gin.Context) (string, error) {
	// Get the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}

	// Check if it's a Bearer token
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("invalid authorization header format")
	}

	// Extract the token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse and verify the token
	return tokenString, nil
}

// VerifyToken validates a token string and returns the parsed token
func (s *IdentityService) VerifyToken(tokenString string) (*jwt.Token, error) {
	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

// VerifyTokenContext validates a token string and returns the parsed token
func (s *IdentityService) VerifyTokenContext(c *gin.Context) (*jwt.Token, error) {
	// Parse and validate the token
	rawToken, err := s.GetTokenString(c)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

// GetClaim extracts a specific claim from a JWT token
func (s *IdentityService) GetClaim(ctx *gin.Context, claimName string) (string, error) {
	token, err := s.GetToken(ctx)
	if err != nil {
		return "", err
	}
	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	// Get the specific claim
	claimValue, exists := claims[claimName]
	if !exists {
		return "", fmt.Errorf("claim %s not found", claimName)
	}

	// Convert claim value to string
	claimStr, ok := claimValue.(string)
	if !ok {
		return "", fmt.Errorf("claim %s is not a string", claimName)
	}

	return claimStr, nil
}
