package utils

import (
	"errors"
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

// ValidateCommand handles binding and validating the input struct
// Returns true if validation passes, false otherwise
func (h *ValidationHelper) ValidateCommand(c *gin.Context, params interface{}) bool {
	if err := c.ShouldBindJSON(params); err != nil {
		resp.ValidationError(c, err.Error())
		return false
	}

	if err := h.validate.Struct(params); err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)
		resp.ValidationError(c, validationErrors.Error())
		return false
	}

	return true
}

// ValidateQuery handles binding and validating query parameters
// Returns true if validation passes, false otherwise
func (h *ValidationHelper) ValidateQuery(c *gin.Context, params interface{}) bool {
	if err := c.ShouldBindQuery(params); err != nil {
		resp.ValidationError(c, err.Error())
		return false
	}

	if err := h.validate.Struct(params); err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)
		resp.ValidationError(c, validationErrors.Error())
		return false
	}

	return true
}
