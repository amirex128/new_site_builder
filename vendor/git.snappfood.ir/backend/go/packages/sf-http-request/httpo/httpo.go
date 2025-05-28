package httpo

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	sfhttprequest "git.snappfood.ir/backend/go/packages/sf-http-request"
	"go.elastic.co/apm/module/apmhttp/v2"
)

// Error constants
var (
	// ErrCircuitOpen is left for backward compatibility
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

// ClientOption is a function to customize the http.Client
type ClientOption func(*http.Client)

// WithTransport sets a custom transport for the client
func WithTransport(transport http.RoundTripper) ClientOption {
	return func(c *http.Client) {
		c.Transport = transport
	}
}

// WithCookieJar sets a custom cookie jar for the client
func WithCookieJar(jar http.CookieJar) ClientOption {
	return func(c *http.Client) {
		c.Jar = jar
	}
}

// WithCheckRedirect sets a custom redirect policy for the client
func WithCheckRedirect(checkRedirect func(req *http.Request, via []*http.Request) error) ClientOption {
	return func(c *http.Client) {
		c.CheckRedirect = checkRedirect
	}
}

// WithTimeout sets a default timeout for all requests made with this client
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *http.Client) {
		c.Timeout = timeout
	}
}

// Request represents an HTTP request with additional features like timeout, retry, fallback
type Request struct {
	client       *http.Client
	request      *http.Request
	method       string
	url          string
	body         []byte
	headers      map[string]string
	queryParams  url.Values
	formData     url.Values
	timeoutMs    int
	retryCount   int
	retryDelayMs int
	fallbackFunc func()
	ctx          context.Context
	logger       sfhttprequest.Logger
}

// Get creates a new GET request using a registered connection
func Get(serviceName, path string) *Request {
	r, err := FromConnection(serviceName)
	if err != nil {
		// Return a request that will fail when sent
		return &Request{
			method:      http.MethodGet,
			headers:     make(map[string]string),
			queryParams: url.Values{},
			formData:    url.Values{},
			ctx:         context.Background(),
		}
	}
	r.method = http.MethodGet
	r.url = r.url + path
	return r
}

// Post creates a new POST request using a registered connection
func Post(serviceName, path string) *Request {
	r, err := FromConnection(serviceName)
	if err != nil {
		// Return a request that will fail when sent
		return &Request{
			method:      http.MethodPost,
			headers:     make(map[string]string),
			queryParams: url.Values{},
			formData:    url.Values{},
			ctx:         context.Background(),
		}
	}
	r.method = http.MethodPost
	r.url = r.url + path
	return r
}

// Put creates a new PUT request using a registered connection
func Put(serviceName, path string) *Request {
	r, err := FromConnection(serviceName)
	if err != nil {
		// Return a request that will fail when sent
		return &Request{
			method:      http.MethodPut,
			headers:     make(map[string]string),
			queryParams: url.Values{},
			formData:    url.Values{},
			ctx:         context.Background(),
		}
	}
	r.method = http.MethodPut
	r.url = r.url + path
	return r
}

// Delete creates a new DELETE request using a registered connection
func Delete(serviceName, path string) *Request {
	r, err := FromConnection(serviceName)
	if err != nil {
		// Return a request that will fail when sent
		return &Request{
			method:      http.MethodDelete,
			headers:     make(map[string]string),
			queryParams: url.Values{},
			formData:    url.Values{},
			ctx:         context.Background(),
		}
	}
	r.method = http.MethodDelete
	r.url = r.url + path
	return r
}

// Patch creates a new PATCH request using a registered connection
func Patch(serviceName, path string) *Request {
	r, err := FromConnection(serviceName)
	if err != nil {
		// Return a request that will fail when sent
		return &Request{
			method:      http.MethodPatch,
			headers:     make(map[string]string),
			queryParams: url.Values{},
			formData:    url.Values{},
			ctx:         context.Background(),
		}
	}
	r.method = http.MethodPatch
	r.url = r.url + path
	return r
}

// Head creates a new HEAD request using a registered connection
func Head(serviceName, path string) *Request {
	r, err := FromConnection(serviceName)
	if err != nil {
		// Return a request that will fail when sent
		return &Request{
			method:      http.MethodHead,
			headers:     make(map[string]string),
			queryParams: url.Values{},
			formData:    url.Values{},
			ctx:         context.Background(),
		}
	}
	r.method = http.MethodHead
	r.url = r.url + path
	return r
}

