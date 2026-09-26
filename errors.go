package paystack

import (
	"context"
	"errors"
	"fmt"
)

// Sentinel errors that callers can check with errors.Is against any error
// returned by this SDK.
var (
	// ErrMissingSecretKey is returned when NewClient is called without a
	// secret key.
	ErrMissingSecretKey = errors.New("paystack: secret key must not be empty")

	// ErrUnauthorized indicates the API rejected the request's credentials.
	ErrUnauthorized = errors.New("paystack: unauthorized")

	// ErrRateLimited indicates the API responded with HTTP 429.
	ErrRateLimited = errors.New("paystack: rate limited")

	// ErrTimeout indicates the request did not complete before its
	// deadline or configured timeout elapsed.
	ErrTimeout = errors.New("paystack: request timed out")

	// ErrNotFound indicates the API responded with HTTP 404.
	ErrNotFound = errors.New("paystack: resource not found")
)

// APIError represents an error response returned by the Paystack API, or a
// transport-level failure encountered while trying to reach it. It carries
// enough structured information for callers to make automated retry and
// alerting decisions without parsing error strings.
type APIError struct {
	// StatusCode is the HTTP status code returned by Paystack. It is zero
	// for errors that never received a response, such as network failures.
	StatusCode int

	// Message is the human-readable error message from the API response,
	// or a description of the underlying transport failure.
	Message string

	// Endpoint is the request path that produced the error.
	Endpoint string

	// RequestID is the Paystack-provided request identifier, when present.
	RequestID string

	// Retryable reports whether the SDK considers this error safe to retry.
	Retryable bool

	// Err is the underlying error, if any (for example a network or JSON
	// decode error). It is never a value that could carry credentials.
	Err error
}

func (e *APIError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("paystack: %s (status %d, endpoint %s)", e.Message, e.StatusCode, e.Endpoint)
	}
	return fmt.Sprintf("paystack: %s (endpoint %s)", e.Message, e.Endpoint)
}

// Unwrap exposes the underlying error so that errors.Is / errors.As can
// traverse into transport or decode failures.
func (e *APIError) Unwrap() error {
	return e.Err
}

// Is implements the errors.Is contract for the package sentinel errors,
// classifying by status code so callers can write:
//
//	if errors.Is(err, paystack.ErrRateLimited) { ... }
func (e *APIError) Is(target error) bool {
	switch target {
	case ErrUnauthorized:
		return e.StatusCode == 401
	case ErrRateLimited:
		return e.StatusCode == 429
	case ErrNotFound:
		return e.StatusCode == 404
	case ErrTimeout:
		return errors.Is(e.Err, context.DeadlineExceeded)
	default:
		return false
	}
}
