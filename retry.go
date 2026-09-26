package paystack

import (
	"math"
	"math/rand"
	"net/http"
	"time"
)

// RetryPolicy controls how the client retries failed requests.
type RetryPolicy struct {
	// MaxRetries is the maximum number of retry attempts after the initial
	// request. A value of 0 disables retries entirely.
	MaxRetries int

	// Backoff computes the delay before each retry attempt.
	Backoff BackoffStrategy
}

// DefaultRetryPolicy returns the retry policy used when none is supplied
// to NewClient: up to two retries with exponential backoff and jitter.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxRetries: 2,
		Backoff:    ExponentialBackoff{Base: 250 * time.Millisecond, Max: 5 * time.Second},
	}
}

// BackoffStrategy computes the delay before a retry attempt. attempt is
// 1-indexed: the first retry is attempt 1, the second is attempt 2, and so
// on.
type BackoffStrategy interface {
	Delay(attempt int) time.Duration
}

// ExponentialBackoff computes delays that grow exponentially with the
// attempt number, capped at Max, with jitter applied so that many
// concurrent clients retrying at once do not collide in a retry storm.
type ExponentialBackoff struct {
	// Base is the delay used for the first retry attempt.
	Base time.Duration

	// Max is the upper bound on any computed delay.
	Max time.Duration
}

// Delay implements BackoffStrategy.
func (b ExponentialBackoff) Delay(attempt int) time.Duration {
	base := b.Base
	if base <= 0 {
		base = 250 * time.Millisecond
	}
	maxDelay := b.Max
	if maxDelay <= 0 {
		maxDelay = 5 * time.Second
	}
	if attempt < 1 {
		attempt = 1
	}

	raw := float64(base) * math.Pow(2, float64(attempt-1))
	if raw > float64(maxDelay) {
		raw = float64(maxDelay)
	}

	// Full jitter: a random delay between half of the computed backoff and
	// the full computed backoff.
	jittered := raw/2 + rand.Float64()*(raw/2)
	return time.Duration(jittered)
}

// isRetryableStatus reports whether an HTTP status code represents a
// transient failure that is generally safe to retry.
func isRetryableStatus(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}