// Options creates a new OPTIONS request using a registered connection
func Options(serviceName, path string) *Request {
	r, err := FromConnection(serviceName)
	if err != nil {
		// Return a request that will fail when sent
		return &Request{
			method:      http.MethodOptions,
			headers:     make(map[string]string),
			queryParams: url.Values{},
			formData:    url.Values{},
			ctx:         context.Background(),
		}
	}
	r.method = http.MethodOptions
	r.url = r.url + path
	return r
}

// Body sets the request body
func (r *Request) Body(body []byte) *Request {
	r.body = body
	return r
}

// JSONBody sets the request body as JSON and sets the Content-Type header to application/json
func (r *Request) JSONBody(data interface{}) *Request {
	jsonData, err := json.Marshal(data)
	if err != nil {
		// Store empty body in case of error
		r.body = []byte{}
		return r
	}
	r.body = jsonData
	r.Header("Content-Type", "application/json")
	return r
}

// FormData adds form data to the request and sets the Content-Type header to application/x-www-form-urlencoded
func (r *Request) FormData(key, value string) *Request {
	r.formData.Add(key, value)
	r.Header("Content-Type", "application/x-www-form-urlencoded")
	return r
}

// Query adds a query parameter to the URL
func (r *Request) Query(key, value string) *Request {
	r.queryParams.Add(key, value)
	return r
}

// Header adds a header to the request
func (r *Request) Header(key, value string) *Request {
	r.headers[key] = value
	return r
}

// SetHeaders adds multiple headers to the request
func (r *Request) SetHeaders(headers map[string]string) *Request {
	for key, value := range headers {
		r.headers[key] = value
	}
	return r
}

// GetHeaders returns all request headers
func (r *Request) GetHeaders() map[string]string {
	return r.headers
}

// SetContext sets a custom context for the request
func (r *Request) SetContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

// Timeout sets the request timeout in milliseconds
func (r *Request) Timeout(timeoutMs int) *Request {
	r.timeoutMs = timeoutMs
	return r
}

// Retry sets the retry count and delay in milliseconds
func (r *Request) Retry(count, delayMs int) *Request {
	r.retryCount = count
	r.retryDelayMs = delayMs
	return r
}

// Fallback sets the fallback function to be called when all retries fail
func (r *Request) Fallback(fn func()) *Request {
	r.fallbackFunc = fn
	return r
}

// BasicAuth sets the basic authentication credentials
func (r *Request) BasicAuth(username, password string) *Request {
	auth := username + ":" + password
	basicAuth := "Basic " + base64Encode(auth)
	r.Header("Authorization", basicAuth)
	return r
}

// BearerAuth sets the bearer token for authentication
func (r *Request) BearerAuth(token string) *Request {
	r.Header("Authorization", "Bearer "+token)
	return r
}

// base64Encode encodes a string to base64
func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// Method sets the HTTP method for the request
func (r *Request) Method(method string) *Request {
	r.method = method
	return r
}

// URL sets the URL for the request
func (r *Request) URL(url string) *Request {
	r.url = url
	return r
}

