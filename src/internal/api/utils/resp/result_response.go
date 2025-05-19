package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Result represents a standardized API response
type Result struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Errors     []string    `json:"errors,omitempty"`
	StatusCode int         `json:"statusCode"`
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
func (r *Result) WithData(data interface{}) *Result {
	r.Data = data
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
func OK(c *gin.Context, data interface{}) {
	NewResponse(true, MsgSuccess, http.StatusOK).
		WithData(data).
		Send(c)
}

// Created creates a 201 Created response
func Created(c *gin.Context, data interface{}) {
	NewResponse(true, MsgCreated, http.StatusCreated).
		WithData(data).
		Send(c)
}

// Updated creates a 200 OK response for updates
func Updated(c *gin.Context, data interface{}) {
	NewResponse(true, MsgUpdated, http.StatusOK).
		WithData(data).
		Send(c)
}

// Deleted creates a 200 OK response for deletions
func Deleted(c *gin.Context, data ...interface{}) {
	result := NewResponse(true, MsgDeleted, http.StatusOK)
	if len(data) > 0 && data[0] != nil {
		result.WithData(data[0])
	}
	result.Send(c)
}

// Retrieved creates a 200 OK response for data retrieval
func Retrieved(c *gin.Context, data interface{}) {
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
func ValidationError(c *gin.Context, errors ...string) {
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
func HandleError(c *gin.Context, err error, statusCode int, message string) {
	if err == nil {
		return
	}

	errMsg := err.Error()
	NewResponse(false, message, statusCode).
		WithErrors(errMsg).
		SendAndAbort(c)
}

// ErrorUnauthorized handles unauthorized errors (for backward compatibility)
func ErrorUnauthorized(c *gin.Context, err error) {
	errMsgs := []string{}
	if err != nil {
		errMsgs = append(errMsgs, err.Error())
	}
	Unauthorized(c, errMsgs...)
}

// ErrorForbidden handles forbidden errors (for backward compatibility)
func ErrorForbidden(c *gin.Context, err error) {
	errMsgs := []string{}
	if err != nil {
		errMsgs = append(errMsgs, err.Error())
	}
	Forbidden(c, errMsgs...)
}
