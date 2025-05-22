package utils

import (
	"errors"
	"net/http"

	"github.com/amirex128/new_site_builder/src/internal/application/utils/nerror"
	"github.com/gin-gonic/gin"
)

// Result represents a standardized API response
type Result struct {
	Success    bool     `json:"success"`
	Message    string   `json:"message,omitempty"`
	Data       any      `json:"data,omitempty"`
	Errors     []string `json:"errors,omitempty"`
	StatusCode int      `json:"statusCode"`
	ErrorData  []any    `json:"errorData,omitempty"`
}

// Standard messages for API responses
const (
	MsgSuccess         = "Operation completed successfully"
	MsgCreated         = "Resource created successfully"
	MsgUpdated         = "Resource updated successfully"
	MsgDeleted         = "Resource deleted successfully"
	MsgRetrieved       = "Data retrieved successfully"
	MsgValidationError = "Validation error"
	MsgInternalError   = "Internal server error"
	MsgNotFound        = "Resource not found"
	MsgUnauthorized    = "Authentication required"
	MsgForbidden       = "Insufficient permissions"
	MsgBadRequest      = "Invalid request"
	MsgConflict        = "Resource conflict"
)

// Response builders

// NewResponse creates a new Result with the given parameters
func NewResponse(success bool, message string, statusCode int) *Result {
	return &Result{
		Success:    success,
		Message:    message,
		StatusCode: statusCode,
	}
}

// WithData adds data to the Result
func (r *Result) WithData(data any) *Result {
	r.Data = data
	return r
}

// WithAnyErrors adds error messages to the Result
func (r *Result) WithAnyErrors(errors ...any) *Result {
	r.ErrorData = append(r.ErrorData, errors...)
	return r
}

// WithErrors adds error messages to the Result
func (r *Result) WithErrors(errors ...string) *Result {
	r.Errors = append(r.Errors, errors...)
	return r
}

// Send sends the Result as a JSON response with the appropriate status code
func (r *Result) Send(c *gin.Context) {
	c.JSON(r.StatusCode, r)
}

// SendAndAbort sends the Result as a JSON response and aborts the request
func (r *Result) SendAndAbort(c *gin.Context) {
	c.AbortWithStatusJSON(r.StatusCode, r)
}

// Success responses

// OK creates a 200 OK response
func OK(c *gin.Context, data any) {
	NewResponse(true, MsgSuccess, http.StatusOK).
		WithData(data).
		Send(c)
}

// Created creates a 201 Created response
func Created(c *gin.Context, data any) {
	NewResponse(true, MsgCreated, http.StatusCreated).
		WithData(data).
		Send(c)
}

// Updated creates a 200 OK response for updates
func Updated(c *gin.Context, data any) {
	NewResponse(true, MsgUpdated, http.StatusOK).
		WithData(data).
		Send(c)
}

// Deleted creates a 200 OK response for deletions
func Deleted(c *gin.Context, data ...any) {
	result := NewResponse(true, MsgDeleted, http.StatusOK)
	if len(data) > 0 && data[0] != nil {
		result.WithData(data[0])
	}
	result.Send(c)
}

// Retrieved creates a 200 OK response for data retrieval
func Retrieved(c *gin.Context, data any) {
	NewResponse(true, MsgRetrieved, http.StatusOK).
		WithData(data).
		Send(c)
}

// NoContent creates a 204 No Content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error responses

// BadRequest creates a 400 Bad Request response
func BadRequest(c *gin.Context, errors ...string) {
	NewResponse(false, MsgBadRequest, http.StatusBadRequest).
		WithErrors(errors...).
		SendAndAbort(c)
}

// ValidationError creates a 400 Bad Request response for validation errors
func ValidationError(c *gin.Context, errors ...ValidationErrorBag) {
	// Convert ValidationErrorBag slice to []any for WithAnyErrors
	anyErrors := make([]any, len(errors))
	for i, err := range errors {
		anyErrors[i] = err
	}

	NewResponse(false, MsgValidationError, http.StatusBadRequest).
		WithAnyErrors(anyErrors...).
		SendAndAbort(c)
}

// ValidationErrorString creates a 400 Bad Request response for validation errors
func ValidationErrorString(c *gin.Context, errors ...string) {
	NewResponse(false, MsgValidationError, http.StatusBadRequest).
		WithErrors(errors...).
		SendAndAbort(c)
}

// Unauthorized creates a 401 Unauthorized response
func Unauthorized(c *gin.Context, errors ...string) {
	NewResponse(false, MsgUnauthorized, http.StatusUnauthorized).
		WithErrors(errors...).
		SendAndAbort(c)
}

// Forbidden creates a 403 Forbidden response
func Forbidden(c *gin.Context, errors ...string) {
	NewResponse(false, MsgForbidden, http.StatusForbidden).
		WithErrors(errors...).
		SendAndAbort(c)
}

// NotFound creates a 404 Not Found response
func NotFound(c *gin.Context, errors ...string) {
	NewResponse(false, MsgNotFound, http.StatusNotFound).
		WithErrors(errors...).
		SendAndAbort(c)
}

// Conflict creates a 409 Conflict response
func Conflict(c *gin.Context, errors ...string) {
	NewResponse(false, MsgConflict, http.StatusConflict).
		WithErrors(errors...).
		SendAndAbort(c)
}

// InternalError creates a 500 Internal Server Error response
func InternalError(c *gin.Context, errors ...string) {
	if len(errors) == 0 {
		errors = []string{MsgInternalError}
	}
	NewResponse(false, MsgInternalError, http.StatusInternalServerError).
		WithErrors(errors...).
		SendAndAbort(c)
}

// Error handling helpers

// HandleError is a convenience function to handle errors with appropriate responses
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var nerr *nerror.NError
	if errors.As(err, &nerr) {
		switch nerr.Type {
		case nerror.NotFound:
			NotFound(c, nerr.Message)
		case nerror.Unauthorized:
			Unauthorized(c, nerr.Message)
		case nerror.BadRequest:
			BadRequest(c, nerr.Message)
		case nerror.Internal:
			InternalError(c, nerr.Message)
		default:
			InternalError(c, nerr.Message)
		}
		return
	}

	// fallback for generic errors
	InternalError(c, err.Error())
}
