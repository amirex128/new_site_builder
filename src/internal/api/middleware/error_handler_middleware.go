package middleware

import (
	"fmt"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware(logger sflogger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic if logger is available
				if logger != nil {
					extraMap := map[string]interface{}{
						sflogger.ExtraKey.Error.ErrorMessage: fmt.Sprintf("%v", err),
					}

					logger.ErrorWithCategory(
						sflogger.Category.Error.Error,
						sflogger.SubCategory.Status.Error,
						"Panic recovered",
						extraMap,
					)
				}

				// Return error response
				if !c.Writer.Written() {
					utils.InternalError(c, fmt.Sprintf("%v", err))
				}
			}
		}()

		c.Next()
	}
}
