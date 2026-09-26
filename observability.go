package paystack

import "time"

// RequestInfo describes an outbound request for observability hooks. Header
// values are always pre-redacted, so implementations can log or export
// RequestInfo without risking credential leakage.
type RequestInfo struct {
	Method   string
	Endpoint string
	Headers  map[string][]string
	Attempt  int
}

// ResponseInfo describes a completed request/response cycle for
// observability hooks.
type ResponseInfo struct {
	RequestInfo
	StatusCode int
	Duration   time.Duration
	RequestID  string
}

// RetryInfo describes a retry decision for observability hooks.
type RetryInfo struct {
	RequestInfo
	Delay time.Duration
	Cause error
}

// TransportObserver receives lifecycle events for every request made by a
// Client, so applications can integrate with OpenTelemetry, Prometheus,
// Datadog, Grafana, or plain structured logging without the SDK depending
// on any particular provider.
//
// Implementations must return quickly: hooks are invoked synchronously on
// the request path, and a slow observer directly slows down every API
// call.
type TransportObserver interface {
	// RequestStarted is invoked immediately before a request is sent.
	RequestStarted(RequestInfo)

	// RequestCompleted is invoked after a response is received, whether or
	// not it represents an API-level error.
	RequestCompleted(ResponseInfo)

	// RequestFailed is invoked when a request could not be completed at
	// all, such as on a network error or context cancellation.
	RequestFailed(RequestInfo, error)

	// Retry is invoked when the client schedules a retry attempt, before
	// the backoff delay is slept.
	Retry(RetryInfo)
}

// noopObserver is the default TransportObserver, used when the caller does
// not supply one via WithObserver.
type noopObserver struct{}

func (noopObserver) RequestStarted(RequestInfo)       {}
func (noopObserver) RequestCompleted(ResponseInfo)    {}
func (noopObserver) RequestFailed(RequestInfo, error) {}
func (noopObserver) Retry(RetryInfo)                  {}
