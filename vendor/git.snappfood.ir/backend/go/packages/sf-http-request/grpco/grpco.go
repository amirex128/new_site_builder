package grpco

import (
	"context"
	"crypto/tls"
	"errors"
	"time"

	sfhttprequest "git.snappfood.ir/backend/go/packages/sf-http-request"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// Error constants
var (
	// ErrCircuitOpen is left for backward compatibility
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

// Request represents a gRPC request with additional features like timeout, retry, and fallback
type Request struct {
	conn             *grpc.ClientConn
	methodName       string
	request          interface{}
	timeoutMs        int
	retryCount       int
	retryDelayMs     int
	fallbackFunc     func()
	ctx              context.Context
	md               metadata.MD
	callOpts         []grpc.CallOption
	responseHeaders  metadata.MD
	responseTrailers metadata.MD
	logger           sfhttprequest.Logger
	err              error // Store errors during initialization to return during Send()
}

// ClientOption is a function to customize the gRPC client connection
type ClientOption func(*grpc.DialOption)

// WithTLS enables TLS for the gRPC connection with customizable TLS config
func WithTLS(config *tls.Config) ClientOption {
	return func(opt *grpc.DialOption) {
		*opt = grpc.WithTransportCredentials(credentials.NewTLS(config))
	}
}

// WithInsecure disables TLS for the gRPC connection
func WithInsecure() ClientOption {
	return func(opt *grpc.DialOption) {
		*opt = grpc.WithTransportCredentials(insecure.NewCredentials())
	}
}

// Method sets the method name for the gRPC request
func (r *Request) Method(method string) *Request {
	r.methodName = method
	return r
}

// Request sets the request payload
func (r *Request) Request(req interface{}) *Request {
	r.request = req
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

// SetContext sets a custom context for the request
func (r *Request) SetContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

// Header adds a metadata header to the request
func (r *Request) Header(key, value string) *Request {
	if r.md == nil {
		r.md = metadata.MD{}
	}
	r.md.Append(key, value)
	return r
}

// SetHeaders adds multiple metadata headers to the request
func (r *Request) SetHeaders(headers map[string]string) *Request {
	if r.md == nil {
		r.md = metadata.MD{}
	}
	for key, value := range headers {
		r.md.Append(key, value)
	}
	return r
}

// GetHeaders returns all request headers
func (r *Request) GetHeaders() metadata.MD {
	return r.md
}

// GetResponseHeaders returns the headers received in the response
func (r *Request) GetResponseHeaders() metadata.MD {
	return r.responseHeaders
}

// GetResponseTrailers returns the trailers received in the response
func (r *Request) GetResponseTrailers() metadata.MD {
	return r.responseTrailers
}

// GetResponseHeader returns a specific response header value
func (r *Request) GetResponseHeader(key string) []string {
	if r.responseHeaders == nil {
		return nil
	}
	return r.responseHeaders.Get(key)
}

// GetResponseTrailer returns a specific response trailer value
func (r *Request) GetResponseTrailer(key string) []string {
	if r.responseTrailers == nil {
		return nil
	}
	return r.responseTrailers.Get(key)
}

// Send is a generic method to execute a gRPC call and return the response
// It allows for a more fluent interface similar to httpo
func (r *Request) Send(response interface{}) error {
	// Check for initialization errors
	if r.err != nil {
		// Execute fallback if available and there's an initialization error
		if r.fallbackFunc != nil {
			r.fallbackFunc()
		}
		return r.err
	}

	if r.methodName == "" {
		if r.fallbackFunc != nil {
			r.fallbackFunc()
		}
		return errors.New("method name is required")
	}

	if r.request == nil {
		if r.fallbackFunc != nil {
			r.fallbackFunc()
		}
		return errors.New("request payload is required")
	}

	// Check if the request implements proto.Message
	_, isProtoMessage := r.request.(proto.Message)
	if !isProtoMessage {
		if r.fallbackFunc != nil {
			r.fallbackFunc()
		}
		return errors.New("request must be a pointer to a struct that implements proto.Message")
	}

	// Create a context with metadata
	ctx := metadata.NewOutgoingContext(r.ctx, r.md)

	// Apply timeout
	var cancel context.CancelFunc
	if r.timeoutMs > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(r.timeoutMs)*time.Millisecond)
		defer cancel()
	}

	// Prepare header and trailer metadata
	var headerMD, trailerMD metadata.MD
	callOpts := append(r.callOpts,
		grpc.Header(&headerMD),
		grpc.Trailer(&trailerMD),
	)

	// Generic invoker that uses the method name to invoke the correct method
	invoker := func(ctx context.Context, conn *grpc.ClientConn, req interface{}, res interface{}, opts ...grpc.CallOption) error {
		return conn.Invoke(ctx, r.methodName, req, res, opts...)
	}

	// Retry logic
	var err error
	for attempt := 0; attempt <= r.retryCount; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(r.retryDelayMs) * time.Millisecond)
		}

		err = invoker(ctx, r.conn, r.request, response, callOpts...)

		// Log the request if logger is available
		if r.logger != nil {
			logExtra := map[string]interface{}{
				sfhttprequest.ExtraKey.Service.ServiceName: r.methodName,
				"attempt":    attempt + 1,
				"maxRetries": r.retryCount,
			}

			if err != nil {
				logExtra[sfhttprequest.ExtraKey.Error.ErrorMessage] = err.Error()
				r.logger.ErrorWithCategory(
					sfhttprequest.Category.API.GRPC,
					sfhttprequest.SubCategory.Status.Error,
					"gRPC request failed",
					logExtra,
				)
			} else {
				r.logger.InfoWithCategory(
					sfhttprequest.Category.API.GRPC,
					sfhttprequest.SubCategory.Status.Success,
					"gRPC request succeeded",
					logExtra,
				)
			}
		}

		// Store response headers and trailers
		r.responseHeaders = headerMD
		r.responseTrailers = trailerMD

		// If success
		if err == nil {
			return nil
		}

		// If this was the last attempt and we still have an error
		if attempt == r.retryCount {
			break
		}
	}

	// Execute fallback if available and all retries have failed
	if err != nil && r.fallbackFunc != nil {
		r.fallbackFunc()
	}

	return err
}

// RequestWithResponse is a helper method that sets the request and returns the Request
// to allow for fluent call chaining
func (r *Request) RequestWithResponse(req interface{}, resp interface{}) *RequestWithResponse {
	r.request = req
	return &RequestWithResponse{
		Request:  r,
		Response: resp,
	}
}

// RequestWithResponse provides a fluent interface for requests with predefined response types
type RequestWithResponse struct {
	*Request
	Response interface{}
}

// Send executes the request and populates the response
func (r *RequestWithResponse) Send() error {
	return r.Request.Send(r.Response)
}
