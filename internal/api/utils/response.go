package utils

import (
	"errors"
	"fmt"
	"github.com/amirex128/new_site_builder/internal/application/utils/resp"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Result represents a standardized API response
type Result struct {
	Success       bool   `json:"success"`
	SystemMessage string `json:"systemMessage,omitempty"`
	Message       string `json:"message,omitempty"`
	Data          any    `json:"data,omitempty"`
	StatusCode    int    `json:"statusCode"`
}

// پیام‌های استاندارد برای پاسخ‌های API
const (
	MsgSuccess         = "عملیات با موفقیت انجام شد"
	MsgCreated         = "منبع با موفقیت ایجاد شد"
	MsgUpdated         = "با موفقیت بروزرسانی شد"
	MsgDeleted         = "با موفقیت حذف شد"
	MsgRetrieved       = "داده‌ها با موفقیت دریافت شدند"
	MsgValidationError = "خطای اعتبارسنجی"
	MsgInternalError   = "خطای داخلی سرور"
	MsgNotFound        = "منبع یافت نشد"
	MsgUnauthorized    = "نیاز به احراز هویت"
	MsgBadRequest      = "درخواست نامعتبر"
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
func (r *Result) withData(data any) *Result {
	r.Data = data
	return r
}

// sendAndAbort sends the Result as a JSON response and aborts the request
func (r *Result) sendAndAbort(c *gin.Context) {
	c.AbortWithStatusJSON(r.StatusCode, r)
}

// success responses

// success creates a 200 success response
func success(c *gin.Context, msg string, data any) {
	newResponse(true, MsgSuccess, msg, http.StatusOK).
		withData(data).
		sendAndAbort(c)
}

// Created creates a 201 Created response
func created(c *gin.Context, msg string, data any) {
	newResponse(true, MsgCreated, msg, http.StatusCreated).
		withData(data).
		sendAndAbort(c)
}

// Updated creates a 200 success response for updates
func updated(c *gin.Context, msg string, data any) {
	newResponse(true, MsgUpdated, msg, http.StatusOK).
		withData(data).
		sendAndAbort(c)
}

// Deleted creates a 200 success response for deletions
func deleted(c *gin.Context, msg string, data any) {
	newResponse(true, MsgDeleted, msg, http.StatusOK).
		withData(data).
		sendAndAbort(c)
}

// Retrieved creates a 200 success response for data retrieval
func retrieved(c *gin.Context, msg string, data any) {
	newResponse(true, MsgRetrieved, msg, http.StatusOK).
		withData(data).
		sendAndAbort(c)
}

// NoContent creates a 204 No Content response
func noContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error responses

// badRequest creates a 400 Bad Request response
func badRequest(c *gin.Context, msg string, data any) {
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

// validationErrorString creates a 400 Bad Request response for validation errors
func validationErrorString(c *gin.Context, msg string, data any) {
	newResponse(false, MsgValidationError, msg, http.StatusBadRequest).
		withData(data).
		sendAndAbort(c)
}

// unauthorized creates a 401 unauthorized response
func unauthorized(c *gin.Context, msg string, data any) {
	newResponse(false, MsgUnauthorized, msg, http.StatusUnauthorized).
		withData(data).
		sendAndAbort(c)
}

// notFound creates a 404 Not Found response
func notFound(c *gin.Context, msg string, data any) {
	newResponse(false, MsgNotFound, msg, http.StatusNotFound).
		withData(data).
		sendAndAbort(c)
}

// internalError creates a 500 Internal Server Error response
func internalError(c *gin.Context, msg string, data any) {
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
			notFound(c, nerr.Message, nerr.Data)
		case resp.Unauthorized:
			unauthorized(c, nerr.Message, nerr.Data)
		case resp.Validation:
			validationErrorString(c, nerr.Message, nerr.Data)
		case resp.BadRequest:
			badRequest(c, nerr.Message, nerr.Data)
		case resp.Internal:
			internalError(c, nerr.Message, nerr.Data)
		default:
			internalError(c, nerr.Message, nerr.Data)
		}
		return
	}

	// fallback for generic errors
	internalError(c, err.Error(), nil)
}

// HandleResponse is a convenience function to handle response with appropriate responses
func HandleResponse(c *gin.Context, response *resp.Response) {
	if response == nil {
		return
	}
	switch response.Type {
	case resp.Success:
		success(c, response.Message, response.Data)
	case resp.Created:
		created(c, response.Message, response.Data)
	case resp.Updated:
		updated(c, response.Message, response.Data)
	case resp.Deleted:
		deleted(c, response.Message, response.Data)
	case resp.Retrieved:
		retrieved(c, response.Message, response.Data)
	case resp.NoContent:
		noContent(c)
	default:
		noContent(c)
	}
	return
}
