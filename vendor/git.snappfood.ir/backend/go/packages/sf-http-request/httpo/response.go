package httpo

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

// Response wraps the standard http.Response with utility methods
type Response struct {
	*http.Response
	bodyCache []byte
}

// Common errors
var (
	ErrNoBody = errors.New("response body is nil or already closed")
)

// WrapResponse wraps an http.Response with our Response type
func WrapResponse(resp *http.Response) *Response {
	if resp == nil {
		return nil
	}
	return &Response{Response: resp}
}

// readBody reads the body once and caches it for subsequent reads
func (r *Response) readBody() ([]byte, error) {
	if r.Response == nil || r.Body == nil {
		return nil, ErrNoBody
	}

	// If we've already read and cached the body, return it
	if r.bodyCache != nil {
		return r.bodyCache, nil
	}

	// Read the body
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// Cache the body
	r.bodyCache = body
	return body, nil
}

// JSON unmarshals the response body into the provided struct
func (r *Response) JSON(v interface{}) error {
	body, err := r.readBody()
	if err != nil {
		return err
	}
	
	return json.Unmarshal(body, v)
}

// String reads the response body and returns it as a string
func (r *Response) String() (string, error) {
	body, err := r.readBody()
	if err != nil {
		return "", err
	}
	
	return string(body), nil
}

// Bytes reads the response body and returns it as a byte slice
func (r *Response) Bytes() ([]byte, error) {
	return r.readBody()
}

// MustString reads the response body as a string, returning empty string on error
func (r *Response) MustString() string {
	s, err := r.String()
	if err != nil {
		return ""
	}
	return s
}

// MustJSON unmarshals the response body into the provided struct, returning whether it succeeded
func (r *Response) MustJSON(v interface{}) bool {
	err := r.JSON(v)
	return err == nil
}

// MustBytes reads the response body as bytes, returning nil on error
func (r *Response) MustBytes() []byte {
	b, err := r.Bytes()
	if err != nil {
		return nil
	}
	return b
}

// IsSuccess returns true if the response status code is in the 2xx range
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsError returns true if the response status code is in the 4xx or 5xx range
func (r *Response) IsError() bool {
	return r.StatusCode >= 400
}

// IsClientError returns true if the response status code is in the 4xx range
func (r *Response) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError returns true if the response status code is in the 5xx range
func (r *Response) IsServerError() bool {
	return r.StatusCode >= 500
}

// IsRedirect returns true if the response status code is in the 3xx range
func (r *Response) IsRedirect() bool {
	return r.StatusCode >= 300 && r.StatusCode < 400
}

// GetHeader returns the value of the specified header
func (r *Response) GetHeader(name string) string {
	if r.Response == nil || r.Header == nil {
		return ""
	}
	return r.Header.Get(name)
}

// GetHeaders returns all response headers as a map
func (r *Response) GetHeaders() map[string][]string {
	if r.Response == nil || r.Header == nil {
		return map[string][]string{}
	}
	return r.Header
}

// GetContentType returns the Content-Type header
func (r *Response) GetContentType() string {
	return r.GetHeader("Content-Type")
}

// IsJSON returns true if the Content-Type header indicates JSON
func (r *Response) IsJSON() bool {
	contentType := r.GetContentType()
	return strings.Contains(contentType, "application/json")
} 