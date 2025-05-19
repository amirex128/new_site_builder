package resp

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// Result represents a standardized API response
type Result struct {
	Succeeded         bool        `json:"succeeded"`
	Messages          []string    `json:"messages,omitempty"`
	SystemMessages    []string    `json:"systemMessages,omitempty"`
	Data              interface{} `json:"data,omitempty"`
	StatusCode        int         `json:"statusCode"`
	Type              string      `json:"type,omitempty"`
	SystemRedirectURL string      `json:"systemRedirectUrl,omitempty"`
}

// newSuccessResult creates a new successful Result
func newSuccessResult(messages ...string) *Result {
	result := &Result{
		Succeeded:  true,
		StatusCode: http.StatusOK,
	}
	result.Messages = append(result.Messages, messages...)
	return result
}

// newFailureResult creates a new failure Result
func newFailureResult(messages ...string) *Result {
	result := &Result{
		Succeeded:  false,
		StatusCode: http.StatusInternalServerError,
	}
	result.Messages = append(result.Messages, messages...)
	return result
}

// newRedirectResult creates a new redirect Result
func newRedirectResult(redirectURL string, queryParams map[string]string) *Result {
	if len(queryParams) > 0 {
		values := url.Values{}
		for k, v := range queryParams {
			values.Add(k, v)
		}

		// Append query parameters to URL
		if strings.Contains(redirectURL, "?") {
			redirectURL = redirectURL + "&" + values.Encode()
		} else {
			redirectURL = redirectURL + "?" + values.Encode()
		}
	}

	return &Result{
		Succeeded:         true,
		SystemRedirectURL: redirectURL,
		Type:              "Redirect",
	}
}

// WithMessage adds message(s) to the Result
func (r *Result) WithMessage(messages ...string) *Result {
	r.Messages = append(r.Messages, messages...)
	return r
}

// WithSystemMessage adds system message(s) to the Result
func (r *Result) WithSystemMessage(messages ...string) *Result {
	r.SystemMessages = append(r.SystemMessages, messages...)
	return r
}

// WithData adds data to the Result
func (r *Result) WithData(data interface{}) *Result {
	r.Data = data
	return r
}

// WithStatusCode sets the status code for the Result
func (r *Result) WithStatusCode(statusCode int) *Result {
	r.StatusCode = statusCode
	return r
}

// WithType sets the type for the Result
func (r *Result) WithType(resultType string) *Result {
	r.Type = resultType
	return r
}

// GetAllMessages returns all messages (system and user-facing)
func (r *Result) GetAllMessages() []string {
	allMessages := make([]string, 0, len(r.SystemMessages)+len(r.Messages))
	allMessages = append(allMessages, r.SystemMessages...)
	allMessages = append(allMessages, r.Messages...)
	return allMessages
}

// resultMessages contains standard response messages
var resultMessages = struct {
	Success           string
	Created           string
	Updated           string
	Deleted           string
	Retrieved         string
	VerifySuccess     string
	FailedCreate      string
	FailedUpdate      string
	FailedDelete      string
	FailedRetrieve    string
	ValidationError   string
	InternalError     string
	NotFoundError     string
	AuthenticateError string
	AuthorizeError    string
	VerifyError       string
	OperationError    string
}{
	Success:           "Operation completed successfully",
	Created:           "Record created successfully",
	Updated:           "Record updated successfully",
	Deleted:           "Record deleted successfully",
	Retrieved:         "Data retrieved successfully",
	VerifySuccess:     "Verification successful",
	FailedCreate:      "Failed to create record",
	FailedUpdate:      "Failed to update record",
	FailedDelete:      "Failed to delete record",
	FailedRetrieve:    "Failed to retrieve data",
	ValidationError:   "Validation error",
	InternalError:     "Internal server error",
	NotFoundError:     "Resource not found",
	AuthenticateError: "Authentication error",
	AuthorizeError:    "You do not have permission for this operation",
	VerifyError:       "Verification failed",
	OperationError:    "An error occurred during the operation",
}

// Factory functions for common response types

// Success creates a standard success response
func Success() *Result {
	return newSuccessResult().
		WithSystemMessage(resultMessages.Success).
		WithStatusCode(http.StatusOK).
		WithType("Success")
}

// Created creates a response for successful resource creation
func Created() *Result {
	return newSuccessResult().
		WithSystemMessage(resultMessages.Created).
		WithStatusCode(http.StatusCreated).
		WithType("Success")
}

// Updated creates a response for successful resource update
func Updated() *Result {
	return newSuccessResult().
		WithSystemMessage(resultMessages.Updated).
		WithStatusCode(http.StatusOK).
		WithType("Success")
}

// Deleted creates a response for successful resource deletion
func Deleted() *Result {
	return newSuccessResult().
		WithSystemMessage(resultMessages.Deleted).
		WithStatusCode(http.StatusOK).
		WithType("Success")
}

// Retrieved creates a response for successful data retrieval
func Retrieved() *Result {
	return newSuccessResult().
		WithSystemMessage(resultMessages.Retrieved).
		WithStatusCode(http.StatusOK).
		WithType("Success")
}

// ValidationError creates a response for validation errors
func ValidationFailed() *Result {
	return newFailureResult().
		WithSystemMessage(resultMessages.ValidationError).
		WithStatusCode(http.StatusBadRequest).
		WithType("ValidationError")
}

// NotFoundError creates a response for resource not found
func NotFoundError() *Result {
	return newFailureResult().
		WithSystemMessage(resultMessages.NotFoundError).
		WithStatusCode(http.StatusNotFound).
		WithType("NotFound")
}

// InternalError creates a response for internal server errors
func InternalError() *Result {
	return newFailureResult().
		WithSystemMessage(resultMessages.InternalError).
		WithStatusCode(http.StatusInternalServerError).
		WithType("Error")
}

// AuthenticateError creates a response for authentication failures
func AuthenticateError() *Result {
	return newFailureResult().
		WithSystemMessage(resultMessages.AuthenticateError).
		WithStatusCode(http.StatusUnauthorized).
		WithType("Error")
}

// AuthorizeError creates a response for authorization failures
func AuthorizeError() *Result {
	return newFailureResult().
		WithSystemMessage(resultMessages.AuthorizeError).
		WithStatusCode(http.StatusUnauthorized).
		WithType("Error")
}

// Redirect creates a redirect response
func Redirect(redirectURL string, queryParams map[string]string) *Result {
	return newRedirectResult(redirectURL, queryParams)
}

// ErrorUnauthorized sends an unauthorized error response
func ErrorUnauthorized(c *gin.Context, err error) {
	result := AuthenticateError()
	if err != nil {
		result.WithSystemMessage(err.Error())
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, result)
}

// ErrorForbidden sends a forbidden error response
func ErrorForbidden(c *gin.Context, err error) {
	result := newFailureResult().
		WithSystemMessage(resultMessages.AuthorizeError).
		WithStatusCode(http.StatusForbidden).
		WithType("Error")

	if err != nil {
		result.WithSystemMessage(err.Error())
	}

	c.AbortWithStatusJSON(http.StatusForbidden, result)
}
