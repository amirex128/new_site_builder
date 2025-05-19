package middleware

import (
	"errors"

	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"
	"github.com/gin-gonic/gin"
)

type Authenticator struct {
	authTransientService func(c *gin.Context) service.IAuthService
	identityService      service.IIdentityService
}

func NewAuthenticator(authTransientService func(c *gin.Context) service.IAuthService, identityService service.IIdentityService) *Authenticator {
	return &Authenticator{
		authTransientService: authTransientService,
		identityService:      identityService,
	}
}

// Authenticate verifies the token without checking roles
func (a *Authenticator) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := a.identityService.VerifyTokenContext(c)
		if err != nil {
			resp.ErrorUnauthorized(c, err)
			c.Abort()
			return
		}

		c.Next()
	}
}

// MustRole requires the user to have ALL specified roles (AND logic)
func (a *Authenticator) MustRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := a.identityService.VerifyTokenContext(c)
		if err != nil {
			resp.ErrorUnauthorized(c, err)
			c.Abort()
			return
		}

		authService := a.authTransientService(c)
		userRoles, err := authService.GetRoles()
		if err != nil {
			resp.ErrorUnauthorized(c, err)
			c.Abort()
			return
		}

		// Check if user is admin (admins have all permissions)
		isAdmin, err := authService.IsAdmin()
		if err == nil && isAdmin {
			c.Next()
			return
		}

		// Check if user has ALL required roles
		for _, requiredRole := range roles {
			hasRole := false
			for _, userRole := range userRoles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if !hasRole {
				resp.ErrorForbidden(c, errors.New("insufficient permissions"))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// ShouldRole requires the user to have AT LEAST ONE of the specified roles (OR logic)
func (a *Authenticator) ShouldRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := a.identityService.VerifyTokenContext(c)
		if err != nil {
			resp.ErrorUnauthorized(c, err)
			c.Abort()
			return
		}

		authService := a.authTransientService(c)
		userRoles, err := authService.GetRoles()
		if err != nil {
			resp.ErrorUnauthorized(c, err)
			c.Abort()
			return
		}

		// Check if user is admin (admins have all permissions)
		isAdmin, err := authService.IsAdmin()
		if err == nil && isAdmin {
			c.Next()
			return
		}

		// Check if user has AT LEAST ONE of the required roles
		if len(roles) > 0 {
			hasAnyRole := false
			for _, requiredRole := range roles {
				for _, userRole := range userRoles {
					if userRole == requiredRole {
						hasAnyRole = true
						break
					}
				}
				if hasAnyRole {
					break
				}
			}

			if !hasAnyRole {
				resp.ErrorForbidden(c, errors.New("insufficient permissions"))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
