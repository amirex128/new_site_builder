package sfconfigmanager

import (
	"fmt"
	"math"
	"time"
)

// RetryOptions defines options for retrying operations
type RetryOptions struct {
	// MaxRetries is the maximum number of retry attempts
	MaxRetries int
	// InitialBackoff is the initial backoff duration
	InitialBackoff time.Duration
	// MaxBackoff is the maximum backoff duration
	MaxBackoff time.Duration
	// BackoffFactor is the multiplier for each subsequent backoff
	BackoffFactor float64
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

// WithRetry executes the provided function with retry logic
func WithRetry(operation func() error, options *RetryOptions, log Logger, operationName string) error {
	if options == nil {
		options = DefaultRetryOptions()
	}

	var lastErr error
	backoff := options.InitialBackoff

	for attempt := 0; attempt <= options.MaxRetries; attempt++ {
		// First attempt is not a retry
		if attempt > 0 && log != nil {
			log.InfoWithCategory(
				Category.System.General,
				SubCategory.Status.Retry,
				fmt.Sprintf("Retrying %s (attempt %d/%d)", operationName, attempt, options.MaxRetries),
				map[string]interface{}{
					"backoff": backoff.String(),
				},
			)

			// Wait before retry
			time.Sleep(backoff)

			// Calculate next backoff with exponential increase
			backoff = time.Duration(math.Min(
				float64(options.MaxBackoff),
				float64(backoff)*options.BackoffFactor,
			))
		}

		// Run the operation
		err := operation()
		if err == nil {
			// Success
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

		lastErr = err
		if log != nil {
			log.WarnWithCategory(
				Category.System.General,
				SubCategory.Status.Failure,
				fmt.Sprintf("%s failed (attempt %d/%d)", operationName, attempt+1, options.MaxRetries),
				map[string]interface{}{
					"error": err.Error(),
				},
			)
		}
	}

	// All retries failed
	if log != nil {
		log.ErrorWithCategory(
			Category.System.General,
			SubCategory.Status.Failure,
			fmt.Sprintf("%s failed after %d attempts", operationName, options.MaxRetries+1),
			map[string]interface{}{
				"error": lastErr.Error(),
			},
		)
	}

	return fmt.Errorf("operation %s failed after %d attempts: %w", operationName, options.MaxRetries+1, lastErr)
}
