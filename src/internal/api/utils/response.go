package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/gin-gonic/gin"
)

// Result represents a standardized API response
type Result struct {
	Success       bool           `json:"success"`
	SystemMessage string         `json:"systemMessage,omitempty"`
	Message       string         `json:"message,omitempty"`
	Data          map[string]any `json:"data,omitempty"`
	StatusCode    int            `json:"statusCode"`
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
	MsgBadRequest      = "Invalid request"
)

// Response builders

// newResponse creates a new Result with the given parameters
func newResponse(success bool, systemMessage string, message string, statusCode int) *Result {
	return &Result{
		Success:       success,
		SystemMessage: systemMessage,
		Message:       message,
		StatusCode:    statusCode,
	}
}

// withData adds data to the Result
func (r *Result) withData(data map[string]any) *Result {
	r.Data = data
	return r
}

// withErrors adds error messages to the Result
func (r *Result) withErrors(error string) *Result {
	r.Message = error
	return r
}

// sendAndAbort sends the Result as a JSON response and aborts the request
func (r *Result) sendAndAbort(c *gin.Context) {
	c.AbortWithStatusJSON(r.StatusCode, r)
}

// Success responses

// Success creates a 200 Success response
func Success(c *gin.Context, msg string, data map[string]any) {
	newResponse(true, MsgSuccess, msg, http.StatusOK).
		withData(data).
		sendAndAbort(c)
}

// Created creates a 201 Created response
func Created(c *gin.Context, msg string, data map[string]any) {
	newResponse(true, MsgCreated, msg, http.StatusCreated).
		withData(data).
		sendAndAbort(c)
}

// Updated creates a 200 Success response for updates
func Updated(c *gin.Context, msg string, data map[string]any) {
	newResponse(true, MsgUpdated, msg, http.StatusOK).
		withData(data).
		sendAndAbort(c)
}

// Deleted creates a 200 Success response for deletions
func Deleted(c *gin.Context, msg string, data map[string]any) {
	newResponse(true, MsgDeleted, msg, http.StatusOK).
		withData(data).
		sendAndAbort(c)
}

// Retrieved creates a 200 Success response for data retrieval
func Retrieved(c *gin.Context, msg string, data map[string]any) {
	newResponse(true, MsgRetrieved, msg, http.StatusOK).
		withData(data).
		sendAndAbort(c)
}

// NoContent creates a 204 No Content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error responses

// BadRequest creates a 400 Bad Request response
func BadRequest(c *gin.Context, msg string, data map[string]any) {
	newResponse(false, MsgBadRequest, msg, http.StatusBadRequest).
		withData(data).
		sendAndAbort(c)
}

// ValidationError creates a 400 Bad Request response for validation errors
func ValidationError(c *gin.Context, errors ...ValidationErrorBag) {
	result := make(map[string]any, len(errors))
	for i, v := range errors {
		key := fmt.Sprintf("%d", i)
		result[key] = v
	}

	newResponse(false, MsgValidationError, "", http.StatusBadRequest).
		withData(result).
		sendAndAbort(c)
}

// ValidationErrorString creates a 400 Bad Request response for validation errors
func ValidationErrorString(c *gin.Context, msg string, data map[string]any) {
	newResponse(false, MsgValidationError, msg, http.StatusBadRequest).
		withData(data).
		sendAndAbort(c)
}

// Unauthorized creates a 401 Unauthorized response
func Unauthorized(c *gin.Context, msg string, data map[string]any) {
	newResponse(false, MsgUnauthorized, msg, http.StatusUnauthorized).
		withData(data).
		sendAndAbort(c)
}

// NotFound creates a 404 Not Found response
func NotFound(c *gin.Context, msg string, data map[string]any) {
	newResponse(false, MsgNotFound, msg, http.StatusNotFound).
		withData(data).
		sendAndAbort(c)
}

// InternalError creates a 500 Internal Server Error response
func InternalError(c *gin.Context, msg string, data map[string]any) {
	newResponse(false, MsgInternalError, msg, http.StatusInternalServerError).
		withData(data).
		sendAndAbort(c)
}

// Error handling helpers

// HandleError is a convenience function to handle errors with appropriate responses
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var nerr *resp.Error
	if errors.As(err, &nerr) {
		switch nerr.Type {
		case resp.NotFound:
			NotFound(c, nerr.Message, nerr.Data)
		case resp.Unauthorized:
			Unauthorized(c, nerr.Message, nerr.Data)
		case resp.Validation:
			ValidationErrorString(c, nerr.Message, nerr.Data)
		case resp.BadRequest:
			BadRequest(c, nerr.Message, nerr.Data)
		case resp.Internal:
			InternalError(c, nerr.Message, nerr.Data)
		default:
			InternalError(c, nerr.Message, nerr.Data)
		}
		return
	}

	// fallback for generic errors
	InternalError(c, err.Error(), nil)
}

// HandleResponse is a convenience function to handle response with appropriate responses
func HandleResponse(c *gin.Context, response *resp.Response) {
	switch response.Type {
	case resp.Success:
		Success(c, response.Message, response.Data)
	case resp.Created:
		Created(c, response.Message, response.Data)
	case resp.Updated:
		Updated(c, response.Message, response.Data)
	case resp.Deleted:
		Deleted(c, response.Message, response.Data)
	case resp.Retrieved:
		Retrieved(c, response.Message, response.Data)
	case resp.NoContent:
		NoContent(c)
	default:
		NoContent(c)
	}
	return
}
