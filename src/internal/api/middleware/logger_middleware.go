package middleware

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/gin-gonic/gin"
	"time"
)

func LoggerMiddleware(logger sflogger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)

		logger.Info("HTTP request", map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"client_ip":  c.ClientIP(),
			"latency":    latency,
			"user_agent": c.Request.UserAgent(),
			"data": map[string]interface{}{
				"test": "Dsadgfds",
				"aa":   "Aa",
			},
		},
		)
	}
}
