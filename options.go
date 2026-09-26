package paystack

import (
	"net/http"
	"strings"
	"time"
)

// Option configures a Client at construction time in NewClient.
type Option func(*Client)

// WithHTTPClient sets the underlying http.Client used for all requests,
// allowing callers to inject custom transports, proxies, connection
// pooling, or DNS behavior. The SDK copies the supplied client so that
// later options (such as WithTimeout) never mutate the caller's original
// http.Client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		if hc == nil {
			return
		}
		copied := *hc
		c.httpClient = &copied
	}
}

// WithBaseURL overrides the API base URL, for example to target a sandbox
// environment or a test double in CI.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		if baseURL != "" {
			c.baseURL = strings.TrimRight(baseURL, "/")
		}
	}
}

// WithTimeout sets the per-request timeout applied by the underlying
// http.Client. It does not override a shorter deadline already present on
// a request's context.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if d > 0 {
			c.httpClient.Timeout = d
		}
	}
}

// WithRetryPolicy overrides the default retry policy used for transient
// failures.
func WithRetryPolicy(policy RetryPolicy) Option {
	return func(c *Client) {
		c.retryPolicy = policy
	}
}

// WithObserver registers a TransportObserver to receive request lifecycle
// events for integration with logging, metrics, or tracing systems.
func WithObserver(observer TransportObserver) Option {
	return func(c *Client) {
		if observer != nil {
			c.observer = observer
		}
	}
}
