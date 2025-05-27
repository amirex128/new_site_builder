package v8

import (
	"fmt"
	"math"
	"time"
)

// RetryOptions defines options for retrying operations
type RetryOptions struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	BackoffFactor  float64
}

// DefaultRetryOptions returns the default retry options
func DefaultRetryOptions() *RetryOptions {
	return &RetryOptions{
		MaxRetries:     3,
		InitialBackoff: 500 * time.Millisecond,
		MaxBackoff:     10 * time.Second,
		BackoffFactor:  2.0,
	}
}

// WithRetry executes the provided function with retry logic and panic recovery
func WithRetry(operation func() error, options *RetryOptions, log Logger, operationName string) (err error) {
	if options == nil {
		options = DefaultRetryOptions()
	}

	var lastErr error
	backoff := options.InitialBackoff

	for attempt := 0; attempt <= options.MaxRetries; attempt++ {
		// Panic recovery for each attempt
		func() {
			defer func() {
				if r := recover(); r != nil {
					switch x := r.(type) {
					case string:
						lastErr = fmt.Errorf("panic: %s", x)
					case error:
						lastErr = fmt.Errorf("panic: %w", x)
					default:
						lastErr = fmt.Errorf("panic: %v", x)
					}
					if log != nil {
						log.ErrorWithCategory(
							Category.System.General,
							SubCategory.Status.Error,
							fmt.Sprintf("Recovered from panic in %s", operationName),
							map[string]interface{}{"error": lastErr.Error()},
						)
					}
				}
			}()
			lastErr = operation()
		}()

		if lastErr == nil {
			if attempt > 0 && log != nil {
				log.InfoWithCategory(
					Category.System.General,
					SubCategory.Status.Success,
					fmt.Sprintf("%s succeeded after %d attempts", operationName, attempt+1),
					nil,
				)
			}
			return nil
		}

		if log != nil {
			log.WarnWithCategory(
				Category.System.General,
				SubCategory.Status.Failure,
				fmt.Sprintf("%s failed (attempt %d/%d)", operationName, attempt+1, options.MaxRetries),
				map[string]interface{}{"error": lastErr.Error()},
			)
		}

		if attempt < options.MaxRetries {
			if log != nil {
				log.InfoWithCategory(
					Category.System.General,
					SubCategory.Status.Retry,
					fmt.Sprintf("Retrying %s (attempt %d/%d)", operationName, attempt+1, options.MaxRetries),
					map[string]interface{}{"backoff": backoff.String()},
				)
			}
			time.Sleep(backoff)
			backoff = time.Duration(math.Min(float64(options.MaxBackoff), float64(backoff)*options.BackoffFactor))
		}
	}

	if log != nil {
		log.ErrorWithCategory(
			Category.System.General,
			SubCategory.Status.Failure,
			fmt.Sprintf("%s failed after %d attempts", operationName, options.MaxRetries+1),
			map[string]interface{}{"error": lastErr.Error()},
		)
	}

	return fmt.Errorf("operation %s failed after %d attempts: %w", operationName, options.MaxRetries+1, lastErr)
}