// Send executes the HTTP request with all configured features
func (r *Request) Send() (*Response, error) {
	req, err := r.buildRequest()
	if err != nil {
		if r.fallbackFunc != nil {
			r.fallbackFunc()
		}

		if r.logger != nil {
			r.logger.ErrorWithCategory(
				sfhttprequest.Category.API.HTTP,
				sfhttprequest.SubCategory.Status.Error,
				"Failed to build HTTP request",
				map[string]interface{}{
					sfhttprequest.ExtraKey.HTTP.URL:           r.url,
					sfhttprequest.ExtraKey.HTTP.Method:        r.method,
					sfhttprequest.ExtraKey.Error.ErrorMessage: err.Error(),
				})
		}

		return nil, err
	}

	// Apply timeout
	var ctx context.Context
	var cancel context.CancelFunc
	if r.timeoutMs > 0 {
		ctx, cancel = context.WithTimeout(r.ctx, time.Duration(r.timeoutMs)*time.Millisecond)
		req = req.WithContext(ctx)
		defer cancel()
	} else {
		req = req.WithContext(r.ctx)
	}

	startTime := time.Now()

	// Log the outgoing request
	if r.logger != nil {
		r.logger.DebugWithCategory(
			sfhttprequest.Category.API.HTTP,
			sfhttprequest.SubCategory.API.Request,
			"Sending HTTP request",
			map[string]interface{}{
				sfhttprequest.ExtraKey.HTTP.URL:    r.url,
				sfhttprequest.ExtraKey.HTTP.Method: r.method,
			})
	}

	// Retry logic
	var resp *http.Response
	for attempt := 0; attempt <= r.retryCount; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(r.retryDelayMs) * time.Millisecond)

			if r.logger != nil {
				r.logger.DebugWithCategory(
					sfhttprequest.Category.API.HTTP,
					sfhttprequest.SubCategory.API.Request,
					"Retrying HTTP request",
					map[string]interface{}{
						sfhttprequest.ExtraKey.HTTP.URL:    r.url,
						sfhttprequest.ExtraKey.HTTP.Method: r.method,
						"attempt":                          attempt,
					})
			}
		}

		resp, err = r.client.Do(req)

		// If success
		if err == nil && resp.StatusCode < 500 {
			requestDuration := time.Since(startTime)

			if r.logger != nil {
				r.logger.InfoWithCategory(
					sfhttprequest.Category.API.HTTP,
					sfhttprequest.SubCategory.API.Response,
					"HTTP request successful",
					map[string]interface{}{
						sfhttprequest.ExtraKey.HTTP.URL:             r.url,
						sfhttprequest.ExtraKey.HTTP.Method:          r.method,
						sfhttprequest.ExtraKey.HTTP.StatusCode:      resp.StatusCode,
						sfhttprequest.ExtraKey.Performance.Duration: requestDuration.String(),
					})
			}

			return WrapResponse(resp), nil
		}

		// Log the failure
		if r.logger != nil {
			errMsg := "unknown error"
			if err != nil {
				errMsg = err.Error()
			} else if resp != nil {
				errMsg = fmt.Sprintf("status code %d", resp.StatusCode)
			}

			r.logger.ErrorWithCategory(
				sfhttprequest.Category.API.HTTP,
				sfhttprequest.SubCategory.Status.Error,
				"HTTP request failed",
				map[string]interface{}{
					sfhttprequest.ExtraKey.HTTP.URL:           r.url,
					sfhttprequest.ExtraKey.HTTP.Method:        r.method,
					sfhttprequest.ExtraKey.Error.ErrorMessage: errMsg,
					"attempt": attempt,
				})
		}

		// If this was the last attempt, trigger fallback
		if attempt == r.retryCount && r.fallbackFunc != nil {
			r.fallbackFunc()

			if r.logger != nil {
				r.logger.InfoWithCategory(
					sfhttprequest.Category.API.HTTP,
					sfhttprequest.SubCategory.API.Request,
					"Using fallback for HTTP request",
					map[string]interface{}{
						sfhttprequest.ExtraKey.HTTP.URL:    r.url,
						sfhttprequest.ExtraKey.HTTP.Method: r.method,
					})
			}
		}
	}

	return WrapResponse(resp), err
}

// buildRequest creates the underlying http.Request
func (r *Request) buildRequest() (*http.Request, error) {
	// Build final URL with query parameters
	finalURL := r.url
	if len(r.queryParams) > 0 {
		if strings.Contains(finalURL, "?") {
			finalURL += "&" + r.queryParams.Encode()
		} else {
			finalURL += "?" + r.queryParams.Encode()
		}
	}

	var bodyReader io.Reader

	// Handle form data if present
	if len(r.formData) > 0 {
		bodyReader = strings.NewReader(r.formData.Encode())
	} else if r.body != nil {
		bodyReader = bytes.NewReader(r.body)
	}

	req, err := http.NewRequest(r.method, finalURL, bodyReader)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range r.headers {
		req.Header.Add(key, value)
	}

	r.request = req
	return req, nil
}

// URL creates a new HTTP request with direct URL (not using the connection registry)
func URL(method, urlStr string) *Request {
	// Create a default client and wrap it with APM tracing
	client := apmhttp.WrapClient(http.DefaultClient)

	r := &Request{
		client:      client,
		method:      method,
		url:         urlStr,
		headers:     make(map[string]string),
		queryParams: url.Values{},
		formData:    url.Values{},
		ctx:         context.Background(),
	}

	// Add global headers to the request
	globalHeaders := GetGlobalHeaders()
	for key, value := range globalHeaders {
		r.headers[key] = value
	}

	return r
}
