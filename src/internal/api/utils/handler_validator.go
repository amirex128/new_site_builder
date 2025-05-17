package utils

import (
	"net/http"

	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidationHelper provides methods for handling validation in handlers
type ValidationHelper struct {
	validate *validator.Validate
}

// NewValidationHelper creates a new ValidationHelper instance
func NewValidationHelper() *ValidationHelper {
	return &ValidationHelper{
		validate: validator.New(),
	}
}

// ValidateRequest handles binding and validating the input struct
// Returns true if validation passes, false otherwise
func (h *ValidationHelper) ValidateRequest(c *gin.Context, params interface{}) bool {
	if err := c.ShouldBindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, resp.ValidationFailed().WithSystemMessage(err.Error()))
		return false
	}

	if err := h.validate.Struct(params); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, resp.ValidationFailed().WithSystemMessage(validationErrors.Error()))
		return false
	}

	return true
}
