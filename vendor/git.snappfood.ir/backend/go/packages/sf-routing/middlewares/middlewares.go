package middlewares

import (
	"github.com/gin-gonic/gin"
)

// Middleware is a function that takes a gin.HandlerFunc and returns a gin.HandlerFunc
type Middleware func(gin.HandlerFunc) gin.HandlerFunc

// LightMode middleware for light mode
func LightMode(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement light mode logic here
		handler(c)
	}
}

// ABTest middleware for A/B testing
func ABTest(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement A/B testing logic here
		handler(c)
	}
}

// CheckDisasterMiddle middleware for disaster checking
func CheckDisasterMiddle(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement disaster checking logic here
		handler(c)
	}
}

// LoggedinMiddle middleware for checking if user is logged in
func LoggedinMiddle(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement logged in check logic here
		handler(c)
	}
}

// Mix combines a handler with middleware
func Mix(handler gin.HandlerFunc, middleware ...Middleware) gin.HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	return handler
}
